package shmpc

import (
	"crypto/rand"
	"math/big"
	"sync"
)

type Share_G struct {
	ShareX *big.Int
	ShareY *big.Int
	Index  int
}

func (system *ECCShareSystem) RandomShareG() *[]Share_G {
	rX, rY := system.RandomG()
	rshares := system.Share_A_G_Offline(rX, rY)
	return rshares
}

func (system *ECCShareSystem) RandomG() (*big.Int, *big.Int) {
	scalar, _ := rand.Int(rand.Reader, system.Order)
	RGX, RGY := system.Curve.ScalarMult(system.Curve.Gx, system.Curve.Gy, scalar.Bytes())
	return RGX, RGY
}

func (system *ECCShareSystem) Share_A_G(elementX, elementY *big.Int) *[]Share_G {
	var wg sync.WaitGroup
	ori_valueX := new(big.Int).Set(elementX)
	ori_valueY := new(big.Int).Set(elementY)
	wg.Add(2)
	go system.Send(&wg, ori_valueX.Bytes())
	go system.Send(&wg, ori_valueY.Bytes())
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].ShareX, shares[i].ShareY = system.RandomG()
			ori_valueX, ori_valueY = system.Curve.Add(ori_valueX, ori_valueY, shares[i].ShareX, new(big.Int).Mod(new(big.Int).Neg(shares[i].ShareY), system.Curve.P))
		} else {
			shares[i].ShareX = ori_valueX
			shares[i].ShareY = ori_valueY
		}
		wg.Add(2)
		go system.Send(&wg, shares[i].ShareX.Bytes())
		go system.Send(&wg, shares[i].ShareY.Bytes())
	}
	wg.Wait()
	return &shares
}

func (system *ECCShareSystem) Share_A_G_Offline(elementX, elementY *big.Int) *[]Share_G {
	var wg sync.WaitGroup
	ori_valueX := new(big.Int).Set(elementX)
	ori_valueY := new(big.Int).Set(elementY)
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].ShareX, shares[i].ShareY = system.RandomG()
			ori_valueX, ori_valueY = system.Curve.Add(ori_valueX, ori_valueY, shares[i].ShareX, new(big.Int).Mod(new(big.Int).Neg(shares[i].ShareY), system.Curve.P))
		} else {
			shares[i].ShareX = ori_valueX
			shares[i].ShareY = ori_valueY
		}
		wg.Add(2)
		go system.OfflineSend(&wg, shares[i].ShareX.Bytes())
		go system.OfflineSend(&wg, shares[i].ShareY.Bytes())
	}
	wg.Wait()
	return &shares
}

func (system *ECCShareSystem) share_EXP_P_G_1(elementX, elementY *big.Int, xshares Share_Fp) Share_G {
	shares := new(Share_G)
	shares.Index = xshares.Index
	shares.ShareX, shares.ShareY = system.Curve.ScalarMult(elementX, elementY, xshares.Share.Bytes())
	return *shares
}

func (system *ECCShareSystem) EXP_P_G_1(elementX, elementY *big.Int, xshares *[]Share_Fp) *[]Share_G {
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_G_1(elementX, elementY, (*xshares)[i])
	}
	return &shares
}

func (system *ECCShareSystem) share_EXP_P_G_2(eshare Share_G, x *big.Int) Share_G {
	shares := new(Share_G)
	shares.Index = eshare.Index
	shares.ShareX, shares.ShareY = system.Curve.ScalarMult(eshare.ShareX, eshare.ShareY, x.Bytes())
	return *shares
}

func (system *ECCShareSystem) EXP_P_G_2(eshares *[]Share_G, x *big.Int) *[]Share_G {
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_G_2((*eshares)[i], x)
	}
	return &shares
}

func (system *ECCShareSystem) shareAdd_G(shares1, shares2 Share_G) Share_G {
	shares := new(Share_G)
	shares.Index = shares1.Index
	shares.ShareX, shares.ShareY = system.Curve.Add(shares1.ShareX, shares1.ShareY, shares2.ShareX, shares2.ShareY)
	return *shares
}

