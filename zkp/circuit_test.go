package zkp

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	bn254eddsa "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
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

func TestRevocationTokenProof_Compile(t *testing.T) {
	var circuit RevocationTokenProof
	_, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Fatalf("Failed to compile circuit: %v", err)
	}
}

func TestRevocationTokenProof_Verify(t *testing.T) {
	issuerSecretKey, err := bn254eddsa.GenerateKey(rand.Reader) // Issuer Secret Key
	require.NoError(t, err)

	vrfKey, err := eddsaFrKeyGen()
	require.NoError(t, err)

	// 1. Issue Credential.
	// 1a. Hash vrf public key via MiMC.
	msgHash, err := hashEddsaPublicKey(vrfKey.Pk)
	require.NoError(t, err)

	// 1b. Sign hash of vrf public key.
	hash := mimc.NewMiMC()
	cred, err := issuerSecretKey.Sign(msgHash, hash)
	require.NoError(t, err)
	require.NotNil(t, cred)

	// 2. Generate Revocation Token for current epoch.
	// token = Hash(epoch || sk)
	token, epoch, err := genCurrentRevocationToken(vrfKey.Sk)
	require.NoError(t, err)

	// 3. Convert parameters into proper format for witness generation.
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)

	// 4. Generate Circuit
	var circuit RevocationTokenProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	assignment := &RevocationTokenProof{
		VrfSecretKey:    vrfKey.Sk,
		VrfPublicKey:    vrfKey.Pk,
		IssuerPubKey:    icIssuerPublicKey,
		CredSignature:   icCredSigInCircuit,
		RevocationToken: token,
		Epoch:           epoch,
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	assert.Nil(t, err)

	// Check constraint satisfaction
	_, err = r1.Solve(witness)
	require.NoError(t, err)
}

func TestRevocationTokenProof_FalseCred(t *testing.T) {
	issuerSecretKey, err := bn254eddsa.GenerateKey(rand.Reader) // Issuer Secret Key
	require.NoError(t, err)

	vrfKey, err := eddsaFrKeyGen()
	require.NoError(t, err)

	// 1. Issue Credential.
	// 1a. Hash vrf public key via MiMC.
	msgHash, err := hashEddsaPublicKey(vrfKey.Pk)
	require.NoError(t, err)

	msgHash[0] = 0x00 // invalidate hash of vrf public key (i.e., "holder tries to use different credential")

	// 1b. Sign hash of vrf public key.
	hash := mimc.NewMiMC()
	cred, err := issuerSecretKey.Sign(msgHash, hash)
	require.NoError(t, err)
	require.NotNil(t, cred)

	// 2. Generate Revocation Token for current epoch.
	// token = Hash(epoch || sk)
	token, epoch, err := genCurrentRevocationToken(vrfKey.Sk)
	require.NoError(t, err)

	// 3. Convert parameters into proper format for witness generation.
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)

	// 4. Generate Circuit
	var circuit RevocationTokenProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	assignment := &RevocationTokenProof{
		VrfSecretKey:    vrfKey.Sk,
		VrfPublicKey:    vrfKey.Pk,
		IssuerPubKey:    icIssuerPublicKey,
		CredSignature:   icCredSigInCircuit,
		RevocationToken: token,
		Epoch:           epoch,
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	assert.Nil(t, err)

	// Check error
	_, err = r1.Solve(witness)
	require.Error(t, err)
}

func TestRevocationTokenProof_FalseSkForToken(t *testing.T) {
	issuerSecretKey, err := bn254eddsa.GenerateKey(rand.Reader) // Issuer Secret Key
	require.NoError(t, err)

	vrfKey, err := eddsaFrKeyGen()
	require.NoError(t, err)

	falseVrfKey, err := eddsaFrKeyGen()
	require.NoError(t, err)

	// 1. Issue Credential.
	// 1a. Hash vrf public key via MiMC.
	msgHash, err := hashEddsaPublicKey(vrfKey.Pk)
	require.NoError(t, err)

	// 1b. Sign hash of vrf public key.
	hash := mimc.NewMiMC()
	cred, err := issuerSecretKey.Sign(msgHash, hash)
	require.NoError(t, err)
	require.NotNil(t, cred)

	// 2. Generate Revocation Token for current epoch.
	// token = Hash(epoch || sk)
	token, epoch, err := genCurrentRevocationToken(falseVrfKey.Sk)
	require.NoError(t, err)

	// 3. Convert parameters into proper format for witness generation.
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)

	// 4. Generate Circuit
	var circuit RevocationTokenProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	assignment := &RevocationTokenProof{
		VrfSecretKey:    vrfKey.Sk,
		VrfPublicKey:    vrfKey.Pk,
		IssuerPubKey:    icIssuerPublicKey,
		CredSignature:   icCredSigInCircuit,
		RevocationToken: token,
		Epoch:           epoch,
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	assert.Nil(t, err)

	// Check constraint satisfaction
	_, err = r1.Solve(witness)
	require.Error(t, err)
}

