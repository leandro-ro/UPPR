package bloom

import (
	"fmt"
	"math"
)

type BloomFilterCascade struct {
	filters  []*BloomFilter
	capacity int
	fprate   float64
}

func NewCascade(capacity int, fprate float64) *BloomFilterCascade {
	// Optimal number of bits for the first layer
	mFloat := -1 * float64(capacity) * math.Log(fprate) / (math.Ln2 * math.Ln2)
	m := uint(math.Ceil(mFloat))

	// Optimal number of hash functions
	kFloat := (float64(m) / float64(capacity)) * math.Ln2
	k := uint(math.Round(kFloat))

	var filters []*BloomFilter
	filters = make([]*BloomFilter, 1)
	filters[0] = NewBloomFilter(m, k)

	return &BloomFilterCascade{filters, capacity, fprate}
}

func (c *BloomFilterCascade) Update(positives [][]byte, negatives [][]byte) error {
	c.reset()

	if len(positives) > c.capacity {
		return fmt.Errorf("bloom filter capacity exceeded")
	}

	// Layer 0: insert actual positives
	for _, p := range positives {
		c.addToLayer(0, p)
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
	for _, fp := range falsePositives {
		c.addToLayer(1, fp)
	}

	prevFalsePositives := negatives // negatives from the start
	layer := 2

	for len(falsePositives) > 0 {
		for _, fp := range falsePositives {
			c.addToLayer(layer, fp)
		}

		var nextFalsePositives [][]byte
		for _, prev := range prevFalsePositives {
			if c.filters[layer].Test(prev) {
				nextFalsePositives = append(nextFalsePositives, prev)
			}
		}

		prevFalsePositives = falsePositives
		falsePositives = nextFalsePositives
		layer++
	}

	return nil
}

func (c *BloomFilterCascade) addToLayer(layer int, element []byte) error {
	if layer < 0 {
		return fmt.Errorf("layer must be >= 0")
	}

	if layer == len(c.filters) {
		// Add one new layer based on layer 0 parameters
		first := c.filters[0] // TODO: Fix parameter choice
		m := first.m
		k := 1
		c.filters = append(c.filters, NewBloomFilter(m, k))
	} else if layer > len(c.filters) {
		return fmt.Errorf("cannot add to layer %d: previous layers missing", layer)
	}

	c.filters[layer].Add(element)
	return nil
}

func (c *BloomFilterCascade) reset() {
	c.filters = NewCascade(c.capacity, c.fprate).filters
}
