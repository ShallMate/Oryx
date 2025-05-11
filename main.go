package main

import (
	_ "net/http/pprof"

	"github.com/Oryx/benckmark"
)

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
	benckmark.TestMaliciousHalfOpenG_WAN(500)
	//benckmark.TwoPartyPII_AIBS_example()
	//benckmark.TestMaliciousSecMul()

}
