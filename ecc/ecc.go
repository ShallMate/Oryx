package ecc

import (
	"crypto/elliptic"
	"math/big"
	"sync"
)

var (
	fieldOne = new(FieldVal).SetInt(1)
)

type KoblitzCurve struct {
	*elliptic.CurveParams

	q *big.Int

	H         int      // cofactor of the curve.
	halfOrder *big.Int // half the order N

	fieldB *FieldVal

	byteSize int

	bytePoints *[32][256][3]FieldVal

	lambda *big.Int

	beta *FieldVal

	// derived.
	a1 *big.Int
	b1 *big.Int
	a2 *big.Int
	b2 *big.Int
}

func (curve *KoblitzCurve) Params() *elliptic.CurveParams {
	return curve.CurveParams
}

func (curve *KoblitzCurve) bigAffineToField(x, y *big.Int) (*FieldVal, *FieldVal) {
	x3, y3 := new(FieldVal), new(FieldVal)
	x3.SetByteSlice(x.Bytes())
	y3.SetByteSlice(y.Bytes())
	return x3, y3
}

func (curve *KoblitzCurve) BigAffineToField(x, y *big.Int) (*FieldVal, *FieldVal, *FieldVal) {
	x3, y3, z3 := new(FieldVal), new(FieldVal), new(FieldVal)
	x3.SetByteSlice(x.Bytes())
	y3.SetByteSlice(y.Bytes())
	z3.SetInt(1)
	return x3, y3, z3
}

func (curve *KoblitzCurve) fieldJacobianToBigAffine(x, y, z *FieldVal) (*big.Int, *big.Int) {
	var zInv, tempZ FieldVal
	zInv.Set(z).Inverse()   // zInv = Z^-1
	tempZ.SquareVal(&zInv)  // tempZ = Z^-2
	x.Mul(&tempZ)           // X = X/Z^2 (mag: 1)
	y.Mul(tempZ.Mul(&zInv)) // Y = Y/Z^3 (mag: 1)
	z.SetInt(1)             // Z = 1 (mag: 1)

	// Normalize the x and y values.
	x.Normalize()
	y.Normalize()

	// Convert the field values for the now affine point to big.Ints.
	x3, y3 := new(big.Int), new(big.Int)
	x3.SetBytes(x.Bytes()[:])
	y3.SetBytes(y.Bytes()[:])
	return x3, y3
}

func (curve *KoblitzCurve) FieldJacobianToBigAffine(x, y, z *FieldVal) (*big.Int, *big.Int) {
	var zInv, tempZ FieldVal
	zInv.Set(z).Inverse()   // zInv = Z^-1
	tempZ.SquareVal(&zInv)  // tempZ = Z^-2
	x.Mul(&tempZ)           // X = X/Z^2 (mag: 1)
	y.Mul(tempZ.Mul(&zInv)) // Y = Y/Z^3 (mag: 1)
	z.SetInt(1)             // Z = 1 (mag: 1)

	// Normalize the x and y values.
	x.Normalize()
	y.Normalize()

	// Convert the field values for the now affine point to big.Ints.
	x3, y3 := new(big.Int), new(big.Int)
	x3.SetBytes(x.Bytes()[:])
	y3.SetBytes(y.Bytes()[:])
	return x3, y3
}

func (curve *KoblitzCurve) IsOnCurve(x, y *big.Int) bool {
	fx, fy := curve.bigAffineToField(x, y)
	y2 := new(FieldVal).SquareVal(fy).Normalize()
	result := new(FieldVal).SquareVal(fx).Mul(fx).AddInt(7).Normalize()
	return y2.Equals(result)
}

func (curve *KoblitzCurve) Getz3(x1, x2 *FieldVal) *FieldVal {
	var z3 FieldVal
	z3.Set(x1).Negate(1).Add(x2) // H = X2-X1 (mag: 3)
	return &z3
}

func (curve *KoblitzCurve) Getx3y3(x1, y1, x2, y2, x3, y3 *FieldVal) {
	var h, i, j, r, v FieldVal
	var negJ, neg2V, negX3 FieldVal
	h.Set(x1).Negate(1).Add(x2)                // H = X2-X1 (mag: 3)
	i.SquareVal(&h).MulInt(4)                  // I = 4*H^2 (mag: 4)
	j.Mul2(&h, &i)                             // J = H*I (mag: 1)
	r.Set(y1).Negate(1).Add(y2).MulInt(2)      // r = 2*(Y2-Y1) (mag: 6)
	v.Mul2(x1, &i)                             // V = X1*I (mag: 1)
	negJ.Set(&j).Negate(1)                     // negJ = -J (mag: 2)
	neg2V.Set(&v).MulInt(2).Negate(2)          // neg2V = -(2*V) (mag: 3)
	x3.Set(&r).Square().Add(&negJ).Add(&neg2V) // X3 = r^2-J-2*V (mag: 6)
	negX3.Set(x3).Negate(6)                    // negX3 = -X3 (mag: 7)
	j.Mul(y1).MulInt(2).Negate(2)              // J = -(2*Y1*J) (mag: 3)
	y3.Set(&v).Add(&negX3).Mul(&r).Add(&j)     // Y3 = r*(V-X3)-2*Y1*J (mag: 4)
	// Normalize the resulting field values to a magnitude of 1 as needed.
	x3.Normalize()
	y3.Normalize()
}

