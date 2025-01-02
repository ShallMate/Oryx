package pii

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
)

func (system *PIISystem) interphase_m(versets []VerSet, seedsets *SeedSet) []*big.Int {
	intersection := make([]*big.Int, 0)
	var wg sync.WaitGroup
	size_except_for_one := 0
	for i := 1; i < system.partynum; i++ {
		size_except_for_one = versets[i].inputsize + size_except_for_one
	}
	verres := make([][]bool, system.partynum)
	var chkpool = sync.Pool{
		New: func() interface{} {
			return new(bool)
		},
	}
	var verpool = sync.Pool{
		New: func() interface{} {
			return new(curve.GT)
		},
	}
	for i := 0; i < system.partynum; i++ {
		verres[i] = make([]bool, versets[i].inputsize)
		for j := 0; j < versets[i].inputsize; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				chk := chkpool.Get().(*bool)
				ver := verpool.Get().(*curve.GT)
				defer chkpool.Put(chk)
				defer verpool.Put(ver)
				ver, *chk = system.PiiSystem.System.OpenGT(*versets[i].Vers[j])
				if !*chk {
					fmt.Println("chk error")
					os.Exit(1)
				}
				verres[i][j] = bytes.Equal(ver.Marshal(), system.PiiSystem.IdentityGTbytes)
			}(i, j)
		}
	}
	wg.Wait()
	var vpool = sync.Pool{
		New: func() interface{} {
			return new([]mpc.Share_Fp)
		},
	}
	var upool = sync.Pool{
		New: func() interface{} {
			return new([]mpc.Share_GT)
		},
	}
	var tpool = sync.Pool{
		New: func() interface{} {
			return new(int)
		},
	}
	var rwMutex sync.RWMutex
	for i := 0; i < versets[0].inputsize; i++ {
		if !verres[0][i] {
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := vpool.Get().(*[]mpc.Share_Fp)
			defer vpool.Put(v)
			v = system.PiiSystem.System.SecSub(*versets[0].HIDs[i], *versets[1].HIDs[0])
			t := tpool.Get().(*int)
			defer tpool.Put(t)
			*t = 0
			for j := 1; j < size_except_for_one; j++ {
				*t = j
				k := 1
				for {
					if *t < versets[k].inputsize {
						break
					}
					*t = *t - versets[k].inputsize
					k = k + 1
				}
				if verres[k][*t] {
					if *t == 0 {
						v = system.PiiSystem.System.SecAdd(*v, *system.PiiSystem.System.SecSub(*versets[0].HIDs[i], *versets[k].HIDs[*t]))
					} else {
						v = system.PiiSystem.System.SecMul(*v, *system.PiiSystem.System.SecSub(*versets[0].HIDs[i], *versets[k].HIDs[*t]))
					}
				}
			}
			u := upool.Get().(*[]mpc.Share_GT)
			defer upool.Put(u)
			u = system.PiiSystem.System.EXP_P_GT_1(system.MK.G, v)
			u = system.PiiSystem.System.EXP_S_GT(*u, *seedsets.Seeds[i])
			uvalue, chk := system.PiiSystem.System.OpenGT(*u)
			if chk {
				if bytes.Equal(uvalue.Marshal(), system.PiiSystem.IdentityGTbytes) {
					interid, chkid := system.PiiSystem.System.OpenFp(*versets[0].HIDs[i])
					if chkid {
						rwMutex.Lock()
						intersection = append(intersection, interid)
						rwMutex.Unlock()
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
	return intersection
}

func (system *PIISystem) PartyPiiRun(inputsets []InputSet, seedsets SeedSet) []*big.Int {
	fmt.Println("ver phase start")
	vertime := time.Now()
	versets := system.verphase(inputsets)
	fmt.Println("ver time:", time.Since(vertime))
	fmt.Println("inter phase start")
	intertime := time.Now()
	intersection := system.interphase_m(versets, &seedsets)
	fmt.Println("inter time:", time.Since(intertime))
	fmt.Println("intersection:", intersection)
	fmt.Printf("Offline Communication: %f MB\n", float64(system.PiiSystem.System.OfflineCom)/1024/1024)
	fmt.Printf("Online Communication: %f MB\n", float64(system.PiiSystem.System.Com)/1024/1024)
	return intersection
}
