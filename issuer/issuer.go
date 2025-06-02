package issuer

import (
	"PrivacyPreservingRevocationCode/bloom"
	crand "crypto/rand"
	"errors"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	mrand "math/rand"
	"runtime"
	"sync"
	"time"
)

// Issuer maintains issued and revoked credentials.
type Issuer struct {
	key                eddsa.PrivateKey             // key is the issuer key pair
	credentialType     CredentialType               // credentialType represents the specific category of CredentialType managed by the issuer.
	issuedCredentials  map[uint]*InternalCredential // issuedCredentials holds all issued credentials (including revoked)
	revokedCredentials map[uint]bool                // revokedCredentials holds the uint ids of revoked creds in issuedCredentials
}

func NewIssuer(credentialType CredentialType) *Issuer {
	key, _ := eddsa.GenerateKey(crand.Reader)

	return &Issuer{
		key:                *key,
		credentialType:     credentialType,
		issuedCredentials:  make(map[uint]*InternalCredential),
		revokedCredentials: make(map[uint]bool),
	}
}

// IssueCredentials generates and stores a number of credentials.
func (i *Issuer) IssueCredentials(amount uint) error {
	for issued := uint(0); issued < amount; {
		id := uint(mrand.Uint32())

		if _, exists := i.issuedCredentials[id]; exists {
			continue // try another if collision (unlikely)
		}

		err := i.IssueCredential(id)
		if err != nil {
			return err
		}
		issued++
	}
	return nil
}

func (i *Issuer) GetCredentialCopy(id uint) (InternalCredential, error) {
	cred, ok := i.issuedCredentials[id]
	if !ok {
		return InternalCredential{}, errors.New("credential not found")
	}
	return *cred, nil
}

func (i *Issuer) IssueCredential(id uint) error {
	cred, ok := i.issuedCredentials[id]
	if ok {
		return errors.New("credential id already assigned")
	}
	cred, err := NewInternalCredential(i.credentialType, id, i.key)
	if err != nil {
		return err
	}
	i.issuedCredentials[id] = cred
	return nil
}

func (i *Issuer) RevokeRandomCredentials(amount uint) error {
	if int(amount) > i.AmountIssued()-i.AmountRevoked() {
		return errors.New("not enough unrevoked credentials")
	}

	// Collect all unrevoked IDs
	unrevoked := make([]uint, 0, i.AmountIssued())
	for id := range i.issuedCredentials {
		if !i.GetRevocationStatus(id) {
			unrevoked = append(unrevoked, id)
		}
	}

	// Shuffle unrevoked IDs
	mrand.Shuffle(len(unrevoked), func(a, b int) {
		unrevoked[a], unrevoked[b] = unrevoked[b], unrevoked[a]
	})

	// Revoke first `amount` credentials
	for j := 0; uint(j) < amount; j++ {
		if err := i.RevokeCredential(unrevoked[j]); err != nil {
			return err // should not happen but good to bubble up
		}
	}

	return nil
}

// RevokeCredential revokes a credential by its ID.
func (i *Issuer) RevokeCredential(id uint) error {
	cred, ok := i.issuedCredentials[id]
	if !ok {
		return errors.New("credential not found")
	}
	cred.Revoked = true
	i.revokedCredentials[id] = true
	return nil
}

func (i *Issuer) genRevocationTokens() (revoked, valid []RevocationToken, epoch int64, err error) {
	epoch = time.Now().UTC().Unix()

	type result struct {
		token   RevocationToken
		revoked bool
		err     error
	}

	numWorkers := runtime.NumCPU()
	jobs := make(chan *InternalCredential, len(i.issuedCredentials))
	results := make(chan result, len(i.issuedCredentials))

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for cred := range jobs {
				token, _, err := cred.GenRevocationToken(epoch)
				if err != nil {
					results <- result{err: err}
					continue
				}
				results <- result{token: token, revoked: cred.Revoked}
			}
		}()
	}

	for _, cred := range i.issuedCredentials {
		jobs <- cred
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	revoked = make([]RevocationToken, 0, i.AmountRevoked())
	valid = make([]RevocationToken, 0, i.AmountIssued()-i.AmountRevoked())

	for res := range results {
		if res.err != nil {
			return nil, nil, -1, res.err
		}
		if res.revoked {
			revoked = append(revoked, res.token)
		} else {
			valid = append(valid, res.token)
		}
	}

	return revoked, valid, epoch, nil
}

// GenRevocationArtifact calls genRevocationTokens() an returns a BloomFilterCascade
func (i *Issuer) GenRevocationArtifact() (artifact *bloom.BloomFilterCascade, revoked, valid []RevocationToken, epoch int64, error error) {
	revoked, valid, epoch, err := i.genRevocationTokens()
	if err != nil {
		return nil, nil, nil, -1, err
	}

	cascade := bloom.NewCascade(i.AmountIssued(), i.AmountRevoked())
	err = cascade.Update(RevocationTokensToByteSlices(revoked), RevocationTokensToByteSlices(valid))
	if err != nil {
		return nil, nil, nil, -1, err
	}
	return cascade, revoked, valid, epoch, nil
}

func (i *Issuer) GetRevocationStatus(id uint) bool {
	return i.revokedCredentials[id]
}

// GetPublicKey returns the issuer's public key.
func (i *Issuer) GetPublicKey() eddsa.PublicKey {
	return i.key.PublicKey
}

// AmountIssued returns the total number of issued credentials.
func (i *Issuer) AmountIssued() int {
	return len(i.issuedCredentials)
}

// AmountRevoked returns the number of revoked credentials.
func (i *Issuer) AmountRevoked() int {
	return len(i.revokedCredentials)
}

// CredentialType returns the type of credentials managed by the issuer.
func (i *Issuer) CredentialType() CredentialType {
	return i.credentialType
}
