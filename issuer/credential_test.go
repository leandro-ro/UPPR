package issuer

import (
	"encoding/binary"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/klayoracle/go-ecvrf"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCredential_OneShow(t *testing.T) {
	issuerKey := NewIssuer(OneShow).key

	cred, err := NewInternalCredential(OneShow, uint(1), issuerKey)
	require.NoError(t, err)
	require.False(t, cred.Revoked)
	require.Equal(t, cred.Type, OneShow)
	require.Equal(t, cred.ID, uint(1))
	require.NotNil(t, cred.PrivateKeyVrf)
	require.NotNil(t, cred.Credential.PublicKeyVrfHash)

	sigValid, err := issuerKey.PublicKey.Verify(cred.Credential.Signature.Bytes(), cred.Credential.PublicKeyVrfHash, mimc.NewMiMC())
	require.NoError(t, err)
	require.True(t, sigValid)

	sigValid, err = cred.Credential.Verify(issuerKey.PublicKey)
	require.NoError(t, err)
	require.True(t, sigValid)
}

func TestCredential_MultiShow(t *testing.T) {
	issuerKey := NewIssuer(MultiShow).key

	cred, err := NewInternalCredential(MultiShow, uint(1), issuerKey)
	require.NoError(t, err)
	require.False(t, cred.Revoked)
	require.Equal(t, cred.Type, MultiShow)
	require.Equal(t, cred.ID, uint(1))
	require.NotNil(t, cred.PrivateKeyVrf)
	require.NotNil(t, cred.Credential.PublicKeyVrfHash)

	sigValid, err := issuerKey.PublicKey.Verify(cred.Credential.Signature.Bytes(), cred.Credential.PublicKeyVrfHash, mimc.NewMiMC())
	require.NoError(t, err)
	require.True(t, sigValid)

	sigValid, err = cred.Credential.Verify(issuerKey.PublicKey)
	require.NoError(t, err)
	require.True(t, sigValid)
}

func TestCredential_RevocationTokenOneShow(t *testing.T) {
	issuerKey := NewIssuer(OneShow)
	cred, err := NewInternalCredential(OneShow, uint(1), issuerKey.key)
	require.NoError(t, err)

	epochUnix := time.Now().UTC().Unix()
	token, proof, err := cred.GenRevocationToken(epochUnix)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotNil(t, proof)

	vrf := ecvrf.Secp256k1Sha256Tai

	epoch := make([]byte, 8)
	binary.BigEndian.PutUint64(epoch, uint64(epochUnix))

	vrfPk := cred.GetVrfPublicKey()
	expectedToken, err := vrf.VerifySecp256k1(vrfPk, epoch, proof)
	require.NoError(t, err)
	require.Equal(t, expectedToken, token)
}
