package ecdsa

import (
	"math/big"

	"github.com/Oryx/mpc"
	"github.com/Oryx/shmpc"
)

type SecureVer struct {
	System     mpc.ECCShareSystem
	SemiSystem shmpc.ECCShareSystem
	Ecdsa      *ECDSA
	Security   bool
}

type Share_Sig struct {
	RX          *big.Int
	RY          *big.Int
	S           *big.Int
	Hmshare     *[]mpc.Share_Fp
	Pkshare     *[]mpc.Share_G
	SemiHmshare *[]shmpc.Share_Fp
	SemiPkshare *[]shmpc.Share_G
}

func SecureVerInit(Partynum int, ismalicious bool) *SecureVer {
	securever := new(SecureVer)
	securever.Ecdsa = NewECDSA()
	if ismalicious {
		securever.System = *mpc.ECCSystemInit(Partynum)
		securever.Security = true
	} else {
		securever.SemiSystem = *shmpc.ECCSystemInit(Partynum)
		securever.Security = false
	}
	return securever
}

func (securever *SecureVer) Share_A_Sig(sig SigInv, pk *PublicKey) *Share_Sig {
	share_sig := new(Share_Sig)
	share_sig.S = sig.S
	share_sig.RX = sig.RX
	share_sig.RY = sig.RY
	if securever.Security {
		share_sig.Pkshare = securever.System.Share_A_G_Offline(pk.PKX, pk.PKY)
		share_sig.Hmshare = securever.System.Share_An_Fp_Offline(sig.HM)
	} else {
		share_sig.SemiPkshare = securever.SemiSystem.Share_A_G_Offline(pk.PKX, pk.PKY)
		share_sig.SemiHmshare = securever.SemiSystem.Share_An_Fp_Offline(sig.HM)
	}
	return share_sig
}

func (securever *SecureVer) SecVer(sigshares *Share_Sig) (bool, bool) {
	if securever.Security {
		u1shares := securever.System.SecMulPlaintext(*sigshares.Hmshare, sigshares.S)
		u2 := new(big.Int).Mul(sigshares.RX, sigshares.S)
		P1 := securever.System.EXP_P_G_2(sigshares.Pkshare, u2)
		P2 := securever.System.EXP_P_G_1(securever.Ecdsa.curve.Gx, securever.Ecdsa.curve.Gy, u1shares)
		P := securever.System.SecAdd_G(*P1, *P2)
		resx, resy, chk := securever.System.OpenG(*P)
		return resx.Cmp(sigshares.RX) == 0 && resy.Cmp(sigshares.RY) == 0, chk
	}
	u1shares := securever.SemiSystem.SecMulPlaintext(*sigshares.SemiHmshare, sigshares.S)
	u2 := new(big.Int).Mul(sigshares.RX, sigshares.S)
	P1 := securever.SemiSystem.EXP_P_G_2(sigshares.SemiPkshare, u2)
	P2 := securever.SemiSystem.EXP_P_G_1(securever.Ecdsa.curve.Gx, securever.Ecdsa.curve.Gy, u1shares)
	P := securever.SemiSystem.SecAdd_G(*P1, *P2)
	resx, resy := securever.SemiSystem.OpenG(*P)
	return resx.Cmp(sigshares.RX) == 0 && resy.Cmp(sigshares.RY) == 0, true
}

func (securever *SecureVer) SecVerWithoutOpen(sigshares *Share_Sig) *[]mpc.Share_G {
	u1shares := securever.System.SecMulPlaintext(*sigshares.Hmshare, sigshares.S)
	u2 := new(big.Int).Mul(sigshares.RX, sigshares.S)
	P1 := securever.System.EXP_P_G_2(sigshares.Pkshare, u2)
	P2 := securever.System.EXP_P_G_1(securever.Ecdsa.curve.Gx, securever.Ecdsa.curve.Gy, u1shares)
	P := securever.System.SecAdd_G(*P1, *P2)
	R := securever.System.SecSubPlaintext_G(*P, sigshares.RX, sigshares.RY)
	return R
}

func (securever *SecureVer) SemiSecVerWithoutOpen(sigshares *Share_Sig) *[]shmpc.Share_G {
	u1shares := securever.SemiSystem.SecMulPlaintext(*sigshares.SemiHmshare, sigshares.S)
	u2 := new(big.Int).Mul(sigshares.RX, sigshares.S)
	P1 := securever.SemiSystem.EXP_P_G_2(sigshares.SemiPkshare, u2)
	P2 := securever.SemiSystem.EXP_P_G_1(securever.Ecdsa.curve.Gx, securever.Ecdsa.curve.Gy, u1shares)
	P := securever.SemiSystem.SecAdd_G(*P1, *P2)
	R := securever.SemiSystem.SecSubPlaintext_G(*P, sigshares.RX, sigshares.RY)
	return R
}
