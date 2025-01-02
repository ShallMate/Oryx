package mpc

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"sync"
)

// Share_G represents a share of a point on an elliptic curve.
type Share_G struct {
	ShareX *big.Int // ShareX is the x-coordinate of the share.
	ShareY *big.Int // ShareY is the y-coordinate of the share.
	GamaX  *big.Int // GamaX is the x-coordinate of the Gama point.
	GamaY  *big.Int // GamaY is the y-coordinate of the Gama point.
	DeltaX *big.Int // DeltaX is the x-coordinate of the Delta point.
	DeltaY *big.Int // DeltaY is the y-coordinate of the Delta point.
	Index  int      // Index is the index of the share.
}

// RandomShareG generates a random share of G.
// It generates random coordinates (rX, rY) using the RandomG function,
// and then computes the shares using the Share_A_G_Offline function.
// The generated shares are returned as a pointer to a slice of Share_G.
func (system *ECCShareSystem) RandomShareG() *[]Share_G {
	rX, rY := system.RandomG()
	rshares := system.Share_A_G_Offline(rX, rY)
	return rshares
}

// RandomG generates a random point on the elliptic curve using the ECCShareSystem.
// It returns the x and y coordinates of the generated point.
func (system *ECCShareSystem) RandomG() (*big.Int, *big.Int) {
	scalar, _ := rand.Int(rand.Reader, system.Order)
	RGX, RGY := system.Curve.ScalarMult(system.Curve.Gx, system.Curve.Gy, scalar.Bytes())
	return RGX, RGY
}

// Share_A_G performs the sharing of an elliptic curve point (elementX, elementY) among multiple parties in the ECCShareSystem.
// It generates random values DeltaX and DeltaY, broadcasts them to all parties, and computes GamaX and GamaY using the curve addition operation.
// It then generates shares for each party, where each share consists of DeltaX, DeltaY, ShareX, ShareY, GamaX, and GamaY.
// The shares are sent to each party using the Send method of the ECCShareSystem.
// Finally, it waits for all the shares to be sent and returns a pointer to the array of shares.
func (system *ECCShareSystem) Share_A_G(elementX, elementY *big.Int) *[]Share_G {
	var wg sync.WaitGroup
	ori_valueX := new(big.Int).Set(elementX)
	ori_valueY := new(big.Int).Set(elementY)
	wg.Add(2)
	go system.Send(&wg, ori_valueX.Bytes())
	go system.Send(&wg, ori_valueY.Bytes())
	DeltaX, DeltaY := system.RandomG()
	wg.Add(2)
	go system.BroadcastN(&wg, DeltaX.Bytes())
	go system.BroadcastN(&wg, DeltaY.Bytes())
	GamaX, GamaY := system.Curve.Add(ori_valueX, ori_valueY, DeltaX, DeltaY)
	GamaX, GamaY = system.Curve.ScalarMult(GamaX, GamaY, system.alpha.Bytes())
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].DeltaX = DeltaX
		shares[i].DeltaY = DeltaY
		if i < system.Partynum-1 {
			shares[i].ShareX, shares[i].ShareY = system.RandomG()
			shares[i].GamaX, shares[i].GamaY = system.RandomG()
			ori_valueX, ori_valueY = system.Curve.Add(ori_valueX, ori_valueY, shares[i].ShareX, new(big.Int).Mod(new(big.Int).Neg(shares[i].ShareY), system.Curve.P))
			GamaX, GamaY = system.Curve.Add(GamaX, GamaY, shares[i].GamaX, new(big.Int).Mod(new(big.Int).Neg(shares[i].GamaY), system.Curve.P))
		} else {
			shares[i].ShareX = ori_valueX
			shares[i].ShareY = ori_valueY
			shares[i].GamaX = GamaX
			shares[i].GamaY = GamaY
		}
		wg.Add(4)
		go system.Send(&wg, shares[i].ShareX.Bytes())
		go system.Send(&wg, shares[i].ShareY.Bytes())
		go system.Send(&wg, shares[i].GamaX.Bytes())
		go system.Send(&wg, shares[i].GamaY.Bytes())
	}
	wg.Wait()
	return &shares
}

