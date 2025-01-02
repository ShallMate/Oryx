package mpc

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"sync"

	curve "github.com/Oryx/curve"
)

type Share_G2 struct {
	Share *curve.G2
	Gama  *curve.G2
	Delta *curve.G2
	Index int
}

func (system *ShareSystem) Share_A_G2(element *curve.G2) *[]Share_G2 {
	var wg sync.WaitGroup
	ori_value := new(curve.G2).Set(element)
	_, Delta, _ := curve.RandomG2(rand.Reader)
	wg.Add(2)
	go system.Send(&wg, ori_value.Marshal())
	go system.BroadcastN(&wg, Delta.Marshal())
	Gama := new(curve.G2).Add(ori_value, Delta)
	Gama = Gama.ScalarMult(Gama, system.alpha)
	shares := make([]Share_G2, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			_, shares[i].Share, _ = curve.RandomG2(rand.Reader)
			_, shares[i].Gama, _ = curve.RandomG2(rand.Reader)
			ori_value = ori_value.Add(ori_value, new(curve.G2).Neg(shares[i].Share))
			Gama = Gama.Add(Gama, new(curve.G2).Neg(shares[i].Gama))
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

func (system *ShareSystem) Share_A_G2_Offline(element *curve.G2) *[]Share_G2 {
	var wg sync.WaitGroup
	ori_value := new(curve.G2).Set(element)
	_, Delta, _ := curve.RandomG2(rand.Reader)
	wg.Add(1)
	go system.OfflineBroadcastN(&wg, Delta.Marshal())
	Gama := new(curve.G2).Add(ori_value, Delta)
	Gama = Gama.ScalarMult(Gama, system.alpha)
	shares := make([]Share_G2, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			_, shares[i].Share, _ = curve.RandomG2(rand.Reader)
			_, shares[i].Gama, _ = curve.RandomG2(rand.Reader)
			ori_value = ori_value.Add(ori_value, new(curve.G2).Neg(shares[i].Share))
			Gama = Gama.Add(Gama, new(curve.G2).Neg(shares[i].Gama))
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

func (system *ShareSystem) RandomShareG2() *[]Share_G2 {
	_, r, _ := curve.RandomG2(rand.Reader)
	rshares := system.Share_A_G2_Offline(r)
	return rshares
}

func (system *ShareSystem) share_EXP_P_G2_1(element *curve.G2, xshares Share_Fp) Share_G2 {
	shares := new(Share_G2)
	shares.Index = xshares.Index
	Delta := new(curve.G2).ScalarMult(element, xshares.Delta)
	shares.Delta = Delta
	shares.Gama = new(curve.G2).ScalarMult(element, xshares.Gama)
	shares.Share = new(curve.G2).ScalarMult(element, xshares.Share)
	return *shares
}

func (system *ShareSystem) EXP_P_G2_1(element *curve.G2, xshares *[]Share_Fp) *[]Share_G2 {
	shares := make([]Share_G2, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_G2_1(element, (*xshares)[i])
	}
	return &shares
}

func (system *ShareSystem) share_EXP_P_G2_2(eshare Share_G2, x *big.Int) Share_G2 {
	shares := new(Share_G2)
	shares.Index = eshare.Index
	Delta := new(curve.G2).ScalarMult(eshare.Delta, x)
	shares.Delta = Delta
	shares.Gama = new(curve.G2).ScalarMult(eshare.Gama, x)
	shares.Share = new(curve.G2).ScalarMult(eshare.Share, x)
	return *shares
}

func (system *ShareSystem) EXP_P_G2_2(eshares *[]Share_G2, x *big.Int) *[]Share_G2 {
	shares := make([]Share_G2, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_G2_2((*eshares)[i], x)
	}
	return &shares
}

func (system *ShareSystem) shareAdd_G2(shares1, shares2 Share_G2) Share_G2 {
	shares := new(Share_G2)
	shares.Index = shares1.Index
	Delta := new(curve.G2).Add(shares1.Delta, shares2.Delta)
	shares.Delta = Delta
	shares.Gama = new(curve.G2).Add(shares1.Gama, shares2.Gama)
	shares.Share = new(curve.G2).Add(shares1.Share, shares2.Share)
	return *shares
}

func (system *ShareSystem) SecAddPlaintext_G2(shares1 []Share_G2, scalar *curve.G2) *[]Share_G2 {
	shares := make([]Share_G2, system.Partynum)
	shares2 := system.Share_A_G2(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd_G2(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) SecAdd_G2(shares1, shares2 []Share_G2) *[]Share_G2 {
	shares := make([]Share_G2, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd_G2(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) shareSub_G2(shares1, shares2 Share_G2) Share_G2 {
	shares := new(Share_G2)
	shares.Index = shares1.Index
	Delta := new(curve.G2).Add(shares1.Delta, new(curve.G2).Neg(shares2.Delta))
	shares.Delta = Delta
	shares.Gama = new(curve.G2).Add(shares1.Gama, new(curve.G2).Neg(shares2.Gama))
	shares.Share = new(curve.G2).Add(shares1.Share, new(curve.G2).Neg(shares2.Share))
	return *shares
}

func (system *ShareSystem) SecSub_G2(shares1, shares2 []Share_G2) *[]Share_G2 {
	shares := make([]Share_G2, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub_G2(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ShareSystem) SecSubPlaintext_G2(shares1 []Share_G2, scalar *curve.G2) *[]Share_G2 {
	shares := make([]Share_G2, system.Partynum)
	shares2 := system.Share_A_G2(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub_G2(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ShareSystem) EXP_S_G2(hshares []Share_G2, xshares []Share_Fp) *[]Share_G2 {
	sharesA, sharesB, sharesC := system.GenTriplets()
	sharesgB := system.EXP_P_G2_1(curve.Gen2, sharesB)
	sharesgC := system.EXP_P_G2_1(curve.Gen2, sharesC)
	XsubAshares := system.SecSub(xshares, *sharesA)
	xsuba := system.HalfOpenFp(*XsubAshares)
	tshares := system.SecSub_G2(hshares, *sharesgB)
	t := system.HalfOpenG2(*tshares)
	t_exp_xsuba_shares := system.EXP_P_G2_1(t, XsubAshares)
	t_exp_a_shares := system.EXP_P_G2_1(t, sharesA)
	gb_exp_xsuba_shares := system.EXP_P_G2_2(sharesgB, xsuba)
	shares := system.SecAdd_G2(*sharesgC, *gb_exp_xsuba_shares)
	shares = system.SecAdd_G2(*shares, *t_exp_a_shares)
	shares = system.SecAdd_G2(*shares, *t_exp_xsuba_shares)
	return shares
}

func (system *ShareSystem) HalfOpenG2(shares []Share_G2) *curve.G2 {
	var wg sync.WaitGroup
	ori_value := new(curve.G2)
	for i := 0; i < system.Partynum; i++ {
		if i == 0 {
			ori_value = new(curve.G2).Set(shares[i].Share)
		} else {
			ori_value = ori_value.Add(ori_value, shares[i].Share)
		}
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Marshal())
	}
	wg.Wait()
	return ori_value
}

func (system *ShareSystem) MacCheckG2(shares []Share_G2, res_value *curve.G2) bool {
	var wg sync.WaitGroup
	chk := new(curve.G2).Set(system.IdentityG2)
	delta := new(curve.G2)
	t := new(curve.G2).Add(res_value, shares[0].Delta)
	for i := 0; i < system.Partynum; i++ {
		delta = delta.ScalarMult(t, system.Alphas[i])
		delta = delta.Add(shares[i].Gama, new(curve.G2).Neg(delta))
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
	wg.Wait()
	return bytes.Equal(chk.Marshal(), system.IdentityG2.Marshal())
}

func (system *ShareSystem) OpenG2(shares []Share_G2) (*curve.G2, bool) {
	ori_value := new(curve.G2)
	var wg sync.WaitGroup
	for i := 0; i < system.Partynum; i++ {
		if i == 0 {
			ori_value = new(curve.G2).Set(shares[i].Share)
		} else {
			ori_value = ori_value.Add(ori_value, shares[i].Share)
		}
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Marshal())
	}
	wg.Wait()
	chk := system.MacCheckG2(shares, ori_value)
	return ori_value, chk
}