func (curve *KoblitzCurve) Getx3(x1, y1, x2, y2, zinv *FieldVal) *big.Int {
	var h, i, r, v, x3 FieldVal
	r.Set(y1).Negate(1).Add(y2) // r = (y2-y1)
	h.Mul2(&r, zinv)            // h = (y2-y1)*(x2-x1)^-1
	i.SquareVal(&h)             // I = H^2 (mag: 4)
	v.Set(x1).Add(x2).Negate(1) // V = -(X1+X2) (mag: 1)
	x3.Set(&i).Add(&v)          // X3 = r^2-v (mag: 6)
	x3.Normalize()
	affx3 := new(big.Int)
	affx3.SetBytes(x3.Bytes()[:])
	return affx3
}

func (curve *KoblitzCurve) NewGetx3(x1, y1, x2, y2, zinv, p *FieldVal) (*big.Int, *big.Int) {
	var h, i, r, v, x3, r1, invy1, h1, i1, invx3 FieldVal
	invy1.Set(y1).Negate(1) //-y1
	r.Set(&invy1).Add(y2)   // r = (y2-y1)
	r1.Set(y1).Add(y2)      //y2-y1
	h.Mul2(&r, zinv)        // h = (y2-y1)*(x2-x1)^-1
	h1.Mul2(&r1, zinv)
	i.SquareVal(&h) // I = H^2 (mag: 4)
	i1.SquareVal(&h1)
	v.Set(x1).Add(x2).Negate(1) // V = -(X1+X2) (mag: 1)
	x3.Set(&i).Add(&v)          // X3 = r^2-v (mag: 6)
	invx3.Set(&i1).Add(&v)
	x3.Normalize()
	invx3.Normalize()
	affx3 := new(big.Int).SetBytes(x3.Bytes()[:])
	invaffx3 := new(big.Int).SetBytes(invx3.Bytes()[:])
	return affx3, invaffx3
}

func (curve *KoblitzCurve) addZ1AndZ2EqualsOne(x1, y1, z1, x2, y2, x3, y3, z3 *FieldVal) {
	x1.Normalize()
	y1.Normalize()
	x2.Normalize()
	y2.Normalize()
	if x1.Equals(x2) {
		if y1.Equals(y2) {
			// Since x1 == x2 and y1 == y2, point doubling must be
			// done, otherwise the addition would end up dividing
			// by zero.
			curve.doubleJacobian(x1, y1, z1, x3, y3, z3)
			return
		}

		// Since x1 == x2 and y1 == -y2, the sum is the point at
		// infinity per the group law.
		x3.SetInt(0)
		y3.SetInt(0)
		z3.SetInt(0)
		return
	}

	var h, i, j, r, v FieldVal
	var negJ, neg2V, negX3 FieldVal
	h.Set(x1).Negate(1).Add(x2)                // H = X2-X1 (mag: 3)
	i.SquareVal(&h).MulInt(4)                  // I = 4*H^2 (mag: 4)
	j.Mul2(&h, &i)                             // J = H*I (mag: 1)
	r.Set(y1).Negate(1).Add(y2).MulInt(2)      // r = 2*(Y2-Y1) (mag: 6)
	v.Mul2(x1, &i)                             // V = X1*I (mag: 1)
	negJ.Set(&j).Negate(1)                     // negJ = -J (mag: 2)
	neg2V.Set(&v).MulInt(2).Negate(2)          // neg2V = -(2*V) (mag: 3)
	x3.Set(&r).Square().Add(&negJ).Add(&neg2V) // X3 = r^2-J-2*V (mag: 6)
	negX3.Set(x3).Negate(6)                    // negX3 = -X3 (mag: 7)
	j.Mul(y1).MulInt(2).Negate(2)              // J = -(2*Y1*J) (mag: 3)
	y3.Set(&v).Add(&negX3).Mul(&r).Add(&j)     // Y3 = r*(V-X3)-2*Y1*J (mag: 4)
	z3.Set(&h).MulInt(2)                       // Z3 = 2*H (mag: 6)

	// Normalize the resulting field values to a magnitude of 1 as needed.
	x3.Normalize()
	y3.Normalize()
	z3.Normalize()
}

func (curve *KoblitzCurve) addZ1EqualsZ2(x1, y1, z1, x2, y2, x3, y3, z3 *FieldVal) {
	x1.Normalize()
	y1.Normalize()
	x2.Normalize()
	y2.Normalize()
	if x1.Equals(x2) {
		if y1.Equals(y2) {
			curve.doubleJacobian(x1, y1, z1, x3, y3, z3)
			return
		}

		// Since x1 == x2 and y1 == -y2, the sum is the point at
		// infinity per the group law.
		x3.SetInt(0)
		y3.SetInt(0)
		z3.SetInt(0)
		return
	}

	// Calculate X3, Y3, and Z3 according to the intermediate elements
	// breakdown above.
	var a, b, c, d, e, f FieldVal
	var negX1, negY1, negE, negX3 FieldVal
	negX1.Set(x1).Negate(1)                // negX1 = -X1 (mag: 2)
	negY1.Set(y1).Negate(1)                // negY1 = -Y1 (mag: 2)
	a.Set(&negX1).Add(x2)                  // A = X2-X1 (mag: 3)
	b.SquareVal(&a)                        // B = A^2 (mag: 1)
	c.Set(&negY1).Add(y2)                  // C = Y2-Y1 (mag: 3)
	d.SquareVal(&c)                        // D = C^2 (mag: 1)
	e.Mul2(x1, &b)                         // E = X1*B (mag: 1)
	negE.Set(&e).Negate(1)                 // negE = -E (mag: 2)
	f.Mul2(x2, &b)                         // F = X2*B (mag: 1)
	x3.Add2(&e, &f).Negate(3).Add(&d)      // X3 = D-E-F (mag: 5)
	negX3.Set(x3).Negate(5).Normalize()    // negX3 = -X3 (mag: 1)
	y3.Set(y1).Mul(f.Add(&negE)).Negate(3) // Y3 = -(Y1*(F-E)) (mag: 4)
	y3.Add(e.Add(&negX3).Mul(&c))          // Y3 = C*(E-X3)+Y3 (mag: 5)
	z3.Mul2(z1, &a)                        // Z3 = Z1*A (mag: 1)

	// Normalize the resulting field values to a magnitude of 1 as needed.
	x3.Normalize()
	y3.Normalize()
}

