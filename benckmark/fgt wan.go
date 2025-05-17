package benckmark

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Oryx/curve"
	mpc "github.com/Oryx/mpc"
)

func TestMaliciousHalfOpenGTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.HalfOpenGT(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening partially [P] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousOpenGTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.OpenGT(*shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Opening [P] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecAddGTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		_, g2, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		shares2 := system.Share_A_GT(g2)
		testnum := 1 << 16
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecAdd_GT(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]+[Q] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecAddPGTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		_, g2, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecAddPlaintext_GT(*shares1, g2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]+Q on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecSubGTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		_, g2, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		shares2 := system.Share_A_GT(g2)
		testnum := 1 << 16
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecSub_GT(*shares1, *shares2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]-[Q] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecSubPGTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		_, g2, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.SecSubPlaintext_GT(*shares1, g2)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [P]-Q on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp1GTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomGTK(rand.Reader)
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
					system.EXP_P_GT_1(g1, shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [x]*P on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp2GTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 15
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_P_GT_2(shares1, x1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing x*[P] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestMaliciousSecExp3GTWAN(bandwidth float64) {
	for partynum := 2; partynum <= 10; partynum++ {
		system := mpc.SystemInitWAN(partynum, bandwidth)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_An_Fp(x1)
		shares2 := system.Share_A_GT(g1)
		testnum := 1 << 13
		worker := 8192
		var wg sync.WaitGroup // 创建 WaitGroup 实例
		t1 := time.Now()
		for i := 0; i < worker; i++ {
			wg.Add(1) // 为每个协程增加计数
			go func() {
				defer wg.Done() // 在协程完成时调用 Done() 递减计数器
				for j := 0; j < testnum/worker; j++ {
					system.EXP_S_GT(*shares2, *shares1)
				}
			}()
		}
		//rese, chk := system.Open(*shares)
		wg.Wait() // 等待所有协程完成
		t2 := time.Since(t1)
		fmt.Printf("n=%d\n", partynum)
		fmt.Printf("In the WAN setting. The bandwidth is %.2f Mbps\n", bandwidth)
		fmt.Printf("Computing [x]*[P] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func BenckmarkFGT() {
	TestMaliciousHalfOpenGTWAN(100)
	TestMaliciousHalfOpenGTWAN(500)
	TestMaliciousOpenGTWAN(100)
	TestMaliciousOpenGTWAN(500)
	TestMaliciousSecAddGTWAN(100)
	TestMaliciousSecAddGTWAN(500)
	TestMaliciousSecAddPGTWAN(100)
	TestMaliciousSecAddPGTWAN(500)
	TestMaliciousSecSubGTWAN(100)
	TestMaliciousSecSubGTWAN(500)
	TestMaliciousSecSubPGTWAN(100)
	TestMaliciousSecSubPGTWAN(500)
	TestMaliciousSecExp1GTWAN(100)
	TestMaliciousSecExp1GTWAN(500)
	TestMaliciousSecExp2GTWAN(100)
	TestMaliciousSecExp2GTWAN(500)
	TestMaliciousSecExp3GTWAN(100)
	TestMaliciousSecExp3GTWAN(500)
}
