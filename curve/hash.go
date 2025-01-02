package curve

import (
	"crypto/sha256"
	"math/big"
)

func HashToG1(msg []byte) *G1 {
	hm := sha256.Sum256(msg)
	hmInt := new(big.Int).SetBytes(hm[:])
	hmInt = hmInt.Mod(hmInt, Order)
	hmpoint := new(G1).ScalarBaseMult(hmInt)
	return hmpoint
}
