# Oryx: An efficient and user-friendly MPC framework without the need for any third party libraries

This repository is the complete implementation of the paper: **<span style="color:red">Privacy-Preserving Authorized Set Matching via Dishonest Majority Multiparty Computation**</span>, which allows readers to understand our protocol and reproduce our results.

$\mathsf{Oryx}$ is a library that was implemented to implement the **<span style="color:red">Private Identity Intersection (PII)**</span> protocols. It implements arbitrary calculations on the prime number field $\mathbb{F}_p$, the elliptic curve field $\mathbb{G}$, and the bilinear group ($\mathbb{G}_1$, $\mathbb{G}_2$, $\mathbb{G}_T$).

**<span style="color:red">PII is a new type of Private Matching, similar to Private Set Intersection (PSI), but with more robust security; it guarantees the authenticity and integrity of input and output, and the entire calculation process is verifiable. The output of PII is sent to all parties. At the same time, we provide two versions. The two-party and multi-party constructions use multiple signature schemes. For details, please see our paper.**</span>

**<span style="color:red">You can find all the experimental code in our paper in the benchmark folder to reproduce our results.**</span>

## Advantages of PII

1. **Authenticity of Input**: Each input in the PII protocol requires certification by a user to ensure the authenticity of the inputs. Therefore, the inputs for computing the intersection in the PII protocol are virtually user IDs. 
2. **Integrity of Output**: Given the inputs of the parties, the PII protocol is guaranteed to output the complete intersection and is resistant to $n-1$ collusion attacks, where $n$ denotes the number of parties.
3. **Mutuality and Scalability**: The output of the PII protocol is provided to all parties, and each party is peer-to-peer, allowing for natural extension from two parties to multiple parties.

## Advantages of Oryx

1. **No Third-Party Libraries Required**: $\mathsf{Oryx}$ does not require any third-party libraries. It can be used directly by installing the Golang environment, avoiding the complex configurations of existing frameworks.
2. **Support for Different Security Settings**: $\mathsf{Oryx}$ supports both semi-honest and malicious settings and does not limit the number of participants. It can conveniently measure the protocol's runtime and communication volume.
3. **User-Friendly**: $\mathsf{Oryx}$ is very friendly to researchers. Even those who are not MPC experts can use $\mathsf{Oryx}$. Additionally, it is easy to convert a cryptographic scheme implemented in the prime field $\mathbb{F}_p$, the elliptic curve field $\mathbb{G}$, or the bilinear groups ($\mathbb{G}_1$, $\mathbb{G}_2$, $\mathbb{G}_T$) into an MPC form.

## Quick start

**<span style="color:red">We assume that the user has a local Golang environment and our language version is Golang1.21.1.**</span> Unfortunately, we have not conducted an in-depth evaluation of other versions to see if there are any version conflicts. However, since we do not use any third-party libraries, I think there should not be any major problems. **<span style="color:red">If there is a version conflict, users can use a Golang version that is close to ours.**</span>

```bash
git clone https://github.com/ShallMate/Oryx.git
cd Oryx
go run main.go
```

## Using Docker

```bash
git clone https://github.com/ShallMate/Oryx.git
cd Oryx
docker build -t pii:latest .
docker run -it pii bash
go run main.go
```

## How to run PII

Below are various examples of how to perform various PII protocols in the benckmark/piibenckmark.go.

```go
// PII based on Our AIBS
func TwoPartyPII_AIBS_example() {
	// Input size
	inputsize := []int{64, 64}

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
	inputsize := []int{64, 64}
	intersize := 10
	t1 := time.Now()
	pii.PIIProtocol(intersize, inputsize, 1)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// If you want to run multi-party PIIv
func MutiPartyPIIv_AIBS_example() {
	inputsize := []int{64, 64, 64}
	intersize := 10
	t1 := time.Now()
	pii.PIIProtocol(intersize, inputsize, 1)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PII based on BLS
func TwoPartyPII_BLS_example() {
	inputsize := []int{64, 64}
	intersize := 10
	t1 := time.Now()
	pii_bls.PIIProtocol(intersize, inputsize, 0)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PIIv based on BLS
func TwoPartyPIIv_BLS_example() {
	inputsize := []int{64, 64}
	intersize := 10
	t1 := time.Now()
	pii_bls.PIIProtocol(intersize, inputsize, 1)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PII based on ECDSA
func TwoPartyPII_ECDSA_example() {
	inputsize := []int{64, 64}
	intersize := 10
	t1 := time.Now()
	pii_ecdsa.PIIProtocol(intersize, inputsize, 0)
	t2 := time.Since(t1)
	fmt.Println(t2)
}

// PIIv based on ECDSA
func TwoPartyPIIv_ECDSA_example() {
	inputsize := []int{64, 64}
	intersize := 10
	t1 := time.Now()
	pii_ecdsa.PIIProtocol(intersize, inputsize, 1)
	t2 := time.Since(t1)
	fmt.Println(t2)
}
```

