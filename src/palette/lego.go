package palette

import (
	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

func GetLEGOColors() []Thread {
	legoColors := []Thread{
		{1, "White", colorConverter.SRGB{255, 255, 255}, colorConverter.CIELab{}},
		{21, "Bright Red", colorConverter.SRGB{255, 0, 0}, colorConverter.CIELab{}},
		{23, "Bright Blue", colorConverter.SRGB{0, 0, 255}, colorConverter.CIELab{}},
		{24, "Bright Yellow", colorConverter.SRGB{255, 255, 0}, colorConverter.CIELab{}},
		{26, "Black", colorConverter.SRGB{0, 0, 0}, colorConverter.CIELab{}},
	}

	for i, c := range legoColors {
		legoColors[i].LAB = colorConverter.SRGBToCIELab(c.RGB)
	}

	return legoColors
}