func (curve *KoblitzCurve) addZ2EqualsOne(x1, y1, z1, x2, y2, x3, y3, z3 *FieldVal) {
	var z1z1, u2, s2 FieldVal
	x1.Normalize()
	y1.Normalize()
	z1z1.SquareVal(z1)                        // Z1Z1 = Z1^2 (mag: 1)
	u2.Set(x2).Mul(&z1z1).Normalize()         // U2 = X2*Z1Z1 (mag: 1)
	s2.Set(y2).Mul(&z1z1).Mul(z1).Normalize() // S2 = Y2*Z1*Z1Z1 (mag: 1)
	if x1.Equals(&u2) {
		if y1.Equals(&s2) {
			// Since x1 == x2 and y1 == y2, point doubling must be
			// done, otherwise the addition would end up dividing
			// by zero.
			curve.doubleJacobian(x1, y1, z1, x3, y3, z3)
			return
		}

		// Since x1 == x2 and y1 == -y2, the sum is the point at
		// infinity per the group law.
		x3.SetInt(0)
		y3.SetInt(0)
		z3.SetInt(0)
		return
	}

	// Calculate X3, Y3, and Z3 according to the intermediate elements
	// breakdown above.
	var h, hh, i, j, r, rr, v FieldVal
	var negX1, negY1, negX3 FieldVal
	negX1.Set(x1).Negate(1)                // negX1 = -X1 (mag: 2)
	h.Add2(&u2, &negX1)                    // H = U2-X1 (mag: 3)
	hh.SquareVal(&h)                       // HH = H^2 (mag: 1)
	i.Set(&hh).MulInt(4)                   // I = 4 * HH (mag: 4)
	j.Mul2(&h, &i)                         // J = H*I (mag: 1)
	negY1.Set(y1).Negate(1)                // negY1 = -Y1 (mag: 2)
	r.Set(&s2).Add(&negY1).MulInt(2)       // r = 2*(S2-Y1) (mag: 6)
	rr.SquareVal(&r)                       // rr = r^2 (mag: 1)
	v.Mul2(x1, &i)                         // V = X1*I (mag: 1)
	x3.Set(&v).MulInt(2).Add(&j).Negate(3) // X3 = -(J+2*V) (mag: 4)
	x3.Add(&rr)                            // X3 = r^2+X3 (mag: 5)
	negX3.Set(x3).Negate(5)                // negX3 = -X3 (mag: 6)
	y3.Set(y1).Mul(&j).MulInt(2).Negate(2) // Y3 = -(2*Y1*J) (mag: 3)
	y3.Add(v.Add(&negX3).Mul(&r))          // Y3 = r*(V-X3)+Y3 (mag: 4)
	z3.Add2(z1, &h).Square()               // Z3 = (Z1+H)^2 (mag: 1)
	z3.Add(z1z1.Add(&hh).Negate(2))        // Z3 = Z3-(Z1Z1+HH) (mag: 4)

	// Normalize the resulting field values to a magnitude of 1 as needed.
	x3.Normalize()
	y3.Normalize()
	z3.Normalize()
}

func (curve *KoblitzCurve) addGeneric(x1, y1, z1, x2, y2, z2, x3, y3, z3 *FieldVal) {

	var z1z1, z2z2, u1, u2, s1, s2 FieldVal
	z1z1.SquareVal(z1)                        // Z1Z1 = Z1^2 (mag: 1)
	z2z2.SquareVal(z2)                        // Z2Z2 = Z2^2 (mag: 1)
	u1.Set(x1).Mul(&z2z2).Normalize()         // U1 = X1*Z2Z2 (mag: 1)
	u2.Set(x2).Mul(&z1z1).Normalize()         // U2 = X2*Z1Z1 (mag: 1)
	s1.Set(y1).Mul(&z2z2).Mul(z2).Normalize() // S1 = Y1*Z2*Z2Z2 (mag: 1)
	s2.Set(y2).Mul(&z1z1).Mul(z1).Normalize() // S2 = Y2*Z1*Z1Z1 (mag: 1)
	if u1.Equals(&u2) {
		if s1.Equals(&s2) {
			// Since x1 == x2 and y1 == y2, point doubling must be
			// done, otherwise the addition would end up dividing
			// by zero.
			curve.doubleJacobian(x1, y1, z1, x3, y3, z3)
			return
		}

		// Since x1 == x2 and y1 == -y2, the sum is the point at
		// infinity per the group law.
		x3.SetInt(0)
		y3.SetInt(0)
		z3.SetInt(0)
		return
	}

	var h, i, j, r, rr, v FieldVal
	var negU1, negS1, negX3 FieldVal
	negU1.Set(&u1).Negate(1)               // negU1 = -U1 (mag: 2)
	h.Add2(&u2, &negU1)                    // H = U2-U1 (mag: 3)
	i.Set(&h).MulInt(2).Square()           // I = (2*H)^2 (mag: 2)
	j.Mul2(&h, &i)                         // J = H*I (mag: 1)
	negS1.Set(&s1).Negate(1)               // negS1 = -S1 (mag: 2)
	r.Set(&s2).Add(&negS1).MulInt(2)       // r = 2*(S2-S1) (mag: 6)
	rr.SquareVal(&r)                       // rr = r^2 (mag: 1)
	v.Mul2(&u1, &i)                        // V = U1*I (mag: 1)
	x3.Set(&v).MulInt(2).Add(&j).Negate(3) // X3 = -(J+2*V) (mag: 4)
	x3.Add(&rr)                            // X3 = r^2+X3 (mag: 5)
	negX3.Set(x3).Negate(5)                // negX3 = -X3 (mag: 6)
	y3.Mul2(&s1, &j).MulInt(2).Negate(2)   // Y3 = -(2*S1*J) (mag: 3)
	y3.Add(v.Add(&negX3).Mul(&r))          // Y3 = r*(V-X3)+Y3 (mag: 4)
	z3.Add2(z1, z2).Square()               // Z3 = (Z1+Z2)^2 (mag: 1)
	z3.Add(z1z1.Add(&z2z2).Negate(2))      // Z3 = Z3-(Z1Z1+Z2Z2) (mag: 4)
	z3.Mul(&h)                             // Z3 = Z3*H (mag: 1)

	// Normalize the resulting field values to a magnitude of 1 as needed.
	x3.Normalize()
	y3.Normalize()
}