func TestRevocationTokenProof_FalseEpoch(t *testing.T) {
	issuerSecretKey, err := bn254eddsa.GenerateKey(rand.Reader) // Issuer Secret Key
	require.NoError(t, err)

	vrfKey, err := eddsaFrKeyGen()
	require.NoError(t, err)

	// 1. Issue Credential.
	// 1a. Hash vrf public key via MiMC.
	msgHash, err := hashEddsaPublicKey(vrfKey.Pk)
	require.NoError(t, err)

	// 1b. Sign hash of vrf public key.
	hash := mimc.NewMiMC()
	cred, err := issuerSecretKey.Sign(msgHash, hash)
	require.NoError(t, err)
	require.NotNil(t, cred)

	// 2. Generate Revocation Token for current epoch.
	// token = Hash(epoch || sk)
	token, epoch, err := genCurrentRevocationToken(vrfKey.Sk)
	require.NoError(t, err)

	epochUnix := time.Now().UTC().Unix()
	epochUnix-- // Use old epoch
	epoch = make([]byte, 8)
	binary.BigEndian.PutUint64(epoch, uint64(epochUnix))

	// 3. Convert parameters into proper format for witness generation.
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)

	// 4. Generate Circuit
	var circuit RevocationTokenProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	assignment := &RevocationTokenProof{
		VrfSecretKey:    vrfKey.Sk,
		VrfPublicKey:    vrfKey.Pk,
		IssuerPubKey:    icIssuerPublicKey,
		CredSignature:   icCredSigInCircuit,
		RevocationToken: token,
		Epoch:           epoch,
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	assert.Nil(t, err)

	// Check constraint satisfaction
	_, err = r1.Solve(witness)
	require.Error(t, err)
}

func BenchmarkRevocationTokenProof_ConstraintCount(b *testing.B) {
	if b.N == 1 {
		var circuit RevocationTokenProof
		ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
		require.NoError(b, err)
		fmt.Printf("Constraints: %d | Public Inputs: %d | Secret Inputs: %d\n", ccs.GetNbConstraints(), ccs.GetNbPublicVariables(), ccs.GetNbSecretVariables())
	} else {
		b.SkipNow()
	}
}

func BenchmarkRevocationTokenProof_Prove(b *testing.B) {
	issuerSecretKey, _ := bn254eddsa.GenerateKey(rand.Reader)
	vrfKey, _ := eddsaFrKeyGen()
	msgHash, _ := hashEddsaPublicKey(vrfKey.Pk)
	hash := mimc.NewMiMC()
	cred, _ := issuerSecretKey.Sign(msgHash, hash)
	token, epoch, _ := genCurrentRevocationToken(vrfKey.Sk)

	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)

	assignment := &RevocationTokenProof{
		VrfSecretKey:    vrfKey.Sk,
		VrfPublicKey:    vrfKey.Pk,
		IssuerPubKey:    icIssuerPublicKey,
		CredSignature:   icCredSigInCircuit,
		RevocationToken: token,
		Epoch:           epoch,
	}

	var circuit RevocationTokenProof
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(b, err)

	// Generate trusted setup
	pk, _, err := groth16.Setup(ccs)
	require.NoError(b, err)

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := groth16.Prove(ccs, pk, witness)
		require.NoError(b, err)
	}
}

func BenchmarkRevocationTokenProof_Verify(b *testing.B) {
	issuerSecretKey, _ := bn254eddsa.GenerateKey(rand.Reader)
	vrfKey, _ := eddsaFrKeyGen()
	msgHash, _ := hashEddsaPublicKey(vrfKey.Pk)
	hash := mimc.NewMiMC()
	cred, _ := issuerSecretKey.Sign(msgHash, hash)
	token, epoch, _ := genCurrentRevocationToken(vrfKey.Sk)

	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)

	assignment := &RevocationTokenProof{
		VrfSecretKey:    vrfKey.Sk,
		VrfPublicKey:    vrfKey.Pk,
		IssuerPubKey:    icIssuerPublicKey,
		CredSignature:   icCredSigInCircuit,
		RevocationToken: token,
		Epoch:           epoch,
	}

	var circuit RevocationTokenProof
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(b, err)

	pk, vk, err := groth16.Setup(ccs)
	require.NoError(b, err)

	fullWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	require.NoError(b, err)

	publicWitness, err := fullWitness.Public()
	require.NoError(b, err)

	proof, err := groth16.Prove(ccs, pk, fullWitness)
	require.NoError(b, err)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := groth16.Verify(proof, vk, publicWitness)
		require.NoError(b, err)
	}
}
