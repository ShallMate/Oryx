package benckmark

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Oryx/curve"
	mpc "github.com/Oryx/mpc"
)

func TestMaliciousHalfOpenG2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G2(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.HalfOpenG2(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening partially [P] on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousOpenG2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G2(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.OpenG2(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening [P] on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecAddG2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG2(rand.Reader)
		_, g2, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G2(g1)
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
					system.SecAdd_G2(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]+[Q] on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecAddPG2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG2(rand.Reader)
		_, g2, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G2(g1)
		testnum := 1 << 16
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecAddPlaintext_G2(*shares1, g2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]+Q on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecSubG2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG2(rand.Reader)
		_, g2, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G2(g1)
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
					system.SecSub_G2(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]-[Q] on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecSubPG2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG2(rand.Reader)
		_, g2, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G2(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecSubPlaintext_G2(*shares1, g2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]-Q on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp1G2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_An_Fp(x1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_P_G2_1(g1, shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [x]*P on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp2G2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_A_G2(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_P_G2_2(shares1, x1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("Computing x*[P] on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp3G2WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomG2(rand.Reader)
		shares1 := system.Share_An_Fp(x1)
		shares2 := system.Share_A_G2(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_S_G2(*shares2, *shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [x]*[P] on the BP group G_2 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func BenckmarkFG2() {
	TestMaliciousHalfOpenG2WAN(100)
	TestMaliciousHalfOpenG2WAN(500)
	TestMaliciousOpenG2WAN(100)
	TestMaliciousOpenG2WAN(500)
	TestMaliciousSecAddG2WAN(100)
	TestMaliciousSecAddG2WAN(500)
	TestMaliciousSecAddPG2WAN(100)
	TestMaliciousSecAddPG2WAN(500)
	TestMaliciousSecSubG2WAN(100)
	TestMaliciousSecSubG2WAN(500)
	TestMaliciousSecSubPG2WAN(100)
	TestMaliciousSecSubPG2WAN(500)
	TestMaliciousSecExp1G2WAN(100)
	TestMaliciousSecExp1G2WAN(500)
	TestMaliciousSecExp2G2WAN(100)
	TestMaliciousSecExp2G2WAN(500)
	TestMaliciousSecExp3G2WAN(100)
	TestMaliciousSecExp3G2WAN(500)
}