func (curve *KoblitzCurve) AddGeneric(x1, y1, z1, x2, y2, z2, x3, y3, z3 *FieldVal) {
	var z1z1, z2z2, u1, u2, s1, s2 FieldVal
	z1z1.SquareVal(z1)                        // Z1Z1 = Z1^2 (mag: 1)
	z2z2.SquareVal(z2)                        // Z2Z2 = Z2^2 (mag: 1)
	u1.Set(x1).Mul(&z2z2).Normalize()         // U1 = X1*Z2Z2 (mag: 1)
	u2.Set(x2).Mul(&z1z1).Normalize()         // U2 = X2*Z1Z1 (mag: 1)
	s1.Set(y1).Mul(&z2z2).Mul(z2).Normalize() // S1 = Y1*Z2*Z2Z2 (mag: 1)
	s2.Set(y2).Mul(&z1z1).Mul(z1).Normalize() // S2 = Y2*Z1*Z1Z1 (mag: 1)
	if u1.Equals(&u2) {
		if s1.Equals(&s2) {
			curve.doubleJacobian(x1, y1, z1, x3, y3, z3)
			return
		}

		// Since x1 == x2 and y1 == -y2, the sum is the point at
		// infinity per the group law.
		x3.SetInt(0)
		y3.SetInt(0)
		z3.SetInt(0)
		return
	}

	// Calculate X3, Y3, and Z3 according to the intermediate elements
	// breakdown above.
	var h, i, j, r, rr, v FieldVal
	var negU1, negS1, negX3 FieldVal
	negU1.Set(&u1).Negate(1)               // negU1 = -U1 (mag: 2)
	h.Add2(&u2, &negU1)                    // H = U2-U1 (mag: 3)
	i.Set(&h).MulInt(2).Square()           // I = (2*H)^2 (mag: 2)
	j.Mul2(&h, &i)                         // J = H*I (mag: 1)
	negS1.Set(&s1).Negate(1)               // negS1 = -S1 (mag: 2)
	r.Set(&s2).Add(&negS1).MulInt(2)       // r = 2*(S2-S1) (mag: 6)
	rr.SquareVal(&r)                       // rr = r^2 (mag: 1)
	v.Mul2(&u1, &i)                        // V = U1*I (mag: 1)
	x3.Set(&v).MulInt(2).Add(&j).Negate(3) // X3 = -(J+2*V) (mag: 4)
	x3.Add(&rr)                            // X3 = r^2+X3 (mag: 5)
	negX3.Set(x3).Negate(5)                // negX3 = -X3 (mag: 6)
	y3.Mul2(&s1, &j).MulInt(2).Negate(2)   // Y3 = -(2*S1*J) (mag: 3)
	y3.Add(v.Add(&negX3).Mul(&r))          // Y3 = r*(V-X3)+Y3 (mag: 4)
	z3.Add2(z1, z2).Square()               // Z3 = (Z1+Z2)^2 (mag: 1)
	z3.Add(z1z1.Add(&z2z2).Negate(2))      // Z3 = Z3-(Z1Z1+Z2Z2) (mag: 4)
	z3.Mul(&h)                             // Z3 = Z3*H (mag: 1)

	// Normalize the resulting field values to a magnitude of 1 as needed.
	x3.Normalize()
	y3.Normalize()
}

// 雅可比坐标相加
func (curve *KoblitzCurve) addJacobian(x1, y1, z1, x2, y2, z2, x3, y3, z3 *FieldVal) {
	// A point at infinity is the identity according to the group law for
	// elliptic curve cryptography.  Thus, ∞ + P = P and P + ∞ = P.
	if (x1.IsZero() && y1.IsZero()) || z1.IsZero() {
		x3.Set(x2)
		y3.Set(y2)
		z3.Set(z2)
		return
	}
	if (x2.IsZero() && y2.IsZero()) || z2.IsZero() {
		x3.Set(x1)
		y3.Set(y1)
		z3.Set(z1)
		return
	}

	z1.Normalize()
	z2.Normalize()
	isZ1One := z1.Equals(fieldOne)
	isZ2One := z2.Equals(fieldOne)
	switch {
	case isZ1One && isZ2One:
		curve.addZ1AndZ2EqualsOne(x1, y1, z1, x2, y2, x3, y3, z3)
		return
	case z1.Equals(z2):
		curve.addZ1EqualsZ2(x1, y1, z1, x2, y2, x3, y3, z3)
		return
	case isZ2One:
		curve.addZ2EqualsOne(x1, y1, z1, x2, y2, x3, y3, z3)
		return
	}

	// None of the above assumptions are true, so fall back to generic
	// point addition.
	curve.addGeneric(x1, y1, z1, x2, y2, z2, x3, y3, z3)
}

