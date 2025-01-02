package pii_ecdsa

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"
)

var zero = big.NewInt(0)

func (system *PIISystem) interphase_m(versets []VerSet, seedsets *SeedSet) []*big.Int {
	intersection := make([]*big.Int, 0)
	var wg sync.WaitGroup
	size_except_for_one := 0
	for i := 1; i < system.partynum; i++ {
		size_except_for_one = versets[i].inputsize + size_except_for_one
	}
	var chkpool = sync.Pool{
		New: func() interface{} {
			return new(bool)
		},
	}
	var verxpool = sync.Pool{
		New: func() interface{} {
			return new(big.Int)
		},
	}
	var verypool = sync.Pool{
		New: func() interface{} {
			return new(big.Int)
		},
	}
	verres := make([][]bool, system.partynum)
	for i := 0; i < system.partynum; i++ {
		verres[i] = make([]bool, versets[i].inputsize)
		for j := 0; j < versets[i].inputsize; j++ {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				chk := chkpool.Get().(*bool)
				verx := verxpool.Get().(*big.Int)
				very := verypool.Get().(*big.Int)
				defer chkpool.Put(chk)
				defer verxpool.Put(verx)
				defer verypool.Put(very)
				verx, _, *chk = system.PiiSystem.System.OpenG(*versets[i].Vers[j])
				if !*chk {
					fmt.Println("chk error")
					os.Exit(1)
				}
				verres[i][j] = (verx.Cmp(zero) == 0)
			}(i, j)
		}
	}
	wg.Wait()
	var rwMutex sync.RWMutex
	for i := 0; i < versets[0].inputsize; i++ {
		if !verres[0][i] {
			continue
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := system.PiiSystem.System.SecSub(*versets[0].PKxshares[i], *versets[1].PKxshares[0])
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
				if verres[k][t] {
					if t == 0 {
						v = system.PiiSystem.System.SecAdd(*v, *system.PiiSystem.System.SecSub(*versets[0].PKxshares[i], *versets[k].PKxshares[t]))
					} else {
						v = system.PiiSystem.System.SecMul(*v, *system.PiiSystem.System.SecSub(*versets[0].PKxshares[i], *versets[k].PKxshares[t]))
					}
				}
			}
			u := system.PiiSystem.System.EXP_P_G_1(system.PiiSystem.System.Curve.Gx, system.PiiSystem.System.Curve.Gy, v)
			u = system.PiiSystem.System.EXP_S_G(*u, *seedsets.Seeds[i])
			uvalueX, _, chk := system.PiiSystem.System.OpenG(*u)
			if chk {
				if bytes.Equal(uvalueX.Bytes(), system.PiiSystem.System.IdentityGx.Bytes()) {
					interid, chkid := system.PiiSystem.System.OpenFp(*versets[0].PKxshares[i])
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
	fmt.Println("ver time: ", time.Since(vertime))
	fmt.Println("inter phase start")
	intertime := time.Now()
	intersection := system.interphase_m(versets, &seedsets)
	fmt.Println("inter time: ", time.Since(intertime))
	fmt.Println("intersection:", intersection)
	fmt.Printf("Offline Communication: %f MB\n", float64(system.PiiSystem.System.OfflineCom)/1024/1024)
	fmt.Printf("Online Communication: %f MB\n", float64(system.PiiSystem.System.Com)/1024/1024)
	return intersection
}
