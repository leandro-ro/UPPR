package issuer

import (
	"PrivacyPreservingRevocationCode/zkp"
	"crypto/ecdsa"
	"errors"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type RevocationToken []byte

func (t RevocationToken) ToBytes() []byte {
	return []byte(t)
}

func RevocationTokensToByteSlices(tokens []RevocationToken) [][]byte {
	result := make([][]byte, len(tokens))
	for i, t := range tokens {
		result[i] = t
	}
	return result
}

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

func (v *VrfKeyPair) GetOneShowPublicKey() (*ecdsa.PublicKey, error) {
	if v.version != OneShow {
		return nil, errors.New("not supported for this credential type")
	}
	return secp256k1.PrivKeyFromBytes(v.PrivateKey).PubKey().ToECDSA(), nil
}

func (v *VrfKeyPair) GetMultiShowPublicKey() (*eddsa.PublicKey, error) {
	if v.version != MultiShow {
		return nil, errors.New("not supported for this credential type")
	}

	x := fr.NewElement(0)
	x.SetBytes(v.publicKey.x)
	y := fr.NewElement(0)
	y.SetBytes(v.publicKey.y)

	return &eddsa.PublicKey{A: twistededwards.NewPointAffine(x, y)}, nil
}

func (v *VrfKeyPair) GetPublicKeyForOnChain() ([]byte, error) {
	if v.version != OneShow {
		return nil, errors.New("GetPublicKeyForOnChain only supported for OneShow credentials")
	}
	pubKey := secp256k1.PrivKeyFromBytes(v.PrivateKey).PubKey()
	return pubKey.SerializeCompressed(), nil
}

func NewVrfKeyPair(version CredentialType) (*VrfKeyPair, error) {
	var privateKey []byte
	var publicKeyHash []byte
	var xBytes, yBytes []byte

	switch version {
	case OneShow:
		sk, err := secp256k1.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		xBytes = sk.PubKey().X().Bytes()
		yBytes = sk.PubKey().Y().Bytes()

		privateKey = sk.Serialize()

		compressed := sk.PubKey().SerializeCompressed() // 33 bytes: 0x02/0x03 || X
		publicKeyHash = crypto.Keccak256(compressed)

	case MultiShow:
		sk, err := zkp.EddsaForCircuitKeyGen()
		if err != nil {
			return nil, err
		}

		var xBig, yBig big.Int
		x := sk.Pk.A.X.(fr.Element)
		x.BigInt(&xBig)
		y := sk.Pk.A.Y.(fr.Element)
		y.BigInt(&yBig)

		xBytes = xBig.Bytes()
		yBytes = yBig.Bytes()
		privateKey = sk.Sk
		publicKeyHash, err = zkp.HashEddsaPublicKey(sk.Pk)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unknown credential type")
	}

	return &VrfKeyPair{
		PrivateKey:       privateKey,
		publicKey:        VrfPublicKey{xBytes, yBytes},
		PublicKeyVrfHash: publicKeyHash,
		version:          version,
	}, nil

}
