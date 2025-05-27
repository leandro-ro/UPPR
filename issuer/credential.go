package issuer

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/klayoracle/go-ecvrf"
	"math/big"
)

type CredentialType uint8

const (
	OneShow   CredentialType = 0
	MultiShow CredentialType = 1
)

// Credential represents a credential containing a single VRF public key hash as attribute and a corresponding signature.
type Credential struct {
	PublicKeyVrfHash []byte          // PublicKeyVrfHash (attribute) is the hash of the VRF public key
	Signature        eddsa.Signature // Signature over the credential attribute (i.e. PublicKeyVrfHash)
}

func (c *Credential) Verify(issuer eddsa.PublicKey) (bool, error) {
	sigValid, err := issuer.Verify(c.Signature.Bytes(), c.PublicKeyVrfHash, mimc.NewMiMC())
	if err != nil {
		return false, err
	}
	return sigValid, nil
}

// InternalCredential represents a structured internal credential containing its id, VRF, status, and associated Credential.
// Instances of this structure are hold by the issuer internally.
type InternalCredential struct {
	ID            uint           // ID is the issuer internal identifier for the Credential.
	Type          CredentialType // Type denotes the specific category of CredentialType used within InternalCredential.
	Revoked       bool           // Revoked is the issuer internal revocation status.
	PrivateKeyVrf []byte         // PrivateKeyVrf is the VRF associated with the Credential.
	Credential    Credential     // Cred is the Credential associated with the InternalCredential.
}

// GenRevocationToken generates a revocation token and its proof based on a given unix epoch and credential type.
// It supports OneShow and MultiShow credential types. Errors if the type is unknown or token generation fails.
func (ic *InternalCredential) GenRevocationToken(unixEpoch int64) (token, proof []byte, error error) {
	epoch := make([]byte, 8)
	binary.BigEndian.PutUint64(epoch, uint64(unixEpoch))

	switch ic.Type {
	case OneShow:
		vrf := ecvrf.Secp256k1Sha256Tai
		vrfSecretKey := secp256k1.PrivKeyFromBytes(ic.PrivateKeyVrf)
		t, p, err := vrf.ProveSecp256k1(vrfSecretKey.ToECDSA(), epoch)
		if err != nil {
			return nil, nil, err
		}

		return t, p, nil
	case MultiShow:
		// Compute revocationToken = Hash(epoch || vrf secret key)
		hf := mimc.NewMiMC()
		_, err := hf.Write(epoch)
		if err != nil {
			return nil, nil, err
		}
		_, err = hf.Write(ic.PrivateKeyVrf)
		if err != nil {
			return nil, nil, err
		}
		return hf.Sum(nil), nil, nil
	default:
		return nil, nil, errors.New("unknown credential type")
	}
}

// TODO: Does this make sense here?
func (ic *InternalCredential) GetVrfPublicKey() *ecdsa.PublicKey {
	vrfSecretKey := secp256k1.PrivKeyFromBytes(ic.PrivateKeyVrf)
	return vrfSecretKey.PubKey().ToECDSA()
}

func NewInternalCredential(version CredentialType, id uint, issuer eddsa.PrivateKey) (*InternalCredential, error) {
	var vrfPrivKey []byte
	var xBig, yBig *big.Int
	var err error

	switch version {
	case OneShow:
		sk, err := secp256k1.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		pub := sk.PubKey()
		xBig, yBig = pub.X(), pub.Y()
		vrfPrivKey = sk.Serialize()

	case MultiShow:
		sk, err := eddsa.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		xBig = big.NewInt(0)
		sk.PublicKey.A.X.BigInt(xBig)
		yBig = big.NewInt(0)
		sk.PublicKey.A.Y.BigInt(yBig)

		vrfPrivKey = sk.Bytes()

	default:
		return nil, errors.New("unknown credential type")
	}

	vrfPubKeyHash, err := hashPublicKeyCoordinates(xBig, yBig)
	if err != nil {
		return nil, err
	}

	credSignature, err := signAttribute(issuer, vrfPubKeyHash)
	if err != nil {
		return nil, err
	}

	return &InternalCredential{
		ID:            id,
		Type:          version,
		Revoked:       false,
		PrivateKeyVrf: vrfPrivKey,
		Credential: Credential{
			PublicKeyVrfHash: vrfPubKeyHash,
			Signature:        credSignature,
		},
	}, nil
}

func hashPublicKeyCoordinates(x, y *big.Int) ([]byte, error) {
	mod := ecc.BN254.ScalarField()

	hfunc := mimc.NewMiMC()
	_, err := hfunc.Write(x.Mod(x, mod).Bytes())
	if err != nil {
		return nil, err
	}
	_, err = hfunc.Write(y.Mod(x, mod).Bytes())
	if err != nil {
		return nil, err
	}
	return hfunc.Sum(nil), nil
}

func signAttribute(issuerKey eddsa.PrivateKey, msg []byte) (eddsa.Signature, error) {
	sigBytes, err := issuerKey.Sign(msg, mimc.NewMiMC())
	if err != nil {
		return eddsa.Signature{}, err
	}

	var sig eddsa.Signature
	_, err = sig.SetBytes(sigBytes)
	if err != nil {
		return eddsa.Signature{}, err
	}

	return sig, nil
}
