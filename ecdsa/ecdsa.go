package ecdsa

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/Oryx/ecc"
)

type ECDSA struct {
	curve *ecc.KoblitzCurve
}

type PublicKey struct {
	PKX *big.Int
	PKY *big.Int
}

type PrivateKey struct {
	sk     *big.Int
	Pubkey *PublicKey
}

type Sig struct {
	R *big.Int
	S *big.Int
}

type SigInv struct {
	RX *big.Int
	RY *big.Int
	S  *big.Int
	HM *big.Int
}

func NewECDSA() *ECDSA {
	return &ECDSA{curve: ecc.S256()}
}

// KeyGen generates a new private key and corresponding public key.
func (ecdsa *ECDSA) KeyGen() (*PrivateKey, *PublicKey) {
	sk := new(PrivateKey)
	sk.Pubkey = new(PublicKey) // Initialize sk.Pubkey
	sk.sk, _ = rand.Int(rand.Reader, ecdsa.curve.N)
	sk.Pubkey.PKX, sk.Pubkey.PKY = ecdsa.curve.ScalarMult(ecdsa.curve.Gx, ecdsa.curve.Gy, sk.sk.Bytes())
	return sk, sk.Pubkey
}

// Sign generates a signature for a message using a private key.
func (ecdsa *ECDSA) Sign(sk *PrivateKey, msg []byte) *Sig {
	hm := sha256.Sum256(msg)
	hmInt := new(big.Int).SetBytes(hm[:])
	hmInt = hmInt.Mod(hmInt, ecdsa.curve.N)
	k, _ := rand.Int(rand.Reader, ecdsa.curve.N)
	r, _ := ecdsa.curve.ScalarMult(ecdsa.curve.Gx, ecdsa.curve.Gy, k.Bytes())
	r = r.Mod(r, ecdsa.curve.N)
	s := new(big.Int).ModInverse(k, ecdsa.curve.N)
	s = s.Mul(s, new(big.Int).Add(hmInt, new(big.Int).Mul(sk.sk, r)))
	s = s.Mod(s, ecdsa.curve.N)
	return &Sig{R: r, S: s}
}

// Sign generates a signature for a message using a private key.
func (ecdsa *ECDSA) SignwithInv(sk *PrivateKey, msg []byte) *SigInv {
	hm := sha256.Sum256(msg)
	hmInt := new(big.Int).SetBytes(hm[:])
	hmInt = hmInt.Mod(hmInt, ecdsa.curve.N)
	k, _ := rand.Int(rand.Reader, ecdsa.curve.N)
	rx, ry := ecdsa.curve.ScalarMult(ecdsa.curve.Gx, ecdsa.curve.Gy, k.Bytes())
	s := new(big.Int).ModInverse(k, ecdsa.curve.N)
	s = s.Mul(s, new(big.Int).Add(hmInt, new(big.Int).Mul(sk.sk, new(big.Int).Mod(rx, ecdsa.curve.N))))
	s = s.ModInverse(s, ecdsa.curve.N)
	return &SigInv{RX: rx, RY: ry, S: s, HM: hmInt}
}

// Verify checks if a signature is valid for a message using a public key.
func (ecdsa *ECDSA) Verify(pk *PublicKey, sig *Sig, msg []byte) bool {
	hm := sha256.Sum256(msg)
	hmInt := new(big.Int).SetBytes(hm[:])
	hmInt = hmInt.Mod(hmInt, ecdsa.curve.N)
	sinv := new(big.Int).ModInverse(sig.S, ecdsa.curve.N)
	u1 := new(big.Int).Mul(hmInt, sinv)
	u1 = u1.Mod(u1, ecdsa.curve.N)
	u2 := new(big.Int).Mul(sig.R, sinv)
	u2 = u2.Mod(u2, ecdsa.curve.N)
	x, y := ecdsa.curve.ScalarMult(pk.PKX, pk.PKY, u2.Bytes())
	x2, y2 := ecdsa.curve.ScalarMult(ecdsa.curve.Gx, ecdsa.curve.Gy, u1.Bytes())
	x, _ = ecdsa.curve.Add(x, y, x2, y2)
	x = x.Mod(x, ecdsa.curve.N)
	return x.Cmp(sig.R) == 0
}

func (ecdsa *ECDSA) VerifyWithoutInv(pk *PublicKey, sig *SigInv) bool {
	u1 := new(big.Int).Mul(sig.HM, sig.S)
	u1 = u1.Mod(u1, ecdsa.curve.N)
	u2 := new(big.Int).Mul(sig.RX, sig.S)
	u2 = u2.Mod(u2, ecdsa.curve.N)
	x, y := ecdsa.curve.ScalarMult(pk.PKX, pk.PKY, u2.Bytes())
	x2, y2 := ecdsa.curve.ScalarMult(ecdsa.curve.Gx, ecdsa.curve.Gy, u1.Bytes())
	x, y = ecdsa.curve.Add(x, y, x2, y2)
	return x.Cmp(sig.RX) == 0 && y.Cmp(sig.RY) == 0
}
