package zkp

import (
	"PrivacyPreservingRevocationCode/zkp"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	bn254eddsa "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/schema"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	eddsaInCicuit "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// Helper to compile Solidity and return paths to .bin and .abi
func compileSolidityContract(solFile string) (string, string, error) {
	cmd := exec.Command(
		"solc",
		"--bin",
		"--overwrite",
		"--abi",
		"--evm-version", "istanbul",
		solFile,
		"-o", "build",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("solc failed: %s\n%s", err, string(output))
	}

	return filepath.Join("build", "Verifier.bin"), filepath.Join("build", "Verifier.abi"), nil
}

// Test function that deploys the contract and prints deployment gas cost
func TestDeployRevocationTokenVerifier(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("key generation failed: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		t.Fatalf("auth setup failed: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	alloc := map[common.Address]core.GenesisAccount{
		address: {Balance: big.NewInt(1e18)},
	}

	sim := backends.NewSimulatedBackend(alloc, 10_000_000)

	_, _, err = compileSolidityContract("revocationTokenVerifier.sol")
	if err != nil {
		t.Fatalf("contract compilation failed: %v", err)
	}

	contractName := "Verifier"
	binPath := filepath.Join("build", contractName+".bin")
	abiPath := filepath.Join("build", contractName+".abi")

	binData, err := os.ReadFile(binPath)
	if err != nil {
		t.Fatalf("reading bin failed: %v", err)
	}
	abiData, err := os.ReadFile(abiPath)
	if err != nil {
		t.Fatalf("reading abi failed: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		t.Fatalf("ABI parsing failed: %v", err)
	}

	contractAddr, tx, _, err := bind.DeployContract(auth, parsedABI, common.FromHex(string(binData)), sim)
	if err != nil {
		t.Fatalf("contract deployment failed: %v", err)
	}

	sim.Commit()

	receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		t.Fatalf("could not get receipt: %v", err)
	}

	t.Logf("Contract deployed at: %s", contractAddr.Hex())
	t.Logf("Deployment gas used: %d", receipt.GasUsed)
}

func TestVerifyRevocationTokenProof(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	alloc := map[common.Address]core.GenesisAccount{address: {Balance: big.NewInt(1e18)}}
	backend := backends.NewSimulatedBackend(alloc, 10_000_000)

	_, _, err = compileSolidityContract("revocationTokenVerifier.sol")
	require.NoError(t, err)

	binData, err := os.ReadFile("build/Verifier.bin")
	require.NoError(t, err)
	abiData, err := os.ReadFile("build/Verifier.abi")
	require.NoError(t, err)

	parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
	require.NoError(t, err)

	contractAddr, _, _, err := bind.DeployContract(auth, parsedABI, common.FromHex(string(binData)), backend)
	require.NoError(t, err)
	backend.Commit()

	// === Setup ===
	var circuit zkp.RevocationTokenProof
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	pkFile, err := os.Open("cubic.g16.pk")
	require.NoError(t, err)
	defer pkFile.Close()
	pk := groth16.NewProvingKey(ecc.BN254)
	_, err = pk.ReadFrom(pkFile)
	require.NoError(t, err)

	issuerSk, err := bn254eddsa.GenerateKey(rand.Reader)
	require.NoError(t, err)
	vrfKey, err := zkp.EddsaForCircuitKeyGen()
	require.NoError(t, err)

	msgHash, err := zkp.HashEddsaPublicKey(vrfKey.Pk)
	require.NoError(t, err)
	hash := mimc.NewMiMC()
	cred, err := issuerSk.Sign(msgHash, hash)
	require.NoError(t, err)

	token, epoch, err := zkp.GenCurrentRevocationToken(vrfKey.Sk)
	require.NoError(t, err)

	assignment := &zkp.RevocationTokenProof{
		VrfSecretKey: vrfKey.Sk,
		VrfPublicKey: vrfKey.Pk,
		CredSignature: func() eddsaInCicuit.Signature {
			var sig eddsaInCicuit.Signature
			sig.Assign(tedwards.BN254, cred)
			return sig
		}(),
		IssuerPubKey: eddsaInCicuit.PublicKey{A: twistededwards.Point{
			X: issuerSk.PublicKey.A.X,
			Y: issuerSk.PublicKey.A.Y,
		}},
		RevocationToken: token,
		Epoch:           epoch,
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)

	publicWitness, err := witness.Public()
	require.NoError(t, err)

	// Convert public witness to JSON using schema
	tLeaf := reflect.TypeOf((*frontend.Variable)(nil)).Elem()
	schema, err := schema.New(assignment, tLeaf)
	require.NoError(t, err)

	jsonPub, err := publicWitness.ToJSON(schema)
	require.NoError(t, err)

	// Unmarshal JSON into a generic map
	var parsed map[string]any
	err = json.Unmarshal(jsonPub, &parsed)
	require.NoError(t, err)

	// Extract and parse values
	pubInputBig := []*big.Int{}

	// IssuerPubKey.A.X
	xStr := parsed["IssuerPubKey"].(map[string]any)["A"].(map[string]any)["X"].(string)
	x := new(big.Int)
	x.SetString(xStr, 10)
	pubInputBig = append(pubInputBig, x)

	// IssuerPubKey.A.Y
	yStr := parsed["IssuerPubKey"].(map[string]any)["A"].(map[string]any)["Y"].(string)
	y := new(big.Int)
	y.SetString(yStr, 10)
	pubInputBig = append(pubInputBig, y)

	// RevocationToken
	rtStr := parsed["RevocationToken"].(string)
	rt := new(big.Int)
	rt.SetString(rtStr, 10)
	pubInputBig = append(pubInputBig, rt)

	// Epoch (int to big.Int)
	epochFloat := parsed["Epoch"].(float64) // JSON numbers are float64 by default
	epochBig := big.NewInt(int64(epochFloat))
	pubInputBig = append(pubInputBig, epochBig)

	proof, err := groth16.Prove(r1cs, pk, witness)
	require.NoError(t, err)

	// === Format Proof ===
	var buf bytes.Buffer
	_, err = proof.WriteRawTo(&buf)
	require.NoError(t, err)
	proofBytes := buf.Bytes()

	const fpSize = 32
	a := [2]*big.Int{
		new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1]),
		new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2]),
	}
	b := [2][2]*big.Int{
		{new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3]), new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])},
		{new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5]), new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])},
	}
	c := [2]*big.Int{
		new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7]),
		new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8]),
	}

	verifierContract := bind.NewBoundContract(contractAddr, parsedABI, backend, backend, backend)
	require.NotNil(t, verifierContract)

	// Flatten a, b, c into proof [8]*big.Int
	a[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	a[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	b[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	b[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	b[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	b[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	c[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	c[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])

	p := [8]*big.Int{a[0], a[1], b[0][0], b[0][1], b[1][0], b[1][1], c[0], c[1]}

	tx, err := verifierContract.Transact(&bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		Context:  context.Background(),
		GasLimit: 0, // Let Ethereum estimate the gas (optional)
	}, "verifyProof", p, pubInputBig)
	require.NoError(t, err)

	// Commit if using simulated backend
	backend.Commit()

	// Get the receipt
	receipt, err := backend.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)

	t.Logf("Gas used by verifyProof: %d", receipt.GasUsed)
}
