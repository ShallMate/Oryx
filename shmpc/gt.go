package shmpc

import (
	"crypto/rand"
	"math/big"
	"sync"

	curve "github.com/Oryx/curve"
)

type Share_GT struct {
	Share *curve.GT
	Index int
}

func (system *ShareSystem) Share_A_GT(element *curve.GT) *[]Share_GT {
	var wg sync.WaitGroup
	ori_value := new(curve.GT).Set(element)
	wg.Add(1)
	go system.Send(&wg, ori_value.Marshal())
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			_, shares[i].Share, _ = curve.RandomGTK(rand.Reader)
			ori_value = ori_value.Add(ori_value, new(curve.GT).Neg(shares[i].Share))
		} else {
			shares[i].Share = ori_value
		}
		wg.Add(1)
		go system.Send(&wg, shares[i].Share.Marshal())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) Share_A_GT_Offline(element *curve.GT) *[]Share_GT {
	var wg sync.WaitGroup
	ori_value := new(curve.GT).Set(element)
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			_, shares[i].Share, _ = curve.RandomGTK(rand.Reader)
			ori_value = ori_value.Add(ori_value, new(curve.GT).Neg(shares[i].Share))
		} else {
			shares[i].Share = ori_value
		}
		wg.Add(1)
		go system.OfflineSend(&wg, shares[i].Share.Marshal())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) RandomShareGT() *[]Share_GT {
	_, r, _ := curve.RandomGTK(rand.Reader)
	rshares := system.Share_A_GT_Offline(r)
	return rshares
}

func (system *ShareSystem) share_EXP_P_GT_1(element *curve.GT, xshares Share_Fp) Share_GT {
	shares := new(Share_GT)
	shares.Index = xshares.Index
	shares.Share = new(curve.GT).ScalarMult(element, xshares.Share)
	return *shares
}

func (system *ShareSystem) EXP_P_GT_1(element *curve.GT, xshares *[]Share_Fp) *[]Share_GT {
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_GT_1(element, (*xshares)[i])
	}
	return &shares
}

func (system *ShareSystem) share_EXP_P_GT_2(eshare Share_GT, x *big.Int) Share_GT {
	shares := new(Share_GT)
	shares.Index = eshare.Index
	shares.Share = new(curve.GT).ScalarMult(eshare.Share, x)
	return *shares
}

func (system *ShareSystem) EXP_P_GT_2(eshares *[]Share_GT, x *big.Int) *[]Share_GT {
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_GT_2((*eshares)[i], x)
	}
	return &shares
}

func (system *ShareSystem) shareAdd_GT(shares1, shares2 Share_GT) Share_GT {
	shares := new(Share_GT)
	shares.Index = shares1.Index
	shares.Share = new(curve.GT).Add(shares1.Share, shares2.Share)
	return *shares
}

func (system *ShareSystem) SecAddPlaintext_GT(shares1 []Share_GT, scalar *curve.GT) *[]Share_GT {
	shares := make([]Share_GT, system.Partynum)
	shares2 := system.Share_A_GT(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd_GT(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) SecAdd_GT(shares1, shares2 []Share_GT) *[]Share_GT {
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd_GT(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) shareSub_GT(shares1, shares2 Share_GT) Share_GT {
	shares := new(Share_GT)
	shares.Index = shares1.Index
	shares.Share = new(curve.GT).Add(shares1.Share, new(curve.GT).Neg(shares2.Share))
	return *shares
}

func (system *ShareSystem) SecSub_GT(shares1, shares2 []Share_GT) *[]Share_GT {
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub_GT(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) SecSubPlaintext_GT(shares1 []Share_GT, scalar *curve.GT) *[]Share_GT {
	shares := make([]Share_GT, system.Partynum)
	shares2 := system.Share_A_GT(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub_GT(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) EXP_S_GT(hshares []Share_GT, xshares []Share_Fp) *[]Share_GT {
	sharesA, sharesB, sharesC := system.GenTriplets()
	sharesgB := system.EXP_P_GT_1(system.GenGT, sharesB)
	sharesgC := system.EXP_P_GT_1(system.GenGT, sharesC)
	XsubAshares := system.SecSub(xshares, *sharesA)
	xsuba := system.OpenFp(*XsubAshares)
	tshares := system.SecSub_GT(hshares, *sharesgB)
	t := system.OpenGT(*tshares)
	t_exp_xsuba_shares := system.EXP_P_GT_1(t, XsubAshares)
	t_exp_a_shares := system.EXP_P_GT_1(t, sharesA)
	gb_exp_xsuba_shares := system.EXP_P_GT_2(sharesgB, xsuba)
	shares := system.SecAdd_GT(*sharesgC, *gb_exp_xsuba_shares)
	shares = system.SecAdd_GT(*shares, *t_exp_a_shares)
	shares = system.SecAdd_GT(*shares, *t_exp_xsuba_shares)
	return shares
}

func (system *ShareSystem) OpenGT(shares []Share_GT) *curve.GT {
	ori_value := new(curve.GT)
	var wg sync.WaitGroup
	for i := 0; i < system.Partynum; i++ {
		if i == 0 {
			ori_value = new(curve.GT).Set(shares[i].Share)
		} else {
			ori_value = ori_value.Add(ori_value, shares[i].Share)
		}
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Marshal())
	}
	wg.Wait()
	return ori_value
}