func (system *ECCShareSystem) SecAddPlaintext_G(shares1 []Share_G, scalarX, scalarY *big.Int) *[]Share_G {
	shares := make([]Share_G, system.Partynum)
	shares2 := system.Share_A_G_Offline(scalarX, scalarY)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd_G(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ECCShareSystem) SecAdd_G(shares1, shares2 []Share_G) *[]Share_G {
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd_G(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ECCShareSystem) shareSub_G(shares1, shares2 Share_G) Share_G {
	shares := new(Share_G)
	shares.Index = shares1.Index
	shares.ShareX, shares.ShareY = system.Curve.Add(shares1.ShareX, shares1.ShareY, shares2.ShareX, new(big.Int).Mod(new(big.Int).Neg(shares2.ShareY), system.Curve.P))
	return *shares
}

func (system *ECCShareSystem) SecSub_G(shares1, shares2 []Share_G) *[]Share_G {
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub_G(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ECCShareSystem) SecSubPlaintext_G(shares1 []Share_G, scalarX, scalarY *big.Int) *[]Share_G {
	shares := make([]Share_G, system.Partynum)
	shares2 := system.Share_A_G_Offline(scalarX, scalarY)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub_G(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ECCShareSystem) EXP_S_G(hshares []Share_G, xshares []Share_Fp) *[]Share_G {
	sharesA, sharesB, sharesC := system.GenTriplets()
	sharesgB := system.EXP_P_G_1(system.Curve.Gx, system.Curve.Gy, sharesB)
	sharesgC := system.EXP_P_G_1(system.Curve.Gx, system.Curve.Gy, sharesC)
	XsubAshares := system.SecSub(xshares, *sharesA)
	xsuba := system.OpenFp(*XsubAshares)
	tshares := system.SecSub_G(hshares, *sharesgB)
	tx, ty := system.OpenG(*tshares)
	t_exp_xsuba_shares := system.EXP_P_G_1(tx, ty, XsubAshares)
	t_exp_a_shares := system.EXP_P_G_1(tx, ty, sharesA)
	gb_exp_xsuba_shares := system.EXP_P_G_2(sharesgB, xsuba)
	shares := system.SecAdd_G(*sharesgC, *gb_exp_xsuba_shares)
	shares = system.SecAdd_G(*shares, *t_exp_a_shares)
	shares = system.SecAdd_G(*shares, *t_exp_xsuba_shares)
	return shares
}

func (system *ECCShareSystem) OpenG(shares []Share_G) (*big.Int, *big.Int) {
	ori_valueX := new(big.Int)
	ori_valueY := new(big.Int)
	var wg sync.WaitGroup
	for i := 0; i < system.Partynum; i++ {
		if i == 0 {
			ori_valueX = new(big.Int).Set(shares[i].ShareX)
			ori_valueY = new(big.Int).Set(shares[i].ShareY)
		} else {
			ori_valueX, ori_valueY = system.Curve.Add(ori_valueX, ori_valueY, shares[i].ShareX, shares[i].ShareY)
		}
		wg.Add(2)
		go system.Broadcast(&wg, shares[i].ShareX.Bytes())
		go system.Broadcast(&wg, shares[i].ShareY.Bytes())
	}
	wg.Wait()
	return ori_valueX, ori_valueY
}

func (system *ECCShareSystem) Share_An_Fp(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	wg.Add(1)
	go system.Send(&wg, ori_value.Bytes())
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].Share, _ = rand.Int(rand.Reader, system.Order)
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

func (system *ECCShareSystem) Share_An_Fp_Offline(element *big.Int) *[]Share_Fp {
	ori_value := new(big.Int).Set(element)
	var wg sync.WaitGroup
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		if i < system.Partynum-1 {
			shares[i].Share, _ = rand.Int(rand.Reader, system.Order)
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

func (system *ECCShareSystem) RandomShareFp() *[]Share_Fp {
	r, _ := rand.Int(rand.Reader, system.Order)
	rshares := system.Share_An_Fp_Offline(r)
	return rshares
}

func (system *ECCShareSystem) shareAdd(shares1, shares2 Share_Fp) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	shares.Share = new(big.Int).Add(shares1.Share, shares2.Share)
	shares.Share = shares.Share.Mod(shares.Share, system.Order)
	return *shares
}

func (system *ECCShareSystem) SecAddPlaintext(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	shares2 := system.Share_An_Fp(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ECCShareSystem) SecAdd(shares1, shares2 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareAdd(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ECCShareSystem) SecSubPlaintext(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	shares2 := system.Share_An_Fp(scalar)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub(shares1[i], (*shares2)[i])
	}
	return &shares
}

func (system *ECCShareSystem) shareSub(shares1, shares2 Share_Fp) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	shares.Share = new(big.Int).Sub(shares1.Share, shares2.Share)
	shares.Share = shares.Share.Mod(shares.Share, system.Order)
	return *shares
}

func (system *ECCShareSystem) SecSub(shares1, shares2 []Share_Fp) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareSub(shares1[i], shares2[i])
	}
	return &shares
}

func (system *ECCShareSystem) shareMulPlaintext(shares1 Share_Fp, scalar *big.Int) Share_Fp {
	shares := new(Share_Fp)
	shares.Index = shares1.Index
	shares.Share = new(big.Int).Mul(shares1.Share, scalar)
	shares.Share = shares.Share.Mod(shares.Share, system.Order)
	return *shares
}

func (system *ECCShareSystem) SecMulPlaintext(shares1 []Share_Fp, scalar *big.Int) *[]Share_Fp {
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.shareMulPlaintext(shares1[i], scalar)
	}
	return &shares
}

func (system *ECCShareSystem) SecMul(shares1, shares2 []Share_Fp) *[]Share_Fp {
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

func (system *ECCShareSystem) SecSquare(shares1 []Share_Fp) *[]Share_Fp {
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

func (system *ECCShareSystem) OpenFp(shares []Share_Fp) *big.Int {
	ori_value := big.NewInt(0)
	var wg sync.WaitGroup
	for i := 0; i < system.Partynum; i++ {
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Bytes())
		ori_value = ori_value.Add(ori_value, shares[i].Share)
	}
	wg.Wait()
	ori_value = ori_value.Mod(ori_value, system.Order)
	return ori_value
}
