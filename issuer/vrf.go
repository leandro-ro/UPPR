package issuer

import (
	"encoding/hex"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/klayoracle/go-ecvrf" // should be compatible with
)

type VerifiableRandomFunction interface {
	GetRevocationToken(epoch uint) (token []byte, proof []byte, err error)
	GetPublicKeyHash() []byte
}

type VrfOneShow struct {
	ecvrf.VRF
}

func NewVrfOneShow(seed []byte) (VerifiableRandomFunction, error) {

	vrfSecretKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	alpha, _ := hex.DecodeString("73616d706c65")

	vrf := ecvrf.Secp256k1Sha256Tai
	_, _, err = vrf.ProveSecp256k1(vrfSecretKey.ToECDSA(), alpha)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (v VrfOneShow) GetRevocationToken(epoch uint) (token []byte, proof []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (v VrfOneShow) GetPublicKeyHash() []byte {
	//TODO implement me
	panic("implement me")
}
