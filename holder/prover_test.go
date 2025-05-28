package holder

import (
	"PrivacyPreservingRevocationCode/issuer"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewRevocationTokenProver(t *testing.T) {
	prover, err := NewRevocationTokenProver()
	require.NoError(t, err)

	iss := issuer.NewIssuer(issuer.MultiShow)
	err = iss.IssueCredential(0)
	require.NoError(t, err)
	cred, err := iss.GetCredentialCopy(0)
	require.NoError(t, err)

	epochUnix := time.Now().UTC().Unix()
	proof, witness, err := prover.GenProof(cred, epochUnix)
	require.NoError(t, err)
	require.NotNil(t, proof)
	require.NotNil(t, witness)

	publicWitness, err := witness.Public()
	require.NoError(t, err)
	require.NotNil(t, publicWitness)

	err = prover.VerifyProof(proof, publicWitness)
	require.NoError(t, err)
}
