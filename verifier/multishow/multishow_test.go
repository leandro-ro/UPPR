package multishow

import (
	onchainBloom "PrivacyPreservingRevocationCode/bloom/sol/build"
	"PrivacyPreservingRevocationCode/holder"
	"PrivacyPreservingRevocationCode/issuer"
	onchainVerifier "PrivacyPreservingRevocationCode/verifier/multishow/build"
	zkpContract "PrivacyPreservingRevocationCode/zkp/sol/build"
	"context"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

/*
func TestMultiShow_CompileAndGenBindings(t *testing.T) {
	buildDir := "build"
	solFile := "MultiShowVerifier.sol"

	// Clean and create build directory
	_ = os.RemoveAll(buildDir)
	require.NoError(t, os.MkdirAll(buildDir, 0755), "failed to create build dir")

	// Compile Solidity file with solc
	cmd := exec.Command(
		"solc",
		"--bin",
		"--abi",
		"--overwrite",
		"--evm-version", "istanbul",
		"--via-ir",
		"--base-path", "../../", // point to project root
		solFile,
		"-o", buildDir,
	)
	out, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "solc failed: %v\n%s", err, string(out))

	// Extract contract name
	contractName := strings.TrimSuffix(filepath.Base(solFile), filepath.Ext(solFile))

	// Paths to output files
	binPath := filepath.Join(buildDir, contractName+".bin")
	abiPath := filepath.Join(buildDir, contractName+".abi")
	require.FileExists(t, binPath, "missing bin file")
	require.FileExists(t, abiPath, "missing abi file")

	// Generate Go bindings
	bindingPath := filepath.Join(buildDir, strings.ToLower(contractName)+"_binding.go")
	abigenCmd := exec.Command(
		"abigen",
		"--abi="+abiPath,
		"--bin="+binPath,
		"--pkg=verifier",
		"--out="+bindingPath,
	)
	abigenOut, err := abigenCmd.CombinedOutput()
	require.NoErrorf(t, err, "abigen failed: %v\n%s", err, string(abigenOut))
	require.FileExists(t, bindingPath, "binding file not created")

	fmt.Printf("Generated bindings:\n- %s\n- %s\n- %s\n", binPath, abiPath, bindingPath)
}
*/

func TestMultiShow_EndToEnd(t *testing.T) {
	testIssuer := issuer.NewIssuer(issuer.MultiShow)
	issuerPubKey := eddsa.PublicKey{}
	_, err := issuerPubKey.SetBytes(testIssuer.GetPublicKey())
	require.NoError(t, err)

	// In the multishow case we use eddsa for cred signing, so we need an additional ecdsa key pair
	// for the eth account of the issuer.
	privKeyContract, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privKeyContract, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000_000)

	// Deploy Bloom filter contract
	bloomAddr, _, bloomContract, err := onchainBloom.DeployBloom(auth, sim)
	require.NoError(t, err)
	sim.Commit()

	// Deploy ZKP Verifier
	zkpVerifierAddr, _, _, err := zkpContract.DeployZkp(auth, sim)
	require.NoError(t, err)
	sim.Commit()

	// Get pk coordinates of issuer public key
	x := big.NewInt(0)
	issuerPubKey.A.X.BigInt(x)
	y := big.NewInt(0)
	issuerPubKey.A.Y.BigInt(y)

	// Deploy verifier contract
	verifierAddress, _, verifierContract, err := onchainVerifier.DeployVerifier(auth, sim, bloomAddr, zkpVerifierAddr, x, y)
	require.NoError(t, err)
	sim.Commit()

	// Transfer ownership of Bloom filter to verifier
	tx, err := bloomContract.TransferOwnership(auth, verifierAddress)
	require.NoError(t, err)
	sim.Commit()
	_, err = sim.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)

	// Use real issuer
	domain := 100
	capacity := 10

	err = testIssuer.IssueCredentials(uint(domain))
	require.NoError(t, err)

	err = testIssuer.RevokeRandomCredentials(uint(capacity))
	require.NoError(t, err)

	artifact, _, _, epoch, err := testIssuer.GenRevocationArtifact()
	require.NoError(t, err)

	filter, hf, bitlen := artifact.GetOnChainFilter()
	_, err = verifierContract.Update(auth, filter, hf, bitlen)
	require.NoError(t, err)
	sim.Commit()

	prover, err := holder.NewRevocationTokenProver("../../zkp/sol/build/verifier.g16.pk", "../../zkp/sol/build/verifier.g16.vk")
	require.NoError(t, err)

	// --- Test valid credentials ---
	for _, cred := range testIssuer.GetAllValidCreds() {
		proof, proofBytes, witness, witnessBytes, err := prover.GenProof(*cred, epoch)
		require.NoError(t, err)
		require.NotNil(t, proof)
		require.NotNil(t, witness)

		result, err := verifierContract.CheckCredential(&bind.CallOpts{}, proofBytes, witnessBytes[2], witnessBytes[3])
		require.NoError(t, err)
		require.Zero(t, result.ErrorCode, "CheckCredential: Expected valid credential, got error code %d", result.ErrorCode)
		require.True(t, result.Valid, "CheckCredential: Expected credential to be valid")
	}

	// --- Test revoked credentials ---
	for _, cred := range testIssuer.GetAllRevokedCreds() {
		proof, proofBytes, witness, witnessBytes, err := prover.GenProof(*cred, epoch)
		require.NoError(t, err)
		require.NotNil(t, proof)
		require.NotNil(t, witness)

		result, err := verifierContract.CheckCredential(&bind.CallOpts{}, proofBytes, witnessBytes[2], witnessBytes[3])
		require.NoError(t, err)
		require.Equal(t, uint8(2), result.ErrorCode, "CheckCredential: Expected revoked credential (code 2), got %d", result.ErrorCode)
		require.False(t, result.Valid, "CheckCredential: Expected credential to be revoked")
	}
}

