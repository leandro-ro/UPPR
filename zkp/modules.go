package zkp

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

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

	return assertSignaturePk(api, curve, p.CredSignature, p.VrfPublicKey, p.IssuerPubKey)
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

	return assertKeyPair(api, curve, p.VrfSecretKey, p.VrfPublicKey)
}

type TokenHashProof struct {
	VrfSecretKey    frontend.Variable
	RevocationToken frontend.Variable `gnark:",public"` // Revocation Token, i.e. vrf output
	Epoch           frontend.Variable `gnark:",public"` // Epoch for Revocation Token}
}

func (p *TokenHashProof) Define(api frontend.API) error {
	return assertRevocationToken(api, p.Epoch, p.VrfSecretKey, p.RevocationToken)
}

// assertKeyPair validates if the provided secretKey and publicKey form a valid key pair for the given elliptic curve.
func assertKeyPair(api frontend.API, curve twistededwards.Curve, secretKey frontend.Variable, publicKey eddsa.PublicKey) error {
	base := twistededwards.Point{X: curve.Params().Base[0].Bytes(), Y: curve.Params().Base[1].Bytes()}
	expectedPublicKey := curve.ScalarMul(base, secretKey)

	api.AssertIsEqual(publicKey.A.X, expectedPublicKey.X)
	api.AssertIsEqual(publicKey.A.Y, expectedPublicKey.Y)
	return nil
}

// assertSignaturePk verifies a signature against a public key message and an issuer public key using the provided curve.
// It uses the MiMC hash function and checks the validity of the signature using the eddsa.Verify method.
func assertSignaturePk(api frontend.API, curve twistededwards.Curve, signature eddsa.Signature, publicKeyMessage, issuerPublicKey eddsa.PublicKey) error {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	h.Write(publicKeyMessage.A.X)
	h.Write(publicKeyMessage.A.Y)
	msg := h.Sum()

	hashsig, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	return eddsa.Verify(curve, signature, msg, issuerPublicKey, &hashsig)
}

// assertRevocationToken ensures the validity of a revocation token by comparing it with a computed hash using MiMC.
func assertRevocationToken(api frontend.API, epoch, secretKey, revocationToken frontend.Variable) error {
	expectedToken, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	expectedToken.Write(epoch, secretKey)

	api.AssertIsEqual(revocationToken, expectedToken.Sum())
	return nil
}
