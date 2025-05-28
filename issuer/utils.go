package issuer

import (
	"PrivacyPreservingRevocationCode/zkp"
	"crypto/ecdsa"
	"errors"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"math/big"
)

type VrfPublicKey struct {
	x []byte
	y []byte
}

type VrfKeyPair struct {
	PrivateKey       []byte
	publicKey        VrfPublicKey
	PublicKeyVrfHash []byte
	version          CredentialType
}

func (v *VrfKeyPair) GetEcdsaPublicKey() (*ecdsa.PublicKey, error) {
	if v.version != OneShow {
		return nil, errors.New("not supported for this credential type")
	}
	return secp256k1.PrivKeyFromBytes(v.PrivateKey).PubKey().ToECDSA(), nil
}

func (v *VrfKeyPair) GetEddsaPublicKey() (*eddsa.PublicKey, error) {
	if v.version != MultiShow {
		return nil, errors.New("not supported for this credential type")
	}

	point := twistededwards.NewPointAffine()
	return eddsa.PublicKey{A: twistededwards.}
}

func NewVrfKeyPair(version CredentialType) (*VrfKeyPair, error) {
	var privateKey []byte
	var xBig, yBig *big.Int

	switch version {
	case OneShow:
		sk, err := secp256k1.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		pub := sk.PubKey()
		xBig, yBig = pub.X(), pub.Y()
		privateKey = sk.Serialize()

	case MultiShow:
		sk, err := zkp.EddsaForCircuitKeyGen()
		if err != nil {
			return nil, err
		}
		xBig = big.NewInt(0)
		xBytes := sk.Pk.A.X.(fr.Element)
		xBytes.BigInt(xBig)
		yBig = big.NewInt(0)
		yBytes := sk.Pk.A.Y.(fr.Element)
		yBytes.BigInt(yBig)

		privateKey = sk.Sk

	default:
		return nil, errors.New("unknown credential type")
	}

	publicKeyHash, err := hashPublicKeyCoordinates(xBig, yBig)
	if err != nil {
		return nil, err
	}

	return &VrfKeyPair{
		PrivateKey:       privateKey,
		publicKey:        VrfPublicKey{xBig.Bytes(), yBig.Bytes()},
		PublicKeyVrfHash: publicKeyHash,
		version:          version,
	}, nil

}
