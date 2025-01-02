package hashtable

import (
	"fmt"
	"math/big"
	"math/rand"
)

type SimpleHash struct {
	// The number of buckets in the hash table
	NumBuckets int
	Elements   [][]*big.Int
	BucketSize int
	Hashnum    int
}

func NewSimpleHash(size int, bucketSize int, hashnum int) *SimpleHash {
	e := 1.3
	numBuckets := int(float64(size) * e)
	table := make([][]*big.Int, numBuckets)
	for i := range table {
		table[i] = make([]*big.Int, 0, bucketSize)
	}
	return &SimpleHash{
		NumBuckets: numBuckets,
		BucketSize: bucketSize,
		Elements:   table,
		Hashnum:    hashnum,
	}
}

func (h *SimpleHash) Insert(key *big.Int) {
	for i := 0; i < h.Hashnum; i++ {
		index := XX64(key.Bytes(), uint64(i)) % uint64(h.NumBuckets)
		h.Elements[index] = append(h.Elements[index], key)
	}
}

func (h *SimpleHash) Find(key *big.Int) bool {
	for i := 0; i < h.Hashnum; i++ {
		index := XX64(key.Bytes(), uint64(i)) % uint64(h.NumBuckets)
		for _, value := range h.Elements[index] {
			if value.Cmp(key) == 0 {
				return true
			}
		}
	}
	return false
}

func (h *SimpleHash) FillBuckets() {
	maxBucketSize := 0
	for _, bucket := range h.Elements {
		if len(bucket) > maxBucketSize {
			maxBucketSize = len(bucket)
		}
	}
	fmt.Println(maxBucketSize)
	h.BucketSize = maxBucketSize
	for i := 0; i < h.NumBuckets; i++ {
		for len(h.Elements[i]) < maxBucketSize {
			randomValue := big.NewInt(rand.Int63())
			h.Elements[i] = append(h.Elements[i], randomValue)
		}
	}
}
