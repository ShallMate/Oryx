package ecc

import (
	"encoding/hex"
)

const (
	twoBitsMask   = 0x3
	fourBitsMask  = 0xf
	sixBitsMask   = 0x3f
	eightBitsMask = 0xff
)

const (
	fieldWords        = 10
	fieldBase         = 26
	fieldOverflowBits = 32 - fieldBase
	fieldBaseMask     = (1 << fieldBase) - 1
	fieldMSBBits      = 256 - (fieldBase * (fieldWords - 1))

	fieldMSBMask       = (1 << fieldMSBBits) - 1
	fieldPrimeWordZero = 0x3fffc2f
	fieldPrimeWordOne  = 0x3ffffbf
)

var (
	fieldQBytes = []byte{
		0x3f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xbf, 0xff, 0xff, 0x0c,
	}
)

type FieldVal struct {
	n [10]uint32
}

func (f FieldVal) String() string {
	t := new(FieldVal).Set(&f).Normalize()
	return hex.EncodeToString(t.Bytes()[:])
}

func (f *FieldVal) Zero() {
	f.n[0] = 0
	f.n[1] = 0
	f.n[2] = 0
	f.n[3] = 0
	f.n[4] = 0
	f.n[5] = 0
	f.n[6] = 0
	f.n[7] = 0
	f.n[8] = 0
	f.n[9] = 0
}

func (f *FieldVal) Set(val *FieldVal) *FieldVal {
	*f = *val
	return f
}

func (f *FieldVal) SetInt(ui uint) *FieldVal {
	f.Zero()
	f.n[0] = uint32(ui)
	return f
}

func (f *FieldVal) SetBytes(b *[32]byte) *FieldVal {
	f.n[0] = uint32(b[31]) | uint32(b[30])<<8 | uint32(b[29])<<16 |
		(uint32(b[28])&twoBitsMask)<<24
	f.n[1] = uint32(b[28])>>2 | uint32(b[27])<<6 | uint32(b[26])<<14 |
		(uint32(b[25])&fourBitsMask)<<22
	f.n[2] = uint32(b[25])>>4 | uint32(b[24])<<4 | uint32(b[23])<<12 |
		(uint32(b[22])&sixBitsMask)<<20
	f.n[3] = uint32(b[22])>>6 | uint32(b[21])<<2 | uint32(b[20])<<10 |
		uint32(b[19])<<18
	f.n[4] = uint32(b[18]) | uint32(b[17])<<8 | uint32(b[16])<<16 |
		(uint32(b[15])&twoBitsMask)<<24
	f.n[5] = uint32(b[15])>>2 | uint32(b[14])<<6 | uint32(b[13])<<14 |
		(uint32(b[12])&fourBitsMask)<<22
	f.n[6] = uint32(b[12])>>4 | uint32(b[11])<<4 | uint32(b[10])<<12 |
		(uint32(b[9])&sixBitsMask)<<20
	f.n[7] = uint32(b[9])>>6 | uint32(b[8])<<2 | uint32(b[7])<<10 |
		uint32(b[6])<<18
	f.n[8] = uint32(b[5]) | uint32(b[4])<<8 | uint32(b[3])<<16 |
		(uint32(b[2])&twoBitsMask)<<24
	f.n[9] = uint32(b[2])>>2 | uint32(b[1])<<6 | uint32(b[0])<<14
	return f
}

func (f *FieldVal) SetByteSlice(b []byte) *FieldVal {
	var b32 [32]byte
	if len(b) > 32 {
		b = b[:32]
	}
	copy(b32[32-len(b):], b)
	return f.SetBytes(&b32)
}

func (f *FieldVal) SetHex(hexString string) *FieldVal {
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}
	bytes, _ := hex.DecodeString(hexString)
	return f.SetByteSlice(bytes)
}