## How to run the secure signature verification

Secure signature verification is one of the essential components of constructing PII. Simply put, all signature verification calculation operations are performed in MPC through the Oryx framework. In PII, secure signature verification is required to be maliciously secure. Here, we also implement the protocols in the semi-honest model because secure signature verification may not be limited to PII.

### 1. How to run secure signature verification using ECDSA
```go
func TestSecVerECDSA() {
	// Parameter 1 is the number of parties, and parameter 2 denotes whether it is the malicious model
	system := ecdsa.SecureVerInit(2, true)

	// Setup
	eccsystem := ecdsa.NewECDSA()

	// KeyGen
	sk, pk := eccsystem.KeyGen()
	msg := "hello world"

	// Sign
	sig := eccsystem.SignwithInv(sk, []byte(msg))

	// How to share
	sigshares := system.Share_A_Sig(*sig, pk)

	// Verify both the signature and MAC
	// If it is the semi-honest model, chk will always be 1
	ver, chk := system.SecVer(sigshares)
	if chk {
		if ver {
			fmt.Println("Signature verification successful")
		} else {
			fmt.Println("Signature verification failed")
		}
	}
	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.System.OfflineCom+system.System.OfflineCom)

	// If it is the semi-honest model
	// fmt.Println("The communication: ", system.SemiSystem.OfflineCom+system.SemiSystem.OfflineCom)
}
```

### 2. How to run secure signature verification using BLS
```go
func TestSecVerBLS() {
	// Parameter 1 is the number of parties, and parameter 2 denotes whether it is the malicious model
	system := bls.SecureVerInit(2, true)

	// KeyGen
	sk, pk := bls.KeyGen()
	msg := "hello world"

	// Sign
	sig, hm := bls.SignwithHm(sk, []byte(msg))

	// How to share
	sigshares := system.Share_A_Sig(sig, hm, pk)

	// Verify both the signature and MAC
	// If it is the semi-honest model, chk will always be 1
	ver, chk := system.SecVer(sigshares)
	if chk {
		if ver {
			fmt.Println("Signature verification successful")
		} else {
			fmt.Println("Signature verification failed")
		}
	}
	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.System.OfflineCom+system.System.OfflineCom)

	// If it is the semi-honest model
	// fmt.Println("The communication: ", system.SemiSystem.OfflineCom+system.SemiSystem.OfflineCom)
}
```

### 3. How to run secure signature verification using AIBS
```go
func TestSecVerAIBS() {
	msk := ibs.MasterKeyGen()

	// Parameter 1 is the number of parties
	// parameter 2 denotes the MPK
	// Parameter 3 denotes whether it is the malicious model
	system := *ibs.SecureVerInit(2, &msk.MasterPubKey, true)

	// id
	userid := big.NewInt(9567)
	// KeyGen
	sk := ibs.UserKeyGen(msk, userid)

	msg := "hello world"

	// Sign
	sig := ibs.Sign(sk, &msk.MasterPubKey, []byte(msg))

	// How to share
	sigshares := system.Share_A_Sig(*sig, []byte(msg), userid)

	// Verify both the signature and MAC
	// If it is the semi-honest model, chk will always be 1
	ver, chk := system.SecVer(sigshares)
	if chk {
		if ver {
			fmt.Println("Signature verification successful")
		} else {
			fmt.Println("Signature verification failed")
		}
	}
	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.System.OfflineCom+system.System.OfflineCom)

	// If it is the semi-honest model
	// fmt.Println("The communication: ", system.SemiSystem.OfflineCom+system.SemiSystem.OfflineCom)
}
```

## How to use Oryx

The following are various simple examples. **<span style="color:red">You can directly replace the content in main.go with the following content and run it directly. Since there are too many functions, we can only select a few to show.**</span> Furthermore, you can also choose to refer to the code in the benckmark folder, which also contains the usage code of various functions of $\mathsf{Oryx}$.


### 1. How to run [x]+[y] on $\mathbb{F}_p$ for two-party setting in the malicious model
```go
package main

import (
	"fmt"
	"math/big"
	_ "net/http/pprof"

	"github.com/Oryx/mpc"
)

func main() {
	e1 := big.NewInt(100)
	e2 := big.NewInt(200)

    // The parameter is the number of parties
	system := mpc.SystemInit(2)

    // If you want to run in the semi-honest model
    //system := shmpc.SystemInit(2)

	// How to Share an Fp with Additive Secret Sharing
	shares1 := system.Share_An_Fp(e1)
	shares2 := system.Share_An_Fp(e2)

	// How to use secure addition on Fp
	shares3 := system.SecAdd(*shares1, *shares2)
	rese, chk := system.OpenFp(*shares3)
	if chk {
		fmt.Printf("[%s]+[%s] = [%s]\n", e1, e2, rese)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}
```

