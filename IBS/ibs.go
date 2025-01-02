package ibs

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/Oryx/curve"
)

// MasterKey contains a master secret key and a master public key.
type MasterKey struct {
	Msk *big.Int
	MasterPubKey
}

type MasterPubKey struct {
	Mpk *curve.G2
	G   *curve.GT
}

// UserKey contains a secret key.
type UserKey struct {
	Sk *curve.G1
}

type Sig struct {
	S1 *curve.GT
	S2 *curve.G1
}

func H1(msg []byte) *big.Int {
	hasher := sha256.New()
	hasher.Write([]byte([]byte(msg)))
	hashbyte := hasher.Sum(nil)[:]
	h1 := new(big.Int).SetBytes(hashbyte)
	h1 = h1.Mod(h1, curve.Order)
	return h1
}

func H2(id *big.Int) *big.Int {
	hasher := sha256.New()
	hasher.Write(id.Bytes())
	hashbyte := hasher.Sum(nil)[:]
	h2 := new(big.Int).SetBytes(hashbyte)
	h2 = h2.Mod(h2, curve.Order)
	return h2
}

func H3(s1 *curve.GT) *big.Int {
	s1byte := s1.Marshal()
	hasher := sha256.New()
	hasher.Write([]byte(s1byte))
	s1byte = hasher.Sum(nil)[:]
	h3 := new(big.Int).SetBytes(s1byte)
	h3 = h3.Mod(h3, curve.Order)
	return h3
}

func UserKeyGen(mk *MasterKey, id *big.Int) *UserKey {
	uk := new(UserKey)
	hid := H2(id)
	hid = hid.Add(hid, mk.Msk)
	hid = hid.ModInverse(hid, curve.Order)
	uk.Sk = new(curve.G1).ScalarBaseMult(hid)
	return uk
}

func MasterKeyGen() *MasterKey {
	s, _ := rand.Int(rand.Reader, curve.Order)
	mk := new(MasterKey)
	mk.Msk = s
	mk.Mpk = new(curve.G2).ScalarBaseMult(s)
	mk.G = curve.Pair(curve.Gen1, curve.Gen2)
	return mk
}

func Sign(uk *UserKey, mpk *MasterPubKey, msg []byte) *Sig {
	sig := new(Sig)
	x, _ := rand.Int(rand.Reader, curve.Order)
	sig.S1 = new(curve.GT).ScalarMult(mpk.G, x)
	hs1 := H3(sig.S1)
	hm := H1(msg)
	hs1 = hs1.Add(hm, hs1)
	hs1 = new(big.Int).Sub(x, hs1)
	hs1 = hs1.Mod(hs1, curve.Order)
	sig.S2 = new(curve.G1).ScalarMult(uk.Sk, hs1)
	return sig
}

func Ver(sig *Sig, msg []byte, id *big.Int, mpk *MasterPubKey) bool {
	hid := H2(id)
	P := new(curve.G2).ScalarBaseMult(hid)
	P = P.Add(P, mpk.Mpk)
	u := curve.Pair(sig.S2, P)
	hm := H1(msg)
	hs1 := H3(sig.S1)
	hm = hm.Add(hm, hs1)
	hm = hm.Mod(hm, curve.Order)
	t := new(curve.GT).ScalarMult(mpk.G, hm)
	t = t.Add(u, t)
	wbytes := t.Marshal()
	s1bytes := sig.S1.Marshal()
	return bytes.Equal(wbytes, s1bytes)
}