func (f *FieldVal) Normalize() *FieldVal {
	t9 := f.n[9]
	m := t9 >> fieldMSBBits
	t9 = t9 & fieldMSBMask
	t0 := f.n[0] + m*977
	t1 := (t0 >> fieldBase) + f.n[1] + (m << 6)
	t0 = t0 & fieldBaseMask
	t2 := (t1 >> fieldBase) + f.n[2]
	t1 = t1 & fieldBaseMask
	t3 := (t2 >> fieldBase) + f.n[3]
	t2 = t2 & fieldBaseMask
	t4 := (t3 >> fieldBase) + f.n[4]
	t3 = t3 & fieldBaseMask
	t5 := (t4 >> fieldBase) + f.n[5]
	t4 = t4 & fieldBaseMask
	t6 := (t5 >> fieldBase) + f.n[6]
	t5 = t5 & fieldBaseMask
	t7 := (t6 >> fieldBase) + f.n[7]
	t6 = t6 & fieldBaseMask
	t8 := (t7 >> fieldBase) + f.n[8]
	t7 = t7 & fieldBaseMask
	t9 = (t8 >> fieldBase) + t9
	t8 = t8 & fieldBaseMask
	m = 1
	if t9 == fieldMSBMask {
		m &= 1
	} else {
		m &= 0
	}
	if t2&t3&t4&t5&t6&t7&t8 == fieldBaseMask {
		m &= 1
	} else {
		m &= 0
	}
	if ((t0+977)>>fieldBase + t1 + 64) > fieldBaseMask {
		m &= 1
	} else {
		m &= 0
	}
	if t9>>fieldMSBBits != 0 {
		m |= 1
	} else {
		m |= 0
	}
	t0 = t0 + m*977
	t1 = (t0 >> fieldBase) + t1 + (m << 6)
	t0 = t0 & fieldBaseMask
	t2 = (t1 >> fieldBase) + t2
	t1 = t1 & fieldBaseMask
	t3 = (t2 >> fieldBase) + t3
	t2 = t2 & fieldBaseMask
	t4 = (t3 >> fieldBase) + t4
	t3 = t3 & fieldBaseMask
	t5 = (t4 >> fieldBase) + t5
	t4 = t4 & fieldBaseMask
	t6 = (t5 >> fieldBase) + t6
	t5 = t5 & fieldBaseMask
	t7 = (t6 >> fieldBase) + t7
	t6 = t6 & fieldBaseMask
	t8 = (t7 >> fieldBase) + t8
	t7 = t7 & fieldBaseMask
	t9 = (t8 >> fieldBase) + t9
	t8 = t8 & fieldBaseMask
	t9 = t9 & fieldMSBMask
	f.n[0] = t0
	f.n[1] = t1
	f.n[2] = t2
	f.n[3] = t3
	f.n[4] = t4
	f.n[5] = t5
	f.n[6] = t6
	f.n[7] = t7
	f.n[8] = t8
	f.n[9] = t9
	return f
}

func (f *FieldVal) PutBytes(b *[32]byte) {
	b[31] = byte(f.n[0] & eightBitsMask)
	b[30] = byte((f.n[0] >> 8) & eightBitsMask)
	b[29] = byte((f.n[0] >> 16) & eightBitsMask)
	b[28] = byte((f.n[0]>>24)&twoBitsMask | (f.n[1]&sixBitsMask)<<2)
	b[27] = byte((f.n[1] >> 6) & eightBitsMask)
	b[26] = byte((f.n[1] >> 14) & eightBitsMask)
	b[25] = byte((f.n[1]>>22)&fourBitsMask | (f.n[2]&fourBitsMask)<<4)
	b[24] = byte((f.n[2] >> 4) & eightBitsMask)
	b[23] = byte((f.n[2] >> 12) & eightBitsMask)
	b[22] = byte((f.n[2]>>20)&sixBitsMask | (f.n[3]&twoBitsMask)<<6)
	b[21] = byte((f.n[3] >> 2) & eightBitsMask)
	b[20] = byte((f.n[3] >> 10) & eightBitsMask)
	b[19] = byte((f.n[3] >> 18) & eightBitsMask)
	b[18] = byte(f.n[4] & eightBitsMask)
	b[17] = byte((f.n[4] >> 8) & eightBitsMask)
	b[16] = byte((f.n[4] >> 16) & eightBitsMask)
	b[15] = byte((f.n[4]>>24)&twoBitsMask | (f.n[5]&sixBitsMask)<<2)
	b[14] = byte((f.n[5] >> 6) & eightBitsMask)
	b[13] = byte((f.n[5] >> 14) & eightBitsMask)
	b[12] = byte((f.n[5]>>22)&fourBitsMask | (f.n[6]&fourBitsMask)<<4)
	b[11] = byte((f.n[6] >> 4) & eightBitsMask)
	b[10] = byte((f.n[6] >> 12) & eightBitsMask)
	b[9] = byte((f.n[6]>>20)&sixBitsMask | (f.n[7]&twoBitsMask)<<6)
	b[8] = byte((f.n[7] >> 2) & eightBitsMask)
	b[7] = byte((f.n[7] >> 10) & eightBitsMask)
	b[6] = byte((f.n[7] >> 18) & eightBitsMask)
	b[5] = byte(f.n[8] & eightBitsMask)
	b[4] = byte((f.n[8] >> 8) & eightBitsMask)
	b[3] = byte((f.n[8] >> 16) & eightBitsMask)
	b[2] = byte((f.n[8]>>24)&twoBitsMask | (f.n[9]&sixBitsMask)<<2)
	b[1] = byte((f.n[9] >> 6) & eightBitsMask)
	b[0] = byte((f.n[9] >> 14) & eightBitsMask)
}

