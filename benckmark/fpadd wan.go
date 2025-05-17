package benckmark

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Oryx/curve"
	mpc "github.com/Oryx/mpc"
)

func TestMaliciousHalfOpenWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp(e1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.HalfOpenFp(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening partially [x] on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousOpenWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp(e1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.OpenFp(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening [x] on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecAddWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		e2, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp(e1)
		shares2 := system.Share_An_Fp(e2)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecAdd(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating [x] + [y] on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecSubWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		e2, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp(e1)
		shares2 := system.Share_An_Fp(e2)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecSub(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating [x] - [y] on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecMulPWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		e2, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp(e1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecMulPlaintext(*shares1, e2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating [x] * y on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecMulWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		e2, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp(e1)
		shares2 := system.Share_An_Fp(e2)
		testnum := 1 << 14
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecMul(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating [x] * [y] on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecSquareWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		e1, _ := curve.RandomK(rand.Reader)
		system := mpc.SystemInitWAN(partynum, bandwidth)
		shares1 := system.Share_An_Fp(e1)
		testnum := 1 << 16
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecSquare(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Calculating [x]^2 on F_p %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func BenckmarkFpadd() {
	TestMaliciousHalfOpenWAN(100)
	TestMaliciousHalfOpenWAN(500)
	TestMaliciousOpenWAN(100)
	TestMaliciousOpenWAN(500)
	TestMaliciousSecAddWAN(100)
	TestMaliciousSecAddWAN(500)
	TestMaliciousSecSubWAN(100)
	TestMaliciousSecSubWAN(500)
	TestMaliciousSecMulPWAN(100)
	TestMaliciousSecMulPWAN(500)
	TestMaliciousSecMulWAN(100)
	TestMaliciousSecMulWAN(500)
	TestMaliciousSecSquareWAN(100)
	TestMaliciousSecSquareWAN(500)

}
