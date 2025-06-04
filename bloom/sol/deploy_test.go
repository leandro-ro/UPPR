package bloom

import (
	"PrivacyPreservingRevocationCode/bloom"
	onchain "PrivacyPreservingRevocationCode/bloom/sol/build"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompileAndGenBindings(t *testing.T) {
	// 1) Remove any existing build directory for a clean slate
	buildDir := "build"
	_ = os.RemoveAll(buildDir)

	// 2) Create build directory
	err := os.MkdirAll(buildDir, 0755)
	require.NoError(t, err, "failed to create build directory")

	// 3) Compile the Solidity contract with solc (via-ir, Istanbul EVM)
	solFile := "cascadingBloomFilter.sol"
	solcCmd := exec.Command(
		"solc",
		"--bin",
		"--abi",
		"--overwrite",
		"--evm-version", "istanbul",
		"--via-ir",
		solFile,
		"-o", buildDir,
	)
	out, err := solcCmd.CombinedOutput()
	require.NoErrorf(t, err, "solc failed: %v\n%s", err, out)

	// 4) Derive contract name from filename
	filename := filepath.Base(solFile)                                   // e.g. "CascadingBloomFilter.sol"
	contractName := strings.TrimSuffix(filename, filepath.Ext(filename)) // "CascadingBloomFilter"

	// 5) Check that .bin and .abi exist
	binPath := filepath.Join(buildDir, contractName+".bin")
	abiPath := filepath.Join(buildDir, contractName+".abi")
	require.FileExists(t, binPath, "Expected binary file at %s", binPath)
	require.FileExists(t, abiPath, "Expected ABI file at %s", abiPath)

	// 6) Run abigen to generate Go bindings in buildDir
	bindingFilename := contractName + "_binding.go"
	bindingPath := filepath.Join(buildDir, bindingFilename)

	abigenCmd := exec.Command(
		"abigen",
		"--abi="+abiPath,
		"--bin="+binPath,
		"--pkg=bloom",
		"--out="+bindingPath,
	)
	out, err = abigenCmd.CombinedOutput()
	require.NoErrorf(t, err, "abigen failed: %v\n%s", err, out)

	// 7) Verify that the binding file was created
	require.FileExists(t, bindingPath, "Expected Go binding file at %s", bindingPath)

	// 8) Sanity checks: ensure all generated files live in buildDir
	require.Equal(t, buildDir, filepath.Dir(binPath), "binary not in build directory")
	require.Equal(t, buildDir, filepath.Dir(abiPath), "ABI not in build directory")
	require.Equal(t, buildDir, filepath.Dir(bindingPath), "binding not in build directory")

	fmt.Printf("Generated files:\n- %s\n- %s\n- %s\n", binPath, abiPath, bindingPath)
}

func TestDeploy(t *testing.T) {
	// 1) Set up a simulated backend and a funded transactor
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)}, // 1 ETH
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000)

	// 2) Deploy the contract via the generated binding (DeployBloom)
	address, tx, contract, err := onchain.DeployBloom(auth, sim)
	require.NoError(t, err, "deployment failed")
	sim.Commit()

	receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status, "deployment reverted")

	t.Logf("Contract deployed at %s, gas used: %d", address.Hex(), receipt.GasUsed)

	// 3) Immediately after deploy, layerCount should be 0
	lc, err := contract.LayerCount(&bind.CallOpts{})
	require.NoError(t, err)
	require.True(t, lc.Cmp(big.NewInt(0)) == 0, "expected layerCount to be 0")
}

func TestOnChainFilterSerialization(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000)

	_, _, contract, err := onchain.DeployBloom(auth, sim)
	require.NoError(t, err)
	sim.Commit()

	domain := 10_000
	capacity := 100

	cascade := bloom.NewCascade(domain, capacity)
	valid, revoked := genRevocationTokens(domain, capacity)
	err = cascade.Update(revoked, valid)
	require.NoError(t, err)

	// Off-chain → On-chain
	onChainFilters, numHf, bitLens := cascade.GetOnChainFilter()
	_, err = contract.UpdateCascade(auth, onChainFilters, numHf, bitLens)
	require.NoError(t, err)
	sim.Commit()

	for i := range onChainFilters {
		expectedFilter := onChainFilters[i]
		expectedBitLen := bitLens[i].Int64()
		expectedK := numHf[i].Int64()

		// Get on-chain metadata + filter
		debug, err := contract.GetLayerMetadata(&bind.CallOpts{}, big.NewInt(int64(i)))
		require.NoError(t, err)

		require.Equal(t, uint64(expectedBitLen), debug.FilterSizeBits.Uint64(), "layer %d: bitLen mismatch", i)
		require.Equal(t, uint64(expectedK), debug.K.Uint64(), "layer %d: k mismatch", i)
		require.Equal(t, expectedFilter, debug.Filter, "layer %d: filter bytes mismatch", i)
	}
}

