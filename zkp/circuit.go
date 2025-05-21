package zkp

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type RevocationTokenProof struct {
	VrfSecretKey frontend.Variable // VRF Secret Key (private input)

	VrfPublicKey eddsa.PublicKey `gnark:",public"` // VRF Public Key, i.e. single Credential Attribute
	IssuerPubKey eddsa.PublicKey `gnark:",public"` // Issuer Public Key

	CredSignature   eddsa.Signature   `gnark:",public"` // Signature on VrfPublicKey by IssuerPubKey
	RevocationToken frontend.Variable `gnark:",public"` // Revocation Token, i.e. vrf output
	Epoch           frontend.Variable `gnark:",public"` // Epoch for Revocation Token
}

func (p *RevocationTokenProof) Define(api frontend.API) error {
	// 1. Step: Verify Vrf Public Key to be signed by the Issuer
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}
	hashsig, _ := mimc.NewMiMC(api)
	err = eddsa.Verify(curve, p.CredSignature, p.VrfSecretKey, p.IssuerPubKey, &hashsig) // Curve already got api
	if err != nil {
		return err
	}

	// 2. Step: Verify Vrf Secret Key (to match VrfPublicKey)
	p.VrfPublicKey.A.Y // TODO: Step 2

	// 3. Step: Check Revocation Token
	hashrt, _ := mimc.NewMiMC(api)
	hashrt.Write(p.Epoch, p.VrfSecretKey)
	api.AssertIsEqual(p.RevocationToken, hashrt.Sum())
}
