package colorConverter

import (
	"math"
)

type CIELab struct {
	L float64
	A float64
	B float64
}

type SRGB struct {
	R uint8
	G uint8
	B uint8
}

type lrgb struct {
	r float64
	g float64
	b float64
}

type xyz struct {
	x float64
	y float64
	z float64
}

func SRGBToCIELab(s SRGB) CIELab {
	l := sRGBTolRGB(s)
	x := lRGBToXYZ(l)
	c := xyzToLab(x)

	return c
}

func sRGBTolRGB(s SRGB) lrgb {
	a := 0.055

	l := lrgb{}

	r, g, b := float64(s.R)/255.0, float64(s.G)/255.0, float64(s.B)/255.0

	if r <= 0.04045 {
		l.r = r/12.92
	} else {
		l.r = math.Pow((r + a)/(1.0 + a), 2.4)
	}

	if g <= 0.04045 {
    l.g = g/12.92
  } else {
    l.g = math.Pow((g + a)/(1.0 + a), 2.4)
  }

	if b <= 0.04045 {
    l.b = b/12.92
  } else {
    l.b = math.Pow((b + a)/(1.0 + a), 2.4)
  }

	return l
}

func lRGBToXYZ(l lrgb) xyz {
	c := xyz{}

	c.x = (0.4124*l.r) + (0.3576*l.g) + (0.1805*l.b)
	c.y = (0.2126*l.r) + (0.7152*l.g) + (0.0722*l.b)
	c.z = (0.0193*l.r) + (0.1192*l.g) + (0.9505*l.b)

	return c
}

func xyzToLab(c xyz) CIELab {
	d := 6.0/29.0
	d2 := math.Pow(d, 2)
	d3 := math.Pow(d, 3)

	xn := 0.95047
	yn := 1.00000
	zn := 1.08883

	xr := c.x / xn
	yr := c.y / yn
	zr := c.z / zn

	var fx, fy, fz float64

	if xr > d3 {
		fx = math.Pow(xr, 1.0/3.0)
	} else {
		fx = xr/(3*d2) + 4.0/29.0
	}

	if yr > d3 {
    fy = math.Pow(yr, 1.0/3.0)
  } else {
    fy = yr/(3*d2) + 4.0/29.0
  }

	if zr > d3 {
    fz = math.Pow(zr, 1.0/3.0)
  } else {
    fz = zr/(3*d2) + 4.0/29.0
  }
	cie := CIELab{}

	cie.L = 116.0*fy - 16.0
	cie.A = 500.0*(fx - fy)
	cie.B = 200.0*(fy - fz)

	return cie
}
