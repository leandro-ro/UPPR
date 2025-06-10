package zkp

import (
	"PrivacyPreservingRevocationCode/zkp"
	zkpContract "PrivacyPreservingRevocationCode/zkp/sol/build"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark/backend/witness"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
func TestVerifier_CompileAndGenBindings(t *testing.T) {
	buildDir := "build"
	solFile := "RevocationTokenVerifier.sol"

	// Compile Solidity file with solc
	cmd := exec.Command(
		"solc",
		"--bin",
		"--abi",
		"--overwrite",
		"--evm-version", "istanbul",
		solFile,
		"-o", buildDir,
	)
	out, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "solc failed: %v\n%s", err, string(out))

	// Extract contract name
	contractName := "Verifier"

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
		"--pkg=zkp",
		"--out="+bindingPath,
	)
	abigenOut, err := abigenCmd.CombinedOutput()
	require.NoErrorf(t, err, "abigen failed: %v\n%s", err, string(abigenOut))
	require.FileExists(t, bindingPath, "binding file not created")

	fmt.Printf("Generated bindings:\n- %s\n- %s\n- %s\n", binPath, abiPath, bindingPath)
}
*/

func TestVerifierProofEndToEnd(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)

	alloc := core.GenesisAlloc{
		auth.From: {Balance: big.NewInt(1_000_000_000_000_000_000)},
	}
	sim := backends.NewSimulatedBackend(alloc, 3_000_000_000)

	_, _, contract, err := zkpContract.DeployZkp(auth, sim)
	require.NoError(t, err)
	sim.Commit()

	var circuit zkp.RevocationTokenProof
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	pkFile, err := os.Open("build/verifier.g16.pk")
	require.NoError(t, err)
	defer pkFile.Close()
	pk := groth16.NewProvingKey(ecc.BN254)
	_, err = pk.ReadFrom(pkFile)
	require.NoError(t, err)

	witness, pubInputs := generateTestAssignment(t)
	proof, err := groth16.Prove(r1cs, pk, witness)
	require.NoError(t, err)

	p := parseGroth16ProofToInputs(proof)

	err = contract.VerifyProof(&bind.CallOpts{}, p, pubInputs)
	require.NoError(t, err)
	sim.Commit()
}

func BenchmarkVerifierGasCosts(b *testing.B) {
	const N = 10

	var totalDeployGas uint64
	var totalVerifyGas uint64

	for i := 0; i < N; i++ {
		key, err := crypto.GenerateKey()
		require.NoError(b, err)
		auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		require.NoError(b, err)
		backend := backends.NewSimulatedBackend(
			map[common.Address]core.GenesisAccount{
				crypto.PubkeyToAddress(key.PublicKey): {Balance: big.NewInt(1e18)},
			},
			10_000_000,
		)

		addr, parsedABI, deployGas, err := deployVerifier(auth, backend)
		require.NoError(b, err)
		totalDeployGas += deployGas

		var circuit zkp.RevocationTokenProof
		r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
		require.NoError(b, err)

		pkFile, err := os.Open("build/verifer.g16.pk")
		require.NoError(b, err)
		defer pkFile.Close()
		pk := groth16.NewProvingKey(ecc.BN254)
		_, err = pk.ReadFrom(pkFile)
		require.NoError(b, err)

		witness, pubInputs := generateTestAssignment(b)
		proof, err := groth16.Prove(r1cs, pk, witness)
		require.NoError(b, err)

		p := parseGroth16ProofToInputs(proof)

		contract := bind.NewBoundContract(addr, parsedABI, backend, backend, backend)
		tx, err := contract.Transact(&bind.TransactOpts{
			From:    auth.From,
			Context: context.Background(),
			Signer:  auth.Signer,
		}, "verifyProof", p, pubInputs)
		require.NoError(b, err)
		backend.Commit()

		receipt, err := backend.TransactionReceipt(context.Background(), tx.Hash())
		require.NoError(b, err)
		totalVerifyGas += receipt.GasUsed
	}

	avgDeployGas := totalDeployGas / N
	avgVerifyGas := totalVerifyGas / N

	fmt.Println()
	fmt.Println("Groth16 zkSNARK Verifier Gas Costs (On-Chain)")
	fmt.Println("| Action            | Avg Gas Used | in ETH (1 Gwei)  |")
	fmt.Println("|-------------------|--------------|------------------|")
	fmt.Printf("| Contract Deploy   | %12d | %.9f ETH |\n", avgDeployGas, float64(avgDeployGas)*1e-9)
	fmt.Printf("| verifyProof()     | %12d | %.9f ETH |\n", avgVerifyGas, float64(avgVerifyGas)*1e-9)
}

func parseGroth16ProofToInputs(proof groth16.Proof) [8]*big.Int {
	var buf bytes.Buffer
	_, err := proof.WriteRawTo(&buf)
	if err != nil {
		panic(err)
	}
	proofBytes := buf.Bytes()
	const fp = 32
	return [8]*big.Int{
		new(big.Int).SetBytes(proofBytes[0*fp : 1*fp]),
		new(big.Int).SetBytes(proofBytes[1*fp : 2*fp]),
		new(big.Int).SetBytes(proofBytes[2*fp : 3*fp]),
		new(big.Int).SetBytes(proofBytes[3*fp : 4*fp]),
		new(big.Int).SetBytes(proofBytes[4*fp : 5*fp]),
		new(big.Int).SetBytes(proofBytes[5*fp : 6*fp]),
		new(big.Int).SetBytes(proofBytes[6*fp : 7*fp]),
		new(big.Int).SetBytes(proofBytes[7*fp : 8*fp]),
	}
}

// generateTestAssignment generates a valid proof assignment and returns the witness and public inputs
func generateTestAssignment(t testing.TB) (witness.Witness, [4]*big.Int) {
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

	pinput := [4]*big.Int{
		x,
		y,
		rt,
		epochBig,
	}

	return witness, pinput
}

// DeployVerifier deploys a zkSNARK Verifier contract to a simulated Ethereum backend and returns its address, ABI, gas used, and error.
func deployVerifier(auth *bind.TransactOpts, backend *backends.SimulatedBackend) (common.Address, abi.ABI, uint64, error) {
	binPath := filepath.Join("build", "Verifier.bin")
	abiPath := filepath.Join("build", "Verifier.abi")

	binData, err := os.ReadFile(binPath)
	if err != nil {
		return common.Address{}, abi.ABI{}, 0, err
	}
	abiData, err := os.ReadFile(abiPath)
	if err != nil {
		return common.Address{}, abi.ABI{}, 0, err
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		return common.Address{}, abi.ABI{}, 0, err
	}

	address, tx, _, err := bind.DeployContract(auth, parsedABI, common.FromHex(string(binData)), backend)
	if err != nil {
		return common.Address{}, abi.ABI{}, 0, err
	}
	backend.Commit()

	receipt, err := backend.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return common.Address{}, abi.ABI{}, 0, err
	}

	return address, parsedABI, receipt.GasUsed, nil
}
