package mpc

import (
	"crypto/rand"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	curve "github.com/Oryx/curve"
	"github.com/Oryx/ecc"
)

var zero = big.NewInt(0)
var one = big.NewInt(1)
var two = big.NewInt(2)

type ShareSystem struct {
	alpha           *big.Int
	Partynum        int
	Alphas          []*big.Int
	AlphasMul       []*big.Int
	IdentityG1      *curve.G1
	IdentityG2      *curve.G2
	GenGT           *curve.GT
	IdentityGT      *curve.GT
	IdentityGTBytes []byte
	Order           *big.Int
	OrderMul        *big.Int
	Com             int64
	OfflineCom      int64
	isWAN           bool
	bandwidthMbps   float64
}

type ECCShareSystem struct {
	alpha         *big.Int
	Partynum      int
	Alphas        []*big.Int
	IdentityGx    *big.Int
	IdentityGy    *big.Int
	Order         *big.Int
	Curve         *ecc.KoblitzCurve
	Com           int64
	OfflineCom    int64
	isWAN         bool
	bandwidthMbps float64
}

type RSAShareSystem struct {
	alpha         *big.Int
	Partynum      int
	Alphas        []*big.Int
	AlphasMul     []*big.Int
	Order         *big.Int
	OrderMul      *big.Int
	Com           int64
	OfflineCom    int64
	isWAN         bool
	bandwidthMbps float64
}

type Triplets struct {
	A *big.Int
	B *big.Int
	C *big.Int
}

type SquarePair struct {
	A *big.Int
	B *big.Int
}

func BandwidthLimitMbps(msg []byte, bandwidthMbps float64) {
	delaySeconds := float64(len(msg)*8) / (bandwidthMbps * 1_000_000)
	delay := time.Duration(delaySeconds * float64(time.Second))
	time.Sleep(delay)
}

func (system *ShareSystem) Send(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64(len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ShareSystem) Broadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum-1)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ShareSystem) BroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ShareSystem) OfflineSend(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64(len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ShareSystem) OfflineBroadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum-1)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ShareSystem) OfflineBroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ECCShareSystem) Send(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64(len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ECCShareSystem) Broadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum-1)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ECCShareSystem) BroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ECCShareSystem) OfflineSend(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64(len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ECCShareSystem) OfflineBroadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum-1)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ECCShareSystem) OfflineBroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *RSAShareSystem) Send(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64(len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *RSAShareSystem) Broadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum-1)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *RSAShareSystem) BroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *RSAShareSystem) OfflineSend(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64(len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *RSAShareSystem) OfflineBroadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum-1)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *RSAShareSystem) OfflineBroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum)*len(msg)))
	if system.isWAN {
		BandwidthLimitMbps(msg, system.bandwidthMbps)
	}
	wg.Done()
}

func (system *ShareSystem) GenTriplets() (*[]Share_Fp, *[]Share_Fp, *[]Share_Fp) {
	A, _ := curve.RandomK(rand.Reader)
	B, _ := curve.RandomK(rand.Reader)
	C := new(big.Int).Mul(A, B)
	C = C.Mod(C, system.Order)
	sharesA := system.Share_An_Fp_Offline(A)
	sharesB := system.Share_An_Fp_Offline(B)
	sharesC := system.Share_An_Fp_Offline(C)
	return sharesA, sharesB, sharesC
}

func (system *ShareSystem) GenTriplets_for_Exp() (*[]Share_Fp, *[]Share_Fp, *[]Share_Fp) {
	A, _ := curve.RandomK(rand.Reader)
	B, _ := curve.RandomK(rand.Reader)
	C := new(big.Int).Mul(A, B)
	C = C.Mod(C, system.OrderMul)
	sharesA := system.Share_An_Fp_for_EXP_Offline(A)
	sharesB := system.Share_An_Fp_for_EXP_Offline(B)
	sharesC := system.Share_An_Fp_for_EXP_Offline(C)
	return sharesA, sharesB, sharesC
}

func (system *ECCShareSystem) GenTriplets() (*[]Share_Fp, *[]Share_Fp, *[]Share_Fp) {
	A, _ := rand.Int(rand.Reader, system.Order)
	B, _ := rand.Int(rand.Reader, system.Order)
	C := new(big.Int).Mul(A, B)
	C = C.Mod(C, system.Order)
	sharesA := system.Share_An_Fp_Offline(A)
	sharesB := system.Share_An_Fp_Offline(B)
	sharesC := system.Share_An_Fp_Offline(C)
	return sharesA, sharesB, sharesC
}

