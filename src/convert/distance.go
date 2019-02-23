package convert

import (
	"math"
	//"fmt"
)

func labDistance76(l1, a1, b1, l2, a2, b2 float64) float64 {
	return square(l2-l1) + square(a2-a1) + square(b2-b1)
}

func labDistance94(l1, a1, b1, l2, a2, b2 float64) float64 {
	dL := l1 - l2
	dA := a1 - a2
	dB := b1 - b2

	c1 := math.Sqrt(a1*a1 + b1*b1)
	c2 := math.Sqrt(a2*a2 + b2*b2)
	dCab := c1 - c2

	dHab := math.Sqrt(square(dA) + square(dB) - square(dCab))

	// constants for textiles
	kL := 2.0
	k1 := 0.048
	k2 := 0.014

	kH := 1.0
	kC := 1.0

	sL := 1.0
	sC := 1.0 + k1*c1
	sH := 1.0 + k2*c1

	e := math.Sqrt(square(dL/(kL*sL)) + square(dCab/(kC*sC)) + square(dHab/(kH*sH)))
	return e
}

func labDistance(l1, a1, b1, l2, a2, b2 float64) float64 {
	//fmt.Printf("%f %f %f\n", l2, a2, b2)
	//l1, a1, b1 = l1*100.0, a1*100.0, b1*100.0
	//l2, a2, b2 = l2*100.0, a2*100.0, b2*100.0

	kL := 1.0
	kC := 1.0
	kH := 1.0

	dLp := l2 - l1
	lL := (l1 + l2) / 2.0

	c1 := math.Sqrt(square(a1) + square(b1))
	c2 := math.Sqrt(square(a2) + square(b2))

	lC := (c1 + c2) / 2.0

	cc := math.Sqrt(math.Pow(lC, 7.0) / (math.Pow(lC, 7.0) + math.Pow(25.0, 7.0)))
	g := (1.0 - cc) / 2.0

	a1p := a1 * (1.0 + g)
	a2p := a2 * (1.0 + g)

	c1p := math.Sqrt(square(a1p) + square(b1))
	c2p := math.Sqrt(square(a2p) + square(b2))

	dCp := c2p - c1p
	lCp := (c1p + c2p) / 2.0

	var h1p float64
	if b1 != a1p || a1p != 0 {
		h1p = math.Atan2(b1, a1p)
		if h1p < 0 {
			h1p += 2.0 * math.Pi
		}
		h1p *= 180.0 / math.Pi
	}

	var h2p float64
	if b2 != a2p || a2p != 0 {
		h2p = math.Atan2(b2, a2p)
		if h2p < 0 {
			h1p += 2.0 * math.Pi
		}
		h2p *= 180.0 / math.Pi
	}

	var dhp float64
	if c1p*c2p != 0 {
		dhp = h2p - h1p
		if dhp > 180.0 {
			dhp -= 360.0
		} else if dhp < 180.0 {
			dhp += 360.0
		}
	}

	dHp := 2.0 * math.Sqrt(c1p*c2p) * math.Sin(dhp/2.0*math.Pi/180.0)

	var lHp float64
	if math.Abs(h1p-h2p) <= 180.0 {
		lHp = (h1p + h2p) / 2.0
	} else if h1p+h2p < 360.0 {
		lHp = (h1p + h2p + 360.0) / 2.0
	} else {
		lHp = (h1p + h2p - 360.0) / 2.0
	}

	t := 1.0 - 0.17*math.Cos((lHp-30.0)*math.Pi/180.0) + 0.24*math.Cos(2.0*lHp*math.Pi/180.0) + 0.32*math.Cos((3.0*lHp+6.0)*math.Pi/180.0) - 0.20*math.Cos((4.0*lHp-63.0)*math.Pi/180.0)

	dT := 30.0 * math.Exp(-square((lHp-275.0)/25.0))

	sL := 1.0 + (0.015 * square(lL-50.0) / math.Sqrt(20.0+square(lL-50.0)))
	sC := 1.0 + 0.045*lCp
	sH := 1.0 + 0.015*lCp*t

	rT := -2.0 * cc * math.Sin(2.0*dT*math.Pi/180.0)

	e := math.Sqrt(square(dLp/(kL*sL)) + square(dCp/(kC*sC)) + square(dHp/(kH*sH)) + rT*(dCp/(kC*sC))*(dHp/(kH*sH)))

	return e
}

func rgbDistance(r1, g1, b1, r2, g2, b2 float64) float64 {
	return math.Sqrt(2*square(r2-r1) + 4*square(g2-g1) + 3*square(b2-b1) + (r2+r1)/2*(square(r2-r1)-square(b2-b1))/256)
}

func square(a float64) float64 {
	return a * a
}
