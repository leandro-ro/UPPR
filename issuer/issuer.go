package issuer

import (
	"PrivacyPreservingRevocationCode/bloom"
	crand "crypto/rand"
	"errors"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	mrand "math/rand"
	"runtime"
	"sync"
	"time"
)

// Issuer maintains issued and revoked credentials.
type Issuer struct {
	key                []byte                       // key is the issuer key pair
	credentialType     CredentialType               // credentialType represents the specific category of CredentialType managed by the issuer.
	issuedCredentials  map[uint]*InternalCredential // issuedCredentials holds all issued credentials (including revoked)
	revokedCredentials map[uint]bool                // revokedCredentials holds the uint ids of revoked creds in issuedCredentials
}

// NewIssuer creates a new Issuer with a generated key appropriate to the credential type.
func NewIssuer(credentialType CredentialType) *Issuer {
	var key []byte
	switch credentialType {
	case OneShow:
		privKey, err := ethcrypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		key = ethcrypto.FromECDSA(privKey) // returns 32-byte secp256k1 private key
	case MultiShow:
		eddsaKey, err := eddsa.GenerateKey(crand.Reader)
		if err != nil {
			panic(err)
		}
		key = eddsaKey.Bytes()
	default:
		panic("unknown credential type")
	}

	return &Issuer{
		key:                key,
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

// GetAllValidCreds returns a slice of all non-revoked credentials.
func (i *Issuer) GetAllValidCreds() []*InternalCredential {
	valid := make([]*InternalCredential, 0)
	for _, cred := range i.issuedCredentials {
		if !cred.Revoked {
			valid = append(valid, cred)
		}
	}
	return valid
}

// GetAllRevokedCreds returns a slice of all revoked credentials.
func (i *Issuer) GetAllRevokedCreds() []*InternalCredential {
	revoked := make([]*InternalCredential, 0)
	for _, cred := range i.issuedCredentials {
		if cred.Revoked {
			revoked = append(revoked, cred)
		}
	}
	return revoked
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

// GetPublicKey returns the issuer's public key (compressed secp256k1 or eddsa encoded).
func (i *Issuer) GetPublicKey() []byte {
	switch i.credentialType {
	case OneShow:
		privKey, err := ethcrypto.ToECDSA(i.key)
		if err != nil {
			return nil
		}
		pubKey := privKey.PublicKey
		return ethcrypto.CompressPubkey(&pubKey)
	case MultiShow:
		var eddsaKey eddsa.PrivateKey
		_, err := eddsaKey.SetBytes(i.key)
		if err != nil {
			return nil
		}
		return eddsaKey.PublicKey.Bytes()
	default:
		return nil
	}
}

func (i *Issuer) GetPrivateKey() []byte {
	return i.key
}

func (i *Issuer) VerifySig(msg []byte, sig []byte) (bool, error) {
	switch i.credentialType {
	case OneShow:
		privKey, err := ethcrypto.ToECDSA(i.key)
		if err != nil {
			return false, err
		}
		pubKey := &privKey.PublicKey
		if len(sig) != 65 {
			return false, errors.New("expected 65-byte signature (r||s||v)")
		}
		return ethcrypto.VerifySignature(
			ethcrypto.FromECDSAPub(pubKey),
			msg,      // must be a 32-byte hash
			sig[:64], // r||s only
		), nil

	case MultiShow:
		pk := eddsa.PublicKey{}
		_, err := pk.SetBytes(i.key)
		if err != nil {
			return false, err
		}
		return pk.Verify(sig, msg, mimc.NewMiMC())

	default:
		return false, errors.New("unsupported credential type")
	}
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
