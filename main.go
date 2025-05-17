package main

import (
	_ "net/http/pprof"

	"github.com/Oryx/benckmark"
)

func TestOryxWAN() {
	benckmark.BenckmarkFpadd()
	benckmark.BenckmarkFpmul()
	benckmark.BenckmarkFG()
	benckmark.BenckmarkFG1()
	benckmark.BenckmarkFG2()
	benckmark.BenckmarkFGT()
	benckmark.BenckmarkPair()
}

func TestPIIWAN() {
	benckmark.BenckmarkTwoPartyPIIv_AIBS_example_WAN(100)
	benckmark.BenckmarkTwoPartyPIIv_AIBS_example_WAN(500)
}

func main() {

	TestPIIWAN()
	//TestOryxWAN()
}
