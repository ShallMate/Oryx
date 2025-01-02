package bls

import (
	"bytes"
	"crypto/rand"
	"math/big"

	"github.com/Oryx/curve"
)

type PublicKey struct {
	PK *curve.G2
}

type PrivateKey struct {
	sk     *big.Int
	Pubkey *PublicKey
}

type Sig struct {
	S *curve.G1
}

func KeyGen() (*PrivateKey, *PublicKey) {
	sk := new(PrivateKey)
	sk.sk, _ = rand.Int(rand.Reader, curve.Order)
	sk.Pubkey = new(PublicKey)
	sk.Pubkey.PK = new(curve.G2).ScalarMult(curve.Gen2, sk.sk)
	return sk, sk.Pubkey
}

func Sign(sk *PrivateKey, msg []byte) *Sig {
	hm := curve.HashToG1(msg)
	sig := new(Sig)
	sig.S = new(curve.G1).ScalarMult(hm, sk.sk)
	return sig
}

func SignwithHm(sk *PrivateKey, msg []byte) (*Sig, *curve.G1) {
	hm := curve.HashToG1(msg)
	sig := new(Sig)
	sig.S = new(curve.G1).ScalarMult(hm, sk.sk)
	return sig, hm
}

func Verify(pk *PublicKey, sig *Sig, msg []byte) bool {
	hm := curve.HashToG1(msg)
	left := curve.Pair(sig.S, curve.Gen2)
	right := curve.Pair(hm, pk.PK)
	return bytes.Equal(left.Marshal(), right.Marshal())
}
