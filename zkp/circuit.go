package zkp

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type RevocationTokenProof struct {
	VrfSecretKey  frontend.Variable // VRF Secret Key
	VrfPublicKey  eddsa.PublicKey   // VRF Public Key, i.e. single Credential Attribute
	CredSignature eddsa.Signature   // Signature on VrfPublicKey by IssuerPubKey

	IssuerPubKey    eddsa.PublicKey   `gnark:",public"` // Issuer Public Key
	RevocationToken frontend.Variable `gnark:",public"` // Revocation Token, i.e. vrf output
	Epoch           frontend.Variable `gnark:",public"` // Epoch for Revocation Token
}

func (p *RevocationTokenProof) Define(api frontend.API) error {
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	// 1. Verify VRF Key Pair, i.e. that the public key is derived from the given secret key.
	err = assertKeyPair(api, curve, p.VrfSecretKey, p.VrfPublicKey)
	if err != nil {
		return err
	}

	// 2. Verify signature of issuer on given public key (i.e., credential presentation)
	err = assertSignaturePk(api, curve, p.CredSignature, p.VrfPublicKey, p.IssuerPubKey)
	if err != nil {
		return err
	}

	// 3. Verify the revocation token.
	err = assertRevocationToken(api, p.Epoch, p.VrfSecretKey, p.RevocationToken)
	if err != nil {
		return err
	}

	return nil
}
