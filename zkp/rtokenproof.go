package zkp

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
	"math/big"
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
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	// 1. Step: Verify Vrf Public Key to be signed by the Issuer
	hashsig, _ := mimc.NewMiMC(api)
	shift := new(big.Int).Lsh(big.NewInt(1), 128)
	api.Add(p.VrfPublicKey.A.X, api.Mul(p.VrfPublicKey.A.Y, shift))

	err = eddsa.Verify(curve, p.CredSignature, shift, p.IssuerPubKey, &hashsig)
	if err != nil {
		return err
	}

	// 2. Step: Verify Vrf Secret Key (to match VrfPublicKey)
	base := twistededwards.Point{X: curve.Params().Base[0], Y: curve.Params().Base[1]}
	pktest := curve.ScalarMul(base, p.VrfSecretKey)
	api.AssertIsEqual(p.VrfPublicKey.A.X, pktest.X)
	api.AssertIsEqual(p.VrfPublicKey.A.Y, pktest.Y)

	// 3. Step: Check Revocation Token
	hashrt, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	hashrt.Write(p.Epoch, p.VrfSecretKey)
	api.AssertIsEqual(p.RevocationToken, hashrt.Sum())
	return nil
}

type CredProof struct {
	VrfPublicKey  eddsa.PublicKey `gnark:",public"` // VRF Public Key, i.e. single Credential Attribute
	IssuerPubKey  eddsa.PublicKey `gnark:",public"` // Issuer Public Key
	CredSignature eddsa.Signature `gnark:",public"` // Signature on VrfPublicKey by IssuerPubKey
}

func (p *CredProof) Define(api frontend.API) error {
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	h, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	h.Write(p.VrfPublicKey.A.X)
	h.Write(p.VrfPublicKey.A.Y)
	msg := h.Sum()

	// 1. Step: Verify Vrf Public Key to be signed by the Issuer
	hashsig, _ := mimc.NewMiMC(api)

	return eddsa.Verify(curve, p.CredSignature, msg, p.IssuerPubKey, &hashsig)
}

type VrfKeyPairProof struct {
	VrfSecretKey frontend.Variable
	VrfPublicKey eddsa.PublicKey `gnark:",public"` // VRF Public Key, i.e. single Credential Attribute
}

func (p *VrfKeyPairProof) Define(api frontend.API) error {
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	base := twistededwards.Point{X: curve.Params().Base[0].Bytes(), Y: curve.Params().Base[1].Bytes()}
	expectedPublicKey := curve.ScalarMul(base, p.VrfSecretKey)

	api.AssertIsEqual(p.VrfPublicKey.A.X, expectedPublicKey.X)
	api.AssertIsEqual(p.VrfPublicKey.A.Y, expectedPublicKey.Y)
	return nil
}

type TokenHashProof struct {
	VrfSecretKey    frontend.Variable
	RevocationToken frontend.Variable `gnark:",public"` // Revocation Token, i.e. vrf output
	Epoch           frontend.Variable `gnark:",public"` // Epoch for Revocation Token}
}

func (p *TokenHashProof) Define(api frontend.API) error {
	expectedHash, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	expectedHash.Write(p.Epoch, p.VrfSecretKey)
	api.AssertIsEqual(p.RevocationToken, expectedHash.Sum())
	return nil
}
