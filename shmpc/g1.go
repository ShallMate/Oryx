package shmpc

import (
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/Oryx/curve"
)

type Share_G1 struct {
	Share *curve.G1
	Index int
}

func (system *ShareSystem) Share_A_G1(element *curve.G1) *[]Share_G1 {
	var wg sync.WaitGroup
	ori_value := new(curve.G1).Set(element)
	wg.Add(1)
	go system.Send(&wg, ori_value.Marshal())
	shares := make([]Share_G1, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			_, shares[i].Share, _ = curve.RandomG1(rand.Reader)
			ori_value = ori_value.Add(ori_value, new(curve.G1).Neg(shares[i].Share))
		} else {
			shares[i].Share = ori_value
		}
		wg.Add(1)
		go system.Send(&wg, shares[i].Share.Marshal())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) Share_A_G1_Offline(element *curve.G1) *[]Share_G1 {
	var wg sync.WaitGroup
	ori_value := new(curve.G1).Set(element)
	shares := make([]Share_G1, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			_, shares[i].Share, _ = curve.RandomG1(rand.Reader)
			ori_value = ori_value.Add(ori_value, new(curve.G1).Neg(shares[i].Share))
		} else {
			shares[i].Share = ori_value
		}
		wg.Add(1)
		go system.OfflineSend(&wg, shares[i].Share.Marshal())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) RandomShareG1() *[]Share_G1 {
	_, r, _ := curve.RandomG1(rand.Reader)
	rshares := system.Share_A_G1_Offline(r)
	return rshares
}

func (system *ShareSystem) share_EXP_P_G1_1(element *curve.G1, xshares Share_Fp) Share_G1 {
	shares := new(Share_G1)
	shares.Index = xshares.Index
	shares.Share = new(curve.G1).ScalarMult(element, xshares.Share)
	return *shares
}

func (system *ShareSystem) EXP_P_G1_1(element *curve.G1, xshares *[]Share_Fp) *[]Share_G1 {
	shares := make([]Share_G1, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_G1_1(element, (*xshares)[i])
	}
	return &shares
}

func (system *ShareSystem) share_EXP_P_G1_2(eshare Share_G1, x *big.Int) Share_G1 {
	shares := new(Share_G1)
	shares.Index = eshare.Index
	shares.Share = new(curve.G1).ScalarMult(eshare.Share, x)
	return *shares
}

func (system *ShareSystem) EXP_P_G1_2(eshares *[]Share_G1, x *big.Int) *[]Share_G1 {
	shares := make([]Share_G1, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_G1_2((*eshares)[i], x)
	}
	return &shares
}

func (system *ShareSystem) shareAdd_G1(shares1, shares2 Share_G1) Share_G1 {
	shares := new(Share_G1)
	shares.Index = shares1.Index
	shares.Share = new(curve.G1).Add(shares1.Share, shares2.Share)
	return *shares
}

func (system *ShareSystem) SecAddPlaintext_G1(shares1 []Share_G1, scalar *curve.G1) *[]Share_G1 {
	shares := make([]Share_G1, system.Partynum)
	shares2 := system.Share_A_G1(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd_G1(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) SecAdd_G1(shares1, shares2 []Share_G1) *[]Share_G1 {
	shares := make([]Share_G1, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd_G1(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) shareSub_G1(shares1, shares2 Share_G1) Share_G1 {
	shares := new(Share_G1)
	shares.Index = shares1.Index
	shares.Share = new(curve.G1).Add(shares1.Share, new(curve.G1).Neg(shares2.Share))
	return *shares
}

func (system *ShareSystem) SecSub_G1(shares1, shares2 []Share_G1) *[]Share_G1 {
	shares := make([]Share_G1, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub_G1(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) SecSubPlaintext_G1(shares1 []Share_G1, scalar *curve.G1) *[]Share_G1 {
	shares := make([]Share_G1, system.Partynum)
	shares2 := system.Share_A_G1(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub_G1(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) EXP_S_G1(hshares []Share_G1, xshares []Share_Fp) *[]Share_G1 {
	sharesA, sharesB, sharesC := system.GenTriplets()
	sharesgB := system.EXP_P_G1_1(curve.Gen1, sharesB)
	sharesgC := system.EXP_P_G1_1(curve.Gen1, sharesC)
	XsubAshares := system.SecSub(xshares, *sharesA)
	xsuba := system.OpenFp(*XsubAshares)
	tshares := system.SecSub_G1(hshares, *sharesgB)
	t := system.OpenG1(*tshares)
	t_exp_xsuba_shares := system.EXP_P_G1_1(t, XsubAshares)
	t_exp_a_shares := system.EXP_P_G1_1(t, sharesA)
	gb_exp_xsuba_shares := system.EXP_P_G1_2(sharesgB, xsuba)
	shares := system.SecAdd_G1(*sharesgC, *gb_exp_xsuba_shares)
	shares = system.SecAdd_G1(*shares, *t_exp_a_shares)
	shares = system.SecAdd_G1(*shares, *t_exp_xsuba_shares)
	return shares
}

func (system *ShareSystem) OpenG1(shares []Share_G1) *curve.G1 {
	var wg sync.WaitGroup
	ori_value := new(curve.G1)
	for i := 0; i < system.Partynum; i++ {
		if i == 0 {
			ori_value = new(curve.G1).Set(shares[i].Share)
		} else {
			ori_value = ori_value.Add(ori_value, shares[i].Share)
		}
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Marshal())
	}
	wg.Wait()
	return ori_value
}
