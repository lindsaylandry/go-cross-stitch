package palette

import (
	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

type Thread struct {
	ID   int
	Name string
	RGB  colorConverter.SRGB
	LAB  colorConverter.CIELab
}

func GreyPalette() ([]Thread, error) {
	return palette("../../palette/black-white-grey.csv")
}

func DMCPalette() []Thread {
	return GetDMCColors()
}

func palette(path string) ([]Thread, error) {
	return GetDMCColors(), nil
}
