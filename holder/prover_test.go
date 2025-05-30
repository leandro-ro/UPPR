package holder

import (
	"PrivacyPreservingRevocationCode/issuer"
	"encoding/binary"
	"github.com/klayoracle/go-ecvrf"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestProver_OneShow(t *testing.T) {
	iss := issuer.NewIssuer(issuer.OneShow)
	err := iss.IssueCredential(0)
	require.NoError(t, err)
	cred, err := iss.GetCredentialCopy(0)
	require.NoError(t, err)

	epochUnix := time.Now().UTC().Unix()
	token, proof, err := cred.GenRevocationToken(epochUnix)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotNil(t, proof)

	vrf := ecvrf.Secp256k1Sha256Tai

	epoch := make([]byte, 8)
	binary.BigEndian.PutUint64(epoch, uint64(epochUnix))

	vrfPk, err := cred.VrfKeyPair.GetOneShowPublicKey()
	require.NoError(t, err)
	expectedToken, err := vrf.Verify(vrfPk, epoch, proof)
	require.NoError(t, err)
	require.Equal(t, issuer.RevocationToken(expectedToken), token)
}

func TestProver_MultiShow(t *testing.T) {
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
