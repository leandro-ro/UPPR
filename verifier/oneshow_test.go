package verifier_test

import (
	onchainBloom "PrivacyPreservingRevocationCode/bloom/sol/build"
	"PrivacyPreservingRevocationCode/issuer"
	onchainVerifier "PrivacyPreservingRevocationCode/verifier/build"
	"context"
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
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

func TestUpdateAndCheckCredential(t *testing.T) {
	issuer := issuer.NewIssuer(issuer.OneShow)
	privKey, err := crypto.ToECDSA(issuer.GetPrivateKey())
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000)

	// Deploy Bloom filter contract
	bloomAddr, _, bloomContract, err := onchainBloom.DeployBloom(auth, sim)
	require.NoError(t, err)
	sim.Commit()

	// Deploy verifier contract
	veriferAddress, _, verifierContract, err := onchainVerifier.DeployVerifier(auth, sim, bloomAddr)
	require.NoError(t, err)
	sim.Commit()

	// Transfer ownership of Bloom filter to verifier
	tx, err := bloomContract.TransferOwnership(auth, veriferAddress)
	require.NoError(t, err)
	sim.Commit()
	_, err = sim.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)

	// Use real issuer
	domain := 1000
	capacity := 100

	err = issuer.IssueCredentials(uint(domain))
	require.NoError(t, err)

	err = issuer.RevokeRandomCredentials(uint(capacity))
	require.NoError(t, err)

	artifact, _, _, epoch, err := issuer.GenRevocationArtifact()
	require.NoError(t, err)

	filter, hf, bitlen := artifact.GetOnChainFilter()
	_, err = verifierContract.Update(auth, filter, hf, bitlen)
	require.NoError(t, err)
	sim.Commit()

	cred := issuer.GetAllValidCreds()[0]
	require.NoError(t, err)
	token, _, err := cred.GenRevocationToken(epoch)
	require.NoError(t, err)

	compressed := secp256k1.PrivKeyFromBytes(cred.VrfKeyPair.PrivateKey).PubKey().SerializeUncompressed()

	result, err := verifierContract.CheckCredential(&bind.CallOpts{}, compressed, cred.Credential.Signature, token.ToBytes())
	require.NoError(t, err)
	require.Zero(t, result.ErrorCode)
	require.True(t, result.Valid)
}
