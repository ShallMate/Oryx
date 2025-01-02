package pm

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Oryx/mpc"
)

type PMSystem struct {
	System *mpc.ShareSystem
	maxID  *big.Int
	zero   *big.Int
}

type IDSet struct {
	IDs        []*big.Int
	Partyindex int
	inputsize  int
}

type InputSet struct {
	HIDs       [](*[]mpc.Share_Fp)
	Partyindex int
	inputsize  int
}

type SeedSet struct {
	Seeds      [](*[]mpc.Share_Fp)
	Partyindex int
}

func (system *PMSystem) generateBigIntSlice(size, intersize int) []*big.Int {
	slice := make([]*big.Int, size)
	for i := intersize; i < size; i++ {
		slice[i], _ = rand.Int(rand.Reader, system.maxID)
	}
	return slice
}

func PMInitSystem(Partynum int) *PMSystem {
	system := new(PMSystem)
	system.System = mpc.SystemInit(Partynum)
	one := big.NewInt(1)
	system.maxID = new(big.Int).Lsh(one, 64)
	system.zero = big.NewInt(0)
	return system
}

func (system *PMSystem) prepareid(intersize int, inputsize []int) []IDSet {
	var wg sync.WaitGroup
	idsets := make([]IDSet, system.System.Partynum)
	for i := 0; i < system.System.Partynum; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			idsets[i].IDs = system.generateBigIntSlice(inputsize[i], intersize)
			idsets[i].Partyindex = i
			idsets[i].inputsize = inputsize[i]
		}(i)
	}
	wg.Wait()
	for j := 0; j < intersize; j++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			commonElement, _ := rand.Int(rand.Reader, system.maxID)
			for i := 0; i < system.System.Partynum; i++ {
				idsets[i].IDs[j] = commonElement
			}
		}(j)
	}
	wg.Wait()
	return idsets
}

func (system *PMSystem) prepareseeds(inputsize []int) []SeedSet {
	var wg sync.WaitGroup
	seedsets := make([]SeedSet, inputsize[0])
	if inputsize[0]*inputsize[1] < 16384 {
		for i := 0; i < inputsize[0]; i++ {
			seedsets[i].Seeds = make([](*[]mpc.Share_Fp), inputsize[1])
			for j := 0; j < inputsize[1]; j++ {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					seedsets[i].Seeds[j] = system.System.RandomShareFp()
				}(i, j)
			}
		}
	} else {
		for i := 0; i < inputsize[0]; i++ {
			seedsets[i].Seeds = make([](*[]mpc.Share_Fp), inputsize[1])
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < inputsize[1]; j++ {
					seedsets[i].Seeds[j] = system.System.RandomShareFp()
				}
			}(i)
		}
	}
	wg.Wait()
	return seedsets
}

func (system *PMSystem) prepareseeds_m(inputsize []int) *SeedSet {
	var wg sync.WaitGroup
	seedsets := new(SeedSet)
	seedsets.Seeds = make([]*[]mpc.Share_Fp, inputsize[0])
	//wg.Wait()
	for i := 0; i < inputsize[0]; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			seedsets.Seeds[i] = system.System.RandomShareFp()
		}(i)
	}
	wg.Wait()
	return seedsets
}

func (system *PMSystem) PrepareData(intersize int, inputsize []int) ([]InputSet, []SeedSet) {
	idsets := system.prepareid(intersize, inputsize)
	//fmt.Println(idsets)
	seedsets := system.prepareseeds(inputsize)
	//fmt.Println(seedsets)
	privatesets := system.prepareinput(idsets)
	//fmt.Println(privatesets)
	return privatesets, seedsets
}

func (system *PMSystem) PrepareData_m(intersize int, inputsize []int) ([]InputSet, *SeedSet) {
	idsets := system.prepareid(intersize, inputsize)
	seedsets := system.prepareseeds_m(inputsize)
	privatesets := system.prepareinput(idsets)
	return privatesets, seedsets
}

func (system *PMSystem) prepareinput(idset []IDSet) []InputSet {
	var wg sync.WaitGroup
	blocknum := 4096
	sigsets := make([]InputSet, system.System.Partynum)
	for i := 0; i < system.System.Partynum; i++ {
		sigsets[i].HIDs = make([](*[]mpc.Share_Fp), idset[i].inputsize)
		sigsets[i].inputsize = idset[i].inputsize
		sigsets[i].Partyindex = i
		eachblock := idset[i].inputsize/blocknum + 1
		for j := 0; ; j = j + eachblock {
			if j+eachblock >= idset[i].inputsize {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					for q := j; q < idset[i].inputsize; q++ {
						sigsets[i].HIDs[q] = system.System.Share_An_Fp_Offline(idset[i].IDs[q])
					}
				}(i, j)
				break
			} else {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					for q := j; q < j+eachblock; q++ {
						sigsets[i].HIDs[q] = system.System.Share_An_Fp_Offline(idset[i].IDs[q])
					}
				}(i, j)
			}

		}
	}
	wg.Wait()
	//fmt.Println(sigsets)
	return sigsets
}

func (system *PMSystem) Run(inputsets []InputSet, seedset []SeedSet) {
	//fmt.Println(inputsets)
	system.twoPartyRun(inputsets, seedset)
}

func (system *PMSystem) Run_m(inputsets []InputSet, seedset *SeedSet) {
	system.PartyRun(inputsets, *seedset)
}

func (system *PMSystem) GetCommunication() (float64, float64) {
	return float64(system.System.OfflineCom) / 1024 / 1024, float64(system.System.Com) / 1024 / 1024
}

func PMProtocol(intersize int, inputsize []int, mode int) *PMSystem {
	partynum := len(inputsize)
	system := PMInitSystem(partynum)
	if mode == 0 && partynum == 2 {
		timepoint := time.Now()
		seedsets, privatesets := system.PrepareData(intersize, inputsize)
		timepoint1 := time.Since(timepoint)
		fmt.Println("Data Preparation Time:", timepoint1)
		system.Run(seedsets, privatesets)
	} else {
		timepoint := time.Now()
		seedsets, privatesets := system.PrepareData_m(intersize, inputsize)
		timepoint1 := time.Since(timepoint)
		fmt.Println("Data Preparation Time:", timepoint1)
		system.Run_m(seedsets, privatesets)
	}
	return system
}
