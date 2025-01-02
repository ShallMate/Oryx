package benckmark

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Oryx/curve"
	"github.com/Oryx/shmpc"
)

func TestSHSecPair1() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		_, g1, _ := curve.RandomG1(rand.Reader)
		_, g2, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		testnum := 1 << 16
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.Pair_P_1(shares1, g2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("Computing e([P],Q) on the BP group %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecPair2() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		_, g1, _ := curve.RandomG1(rand.Reader)
		_, g2, _ := curve.RandomG2(rand.Reader)
		shares2 := system.Share_A_G2(g2)
		testnum := 1 << 16
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.Pair_P_2(g1, shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("Computing e(P,[Q]) on the BP group %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecPair3() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		_, g1, _ := curve.RandomG1(rand.Reader)
		_, g2, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		shares2 := system.Share_A_G2(g2)
		testnum := 1 << 14
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.Pair_S(shares1, shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("Computing e([P],[Q]) on the BP group %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}
