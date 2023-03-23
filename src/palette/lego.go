package palette

import (
	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

func GetLEGOColors() []Thread {
	legoColors := []Thread{
		{ID: 1, StringID: "1", Name: "White", RGB: colorConverter.SRGB{R: 255, G: 255, B: 255}, LAB: colorConverter.CIELab{}},
		{ID: 21, StringID: "21", Name: "Bright Red", RGB: colorConverter.SRGB{R: 255, G: 0, B: 0}, LAB: colorConverter.CIELab{}},
		{ID: 23, StringID: "23", Name: "Bright Blue", RGB: colorConverter.SRGB{R: 0, G: 0, B: 255}, LAB: colorConverter.CIELab{}},
		{ID: 24, StringID: "24", Name: "Bright Yellow", RGB: colorConverter.SRGB{R: 255, G: 255, B: 0}, LAB: colorConverter.CIELab{}},
		{ID: 26, StringID: "26", Name: "Black", RGB: colorConverter.SRGB{R: 0, G: 0, B: 0}, LAB: colorConverter.CIELab{}},
	}

	for i, c := range legoColors {
		legoColors[i].LAB = colorConverter.SRGBToCIELab(c.RGB)
	}

	return legoColors
}
