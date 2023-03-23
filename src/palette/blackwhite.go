package palette

import (
	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

func GetBWColors() []Thread {
	bwColors := []Thread{
		{ID: 1, StringID: "1", Name: "White", RGB: colorConverter.SRGB{R: 255, G: 255, B: 255}, LAB: colorConverter.CIELab{}},
		{ID: 26, StringID: "26", Name: "Black", RGB: colorConverter.SRGB{R: 0, G: 0, B: 0}, LAB: colorConverter.CIELab{}},
	}

	for i, c := range bwColors {
		bwColors[i].LAB = colorConverter.SRGBToCIELab(c.RGB)
	}

	return bwColors
}
