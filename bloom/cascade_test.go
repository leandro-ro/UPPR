package bloom

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCascade_Update(t *testing.T) {
	domain := 100_000
	capacity := 10_000

	cascade := NewCascade(domain, capacity)
	require.Equal(t, capacity, cascade.capacity)

	valid, revoked := genRevocationTokens(domain, capacity)

	// Initial update: populate with revoked tokens
	err := cascade.Update(revoked, valid)
	require.NoError(t, err, "Failed to update cascade with revoked tokens")
	require.Greater(t, len(cascade.filters), 5, "Cascade should contain multiple layers after initial update")

	// Update with no revocations: should reset to one empty filter
	err = cascade.Update(nil, valid)
	require.NoError(t, err, "Update with no revocations should not fail")
	require.Len(t, cascade.filters, 1, "Cascade should reset to one filter if no positives are provided")
	require.Equal(t, uint(0), cascade.filters[0].BitSet().Count(), "Reset filter should be empty")
}

func TestCascade_Check(t *testing.T) {
	domain := 100_000
	capacity := 1_000

	cascade := NewCascade(domain, capacity)
	valid, revoked := genRevocationTokens(domain, capacity)

	err := cascade.Update(revoked, valid)
	require.NoError(t, err, "Failed to update cascade with test data")

	// Test valid tokens: must not be detected as revoked
	for i, element := range valid {
		ok, layer := cascade.Test(element)
		require.Falsef(t, ok, "Valid element falsely identified as revoked at layer %d (index %d, data %s)", layer, i, hex.EncodeToString(element))
	}

	// Test revoked tokens: must be detected as revoked
	for i, element := range revoked {
		ok, layer := cascade.Test(element)
		require.Truef(t, ok, "Revoked element not identified as revoked (index %d, data %s, last matched layer %d)", i, hex.EncodeToString(element), layer)
	}
}

func TestCascade_EmptyDomain(t *testing.T) {
	cascade := NewCascade(1000, 0)

	err := cascade.Update(nil, nil)
	require.NoError(t, err)
	require.Len(t, cascade.filters, 1)
	require.Equal(t, uint(0), cascade.filters[0].BitSet().Count())
}

func TestCascade_SingleElement(t *testing.T) {
	token := generateRandom128BitSlices(1)[0]

	cascade := NewCascade(10, 1)
	err := cascade.Update([][]byte{token}, nil)
	require.NoError(t, err)

	ok, layer := cascade.Test(token)
	require.Truef(t, ok, "Single revoked element not detected, last layer %d", layer)
}

func TestCascade_StabilityAcrossUpdates(t *testing.T) {
	domain := 1000
	capacity := 100

	tokens1 := generateRandom128BitSlices(capacity)
	tokens2 := generateRandom128BitSlices(capacity)

	cascade := NewCascade(domain, capacity)
	require.NoError(t, cascade.Update(tokens1, nil))

	for _, tok := range tokens1 {
		ok, _ := cascade.Test(tok)
		require.True(t, ok)
	}

	// New update with entirely new tokens
	require.NoError(t, cascade.Update(tokens2, nil))

	for _, tok := range tokens2 {
		ok, _ := cascade.Test(tok)
		require.True(t, ok)
	}
}

func TestCascade_LayerCountSanity(t *testing.T) {
	domain := 10_000
	capacity := 2_000

	valid, revoked := genRevocationTokens(domain, capacity)
	cascade := NewCascade(domain, capacity)

	require.NoError(t, cascade.Update(revoked, valid))
	require.Greater(t, len(cascade.filters), 1)
	require.Less(t, len(cascade.filters), 100, "Too many layers: possible non-converging cascade")
}

func BenchmarkCascadeGeneration(b *testing.B) {
	domainSizes := []int{50_000, 100_000, 200_000, 300_000, 400_000, 500_000, 600_000, 700_000, 800_000, 900_000, 1_000_000}
	revocationRates := []float64{0.05, 0.1}

	for _, domain := range domainSizes {
		for _, rate := range revocationRates {
			name := fmt.Sprintf("Domain_%d_Rate_%.2f", domain, rate)
			b.Run(name, func(b *testing.B) {
				capacity := int(float64(domain) * rate)
				valid := generateRandom128BitSlices(domain - capacity)
				revoked := generateRandom128BitSlices(capacity)

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					cascade := getCascadeFromRate(domain, rate)
					if err := cascade.Update(revoked, valid); err != nil {
						b.Fatalf("Update failed: %v", err)
					}
				}
			})
		}
	}
}

func getCascadeFromRate(domain int, maxRevocationRate float64) *BloomFilterCascade {
	capacity := int(float64(domain) * maxRevocationRate)
	return NewCascade(domain, capacity)
}

// genRevocationTokens generates random 128-bit tokens and splits them into valid and revoked sets.
func genRevocationTokens(domain, revocCapacity int) (valid, revoked [][]byte) {
	valid = generateRandom128BitSlices(domain - revocCapacity)
	revoked = generateRandom128BitSlices(revocCapacity)
	return
}

// generateRandom128BitSlices returns 'count' number of random 128-bit (16-byte) slices.
func generateRandom128BitSlices(count int) [][]byte {
	result := make([][]byte, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 16)
		_, err := rand.Read(b)
		if err != nil {
			panic(fmt.Sprintf("failed to generate random bytes: %v", err))
		}
		result[i] = b
	}
	return result
}