// 两个访射坐标相加，本质是变成雅可比再变回来
func (curve *KoblitzCurve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	// A point at infinity is the identity according to the group law for
	// elliptic curve cryptography.  Thus, ∞ + P = P and P + ∞ = P.
	if x1.Sign() == 0 && y1.Sign() == 0 {
		return x2, y2
	}
	if x2.Sign() == 0 && y2.Sign() == 0 {
		return x1, y1
	}

	// Convert the affine coordinates from big integers to field values
	// and do the point addition in Jacobian projective space.
	fx1, fy1 := curve.bigAffineToField(x1, y1)
	fx2, fy2 := curve.bigAffineToField(x2, y2)
	fx3, fy3, fz3 := new(FieldVal), new(FieldVal), new(FieldVal)
	fOne := new(FieldVal).SetInt(1)
	curve.addJacobian(fx1, fy1, fOne, fx2, fy2, fOne, fx3, fy3, fz3)

	// Convert the Jacobian coordinate field values back to affine big
	// integers.
	return curve.fieldJacobianToBigAffine(fx3, fy3, fz3)
}

func (curve *KoblitzCurve) AddToF(x1, y1, x2, y2 *big.Int) (*FieldVal, *FieldVal) {

	// Convert the affine coordinates from big integers to field values
	// and do the point addition in Jacobian projective space.
	fx1, fy1 := curve.bigAffineToField(x1, y1)
	fx2, fy2 := curve.bigAffineToField(x2, y2)
	fx3, fy3, fz3 := new(FieldVal), new(FieldVal), new(FieldVal)
	fOne := new(FieldVal).SetInt(1)
	curve.addJacobian(fx1, fy1, fOne, fx2, fy2, fOne, fx3, fy3, fz3)
	// Convert the Jacobian coordinate field values back to affine big
	// integers.
	bx3, by3 := curve.fieldJacobianToBigAffine(fx3, fy3, fz3)
	fbx3, fby3 := curve.bigAffineToField(bx3, by3)
	return fbx3, fby3
}

func (curve *KoblitzCurve) Add1(x1, y1, x2, y2 *big.Int) (*FieldVal, *FieldVal, *FieldVal) {
	// A point at infinity is the identity according to the group law for
	// elliptic curve cryptography.  Thus, ∞ + P = P and P + ∞ = P.

	// Convert the affine coordinates from big integers to field values
	// and do the point addition in Jacobian projective space.
	fx1, fy1 := curve.bigAffineToField(x1, y1)
	fx2, fy2 := curve.bigAffineToField(x2, y2)
	fx3, fy3, fz3 := new(FieldVal), new(FieldVal), new(FieldVal)
	fOne := new(FieldVal).SetInt(1)
	curve.addJacobian(fx1, fy1, fOne, fx2, fy2, fOne, fx3, fy3, fz3)
	// Convert the Jacobian coordinate field values back to affine big
	// integers.
	return fx3, fy3, fz3
}

// 雅可比坐标翻倍，当z =1
func (curve *KoblitzCurve) doubleZ1EqualsOne(x1, y1, x3, y3, z3 *FieldVal) {
	var a, b, c, d, e, f FieldVal
	z3.Set(y1).MulInt(2)                     // Z3 = 2*Y1 (mag: 2)
	a.SquareVal(x1)                          // A = X1^2 (mag: 1)
	b.SquareVal(y1)                          // B = Y1^2 (mag: 1)
	c.SquareVal(&b)                          // C = B^2 (mag: 1)
	b.Add(x1).Square()                       // B = (X1+B)^2 (mag: 1)
	d.Set(&a).Add(&c).Negate(2)              // D = -(A+C) (mag: 3)
	d.Add(&b).MulInt(2)                      // D = 2*(B+D)(mag: 8)
	e.Set(&a).MulInt(3)                      // E = 3*A (mag: 3)
	f.SquareVal(&e)                          // F = E^2 (mag: 1)
	x3.Set(&d).MulInt(2).Negate(16)          // X3 = -(2*D) (mag: 17)
	x3.Add(&f)                               // X3 = F+X3 (mag: 18)
	f.Set(x3).Negate(18).Add(&d).Normalize() // F = D-X3 (mag: 1)
	y3.Set(&c).MulInt(8).Negate(8)           // Y3 = -(8*C) (mag: 9)
	y3.Add(f.Mul(&e))                        // Y3 = E*F+Y3 (mag: 10)

	// Normalize the field values back to a magnitude of 1.
	x3.Normalize()
	y3.Normalize()
	z3.Normalize()
}

// 雅可比坐标翻倍
func (curve *KoblitzCurve) doubleGeneric(x1, y1, z1, x3, y3, z3 *FieldVal) {
	var a, b, c, d, e, f FieldVal
	z3.Mul2(y1, z1).MulInt(2)                // Z3 = 2*Y1*Z1 (mag: 2)
	a.SquareVal(x1)                          // A = X1^2 (mag: 1)
	b.SquareVal(y1)                          // B = Y1^2 (mag: 1)
	c.SquareVal(&b)                          // C = B^2 (mag: 1)
	b.Add(x1).Square()                       // B = (X1+B)^2 (mag: 1)
	d.Set(&a).Add(&c).Negate(2)              // D = -(A+C) (mag: 3)
	d.Add(&b).MulInt(2)                      // D = 2*(B+D)(mag: 8)
	e.Set(&a).MulInt(3)                      // E = 3*A (mag: 3)
	f.SquareVal(&e)                          // F = E^2 (mag: 1)
	x3.Set(&d).MulInt(2).Negate(16)          // X3 = -(2*D) (mag: 17)
	x3.Add(&f)                               // X3 = F+X3 (mag: 18)
	f.Set(x3).Negate(18).Add(&d).Normalize() // F = D-X3 (mag: 1)
	y3.Set(&c).MulInt(8).Negate(8)           // Y3 = -(8*C) (mag: 9)
	y3.Add(f.Mul(&e))                        // Y3 = E*F+Y3 (mag: 10)

	// Normalize the field values back to a magnitude of 1.
	x3.Normalize()
	y3.Normalize()
	z3.Normalize()
}