### 2. How to run [x]*[y] on $\mathbb{F}_p$ for two-party setting in the semi-honest model
```go
package main

import (
	"fmt"
	"math/big"
	_ "net/http/pprof"

	"github.com/Oryx/shmpc"
)

func main() {
	e1 := big.NewInt(100)
	e2 := big.NewInt(200)

	// The parameter is the number of parties
	system := shmpc.SystemInit(2)

	// How to Share an Fp with Additive Secret Sharing
	shares1 := system.Share_An_Fp(e1)
	shares2 := system.Share_An_Fp(e2)

	// How to use secure multiplication on Fp
	shares3 := system.SecMul(*shares1, *shares2)
	rese := system.OpenFp(*shares3)

	fmt.Printf("[%s]*[%s] = [%s]\n", e1, e2, rese)

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}
```

### 3. How to run $\langle x\rangle*\langle y\rangle$ on $\mathbb{F}_p$ for multiparty setting in the malicious model
```go
package main

import (
	"fmt"
	"math/big"
	_ "net/http/pprof"

	"github.com/Oryx/mpc"
)

func main() {
	e1 := big.NewInt(100)
	e2 := big.NewInt(200)

	// The parameter is the number of parties
	system := mpc.SystemInit(10)

	// If you want to run in the semi-honest model
	//system := shmpc.SystemInit(2)

	// How to Share an Fp with Multiplication Secret Sharing
	shares1 := system.Share_An_Fp_Mul(e1)
	shares2 := system.Share_An_Fp_Mul(e2)

	// How to use secure multiplication using MSS on Fp
	shares3 := system.SecMul_Mul(*shares1, *shares2)

	// We need to use the corresponding opening protocol here
	rese, chk := system.OpenFp_Mul(*shares3)
	if chk {
		fmt.Printf("<%s>+<%s> = <%s>\n", e1, e2, rese)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}

```

### 4. How to run $[P]+[Q]$ on the ECC group for two-party setting in the malicious model
```go
package main

import (
	"fmt"
	_ "net/http/pprof"

	"github.com/Oryx/mpc"
)

func main() {

	// How to initialize ECC system
	system := mpc.ECCSystemInit(2)

	e1x, e1y := system.RandomG()
	e2x, e2y := system.RandomG()

	// If you want to run in the semi-honest model
	//system := shmpc.ECCSystemInit(2)

	// How to Share an element on the ECC group
	shares1 := system.Share_A_G(e1x, e1y)
	shares2 := system.Share_A_G(e2x, e2y)

	// How to use secure point addition
	shares3 := system.SecAdd_G(*shares1, *shares2)

	// We need to use the corresponding opening protocol here
	resex, resey, chk := system.OpenG(*shares3)
	if chk {
		fmt.Printf("[(%s,%s)]+[(%s,%s)] = [(%s,%s)]\n", e1x, e1y, e2x, e2y, resex, resey)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}
```

### 5. How to run $[x]*[P]$ on the ECC group for two-party setting in the malicious model
```go
package main

import (
	"crypto/rand"
	"fmt"
	_ "net/http/pprof"

	"github.com/Oryx/mpc"
)

func main() {

	// How to initialize ECC system
	system := mpc.ECCSystemInit(2)

	x1, _ := rand.Int(rand.Reader, system.Order)
	e2x, e2y := system.RandomG()

	// If you want to run in the semi-honest model
	//system := shmpc.ECCSystemInit(2)

	// How to share an element on Fp
	shares1 := system.Share_An_Fp(x1)

	// How to share an element on the ECC group
	shares2 := system.Share_A_G(e2x, e2y)

	// How to use secure scalar multiplication
	shares3 := system.EXP_S_G(*shares2, *shares1)

	// We need to use the corresponding opening protocol here
	resex, resey, chk := system.OpenG(*shares3)
	if chk {
		fmt.Printf("[%s]*[(%s,%s)] = [(%s,%s)]\n", x1, e2x, e2y, resex, resey)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}

```

### 6. How to run $[x]*[P]$ on the BP group G1 for two-party setting in the malicious model
```go
package main

import (
	"crypto/rand"
	"fmt"
	_ "net/http/pprof"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
)

func main() {

	// How to initialize BP system
	system := mpc.SystemInit(2)

	x, _ := rand.Int(rand.Reader, system.Order)
	_, e, _ := curve.RandomG1(rand.Reader)

	// If you want to run in the semi-honest model
	//system := shmpc.SystemInit(2)

	// How to share an element on Fp
	shares1 := system.Share_An_Fp(x)

	// How to share an element on the BP group G_1
	shares2 := system.Share_A_G1(e)

	// How to use secure scalar multiplication on G_1
	shares3 := system.EXP_S_G1(*shares2, *shares1)

	// We need to use the corresponding opening protocol here
	rese, chk := system.OpenG1(*shares3)
	if chk {
		fmt.Printf("[%s]*[%s] = [%s]\n", x, e, rese)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}
```

