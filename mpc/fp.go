package mpc

import (
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/Oryx/curve"
)

type Share_Fp struct {
	Share *big.Int
	Gama  *big.Int
	Delta *big.Int
	Index int
}

func (system *ShareSystem) Share_An_Fp(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	Delta, _ := curve.RandomK(rand.Reader)
	wg.Add(2)
	go system.Send(&wg, ori_value.Bytes())
	go system.BroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Add(ori_value, Delta)
	Gama = Gama.Mul(system.alpha, Gama)
	Gama = Gama.Mod(Gama, system.Order)
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = curve.RandomK(rand.Reader)
			shares[i].Gama, _ = curve.RandomK(rand.Reader)
			ori_value = ori_value.Sub(ori_value, shares[i].Share)
			ori_value = ori_value.Mod(ori_value, system.Order)
			Gama = Gama.Sub(Gama, shares[i].Gama)
			Gama = Gama.Mod(Gama, system.Order)
		} else {
			shares[i].Share = ori_value
			shares[i].Gama = Gama
		}
		wg.Add(2)
		go system.Send(&wg, shares[i].Share.Bytes())
		go system.Send(&wg, shares[i].Gama.Bytes())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) Share_An_Fp_Offline(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	Delta, _ := curve.RandomK(rand.Reader)
	wg.Add(1)
	go system.OfflineBroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Add(ori_value, Delta)
	Gama = Gama.Mul(system.alpha, Gama)
	Gama = Gama.Mod(Gama, system.Order)
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = curve.RandomK(rand.Reader)
			shares[i].Gama, _ = curve.RandomK(rand.Reader)
			ori_value = ori_value.Sub(ori_value, shares[i].Share)
			ori_value = ori_value.Mod(ori_value, system.Order)
			Gama = Gama.Sub(Gama, shares[i].Gama)
			Gama = Gama.Mod(Gama, system.Order)
		} else {
			shares[i].Share = ori_value
			shares[i].Gama = Gama
		}
		wg.Add(2)
		go system.OfflineSend(&wg, shares[i].Share.Bytes())
		go system.OfflineSend(&wg, shares[i].Gama.Bytes())
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
	Delta := new(big.Int).Add(shares1.Delta, shares2.Delta)
	Delta = Delta.Mod(Delta, system.Order)
	shares.Delta = Delta
	shares.Gama = new(big.Int).Add(shares1.Gama, shares2.Gama)
	shares.Gama = shares.Gama.Mod(shares.Gama, system.Order)
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
	Delta := new(big.Int).Sub(shares1.Delta, shares2.Delta)
	Delta = Delta.Mod(Delta, system.Order)
	shares.Delta = Delta
	shares.Gama = new(big.Int).Sub(shares1.Gama, shares2.Gama)
	shares.Gama = shares.Gama.Mod(shares.Gama, system.Order)
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
	Delta := new(big.Int).Mul(shares1.Delta, scalar)
	Delta = Delta.Mod(Delta, system.Order)
	shares.Delta = Delta
	shares.Gama = new(big.Int).Mul(shares1.Gama, scalar)
	shares.Gama = shares.Gama.Mod(shares.Gama, system.Order)
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
	e := system.HalfOpenFp(*eshares)
	f := system.HalfOpenFp(*fshares)
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
	e := system.HalfOpenFp(*eshares)
	e2 := new(big.Int).Mul(e, two)
	e2 = e2.Mod(e2, system.Order)
	ee := new(big.Int).Exp(e, two, system.Order)
	shares := system.SecMulPlaintext(shares1, e2)
	shares = system.SecAdd(*shares, *sharesB)
	shares = system.SecSubPlaintext(*shares, ee)
	return shares
}

func (system *ShareSystem) HalfOpenFp(shares []Share_Fp) *big.Int {
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

func (system *ShareSystem) MacCheckFp(shares []Share_Fp, res_value *big.Int) bool {
	var wg sync.WaitGroup
	chk := big.NewInt(0)
	delta := new(big.Int)
	t := new(big.Int).Add(res_value, shares[0].Delta)
	t = t.Mod(t, system.Order)
	for i := 0; i < system.Partynum; i++ {
		delta = delta.Mul(system.Alphas[i], t)
		delta = delta.Sub(shares[i].Gama, delta)
		commit, r := Com(delta.Bytes())
		wg.Add(3)
		go system.Broadcast(&wg, commit)
		go system.Broadcast(&wg, r.Bytes())
		go system.Broadcast(&wg, delta.Bytes())
		opencommit := OpenComit(delta.Bytes(), commit, r)
		if !opencommit {
			return false
		}
		chk = chk.Add(chk, delta)
	}
	chk = chk.Mod(chk, system.Order)
	wg.Wait()
	return chk.Cmp(zero) == 0
}

func (system *ShareSystem) OpenFp(shares []Share_Fp) (*big.Int, bool) {
	ori_value := big.NewInt(0)
	var wg sync.WaitGroup
	for i := 0; i < system.Partynum; i++ {
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Bytes())
		ori_value = ori_value.Add(ori_value, shares[i].Share)
	}
	ori_value = ori_value.Mod(ori_value, system.Order)
	wg.Wait()
	chk := system.MacCheckFp(shares, ori_value)
	return ori_value, chk
}
