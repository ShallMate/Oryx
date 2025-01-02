package pii_bls

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
)

func (system *PIISystem) verphase(inputsets []InputSet) []VerSet {
	versets := make([]VerSet, system.partynum)
	var wg sync.WaitGroup // Add missing import for "sync" package
	blocknum := 4096
	for i := 0; i < system.partynum; i++ {
		versets[i].Vers = make([](*[]mpc.Share_GT), inputsets[i].inputsize)
		versets[i].PKshares = make([](*[]mpc.Share_G2), inputsets[i].inputsize)
		//versets[i].PKxshares = make([](*[]mpc.Share_Fp), inputsets[i].inputsize)
		versets[i].inputsize = inputsets[i].inputsize
		eachblock := inputsets[i].inputsize/blocknum + 1
		for j := 0; ; j = j + eachblock {
			if j+eachblock >= inputsets[i].inputsize {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					for q := j; q < inputsets[i].inputsize; q++ {
						versets[i].Vers[q] = system.PiiSystem.SecVerWithoutOpen(inputsets[i].Sigs[q])
						//fmt.Println(system.PiiSystem.System.HalfOpenGT(*versets[i].Vers[q]))
						versets[i].PKshares[q] = inputsets[i].Sigs[q].Pkshare
						//versets[i].PKxshares[q] = inputsets[i].PKxshares[q]
					}
				}(i, j)
				break
			} else {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					for q := j; q < j+eachblock; q++ {
						versets[i].Vers[q] = system.PiiSystem.SecVerWithoutOpen(inputsets[i].Sigs[q])
						//fmt.Println(system.PiiSystem.System.HalfOpenGT(*versets[i].Vers[q]))
						versets[i].PKshares[q] = inputsets[i].Sigs[q].Pkshare
						//versets[i].PKxshares[q] = inputsets[i].PKxshares[q]
					}
				}(i, j)
			}

		}
	}
	wg.Wait()
	return versets
}

func (system *PIISystem) interphase(versets []VerSet, seedsets []SeedSet) []*curve.G2 {
	intersection := make([]*curve.G2, 0)
	interChan := make(chan *curve.G2, 100)
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
					v := system.PiiSystem.System.SecSub_G2(*versets[0].PKshares[i], *versets[1].PKshares[j])
					w := system.PiiSystem.System.Pair_P_2(system.PiiSystem.System.IdentityG1, v)
					//w := system.PiiSystem.System.EXP_P_GT_1(system.PiiSystem.System.GenGT, v)
					w = system.PiiSystem.System.SecAdd_GT(*w, *versets[0].Vers[i])
					w = system.PiiSystem.System.SecAdd_GT(*w, *versets[1].Vers[j])
					w = system.PiiSystem.System.EXP_S_GT(*w, *seedsets[i].Seeds[j])
					wvalue, chk := system.PiiSystem.System.OpenGT(*w)
					if chk {
						if bytes.Equal(wvalue.Marshal(), system.PiiSystem.System.IdentityGT.Marshal()) {
							interid, chkid := system.PiiSystem.System.OpenG2(*versets[0].PKshares[i])
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
					v := system.PiiSystem.System.SecSub_G2(*versets[0].PKshares[i], *versets[1].PKshares[j])
					w := system.PiiSystem.System.Pair_P_2(system.PiiSystem.System.IdentityG1, v)
					w = system.PiiSystem.System.SecAdd_GT(*w, *versets[0].Vers[i])
					w = system.PiiSystem.System.SecAdd_GT(*w, *versets[1].Vers[j])
					w = system.PiiSystem.System.EXP_S_GT(*w, *seedsets[i].Seeds[j])
					wvalue, chk := system.PiiSystem.System.OpenGT(*w)
					if chk {
						if bytes.Equal(wvalue.Marshal(), system.PiiSystem.System.IdentityGT.Marshal()) {
							interid, chkid := system.PiiSystem.System.OpenG2(*versets[0].PKshares[i])
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

func (system *PIISystem) twoPartyPiiRun(inputsets []InputSet, seedsets []SeedSet) []*curve.G2 {
	fmt.Println("ver phase start")
	vertime := time.Now()
	versets := system.verphase(inputsets)
	fmt.Println("ver time:", time.Since(vertime))
	fmt.Println("inter phase start")
	intertime := time.Now()
	intersection := system.interphase(versets, seedsets)
	fmt.Println("inter time:", time.Since(intertime))
	fmt.Println("intersection:", intersection)
	fmt.Printf("Offline Communication: %f MB\n", float64(system.PiiSystem.System.OfflineCom)/1024/1024)
	fmt.Printf("Online Communication: %f MB\n", float64(system.PiiSystem.System.Com)/1024/1024)
	return intersection
}