// Share_A_G_Offline performs offline sharing of an elliptic curve point (elementX, elementY) among multiple parties.
// It generates random values DeltaX and DeltaY, and broadcasts them to all parties.
// It then computes GamaX and GamaY by adding (elementX, elementY) with DeltaX and DeltaY, and multiplies the result with the system's alpha value.
// Finally, it generates shares for each party, where each share consists of DeltaX, DeltaY, ShareX, ShareY, GamaX, and GamaY.
// The shares are sent to the respective parties using the OfflineSend method.
// The function returns a pointer to the array of shares.
func (system *ECCShareSystem) Share_A_G_Offline(elementX, elementY *big.Int) *[]Share_G {
	var wg sync.WaitGroup
	ori_valueX := new(big.Int).Set(elementX)
	ori_valueY := new(big.Int).Set(elementY)
	DeltaX, DeltaY := system.RandomG()
	wg.Add(2)
	go system.OfflineBroadcastN(&wg, DeltaX.Bytes())
	go system.OfflineBroadcastN(&wg, DeltaY.Bytes())
	GamaX, GamaY := system.Curve.Add(ori_valueX, ori_valueY, DeltaX, DeltaY)
	GamaX, GamaY = system.Curve.ScalarMult(GamaX, GamaY, system.alpha.Bytes())
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].DeltaX = DeltaX
		shares[i].DeltaY = DeltaY
		if i < system.Partynum-1 {
			shares[i].ShareX, shares[i].ShareY = system.RandomG()
			shares[i].GamaX, shares[i].GamaY = system.RandomG()
			ori_valueX, ori_valueY = system.Curve.Add(ori_valueX, ori_valueY, shares[i].ShareX, new(big.Int).Mod(new(big.Int).Neg(shares[i].ShareY), system.Curve.P))
			GamaX, GamaY = system.Curve.Add(GamaX, GamaY, shares[i].GamaX, new(big.Int).Mod(new(big.Int).Neg(shares[i].GamaY), system.Curve.P))
		} else {
			shares[i].ShareX = ori_valueX
			shares[i].ShareY = ori_valueY
			shares[i].GamaX = GamaX
			shares[i].GamaY = GamaY
		}
		wg.Add(4)
		go system.OfflineSend(&wg, shares[i].ShareX.Bytes())
		go system.OfflineSend(&wg, shares[i].ShareY.Bytes())
		go system.OfflineSend(&wg, shares[i].GamaX.Bytes())
		go system.OfflineSend(&wg, shares[i].GamaY.Bytes())
	}
	wg.Wait()
	return &shares
}

// share_EXP_P_G_1 calculates and returns the share of the point (elementX, elementY) in the ECCShareSystem.
// It takes an elementX and elementY as input, which are the coordinates of the point to be shared.
// The xshares parameter is of type Share_Fp and contains the shares of the secret value x.
// The function returns a Share_G struct that contains the shares of the point (elementX, elementY) in the ECCShareSystem.
func (system *ECCShareSystem) share_EXP_P_G_1(elementX, elementY *big.Int, xshares Share_Fp) Share_G {
	shares := new(Share_G)
	shares.Index = xshares.Index
	DeltaX, DeltaY := system.Curve.ScalarMult(elementX, elementY, xshares.Delta.Bytes())
	shares.DeltaX = DeltaX
	shares.DeltaY = DeltaY
	shares.GamaX, shares.GamaY = system.Curve.ScalarMult(elementX, elementY, xshares.Gama.Bytes())
	shares.ShareX, shares.ShareY = system.Curve.ScalarMult(elementX, elementY, xshares.Share.Bytes())
	return *shares
}

// EXP_P_G_1 calculates the shares of the point multiplication of a point (elementX, elementY) with a scalar value for each party.
// It takes the xshares as input, which are the shares of the scalar value for each party.
// It returns the shares of the resulting point multiplication for each party.
func (system *ECCShareSystem) EXP_P_G_1(elementX, elementY *big.Int, xshares *[]Share_Fp) *[]Share_G {
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_G_1(elementX, elementY, (*xshares)[i])
	}
	return &shares
}