func (system *RSAShareSystem) GenTriplets() (*[]Share_Fn, *[]Share_Fn, *[]Share_Fn) {
	A, _ := rand.Int(rand.Reader, system.Order)
	B, _ := rand.Int(rand.Reader, system.Order)
	C := new(big.Int).Mul(A, B)
	C = C.Mod(C, system.Order)
	sharesA := system.Share_An_Fn_Offline(A)
	sharesB := system.Share_An_Fn_Offline(B)
	sharesC := system.Share_An_Fn_Offline(C)
	return sharesA, sharesB, sharesC
}

func (system *ShareSystem) GenSquarePair() (*[]Share_Fp, *[]Share_Fp) {
	A, _ := curve.RandomK(rand.Reader)
	B := new(big.Int).Mul(A, A)
	B = B.Mod(B, system.Order)
	sharesA := system.Share_An_Fp_Offline(A)
	sharesB := system.Share_An_Fp_Offline(B)
	return sharesA, sharesB
}

func (system *ECCShareSystem) GenSquarePair() (*[]Share_Fp, *[]Share_Fp) {
	A, _ := rand.Int(rand.Reader, system.Order)
	B := new(big.Int).Mul(A, A)
	B = B.Mod(B, system.Order)
	sharesA := system.Share_An_Fp_Offline(A)
	sharesB := system.Share_An_Fp_Offline(B)
	return sharesA, sharesB
}

func (system *RSAShareSystem) GenSquarePair() (*[]Share_Fn, *[]Share_Fn) {
	A, _ := rand.Int(rand.Reader, system.Order)
	B := new(big.Int).Mul(A, A)
	B = B.Mod(B, system.Order)
	sharesA := system.Share_An_Fn_Offline(A)
	sharesB := system.Share_An_Fn_Offline(B)
	return sharesA, sharesB
}

func OrderOfElement(element, modulus *big.Int) *big.Int {
	one := big.NewInt(1)
	order := big.NewInt(1)
	current := new(big.Int).Set(element)

	for current.Cmp(one) != 0 {
		current.Mul(current, element)
		current.Mod(current, modulus)
		order.Add(order, one)
	}

	return order
}

func SystemInit(Partynum int) *ShareSystem {
	system := new(ShareSystem)
	system.Partynum = Partynum
	system.alpha, _ = curve.RandomK(rand.Reader)
	orialpha := new(big.Int).Set(system.alpha)
	orialphamul := new(big.Int).Set(system.alpha)
	system.OrderMul = new(big.Int).Sub(curve.Order, one)
	system.Alphas = make([]*big.Int, Partynum)
	system.AlphasMul = make([]*big.Int, Partynum)
	for i := 0; i < system.Partynum; i++ {
		if i < system.Partynum-1 {
			system.Alphas[i], _ = curve.RandomK(rand.Reader)
			orialpha = orialpha.Sub(orialpha, system.Alphas[i])
			orialpha = orialpha.Mod(orialpha, curve.Order)
			system.AlphasMul[i], _ = curve.RandomK(rand.Reader)
			orialphamul = orialphamul.Sub(orialphamul, system.AlphasMul[i])
			orialphamul = orialphamul.Mod(orialphamul, system.OrderMul)
		} else {
			system.Alphas[i] = orialpha
			system.AlphasMul[i] = orialphamul
		}
	}
	system.IdentityG1 = new(curve.G1).ScalarBaseMult(zero)
	system.IdentityG2 = new(curve.G2).ScalarBaseMult(zero)
	system.GenGT = curve.Pair(curve.Gen1, curve.Gen2)
	system.IdentityGT = new(curve.GT).ScalarMult(system.GenGT, zero)
	system.IdentityGTBytes = system.IdentityGT.Marshal()
	system.Order = new(big.Int).Set(curve.Order)
	return system
}

func ECCSystemInit(Partynum int) *ECCShareSystem {
	system := new(ECCShareSystem)
	system.Partynum = Partynum
	s := ecc.S256()
	system.alpha, _ = rand.Int(rand.Reader, s.N)
	orialpha := new(big.Int).Set(system.alpha)
	system.Alphas = make([]*big.Int, Partynum)
	for i := 0; i < system.Partynum; i++ {
		if i < system.Partynum-1 {
			system.Alphas[i], _ = curve.RandomK(rand.Reader)
			orialpha = orialpha.Sub(orialpha, system.Alphas[i])
			orialpha = orialpha.Mod(orialpha, s.N)
		} else {
			system.Alphas[i] = orialpha
		}
	}
	system.IdentityGx, system.IdentityGy = s.ScalarMult(s.Gx, s.Gy, zero.Bytes())
	system.Order = new(big.Int).Set(s.N)
	system.Curve = s
	return system
}

