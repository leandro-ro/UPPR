package holder

import (
	"PrivacyPreservingRevocationCode/issuer"
	"PrivacyPreservingRevocationCode/zkp"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/schema"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	eddsaInCicuit "github.com/consensys/gnark/std/signature/eddsa"
	"math/big"
	"os"
	"reflect"
)

// RevocationTokenProver is a struct for generating and verifying zero-knowledge proofs for credential revocation tokens.
// It uses Groth16 proving and verifying keys, as well as a constraint system for proof construction and validation.
type RevocationTokenProver struct {
	cs constraint.ConstraintSystem // cs represents the constraint system used in zero-knowledge proof generation.
	pk groth16.ProvingKey          // pk represents the Groth16 proving key used for generating zero-knowledge proofs.
	vk groth16.VerifyingKey        // vk represents the Groth16 verifying key used for verifying zero-knowledge proofs.
}

// NewRevocationTokenProver initializes and returns a new RevocationTokenProver instance for generating and verifying proofs.
func NewRevocationTokenProver(pkPath, vkPath string) (*RevocationTokenProver, error) {
	var circuit zkp.RevocationTokenProof
	r1, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return nil, err
	}

	pkFile, err := os.Open(pkPath)
	if err != nil {
		return nil, err
	}
	pk := groth16.NewProvingKey(ecc.BN254)
	_, err = pk.ReadFrom(pkFile)

	vkFile, err := os.Open(vkPath)
	if err != nil {
		return nil, err
	}
	vk := groth16.NewVerifyingKey(ecc.BN254)
	_, err = vk.ReadFrom(vkFile)

	return &RevocationTokenProver{r1, pk, vk}, nil
}

// GenProof generates a zero-knowledge proof for a credential's revocation token based on the provided epoch timestamp.
// It supports MultiShow credential types and returns the proof, proof in byte array, witness, and an error if any occur.
func (r *RevocationTokenProver) GenProof(cred issuer.InternalCredential, epochUnix int64) (proof groth16.Proof, proofBytes [8]*big.Int, witness witness.Witness, witnessBytes [4]*big.Int, err error) {
	if cred.Credential.Type != issuer.MultiShow {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, fmt.Errorf("credential type is not supported")
	}

	token, _, err := cred.GenRevocationToken(epochUnix)
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	pkVrf, err := cred.VrfKeyPair.GetMultiShowPublicKey()
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	icCredSigInCircuit := eddsaInCicuit.Signature{}
	icCredSigInCircuit.Assign(tedwards.BN254, cred.Credential.Signature)

	icEpoch := make([]byte, 8)
	binary.BigEndian.PutUint64(icEpoch, uint64(epochUnix))

	issPubKey := eddsa.PublicKey{}
	_, err = issPubKey.SetBytes(cred.IssuerPublicKey)
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	icVrfPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: pkVrf.A.X, Y: pkVrf.A.Y}}
	icIssuerPublicKey := eddsaInCicuit.PublicKey{A: twistededwards.Point{X: issPubKey.A.X, Y: issPubKey.A.Y}}
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
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	proof, err = groth16.Prove(r.cs, r.pk, fullWitness)
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	proofBytes, err = groth16ProofToOnChainInput(proof)
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	tLeaf := reflect.TypeOf((*frontend.Variable)(nil)).Elem()
	sch, err := schema.New(assignment, tLeaf)
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	publicWitness, err := fullWitness.Public()
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	jsonPub, err := publicWitness.ToJSON(sch)
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	var parsed map[string]any
	err = json.Unmarshal(jsonPub, &parsed)
	if err != nil {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, err
	}

	xStr := parsed["IssuerPubKey"].(map[string]any)["A"].(map[string]any)["X"].(string)
	x := new(big.Int)
	_, ok := x.SetString(xStr, 10)
	if !ok {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, fmt.Errorf("invalid revocation token")
	}

	yStr := parsed["IssuerPubKey"].(map[string]any)["A"].(map[string]any)["Y"].(string)
	y := new(big.Int)
	_, ok = y.SetString(yStr, 10)
	if !ok {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, fmt.Errorf("invalid revocation token")
	}

	rtStr := parsed["RevocationToken"].(string)
	rt := new(big.Int)
	_, ok = rt.SetString(rtStr, 10)
	if !ok {
		return nil, [8]*big.Int{}, nil, [4]*big.Int{}, fmt.Errorf("invalid revocation token")
	}

	epochFloat := parsed["Epoch"].(float64)
	epochBig := big.NewInt(int64(epochFloat))

	return proof, proofBytes, fullWitness, [4]*big.Int{x, y, rt, epochBig}, nil
}

func (r *RevocationTokenProver) VerifyProof(proof groth16.Proof, publicWitness witness.Witness) error {
	return groth16.Verify(proof, r.vk, publicWitness)
}

// groth16ProofToOnChainInput converts a Groth16 proof to a [8]*big.Int array suitable for on-chain verification.
// It writes the proof in raw format and slices it into 8 segments, each representing a 256-bit field element.
// Returns the converted [8]*big.Int array and an error if the conversion fails.
func groth16ProofToOnChainInput(proof groth16.Proof) ([8]*big.Int, error) {
	var buf bytes.Buffer
	_, err := proof.WriteRawTo(&buf)
	if err != nil {
		return [8]*big.Int{}, err
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
	}, nil
}
