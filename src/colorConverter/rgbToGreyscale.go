package colorConverter

func Greyscale(r, g, b uint8) uint8 {
	return uint8(0.3*float64(r) + 0.59*float64(g) + 0.11*float64(b))
}
