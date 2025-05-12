package shmpc

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
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
	alpha         *big.Int
	Partynum      int
	Alphas        []*big.Int
	IdentityG1    *curve.G1
	IdentityG2    *curve.G2
	GenGT         *curve.GT
	IdentityGT    *curve.GT
	Order         *big.Int
	OrderMul      *big.Int
	Com           int64
	OfflineCom    int64
	isWAN         bool
	bandwidth     float64
	BandwidthCtrl *BandwidthSimulator
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
	bandwidth     float64
	BandwidthCtrl *BandwidthSimulator
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

type BandwidthSimulator struct {
	bandwidthMbps float64
	nextFreeAt    int64 // 以纳秒为单位的 Unix 时间戳
}

func NewBandwidthSimulator(bandwidthMbps float64) *BandwidthSimulator {
	return &BandwidthSimulator{
		bandwidthMbps: bandwidthMbps,
		nextFreeAt:    time.Now().UnixNano(),
	}
}

func (b *BandwidthSimulator) SimulateSend(msg []byte) {
	delayNs := int64(float64(len(msg)*8) / (b.bandwidthMbps * 1_000_000) * 1e9)

	for {
		now := time.Now().UnixNano()
		prev := atomic.LoadInt64(&b.nextFreeAt)
		start := max(now, prev)
		newFree := start + delayNs
		if atomic.CompareAndSwapInt64(&b.nextFreeAt, prev, newFree) {
			sleepDuration := time.Until(time.Unix(0, start))
			if sleepDuration > 0 {
				time.Sleep(sleepDuration)
			}
			time.Sleep(time.Duration(delayNs))
			break
		}
	}
}

func (b *BandwidthSimulator) SimulateBroadcast(msg []byte, receiverCount int) {
	totalBytes := len(msg) * receiverCount

	// 同之前的无锁限速逻辑
	delayNs := int64(float64(totalBytes*8) / (b.bandwidthMbps * 1_000_000) * 1e9)

	for {
		now := time.Now().UnixNano()
		prev := atomic.LoadInt64(&b.nextFreeAt)

		start := max(now, prev)
		newFree := start + delayNs

		if atomic.CompareAndSwapInt64(&b.nextFreeAt, prev, newFree) {
			sleepUntil := time.Until(time.Unix(0, start))
			if sleepUntil > 0 {
				time.Sleep(sleepUntil)
			}
			time.Sleep(time.Duration(delayNs))
			break
		}
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (system *ShareSystem) Send(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64(len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateSend(msg)
	}
	wg.Done()
}

func (system *ShareSystem) Broadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum-1)*len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateBroadcast(msg, system.Partynum-1)
	}
	wg.Done()
}

func (system *ShareSystem) BroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum)*len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateBroadcast(msg, system.Partynum)
	}
	wg.Done()
}

func (system *ShareSystem) OfflineSend(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64(len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateSend(msg)
	}
	wg.Done()
}

func (system *ShareSystem) OfflineBroadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum-1)*len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateBroadcast(msg, system.Partynum-1)
	}
	wg.Done()
}

func (system *ShareSystem) OfflineBroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum)*len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateBroadcast(msg, system.Partynum)
	}
	wg.Done()
}

func (system *ECCShareSystem) Send(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64(len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateSend(msg)
	}
	wg.Done()
}

func (system *ECCShareSystem) Broadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum-1)*len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateBroadcast(msg, system.Partynum-1)
	}
	wg.Done()
}

func (system *ECCShareSystem) BroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.Com, int64((system.Partynum)*len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateBroadcast(msg, system.Partynum)
	}
	wg.Done()
}

func (system *ECCShareSystem) OfflineSend(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64(len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateSend(msg)
	}
	wg.Done()
}

func (system *ECCShareSystem) OfflineBroadcast(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum-1)*len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateBroadcast(msg, system.Partynum-1)
	}
	wg.Done()
}

func (system *ECCShareSystem) OfflineBroadcastN(wg *sync.WaitGroup, msg []byte) {
	atomic.AddInt64(&system.OfflineCom, int64((system.Partynum)*len(msg)))
	if system.isWAN && system.BandwidthCtrl != nil {
		system.BandwidthCtrl.SimulateBroadcast(msg, system.Partynum)
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

func SystemInit(Partynum int) *ShareSystem {
	system := new(ShareSystem)
	system.Partynum = Partynum
	system.alpha, _ = curve.RandomK(rand.Reader)
	orialpha := new(big.Int).Set(system.alpha)
	system.Alphas = make([]*big.Int, Partynum)
	for i := 0; i < system.Partynum; i++ {
		if i < system.Partynum-1 {
			system.Alphas[i], _ = curve.RandomK(rand.Reader)
			orialpha = orialpha.Sub(orialpha, system.Alphas[i])
			orialpha = orialpha.Mod(orialpha, curve.Order)
		} else {
			system.Alphas[i] = orialpha
		}
	}
	system.IdentityG1 = new(curve.G1).ScalarBaseMult(zero)
	system.IdentityG2 = new(curve.G2).ScalarBaseMult(zero)
	system.GenGT = curve.Pair(curve.Gen1, curve.Gen2)
	system.IdentityGT = new(curve.GT).ScalarMult(system.GenGT, zero)
	system.Order = new(big.Int).Set(curve.Order)
	system.OrderMul = new(big.Int).Sub(curve.Order, one)
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

func SystemInitWAN(Partynum int, bandwidth float64) *ShareSystem {
	system := new(ShareSystem)
	system.Partynum = Partynum
	system.alpha, _ = curve.RandomK(rand.Reader)
	orialpha := new(big.Int).Set(system.alpha)
	system.Alphas = make([]*big.Int, Partynum)
	for i := 0; i < system.Partynum; i++ {
		if i < system.Partynum-1 {
			system.Alphas[i], _ = curve.RandomK(rand.Reader)
			orialpha = orialpha.Sub(orialpha, system.Alphas[i])
			orialpha = orialpha.Mod(orialpha, curve.Order)
		} else {
			system.Alphas[i] = orialpha
		}
	}
	system.IdentityG1 = new(curve.G1).ScalarBaseMult(zero)
	system.IdentityG2 = new(curve.G2).ScalarBaseMult(zero)
	system.GenGT = curve.Pair(curve.Gen1, curve.Gen2)
	system.IdentityGT = new(curve.GT).ScalarMult(system.GenGT, zero)
	system.Order = new(big.Int).Set(curve.Order)
	system.OrderMul = new(big.Int).Sub(curve.Order, one)
	system.isWAN = true
	if bandwidth > 0 {
		system.BandwidthCtrl = NewBandwidthSimulator(bandwidth)
	} else {
		fmt.Println("The bindwidth < 0")
		os.Exit(1)
	}
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
	if bandwidth > 0 {
		system.BandwidthCtrl = NewBandwidthSimulator(bandwidth)
	} else {
		fmt.Println("The bindwidth < 0")
		os.Exit(1)
	}
	return system
}
