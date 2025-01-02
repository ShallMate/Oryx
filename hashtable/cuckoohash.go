package hashtable

import (
	"fmt"
	"math/big"
)

// CuckooHashTable 是 Cuckoo 哈希表结构
type CuckooHashTable struct {
	table      []*big.Int
	size       int
	cuckoosize int
	maxLoops   int
}

// NewCuckooHashTable 创建一个新的 Cuckoo 哈希表
func NewCuckooHashTable(size int) *CuckooHashTable {
	e := 1.27
	cuckoosize := int(float64(size) * e)
	return &CuckooHashTable{
		table:      make([]*big.Int, cuckoosize),
		size:       size,
		cuckoosize: cuckoosize,
		maxLoops:   500,
	}
}

func (cht *CuckooHashTable) Insert(key []*big.Int) bool {
	Hash_index := make([]uint64, cht.cuckoosize)
	for i := 0; i < cht.size; i++ {
		var old_hash_id uint64 = 1
		j := 0
		x_key := key[i]
		MAXITER := cht.maxLoops
		for ; j < MAXITER; j++ {
			h := XX64(x_key.Bytes(), old_hash_id) % uint64(cht.cuckoosize)
			hash_id_address := &Hash_index[h]
			key_address := &cht.table[h]
			if *hash_id_address == 0 {
				*hash_id_address = old_hash_id
				*key_address = x_key
				break
			} else {
				old_hash_id, *hash_id_address = *hash_id_address, old_hash_id
				x_key, *key_address = *key_address, x_key
				old_hash_id = old_hash_id%3 + 1
			}
		}
		if j == MAXITER-1 {
			fmt.Println("insert failed, ", i)
			return false
		}
	}
	return true
}

func (cht *CuckooHashTable) Get(key *big.Int) (*big.Int, bool) {
	index1 := XX64(key.Bytes(), 1) % uint64(cht.cuckoosize)
	if cht.table[index1] != nil && cht.table[index1].Cmp(key) == 0 {
		return cht.table[index1], true
	}
	index2 := XX64(key.Bytes(), 2) % uint64(cht.cuckoosize)
	if cht.table[index2] != nil && cht.table[index2].Cmp(key) == 0 {
		return cht.table[index2], true
	}
	index3 := XX64(key.Bytes(), 3) % uint64(cht.cuckoosize)
	if cht.table[index3] != nil && cht.table[index3].Cmp(key) == 0 {
		return cht.table[index3], true
	}
	return nil, false
}

func (cht *CuckooHashTable) Iterate() {
	for _, value := range cht.table {
		if value != nil {
			fmt.Printf("Value: %s\n", value.String())
		}
	}
}

func (cht *CuckooHashTable) GetSize() (int, int) {
	return cht.size, cht.cuckoosize
}
