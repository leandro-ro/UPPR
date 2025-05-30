package holder

import (
	"PrivacyPreservingRevocationCode/issuer"
	"PrivacyPreservingRevocationCode/zkp"
	"encoding/binary"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	eddsaInCicuit "github.com/consensys/gnark/std/signature/eddsa"
	"math/big"
)

type RevocationTokenProver struct {
	cs constraint.ConstraintSystem // cs represents the constraint system used in zero-knowledge proof generation.
	pk groth16.ProvingKey          // pk represents the Groth16 proving key used for generating zero-knowledge proofs.
	vk groth16.VerifyingKey        // vk represents the Groth16 verifying key used for verifying zero-knowledge proofs.
}

func NewRevocationTokenProver() (*RevocationTokenProver, error) {
	var circuit zkp.RevocationTokenProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return nil, err
	}

	pk, vk, err := groth16.Setup(r1)
	if err != nil {
		return nil, err
	}

	return &RevocationTokenProver{r1, pk, vk}, nil
}

func (r *RevocationTokenProver) GenProof(cred issuer.InternalCredential, epochUnix int64) (groth16.Proof, witness.Witness, error) {
	if cred.Type != issuer.MultiShow {
		return nil, nil, fmt.Errorf("credential type is not supported")
	}

	token, _, err := cred.GenRevocationToken(epochUnix)
	if err != nil {
		return nil, nil, err
	}

	pkVrf, err := cred.VrfKeyPair.GetMultiShowPublicKey()
	if err != nil {
		return nil, nil, err
	}

	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred.Credential.Signature.Bytes())

	icEpoch := make([]byte, 8)
	binary.BigEndian.PutUint64(icEpoch, uint64(epochUnix))

	icVrfPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: pkVrf.A.X, Y: pkVrf.A.Y}}
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: cred.IssuerPublicKey.A.X, Y: cred.IssuerPublicKey.A.Y}}
	icToken := big.NewInt(0).SetBytes(token)

	assignment := &zkp.RevocationTokenProof{
		VrfSecretKey:    cred.VrfKeyPair.PrivateKey,
		VrfPublicKey:    icVrfPublicKey,
		IssuerPubKey:    icIssuerPublicKey,
		CredSignature:   icCredSigInCircuit,
		RevocationToken: icToken,
		Epoch:           icEpoch,
	}

	fullWitness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, err
	}

	proof, err := groth16.Prove(r.cs, r.pk, fullWitness)
	if err != nil {
		return nil, nil, err
	}

	return proof, fullWitness, nil
}

func (r *RevocationTokenProver) VerifyProof(proof groth16.Proof, publicWitness witness.Witness) error {
	return groth16.Verify(proof, r.vk, publicWitness)
}
