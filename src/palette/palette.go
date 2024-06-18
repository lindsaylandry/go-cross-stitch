package palette

import (
	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

type Thread struct {
	ID       int 
	StringID string `csv:"id"`
	Name     string `csv:"name"`
	R        uint8  `csv:"r"`
  G        uint8  `csv:"g"`
	B        uint8  `csv:"b"`
	RGB      colorConverter.SRGB `csv:"-"`
	LAB      colorConverter.CIELab `csv:"-"`
}

func GreyPalette() ([]Thread, error) {
	return palette("../../palette/black-white-grey.csv")
}

func DMCPalette() ([]Thread, error) {
	return GetDMCColors()
}

func palette(path string) ([]Thread, error) {
	return GetDMCColors()
}
