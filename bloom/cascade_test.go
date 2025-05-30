package bloom

import (
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestCascade(t *testing.T) {
	cascade := NewCascade(1000, 0.01)

	pos := generateRandom128BitSlices(1000) // fill to max capacity
	neg := generateRandom128BitSlices(10000)

	err := cascade.Update(pos, neg)
	require.NoError(t, err)
}

func generateRandom128BitSlices(count int) [][]byte {
	result := make([][]byte, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 16) // 128 bits = 16 bytes
		_, err := rand.Read(b)
		if err != nil {
			log.Fatalf("failed to generate random bytes: %v", err)
		}
		result[i] = b
	}
	return result
}
