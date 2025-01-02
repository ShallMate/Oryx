package benckmark

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	ibs "github.com/Oryx/IBS"
)

func BenckmarkAIBS() {
	var totalSetup, totalKeyGen, totalSign, totalVer time.Duration
	iterations := 1000

	for i := 0; i < iterations; i++ {
		// Setup
		startSetup := time.Now()
		msk := ibs.MasterKeyGen()
		elapsedSetup := time.Since(startSetup)
		totalSetup += elapsedSetup

		// KeyGen
		startKeyGen := time.Now()
		userid := big.NewInt(9567)
		sk := ibs.UserKeyGen(msk, userid)
		elapsedKeyGen := time.Since(startKeyGen)
		totalKeyGen += elapsedKeyGen

		// Sign
		startSign := time.Now()
		msg := "hello world"
		sig := ibs.Sign(sk, &msk.MasterPubKey, []byte(msg))
		elapsedSign := time.Since(startSign)
		totalSign += elapsedSign

		// Ver
		startVer := time.Now()
		ver := ibs.Ver(sig, []byte(msg), userid, &msk.MasterPubKey)
		elapsedVer := time.Since(startVer)
		totalVer += elapsedVer

		if !ver {
			fmt.Println("Signature verification failed")
		}
	}

	averageSetup := totalSetup / time.Duration(iterations)
	averageKeyGen := totalKeyGen / time.Duration(iterations)
	averageSign := totalSign / time.Duration(iterations)
	averageVer := totalVer / time.Duration(iterations)

	fmt.Printf("Average Setup time: %v\n", averageSetup)
	fmt.Printf("Average KeyGen time: %v\n", averageKeyGen)
	fmt.Printf("Average Sign time: %v\n", averageSign)
	fmt.Printf("Average Ver time: %v\n", averageVer)
}

func BenckmarkMutiThreadAIBS() {
	var totalSetup, totalKeyGen, totalSign, totalVer time.Duration
	iterations := 1 << 17
	msk := make([]*ibs.MasterKey, iterations)
	sk := make([]*ibs.UserKey, iterations)
	sig := make([]*ibs.Sig, iterations)
	var wg sync.WaitGroup
	// Setup
	startSetup := time.Now()
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			msk[i] = ibs.MasterKeyGen()
		}(i)
	}
	wg.Wait()
	elapsedSetup := time.Since(startSetup)
	totalSetup += elapsedSetup

	// KeyGen
	startKeyGen := time.Now()
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			userid := big.NewInt(9567)
			sk[i] = ibs.UserKeyGen(msk[i], userid)
		}(i)
	}
	wg.Wait()
	elapsedKeyGen := time.Since(startKeyGen)
	totalKeyGen += elapsedKeyGen

	// Sign
	startSign := time.Now()
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			msg := "hello world"
			sig[i] = ibs.Sign(sk[i], &msk[i].MasterPubKey, []byte(msg))
		}(i)
	}
	wg.Wait()
	elapsedSign := time.Since(startSign)
	totalSign += elapsedSign

	// Ver
	startVer := time.Now()
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			msg := "hello world"
			userid := big.NewInt(9567)
			ibs.Ver(sig[i], []byte(msg), userid, &msk[i].MasterPubKey)
		}(i)
	}
	wg.Wait()
	elapsedVer := time.Since(startVer)
	totalVer += elapsedVer

	averageSetup := totalSetup / time.Duration(iterations)
	averageKeyGen := totalKeyGen / time.Duration(iterations)
	averageSign := totalSign / time.Duration(iterations)
	averageVer := totalVer / time.Duration(iterations)

	fmt.Printf("Average Setup time: %v\n", averageSetup)
	fmt.Printf("Average KeyGen time: %v\n", averageKeyGen)
	fmt.Printf("Average Sign time: %v\n", averageSign)
	fmt.Printf("Average Ver time: %v\n", averageVer)
}