func (f *FieldVal) Bytes() *[32]byte {
	b := new([32]byte)
	f.PutBytes(b)
	return b
}

func (f *FieldVal) IsZero() bool {
	bits := f.n[0] | f.n[1] | f.n[2] | f.n[3] | f.n[4] |
		f.n[5] | f.n[6] | f.n[7] | f.n[8] | f.n[9]

	return bits == 0
}

func (f *FieldVal) IsOdd() bool {
	return f.n[0]&1 == 1
}

func (f *FieldVal) Equals(val *FieldVal) bool {
	bits := (f.n[0] ^ val.n[0]) | (f.n[1] ^ val.n[1]) | (f.n[2] ^ val.n[2]) |
		(f.n[3] ^ val.n[3]) | (f.n[4] ^ val.n[4]) | (f.n[5] ^ val.n[5]) |
		(f.n[6] ^ val.n[6]) | (f.n[7] ^ val.n[7]) | (f.n[8] ^ val.n[8]) |
		(f.n[9] ^ val.n[9])

	return bits == 0
}

func (f *FieldVal) NegateVal(val *FieldVal, magnitude uint32) *FieldVal {
	f.n[0] = (magnitude+1)*fieldPrimeWordZero - val.n[0]
	f.n[1] = (magnitude+1)*fieldPrimeWordOne - val.n[1]
	f.n[2] = (magnitude+1)*fieldBaseMask - val.n[2]
	f.n[3] = (magnitude+1)*fieldBaseMask - val.n[3]
	f.n[4] = (magnitude+1)*fieldBaseMask - val.n[4]
	f.n[5] = (magnitude+1)*fieldBaseMask - val.n[5]
	f.n[6] = (magnitude+1)*fieldBaseMask - val.n[6]
	f.n[7] = (magnitude+1)*fieldBaseMask - val.n[7]
	f.n[8] = (magnitude+1)*fieldBaseMask - val.n[8]
	f.n[9] = (magnitude+1)*fieldMSBMask - val.n[9]

	return f
}

func (f *FieldVal) Negate(magnitude uint32) *FieldVal {
	return f.NegateVal(f, magnitude)
}

func (f *FieldVal) AddInt(ui uint) *FieldVal {

	f.n[0] += uint32(ui)

	return f
}

func (f *FieldVal) Add(val *FieldVal) *FieldVal {
	f.n[0] += val.n[0]
	f.n[1] += val.n[1]
	f.n[2] += val.n[2]
	f.n[3] += val.n[3]
	f.n[4] += val.n[4]
	f.n[5] += val.n[5]
	f.n[6] += val.n[6]
	f.n[7] += val.n[7]
	f.n[8] += val.n[8]
	f.n[9] += val.n[9]

	return f
}

