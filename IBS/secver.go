package ibs

import (
	"bytes"
	"math/big"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
	"github.com/Oryx/shmpc"
)

type SecureVer struct {
	mpk             *MasterPubKey
	IdentityGTbytes []byte
	System          mpc.ShareSystem
	mpkshare        *[]mpc.Share_G2
	SemiSystem      shmpc.ShareSystem
	Semimpkshare    *[]shmpc.Share_G2
	Security        bool
}

type Share_Sig struct {
	S       Sig
	HM      *[]mpc.Share_Fp
	HID     *[]mpc.Share_Fp
	HS1     *[]mpc.Share_Fp
	SemiHM  *[]shmpc.Share_Fp
	SemiHID *[]shmpc.Share_Fp
	SemiHS1 *[]shmpc.Share_Fp
}

func SecureVerInit(Partynum int, mpk *MasterPubKey, ismalicious bool, isWAN bool, bandwidth float64) *SecureVer {
	securever := new(SecureVer)
	securever.Security = ismalicious
	securever.mpk = mpk
	if ismalicious {
		if isWAN {
			securever.System = *mpc.SystemInitWAN(Partynum, bandwidth)
		} else {
			securever.System = *mpc.SystemInit(Partynum)
		}
		securever.mpkshare = securever.System.Share_A_G2(mpk.Mpk)
		securever.IdentityGTbytes = securever.System.IdentityGT.Marshal()
	} else {
		securever.SemiSystem = *shmpc.SystemInit(Partynum)
		securever.Semimpkshare = securever.SemiSystem.Share_A_G2(mpk.Mpk)
		securever.IdentityGTbytes = securever.SemiSystem.IdentityGT.Marshal()
	}
	return securever
}

func (securever *SecureVer) Share_A_Sig(sig Sig, msg []byte, id *big.Int) *Share_Sig {
	share_sig := new(Share_Sig)
	share_sig.S = sig
	hm := H1(msg)
	hid := H2(id)
	hs1 := H3(sig.S1)
	if securever.Security {
		share_sig.HM = securever.System.Share_An_Fp_Offline(hm)
		share_sig.HID = securever.System.Share_An_Fp_Offline(hid)
		share_sig.HS1 = securever.System.Share_An_Fp_Offline(hs1)
	} else {
		share_sig.SemiHM = securever.SemiSystem.Share_An_Fp_Offline(hm)
		share_sig.SemiHID = securever.SemiSystem.Share_An_Fp_Offline(hid)
		share_sig.SemiHS1 = securever.SemiSystem.Share_An_Fp_Offline(hs1)
	}
	return share_sig
}

func (securever *SecureVer) SecVer(sigshares *Share_Sig) (bool, bool) {
	if securever.Security {
		g2hidshares := securever.System.EXP_P_G2_1(curve.Gen2, sigshares.HID)
		mpkg2hidshares := securever.System.SecAdd_G2(*g2hidshares, *securever.mpkshare)
		wshares := securever.System.Pair_P_2(sigshares.S.S2, mpkg2hidshares)
		hidaddhs1shares := securever.System.SecAdd(*sigshares.HM, *sigshares.HS1)
		hshares := securever.System.EXP_P_GT_1(securever.mpk.G, hidaddhs1shares)
		haddwshares := securever.System.SecAdd_GT(*wshares, *hshares)
		resshares := securever.System.SecSubPlaintext_GT(*haddwshares, sigshares.S.S1)
		res, chk := securever.System.OpenGT(*resshares)
		return bytes.Equal(res.Marshal(), securever.IdentityGTbytes), chk
	}
	g2hidshares := securever.SemiSystem.EXP_P_G2_1(curve.Gen2, sigshares.SemiHID)
	mpkg2hidshares := securever.SemiSystem.SecAdd_G2(*g2hidshares, *securever.Semimpkshare)
	wshares := securever.SemiSystem.Pair_P_2(sigshares.S.S2, mpkg2hidshares)
	hidaddhs1shares := securever.SemiSystem.SecAdd(*sigshares.SemiHM, *sigshares.SemiHS1)
	hshares := securever.SemiSystem.EXP_P_GT_1(securever.mpk.G, hidaddhs1shares)
	haddwshares := securever.SemiSystem.SecAdd_GT(*wshares, *hshares)
	resshares := securever.SemiSystem.SecSubPlaintext_GT(*haddwshares, sigshares.S.S1)
	res := securever.SemiSystem.OpenGT(*resshares)
	return bytes.Equal(res.Marshal(), securever.IdentityGTbytes), true
}

func (securever *SecureVer) SecVerWithoutOpen(sigshares *Share_Sig) *[]mpc.Share_GT {
	g2hidshares := securever.System.EXP_P_G2_1(curve.Gen2, sigshares.HID)
	mpkg2hidshares := securever.System.SecAdd_G2(*g2hidshares, *securever.mpkshare)
	wshares := securever.System.Pair_P_2(sigshares.S.S2, mpkg2hidshares)
	hidaddhs1shares := securever.System.SecAdd(*sigshares.HM, *sigshares.HS1)
	hshares := securever.System.EXP_P_GT_1(securever.mpk.G, hidaddhs1shares)
	haddwshares := securever.System.SecAdd_GT(*wshares, *hshares)
	resshares := securever.System.SecSubPlaintext_GT(*haddwshares, sigshares.S.S1)
	return resshares
}

func (securever *SecureVer) SemiSecVerWithoutOpen(sigshares *Share_Sig) *[]shmpc.Share_GT {
	g2hidshares := securever.SemiSystem.EXP_P_G2_1(curve.Gen2, sigshares.SemiHID)
	mpkg2hidshares := securever.SemiSystem.SecAdd_G2(*g2hidshares, *securever.Semimpkshare)
	wshares := securever.SemiSystem.Pair_P_2(sigshares.S.S2, mpkg2hidshares)
	hidaddhs1shares := securever.SemiSystem.SecAdd(*sigshares.SemiHM, *sigshares.SemiHS1)
	hshares := securever.SemiSystem.EXP_P_GT_1(securever.mpk.G, hidaddhs1shares)
	haddwshares := securever.SemiSystem.SecAdd_GT(*wshares, *hshares)
	resshares := securever.SemiSystem.SecSubPlaintext_GT(*haddwshares, sigshares.S.S1)
	return resshares
}
