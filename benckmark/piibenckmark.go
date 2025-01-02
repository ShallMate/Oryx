package benckmark

import (
	"fmt"
	"time"

	"github.com/Oryx/pii"
	"github.com/Oryx/pii_bls"
	"github.com/Oryx/pii_ecdsa"
)

// PII based on Our AIBS
func TwoPartyPII_AIBS_example() {
	// Input size
	inputsize := []int{32, 32}

	// Intersection size
	intersize := 10
	t1 := time.Now()

	// If the third parameter is 0, it is the PII protocol, and if it is 1, it is the PIIv protocol.
	pii.PIIProtocol(intersize, inputsize, 0)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PII based on Our AIBS
func TwoPartyPIIv_AIBS_example() {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii.PIIProtocol(intersize, inputsize, 1)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// If you want to run multi-party PIIv
func MutiPartyPIIv_AIBS_example() {
	inputsize := []int{32, 32, 32}
	intersize := 10
	t1 := time.Now()
	pii.PIIProtocol(intersize, inputsize, 1)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PII based on BLS
func TwoPartyPII_BLS_example() {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii_bls.PIIProtocol(intersize, inputsize, 0)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PIIv based on BLS
func TwoPartyPIIv_BLS_example() {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii_bls.PIIProtocol(intersize, inputsize, 1)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PII based on ECDSA
func TwoPartyPII_ECDSA_example() {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii_ecdsa.PIIProtocol(intersize, inputsize, 0)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PIIv based on ECDSA
func TwoPartyPIIv_ECDSA_example() {
	inputsize := []int{32, 32}
	intersize := 10
	t1 := time.Now()
	pii_ecdsa.PIIProtocol(intersize, inputsize, 1)
	t2 := time.Since(t1)
	fmt.Println(t2)
}