// share_EXP_P_G_2 performs the exponentiation of a point G in the elliptic curve
// and returns a new Share_G structure with the computed shares.
// The function takes an input Share_G structure, `eshare`, and a big.Int value, `x`,
// which is used as the exponent for the scalar multiplication.
// The function computes the scalar multiplication of `eshare.DeltaX` and `eshare.DeltaY`
// with `x` as the scalar, and assigns the result to `shares.DeltaX` and `shares.DeltaY`.
// Similarly, it computes the scalar multiplication of `eshare.GamaX` and `eshare.GamaY`
// with `x` as the scalar, and assigns the result to `shares.GamaX` and `shares.GamaY`.
// Finally, it computes the scalar multiplication of `eshare.ShareX` and `eshare.ShareY`
// with `x` as the scalar, and assigns the result to `shares.ShareX` and `shares.ShareY`.
// The function returns the computed shares as a new Share_G structure.
func (system *ECCShareSystem) share_EXP_P_G_2(eshare Share_G, x *big.Int) Share_G {
	shares := new(Share_G)
	shares.Index = eshare.Index
	DeltaX, DeltaY := system.Curve.ScalarMult(eshare.DeltaX, eshare.DeltaY, x.Bytes())
	shares.DeltaX = DeltaX
	shares.DeltaY = DeltaY
	shares.GamaX, shares.GamaY = system.Curve.ScalarMult(eshare.GamaX, eshare.GamaY, x.Bytes())
	shares.ShareX, shares.ShareY = system.Curve.ScalarMult(eshare.ShareX, eshare.ShareY, x.Bytes())
	return *shares
}

// EXP_P_G_2 calculates the exponentiation of a point G in the elliptic curve
// using the provided shares and the exponent x.
// It returns a pointer to a slice of Share_G containing the resulting shares.
func (system *ECCShareSystem) EXP_P_G_2(eshares *[]Share_G, x *big.Int) *[]Share_G {
	shares := make([]Share_G, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_EXP_P_G_2((*eshares)[i], x)
	}
	return &shares
}

// shareAdd_G performs addition of two Share_G objects and returns the result.
// It takes two Share_G objects as input and returns a new Share_G object.
func (system *ECCShareSystem) shareAdd_G(shares1, shares2 Share_G) Share_G {
	shares := new(Share_G)
	shares.Index = shares1.Index
	DeltaX, DeltaY := system.Curve.Add(shares1.DeltaX, shares1.DeltaY, shares2.DeltaX, shares2.DeltaY)
	shares.DeltaX = DeltaX
	shares.DeltaY = DeltaY
	shares.GamaX, shares.GamaY = system.Curve.Add(shares1.GamaX, shares1.GamaY, shares2.GamaX, shares2.GamaY)
	shares.ShareX, shares.ShareY = system.Curve.Add(shares1.ShareX, shares1.ShareY, shares2.ShareX, shares2.ShareY)
	return *shares
}

// SecAddPlaintext_G performs secure addition of plaintext values to the given shares.
// It takes in an array of shares `shares1`, representing the shares of the first operand,
// and two scalar values `scalarX` and `scalarY`, representing the plaintext values to be added.
// It returns a pointer to an array of shares, where each share represents the result of the addition.
// The function internally generates a second set of shares using the `Share_A_G_Offline` function,
// and then performs the share-wise addition of the two sets of shares using the `shareAdd_G` function.
// The resulting shares are stored in a new array and returned as a pointer.
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
	DeltaX, DeltaY := system.Curve.Add(shares1.DeltaX, shares1.DeltaY, shares2.DeltaX, new(big.Int).Mod(new(big.Int).Neg(shares2.DeltaY), system.Curve.P))
	shares.DeltaX = DeltaX
	shares.DeltaY = DeltaY
	shares.GamaX, shares.GamaY = system.Curve.Add(shares1.GamaX, shares1.GamaY, shares2.GamaX, new(big.Int).Mod(new(big.Int).Neg(shares2.GamaY), system.Curve.P))
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
	xsuba := system.HalfOpenFp(*XsubAshares)
	tshares := system.SecSub_G(hshares, *sharesgB)
	tx, ty := system.HalfOpenG(*tshares)
	t_exp_xsuba_shares := system.EXP_P_G_1(tx, ty, XsubAshares)
	t_exp_a_shares := system.EXP_P_G_1(tx, ty, sharesA)
	gb_exp_xsuba_shares := system.EXP_P_G_2(sharesgB, xsuba)
	shares := system.SecAdd_G(*sharesgC, *gb_exp_xsuba_shares)
	shares = system.SecAdd_G(*shares, *t_exp_a_shares)
	shares = system.SecAdd_G(*shares, *t_exp_xsuba_shares)
	return shares
}

