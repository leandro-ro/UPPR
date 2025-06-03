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

func TestUpdate(t *testing.T) {
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

	domain := 10_000
	capacity := 100

	cascade := bloom.NewCascade(domain, capacity)
	valid, revoked := genRevocationTokens(domain, capacity)
	err = cascade.Update(revoked, valid)
	require.NoError(t, err, "failed to update cascade with test data")

	contract.UpdateCascade(auth)

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
