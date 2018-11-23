package convert

import (
	"math"
)

func labDistance76(l1, a1, b1, l2, a2, b2 float64) float64 {
  return math.Pow((l2-l1), 2) + math.Pow((a2-a1), 2) + math.Pow((b2-b1), 2)
}

func labDistance94(l1, a1, b1, l2, a2, b2 float64) float64 {
  dL := l1 - l2
  dA := a1 - a2
  dB := b1 - b2

  c1 := math.Sqrt(math.Pow(a1, 2) + math.Pow(b1, 2))
  c2 := math.Sqrt(math.Pow(a2, 2) + math.Pow(b2, 2))
  dCab := c1 - c2

  dHab := math.Sqrt(math.Pow(dA, 2) + math.Pow(dB, 2) - math.Pow(dCab, 2))

  // constants for textiles
  kL := 2.0
  k1 := 0.048
  k2 := 0.014

  kH := 1.0
  kC := 1.0

  sL := 1.0
  sC := 1.0 + k1*c1
  sH := 1.0 + k2*c1

  e := math.Sqrt(math.Pow(dL/(kL*sL), 2) + math.Pow(dCab/(kC*sC), 2) + math.Pow(dHab/(kH*sH), 2))
  return e
}

func labDistance(l1, a1, b1, l2, a2, b2 float64) float64 {
	kL := 1.0
	kC := 1.0
	kH := 1.0

	dLp := l1 - l2
	lL := (l1 + l2)/2.0

	c1 := math.Sqrt(math.Pow(a1, 2) + math.Pow(b1, 2))
  c2 := math.Sqrt(math.Pow(a2, 2) + math.Pow(b2, 2))

	lC := (c1 + c2)/2.0

	a1p := a1 + (a1/2.0)*(1.0 - math.Sqrt(math.Pow(lC, 7)/(math.Pow(lC, 7) + math.Pow(25.0, 7))))
	a2p := a2 + (a2/2.0)*(1.0 - math.Sqrt(math.Pow(lC, 7)/(math.Pow(lC, 7) + math.Pow(25.0, 7))))

	c1p := math.Sqrt(math.Pow(a1p, 2) + math.Pow(b1, 2))
	c2p := math.Sqrt(math.Pow(a2p, 2) + math.Pow(b2, 2))

	dCp := c2p - c1p
	lCp := (c1p + c2p)/2.0

	// b1, b2, a1p, and a2p are in degrees, convert to radians
	h1p := math.Mod(math.Atan2(math.Pi/180.0*b1, math.Pi/180.0*a1p), 2.0*math.Pi)
	h2p := math.Mod(math.Atan2(math.Pi/180.0*b2, math.Pi/180.0*a2p), 2.0*math.Pi)

	var dhp float64

	if math.Abs(h1p - h2p) <= math.Pi {
		dhp = h2p - h1p
	} else if h2p <= h1p {
		dhp = h2p - h1p + 2.0*math.Pi
	} else {
		dhp = h2p - h1p - 2.0*math.Pi
	}

	dHp := 2.0*math.Sqrt(c1p*c2p)*math.Sin(dhp/2.0)

	var lHp float64
	if math.Abs(h1p - h2p) <= math.Pi {
		lHp = (h1p + h2p)/2.0
	} else if h2p + h1p < 2.0*math.Pi {
		lHp = (h1p + h2p + 2.0*math.Pi)/2.0
  } else {
		lHp = (h1p + h2p - 2.0*math.Pi)/2.0
	}

	t := 1.0 - 0.17*math.Cos(lHp - math.Pi/6.0) + 0.24*math.Cos(2.0*lHp) + 0.32*math.Cos(3.0*lHp + math.Pi/30.0) - 0.20*math.Cos(4.0*lHp - math.Pi*21.0/60.0)

	sL := 1.0 + (0.015*math.Pow(lL - 50.0, 2)/math.Sqrt(20.0 + math.Pow(lL - 50.0, 2)))
	sC := 1.0 + 0.045*lCp
	sH := 1.0 + 0.015*lCp*t

	rT := -2.0*math.Sqrt(math.Pow(lCp, 7)/(math.Pow(lCp, 7) + math.Pow(25.0, 7)))*math.Sin(math.Pi/3.0 * math.Exp(-math.Pow((lHp - math.Pi*55.0/36.0)/(math.Pi*5.0/36.0), 2)))

	e := math.Sqrt(math.Pow(dLp/(kL*sL), 2) + math.Pow(dCp/(kC*sC), 2) + math.Pow(dHp/(kH*sH), 2) + rT*(dCp/(kC*sC))*(dHp/(kH*sH)))

	return e
}

func rgbDistance(r1, g1, b1, r2, g2, b2 float64) float64 {
  dist := 2*math.Pow((r2-r1), 2) + 4*math.Pow((g2-g1), 2) + 3*math.Pow((b2-b1), 2) + (r2+r1)/2*(math.Pow((r2-r1), 2)-math.Pow((b2-b1), 2))/256
  return dist
}
