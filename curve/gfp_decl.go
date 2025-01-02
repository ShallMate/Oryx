//go:build (amd64 && !generic) || (arm64 && !generic)
// +build amd64,!generic arm64,!generic

package curve

// This file contains forward declarations for the architecture-specific
// assembly implementations of these functions, provided that they exist.

var hasBMI2 = true

// go:noescape
func gfpNeg(c, a *gfP)

//go:noescape
func gfpAdd(c, a, b *gfP)

//go:noescape
func gfpSub(c, a, b *gfP)

//go:noescape
func gfpMul(c, a, b *gfP)
