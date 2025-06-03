package zkp

import (
	"PrivacyPreservingRevocationCode/zkp"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark/backend/witness"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	bn254eddsa "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/schema"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	edddsaInCircuit "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// compileSolidityContract compiles the verifier contract using solc
func compileSolidityContract(solFile string) (binPath, abiPath string, err error) {
	cmd := exec.Command("solc", "--bin", "--overwrite", "--abi", "--evm-version", "istanbul", solFile, "-o", "build")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("solc failed: %v\n%s", err, out)
	}
	contractName := "Verifier"
	return filepath.Join("build", contractName+".bin"), filepath.Join("build", contractName+".abi"), nil
}

// deployVerifierContract compiles and deploys the verifier to a simulated backend
func deployVerifierContract(t *testing.T, auth *bind.TransactOpts, backend *backends.SimulatedBackend) (common.Address, abi.ABI) {
	binPath, abiPath, err := compileSolidityContract("revocationTokenVerifier.sol")
	require.NoError(t, err)

	binData, err := os.ReadFile(binPath)
	require.NoError(t, err)
	abiData, err := os.ReadFile(abiPath)
	require.NoError(t, err)

	parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
	require.NoError(t, err)

	address, tx, _, err := bind.DeployContract(auth, parsedABI, common.FromHex(string(binData)), backend)
	require.NoError(t, err)
	backend.Commit()

	receipt, err := backend.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	t.Logf("Contract deployed at %s, gas used: %d", address.Hex(), receipt.GasUsed)

	return address, parsedABI
}

// generateTestAssignment generates a valid proof assignment and returns the witness and public inputs
func generateTestAssignment(t *testing.T) (witness.Witness, []*big.Int) {
	issuerSk, err := bn254eddsa.GenerateKey(rand.Reader)
	require.NoError(t, err)
	vrfKey, err := zkp.EddsaForCircuitKeyGen()
	require.NoError(t, err)

	hash := mimc.NewMiMC()
	msgHash, err := zkp.HashEddsaPublicKey(vrfKey.Pk)
	require.NoError(t, err)
	cred, err := issuerSk.Sign(msgHash, hash)
	require.NoError(t, err)
	token, epoch, err := zkp.GenCurrentRevocationToken(vrfKey.Sk)
	require.NoError(t, err)

	assignment := &zkp.RevocationTokenProof{
		VrfSecretKey: vrfKey.Sk,
		VrfPublicKey: vrfKey.Pk,
		CredSignature: func() edddsaInCircuit.Signature {
			sig := edddsaInCircuit.Signature{}
			sig.Assign(tedwards.BN254, cred)
			return sig
		}(),
		IssuerPubKey: edddsaInCircuit.PublicKey{
			A: twistededwards.Point{X: issuerSk.PublicKey.A.X, Y: issuerSk.PublicKey.A.Y},
		},
		RevocationToken: token,
		Epoch:           epoch,
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)
	publicWitness, err := witness.Public()
	require.NoError(t, err)

	tLeaf := reflect.TypeOf((*frontend.Variable)(nil)).Elem()
	sch, err := schema.New(assignment, tLeaf)
	require.NoError(t, err)
	jsonPub, err := publicWitness.ToJSON(sch)
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(jsonPub, &parsed))

	xStr := parsed["IssuerPubKey"].(map[string]any)["A"].(map[string]any)["X"].(string)
	x := new(big.Int)
	_, ok := x.SetString(xStr, 10)
	require.True(t, ok)

	yStr := parsed["IssuerPubKey"].(map[string]any)["A"].(map[string]any)["Y"].(string)
	y := new(big.Int)
	_, ok = y.SetString(yStr, 10)
	require.True(t, ok)

	rtStr := parsed["RevocationToken"].(string)
	rt := new(big.Int)
	_, ok = rt.SetString(rtStr, 10)
	require.True(t, ok)

	epochFloat := parsed["Epoch"].(float64)
	epochBig := big.NewInt(int64(epochFloat))

	pinput := []*big.Int{
		x,
		y,
		rt,
		epochBig,
	}

	return witness, pinput
}

func TestVerifierGasCosts(t *testing.T) {
	// === Environment Setup ===
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	backend := backends.NewSimulatedBackend(
		map[common.Address]core.GenesisAccount{
			crypto.PubkeyToAddress(key.PublicKey): {Balance: big.NewInt(1e18)},
		},
		10_000_000,
	)

	// === Contract Deployment ===
	addr, parsedABI := deployVerifierContract(t, auth, backend)

	// === Proof Generation ===
	var circuit zkp.RevocationTokenProof
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	pkFile, err := os.Open("build/verifer.g16.pk")
	require.NoError(t, err)
	defer pkFile.Close()
	pk := groth16.NewProvingKey(ecc.BN254)
	_, err = pk.ReadFrom(pkFile)
	require.NoError(t, err)

	witness, pubInputs := generateTestAssignment(t)
	proof, err := groth16.Prove(r1cs, pk, witness)
	require.NoError(t, err)

	// === Format Proof ===
	var buf bytes.Buffer
	_, err = proof.WriteRawTo(&buf)
	require.NoError(t, err)
	proofBytes := buf.Bytes()
	const fp = 32
	p := [8]*big.Int{
		new(big.Int).SetBytes(proofBytes[0*fp : 1*fp]),
		new(big.Int).SetBytes(proofBytes[1*fp : 2*fp]),
		new(big.Int).SetBytes(proofBytes[2*fp : 3*fp]),
		new(big.Int).SetBytes(proofBytes[3*fp : 4*fp]),
		new(big.Int).SetBytes(proofBytes[4*fp : 5*fp]),
		new(big.Int).SetBytes(proofBytes[5*fp : 6*fp]),
		new(big.Int).SetBytes(proofBytes[6*fp : 7*fp]),
		new(big.Int).SetBytes(proofBytes[7*fp : 8*fp]),
	}

	contract := bind.NewBoundContract(addr, parsedABI, backend, backend, backend)
	tx, err := contract.Transact(&bind.TransactOpts{
		From:    auth.From,
		Context: context.Background(),
		Signer:  auth.Signer,
	}, "verifyProof", p, pubInputs)
	require.NoError(t, err)
	backend.Commit()

	receipt, err := backend.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	t.Logf("Gas used by verifyProof: %d", receipt.GasUsed)
}
