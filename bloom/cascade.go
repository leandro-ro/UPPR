package bloom

import (
	"fmt"
	"math"
)

// BloomFilterCascade represents a cascade of Bloom filters.
// It is constructed by iteratively filtering false positives from prior layers.
type BloomFilterCascade struct {
	filters          []*BloomFilter // filters holds a slice of BloomFilter layers in the Bloom filter cascade.
	capacity         int            // capacity defines the maximum number of elements expected to be processed by the Bloom filter cascade.
	falsePosRate     float64        // falsePosRate represents the false positive rate for the first layer.
	falsePosRateSucc float64        // falsePosRateSucc represents the false positive rate for subsequent layers.
}

// NewCascade creates a new BloomFilterCascade with an initial layer based on the given domain and capacity.
// The false positive rate for the first layer is computed based on the ratio of capacity to domain size.
func NewCascade(domain, capacity int) *BloomFilterCascade {
	falsePosRate := float64(capacity) * math.Sqrt(0.5) / float64(domain-capacity)
	falsePosRateSucc := 0.5

	m, k := getOptimalFilterParameters(capacity, falsePosRate)

	filters := make([]*BloomFilter, 1)
	filters[0] = NewBloomFilter(m, k)

	return &BloomFilterCascade{filters, capacity, falsePosRate, falsePosRateSucc}
}

// Update constructs the cascade from a set of true positives and known negatives.
// Each layer stores false positives of the previous layer, alternating between accepting and rejecting layers.
// Terminates once no new false positives are found or a maximum depth is reached.
func (c *BloomFilterCascade) Update(positives [][]byte, negatives [][]byte) error {
	c.reset()

	if len(positives) > c.capacity {
		return fmt.Errorf("bloom filter capacity exceeded")
	}

	// Layer 0: insert actual positives
	for _, p := range positives {
		c.filters[0].Add(p)
	}

	// Find false positives at layer 0
	var falsePositives [][]byte
	for _, n := range negatives {
		if c.filters[0].Test(n) {
			falsePositives = append(falsePositives, n)
		}
	}

	if len(falsePositives) == 0 {
		return nil
	}

	// Layer 1: insert false positives
	c.addNextLayer(&falsePositives, c.falsePosRateSucc)

	// Generate succeeding layers
	prevPrevFalsePositives := &positives
	prevFalsePositives := &falsePositives
	layer := 1

	for {
		var nextFalsePositives [][]byte
		for _, prev := range *prevPrevFalsePositives {
			if c.filters[layer].Test(prev) {
				nextFalsePositives = append(nextFalsePositives, prev)
			}
		}
		prevPrevFalsePositives = prevFalsePositives
		prevFalsePositives = &nextFalsePositives

		if len(nextFalsePositives) == 0 {
			return nil
		} else if len(nextFalsePositives) > 200 {
			c.addNextLayer(&nextFalsePositives, c.falsePosRateSucc)
		} else {
			c.addNextLayer(&nextFalsePositives, 0.1) // Ensure termination
		}

		layer++
		if layer > 100 {
			return fmt.Errorf("over 100 layers. bloom filter cascade too deep — probably cyclic data")
		}
	}
}

// Test determines whether the given element is accepted by the Bloom Filter Cascade.
//
// The cascade alternates between positive (even) and filtering (odd) layers.
// - Even-numbered layers (0, 2, ...) are expected to return true for accepted elements.
// - Odd-numbered layers (1, 3, ...) filter out false positives and are expected to return false.
//
// If the element mismatches the expected behavior at any non-final layer, it is either accepted or rejected early.
// The final layer determines the classification if no early termination occurs.
//
// Returns:
// - bool: whether the element is accepted
// - int: the index of the layer that determined the result
func (c *BloomFilterCascade) Test(element []byte) (bool, int) {
	last := len(c.filters) - 1

	for layer, f := range c.filters {
		match := f.Test(element)

		if layer == last {
			expected := layer%2 == 0
			return match == expected, layer
		}

		if !match {
			accept := layer%2 == 1
			return accept, layer
		}
	}

	panic("unreachable: all layers exhausted without return")
}

// addNextLayer adds a new Bloom filter layer to the cascade based on the provided elements and false positive rate.
// The filter stores only the elements passed in and is appended to the internal filter list.
func (c *BloomFilterCascade) addNextLayer(elements *[][]byte, fprate float64) {
	capacity := int(max(uint(len(*elements)), 100))
	m, k := getOptimalFilterParameters(capacity, fprate)

	nextLayer := NewBloomFilter(m, k)

	for _, element := range *elements {
		nextLayer.Add(element)
	}

	c.filters = append(c.filters, nextLayer)
}

// reset clears the cascade and reinitializes the first filter layer with original parameters.
func (c *BloomFilterCascade) reset() {
	m, k := getOptimalFilterParameters(c.capacity, c.falsePosRate)
	c.filters = []*BloomFilter{NewBloomFilter(m, k)}
}

// printStats prints the size and number of hash functions for each layer in the cascade.
// It also reports the total size in bits and cumulative number of hash functions.
func (c *BloomFilterCascade) printStats() {
	fmt.Println("Bloom Filter Cascade Statistics:")

	var totalSizeBits uint
	var totalHashFuncs uint

	for i, f := range c.filters {
		size := f.Cap()
		k := f.K()
		totalSizeBits += size
		totalHashFuncs += k
		fmt.Printf("Layer %d: size = %d bits, hash functions = %d\n", i, size, k)
	}

	fmt.Printf("Total: size = %d bits, total hash functions = %d\n", totalSizeBits, totalHashFuncs)
}

// getOptimalFilterParameters calculates optimal parameters (m: filter size in bits, k: number of hash functions)
// for a Bloom filter given the expected capacity and desired false positive rate.
func getOptimalFilterParameters(capacity int, fprate float64) (m, k uint) {
	mFloat := -1 * float64(capacity) * math.Log(fprate) / (math.Ln2 * math.Ln2)
	m = uint(math.Ceil(mFloat))

	kFloat := (float64(m) / float64(capacity)) * math.Ln2
	k = uint(math.Round(kFloat))
	return m, k
}