func BenchmarkMultiShow_GasCheckCredential(b *testing.B) {
	configs := []struct {
		name     string
		domain   int
		capacity int
	}{
		{"D1k_C50", 1_000, 50},
		{"D1k_C100", 1_000, 100},
		{"D10k_C500", 10_000, 500},
		{"D10k_C1k", 10_000, 1_000},
		{"D100k_C5k", 100_000, 5_000},
		{"D100k_C10k", 100_000, 10_000},
		{"D100k_C50k", 1_000_000, 50_000},
		{"D1M_C100k", 1_000_000, 100_000},
	}

	fmt.Println("Benchmark Gas Consumption of MultiShow CheckCredential (N = 500 credentials):")
	fmt.Println("| Domain   | Capacity | Avg Gas Used | ETH (1 Gwei) |")
	fmt.Println("|----------|----------|--------------|--------------|")

	for _, cfg := range configs {
		avgGas, err := runMultiShowBenchmark(cfg.domain, cfg.capacity)
		if err != nil {
			b.Fatalf("Benchmark failed for domain=%d, capacity=%d: %v", cfg.domain, cfg.capacity, err)
		}

		ethCost := float64(avgGas) * 1e-9 // gas price = 1 Gwei
		fmt.Printf("| %8d | %8d | %12d | %.9f |\n", cfg.domain, cfg.capacity, avgGas, ethCost)
	}
}

func runMultiShowBenchmark(domain, capacity int) (uint64, error) {
	testIssuer := issuer.NewIssuer(issuer.MultiShow)

	issuerPubKey := eddsa.PublicKey{}
	_, err := issuerPubKey.SetBytes(testIssuer.GetPublicKey())
	if err != nil {
		return 0, err
	}

	privKeyContract, err := crypto.GenerateKey()
	if err != nil {
		return 0, err
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privKeyContract, big.NewInt(1337))
	if err != nil {
		return 0, err
	}

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 30_000_000_000)

	// Deploy Bloom filter
	bloomAddr, _, bloomContract, err := onchainBloom.DeployBloom(auth, sim)
	if err != nil {
		return 0, err
	}
	sim.Commit()

	// Deploy Verifier
	x := big.NewInt(0)
	y := big.NewInt(0)
	issuerPubKey.A.X.BigInt(x)
	issuerPubKey.A.Y.BigInt(y)

	zkpVerifierAddr, _, _, err := zkpContract.DeployZkp(auth, sim)
	if err != nil {
		return 0, err
	}
	sim.Commit()

	verifierAddr, _, verifier, err := onchainVerifier.DeployVerifier(auth, sim, bloomAddr, zkpVerifierAddr, x, y)
	if err != nil {
		return 0, err
	}
	sim.Commit()

	// Transfer ownership
	tx, err := bloomContract.TransferOwnership(auth, verifierAddr)
	if err != nil {
		return 0, err
	}
	sim.Commit()
	_, err = sim.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return 0, err
	}

	// Issue credentials
	err = testIssuer.IssueCredentials(uint(domain))
	if err != nil {
		return 0, err
	}
	err = testIssuer.RevokeRandomCredentials(uint(capacity))
	if err != nil {
		return 0, err
	}

	artifact, _, _, epoch, err := testIssuer.GenRevocationArtifact()
	if err != nil {
		return 0, err
	}

	filter, hf, bitlen := artifact.GetOnChainFilter()
	_, err = verifier.Update(auth, filter, hf, bitlen)
	if err != nil {
		return 0, err
	}
	sim.Commit()

	prover, err := holder.NewRevocationTokenProver("../../zkp/sol/build/verifier.g16.pk", "../../zkp/sol/build/verifier.g16.vk")
	if err != nil {
		return 0, err
	}

	// Use first 500 valid credentials
	validCreds := testIssuer.GetAllValidCreds()
	n := 500

	if len(validCreds) < n {
		return 0, fmt.Errorf("expected at least %d valid credentials", n)
	}

	var totalGas uint64
	for i := 0; i < n; i++ {
		cred := validCreds[i]
		_, proofBytes, _, witnessBytes, err := prover.GenProof(*cred, epoch)
		if err != nil {
			return 0, err
		}

		tx, err := verifier.MeasureCheckCredentialGas(auth, proofBytes, witnessBytes[2], witnessBytes[3])
		if err != nil {
			return 0, err
		}
		sim.Commit()

		receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return 0, err
		}
		totalGas += receipt.GasUsed
	}

	return totalGas / uint64(n), nil
}
