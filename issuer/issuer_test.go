package issuer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIssuer_IssueAndRevokeCredentials(t *testing.T) {
	issuer := NewIssuer(OneShow)

	// Issue 10 credentials
	err := issuer.IssueCredentials(100)
	require.NoError(t, err)

	require.Equal(t, 100, issuer.AmountIssued())

	err = issuer.RevokeRandomCredentials(40)
	require.NoError(t, err)
	require.Equal(t, 100, issuer.AmountIssued())
	require.Equal(t, 40, issuer.AmountRevoked())

	err = issuer.RevokeRandomCredentials(20)
	require.NoError(t, err)
	require.Equal(t, 100, issuer.AmountIssued())
	require.Equal(t, 60, issuer.AmountRevoked())
}

func TestIssuer_GetRevocationStatus(t *testing.T) {
	issuer := NewIssuer(OneShow)
	id := uint(123456)

	err := issuer.IssueCredential(id)
	require.NoError(t, err)

	// Should not be revoked yet
	require.False(t, issuer.GetRevocationStatus(id))

	// Revoke it
	err = issuer.RevokeCredential(id)
	require.NoError(t, err)

	require.True(t, issuer.GetRevocationStatus(id))
	require.Equal(t, 1, issuer.AmountRevoked())

	cred, err := issuer.GetCredentialCopy(id)
	require.NoError(t, err)
	require.Equal(t, id, cred.ID)
	require.True(t, cred.Revoked)
}
