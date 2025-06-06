package issuer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vechain/go-ecvrf"
	"math/big"
)

type CredentialType uint8

const (
	OneShow   CredentialType = 0
	MultiShow CredentialType = 1
)

func (ct CredentialType) String() string {
	switch ct {
	case OneShow:
		return "OneShow"
	case MultiShow:
		return "MultiShow"
	default:
		return fmt.Sprintf("CredentialType(%d)", ct)
	}
}

// Credential represents a credential containing a single VRF public key hash as attribute and a corresponding signature.
type Credential struct {
	PublicKeyVrfHash []byte         // PublicKeyVrfHash (attribute) is the hash of the VRF public key
	Signature        []byte         // Signature over the credential attribute (i.e. PublicKeyVrfHash)
	Type             CredentialType // Type denotes the specific category of CredentialType used within InternalCredential.
}

func (c *Credential) Verify(issuerPublicKey []byte) (bool, error) {
	switch c.Type {
	case OneShow:
		if len(c.Signature) != 65 {
			return false, errors.New("signature must be 65 bytes")
		}

		sig := make([]byte, len(c.Signature))
		copy(sig, c.Signature)

		// Reverse the v normalization for Go's crypto.Ecrecover (expects v ∈ {0,1})
		if sig[64] >= 27 {
			sig[64] -= 27
		}

		// Recover uncompressed pubkey (65 bytes) from signature
		recoveredPubkeyBytes, err := crypto.Ecrecover(c.PublicKeyVrfHash, sig)
		if err != nil {
			return false, fmt.Errorf("ecrecover failed: %w", err)
		}

		recoveredPubkey, err := crypto.UnmarshalPubkey(recoveredPubkeyBytes)
		if err != nil {
			return false, fmt.Errorf("failed to unmarshal recovered pubkey: %w", err)
		}
		recoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)

		issuerPubkey, err := crypto.DecompressPubkey(issuerPublicKey)
		if err != nil {
			return false, fmt.Errorf("failed to decompress issuer pubkey: %w", err)
		}
		issuerAddress := crypto.PubkeyToAddress(*issuerPubkey)

		return recoveredAddress == issuerAddress, nil
	case MultiShow:
		issuerPk := eddsa.PublicKey{}
		_, err := issuerPk.SetBytes(issuerPublicKey)
		if err != nil {
			return false, err
		}

		sigValid, err := issuerPk.Verify(c.Signature, c.PublicKeyVrfHash, mimc.NewMiMC())
		if err != nil {
			return false, err
		}
		return sigValid, nil
	default:
		return false, errors.New("unknown credential type")
	}
}

// InternalCredential represents a structured internal credential containing its id, VRF, status, and associated Credential.
// Instances of this structure are hold by the issuer internally.
type InternalCredential struct {
	ID              uint        // ID is the issuer internal identifier for the Credential.
	Revoked         bool        // Revoked is the issuer internal revocation status.
	VrfKeyPair      *VrfKeyPair // PrivateKeyVrf is the VRF associated with the Credential.
	Credential      Credential  // Credential is the Credential associated with the InternalCredential.
	IssuerPublicKey []byte      // IssuerPublicKey is the public key of the credential issuer used to verify Credential.
}

func NewInternalCredential(version CredentialType, id uint, issuerPrivateKey []byte) (*InternalCredential, error) {
	vrfKeyPair, err := NewVrfKeyPair(version)
	if err != nil {
		return nil, err
	}

	credSignature, err := signAttribute(issuerPrivateKey, vrfKeyPair.PublicKeyVrfHash, version)
	if err != nil {
		return nil, err
	}

	issuerPkBytes := []byte{}
	if version == OneShow {
		privKey, err := crypto.ToECDSA(issuerPrivateKey)
		if err != nil {
			return nil, err
		}
		pubKey := crypto.FromECDSAPub(&privKey.PublicKey) // uncompressed [X || Y] (64 bytes)
		issuerPkBytes = pubKey
	} else if version == MultiShow {
		issuerSk := eddsa.PrivateKey{}
		_, err = issuerSk.SetBytes(issuerPrivateKey)
		if err != nil {
			return nil, err
		}
		issuerPkBytes = issuerSk.PublicKey.Bytes()
	}

	return &InternalCredential{
		ID:         id,
		Revoked:    false,
		VrfKeyPair: vrfKeyPair,
		Credential: Credential{
			PublicKeyVrfHash: vrfKeyPair.PublicKeyVrfHash,
			Signature:        credSignature,
			Type:             version,
		},
		IssuerPublicKey: issuerPkBytes,
	}, nil
}

// GenRevocationToken generates a revocation token and its proof based on a given unix epoch and credential type.
// It supports OneShow and MultiShow credential types. Errors if the type is unknown or token generation fails.
func (ic *InternalCredential) GenRevocationToken(unixEpoch int64) (token RevocationToken, proof []byte, error error) {
	epoch := make([]byte, 8)
	binary.BigEndian.PutUint64(epoch, uint64(unixEpoch))

	switch ic.Credential.Type {
	case OneShow:
		vrf := ecvrf.Secp256k1Sha256Tai
		ecdsaKey := secp256k1.PrivKeyFromBytes(ic.VrfKeyPair.PrivateKey).ToECDSA()

		t, p, err := vrf.Prove(ecdsaKey, epoch)
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
		_, err = hf.Write(ic.VrfKeyPair.PrivateKey)
		if err != nil {
			return nil, nil, err
		}
		return hf.Sum(nil), nil, nil
	default:
		return nil, nil, errors.New("unknown credential type")
	}
}

func hashPublicKeyCoordinates(x, y *big.Int) ([]byte, error) {
	// MiMC hash for zk-friendly applications
	mod := ecc.BN254.ScalarField()
	xCopy := new(big.Int).Mod(x, mod)
	yCopy := new(big.Int).Mod(y, mod)

	hfunc := mimc.NewMiMC()
	_, err := hfunc.Write(xCopy.Bytes())
	if err != nil {
		return nil, err
	}
	_, err = hfunc.Write(yCopy.Bytes())
	if err != nil {
		return nil, err
	}
	return hfunc.Sum(nil), nil
}

func signAttribute(issuerKey []byte, msg []byte, version CredentialType) ([]byte, error) {
	switch version {
	case OneShow:
		privKey, err := crypto.ToECDSA(issuerKey)
		if err != nil {
			return nil, err
		}

		sig, err := crypto.Sign(msg, privKey) // 65 bytes: R || S || V
		if err != nil {
			return nil, err
		}

		if sig[64] < 27 {
			sig[64] += 27
		}
		return sig, nil
	case MultiShow:
		key := eddsa.PrivateKey{}
		_, err := key.SetBytes(issuerKey)
		if err != nil {
			return nil, err
		}

		sigBytes, err := key.Sign(msg, mimc.NewMiMC())
		if err != nil {
			return nil, err
		}

		return sigBytes, nil
	default:
		return nil, errors.New("unknown credential type")
	}

}
