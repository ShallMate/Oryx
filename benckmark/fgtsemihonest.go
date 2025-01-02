package benckmark

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/Oryx/curve"
	"github.com/Oryx/shmpc"
)

func TestSHOpenGT() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 16
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
		fmt.Printf("Opening [P] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecAddGT() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		_, g2, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		shares2 := system.Share_A_GT(g2)
		testnum := 1 << 18
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
		fmt.Printf("Computing [P]+[Q] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecAddPGT() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		_, g2, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 16
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
		fmt.Printf("Computing [P]+Q on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecSubGT() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		_, g2, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		shares2 := system.Share_A_GT(g2)
		testnum := 1 << 18
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
		fmt.Printf("Computing [P]-[Q] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecSubPGT() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		_, g2, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 16
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
		fmt.Printf("Computing [P]-Q on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecExp1GT() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_An_Fp(x1)
		testnum := 1 << 16
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
		fmt.Printf("Computing [x]*P on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecExp2GT() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_A_GT(g1)
		testnum := 1 << 16
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
		fmt.Printf("Computing x*[P] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}

func TestSHSecExp3GT() {
	for partynum := 2; partynum <= 10; partynum++ {
		system := shmpc.SystemInit(partynum)
		x1, _ := rand.Int(rand.Reader, system.Order)
		_, g1, _ := curve.RandomGTK(rand.Reader)
		shares1 := system.Share_An_Fp(x1)
		shares2 := system.Share_A_GT(g1)
		testnum := 1 << 14
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
		fmt.Printf("Computing [x]*[P] on the BP group G_T %d times took %s, averaging %s\n", testnum, t2, t2/time.Duration(testnum))
		fmt.Printf("The communication is %.12f MB, averaging %.12f MB\n", float64(system.Com+system.OfflineCom)/1024/1024, float64(system.Com+system.OfflineCom)/1024/1024/float64(testnum))
	}
}
