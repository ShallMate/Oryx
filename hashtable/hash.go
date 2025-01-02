package hashtable

import (
	"encoding/binary"
	"math/bits"
)

const (
	xxPrime64_1 = 11400714785074694791
	xxPrime64_2 = 14029467366897019727
	xxPrime64_3 = 1609587929392839161
	xxPrime64_4 = 9650029242287828579
	xxPrime64_5 = 2870177450012600261
)

func u64(buf []byte) uint64 {
	return binary.LittleEndian.Uint64(buf)
}

func u32(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf)
}

func rol64(u uint64, k int) uint64 { return bits.RotateLeft64(u, k) }

// XX64 returns the 64bit xxHash value.
func XX64(data []byte, seed uint64) uint64 {
	n := len(data)
	if n == 0 {
		return 0
	}

	var h64 uint64
	if n >= 32 {
		v1 := seed + xxPrime64_1 + xxPrime64_2
		v2 := seed + xxPrime64_2
		v3 := seed
		v4 := seed - xxPrime64_1
		p := 0
		for n := n - 32; p <= n; p += 32 {
			sub := data[p:][:32] //BCE hint for compiler
			v1 = rol64(v1+u64(sub[:])*xxPrime64_2, 31) * xxPrime64_1
			v2 = rol64(v2+u64(sub[8:])*xxPrime64_2, 31) * xxPrime64_1
			v3 = rol64(v3+u64(sub[16:])*xxPrime64_2, 31) * xxPrime64_1
			v4 = rol64(v4+u64(sub[24:])*xxPrime64_2, 31) * xxPrime64_1
		}

		h64 = rol64(v1, 1) + rol64(v2, 7) + rol64(v3, 12) + rol64(v4, 18)

		v1 *= xxPrime64_2
		v2 *= xxPrime64_2
		v3 *= xxPrime64_2
		v4 *= xxPrime64_2

		h64 = (h64^(rol64(v1, 31)*xxPrime64_1))*xxPrime64_1 + xxPrime64_4
		h64 = (h64^(rol64(v2, 31)*xxPrime64_1))*xxPrime64_1 + xxPrime64_4
		h64 = (h64^(rol64(v3, 31)*xxPrime64_1))*xxPrime64_1 + xxPrime64_4
		h64 = (h64^(rol64(v4, 31)*xxPrime64_1))*xxPrime64_1 + xxPrime64_4

		h64 += uint64(n)

		data = data[p:]
		n -= p
	} else {
		h64 = seed + xxPrime64_5 + uint64(n)
	}

	p := 0
	for n := n - 8; p <= n; p += 8 {
		sub := data[p : p+8]
		h64 ^= rol64(u64(sub)*xxPrime64_2, 31) * xxPrime64_1
		h64 = rol64(h64, 27)*xxPrime64_1 + xxPrime64_4
	}
	if p+4 <= n {
		sub := data[p : p+4]
		h64 ^= uint64(u32(sub)) * xxPrime64_1
		h64 = rol64(h64, 23)*xxPrime64_2 + xxPrime64_3
		p += 4
	}
	for ; p < n; p++ {
		h64 ^= uint64(data[p]) * xxPrime64_5
		h64 = rol64(h64, 11) * xxPrime64_1
	}

	h64 ^= h64 >> 33
	h64 *= xxPrime64_2
	h64 ^= h64 >> 29
	h64 *= xxPrime64_3
	h64 ^= h64 >> 32

	return h64
}
