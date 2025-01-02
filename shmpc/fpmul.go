package shmpc

import (
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/Oryx/curve"
)

func (system *ShareSystem) Share_An_Fp_Mul(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	wg.Add(1)
	go system.Send(&wg, ori_value.Bytes())
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].Share, _ = rand.Int(rand.Reader, system.Order)
			shareinv := new(big.Int).ModInverse(shares[i].Share, system.Order)
			ori_value = ori_value.Mul(ori_value, shareinv)
			ori_value = ori_value.Mod(ori_value, system.Order)
		} else {
			shares[i].Share = ori_value
		}
		wg.Add(1)
		go system.Send(&wg, shares[i].Share.Bytes())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) Share_An_Fp_Mul_Offline(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	wg.Add(1)
	go system.OfflineBroadcastN(&wg, ori_value.Bytes())
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].Share, _ = rand.Int(rand.Reader, system.Order)
			shareinv := new(big.Int).ModInverse(shares[i].Share, system.Order)
			ori_value = ori_value.Mul(ori_value, shareinv)
			ori_value = ori_value.Mod(ori_value, system.Order)
		} else {
			shares[i].Share = ori_value
		}
		wg.Add(1)
		go system.OfflineSend(&wg, shares[i].Share.Bytes())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) Share_An_Fp_for_EXP(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	wg.Add(1)
	go system.Send(&wg, ori_value.Bytes())
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].Share, _ = curve.RandomK(rand.Reader)
			ori_value = ori_value.Sub(ori_value, shares[i].Share)
			ori_value = ori_value.Mod(ori_value, system.OrderMul)
		} else {
			shares[i].Share = ori_value
		}
		wg.Add(1)
		go system.Send(&wg, shares[i].Share.Bytes())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) Share_An_Fp_for_EXP_Offline(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	wg.Add(1)
	go system.OfflineBroadcastN(&wg, ori_value.Bytes())
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].Share, _ = curve.RandomK(rand.Reader)
			ori_value = ori_value.Sub(ori_value, shares[i].Share)
			ori_value = ori_value.Mod(ori_value, system.OrderMul)
		} else {
			shares[i].Share = ori_value
		}
		wg.Add(1)
		go system.OfflineSend(&wg, shares[i].Share.Bytes())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) RandomShareFp_Mul() *[]Share_Fp {
	r, _ := curve.RandomK(rand.Reader)
	rshares := system.Share_An_Fp_Mul_Offline(r)
	return rshares
}

func (system *ShareSystem) shareMul(shares1, shares2 Share_Fp) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	shares.Share = new(big.Int).Mul(shares1.Share, shares2.Share)
	shares.Share = shares.Share.Mod(shares.Share, system.Order)
	return *shares
}

func (system *ShareSystem) SecMulPlaintext_Mul(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	shares2 := system.Share_An_Fp_Mul(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareMul(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) SecMul_Mul(shares1, shares2 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareMul(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) shareDiv(shares1, shares2 Share_Fp) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	shares.Share = new(big.Int).Mul(shares1.Share, new(big.Int).ModInverse(shares2.Share, system.Order))
	shares.Share = shares.Share.Mod(shares.Share, system.Order)
	return *shares
}

func (system *ShareSystem) SecDiv_Plaintext_1(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	shares2 := system.Share_An_Fp_Mul(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareDiv(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) SecDiv_Plaintext_2(scalar *big.Int, shares1 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	shares2 := system.Share_An_Fp_Mul(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareDiv((*shares2)[i], shares1[i])
	}
	return &shares
}

func (system *ShareSystem) SecDiv(shares1, shares2 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareDiv(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) EXP_P_Fp_1(scalar *big.Int, shares1 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Share = new(big.Int).Exp(scalar, shares1[i].Share, system.Order)
	}
	return &shares
}

func (system *ShareSystem) EXP_P_Fp_2(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Share = new(big.Int).Exp(shares1[i].Share, scalar, system.Order)
	}
	return &shares
}

func (system *ShareSystem) HalfOpenFp_for_Exp(shares []Share_Fp) *big.Int {
	var wg sync.WaitGroup
	ori_value := big.NewInt(0)
	for i := 0; i < system.Partynum; i++ {
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Bytes())
		ori_value = ori_value.Add(ori_value, shares[i].Share)
	}
	ori_value = ori_value.Mod(ori_value, system.OrderMul)
	wg.Wait()
	return ori_value
}

func (system *ShareSystem) OpenFp_Mul(shares []Share_Fp) *big.Int {
	var wg sync.WaitGroup
	ori_value := big.NewInt(1)
	for i := 0; i < system.Partynum; i++ {
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Bytes())
		ori_value = ori_value.Mul(ori_value, shares[i].Share)
	}
	ori_value = ori_value.Mod(ori_value, system.Order)
	wg.Wait()
	return ori_value
}
