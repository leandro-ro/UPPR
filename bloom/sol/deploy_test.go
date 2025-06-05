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
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

// Uncomment the following if making changes to cascadingBloomFilter.sol
/*func TestCompileAndGenBindings(t *testing.T) {
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
*/

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

// BenchmarkUpdateCascade benchmarks the gas consumption during the update process of an on-chain Bloom filter cascade.
func BenchmarkUpdateCascade(b *testing.B) {
	configs := []struct {
		name     string
		domain   int
		capacity int
	}{
		{"D1k_C10", 1_000, 10},
		{"D1k_C50", 1_000, 50},
		{"D1k_C100", 1_000, 100},
		{"D10k_C100", 10_000, 100},
		{"D10k_C500", 10_000, 500},
		{"D10k_C1k", 10_000, 1_000},
		{"D100k_C1k", 100_000, 1_000},
		{"D100k_C5k", 100_000, 5_000},
		{"D100k_C10k", 100_000, 10_000},
		{"D1M_C10k", 1_000_000, 10_000},
		{"D1M_C50k", 1_000_000, 50_000},
		{"D1M_C100k", 1_000_000, 100_000},
	}

	fmt.Println("Benchmark Gas Consumption of UpdateCascade (1st: initial, 2nd: update) for 1 Gas = 1 Gwei and N=10:")
	fmt.Println("| Domain   | Capacity | 1st Avg Gas | 1st ETH     | 2nd Avg Gas | 2nd ETH     |")
	fmt.Println("|----------|----------|-------------|-------------|-------------|-------------|")

	for _, cfg := range configs {
		var sumGas1, sumGas2 uint64

		for i := 0; i < 10; i++ {
			gas1, gas2, err := runUpdateCascadeBenchmark(cfg.domain, cfg.capacity)
			if err != nil {
				b.Fatalf("Benchmark failed for domain=%d, capacity=%d: %v", cfg.domain, cfg.capacity, err)
			}
			sumGas1 += gas1
			sumGas2 += gas2
		}

		avgGas1 := sumGas1 / 10
		avgGas2 := sumGas2 / 10
		eth1 := float64(avgGas1) * 1e-9 // Assume gas price of 1 Gwei
		eth2 := float64(avgGas2) * 1e-9

		fmt.Printf("| %8d | %8d | %11d | %.9f | %11d | %.9f |\n",
			cfg.domain, cfg.capacity,
			avgGas1, eth1,
			avgGas2, eth2,
		)
	}
}

func runUpdateCascadeBenchmark(domain, capacity int) (uint64, uint64, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return 0, 0, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	if err != nil {
		return 0, 0, err
	}

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 30_000_000_000)

	_, tx, contract, err := onchain.DeployBloom(auth, sim)
	if err != nil {
		return 0, 0, err
	}
	sim.Commit()

	_, err = sim.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return 0, 0, err
	}

	// First update
	cascade1 := bloom.NewCascade(domain, capacity)
	valid1, revoked1 := genRevocationTokens(domain, capacity)
	if err := cascade1.Update(revoked1, valid1); err != nil {
		return 0, 0, err
	}
	onChainFilter1, numHf1, bitlen1 := cascade1.GetOnChainFilter()

	tx1, err := contract.UpdateCascade(auth, onChainFilter1, numHf1, bitlen1)
	if err != nil {
		return 0, 0, err
	}
	sim.Commit()
	receipt1, err := sim.TransactionReceipt(context.Background(), tx1.Hash())
	if err != nil {
		return 0, 0, err
	}

	// Wait and second update
	time.Sleep(1 * time.Second)
	cascade2 := bloom.NewCascade(domain, capacity)
	valid2, revoked2 := genRevocationTokens(domain, capacity)
	if err := cascade2.Update(revoked2, valid2); err != nil {
		return 0, 0, err
	}
	onChainFilter2, numHf2, bitlen2 := cascade2.GetOnChainFilter()

	tx2, err := contract.UpdateCascade(auth, onChainFilter2, numHf2, bitlen2)
	if err != nil {
		return 0, 0, err
	}
	sim.Commit()
	receipt2, err := sim.TransactionReceipt(context.Background(), tx2.Hash())
	if err != nil {
		return 0, 0, err
	}

	return receipt1.GasUsed, receipt2.GasUsed, nil
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