func (f *FieldVal) Add2(val *FieldVal, val2 *FieldVal) *FieldVal {
	f.n[0] = val.n[0] + val2.n[0]
	f.n[1] = val.n[1] + val2.n[1]
	f.n[2] = val.n[2] + val2.n[2]
	f.n[3] = val.n[3] + val2.n[3]
	f.n[4] = val.n[4] + val2.n[4]
	f.n[5] = val.n[5] + val2.n[5]
	f.n[6] = val.n[6] + val2.n[6]
	f.n[7] = val.n[7] + val2.n[7]
	f.n[8] = val.n[8] + val2.n[8]
	f.n[9] = val.n[9] + val2.n[9]

	return f
}

func (f *FieldVal) MulInt(val uint) *FieldVal {
	ui := uint32(val)
	f.n[0] *= ui
	f.n[1] *= ui
	f.n[2] *= ui
	f.n[3] *= ui
	f.n[4] *= ui
	f.n[5] *= ui
	f.n[6] *= ui
	f.n[7] *= ui
	f.n[8] *= ui
	f.n[9] *= ui

	return f
}

func (f *FieldVal) Mul(val *FieldVal) *FieldVal {
	return f.Mul2(f, val)
}

func (f *FieldVal) Mul2(val *FieldVal, val2 *FieldVal) *FieldVal {
	m := uint64(val.n[0]) * uint64(val2.n[0])
	t0 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[1]) +
		uint64(val.n[1])*uint64(val2.n[0])
	t1 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[2]) +
		uint64(val.n[1])*uint64(val2.n[1]) +
		uint64(val.n[2])*uint64(val2.n[0])
	t2 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[3]) +
		uint64(val.n[1])*uint64(val2.n[2]) +
		uint64(val.n[2])*uint64(val2.n[1]) +
		uint64(val.n[3])*uint64(val2.n[0])
	t3 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[4]) +
		uint64(val.n[1])*uint64(val2.n[3]) +
		uint64(val.n[2])*uint64(val2.n[2]) +
		uint64(val.n[3])*uint64(val2.n[1]) +
		uint64(val.n[4])*uint64(val2.n[0])
	t4 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[5]) +
		uint64(val.n[1])*uint64(val2.n[4]) +
		uint64(val.n[2])*uint64(val2.n[3]) +
		uint64(val.n[3])*uint64(val2.n[2]) +
		uint64(val.n[4])*uint64(val2.n[1]) +
		uint64(val.n[5])*uint64(val2.n[0])
	t5 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[6]) +
		uint64(val.n[1])*uint64(val2.n[5]) +
		uint64(val.n[2])*uint64(val2.n[4]) +
		uint64(val.n[3])*uint64(val2.n[3]) +
		uint64(val.n[4])*uint64(val2.n[2]) +
		uint64(val.n[5])*uint64(val2.n[1]) +
		uint64(val.n[6])*uint64(val2.n[0])
	t6 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[7]) +
		uint64(val.n[1])*uint64(val2.n[6]) +
		uint64(val.n[2])*uint64(val2.n[5]) +
		uint64(val.n[3])*uint64(val2.n[4]) +
		uint64(val.n[4])*uint64(val2.n[3]) +
		uint64(val.n[5])*uint64(val2.n[2]) +
		uint64(val.n[6])*uint64(val2.n[1]) +
		uint64(val.n[7])*uint64(val2.n[0])
	t7 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[8]) +
		uint64(val.n[1])*uint64(val2.n[7]) +
		uint64(val.n[2])*uint64(val2.n[6]) +
		uint64(val.n[3])*uint64(val2.n[5]) +
		uint64(val.n[4])*uint64(val2.n[4]) +
		uint64(val.n[5])*uint64(val2.n[3]) +
		uint64(val.n[6])*uint64(val2.n[2]) +
		uint64(val.n[7])*uint64(val2.n[1]) +
		uint64(val.n[8])*uint64(val2.n[0])
	t8 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[0])*uint64(val2.n[9]) +
		uint64(val.n[1])*uint64(val2.n[8]) +
		uint64(val.n[2])*uint64(val2.n[7]) +
		uint64(val.n[3])*uint64(val2.n[6]) +
		uint64(val.n[4])*uint64(val2.n[5]) +
		uint64(val.n[5])*uint64(val2.n[4]) +
		uint64(val.n[6])*uint64(val2.n[3]) +
		uint64(val.n[7])*uint64(val2.n[2]) +
		uint64(val.n[8])*uint64(val2.n[1]) +
		uint64(val.n[9])*uint64(val2.n[0])
	t9 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[1])*uint64(val2.n[9]) +
		uint64(val.n[2])*uint64(val2.n[8]) +
		uint64(val.n[3])*uint64(val2.n[7]) +
		uint64(val.n[4])*uint64(val2.n[6]) +
		uint64(val.n[5])*uint64(val2.n[5]) +
		uint64(val.n[6])*uint64(val2.n[4]) +
		uint64(val.n[7])*uint64(val2.n[3]) +
		uint64(val.n[8])*uint64(val2.n[2]) +
		uint64(val.n[9])*uint64(val2.n[1])
	t10 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[2])*uint64(val2.n[9]) +
		uint64(val.n[3])*uint64(val2.n[8]) +
		uint64(val.n[4])*uint64(val2.n[7]) +
		uint64(val.n[5])*uint64(val2.n[6]) +
		uint64(val.n[6])*uint64(val2.n[5]) +
		uint64(val.n[7])*uint64(val2.n[4]) +
		uint64(val.n[8])*uint64(val2.n[3]) +
		uint64(val.n[9])*uint64(val2.n[2])
	t11 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[3])*uint64(val2.n[9]) +
		uint64(val.n[4])*uint64(val2.n[8]) +
		uint64(val.n[5])*uint64(val2.n[7]) +
		uint64(val.n[6])*uint64(val2.n[6]) +
		uint64(val.n[7])*uint64(val2.n[5]) +
		uint64(val.n[8])*uint64(val2.n[4]) +
		uint64(val.n[9])*uint64(val2.n[3])
	t12 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[4])*uint64(val2.n[9]) +
		uint64(val.n[5])*uint64(val2.n[8]) +
		uint64(val.n[6])*uint64(val2.n[7]) +
		uint64(val.n[7])*uint64(val2.n[6]) +
		uint64(val.n[8])*uint64(val2.n[5]) +
		uint64(val.n[9])*uint64(val2.n[4])
	t13 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[5])*uint64(val2.n[9]) +
		uint64(val.n[6])*uint64(val2.n[8]) +
		uint64(val.n[7])*uint64(val2.n[7]) +
		uint64(val.n[8])*uint64(val2.n[6]) +
		uint64(val.n[9])*uint64(val2.n[5])
	t14 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[6])*uint64(val2.n[9]) +
		uint64(val.n[7])*uint64(val2.n[8]) +
		uint64(val.n[8])*uint64(val2.n[7]) +
		uint64(val.n[9])*uint64(val2.n[6])
	t15 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[7])*uint64(val2.n[9]) +
		uint64(val.n[8])*uint64(val2.n[8]) +
		uint64(val.n[9])*uint64(val2.n[7])
	t16 := m & fieldBaseMask
	m = (m >> fieldBase) +
		uint64(val.n[8])*uint64(val2.n[9]) +
		uint64(val.n[9])*uint64(val2.n[8])
	t17 := m & fieldBaseMask

	m = (m >> fieldBase) + uint64(val.n[9])*uint64(val2.n[9])
	t18 := m & fieldBaseMask

	t19 := m >> fieldBase
	m = t0 + t10*15632
	t0 = m & fieldBaseMask
	m = (m >> fieldBase) + t1 + t10*1024 + t11*15632
	t1 = m & fieldBaseMask
	m = (m >> fieldBase) + t2 + t11*1024 + t12*15632
	t2 = m & fieldBaseMask
	m = (m >> fieldBase) + t3 + t12*1024 + t13*15632
	t3 = m & fieldBaseMask
	m = (m >> fieldBase) + t4 + t13*1024 + t14*15632
	t4 = m & fieldBaseMask
	m = (m >> fieldBase) + t5 + t14*1024 + t15*15632
	t5 = m & fieldBaseMask
	m = (m >> fieldBase) + t6 + t15*1024 + t16*15632
	t6 = m & fieldBaseMask
	m = (m >> fieldBase) + t7 + t16*1024 + t17*15632
	t7 = m & fieldBaseMask
	m = (m >> fieldBase) + t8 + t17*1024 + t18*15632
	t8 = m & fieldBaseMask
	m = (m >> fieldBase) + t9 + t18*1024 + t19*68719492368
	t9 = m & fieldMSBMask
	m = m >> fieldMSBBits
	d := t0 + m*977
	f.n[0] = uint32(d & fieldBaseMask)
	d = (d >> fieldBase) + t1 + m*64
	f.n[1] = uint32(d & fieldBaseMask)
	f.n[2] = uint32((d >> fieldBase) + t2)
	f.n[3] = uint32(t3)
	f.n[4] = uint32(t4)
	f.n[5] = uint32(t5)
	f.n[6] = uint32(t6)
	f.n[7] = uint32(t7)
	f.n[8] = uint32(t8)
	f.n[9] = uint32(t9)

	return f
}

