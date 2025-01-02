package pm

import (
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"
)

func (system *PMSystem) interphase(versets []InputSet, seedsets []SeedSet) []*big.Int {
	intersection := make([]*big.Int, 0)
	interChan := make(chan *big.Int, 100)
	var wg sync.WaitGroup
	go func() {
		for inter := range interChan {
			intersection = append(intersection, inter)
		}
	}()
	if versets[0].inputsize*versets[1].inputsize <= 16384 {
		for i := 0; i < versets[0].inputsize; i++ {
			for j := 0; j < versets[1].inputsize; j++ {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					v := system.System.SecSub(*versets[0].HIDs[i], *versets[1].HIDs[j])
					v = system.System.SecMul(*v, *seedsets[i].Seeds[j])
					vvalue, chk := system.System.OpenFp(*v)
					if chk {
						if vvalue.Cmp(system.zero) == 0 {
							interid, chkid := system.System.OpenFp(*versets[0].HIDs[i])
							if chkid {
								interChan <- interid
							} else {
								fmt.Println("chk error")
								os.Exit(1)
							}
						}
					} else {
						fmt.Println("chk error")
						os.Exit(1)
					}
				}(i, j)

			}
		}
	} else {
		for i := 0; i < versets[0].inputsize; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < versets[1].inputsize; j++ {
					v := system.System.SecSub(*versets[0].HIDs[i], *versets[1].HIDs[j])
					v = system.System.SecMul(*v, *seedsets[i].Seeds[j])
					vvalue, chk := system.System.OpenFp(*v)
					if chk {
						if vvalue.Cmp(system.zero) == 0 {
							interid, chkid := system.System.OpenFp(*versets[0].HIDs[i])
							if chkid {
								interChan <- interid
							} else {
								fmt.Println("chk error")
								os.Exit(1)
							}
						}
					} else {
						fmt.Println("chk error")
						os.Exit(1)
					}
				}
			}(i)
		}
	}
	wg.Wait()
	close(interChan) // Close the channel after all goroutines are done
	return intersection
}

func (system *PMSystem) twoPartyRun(inputsets []InputSet, seedsets []SeedSet) []*big.Int {
	fmt.Println("inter phase start")
	intertime := time.Now()
	intersection := system.interphase(inputsets, seedsets)
	fmt.Println("inter time:", time.Since(intertime))
	fmt.Println("intersection:", intersection)
	fmt.Printf("Offline Communication: %f MB\n", float64(system.System.OfflineCom)/1024/1024)
	fmt.Printf("Online Communication: %f MB\n", float64(system.System.Com)/1024/1024)
	return intersection
}
