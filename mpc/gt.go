package mpc

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"sync"

	curve "github.com/Oryx/curve"
)

type Share_GT struct {
	Share *curve.GT
	Gama  *curve.GT
	Delta *curve.GT
	Index int
}

func (system *ShareSystem) Share_A_GT(element *curve.GT) *[]Share_GT {
	var wg sync.WaitGroup
	ori_value := new(curve.GT).Set(element)
	_, Delta, _ := curve.RandomGTK(rand.Reader)
	wg.Add(2)
	go system.Send(&wg, ori_value.Marshal())
	go system.BroadcastN(&wg, Delta.Marshal())
	Gama := new(curve.GT).Add(ori_value, Delta)
	Gama = Gama.ScalarMult(Gama, system.alpha)
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			_, shares[i].Share, _ = curve.RandomGTK(rand.Reader)
			_, shares[i].Gama, _ = curve.RandomGTK(rand.Reader)
			ori_value = ori_value.Add(ori_value, new(curve.GT).Neg(shares[i].Share))
			Gama = Gama.Add(Gama, new(curve.GT).Neg(shares[i].Gama))
		} else {
			shares[i].Share = ori_value
			shares[i].Gama = Gama
		}
		wg.Add(2)
		go system.Send(&wg, shares[i].Share.Marshal())
		go system.Send(&wg, shares[i].Gama.Marshal())
	}
	wg.Wait()
	return &shares
}

func (system *ShareSystem) Share_A_GT_Offline(element *curve.GT) *[]Share_GT {
	var wg sync.WaitGroup
	ori_value := new(curve.GT).Set(element)
	_, Delta, _ := curve.RandomGTK(rand.Reader)
	wg.Add(1)
	go system.OfflineBroadcastN(&wg, Delta.Marshal())
	Gama := new(curve.GT).Add(ori_value, Delta)
	Gama = Gama.ScalarMult(Gama, system.alpha)
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			_, shares[i].Share, _ = curve.RandomGTK(rand.Reader)
			_, shares[i].Gama, _ = curve.RandomGTK(rand.Reader)
			ori_value = ori_value.Add(ori_value, new(curve.GT).Neg(shares[i].Share))
			Gama = Gama.Add(Gama, new(curve.GT).Neg(shares[i].Gama))
		} else {
			shares[i].Share = ori_value
			shares[i].Gama = Gama
		}
		wg.Add(2)
		go system.OfflineSend(&wg, shares[i].Share.Marshal())
		go system.OfflineSend(&wg, shares[i].Gama.Marshal())
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
	Delta := new(curve.GT).ScalarMult(element, xshares.Delta)
	shares.Delta = Delta
	shares.Gama = new(curve.GT).ScalarMult(element, xshares.Gama)
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
	Delta := new(curve.GT).ScalarMult(eshare.Delta, x)
	shares.Delta = Delta
	shares.Gama = new(curve.GT).ScalarMult(eshare.Gama, x)
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
	Delta := new(curve.GT).Add(shares1.Delta, shares2.Delta)
	shares.Delta = Delta
	shares.Gama = new(curve.GT).Add(shares1.Gama, shares2.Gama)
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
	Delta := new(curve.GT).Add(shares1.Delta, new(curve.GT).Neg(shares2.Delta))
	shares.Delta = Delta
	shares.Gama = new(curve.GT).Add(shares1.Gama, new(curve.GT).Neg(shares2.Gama))
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
	xsuba := system.HalfOpenFp(*XsubAshares)
	tshares := system.SecSub_GT(hshares, *sharesgB)
	t := system.HalfOpenGT(*tshares)
	t_exp_xsuba_shares := system.EXP_P_GT_1(t, XsubAshares)
	t_exp_a_shares := system.EXP_P_GT_1(t, sharesA)
	gb_exp_xsuba_shares := system.EXP_P_GT_2(sharesgB, xsuba)
	shares := system.SecAdd_GT(*sharesgC, *gb_exp_xsuba_shares)
	shares = system.SecAdd_GT(*shares, *t_exp_a_shares)
	shares = system.SecAdd_GT(*shares, *t_exp_xsuba_shares)
	return shares
}

func (system *ShareSystem) HalfOpenGT(shares []Share_GT) *curve.GT {
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

func (system *ShareSystem) MacCheckGT(shares []Share_GT, res_value *curve.GT) bool {
	var wg sync.WaitGroup
	chk := new(curve.GT).Set(system.IdentityGT)
	delta := new(curve.GT)
	t := new(curve.GT).Add(res_value, shares[0].Delta)
	for i := 0; i < system.Partynum; i++ {
		delta = delta.ScalarMult(t, system.Alphas[i])
		delta = delta.Add(shares[i].Gama, new(curve.GT).Neg(delta))
		commit, r := Com(delta.Marshal())
		wg.Add(3)
		go system.Broadcast(&wg, delta.Marshal())
		go system.Broadcast(&wg, commit)
		go system.Broadcast(&wg, r.Bytes())
		opencommit := OpenComit(delta.Marshal(), commit, r)
		if !opencommit {
			return false
		}
		chk = chk.Add(chk, delta)
	}
	return bytes.Equal(chk.Marshal(), system.IdentityGT.Marshal())
}

func (system *ShareSystem) OpenGT(shares []Share_GT) (*curve.GT, bool) {
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
	chk := system.MacCheckGT(shares, ori_value)
	return ori_value, chk
}
