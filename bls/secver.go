package bls

import (
	"bytes"

	"github.com/Oryx/curve"
	"github.com/Oryx/mpc"
	"github.com/Oryx/shmpc"
)

type SecureVer struct {
	System     mpc.ShareSystem
	SemiSystem shmpc.ShareSystem
	Security   bool
}

type Share_Sig struct {
	Sig         *Sig
	Hmshare     *[]mpc.Share_G1
	Pkshare     *[]mpc.Share_G2
	SemiHmshare *[]shmpc.Share_G1
	SemiPkshare *[]shmpc.Share_G2
}

func SecureVerInit(Partynum int, ismalicious bool) *SecureVer {
	securever := new(SecureVer)
	securever.Security = ismalicious
	if ismalicious {
		securever.System = *mpc.SystemInit(Partynum)
	} else {
		securever.SemiSystem = *shmpc.SystemInit(Partynum)
	}
	return securever
}

func (securever *SecureVer) Share_A_Sig(sig *Sig, HM *curve.G1, pk *PublicKey) *Share_Sig {
	share_sig := new(Share_Sig)
	share_sig.Sig = sig
	if securever.Security {
		share_sig.Hmshare = securever.System.Share_A_G1_Offline(HM)
		share_sig.Pkshare = securever.System.Share_A_G2_Offline(pk.PK)
	} else {
		share_sig.SemiHmshare = securever.SemiSystem.Share_A_G1_Offline(HM)
		share_sig.SemiPkshare = securever.SemiSystem.Share_A_G2_Offline(pk.PK)
	}
	return share_sig
}

func (securever *SecureVer) SecVer(sigshares *Share_Sig) (bool, bool) {
	left := curve.Pair(sigshares.Sig.S, curve.Gen2)
	if securever.Security {
		rightshares := securever.System.Pair_S(sigshares.Hmshare, sigshares.Pkshare)
		Q := securever.System.SecSubPlaintext_GT(*rightshares, left)
		res, check := securever.System.OpenGT(*Q)
		return bytes.Equal(res.Marshal(), securever.System.IdentityGTBytes), check
	}
	rightshares := securever.SemiSystem.Pair_S(sigshares.SemiHmshare, sigshares.SemiPkshare)
	Q := securever.SemiSystem.SecSubPlaintext_GT(*rightshares, left)
	res := securever.SemiSystem.OpenGT(*Q)
	return bytes.Equal(res.Marshal(), securever.System.IdentityGTBytes), true
}

func (securever *SecureVer) SecVerWithoutOpen(sigshares *Share_Sig) *[]mpc.Share_GT {
	left := curve.Pair(sigshares.Sig.S, curve.Gen2)
	rightshares := securever.System.Pair_S(sigshares.Hmshare, sigshares.Pkshare)
	Q := securever.System.SecSubPlaintext_GT(*rightshares, left)
	return Q
}

func (securever *SecureVer) SemiSecVerWithoutOpen(sigshares *Share_Sig) *[]shmpc.Share_GT {
	left := curve.Pair(sigshares.Sig.S, curve.Gen2)
	rightshares := securever.SemiSystem.Pair_S(sigshares.SemiHmshare, sigshares.SemiPkshare)
	Q := securever.SemiSystem.SecSubPlaintext_GT(*rightshares, left)
	return Q
}
