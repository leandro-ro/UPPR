package issuer

import (
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCredential_OneShow(t *testing.T) {
	issuerKey := NewIssuer(OneShow).key

	cred, err := NewInternalCredential(OneShow, uint(1), issuerKey)
	require.NoError(t, err)
	require.False(t, cred.Revoked)
	require.Equal(t, cred.Type, OneShow)
	require.Equal(t, cred.ID, uint(1))
	require.NotNil(t, cred.VrfKeyPair.PrivateKey)
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
	require.NotNil(t, cred.VrfKeyPair.PrivateKey)
	require.NotNil(t, cred.Credential.PublicKeyVrfHash)

	sigValid, err := issuerKey.PublicKey.Verify(cred.Credential.Signature.Bytes(), cred.Credential.PublicKeyVrfHash, mimc.NewMiMC())
	require.NoError(t, err)
	require.True(t, sigValid)

	sigValid, err = cred.Credential.Verify(issuerKey.PublicKey)
	require.NoError(t, err)
	require.True(t, sigValid)
}