/*func TestTokenHash(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)}, // 1 ETH
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000_000)

	_, tx, contract, err := onchain.DeployBloom(auth, sim)
	require.NoError(t, err, "deployment failed")
	sim.Commit()

	receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status, "deployment reverted")

	domain := 1000
	capacity := 100

	cascade := bloom.NewCascade(domain, capacity)
	valid, revoked := genRevocationTokens(domain, capacity)
	err = cascade.Update(revoked, valid)
	require.NoError(t, err, "failed to update cascade with test data")

	onChainFilter, numHf, bitlen := cascade.GetOnChainFilter()
	tx, err = contract.UpdateCascade(auth, onChainFilter, numHf, bitlen)
	require.NoError(t, err, "failed to update on-chain filter")
	sim.Commit()

	receipt, err = sim.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status, "update reverted")

	// Select one token to test (e.g., first rev token)
	tok := revoked[0]
	offChainResult, offChainLayer := cascade.Test(tok)

	k := cascade.GetFilters()[0].K()
	m := cascade.GetFilters()[0].BitLen()
	t.Logf("Mod expected to be %v", m)
	debug, err := contract.DebugToken(&bind.CallOpts{}, tok, big.NewInt(int64(k)), big.NewInt(int64(m)))
	require.NoError(t, err, "failed to get debug hash values and locations")

	t.Logf("Off-chain test result: %v (stopped at layer %d)", offChainResult, offChainLayer)
	t.Logf("On-chain extracted hashes: %v", debug.Hashes)
	t.Logf("On-chain raw locations: %v", debug.RawLocations)
	t.Logf("On-chain mod locations: %v", debug.ModLocations)

	debug2, err := contract.DebugTestToken(&bind.CallOpts{}, tok)
	require.NoError(t, err, "failed to get debug hash values and locations")

	t.Logf("On-chain: Layer %v we got filter size %v", debug2.LayerIndex, debug2.FilterSizeBits)

	debug5, err := contract.GetLayerMetadata(&bind.CallOpts{}, big.NewInt(0))
	t.Logf("Layer metadata: %v", debug5)

	hexpected := bloom.BaseHashesDebug(tok)
	t.Logf("Expected hashes: %v", hexpected)

	for i := uint(0); i < k; i++ {
		locrawexpected := cascade.GetFilters()[0].LocationRawDebug(hexpected, i)
		t.Logf("Expected raw loc %d: %v", i, locrawexpected)
		locexpected := cascade.GetFilters()[0].LocationDebug(hexpected, i)
		t.Logf("Expected mod loc %d: %v", i, locexpected)
	}

	res, layer, err := contract.TestToken(&bind.CallOpts{}, tok)
	require.NoError(t, err, "failed to test token on-chain")
	t.Logf("On-chain test result: %v (stopped at layer %d)", res, layer)
}*/

func TestUpdate(t *testing.T) {
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)}, // 1 ETH
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000_000)

	// 2) Deploy the contract via the generated binding (DeployBloom)
	_, tx, contract, err := onchain.DeployBloom(auth, sim)
	require.NoError(t, err, "deployment failed")
	sim.Commit()

	receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status, "deployment reverted")

	domain := 1000
	capacity := 100

	cascade := bloom.NewCascade(domain, capacity)
	valid, revoked := genRevocationTokens(domain, capacity)
	err = cascade.Update(revoked, valid)
	require.NoError(t, err, "failed to update cascade with test data")

	onChainFilter, numHf, bitlen := cascade.GetOnChainFilter()
	require.NotNil(t, onChainFilter, "expected on-chain filter to be non-nil")
	require.NotNil(t, numHf, "expected number of hash functions to be non-nil")

	tx, err = contract.UpdateCascade(auth, onChainFilter, numHf, bitlen)
	require.NoError(t, err, "failed to update on-chain filter")
	sim.Commit()

	receipt, err = sim.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status, "update reverted")
	// test all non-revocations
	for i, tok := range valid {
		b, lexpected := cascade.Test(tok)
		require.False(t, b, "non-revoked token detected as revoked")

		res, layer, err := contract.TestToken(&bind.CallOpts{}, tok)
		require.Equal(t, lexpected, int(layer.Int64()), "expected layer %d, got %d at index %d", lexpected, int(layer.Int64()), i)
		require.NoError(t, err, "failed to test token")
		require.False(t, res, "non-revoked token detected as revoked")
	}
	// test all revocations
	for i, tok := range revoked {
		b, lexpected := cascade.Test(tok)
		require.True(t, b, "revoked token not detected as revoked. Got to index %d", i)

		res, layer, err := contract.TestToken(&bind.CallOpts{}, tok)
		require.NoError(t, err, "failed to test token")
		require.True(t, res, "revoked token not detected as revoked. Got to index %d", i)
		require.Equal(t, lexpected, int(layer.Int64()), "expected layer %d, got %d at index %d", lexpected, int(layer.Int64()), i)
	}

}

// genRevocationTokens generates random 128-bit tokens and splits them into valid and revoked sets.
func genRevocationTokens(domain, revocCapacity int) (valid, revoked [][]byte) {
	valid = generateRandom128BitSlices(domain - revocCapacity)
	revoked = generateRandom128BitSlices(revocCapacity)
	return
}

// generateRandom128BitSlices returns 'count' number of random 128-bit (16-byte) slices.
func generateRandom128BitSlices(count int) [][]byte {
	result := make([][]byte, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			panic(fmt.Sprintf("failed to generate random bytes: %v", err))
		}
		result[i] = b
	}
	return result
}