// doubleJacobian doubles the passed Jacobian point (x1, y1, z1) and stores the
// result in (x3, y3, z3).
func (curve *KoblitzCurve) doubleJacobian(x1, y1, z1, x3, y3, z3 *FieldVal) {
	// Doubling a point at infinity is still infinity.
	if y1.IsZero() || z1.IsZero() {
		x3.SetInt(0)
		y3.SetInt(0)
		z3.SetInt(0)
		return
	}

	if z1.Normalize().Equals(fieldOne) {
		curve.doubleZ1EqualsOne(x1, y1, x3, y3, z3)
		return
	}

	// Fall back to generic point doubling which works with arbitrary z
	// values.
	curve.doubleGeneric(x1, y1, z1, x3, y3, z3)
}

// Double returns 2*(x1,y1). Part of the elliptic.Curve interface.
func (curve *KoblitzCurve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	if y1.Sign() == 0 {
		return new(big.Int), new(big.Int)
	}

	// Convert the affine coordinates from big integers to field values
	// and do the point doubling in Jacobian projective space.
	fx1, fy1 := curve.bigAffineToField(x1, y1)
	fx3, fy3, fz3 := new(FieldVal), new(FieldVal), new(FieldVal)
	fOne := new(FieldVal).SetInt(1)
	curve.doubleJacobian(fx1, fy1, fOne, fx3, fy3, fz3)

	// Convert the Jacobian coordinate field values back to affine big
	// integers.
	return curve.fieldJacobianToBigAffine(fx3, fy3, fz3)
}

func (curve *KoblitzCurve) splitK(k []byte) ([]byte, []byte, int, int) {
	// All math here is done with big.Int, which is slow.
	// At some point, it might be useful to write something similar to
	// FieldVal but for N instead of P as the prime field if this ends up
	// being a bottleneck.
	bigIntK := new(big.Int)
	c1, c2 := new(big.Int), new(big.Int)
	tmp1, tmp2 := new(big.Int), new(big.Int)
	k1, k2 := new(big.Int), new(big.Int)

	bigIntK.SetBytes(k)
	// c1 = round(b2 * k / n) from step 4.
	// Rounding isn't really necessary and costs too much, hence skipped
	c1.Mul(curve.b2, bigIntK)
	c1.Div(c1, curve.N)
	// c2 = round(b1 * k / n) from step 4 (sign reversed to optimize one step)
	// Rounding isn't really necessary and costs too much, hence skipped
	c2.Mul(curve.b1, bigIntK)
	c2.Div(c2, curve.N)
	// k1 = k - c1 * a1 - c2 * a2 from step 5 (note c2's sign is reversed)
	tmp1.Mul(c1, curve.a1)
	tmp2.Mul(c2, curve.a2)
	k1.Sub(bigIntK, tmp1)
	k1.Add(k1, tmp2)
	// k2 = - c1 * b1 - c2 * b2 from step 5 (note c2's sign is reversed)
	tmp1.Mul(c1, curve.b1)
	tmp2.Mul(c2, curve.b2)
	k2.Sub(tmp2, tmp1)

	// Note Bytes() throws out the sign of k1 and k2. This matters
	// since k1 and/or k2 can be negative. Hence, we pass that
	// back separately.
	return k1.Bytes(), k2.Bytes(), k1.Sign(), k2.Sign()
}

// moduloReduce reduces k from more than 32 bytes to 32 bytes and under.  This
// is done by doing a simple modulo curve.N.  We can do this since G^N = 1 and
// thus any other valid point on the elliptic curve has the same order.
func (curve *KoblitzCurve) moduloReduce(k []byte) []byte {
	// Since the order of G is curve.N, we can use a much smaller number
	// by doing modulo curve.N
	if len(k) > curve.byteSize {
		// Reduce k by performing modulo curve.N.
		tmpK := new(big.Int).SetBytes(k)
		tmpK.Mod(tmpK, curve.N)
		return tmpK.Bytes()
	}

	return k
}

func NAF(k []byte) ([]byte, []byte) {
	var carry, curIsOne, nextIsOne bool
	// these default to zero
	retPos := make([]byte, len(k)+1)
	retNeg := make([]byte, len(k)+1)
	for i := len(k) - 1; i >= 0; i-- {
		curByte := k[i]
		for j := uint(0); j < 8; j++ {
			curIsOne = curByte&1 == 1
			if j == 7 {
				if i == 0 {
					nextIsOne = false
				} else {
					nextIsOne = k[i-1]&1 == 1
				}
			} else {
				nextIsOne = curByte&2 == 2
			}
			if carry {
				if curIsOne {
					// This bit is 1, so continue to carry
					// and don't need to do anything.
				} else {
					// We've hit a 0 after some number of
					// 1s.
					if nextIsOne {
						// Start carrying again since
						// a new sequence of 1s is
						// starting.
						retNeg[i+1] += 1 << j
					} else {
						// Stop carrying since 1s have
						// stopped.
						carry = false
						retPos[i+1] += 1 << j
					}
				}
			} else if curIsOne {
				if nextIsOne {
					// If this is the start of at least 2
					// consecutive 1s, set the current one
					// to -1 and start carrying.
					retNeg[i+1] += 1 << j
					carry = true
				} else {
					// This is a singleton, not consecutive
					// 1s.
					retPos[i+1] += 1 << j
				}
			}
			curByte >>= 1
		}
	}
	if carry {
		retPos[0] = 1
		return retPos, retNeg
	}
	return retPos[1:], retNeg[1:]
}

