package benckmark

import (
	"fmt"
	"time"

	"github.com/Oryx/pii"
	"github.com/Oryx/pii_bls"
	"github.com/Oryx/pii_ecdsa"
)

// PII based on Our AIBS
func TwoPartyPII_AIBS_example_WAN(bandwidth float64) {
	// Input size
	inputsize := []int{32, 32}

	// Intersection size
	intersize := 10
	t1 := time.Now()

	// If the third parameter is 0, it is the PII protocol, and if it is 1, it is the PIIv protocol.
	pii.PIIProtocol(intersize, inputsize, 0, true, bandwidth)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PII based on Our AIBS
func TwoPartyPIIv_AIBS_example_WAN(bandwidth float64) {
	inputsize := []int{10, 10}
	intersize := 5
	t1 := time.Now()
	pii.PIIProtocol(intersize, inputsize, 1, true, bandwidth)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PII based on Our AIBS
func BenckmarkTwoPartyPIIv_AIBS_example_WAN(bandwidth float64) {
	inputsizetests := []int{10, 20, 50, 100, 200, 500, 1000}
	for i := 0; i < len(inputsizetests); i++ {
		inputsize := []int{inputsizetests[i], inputsizetests[i]}
		intersize := inputsizetests[i] / 2
		t1 := time.Now()
		pii.PIIProtocol(intersize, inputsize, 1, true, bandwidth)
		t2 := time.Since(t1)
		fmt.Println(t2)
	}
}

// If you want to run multi-party PIIv
func BenckmarkMutiPartyPIIv_AIBS_example_WAN(bandwidth float64) {
	inputsizetests := []int{10, 20, 50, 100, 200, 500, 1000}
	for i := 0; i < len(inputsizetests); i++ {
		for j := 3; j <= 10; j++ {
			inputsize := make([]int, 0, j)
			for k := 0; k < j; k++ {
				inputsize = append(inputsize, inputsizetests[i])
			}
			intersize := inputsizetests[i] / 2
			t1 := time.Now()
			pii.PIIProtocol(intersize, inputsize, 1, true, bandwidth)
			t2 := time.Since(t1)
			fmt.Println(t2)
		}
	}
}

// PII based on BLS
func TwoPartyPII_BLS_example_WAN(bandwidth float64) {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii_bls.PIIProtocol(intersize, inputsize, 0, true, bandwidth)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PIIv based on BLS
func TwoPartyPIIv_BLS_example_WAN(bandwidth float64) {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii_bls.PIIProtocol(intersize, inputsize, 1, true, bandwidth)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PIIv based on BLS
func BenchmarkTwoPartyPIIv_BLS_example_WAN(bandwidth float64) {
	inputsizetests := []int{10, 20, 50, 100, 200, 500, 1000}
	for i := 0; i < len(inputsizetests); i++ {
		inputsize := []int{inputsizetests[i], inputsizetests[i]}
		intersize := inputsizetests[i] / 2
		t1 := time.Now()
		pii_bls.PIIProtocol(intersize, inputsize, 1, true, bandwidth)
		t2 := time.Since(t1)
		fmt.Println(t2)
	}
}

// PII based on ECDSA
func TwoPartyPII_ECDSA_example_WAN(bandwidth float64) {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii_ecdsa.PIIProtocol(intersize, inputsize, 0, true, bandwidth)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PII based on ECDSA
func BenchmarkTwoPartyPII_ECDSA_example_WAN(bandwidth float64) {
	inputsizetests := []int{10, 20, 50, 100, 200, 500, 1000}
	for i := 0; i < len(inputsizetests); i++ {
		inputsize := []int{inputsizetests[i], inputsizetests[i]}
		intersize := inputsizetests[i] / 2
		t1 := time.Now()
		pii_ecdsa.PIIProtocol(intersize, inputsize, 1, true, bandwidth)
		t2 := time.Since(t1)
		fmt.Println(t2)
	}
}

// PIIv based on ECDSA
func TwoPartyPIIv_ECDSA_example_WAN(bandwidth float64) {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii_ecdsa.PIIProtocol(intersize, inputsize, 1, true, bandwidth)
	t2 := time.Since(t1)
	fmt.Println(t2)
}
