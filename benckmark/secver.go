package benckmark

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	ibs "github.com/Oryx/IBS"
	"github.com/Oryx/bls"
	"github.com/Oryx/ecdsa"
)

func TestSecVerECDSA() {
	// Parameter 1 is the number of parties, and parameter 2 denotes whether it is the malicious model
	system := ecdsa.SecureVerInit(2, true, false, 0)

	// Setup
	eccsystem := ecdsa.NewECDSA()

	// KeyGen
	sk, pk := eccsystem.KeyGen()
	msg := "hello world"

	// Sign
	sig := eccsystem.SignwithInv(sk, []byte(msg))

	// How to share
	sigshares := system.Share_A_Sig(*sig, pk)

	// Verify both the signature and MAC
	// If it is the semi-honest model, chk will always be 1
	ver, chk := system.SecVer(sigshares)
	if chk {
		if ver {
			fmt.Println("Signature verification successful")
		} else {
			fmt.Println("Signature verification failed")
		}
	}
	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.System.OfflineCom+system.System.OfflineCom)

	// If it is the semi-honest model
	// fmt.Println("The communication: ", system.SemiSystem.OfflineCom+system.SemiSystem.OfflineCom)
}

func BenckmarkSecVerECDSA() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := ecdsa.SecureVerInit(partynum, false, false, 0)
		eccsystem := ecdsa.NewECDSA()
		sk, pk := eccsystem.KeyGen()
		msg := "hello world"
		sig := eccsystem.SignwithInv(sk, []byte(msg))
		sigshares := system.Share_A_Sig(*sig, pk)
		testnum := 1 << 16
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SemiSecVerWithoutOpen(sigshares)
				}
			}()
		}
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("SecVer based on ECDSA %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		//fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.System.Com+system.System.OfflineCom)/1024/1024, float64(system.System.Com+system.System.OfflineCom)/1024/1024/float64(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.SemiSystem.Com+system.SemiSystem.OfflineCom)/1024/1024, float64(system.SemiSystem.Com+system.SemiSystem.OfflineCom)/1024/1024/float64(testnum))
	}

}

func TestSecVerBLS() {
	// Parameter 1 is the number of parties, and parameter 2 denotes whether it is the malicious model
	system := bls.SecureVerInit(2, true, false, 0)

	// KeyGen
	sk, pk := bls.KeyGen()
	msg := "hello world"

	// Sign
	sig, hm := bls.SignwithHm(sk, []byte(msg))

	// How to share
	sigshares := system.Share_A_Sig(sig, hm, pk)

	// Verify both the signature and MAC
	// If it is the semi-honest model, chk will always be 1
	ver, chk := system.SecVer(sigshares)
	if chk {
		if ver {
			fmt.Println("Signature verification successful")
		} else {
			fmt.Println("Signature verification failed")
		}
	}
	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.System.OfflineCom+system.System.OfflineCom)

	// If it is the semi-honest model
	// fmt.Println("The communication: ", system.SemiSystem.OfflineCom+system.SemiSystem.OfflineCom)
}

func BenckmarkSecVerBLS() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := bls.SecureVerInit(partynum, false, false, 0)
		sk, pk := bls.KeyGen()
		msg := "hello world"
		sig, hm := bls.SignwithHm(sk, []byte(msg))
		sigshares := system.Share_A_Sig(sig, hm, pk)
		testnum := 1 << 13
		var wg sync.WaitGroup
		t1 := time.Now()
		for i := 0; i < testnum; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				system.SemiSecVerWithoutOpen(sigshares)
			}()
		}
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("SecVer based on BLS %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		//fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.System.Com+system.System.OfflineCom)/1024/1024, float64(system.System.Com+system.System.OfflineCom)/1024/1024/float64(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.SemiSystem.Com+system.SemiSystem.OfflineCom)/1024/1024, float64(system.SemiSystem.Com+system.SemiSystem.OfflineCom)/1024/1024/float64(testnum))
	}

}

func TestSecVerAIBS() {
	msk := ibs.MasterKeyGen()

	// Parameter 1 is the number of parties
	// parameter 2 denotes the MPK
	// Parameter 3 denotes whether it is the malicious model
	system := *ibs.SecureVerInit(2, &msk.MasterPubKey, true, false, 0)

	// id
	userid := big.NewInt(9567)
	// KeyGen
	sk := ibs.UserKeyGen(msk, userid)

	msg := "hello world"

	// Sign
	sig := ibs.Sign(sk, &msk.MasterPubKey, []byte(msg))

	// How to share
	sigshares := system.Share_A_Sig(*sig, []byte(msg), userid)

	// Verify both the signature and MAC
	// If it is the semi-honest model, chk will always be 1
	ver, chk := system.SecVer(sigshares)
	if chk {
		if ver {
			fmt.Println("Signature verification successful")
		} else {
			fmt.Println("Signature verification failed")
		}
	}
	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.System.OfflineCom+system.System.OfflineCom)

	// If it is the semi-honest model
	// fmt.Println("The communication: ", system.SemiSystem.OfflineCom+system.SemiSystem.OfflineCom)
}

func BenckmarkSecVerAIBS() {
	for partynum := 2; partynum <= 10; partynum++ {
		msk := ibs.MasterKeyGen()
		system := *ibs.SecureVerInit(partynum, &msk.MasterPubKey, false, false, 0)
		userid := big.NewInt(9567)
		sk := ibs.UserKeyGen(msk, userid)
		msg := "hello world"
		sig := ibs.Sign(sk, &msk.MasterPubKey, []byte(msg))
		sigshares := system.Share_A_Sig(*sig, []byte(msg), userid)
		testnum := 1 << 13
		var wg sync.WaitGroup
		t1 := time.Now()
		for i := 0; i < testnum; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				system.SemiSecVerWithoutOpen(sigshares)
			}()
		}
		wg.Wait()
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("SecVer based on AIBS %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		//fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.System.Com+system.System.OfflineCom)/1024/1024, float64(system.System.Com+system.System.OfflineCom)/1024/1024/float64(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.SemiSystem.Com+system.SemiSystem.OfflineCom)/1024/1024, float64(system.SemiSystem.Com+system.SemiSystem.OfflineCom)/1024/1024/float64(testnum))
	}

}
