package benckmark

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Oryx/curve"
	mpc "github.com/Oryx/mpc"
)

func TestMaliciousHalfOpenG1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG1(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.HalfOpenG1(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening partially [P] on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousOpenG1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG1(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.OpenG1(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening [P] on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecAddG1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG1(rand.Reader)
		_, g2, _ := curve.RandomG1(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		shares2 := system.Share_A_G1(g2)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecAdd_G1(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]+[Q] on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecAddPG1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG1(rand.Reader)
		_, g2, _ := curve.RandomG1(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecAddPlaintext_G1(*shares1, g2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]+Q on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecSubG1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG1(rand.Reader)
		_, g2, _ := curve.RandomG1(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		shares2 := system.Share_A_G1(g2)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecSub_G1(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]-[Q] on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecSubPG1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomG1(rand.Reader)
		_, g2, _ := curve.RandomG1(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		testnum := 1 << 18
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecSubPlaintext_G1(*shares1, g2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]-Q on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp1G1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomG1(rand.Reader)
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
					system.EXP_P_G1_1(g1, shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [x]*P on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp2G1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomG1(rand.Reader)
		shares1 := system.Share_A_G1(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_P_G1_2(shares1, x1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing x*[P] on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp3G1WAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomG1(rand.Reader)
		shares1 := system.Share_An_Fp(x1)
		shares2 := system.Share_A_G1(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_S_G1(*shares2, *shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [x]*[P] on the BP group G_1 %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func BenckmarkFG1() {
	TestMaliciousHalfOpenG1WAN(100)
	TestMaliciousHalfOpenG1WAN(500)
	TestMaliciousOpenG1WAN(100)
	TestMaliciousOpenG1WAN(500)
	TestMaliciousSecAddG1WAN(100)
	TestMaliciousSecAddG1WAN(500)
	TestMaliciousSecAddPG1WAN(100)
	TestMaliciousSecAddPG1WAN(500)
	TestMaliciousSecSubG1WAN(100)
	TestMaliciousSecSubG1WAN(500)
	TestMaliciousSecSubPG1WAN(100)
	TestMaliciousSecSubPG1WAN(500)
	TestMaliciousSecExp1G1WAN(100)
	TestMaliciousSecExp1G1WAN(500)
	TestMaliciousSecExp2G1WAN(100)
	TestMaliciousSecExp2G1WAN(500)
	TestMaliciousSecExp3G1WAN(100)
	TestMaliciousSecExp3G1WAN(500)
}
