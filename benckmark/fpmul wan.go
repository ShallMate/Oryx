package benckmark

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Oryx/curve"
	mpc "github.com/Oryx/mpc"
)

func TestMaliciousHalfOpenMulWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp_Mul(e1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.HalfOpenFp_Mul(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening partially <x> on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousOpenMulWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp_Mul(e1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.OpenFp_Mul(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening <x> on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecMul_MulWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		e2, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp_Mul(e1)
		shares2 := system.Share_An_Fp_Mul(e2)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecMul_Mul(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating <x> * <y> on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecDiv_MulWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		e2, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp_Mul(e1)
		shares2 := system.Share_An_Fp_Mul(e2)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecDiv(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating <x> / <y> on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		e2, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares2 := system.Share_An_Fp_for_EXP(e2)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_P_Fp_1(e1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating y^[x] = <y^x> on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		e2, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp_Mul(e1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_P_Fp_2(*shares1, e2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating <x>^y on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func BenckmarkFpmul() {
	TestMaliciousHalfOpenMulWAN(100)
	TestMaliciousHalfOpenMulWAN(500)
	TestMaliciousOpenMulWAN(100)
	TestMaliciousOpenMulWAN(500)
	TestMaliciousSecMul_MulWAN(100)
	TestMaliciousSecMul_MulWAN(500)
	TestMaliciousSecDiv_MulWAN(100)
	TestMaliciousSecDiv_MulWAN(500)
	TestMaliciousSecExp1WAN(100)
	TestMaliciousSecExp1WAN(500)
	TestMaliciousSecExp2WAN(100)
	TestMaliciousSecExp2WAN(500)

}
