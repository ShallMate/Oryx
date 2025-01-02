package curve

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
)

// Extract generates a pseudorandom key for use with Expand from an input secret
// and an optional independent salt.
//
// Only use this function if you need to reuse the extracted key with multiple
// Expand invocations and different context values. Most common scenarios,
// including the generation of multiple keys, should use New instead.
func Extract(hash func() hash.Hash, secret, salt []byte) []byte {
	if salt == nil {
		salt = make([]byte, hash().Size())
	}
	extractor := hmac.New(hash, salt)
	extractor.Write(secret)
	return extractor.Sum(nil)
}

type hkdf struct {
	expander hash.Hash
	size     int

	info    []byte
	counter byte

	prev []byte
	buf  []byte
}

func (f *hkdf) Read(p []byte) (int, error) {
	// Check whether enough data can be generated
	need := len(p)
	remains := len(f.buf) + int(255-f.counter+1)*f.size
	if remains < need {
		return 0, errors.New("hkdf: entropy limit reached")
	}
	// Read any leftover from the buffer
	n := copy(p, f.buf)
	p = p[n:]

	// Fill the rest of the buffer
	for len(p) > 0 {
		if f.counter > 1 {
			f.expander.Reset()
		}
		f.expander.Write(f.prev)
		f.expander.Write(f.info)
		f.expander.Write([]byte{f.counter})
		f.prev = f.expander.Sum(f.prev[:0])
		f.counter++

		// Copy the new batch into p
		f.buf = f.prev
		n = copy(p, f.buf)
		p = p[n:]
	}
	// Save leftovers for next run
	f.buf = f.buf[n:]

	return need, nil
}

// Expand returns a Reader, from which keys can be read, using the given
// pseudorandom key and optional context info, skipping the extraction step.
//
// The pseudorandomKey should have been generated by Extract, or be a uniformly
// random or pseudorandom cryptographically strong key. See RFC 5869, Section
// 3.3. Most common scenarios will want to use New instead.
func Expand(hash func() hash.Hash, pseudorandomKey, info []byte) io.Reader {
	expander := hmac.New(hash, pseudorandomKey)
	return &hkdf{expander, expander.Size(), info, 1, nil, nil}
}

// New returns a Reader, from which keys can be read, using the given hash,
// secret, salt and context info. Salt and info can be nil.
func New(hash func() hash.Hash, secret, salt, info []byte) io.Reader {
	prk := Extract(hash, secret, salt)
	return Expand(hash, prk, info)
}

type gfP [4]uint64

func newGFp(x int64) (out *gfP) {
	if x >= 0 {
		out = &gfP{uint64(x)}
	} else {
		out = &gfP{uint64(-x)}
		gfpNeg(out, out)
	}

	montEncode(out, out)
	return out
}

// hashToBase implements hashing a message to an element of the field.
//
// L = ceil((256+128)/8)=48, ctr = 0, i = 1
func hashToBase(msg, dst []byte) *gfP {
	var t [48]byte
	info := []byte{'H', '2', 'C', byte(0), byte(1)}
	r := New(sha256.New, msg, dst, info)
	if _, err := r.Read(t[:]); err != nil {
		panic(err)
	}
	var x big.Int
	v := x.SetBytes(t[:]).Mod(&x, p).Bytes()
	v32 := [32]byte{}
	for i := len(v) - 1; i >= 0; i-- {
		v32[len(v)-1-i] = v[i]
	}
	u := &gfP{
		binary.LittleEndian.Uint64(v32[0*8 : 1*8]),
		binary.LittleEndian.Uint64(v32[1*8 : 2*8]),
		binary.LittleEndian.Uint64(v32[2*8 : 3*8]),
		binary.LittleEndian.Uint64(v32[3*8 : 4*8]),
	}
	montEncode(u, u)
	return u
}

func (e *gfP) String() string {
	return fmt.Sprintf("%16.16x%16.16x%16.16x%16.16x", e[3], e[2], e[1], e[0])
}

func (e *gfP) Set(f *gfP) {
	e[0] = f[0]
	e[1] = f[1]
	e[2] = f[2]
	e[3] = f[3]
}

func (e *gfP) exp(f *gfP, bits [4]uint64) {
	sum, power := &gfP{}, &gfP{}
	sum.Set(rN1)
	power.Set(f)

	for word := 0; word < 4; word++ {
		for bit := uint(0); bit < 64; bit++ {
			if (bits[word]>>bit)&1 == 1 {
				gfpMul(sum, sum, power)
			}
			gfpMul(power, power, power)
		}
	}

	gfpMul(sum, sum, r3)
	e.Set(sum)
}

func (e *gfP) Invert(f *gfP) {
	e.exp(f, pMinus2)
}

func (e *gfP) Sqrt(f *gfP) {
	// Since p = 8k+5,
	// if f^((k+1)/4) = 1 mod p, then
	// e = f^(k+1) is a root of f;
	//else if f^((k+1)/4) = -1 mod p, then
	//e = 2^(2k+1)*f^(k+1) is a root of f.
	one := newGFp(1)
	tmp := new(gfP)
	tmp.exp(f, pMinus1Over4)

	if *tmp == *one {
		e.exp(f, pPlus3Over4)
	} else if *tmp == pMinus1 {
		e.exp(f, pPlus3Over4)
		gfpMul(e, e, &twoTo2kPlus1)
	}
}

func (e *gfP) Marshal(out []byte) {
	for w := uint(0); w < 4; w++ {
		for b := uint(0); b < 8; b++ {
			out[8*w+b] = byte(e[3-w] >> (56 - 8*b))
		}
	}
}

func (e *gfP) Unmarshal(in []byte) {
	for w := uint(0); w < 4; w++ {
		e[3-w] = 0
		for b := uint(0); b < 8; b++ {
			e[3-w] += uint64(in[8*w+b]) << (56 - 8*b)
		}
	}
}

func montEncode(c, a *gfP) { gfpMul(c, a, r2) }
func montDecode(c, a *gfP) { gfpMul(c, a, &gfP{1}) }

func sign0(e *gfP) int {
	x := &gfP{}
	montDecode(x, e)
	for w := 3; w >= 0; w-- {
		if x[w] > pMinus1Over2[w] {
			return 1
		} else if x[w] < pMinus1Over2[w] {
			return -1
		}
	}
	return 1
}

func legendre(e *gfP) int {
	f := &gfP{}
	// Since p = 8k+5, then e^(4k+2) is the Legendre symbol of e.
	f.exp(e, pMinus1Over2)

	montDecode(f, f)

	if *f != [4]uint64{} {
		return 2*int(f[0]&1) - 1
	}

	return 0
}

func (c *gfP) Println() {
	fmt.Print("&gfP{")
	y, _ := new(big.Int).SetString(c.String(), 16)
	words := y.Bits()
	for _, word := range words[:len(words)-1] {
		fmt.Printf("%#x, ", word)
	}
	fmt.Printf("%#x}\n\n", words[len(words)-1])
}

func gfpFromString(s string) gfP {
	y, _ := new(big.Int).SetString(s, 16)
	words := y.Bits()
	var a = gfP{}
	for i := 0; i < len(words); i++ {
		a[i] = uint64(words[i])
	}
	return a
}