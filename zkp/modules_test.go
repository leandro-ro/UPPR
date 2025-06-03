package zkp

import (
	"crypto/rand"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	bn254eddsa "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	eddsaInCicuit "github.com/consensys/gnark/std/signature/eddsa"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModule_CredProof(t *testing.T) {
	issuerSecretKey, err := bn254eddsa.GenerateKey(rand.Reader) // Issuer Secret Key
	require.NoError(t, err)

	vrfKey, err := EddsaForCircuitKeyGen()
	require.NoError(t, err)

	// 1. Hash vrf public key via MiMC.
	msgHash, err := HashEddsaPublicKey(vrfKey.Pk)
	require.NoError(t, err)

	// Sign hash of vrf public key.
	hash := mimc.NewMiMC()
	cred, err := issuerSecretKey.Sign(msgHash, hash)
	require.NoError(t, err)
	require.NotNil(t, cred)

	var circuit CredProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	// In Circuit Type Conversion
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issuerSecretKey.PublicKey.A.X, Y: issuerSecretKey.PublicKey.A.Y}}
	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred)

	assignment := &CredProof{
		VrfPublicKey:  vrfKey.Pk,
		IssuerPubKey:  icIssuerPublicKey,
		CredSignature: icCredSigInCircuit,
	}

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)

	// Check constraint satisfaction
	_, err = r1.Solve(witness)
	require.NoError(t, err)
}

func TestModule_VrfKeyPairProof(t *testing.T) {
	vrfKeyPair, err := EddsaForCircuitKeyGen()
	require.NoError(t, err)

	// Compile circuit
	var circuit VrfKeyPairProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	// Build witness
	assignment := &VrfKeyPairProof{
		VrfSecretKey: vrfKeyPair.Sk,
		VrfPublicKey: vrfKeyPair.Pk,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)

	// Check constraint satisfaction
	_, err = r1.Solve(witness)
	require.NoError(t, err)
}

func TestModule_TokenHashProof(t *testing.T) {
	vrfKeyPair, err := EddsaForCircuitKeyGen()
	require.NoError(t, err)

	// Compute revocationToken = Hash(epoch || sk)
	token, epoch, err := GenCurrentRevocationToken(vrfKeyPair.Sk)
	require.NoError(t, err)

	// Compile circuit
	var circuit TokenHashProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	require.NoError(t, err)

	// Build witness
	assignment := &TokenHashProof{
		VrfSecretKey:    vrfKeyPair.Sk,
		RevocationToken: token,
		Epoch:           epoch,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	require.NoError(t, err)

	// Check constraint satisfaction
	_, err = r1.Solve(witness)
	require.NoError(t, err)
}
