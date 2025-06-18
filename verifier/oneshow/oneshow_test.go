package oneshow

import (
	onchainBloom "PrivacyPreservingRevocationCode/bloom/sol/build"
	"PrivacyPreservingRevocationCode/issuer"
	onchainVerifier "PrivacyPreservingRevocationCode/verifier/oneshow/build"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

/*
func TestOneShow_CompileAndGenBindings(t *testing.T) {
	buildDir := "build"
	solFile := "OneShowVerifier.sol"

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

func TestOneShow_EndToEnd(t *testing.T) {
	iss := issuer.NewIssuer(issuer.OneShow)
	privKey, err := crypto.ToECDSA(iss.GetPrivateKey())
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000_000)

	// Deploy Bloom filter contract
	bloomAddr, _, bloomContract, err := onchainBloom.DeployBloom(auth, sim)
	require.NoError(t, err)
	sim.Commit()

	// Deploy verifier contract
	verifierAddress, _, verifierContract, err := onchainVerifier.DeployVerifier(auth, sim, bloomAddr)
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

	err = iss.IssueCredentials(uint(domain))
	require.NoError(t, err)

	err = iss.RevokeRandomCredentials(uint(capacity))
	require.NoError(t, err)

	artifact, _, _, epoch, err := iss.GenRevocationArtifact()
	require.NoError(t, err)

	filter, hf, bitlen := artifact.GetOnChainFilter()
	_, err = verifierContract.Update(auth, filter, hf, bitlen)
	require.NoError(t, err)
	sim.Commit()

	// --- Test valid credentials ---
	for _, cred := range iss.GetAllValidCreds() {
		_, proof, err := cred.GenRevocationToken(epoch)
		require.NoError(t, err)

		ok, err := cred.Credential.Verify(iss.GetPublicKey())
		require.NoError(t, err)
		require.True(t, ok)

		pubkey, err := cred.VrfKeyPair.GetPublicKeyForOnChain()
		require.NoError(t, err)
		result, err := verifierContract.CheckCredential(&bind.CallOpts{}, pubkey, cred.Credential.Signature, proof, big.NewInt(epoch))
		require.NoError(t, err)
		require.Zero(t, result.ErrorCode, "CheckCredential: Expected valid credential, got error code %d", result.ErrorCode)
		require.True(t, result.Valid, "CheckCredential: Expected credential to be valid")
	}

	// --- Test revoked credentials ---
	for _, cred := range iss.GetAllRevokedCreds() {
		_, proof, err := cred.GenRevocationToken(epoch)
		require.NoError(t, err)

		ok, err := cred.Credential.Verify(iss.GetPublicKey())
		require.NoError(t, err)
		require.True(t, ok)

		pubkey, err := cred.VrfKeyPair.GetPublicKeyForOnChain()
		require.NoError(t, err)

		result, err := verifierContract.CheckCredential(&bind.CallOpts{}, pubkey, cred.Credential.Signature, proof, big.NewInt(epoch))
		require.NoError(t, err)
		require.Equal(t, uint8(4), result.ErrorCode, "CheckCredential: Expected revoked credential (code 3), got %d", result.ErrorCode)
		require.False(t, result.Valid, "CheckCredential: Expected credential to be revoked")
	}

	otherIss := issuer.NewIssuer(issuer.OneShow)
	err = otherIss.IssueCredentials(1)
	require.NoError(t, err)
	foreignCred := otherIss.GetAllValidCreds()[0]
	_, proof, err := foreignCred.GenRevocationToken(epoch)
	require.NoError(t, err)

	onchainKey, err := foreignCred.VrfKeyPair.GetPublicKeyForOnChain()
	require.NoError(t, err)

	result, err := verifierContract.CheckCredential(&bind.CallOpts{}, onchainKey, foreignCred.Credential.Signature, proof, big.NewInt(epoch))
	require.NoError(t, err)
	require.Equal(t, uint8(2), result.ErrorCode, "Credential from another issuer should fail signature check")
	require.False(t, result.Valid)
}

func TestOneShow_EndToEndFast(t *testing.T) {
	iss := issuer.NewIssuer(issuer.OneShow)
	privKey, err := crypto.ToECDSA(iss.GetPrivateKey())
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000_000)

	// Deploy Bloom filter contract
	bloomAddr, _, bloomContract, err := onchainBloom.DeployBloom(auth, sim)
	require.NoError(t, err)
	sim.Commit()

	// Deploy verifier contract
	verifierAddress, _, verifierContract, err := onchainVerifier.DeployVerifier(auth, sim, bloomAddr)
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

	err = iss.IssueCredentials(uint(domain))
	require.NoError(t, err)

	err = iss.RevokeRandomCredentials(uint(capacity))
	require.NoError(t, err)

	artifact, _, _, epoch, err := iss.GenRevocationArtifact()
	require.NoError(t, err)

	filter, hf, bitlen := artifact.GetOnChainFilter()
	_, err = verifierContract.Update(auth, filter, hf, bitlen)
	require.NoError(t, err)
	sim.Commit()

	// --- Test valid credentials ---
	for _, cred := range iss.GetAllValidCreds() {
		_, proof, err := cred.GenRevocationToken(epoch)
		require.NoError(t, err)

		ok, err := cred.Credential.Verify(iss.GetPublicKey())
		require.NoError(t, err)
		require.True(t, ok)

		pubkey, err := cred.VrfKeyPair.GetPublicKeyForOnChain()
		require.NoError(t, err)

		precomputedParams, err := verifierContract.GetFastVerifyParams(&bind.CallOpts{}, pubkey, proof, big.NewInt(epoch))
		require.NoError(t, err)

		resultFast, err := verifierContract.CheckCredentialFast(&bind.CallOpts{}, pubkey, cred.Credential.Signature, proof, big.NewInt(epoch), precomputedParams.UPoint, precomputedParams.VComponents)
		require.NoError(t, err)
		require.Zero(t, resultFast.ErrorCode, "CheckCredentialFast: Expected valid credential, got error code %d", resultFast.ErrorCode)
		require.True(t, resultFast.Valid, "CheckCredentialFast: Expected credential to be valid")
	}

	// --- Test revoked credentials ---
	for _, cred := range iss.GetAllRevokedCreds() {
		_, proof, err := cred.GenRevocationToken(epoch)
		require.NoError(t, err)

		ok, err := cred.Credential.Verify(iss.GetPublicKey())
		require.NoError(t, err)
		require.True(t, ok)

		pubkey, err := cred.VrfKeyPair.GetPublicKeyForOnChain()
		require.NoError(t, err)

		precomputedParams, err := verifierContract.GetFastVerifyParams(&bind.CallOpts{}, pubkey, proof, big.NewInt(epoch))
		require.NoError(t, err)

		resultFast, err := verifierContract.CheckCredentialFast(&bind.CallOpts{}, pubkey, cred.Credential.Signature, proof, big.NewInt(epoch), precomputedParams.UPoint, precomputedParams.VComponents)
		require.NoError(t, err)
		require.Equal(t, uint8(4), resultFast.ErrorCode, "CheckCredentialFast: Expected revoked credential (code 3), got %d", resultFast.ErrorCode)
		require.False(t, resultFast.Valid, "CheckCredentialFast: Expected credential to be valid")
	}

	otherIss := issuer.NewIssuer(issuer.OneShow)
	err = otherIss.IssueCredentials(1)
	require.NoError(t, err)
	foreignCred := otherIss.GetAllValidCreds()[0]
	_, proof, err := foreignCred.GenRevocationToken(epoch)
	require.NoError(t, err)

	onchainKey, err := foreignCred.VrfKeyPair.GetPublicKeyForOnChain()
	require.NoError(t, err)

	precomputedParams, err := verifierContract.GetFastVerifyParams(&bind.CallOpts{}, onchainKey, proof, big.NewInt(epoch))
	require.NoError(t, err)

	result, err := verifierContract.CheckCredentialFast(&bind.CallOpts{}, onchainKey, foreignCred.Credential.Signature, proof, big.NewInt(epoch), precomputedParams.UPoint, precomputedParams.VComponents)
	require.NoError(t, err)
	require.Equal(t, uint8(2), result.ErrorCode, "Credential from another issuer should fail signature check")
	require.False(t, result.Valid)
}

// BenchmarkOneShow_PrecomputeFastParams benchmarks the generation of fast verification parameters off-chain.
func BenchmarkOneShow_PrecomputeFastParams(b *testing.B) {
	domain := 10_000
	capacity := 1_000

	testIssuer := issuer.NewIssuer(issuer.OneShow)
	privKey, err := crypto.ToECDSA(testIssuer.GetPrivateKey())
	if err != nil {
		b.Fatalf("Failed to get private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	if err != nil {
		b.Fatalf("Failed to create transactor: %v", err)
	}

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 30_000_000_000)

	bloomAddr, _, bloomContract, err := onchainBloom.DeployBloom(auth, sim)
	if err != nil {
		b.Fatalf("Failed to deploy bloom: %v", err)
	}
	sim.Commit()

	verifierAddr, _, verifier, err := onchainVerifier.DeployVerifier(auth, sim, bloomAddr)
	if err != nil {
		b.Fatalf("Failed to deploy verifier: %v", err)
	}
	sim.Commit()

	tx, err := bloomContract.TransferOwnership(auth, verifierAddr)
	if err != nil {
		b.Fatalf("Failed to transfer ownership: %v", err)
	}
	sim.Commit()
	_, err = sim.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		b.Fatalf("Failed to get tx receipt: %v", err)
	}

	if err := testIssuer.IssueCredentials(uint(domain)); err != nil {
		b.Fatalf("IssueCredentials failed: %v", err)
	}
	if err := testIssuer.RevokeRandomCredentials(uint(capacity)); err != nil {
		b.Fatalf("RevokeRandomCredentials failed: %v", err)
	}

	artifact, _, _, epoch, err := testIssuer.GenRevocationArtifact()
	if err != nil {
		b.Fatalf("GenRevocationArtifact failed: %v", err)
	}

	filter, ks, lens := artifact.GetOnChainFilter()
	_, err = verifier.Update(auth, filter, ks, lens)
	if err != nil {
		b.Fatalf("Update bloom failed: %v", err)
	}
	sim.Commit()

	validCreds := testIssuer.GetAllValidCreds()
	if len(validCreds) == 0 {
		b.Fatalf("No valid credentials found")
	}
	cred := validCreds[0]

	_, proof, err := cred.GenRevocationToken(epoch)
	if err != nil {
		b.Fatalf("GenProof failed: %v", err)
	}
	pubKey, err := cred.VrfKeyPair.GetPublicKeyForOnChain()
	if err != nil {
		b.Fatalf("GetPublicKeyForOnChain failed: %v", err)
	}

	// Benchmark loop
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := verifier.GetFastVerifyParams(&bind.CallOpts{}, pubKey, proof, big.NewInt(epoch))
		if err != nil {
			b.Fatalf("GetFastVerifyParams failed: %v", err)
		}
	}
}

func BenchmarkOneShow_GasCheckCredential(b *testing.B) {
	configs := []struct {
		name     string
		domain   int
		capacity int
		fast     bool
	}{
		{"D1k_C50_Fast", 1_000, 50, true},
		{"D1k_C100_Fast", 1_000, 100, true},
		{"D10k_C500_Fast", 10_000, 500, true},
		{"D10k_C1k_Fast", 10_000, 1_000, true},
		{"D100k_C5k_Fast", 100_000, 5_000, true},
		{"D100k_C10k_Fast", 100_000, 10_000, true},
		{"D1M_C50k_Fast", 1_000_000, 50_000, true},
		{"D1M_C100k_Fast", 1_000_000, 100_000, true},
	}

	fmt.Println("Benchmark Gas Consumption and Local Runtime of OneShow CheckCredential (N = 500 credentials):")
	fmt.Println("| Domain   | Capacity |   Mode   | Avg Gas Used | ETH (1 Gwei) | Avg Local Time [ms] |")
	fmt.Println("|----------|----------|----------|--------------|--------------|---------------------|")

	for _, cfg := range configs {
		avgGas, avgTime, err := runOneShowBenchmark(cfg.domain, cfg.capacity, cfg.fast)
		if err != nil {
			b.Fatalf("Benchmark failed for domain=%d, capacity=%d: %v", cfg.domain, cfg.capacity, err)
		}

		ethCost := float64(avgGas) * 1e-9 // gas price = 1 Gwei
		mode := "Fast"
		if !cfg.fast {
			mode = "Slow"
		}

		fmt.Printf("| %8d | %8d | %8s | %12d | %.9f | %19.3f |\n",
			cfg.domain, cfg.capacity, mode, avgGas, ethCost, float64(avgTime.Microseconds())/1000.0)
	}
}

func runOneShowBenchmark(domain, capacity int, fast bool) (uint64, time.Duration, error) {
	testIssuer := issuer.NewIssuer(issuer.OneShow)
	privKey, err := crypto.ToECDSA(testIssuer.GetPrivateKey())
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

	bloomAddr, _, bloomContract, err := onchainBloom.DeployBloom(auth, sim)
	if err != nil {
		return 0, 0, err
	}
	sim.Commit()

	verifierAddr, _, verifier, err := onchainVerifier.DeployVerifier(auth, sim, bloomAddr)
	if err != nil {
		return 0, 0, err
	}
	sim.Commit()

	tx, err := bloomContract.TransferOwnership(auth, verifierAddr)
	if err != nil {
		return 0, 0, err
	}
	sim.Commit()
	_, err = sim.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return 0, 0, err
	}

	err = testIssuer.IssueCredentials(uint(domain))
	if err != nil {
		return 0, 0, err
	}
	err = testIssuer.RevokeRandomCredentials(uint(capacity))
	if err != nil {
		return 0, 0, err
	}

	artifact, _, _, epoch, err := testIssuer.GenRevocationArtifact()
	if err != nil {
		return 0, 0, err
	}

	filter, ks, lens := artifact.GetOnChainFilter()
	_, err = verifier.Update(auth, filter, ks, lens)
	if err != nil {
		return 0, 0, err
	}
	sim.Commit()

	validCreds := testIssuer.GetAllValidCreds()
	n := 500
	if len(validCreds) < n {
		return 0, 0, fmt.Errorf("expected %d valid creds, got %d", n, len(validCreds))
	}

	var totalGas uint64
	var totalTime time.Duration

	for i := 0; i < n; i++ {
		cred := validCreds[i]
		_, proof, err := cred.GenRevocationToken(epoch)
		if err != nil {
			return 0, 0, err
		}
		pubKey, err := cred.VrfKeyPair.GetPublicKeyForOnChain()
		if err != nil {
			return 0, 0, err
		}

		var tx *types.Transaction
		start := time.Now()
		elapsed := time.Since(start)
		if fast {
			precomputedParams, err := verifier.GetFastVerifyParams(&bind.CallOpts{}, pubKey, proof, big.NewInt(epoch))
			if err != nil {
				return 0, 0, err
			}

			tx, err = verifier.MeasureCheckCredentialFastGas(auth, pubKey, cred.Credential.Signature, proof, big.NewInt(epoch), precomputedParams.UPoint, precomputedParams.VComponents)
			if err != nil {
				return 0, 0, err
			}

			// Local runtime
			start = time.Now()
			_, err = verifier.CheckCredentialFast(&bind.CallOpts{}, pubKey, cred.Credential.Signature, proof, big.NewInt(epoch), precomputedParams.UPoint, precomputedParams.VComponents)
			elapsed = time.Since(start)
			if err != nil {
				return 0, 0, err
			}
		} else {
			tx, err = verifier.MeasureCheckCredentialGas(auth, pubKey, cred.Credential.Signature, proof, big.NewInt(epoch))
			if err != nil {
				return 0, 0, err
			}

			// Local runtime
			start = time.Now()
			_, err = verifier.CheckCredential(&bind.CallOpts{}, pubKey, cred.Credential.Signature, proof, big.NewInt(epoch))
			elapsed = time.Since(start)
			if err != nil {
				return 0, 0, err
			}
		}

		totalTime += elapsed

		sim.Commit()
		receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return 0, 0, err
		}
		totalGas += receipt.GasUsed
	}

	avgGas := totalGas / uint64(n)
	avgTime := totalTime / time.Duration(n)
	return avgGas, avgTime, nil
}