func (f *FieldVal) Square() *FieldVal {
	return f.SquareVal(f)
}

func (f *FieldVal) SquareVal(val *FieldVal) *FieldVal {
	m := uint64(val.n[0]) * uint64(val.n[0])
	t0 := m & fieldBaseMask

	// Terms for 2^(fieldBase*1).
	m = (m >> fieldBase) + 2*uint64(val.n[0])*uint64(val.n[1])
	t1 := m & fieldBaseMask

	// Terms for 2^(fieldBase*2).
	m = (m >> fieldBase) +
		2*uint64(val.n[0])*uint64(val.n[2]) +
		uint64(val.n[1])*uint64(val.n[1])
	t2 := m & fieldBaseMask

	// Terms for 2^(fieldBase*3).
	m = (m >> fieldBase) +
		2*uint64(val.n[0])*uint64(val.n[3]) +
		2*uint64(val.n[1])*uint64(val.n[2])
	t3 := m & fieldBaseMask

	// Terms for 2^(fieldBase*4).
	m = (m >> fieldBase) +
		2*uint64(val.n[0])*uint64(val.n[4]) +
		2*uint64(val.n[1])*uint64(val.n[3]) +
		uint64(val.n[2])*uint64(val.n[2])
	t4 := m & fieldBaseMask

	// Terms for 2^(fieldBase*5).
	m = (m >> fieldBase) +
		2*uint64(val.n[0])*uint64(val.n[5]) +
		2*uint64(val.n[1])*uint64(val.n[4]) +
		2*uint64(val.n[2])*uint64(val.n[3])
	t5 := m & fieldBaseMask

	// Terms for 2^(fieldBase*6).
	m = (m >> fieldBase) +
		2*uint64(val.n[0])*uint64(val.n[6]) +
		2*uint64(val.n[1])*uint64(val.n[5]) +
		2*uint64(val.n[2])*uint64(val.n[4]) +
		uint64(val.n[3])*uint64(val.n[3])
	t6 := m & fieldBaseMask

	// Terms for 2^(fieldBase*7).
	m = (m >> fieldBase) +
		2*uint64(val.n[0])*uint64(val.n[7]) +
		2*uint64(val.n[1])*uint64(val.n[6]) +
		2*uint64(val.n[2])*uint64(val.n[5]) +
		2*uint64(val.n[3])*uint64(val.n[4])
	t7 := m & fieldBaseMask

	// Terms for 2^(fieldBase*8).
	m = (m >> fieldBase) +
		2*uint64(val.n[0])*uint64(val.n[8]) +
		2*uint64(val.n[1])*uint64(val.n[7]) +
		2*uint64(val.n[2])*uint64(val.n[6]) +
		2*uint64(val.n[3])*uint64(val.n[5]) +
		uint64(val.n[4])*uint64(val.n[4])
	t8 := m & fieldBaseMask

	// Terms for 2^(fieldBase*9).
	m = (m >> fieldBase) +
		2*uint64(val.n[0])*uint64(val.n[9]) +
		2*uint64(val.n[1])*uint64(val.n[8]) +
		2*uint64(val.n[2])*uint64(val.n[7]) +
		2*uint64(val.n[3])*uint64(val.n[6]) +
		2*uint64(val.n[4])*uint64(val.n[5])
	t9 := m & fieldBaseMask

	// Terms for 2^(fieldBase*10).
	m = (m >> fieldBase) +
		2*uint64(val.n[1])*uint64(val.n[9]) +
		2*uint64(val.n[2])*uint64(val.n[8]) +
		2*uint64(val.n[3])*uint64(val.n[7]) +
		2*uint64(val.n[4])*uint64(val.n[6]) +
		uint64(val.n[5])*uint64(val.n[5])
	t10 := m & fieldBaseMask

	// Terms for 2^(fieldBase*11).
	m = (m >> fieldBase) +
		2*uint64(val.n[2])*uint64(val.n[9]) +
		2*uint64(val.n[3])*uint64(val.n[8]) +
		2*uint64(val.n[4])*uint64(val.n[7]) +
		2*uint64(val.n[5])*uint64(val.n[6])
	t11 := m & fieldBaseMask

	// Terms for 2^(fieldBase*12).
	m = (m >> fieldBase) +
		2*uint64(val.n[3])*uint64(val.n[9]) +
		2*uint64(val.n[4])*uint64(val.n[8]) +
		2*uint64(val.n[5])*uint64(val.n[7]) +
		uint64(val.n[6])*uint64(val.n[6])
	t12 := m & fieldBaseMask

	// Terms for 2^(fieldBase*13).
	m = (m >> fieldBase) +
		2*uint64(val.n[4])*uint64(val.n[9]) +
		2*uint64(val.n[5])*uint64(val.n[8]) +
		2*uint64(val.n[6])*uint64(val.n[7])
	t13 := m & fieldBaseMask

	// Terms for 2^(fieldBase*14).
	m = (m >> fieldBase) +
		2*uint64(val.n[5])*uint64(val.n[9]) +
		2*uint64(val.n[6])*uint64(val.n[8]) +
		uint64(val.n[7])*uint64(val.n[7])
	t14 := m & fieldBaseMask

	// Terms for 2^(fieldBase*15).
	m = (m >> fieldBase) +
		2*uint64(val.n[6])*uint64(val.n[9]) +
		2*uint64(val.n[7])*uint64(val.n[8])
	t15 := m & fieldBaseMask

	// Terms for 2^(fieldBase*16).
	m = (m >> fieldBase) +
		2*uint64(val.n[7])*uint64(val.n[9]) +
		uint64(val.n[8])*uint64(val.n[8])
	t16 := m & fieldBaseMask

	// Terms for 2^(fieldBase*17).
	m = (m >> fieldBase) + 2*uint64(val.n[8])*uint64(val.n[9])
	t17 := m & fieldBaseMask

	// Terms for 2^(fieldBase*18).
	m = (m >> fieldBase) + uint64(val.n[9])*uint64(val.n[9])
	t18 := m & fieldBaseMask

	t19 := m >> fieldBase
	m = t0 + t10*15632
	t0 = m & fieldBaseMask
	m = (m >> fieldBase) + t1 + t10*1024 + t11*15632
	t1 = m & fieldBaseMask
	m = (m >> fieldBase) + t2 + t11*1024 + t12*15632
	t2 = m & fieldBaseMask
	m = (m >> fieldBase) + t3 + t12*1024 + t13*15632
	t3 = m & fieldBaseMask
	m = (m >> fieldBase) + t4 + t13*1024 + t14*15632
	t4 = m & fieldBaseMask
	m = (m >> fieldBase) + t5 + t14*1024 + t15*15632
	t5 = m & fieldBaseMask
	m = (m >> fieldBase) + t6 + t15*1024 + t16*15632
	t6 = m & fieldBaseMask
	m = (m >> fieldBase) + t7 + t16*1024 + t17*15632
	t7 = m & fieldBaseMask
	m = (m >> fieldBase) + t8 + t17*1024 + t18*15632
	t8 = m & fieldBaseMask
	m = (m >> fieldBase) + t9 + t18*1024 + t19*68719492368
	t9 = m & fieldMSBMask
	m = m >> fieldMSBBits
	n := t0 + m*977
	f.n[0] = uint32(n & fieldBaseMask)
	n = (n >> fieldBase) + t1 + m*64
	f.n[1] = uint32(n & fieldBaseMask)
	f.n[2] = uint32((n >> fieldBase) + t2)
	f.n[3] = uint32(t3)
	f.n[4] = uint32(t4)
	f.n[5] = uint32(t5)
	f.n[6] = uint32(t6)
	f.n[7] = uint32(t7)
	f.n[8] = uint32(t8)
	f.n[9] = uint32(t9)

	return f
}