// ScalarMult returns k*(Bx, By) where k is a big endian integer.
// Part of the elliptic.Curve interface.
func (curve *KoblitzCurve) ScalarMult(Bx, By *big.Int, k []byte) (*big.Int, *big.Int) {
	//start1 := time.Now()
	// Point Q = ∞ (point at infinity).
	qx, qy, qz := new(FieldVal), new(FieldVal), new(FieldVal)

	// Decompose K into k1 and k2 in order to halve the number of EC ops.
	// See Algorithm 3.74 in [GECC].
	k1, k2, signK1, signK2 := curve.splitK(curve.moduloReduce(k))

	// The main equation here to remember is:
	//   k * P = k1 * P + k2 * ϕ(P)
	//
	// P1 below is P in the equation, P2 below is ϕ(P) in the equation
	p1x, p1y := curve.bigAffineToField(Bx, By)
	p1yNeg := new(FieldVal).NegateVal(p1y, 1)
	p1z := new(FieldVal).SetInt(1)

	// NOTE: ϕ(x,y) = (βx,y).  The Jacobian z coordinate is 1, so this math
	// goes through.
	p2x := new(FieldVal).Mul2(p1x, curve.beta)
	p2y := new(FieldVal).Set(p1y)
	p2yNeg := new(FieldVal).NegateVal(p2y, 1)
	p2z := new(FieldVal).SetInt(1)

	if signK1 == -1 {
		p1y, p1yNeg = p1yNeg, p1y
	}
	if signK2 == -1 {
		p2y, p2yNeg = p2yNeg, p2y
	}

	// NAF versions of k1 and k2 should have a lot more zeros.
	//
	// The Pos version of the bytes contain the +1s and the Neg versions
	// contain the -1s.
	k1PosNAF, k1NegNAF := NAF(k1)
	k2PosNAF, k2NegNAF := NAF(k2)
	k1Len := len(k1PosNAF)
	k2Len := len(k2PosNAF)

	m := k1Len
	if m < k2Len {
		m = k2Len
	}

	var k1BytePos, k1ByteNeg, k2BytePos, k2ByteNeg byte
	for i := 0; i < m; i++ {
		// Since we're going left-to-right, pad the front with 0s.
		if i < m-k1Len {
			k1BytePos = 0
			k1ByteNeg = 0
		} else {
			k1BytePos = k1PosNAF[i-(m-k1Len)]
			k1ByteNeg = k1NegNAF[i-(m-k1Len)]
		}
		if i < m-k2Len {
			k2BytePos = 0
			k2ByteNeg = 0
		} else {
			k2BytePos = k2PosNAF[i-(m-k2Len)]
			k2ByteNeg = k2NegNAF[i-(m-k2Len)]
		}

		for j := 7; j >= 0; j-- {
			// Q = 2 * Q
			curve.doubleJacobian(qx, qy, qz, qx, qy, qz)

			if k1BytePos&0x80 == 0x80 {
				curve.addJacobian(qx, qy, qz, p1x, p1y, p1z,
					qx, qy, qz)
			} else if k1ByteNeg&0x80 == 0x80 {
				curve.addJacobian(qx, qy, qz, p1x, p1yNeg, p1z,
					qx, qy, qz)
			}

			if k2BytePos&0x80 == 0x80 {
				curve.addJacobian(qx, qy, qz, p2x, p2y, p2z,
					qx, qy, qz)
			} else if k2ByteNeg&0x80 == 0x80 {
				curve.addJacobian(qx, qy, qz, p2x, p2yNeg, p2z,
					qx, qy, qz)
			}
			k1BytePos <<= 1
			k1ByteNeg <<= 1
			k2BytePos <<= 1
			k2ByteNeg <<= 1
		}
	}

	// Convert the Jacobian coordinate field values back to affine big.Ints.
	return curve.fieldJacobianToBigAffine(qx, qy, qz)
}

