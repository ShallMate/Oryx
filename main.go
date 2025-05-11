package main

import (
	_ "net/http/pprof"

	"github.com/Oryx/benckmark"
)

func testG() {
	benckmark.TestMaliciousHalfOpenG_WAN(500)
	//benckmark.TestMaliciousHalfOpenG_WAN(500)
	//benckmark.TestMaliciousHalfOpenG_WAN(1000)
	//benckmark.TestMaliciousOpenG_WAN(100)
	//benckmark.TestMaliciousOpenG_WAN(500)
	//benckmark.TestMaliciousOpenG_WAN(1000)

}

func TestPII_AIBS() {
	benckmark.TwoPartyPIIv_AIBS_example_WAN(1000)
}

func main() {

	//benckmark.BenckmarkSecVerAIBS()
	//benckmark.TestMaliciousOpenG()
	//benckmark.TestSecVerECDSA()
	//benckmark.TestSecVerBLS()
	//benckmark.TestSecVerAIBS()
	//benckmark.BenckmarkAIBS()
	//benckmark.TwoPartyPII_BLS_example()
	//benckmark.TwoPartyPII_ECDSA_example()
	//benckmark.MutiPartyPIIv_AIBS_example()
	//benckmark.TwoPartyPII_AIBS_example()
	//benckmark.TestMaliciousSecMul()
	TestPII_AIBS()
}