func (f *FieldVal) Inverse() *FieldVal {
	var a2, a3, a4, a10, a11, a21, a42, a45, a63, a1019, a1023 FieldVal
	a2.SquareVal(f)
	a3.Mul2(&a2, f)
	a4.SquareVal(&a2)
	a10.SquareVal(&a4).Mul(&a2)
	a11.Mul2(&a10, f)
	a21.Mul2(&a10, &a11)
	a42.SquareVal(&a21)
	a45.Mul2(&a42, &a3)
	a63.Mul2(&a42, &a21)
	a1019.SquareVal(&a63).Square().Square().Square().Mul(&a11)
	a1023.Mul2(&a1019, &a4)
	f.Set(&a63)                                    // f = a^(2^6 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^11 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^16 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^16 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^21 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^26 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^26 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^31 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^36 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^36 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^41 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^46 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^46 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^51 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^56 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^56 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^61 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^66 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^66 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^71 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^76 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^76 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^81 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^86 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^86 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^91 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^96 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^96 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^101 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^106 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^106 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^111 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^116 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^116 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^121 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^126 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^126 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^131 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^136 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^136 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^141 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^146 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^146 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^151 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^156 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^156 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^161 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^166 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^166 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^171 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^176 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^176 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^181 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^186 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^186 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^191 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^196 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^196 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^201 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^206 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^206 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^211 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^216 - 1024)
	f.Mul(&a1023)                                  // f = a^(2^216 - 1)
	f.Square().Square().Square().Square().Square() // f = a^(2^221 - 32)
	f.Square().Square().Square().Square().Square() // f = a^(2^226 - 1024)
	f.Mul(&a1019)                                  // f = a^(2^226 - 5)
	f.Square().Square().Square().Square().Square() // f = a^(2^231 - 160)
	f.Square().Square().Square().Square().Square() // f = a^(2^236 - 5120)
	f.Mul(&a1023)                                  // f = a^(2^236 - 4097)
	f.Square().Square().Square().Square().Square() // f = a^(2^241 - 131104)
	f.Square().Square().Square().Square().Square() // f = a^(2^246 - 4195328)
	f.Mul(&a1023)                                  // f = a^(2^246 - 4194305)
	f.Square().Square().Square().Square().Square() // f = a^(2^251 - 134217760)
	f.Square().Square().Square().Square().Square() // f = a^(2^256 - 4294968320)
	return f.Mul(&a45)                             // f = a^(2^256 - 4294968275) = a^(p-2)
}

func (f *FieldVal) SqrtVal(x *FieldVal) *FieldVal {
	f.SetInt(1)
	for _, b := range fieldQBytes {
		switch b {

		// Most common case, where all 8 bits are set.
		case 0xff:
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
		case 0x3f:
			f.Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)

		// Byte 28 of Q (0xbf), where only bit 7 is unset.
		case 0xbf:
			f.Square().Mul(x)
			f.Square()
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square().Mul(x)

		// Byte 31 of Q (0x0c), where only bits 3 and 4 are set.
		default:
			f.Square()
			f.Square()
			f.Square()
			f.Square()
			f.Square().Mul(x)
			f.Square().Mul(x)
			f.Square()
			f.Square()
		}
	}

	return f
}

func (f *FieldVal) Sqrt() *FieldVal {
	return f.SqrtVal(f)
}
