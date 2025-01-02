package mpc

import (
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/Oryx/curve"
)

func (system *ShareSystem) Share_An_Fp_Mul(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	Delta, _ := rand.Int(rand.Reader, system.Order)
	wg.Add(2)
	go system.Send(&wg, ori_value.Bytes())
	go system.BroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Mul(ori_value, Delta)
	Gama = Gama.Mod(Gama, system.Order)
	Gama = Gama.Exp(Gama, system.alpha, system.Order)
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = rand.Int(rand.Reader, system.Order)
			shares[i].Gama, _ = rand.Int(rand.Reader, system.Order)
			shareinv := new(big.Int).ModInverse(shares[i].Share, system.Order)
			ori_value = ori_value.Mul(ori_value, shareinv)
			ori_value = ori_value.Mod(ori_value, system.Order)
			gamainv := new(big.Int).ModInverse(shares[i].Gama, system.Order)
			Gama = Gama.Mul(Gama, gamainv)
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

func (system *ShareSystem) Share_An_Fp_Mul_Offline(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	Delta, _ := rand.Int(rand.Reader, system.Order)
	wg.Add(1)
	go system.OfflineBroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Mul(ori_value, Delta)
	Gama = Gama.Exp(Gama, system.alpha, system.Order)
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = rand.Int(rand.Reader, system.Order)
			shares[i].Gama, _ = rand.Int(rand.Reader, system.Order)
			shareinv := new(big.Int).ModInverse(shares[i].Share, system.Order)
			ori_value = ori_value.Mul(ori_value, shareinv)
			ori_value = ori_value.Mod(ori_value, system.Order)
			gamainv := new(big.Int).ModInverse(shares[i].Gama, system.Order)
			Gama = Gama.Mul(Gama, gamainv)
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

func (system *ShareSystem) Share_An_Fp_for_EXP(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	Delta, _ := curve.RandomK(rand.Reader)
	wg.Add(2)
	go system.Send(&wg, ori_value.Bytes())
	go system.BroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Add(ori_value, Delta)
	Gama = Gama.Mul(system.alpha, Gama)
	Gama = Gama.Mod(Gama, system.OrderMul)
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = curve.RandomK(rand.Reader)
			shares[i].Gama, _ = curve.RandomK(rand.Reader)
			ori_value = ori_value.Sub(ori_value, shares[i].Share)
			ori_value = ori_value.Mod(ori_value, system.OrderMul)
			Gama = Gama.Sub(Gama, shares[i].Gama)
			Gama = Gama.Mod(Gama, system.OrderMul)
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

func (system *ShareSystem) Share_An_Fp_for_EXP_Offline(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	Delta, _ := curve.RandomK(rand.Reader)
	wg.Add(1)
	go system.OfflineBroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Add(ori_value, Delta)
	Gama = Gama.Mul(system.alpha, Gama)
	Gama = Gama.Mod(Gama, system.OrderMul)
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = curve.RandomK(rand.Reader)
			shares[i].Gama, _ = curve.RandomK(rand.Reader)
			ori_value = ori_value.Sub(ori_value, shares[i].Share)
			ori_value = ori_value.Mod(ori_value, system.OrderMul)
			Gama = Gama.Sub(Gama, shares[i].Gama)
			Gama = Gama.Mod(Gama, system.OrderMul)
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

func (system *ShareSystem) RandomShareFp_Mul() *[]Share_Fp {
	r, _ := curve.RandomK(rand.Reader)
	rshares := system.Share_An_Fp_Mul_Offline(r)
	return rshares
}

func (system *ShareSystem) shareMul(shares1, shares2 Share_Fp) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	Delta := new(big.Int).Mul(shares1.Delta, shares2.Delta)
	Delta = Delta.Mod(Delta, system.Order)
	shares.Delta = Delta
	shares.Gama = new(big.Int).Mul(shares1.Gama, shares2.Gama)
	shares.Gama = shares.Gama.Mod(shares.Gama, system.Order)
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
	Delta := new(big.Int).Mul(shares1.Delta, new(big.Int).ModInverse(shares2.Delta, system.Order))
	Delta = Delta.Mod(Delta, system.Order)
	shares.Delta = Delta
	shares.Gama = new(big.Int).Mul(shares1.Gama, new(big.Int).ModInverse(shares2.Gama, system.Order))
	shares.Gama = shares.Gama.Mod(shares.Gama, system.Order)
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
		shares[i].Delta = new(big.Int).Exp(scalar, shares1[i].Delta, system.Order)
		shares[i].Gama = new(big.Int).Exp(scalar, shares1[i].Gama, system.Order)
		shares[i].Share = new(big.Int).Exp(scalar, shares1[i].Share, system.Order)
	}
	return &shares
}

func (system *ShareSystem) EXP_P_Fp_2(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = new(big.Int).Exp(shares1[i].Delta, scalar, system.Order)
		shares[i].Gama = new(big.Int).Exp(shares1[i].Gama, scalar, system.Order)
		shares[i].Share = new(big.Int).Exp(shares1[i].Share, scalar, system.Order)
	}
	return &shares
}

/*
func (system *ShareSystem) EXP_S_Fp(hshares, xshares []Share_Fp) *[]Share_Fp {
	sharesA, sharesB, sharesC := system.GenTriplets_for_Exp()
	G, _ := curve.RandomK(rand.Reader)
	sharesgB := system.EXP_P_Fp_1(G, *sharesB)
	sharesgC := system.EXP_P_Fp_1(G, *sharesC)
	XsubAshares := system.SecSub__for_EXP(xshares, *sharesA)
	xsuba := system.HalfOpenFp_for_Exp(*XsubAshares)
	tshares := system.SecDiv(hshares, *sharesgB)
	t := system.HalfOpenFp_Mul(*tshares)
	t_exp_xsuba_shares := system.EXP_P_Fp_1(t, *XsubAshares)
	t_exp_a_shares := system.EXP_P_Fp_1(t, *sharesA)
	gb_exp_xsuba_shares := system.EXP_P_Fp_2(*sharesgB, xsuba)
	shares := system.SecMul_Mul(*sharesgC, *gb_exp_xsuba_shares)
	shares = system.SecMul_Mul(*shares, *t_exp_a_shares)
	shares = system.SecMul_Mul(*shares, *t_exp_xsuba_shares)
	return shares
}
*/

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

func (system *ShareSystem) HalfOpenFp_Mul(shares []Share_Fp) *big.Int {
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

func (system *ShareSystem) MacCheckFp_Mul(shares []Share_Fp, res_value *big.Int) bool {
	var wg sync.WaitGroup
	chk := big.NewInt(1)
	delta := new(big.Int)
	t := new(big.Int).Mul(res_value, shares[0].Delta)
	t = t.Mod(t, system.Order)
	for i := 0; i < system.Partynum; i++ {
		delta = delta.Exp(t, system.AlphasMul[i], system.Order)
		delta = delta.ModInverse(delta, system.Order)
		delta = delta.Mul(shares[i].Gama, delta)
		delta = delta.Mod(delta, system.Order)
		commit, r := Com(delta.Bytes())
		wg.Add(3)
		go system.Broadcast(&wg, commit)
		go system.Broadcast(&wg, r.Bytes())
		go system.Broadcast(&wg, delta.Bytes())
		opencommit := OpenComit(delta.Bytes(), commit, r)
		if !opencommit {
			return false
		}
		chk = chk.Mul(chk, delta)
	}
	chk = chk.Mod(chk, system.Order)
	wg.Wait()
	return chk.Cmp(one) == 0

}

func (system *ShareSystem) OpenFp_Mul(shares []Share_Fp) (*big.Int, bool) {
	var wg sync.WaitGroup
	ori_value := big.NewInt(1)
	for i := 0; i < system.Partynum; i++ {
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Bytes())
		ori_value = ori_value.Mul(ori_value, shares[i].Share)
	}
	ori_value = ori_value.Mod(ori_value, system.Order)
	wg.Wait()
	chk := system.MacCheckFp_Mul(shares, ori_value)
	return ori_value, chk
}
