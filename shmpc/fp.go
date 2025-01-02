package shmpc

import (
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/Oryx/curve"
)

type Share_Fp struct {
	Share *big.Int
	Index int
}

func (system *ShareSystem) Share_An_Fp(element *big.Int) *[]Share_Fp {
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

func (system *ShareSystem) Share_An_Fp_Offline(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].Share, _ = curve.RandomK(rand.Reader)
			ori_value = ori_value.Sub(ori_value, shares[i].Share)
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

func (system *ShareSystem) RandomShareFp() *[]Share_Fp {
	r, _ := curve.RandomK(rand.Reader)
	rshares := system.Share_An_Fp_Offline(r)
	return rshares
}

func (system *ShareSystem) shareAdd(shares1, shares2 Share_Fp) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	shares.Share = new(big.Int).Add(shares1.Share, shares2.Share)
	shares.Share = shares.Share.Mod(shares.Share, system.Order)
	return *shares
}

func (system *ShareSystem) SecAddPlaintext(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	shares2 := system.Share_An_Fp(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) SecAdd(shares1, shares2 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) SecSubPlaintext(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	shares2 := system.Share_An_Fp(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) shareSub(shares1, shares2 Share_Fp) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	shares.Share = new(big.Int).Sub(shares1.Share, shares2.Share)
	shares.Share = shares.Share.Mod(shares.Share, system.Order)
	return *shares
}

func (system *ShareSystem) SecSub(shares1, shares2 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) shareMulPlaintext(shares1 Share_Fp, scalar *big.Int) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	shares.Share = new(big.Int).Mul(shares1.Share, scalar)
	shares.Share = shares.Share.Mod(shares.Share, system.Order)
	return *shares
}

func (system *ShareSystem) SecMulPlaintext(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareMulPlaintext(shares1[i], scalar)
	}
	return &shares
}

func (system *ShareSystem) SecMul(shares1, shares2 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	sharesA, sharesB, sharesC := system.GenTriplets()
	eshares := system.SecSub(shares1, *sharesA)
	fshares := system.SecSub(shares2, *sharesB)
	e := system.OpenFp(*eshares)
	f := system.OpenFp(*fshares)
	ef := new(big.Int).Mul(e, f)
	ef = ef.Mod(ef, system.Order)
	efshares := system.Share_An_Fp(ef)
	for i := 0; i < system.Partynum; i++ {
		beshare := system.shareMulPlaintext((*sharesB)[i], e)
		afshare := system.shareMulPlaintext((*sharesA)[i], f)
		shares[i] = system.shareAdd((*sharesC)[i], (*efshares)[i])
		shares[i] = system.shareAdd(shares[i], beshare)
		shares[i] = system.shareAdd(shares[i], afshare)
	}
	return &shares
}

func (system *ShareSystem) SecSquare(shares1 []Share_Fp) *[]Share_Fp {
	sharesA, sharesB := system.GenSquarePair()
	eshares := system.SecSub(shares1, *sharesA)
	e := system.OpenFp(*eshares)
	e2 := new(big.Int).Mul(e, two)
	e2 = e2.Mod(e2, system.Order)
	ee := new(big.Int).Exp(e, two, system.Order)
	shares := system.SecMulPlaintext(shares1, e2)
	shares = system.SecAdd(*shares, *sharesB)
	shares = system.SecSubPlaintext(*shares, ee)
	return shares
}

func (system *ShareSystem) OpenFp(shares []Share_Fp) *big.Int {
	var wg sync.WaitGroup
	ori_value := big.NewInt(0)
	for i := 0; i < system.Partynum; i++ {
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Bytes())
		ori_value = ori_value.Add(ori_value, shares[i].Share)
	}
	ori_value = ori_value.Mod(ori_value, system.Order)
	wg.Wait()
	return ori_value
}
