package zkp

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	eddsaInCicuit "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestRevocationTokenProofCompile(t *testing.T) {
	var circuit RevocationTokenProof
	startCompile := time.Now()
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Fatalf("Failed to compile circuit: %v", err)
	}
	compileDuration := time.Since(startCompile)
	fmt.Printf("Compile time: %v, Number of Constrains: %v\n", compileDuration, r1.GetNbConstraints())
}

func TestRevocationTokenProofVerify(t *testing.T) {
	issuerSecretKey, err := eddsa.GenerateKey(rand.Reader) // Issuer Secret Key
	require.NoError(t, err)

	vrfSecretKey, err := eddsa.GenerateKey(rand.Reader) // VRF Secret Key
	require.NoError(t, err)

	// Get epoch as Unix timestamp
	epoch := time.Now().UTC().Unix()
	epochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(epochBytes, uint64(epoch))

	var vrfSkFr fr.Element
	vrfSkFr.SetBytes(vrfSecretKey.Bytes()[32:64]) // assuming 32–64 is scalar, adjust if needed
	skBytes := vrfSkFr.Bytes()

	// Hash with MiMC
	hf := mimc.NewMiMC()
	hashInput := make([]byte, 2*hf.BlockSize()) // hash input needs to be multiple of block size
	copy(hashInput[0:8], epochBytes)
	copy(hashInput[len(hashInput)-32:], skBytes[:]) // assuming scalar is at [32:64]
	_, err = hf.Write(hashInput)
	require.NoError(t, err)

	revocationToken := hf.Sum(nil) // 32 bytes
	require.NotNil(t, revocationToken)

	// Issuer signs vrfPublicKey
	vrfPublicKeyX := vrfSecretKey.PublicKey.A.X.Bytes()
	vrfPublicKeyY := vrfSecretKey.PublicKey.A.Y.Bytes()
	vrfPubKeyBytes := make([]byte, 2*fr.Bytes)
	copy(vrfPubKeyBytes[:32], vrfPublicKeyX[:])
	copy(vrfPubKeyBytes[32:], vrfPublicKeyY[:])

	cred, err := issuerSecretKey.Sign(vrfPubKeyBytes, mimc.NewMiMC())
	require.NoError(t, err)
	require.NotNil(t, cred)

	var circuit RevocationTokenProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	// In Circuit Type Conversion
	icVrfSecretKey := frontend.Variable(vrfSkFr)
	icVrfPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: vrfSecretKey.PublicKey.A.X, Y: vrfSecretKey.PublicKey.A.Y}}
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)
	icRevocationToken := frontend.Variable(revocationToken)
	icEpoch := frontend.Variable(epochBytes)

	assignment := &RevocationTokenProof{
		VrfSecretKey:    icVrfSecretKey,
		VrfPublicKey:    icVrfPublicKey,
		IssuerPubKey:    icIssuerPublicKey,
		CredSignature:   icCredSigInCircuit,
		RevocationToken: icRevocationToken,
		Epoch:           icEpoch,
	}

	pk, vk, err := groth16.Setup(r1)
	assert.Nil(t, err)

	w, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	assert.Nil(t, err)

	startProve := time.Now()
	proof, err := groth16.Prove(r1, pk, w)
	assert.Nil(t, err)
	proveDuration := time.Since(startProve)
	fmt.Printf("Prove time: %v\n", proveDuration)

	wpub, err := w.Public()
	assert.Nil(t, err)

	startVerify := time.Now()
	err = groth16.Verify(proof, vk, wpub)
	assert.Nil(t, err)
	verifyDuration := time.Since(startVerify)
	fmt.Printf("Verify time: %v\n", verifyDuration)
}

func TestCredProof(t *testing.T) {
	issuerSecretKey, err := eddsa.GenerateKey(rand.Reader) // Issuer Secret Key
	require.NoError(t, err)

	vrfSecretKey, err := eddsa.GenerateKey(rand.Reader) // VRF Secret Key
	require.NoError(t, err)

	// 1. Get X and Y as bytes (each is [32]byte)
	xBytes := vrfSecretKey.PublicKey.A.X.Bytes()
	yBytes := vrfSecretKey.PublicKey.A.Y.Bytes()

	// 2. Hash them using MiMC. Double-hash for simplicity of circut input
	h := mimc.NewMiMC()
	_, err = h.Write(xBytes[:])
	require.NoError(t, err)
	_, err = h.Write(yBytes[:])
	require.NoError(t, err)

	msgHash := h.Sum(nil) // []byte, 32 bytes

	// Sign packed message
	hash := mimc.NewMiMC()
	cred, err := issuerSecretKey.Sign(msgHash, hash)
	require.NoError(t, err)
	require.NotNil(t, cred)

	var circuit CredProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	// In Circuit Type Conversion
	icVrfPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: vrfSecretKey.PublicKey.A.X, Y: vrfSecretKey.PublicKey.A.Y}}
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)

	assignment := &CredProof{
		VrfPublicKey:  icVrfPublicKey,
		IssuerPubKey:  icIssuerPublicKey,
		CredSignature: icCredSigInCircuit,
	}

	pk, vk, err := groth16.Setup(r1)
	assert.Nil(t, err)

	w, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	assert.Nil(t, err)

	startProve := time.Now()
	proof, err := groth16.Prove(r1, pk, w)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)
	proveDuration := time.Since(startProve)
	fmt.Printf("Prove time: %v\n", proveDuration)

	wpub, err := w.Public()
	assert.Nil(t, err)

	startVerify := time.Now()
	err = groth16.Verify(proof, vk, wpub)
	assert.Nil(t, err)
	verifyDuration := time.Since(startVerify)
	fmt.Printf("Verify time: %v\n", verifyDuration)
}

func TestVrfKeyPairCheck(t *testing.T) {
	vrfSecretKey, err := eddsa.GenerateKey(rand.Reader) // VRF Secret Key
	require.NoError(t, err)

	var vrfSkFr fr.Element
	vrfSkFr.SetBytes(vrfSecretKey.Bytes()[32:64]) // assuming 32–64 is scalar, adjust if needed

	// In Circuit Type Conversion
	icVrfSecretKey := vrfSkFr
	icVrfPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: vrfSecretKey.PublicKey.A.X, Y: vrfSecretKey.PublicKey.A.Y}}

	fmt.Printf("VrfSecretKey: %v\n", vrfSkFr)
	fmt.Printf("VrfPublicKey: %v\n", icVrfPublicKey)

	var circuit VrfKeyPairCheck
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	pk, vk, err := groth16.Setup(r1)
	assert.Nil(t, err)

	assignment := &VrfKeyPairCheck{
		VrfSecretKey: icVrfSecretKey,
		VrfPublicKey: icVrfPublicKey,
	}
	w, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	assert.Nil(t, err)

	startProve := time.Now()
	proof, err := groth16.Prove(r1, pk, w)
	assert.Nil(t, err)
	if err != nil {
		fmt.Println(err)
	}
	proveDuration := time.Since(startProve)
	fmt.Printf("Prove time: %v\n", proveDuration)

	wpub, err := w.Public()
	assert.Nil(t, err)

	startVerify := time.Now()
	err = groth16.Verify(proof, vk, wpub)
	assert.Nil(t, err)
	verifyDuration := time.Since(startVerify)
	fmt.Printf("Verify time: %v\n", verifyDuration)
}