func RSASystemInit(Partynum int, Element *big.Int, Order *big.Int) *RSAShareSystem {
	system := new(RSAShareSystem)
	system.Partynum = Partynum
	system.OrderMul = OrderOfElement(Element, Order)
	system.alpha, _ = rand.Int(rand.Reader, system.OrderMul)
	orialpha := new(big.Int).Set(system.alpha)
	system.Alphas = make([]*big.Int, Partynum)
	for i := 0; i < system.Partynum; i++ {
		if i < system.Partynum-1 {
			system.Alphas[i], _ = rand.Int(rand.Reader, Order)
			orialpha = orialpha.Sub(orialpha, system.Alphas[i])
			orialpha = orialpha.Mod(orialpha, Order)
		} else {
			system.Alphas[i] = orialpha
		}
	}
	system.Order = new(big.Int).Set(Order)
	return system
}

func SystemInitWAN(Partynum int, bandwidth float64) *ShareSystem {
	system := new(ShareSystem)
	system.Partynum = Partynum
	system.alpha, _ = curve.RandomK(rand.Reader)
	orialpha := new(big.Int).Set(system.alpha)
	orialphamul := new(big.Int).Set(system.alpha)
	system.OrderMul = new(big.Int).Sub(curve.Order, one)
	system.Alphas = make([]*big.Int, Partynum)
	system.AlphasMul = make([]*big.Int, Partynum)
	for i := 0; i < system.Partynum; i++ {
		if i < system.Partynum-1 {
			system.Alphas[i], _ = curve.RandomK(rand.Reader)
			orialpha = orialpha.Sub(orialpha, system.Alphas[i])
			orialpha = orialpha.Mod(orialpha, curve.Order)
			system.AlphasMul[i], _ = curve.RandomK(rand.Reader)
			orialphamul = orialphamul.Sub(orialphamul, system.AlphasMul[i])
			orialphamul = orialphamul.Mod(orialphamul, system.OrderMul)
		} else {
			system.Alphas[i] = orialpha
			system.AlphasMul[i] = orialphamul
		}
	}
	system.IdentityG1 = new(curve.G1).ScalarBaseMult(zero)
	system.IdentityG2 = new(curve.G2).ScalarBaseMult(zero)
	system.GenGT = curve.Pair(curve.Gen1, curve.Gen2)
	system.IdentityGT = new(curve.GT).ScalarMult(system.GenGT, zero)
	system.IdentityGTBytes = system.IdentityGT.Marshal()
	system.Order = new(big.Int).Set(curve.Order)
	system.isWAN = true
	system.bandwidthMbps = bandwidth
	return system
}

func ECCSystemInitWAN(Partynum int, bandwidth float64) *ECCShareSystem {
	system := new(ECCShareSystem)
	system.Partynum = Partynum
	s := ecc.S256()
	system.alpha, _ = rand.Int(rand.Reader, s.N)
	orialpha := new(big.Int).Set(system.alpha)
	system.Alphas = make([]*big.Int, Partynum)
	for i := 0; i < system.Partynum; i++ {
		if i < system.Partynum-1 {
			system.Alphas[i], _ = curve.RandomK(rand.Reader)
			orialpha = orialpha.Sub(orialpha, system.Alphas[i])
			orialpha = orialpha.Mod(orialpha, s.N)
		} else {
			system.Alphas[i] = orialpha
		}
	}
	system.IdentityGx, system.IdentityGy = s.ScalarMult(s.Gx, s.Gy, zero.Bytes())
	system.Order = new(big.Int).Set(s.N)
	system.Curve = s
	system.isWAN = true
	system.bandwidthMbps = bandwidth
	return system
}

func RSASystemInitWAN(Partynum int, Element *big.Int, Order *big.Int, bandwidth float64) *RSAShareSystem {
	system := new(RSAShareSystem)
	system.Partynum = Partynum
	system.OrderMul = OrderOfElement(Element, Order)
	system.alpha, _ = rand.Int(rand.Reader, system.OrderMul)
	orialpha := new(big.Int).Set(system.alpha)
	system.Alphas = make([]*big.Int, Partynum)
	for i := 0; i < system.Partynum; i++ {
		if i < system.Partynum-1 {
			system.Alphas[i], _ = rand.Int(rand.Reader, Order)
			orialpha = orialpha.Sub(orialpha, system.Alphas[i])
			orialpha = orialpha.Mod(orialpha, Order)
		} else {
			system.Alphas[i] = orialpha
		}
	}
	system.Order = new(big.Int).Set(Order)
	system.isWAN = true
	system.bandwidthMbps = bandwidth
	return system
}