func (system *ECCShareSystem) HalfOpenG(shares []Share_G) (*big.Int, *big.Int) {
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

func (system *ECCShareSystem) MacCheckG(shares []Share_G, res_valueX, res_valueY *big.Int) bool {
	chkX := new(big.Int).Set(system.IdentityGx)
	chkY := new(big.Int).Set(system.IdentityGy)
	deltaX := new(big.Int)
	deltaY := new(big.Int)
	var wg sync.WaitGroup
	tx, ty := system.Curve.Add(res_valueX, res_valueY, shares[0].DeltaX, shares[0].DeltaY)
	for i := 0; i < system.Partynum; i++ {
		deltaX, deltaY = system.Curve.ScalarMult(tx, ty, system.Alphas[i].Bytes())
		deltaX, deltaY = system.Curve.Add(shares[i].GamaX, shares[i].GamaY, deltaX, new(big.Int).Mod(new(big.Int).Neg(deltaY), system.Curve.P))
		commit, r := Com(deltaX.Bytes())
		wg.Add(4)
		go system.Broadcast(&wg, deltaX.Bytes())
		go system.Broadcast(&wg, deltaX.Bytes())
		go system.Broadcast(&wg, commit)
		go system.Broadcast(&wg, r.Bytes())
		opencommit := OpenComit(deltaX.Bytes(), commit, r)
		if !opencommit {
			return false
		}
		chkX, chkY = system.Curve.Add(chkX, chkY, deltaX, deltaY)
	}
	wg.Wait()
	return bytes.Equal(chkX.Bytes(), system.IdentityGx.Bytes())
}

func (system *ECCShareSystem) OpenG(shares []Share_G) (*big.Int, *big.Int, bool) {
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
	chk := system.MacCheckG(shares, ori_valueX, ori_valueY)
	return ori_valueX, ori_valueY, chk
}

func (system *ECCShareSystem) Share_An_Fp(element *big.Int) *[]Share_Fp {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	wg.Add(1)
	go system.Send(&wg, ori_value.Bytes())
	Delta, _ := rand.Int(rand.Reader, system.Order)
	wg.Add(1)
	go system.BroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Add(ori_value, Delta)
	Gama = Gama.Mul(system.alpha, Gama)
	Gama = Gama.Mod(Gama, system.Order)
	shares := make([]Share_Fp, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = rand.Int(rand.Reader, system.Order)
			shares[i].Gama, _ = rand.Int(rand.Reader, system.Order)
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

func (system *ECCShareSystem) Share_An_Fp_Offline(element *big.Int) *[]Share_Fp {
	ori_value := new(big.Int).Set(element)
	var wg sync.WaitGroup
	Delta, _ := rand.Int(rand.Reader, system.Order)
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
			shares[i].Share, _ = rand.Int(rand.Reader, system.Order)
			shares[i].Gama, _ = rand.Int(rand.Reader, system.Order)
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

func (system *ECCShareSystem) RandomShareFp() *[]Share_Fp {
	r, _ := rand.Int(rand.Reader, system.Order)
	rshares := system.Share_An_Fp_Offline(r)
	return rshares
}

func (system *ECCShareSystem) shareAdd(shares1, shares2 Share_Fp) Share_Fp {
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
	Delta := new(big.Int).Sub(shares1.Delta, shares2.Delta)
	Delta = Delta.Mod(Delta, system.Order)
	shares.Delta = Delta
	shares.Gama = new(big.Int).Sub(shares1.Gama, shares2.Gama)
	shares.Gama = shares.Gama.Mod(shares.Gama, system.Order)
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
	Delta := new(big.Int).Mul(shares1.Delta, scalar)
	Delta = Delta.Mod(Delta, system.Order)
	shares.Delta = Delta
	shares.Gama = new(big.Int).Mul(shares1.Gama, scalar)
	shares.Gama = shares.Gama.Mod(shares.Gama, system.Order)
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

func (system *ECCShareSystem) SecSquare(shares1 []Share_Fp) *[]Share_Fp {
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

func (system *ECCShareSystem) HalfOpenFp(shares []Share_Fp) *big.Int {
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

func (system *ECCShareSystem) MacCheckFp(shares []Share_Fp, res_value *big.Int) bool {
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
		go system.Broadcast(&wg, delta.Bytes())
		go system.Broadcast(&wg, commit)
		go system.Broadcast(&wg, r.Bytes())
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

func (system *ECCShareSystem) OpenFp(shares []Share_Fp) (*big.Int, bool) {
	ori_value := big.NewInt(0)
	var wg sync.WaitGroup
	for i := 0; i < system.Partynum; i++ {
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Bytes())
		ori_value = ori_value.Add(ori_value, shares[i].Share)
	}
	wg.Wait()
	ori_value = ori_value.Mod(ori_value, system.Order)
	chk := system.MacCheckFp(shares, ori_value)
	return ori_value, chk
}
