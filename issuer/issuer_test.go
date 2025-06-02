package issuer

import (
	"fmt"
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

func TestIssuer_GenRevocationTokensOneShow(t *testing.T) {
	issuer := NewIssuer(OneShow)

	err := issuer.IssueCredentials(1000)
	require.NoError(t, err)

	err = issuer.RevokeRandomCredentials(20)
	require.NoError(t, err)

	revokedTokens, validTokens, _, err := issuer.genRevocationTokens()
	require.NoError(t, err)
	require.Equal(t, 20, len(revokedTokens))
	require.Equal(t, 980, len(validTokens))
}

func TestIssuer_GenRevocationTokensMultiShow(t *testing.T) {
	issuer := NewIssuer(MultiShow)

	err := issuer.IssueCredentials(1000)
	require.NoError(t, err)

	err = issuer.RevokeRandomCredentials(20)
	require.NoError(t, err)

	revokedTokens, validTokens, _, err := issuer.genRevocationTokens()
	require.NoError(t, err)
	require.Equal(t, 20, len(revokedTokens))
	require.Equal(t, 980, len(validTokens))
}

func TestIssuer_GenRevocationArtifactOneShow(t *testing.T) {
	issuer := NewIssuer(OneShow)
	err := issuer.IssueCredentials(1000)
	require.NoError(t, err)
	err = issuer.RevokeRandomCredentials(20)

	filter, revokedTokens, validTokens, _, err := issuer.GenRevocationArtifact()
	require.NoError(t, err)

	for _, token := range revokedTokens {
		b, _ := filter.Test(token.ToBytes())
		require.True(t, b)
	}

	for _, token := range validTokens {
		b, _ := filter.Test(token.ToBytes())
		require.False(t, b)
	}
}

func TestIssuer_GenRevocationArtifactMultiShow(t *testing.T) {
	issuer := NewIssuer(MultiShow)
	err := issuer.IssueCredentials(1000)
	require.NoError(t, err)
	err = issuer.RevokeRandomCredentials(20)

	filter, revokedTokens, validTokens, _, err := issuer.GenRevocationArtifact()
	require.NoError(t, err)

	for _, token := range revokedTokens {
		b, _ := filter.Test(token.ToBytes())
		require.True(t, b)
	}

	for _, token := range validTokens {
		b, _ := filter.Test(token.ToBytes())
		require.False(t, b)
	}
}

func BenchmarkIssuer_GenRevocationArtifact(b *testing.B) {
	domains := []int{10_000, 100_000, 1_000_000}
	rates := []float64{0.10, 0.05, 0.01}
	modes := []CredentialType{OneShow, MultiShow}

	for _, mode := range modes {
		for _, domain := range domains {
			for _, rate := range rates {
				issued := domain
				revoked := int(float64(domain) * rate)

				name := fmt.Sprintf("%s_Domain_%d_Rate_%.2f", mode.String(), domain, rate)
				b.Run(name, func(b *testing.B) {
					issuer := NewIssuer(mode)

					err := issuer.IssueCredentials(uint(issued))
					if err != nil {
						b.Fatalf("IssueCredentials failed: %v", err)
					}

					err = issuer.RevokeRandomCredentials(uint(revoked))
					if err != nil {
						b.Fatalf("RevokeRandomCredentials failed: %v", err)
					}

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_, _, _, _, err := issuer.GenRevocationArtifact() // includes gen of revocation tokens + cascade generation
						if err != nil {
							b.Fatalf("GenRevocationArtifact failed: %v", err)
						}
					}
					b.StopTimer()
				})
			}
		}
	}
}