### 7. How to run $[x]*[P]$ on the BP group G2 for two-party setting in the malicious model
```go
package main

import (
	"crypto/rand"
	"fmt"
	_ "net/http/pprof"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
)

func main() {

	// How to initialize BP system
	system := mpc.SystemInit(2)

	x, _ := rand.Int(rand.Reader, system.Order)
	_, e, _ := curve.RandomG2(rand.Reader)

	// If you want to run in the semi-honest model
	//system := shmpc.SystemInit(2)

	// How to share an element on Fp
	shares1 := system.Share_An_Fp(x)

	// How to share an element on the BP group G_2
	shares2 := system.Share_A_G2(e)

	// How to use secure scalar multiplication on G_2
	shares3 := system.EXP_S_G2(*shares2, *shares1)

	// We need to use the corresponding opening protocol here
	rese, chk := system.OpenG2(*shares3)
	if chk {
		fmt.Printf("[%s]*[%s] = [%s]\n", x, e, rese)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}

```

### 8. How to run $[x]*[P]$ on the BP group GT for two-party setting in the malicious model
```go
package main

import (
	"crypto/rand"
	"fmt"
	_ "net/http/pprof"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
)

func main() {

	// How to initialize BP system
	system := mpc.SystemInit(2)

	_, g1, _ := curve.RandomG1(rand.Reader)
    _, g2, _ := curve.RandomG1(rand.Reader)


	// If you want to run in the semi-honest model
	//system := shmpc.SystemInit(2)

	// How to share an element on G1
	shares1 := system.Share_An_G1(g1)

	//  How to share an element on G2
	shares2 := system.Share_A_G2(g2)

	// How to use secure scalar multiplication on G_T
	shares3 := system.EXP_S_GT(*shares2, *shares1)

	// We need to use the corresponding opening protocol here
	rese, chk := system.OpenGT(*shares3)
	if chk {
		fmt.Printf("[%s]*[%s] = [%s]\n", x, e, rese)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}
```

### 9. How to run the secure pairing $e([P],Q)$  for two-party setting in the malicious model
```go
package main

import (
	"crypto/rand"
	"fmt"
	_ "net/http/pprof"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
)

func main() {

	// How to initialize BP system
	system := mpc.SystemInit(2)

	_, g1, _ := curve.RandomG1(rand.Reader)
	_, g2, _ := curve.RandomG2(rand.Reader)

	// If you want to run in the semi-honest model
	//system := shmpc.SystemInit(2)

	// How to share an element on G1
	shares1 := system.Share_A_G1(g1)

	// How to use secure pairing e([P],Q) on G_T
	shares3 := system.Pair_P_1(shares1, g2)

	// We need to use the corresponding opening protocol here
	rese, chk := system.OpenGT(*shares3)
	if chk {
		fmt.Printf("e([%s],[%s]) = [%s]\n", g1, g2, rese)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}
```

### 10. How to run the secure pairing $e([P],[Q])$  for two-party setting in the malicious model
```go
package main

import (
	"crypto/rand"
	"fmt"
	_ "net/http/pprof"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
)

func main() {

	// How to initialize BP system
	system := mpc.SystemInit(2)

	_, g1, _ := curve.RandomG1(rand.Reader)
	_, g2, _ := curve.RandomG2(rand.Reader)

	// If you want to run in the semi-honest model
	//system := shmpc.SystemInit(2)

	// How to share an element on G1
	shares1 := system.Share_A_G1(g1)

	//  How to share an element on G2
	shares2 := system.Share_A_G2(g2)

	// How to use secure pairing e([P],[Q]) on G_T
	shares3 := system.Pair_S(shares1, shares2)

	// We need to use the corresponding opening protocol here
	rese, chk := system.OpenGT(*shares3)
	if chk {
		fmt.Printf("e([%s],[%s]) = [%s]\n", g1, g2, rese)
	}

	// If you want to print the communication, the unit is bytes
	fmt.Println("The communication: ", system.OfflineCom+system.Com)
}
```

## NOTE

Oryx is mainly used for scientific research. Please do not use it in production environments. In addition, due to my limited knowledge level, please forgive me if there are a few bugs or non-standard programming here. If you encounter any problems when using this library, you can ask questions about the issues or contact me directly at gw_ling@sjtu.edu.cn. If you use Oryx in your research, please cite this library. 