package verifier_test

import (
	onchainBloom "PrivacyPreservingRevocationCode/bloom/sol/build"
	"PrivacyPreservingRevocationCode/issuer"
	onchainVerifier "PrivacyPreservingRevocationCode/verifier/build"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompileAndGenBindings(t *testing.T) {
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
		"--base-path", "../", // point to project root
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

	fmt.Printf("✅ Generated bindings:\n- %s\n- %s\n- %s\n", binPath, abiPath, bindingPath)
}

func TestEndToEnd(t *testing.T) {
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
	domain := 1000
	capacity := 100

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
		require.Zero(t, result.ErrorCode, "Expected valid credential, got error code %d", result.ErrorCode)
		require.True(t, result.Valid, "Expected credential to be valid")
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
		require.Equal(t, uint8(4), result.ErrorCode, "Expected revoked credential (code 3), got %d", result.ErrorCode)
		require.False(t, result.Valid, "Expected credential to be revoked")
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
