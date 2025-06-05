package issuer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCredential_OneShow(t *testing.T) {
	issuer := NewIssuer(OneShow)

	cred, err := NewInternalCredential(OneShow, uint(1), issuer.key)
	require.NoError(t, err)
	require.False(t, cred.Revoked)
	require.Equal(t, cred.Credential.Type, OneShow)
	require.Equal(t, cred.ID, uint(1))
	require.NotNil(t, cred.VrfKeyPair.PrivateKey)
	require.NotNil(t, cred.Credential.PublicKeyVrfHash)

	sigValid, err := issuer.VerifySig(cred.Credential.PublicKeyVrfHash, cred.Credential.Signature)
	require.NoError(t, err)
	require.True(t, sigValid)

	sigValid, err = cred.Credential.Verify(issuer.GetPublicKey())
	require.NoError(t, err)
	require.True(t, sigValid)
}

func TestCredential_MultiShow(t *testing.T) {
	issuer := NewIssuer(MultiShow)

	cred, err := NewInternalCredential(MultiShow, uint(1), issuer.key)
	require.NoError(t, err)
	require.False(t, cred.Revoked)
	require.Equal(t, cred.Credential.Type, MultiShow)
	require.Equal(t, cred.ID, uint(1))
	require.NotNil(t, cred.VrfKeyPair.PrivateKey)
	require.NotNil(t, cred.Credential.PublicKeyVrfHash)

	sigValid, err := issuer.VerifySig(cred.Credential.PublicKeyVrfHash, cred.Credential.Signature)
	require.NoError(t, err)
	require.True(t, sigValid)

	sigValid, err = cred.Credential.Verify(issuer.GetPublicKey())
	require.NoError(t, err)
	require.True(t, sigValid)
}
