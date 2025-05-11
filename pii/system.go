package pii

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	ibs "github.com/Oryx/IBS"
	"github.com/Oryx/mpc"
)

type PIISystem struct {
	PiiSystem ibs.SecureVer
	partynum  int
	MK        *ibs.MasterKey
	maxID     *big.Int
}

type IDSet struct {
	IDs        []*big.Int
	Partyindex int
	inputsize  int
}

type InputSet struct {
	Sigs       []*ibs.Share_Sig
	Partyindex int
	inputsize  int
}

type VerSet struct {
	Vers       [](*[]mpc.Share_GT)
	HIDs       [](*[]mpc.Share_Fp)
	Partyindex int
	inputsize  int
}

type MsgSet struct {
	Msgs       *[][]byte
	Partyindex int
}

type SeedSet struct {
	Seeds      [](*[]mpc.Share_Fp)
	Partyindex int
}

func (system *PIISystem) generateBigIntSlice(size, intersize int) []*big.Int {
	slice := make([]*big.Int, size)
	for i := intersize; i < size; i++ {
		slice[i], _ = rand.Int(rand.Reader, system.maxID)
	}
	return slice
}

func PiiInitSystem(Partynum int, isWAN bool, bandwidth float64) *PIISystem {
	Mk := ibs.MasterKeyGen()
	piisystem := new(PIISystem)
	piisystem.PiiSystem = *ibs.SecureVerInit(Partynum, &Mk.MasterPubKey, true, isWAN, bandwidth)
	piisystem.partynum = Partynum
	one := big.NewInt(1)
	piisystem.maxID = new(big.Int).Lsh(one, 64)
	piisystem.MK = Mk
	return piisystem
}

func (system *PIISystem) prepareid(intersize int, inputsize []int) []IDSet {
	var wg sync.WaitGroup
	idsets := make([]IDSet, system.partynum)
	for i := 0; i < system.partynum; i++ {
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
			for i := 0; i < system.partynum; i++ {
				idsets[i].IDs[j] = commonElement
			}
		}(j)
	}
	wg.Wait()
	return idsets
}

func (system *PIISystem) prepareseeds(inputsize []int) []SeedSet {
	var wg sync.WaitGroup
	seedsets := make([]SeedSet, inputsize[0])
	if inputsize[0]*inputsize[1] < 16384 {
		for i := 0; i < inputsize[0]; i++ {
			seedsets[i].Seeds = make([](*[]mpc.Share_Fp), inputsize[1])
			for j := 0; j < inputsize[1]; j++ {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					seedsets[i].Seeds[j] = system.PiiSystem.System.RandomShareFp()
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
					seedsets[i].Seeds[j] = system.PiiSystem.System.RandomShareFp()
				}
			}(i)
		}
	}
	wg.Wait()
	return seedsets
}

func (system *PIISystem) prepareseeds_m(inputsize []int) *SeedSet {
	var wg sync.WaitGroup
	seedsets := new(SeedSet)
	seedsets.Seeds = make([]*[]mpc.Share_Fp, inputsize[0])
	//wg.Wait()
	for i := 0; i < inputsize[0]; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			seedsets.Seeds[i] = system.PiiSystem.System.RandomShareFp()
		}(i)
	}
	wg.Wait()
	return seedsets
}

func (system *PIISystem) prepareinput(idset []IDSet) []InputSet {
	var wg sync.WaitGroup
	blocknum := 4096
	sigsets := make([]InputSet, system.partynum)
	for i := 0; i < system.partynum; i++ {
		sigsets[i].Sigs = make([]*ibs.Share_Sig, idset[i].inputsize)
		sigsets[i].inputsize = idset[i].inputsize
		eachblock := idset[i].inputsize/blocknum + 1
		for j := 0; ; j = j + eachblock {
			if j+eachblock >= idset[i].inputsize {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					for q := j; q < idset[i].inputsize; q++ {
						uk := ibs.UserKeyGen(system.MK, idset[i].IDs[q])
						mbytes := idset[i].IDs[q].Bytes()
						sig := ibs.Sign(uk, &system.MK.MasterPubKey, mbytes)
						sigshares := system.PiiSystem.Share_A_Sig(*sig, mbytes, idset[i].IDs[q])
						sigsets[i].Sigs[q] = sigshares
					}
				}(i, j)
				break
			} else {
				wg.Add(1)
				go func(i, j int) {
					defer wg.Done()
					for q := j; q < j+eachblock; q++ {
						uk := ibs.UserKeyGen(system.MK, idset[i].IDs[q])
						mbytes := idset[i].IDs[q].Bytes()
						sig := ibs.Sign(uk, &system.MK.MasterPubKey, mbytes)
						sigshares := system.PiiSystem.Share_A_Sig(*sig, mbytes, idset[i].IDs[q])
						sigsets[i].Sigs[q] = sigshares
					}
				}(i, j)
			}

		}
	}
	wg.Wait()
	return sigsets
}

func (system *PIISystem) PrepareData(intersize int, inputsize []int) ([]InputSet, []SeedSet) {
	idsets := system.prepareid(intersize, inputsize)
	seedsets := system.prepareseeds(inputsize)
	privatesets := system.prepareinput(idsets)
	return privatesets, seedsets
}

func (system *PIISystem) PrepareData_m(intersize int, inputsize []int) ([]InputSet, *SeedSet) {
	idsets := system.prepareid(intersize, inputsize)
	seedsets := system.prepareseeds_m(inputsize)
	privatesets := system.prepareinput(idsets)
	return privatesets, seedsets
}

func (system *PIISystem) Run(inputsets []InputSet, seedset []SeedSet) {
	system.twoPartyPiiRun(inputsets, seedset)
}

func (system *PIISystem) Run_m(inputsets []InputSet, seedset *SeedSet) {
	system.PartyPiiRun(inputsets, *seedset)
}

func (system *PIISystem) GetCommunication() (float64, float64) {
	return float64(system.PiiSystem.System.OfflineCom) / 1024 / 1024, float64(system.PiiSystem.System.Com) / 1024 / 1024
}

func PIIProtocol(intersize int, inputsize []int, mode int, isWAN bool, bandwidth float64) *PIISystem {
	partynum := len(inputsize)
	piisystem := PiiInitSystem(partynum, isWAN, bandwidth)
	if mode == 0 && partynum == 2 {
		timepoint := time.Now()
		seedsets, privatesets := piisystem.PrepareData(intersize, inputsize)
		timepoint1 := time.Since(timepoint)
		fmt.Println("Data Preparation Time:", timepoint1)
		piisystem.Run(seedsets, privatesets)
	} else {
		timepoint := time.Now()
		seedsets, privatesets := piisystem.PrepareData_m(intersize, inputsize)
		timepoint1 := time.Since(timepoint)
		fmt.Println("Data Preparation Time:", timepoint1)
		piisystem.Run_m(seedsets, privatesets)
	}
	return piisystem
}
