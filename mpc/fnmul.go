package mpc

import (
	"crypto/rand"
	"math/big"
	"sync"
)

// isCoprime checks if two big integers are coprime.
// Two integers are coprime if their greatest common divisor (GCD) is 1.
func isCoprime(a, b *big.Int) bool {
	return new(big.Int).GCD(nil, nil, a, b).Cmp(big.NewInt(1)) == 0
}

// GenerateCoprimeNumber generates a random number that is coprime with the given order.
// It takes a pointer to a big.Int representing the order and returns a pointer to a big.Int representing the generated coprime number.
// If an error occurs during the generation process, it returns nil and the error.
func GenerateCoprimeNumber(order *big.Int) (*big.Int, error) {
	for {
		n, err := rand.Int(rand.Reader, order)
		if err != nil {
			return nil, err
		}
		if isCoprime(n, order) {
			return n, nil
		}
	}
}

// Share_An_Fn_Mul generates multiplication secret shares.
// It takes an element of type *big.Int as input and returns a pointer to an array of Share_Fn.
// Each Share_Fn struct in the array represents a secret share of the multiplication operation.
// The function uses a RSAShareSystem instance to perform the secret sharing.
// It generates a random value Delta and broadcasts it to all parties in the system.
// It then computes Gama as the product of the input element and Delta, raised to the power of system.alpha modulo system.Order.
// For each party in the system, it generates a random coprime number as the share and another random coprime number as Gama.
// It updates the input element and Gama based on the generated shares and computes the final share and Gama for the last party.
// Finally, it sends the shares and Gama values to all parties and returns a pointer to the array of shares.
func (system *RSAShareSystem) Share_An_Fn_Mul(element *big.Int) *[]Share_Fn {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	Delta, _ := rand.Int(rand.Reader, system.Order)
	wg.Add(2)
	go system.Send(&wg, ori_value.Bytes())
	go system.BroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Mul(ori_value, Delta)
	Gama = Gama.Exp(Gama, system.alpha, system.Order)
	shares := make([]Share_Fn, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = GenerateCoprimeNumber(system.Order)
			shares[i].Gama, _ = GenerateCoprimeNumber(system.Order)
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

// Share_An_Fn_Mul_Offline performs offline sharing of an Fn multiplication operation.
// It takes an element of type *big.Int as input and returns a pointer to an array of Share_Fn.
// Each Share_Fn contains the index, Delta, Share, and Gama values.
// The function uses a sync.WaitGroup to synchronize the goroutines.
// It generates a random Delta value and broadcasts it to all parties using the OfflineBroadcastN method.
// It then calculates Gama by multiplying the original value with Delta and exponentiating it with the system's alpha value.
// For each party, it generates a Share and Gama value, except for the last party where it assigns the original value and Gama.
// It sends the Share and Gama values to each party using the OfflineSend method.
// Finally, it waits for all goroutines to finish and returns a pointer to the array of Share_Fn.
func (system *RSAShareSystem) Share_An_Fn_Mul_Offline(element *big.Int) *[]Share_Fn {
	var wg sync.WaitGroup
	ori_value := new(big.Int).Set(element)
	Delta, _ := rand.Int(rand.Reader, system.Order)
	wg.Add(1)
	go system.OfflineBroadcastN(&wg, Delta.Bytes())
	Gama := new(big.Int).Mul(ori_value, Delta)
	Gama = Gama.Exp(Gama, system.alpha, system.Order)
	shares := make([]Share_Fn, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i].Index = i
		shares[i].Delta = Delta
		if i < system.Partynum-1 {
			shares[i].Share, _ = GenerateCoprimeNumber(system.Order)
			shares[i].Gama, _ = GenerateCoprimeNumber(system.Order)
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

// HalfOpenFn_Mul performs a half-open functional multiplication on the given shares.
// It broadcasts each share to all parties in the system and computes the product of all shares.
func (system *RSAShareSystem) HalfOpenFn_Mul(shares []Share_Fn) *big.Int {
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

func (system *RSAShareSystem) MacCheckFn_Mul(shares []Share_Fn, res_value *big.Int) bool {
	var wg sync.WaitGroup
	chkdelta := big.NewInt(1)
	chkgama := big.NewInt(1)
	delta := new(big.Int)
	t := new(big.Int).Mul(res_value, shares[0].Delta)
	t = t.Mod(t, system.Order)
	for i := 0; i < system.Partynum; i++ {
		delta = delta.Exp(t, system.Alphas[i], system.Order)
		commit1, r1 := Com(delta.Bytes())
		commit2, r2 := Com(shares[i].Gama.Bytes())
		wg.Add(6)
		go system.Broadcast(&wg, commit1)
		go system.Broadcast(&wg, commit2)
		go system.Broadcast(&wg, r1.Bytes())
		go system.Broadcast(&wg, r2.Bytes())
		go system.Broadcast(&wg, delta.Bytes())
		go system.Broadcast(&wg, shares[i].Gama.Bytes())
		opencommit1 := OpenComit(delta.Bytes(), commit1, r1)
		opencommit2 := OpenComit(shares[i].Gama.Bytes(), commit2, r2)
		if !opencommit1 || !opencommit2 {
			return false
		}
		chkdelta = chkdelta.Mul(chkdelta, delta)
		chkgama = chkgama.Mul(chkgama, shares[i].Gama)
	}
	chkdelta = chkdelta.Mod(chkdelta, system.Order)
	chkgama = chkgama.Mod(chkgama, system.Order)
	wg.Wait()
	return chkdelta.Cmp(chkgama) == 0
}

func (system *RSAShareSystem) OpenFn_Mul(shares []Share_Fn) (*big.Int, bool) {
	var wg sync.WaitGroup
	ori_value := big.NewInt(1)
	for i := 0; i < system.Partynum; i++ {
		wg.Add(1)
		go system.Broadcast(&wg, shares[i].Share.Bytes())
		ori_value = ori_value.Mul(ori_value, shares[i].Share)
	}
	ori_value = ori_value.Mod(ori_value, system.Order)
	wg.Wait()
	chk := system.MacCheckFn_Mul(shares, ori_value)
	return ori_value, chk
}