// ScalarMult returns k*(Bx, By) where k is a big endian integer.
// Part of the elliptic.Curve interface.
func (curve *KoblitzCurve) ScalarMultField(Bx, By *big.Int, k []byte) (*FieldVal, *FieldVal, *FieldVal) {
	//start1 := time.Now()
	// Point Q = ∞ (point at infinity).
	qx, qy, qz := new(FieldVal), new(FieldVal), new(FieldVal)

	// Decompose K into k1 and k2 in order to halve the number of EC ops.
	// See Algorithm 3.74 in [GECC].
	k1, k2, signK1, signK2 := curve.splitK(curve.moduloReduce(k))

	// The main equation here to remember is:
	//   k * P = k1 * P + k2 * ϕ(P)
	//
	// P1 below is P in the equation, P2 below is ϕ(P) in the equation
	p1x, p1y := curve.bigAffineToField(Bx, By)
	p1yNeg := new(FieldVal).NegateVal(p1y, 1)
	p1z := new(FieldVal).SetInt(1)

	// NOTE: ϕ(x,y) = (βx,y).  The Jacobian z coordinate is 1, so this math
	// goes through.
	p2x := new(FieldVal).Mul2(p1x, curve.beta)
	p2y := new(FieldVal).Set(p1y)
	p2yNeg := new(FieldVal).NegateVal(p2y, 1)
	p2z := new(FieldVal).SetInt(1)

	if signK1 == -1 {
		p1y, p1yNeg = p1yNeg, p1y
	}
	if signK2 == -1 {
		p2y, p2yNeg = p2yNeg, p2y
	}

	// NAF versions of k1 and k2 should have a lot more zeros.
	//
	// The Pos version of the bytes contain the +1s and the Neg versions
	// contain the -1s.
	k1PosNAF, k1NegNAF := NAF(k1)
	k2PosNAF, k2NegNAF := NAF(k2)
	k1Len := len(k1PosNAF)
	k2Len := len(k2PosNAF)

	m := k1Len
	if m < k2Len {
		m = k2Len
	}

	var k1BytePos, k1ByteNeg, k2BytePos, k2ByteNeg byte
	for i := 0; i < m; i++ {
		// Since we're going left-to-right, pad the front with 0s.
		if i < m-k1Len {
			k1BytePos = 0
			k1ByteNeg = 0
		} else {
			k1BytePos = k1PosNAF[i-(m-k1Len)]
			k1ByteNeg = k1NegNAF[i-(m-k1Len)]
		}
		if i < m-k2Len {
			k2BytePos = 0
			k2ByteNeg = 0
		} else {
			k2BytePos = k2PosNAF[i-(m-k2Len)]
			k2ByteNeg = k2NegNAF[i-(m-k2Len)]
		}

		for j := 7; j >= 0; j-- {
			// Q = 2 * Q
			curve.doubleJacobian(qx, qy, qz, qx, qy, qz)

			if k1BytePos&0x80 == 0x80 {
				curve.addJacobian(qx, qy, qz, p1x, p1y, p1z,
					qx, qy, qz)
			} else if k1ByteNeg&0x80 == 0x80 {
				curve.addJacobian(qx, qy, qz, p1x, p1yNeg, p1z,
					qx, qy, qz)
			}

			if k2BytePos&0x80 == 0x80 {
				curve.addJacobian(qx, qy, qz, p2x, p2y, p2z,
					qx, qy, qz)
			} else if k2ByteNeg&0x80 == 0x80 {
				curve.addJacobian(qx, qy, qz, p2x, p2yNeg, p2z,
					qx, qy, qz)
			}
			k1BytePos <<= 1
			k1ByteNeg <<= 1
			k2BytePos <<= 1
			k2ByteNeg <<= 1
		}
	}

	// Convert the Jacobian coordinate field values back to affine big.Ints.
	return qx, qy, qz
}

func (curve *KoblitzCurve) ScalarBaseMultField(k []byte) (*FieldVal, *FieldVal, *FieldVal) {
	newK := curve.moduloReduce(k)
	diff := len(curve.bytePoints) - len(newK)

	// Point Q = ∞ (point at infinity).
	qx, qy, qz := new(FieldVal), new(FieldVal), new(FieldVal)

	for i, byteVal := range newK {
		p := curve.bytePoints[diff+i][byteVal]
		curve.addJacobian(qx, qy, qz, &p[0], &p[1], &p[2], qx, qy, qz)
	}
	return qx, qy, qz
}

func (curve *KoblitzCurve) QPlus1Div4() *big.Int {
	return curve.q
}

// Q returns the (P+1)/4 constant for the curve for use in calculating square
// roots via exponentiation.
func (curve *KoblitzCurve) Q() *big.Int {
	return curve.q
}

var initonce sync.Once
var secp256k1 KoblitzCurve

func initAll() {
	initS256()
}

func fromHex(s string) *big.Int {
	r, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid hex in source file: " + s)
	}
	return r
}

func initS256() {
	// Curve parameters taken from [SECG] section 2.4.1.
	secp256k1.CurveParams = new(elliptic.CurveParams)
	secp256k1.P = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
	secp256k1.N = fromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")
	secp256k1.B = fromHex("0000000000000000000000000000000000000000000000000000000000000007")
	secp256k1.Gx = fromHex("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798")
	secp256k1.Gy = fromHex("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8")
	secp256k1.BitSize = 256
	secp256k1.q = new(big.Int).Div(new(big.Int).Add(secp256k1.P,
		big.NewInt(1)), big.NewInt(4))
	secp256k1.H = 1
	secp256k1.halfOrder = new(big.Int).Rsh(secp256k1.N, 1)
	secp256k1.fieldB = new(FieldVal).SetByteSlice(secp256k1.B.Bytes())

	// Provided for convenience since this gets computed repeatedly.
	secp256k1.byteSize = secp256k1.BitSize / 8

	if err := loadS256BytePoints(); err != nil {
		panic(err)
	}
	secp256k1.lambda = fromHex("5363AD4CC05C30E0A5261C028812645A122E22EA20816678DF02967C1B23BD72")
	secp256k1.beta = new(FieldVal).SetHex("7AE96A2B657C07106E64479EAC3434E99CF0497512F58995C1396C28719501EE")
	secp256k1.a1 = fromHex("3086D221A7D46BCDE86C90E49284EB15")
	secp256k1.b1 = fromHex("-E4437ED6010E88286F547FA90ABFE4C3")
	secp256k1.a2 = fromHex("114CA50F7A8E2F3F657C1108D9D44CFD8")
	secp256k1.b2 = fromHex("3086D221A7D46BCDE86C90E49284EB15")
}

// S256 returns a Curve which implements secp256k1.
func S256() *KoblitzCurve {
	initonce.Do(initAll)
	return &secp256k1
}
