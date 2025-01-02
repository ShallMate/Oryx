package mpc

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/Oryx/curve"
)

func Com(msg []byte) ([]byte, *big.Int) {
	r, _ := curve.RandomK(rand.Reader)
	hasher := sha256.New()
	msg = append(msg, r.Bytes()...)
	hasher.Write(msg)
	commit := hasher.Sum(nil)[:]
	return commit, r
}

func OpenComit(msg []byte, commit []byte, r *big.Int) bool {
	hasher := sha256.New()
	msg = append(msg, r.Bytes()...)
	hasher.Write(msg)
	return bytes.Equal(hasher.Sum(nil)[:], commit)
}
