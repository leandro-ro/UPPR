package zkp

import (
	"encoding/binary"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	bn254ted "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	eddsaInCicuit "github.com/consensys/gnark/std/signature/eddsa"
	"math/big"
	"time"
)

// EddsaKeyPair represents an EdDSA key pair consisting of a secret key
// and a public key in a zkp circuit compliant format.
type EddsaKeyPair struct {
	Sk *big.Int
	Pk eddsaInCicuit.PublicKey
}

// eddsaFrKeyGen generates an EdDSA key pair with the secret key sampled within the BN254 scalar field.
// It computes the corresponding public key using the elliptic curve's base point and formats both keys for compatibility.
func eddsaFrKeyGen() (EddsaKeyPair, error) {
	var maxFr fr.Element
	(*fr.Element).SetBigInt(&maxFr, ecc.BN254.ScalarField())

	// We generate a secret key < Fr to avoid implicit reduction inside the circuit.
	// gnark reduces big.Int values mod Fr when assigning them to frontend.Variable,
	// which would alter the original scalar and break the key pair check pk = sk × G.
	// Therefore, we explicitly sample the scalar within the field Fr to ensure
	// circuit and off-circuit behavior remain consistent.
	sk, err := (*fr.Element).SetRandom(&maxFr)
	if err != nil {
		return EddsaKeyPair{}, err
	}

	// Convert secret key to big.Int
	skBig := new(big.Int)
	sk.BigInt(skBig)

	// Compute public key: pk = sk × G
	var pk bn254ted.PointAffine
	base := bn254ted.GetEdwardsCurve().Base
	pk.ScalarMultiplication(&base, skBig)

	// Convert to in-circuit types
	icSk := skBig
	icPk := eddsaInCicuit.PublicKey{
		A: twistededwards.Point{
			X: pk.X,
			Y: pk.Y,
		},
	}

	return EddsaKeyPair{Sk: icSk, Pk: icPk}, nil
}

// hashEddsaPublicKey hashes an EdDSA public key using the MiMC hashing algorithm.
// It takes an EdDSA public key as input and returns the corresponding hash or an error if hashing fails.
func hashEddsaPublicKey(pk eddsaInCicuit.PublicKey) ([]byte, error) {
	var xBig, yBig big.Int
	xBytes := pk.A.X.(fr.Element)
	xBytes.BigInt(&xBig)
	yBytes := pk.A.Y.(fr.Element)
	yBytes.BigInt(&yBig)

	h := mimc.NewMiMC()
	_, err := h.Write(xBig.Bytes())
	if err != nil {
		return nil, err
	}
	_, err = h.Write(yBig.Bytes())
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// genCurrentRevocationToken generates a revocation token  Hash(epoch || sk) for the current epoch using a VRF secret key.
// It returns the token as a big.Int, the epoch as a byte slice, and an error if any occurs during execution.
func genCurrentRevocationToken(vrfSecretKey *big.Int) (token *big.Int, epoch []byte, err error) {
	epochUnix := time.Now().UTC().Unix()
	epoch = make([]byte, 8)
	binary.BigEndian.PutUint64(epoch, uint64(epochUnix))

	// Compute revocationToken = Hash(epoch || sk)
	hf := mimc.NewMiMC()
	_, err = hf.Write(epoch)
	if err != nil {
		return nil, nil, err
	}
	_, err = hf.Write(vrfSecretKey.Bytes())
	if err != nil {
		return nil, nil, err
	}
	revocationToken := hf.Sum(nil)
	return new(big.Int).SetBytes(revocationToken), epoch, nil
}
