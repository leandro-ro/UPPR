package multishow

import (
	onchainBloom "PrivacyPreservingRevocationCode/bloom/sol/build"
	"PrivacyPreservingRevocationCode/holder"
	"PrivacyPreservingRevocationCode/issuer"
	onchainVerifier "PrivacyPreservingRevocationCode/verifier/multishow/build"
	"PrivacyPreservingRevocationCode/zkp"
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
	zkpVerifierAddr, abi, _, err := zkp.DeployVerifier(auth, sim)
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

	validTestCred := testIssuer.GetAllRevokedCreds()[0]
	prover, err := holder.NewRevocationTokenProver()
	require.NoError(t, err)

	proof, proofBytes, witness, witnessBytes, err := prover.GenProof(*validTestCred, epoch)
	require.NoError(t, err)
	require.NotNil(t, proof)
	require.NotNil(t, witness)

	publicWitness, err := witness.Public()
	require.NoError(t, err)
	require.NoError(t, prover.VerifyProof(proof, publicWitness))

	contract := bind.NewBoundContract(zkpVerifierAddr, abi, sim, sim, sim)
	tx, err = contract.Transact(&bind.TransactOpts{
		From:    auth.From,
		Context: context.Background(),
		Signer:  auth.Signer,
	}, "verifyProof", proofBytes, witnessBytes)
	require.NoError(t, err)
	sim.Commit()

	fmt.Printf("Key x input %v\n", x)
	fmt.Printf("Key y input %v\n", y)

	result, err := verifierContract.CheckCredential(&bind.CallOpts{}, proofBytes, witnessBytes[0], witnessBytes[1], witnessBytes[2], witnessBytes[3])
	fmt.Printf("CheckCredential result: %v\n", result)
	require.NoError(t, err)
	require.Zero(t, result.ErrorCode, "CheckCredential: Expected valid credential, got error code %d", result.ErrorCode)
	require.True(t, result.Valid, "CheckCredential: Expected credential to be valid")
}
