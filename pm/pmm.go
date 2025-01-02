package pm

import (
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"
)

func (system *PMSystem) interphase_m(versets []InputSet, seedsets *SeedSet) []*big.Int {
	intersection := make([]*big.Int, 0)
	interChan := make(chan *big.Int, 100)
	var wg sync.WaitGroup
	go func() {
		for inter := range interChan {
			intersection = append(intersection, inter)
		}
	}()
	size_except_for_one := 0
	for i := 1; i < system.System.Partynum; i++ {
		size_except_for_one = versets[i].inputsize + size_except_for_one
	}
	for i := 0; i < versets[0].inputsize; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := system.System.SecSub(*versets[0].HIDs[i], *versets[1].HIDs[0])
			t := 0
			for j := 1; j < size_except_for_one; j++ {
				t = j
				k := 1
				for {
					if t < versets[k].inputsize {
						break
					}
					t = t - versets[k].inputsize
					k = k + 1
				}
				if t == 0 {
					v = system.System.SecAdd(*v, *system.System.SecSub(*versets[0].HIDs[i], *versets[k].HIDs[t]))
				} else {
					v = system.System.SecMul(*v, *system.System.SecSub(*versets[0].HIDs[i], *versets[k].HIDs[t]))
				}
			}
			v = system.System.SecMul(*v, *seedsets.Seeds[i])
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
		}(i)
	}
	wg.Wait()
	close(interChan) // Close the channel after all goroutines are done
	return intersection
}

func (system *PMSystem) PartyRun(inputsets []InputSet, seedsets SeedSet) []*big.Int {
	fmt.Println("inter phase start")
	intertime := time.Now()
	intersection := system.interphase_m(inputsets, &seedsets)
	fmt.Println("inter time:", time.Since(intertime))
	fmt.Println("intersection:", intersection)
	fmt.Printf("Offline Communication: %f MB\n", float64(system.System.OfflineCom)/1024/1024)
	fmt.Printf("Online Communication: %f MB\n", float64(system.System.Com)/1024/1024)
	return intersection
}
