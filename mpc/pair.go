package mpc

import (
	curve "github.com/Oryx/curve"
)

func (system *ShareSystem) share_Pair_P_1(g1shares Share_G1, g2 *curve.G2) Share_GT {
	shares := new(Share_GT)
	shares.Index = g1shares.Index
	shares.Delta = curve.Pair(g1shares.Delta, g2)
	shares.Gama = curve.Pair(g1shares.Gama, g2)
	shares.Share = curve.Pair(g1shares.Share, g2)
	return *shares
}

func (system *ShareSystem) Pair_P_1(g1shares *[]Share_G1, g2 *curve.G2) *[]Share_GT {
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_Pair_P_1((*g1shares)[i], g2)
	}
	return &shares
}

func (system *ShareSystem) share_Pair_P_2(g1 *curve.G1, g2shares Share_G2) Share_GT {
	shares := new(Share_GT)
	shares.Index = g2shares.Index
	shares.Delta = curve.Pair(g1, g2shares.Delta)
	shares.Gama = curve.Pair(g1, g2shares.Gama)
	shares.Share = curve.Pair(g1, g2shares.Share)
	return *shares
}

func (system *ShareSystem) Pair_P_2(g1 *curve.G1, g2shares *[]Share_G2) *[]Share_GT {
	shares := make([]Share_GT, system.Partynum)
	for i := 0; i < system.Partynum; i++ {
		shares[i] = system.share_Pair_P_2(g1, (*g2shares)[i])
	}
	return &shares
}

func (system *ShareSystem) Pair_S(g1shares *[]Share_G1, g2shares *[]Share_G2) *[]Share_GT {
	sharesA, sharesB, sharesC := system.GenTriplets()
	sharesgA := system.EXP_P_G1_1(curve.Gen1, sharesA)
	sharesgB := system.EXP_P_G2_1(curve.Gen2, sharesB)
	sharesgC := system.EXP_P_G1_1(curve.Gen1, sharesC)
	vshares := system.SecSub_G1(*g1shares, *sharesgA)
	wshares := system.SecSub_G2(*g2shares, *sharesgB)
	v := system.HalfOpenG1(*vshares)
	w := system.HalfOpenG2(*wshares)
	a := system.Pair_P_1(vshares, w)
	b := system.Pair_P_1(sharesgC, curve.Gen2)
	c := system.Pair_P_1(sharesgA, w)
	d := system.Pair_P_2(v, sharesgB)
	shares := system.SecAdd_GT(*a, *b)
	shares = system.SecAdd_GT(*shares, *c)
	shares = system.SecAdd_GT(*shares, *d)
	return shares
}
