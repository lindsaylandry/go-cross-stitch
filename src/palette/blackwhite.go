package palette

import (
	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

func GetBWColors() []Thread {
	bwColors := []Thread{
		{1, "White", colorConverter.SRGB{255, 255, 255}, colorConverter.CIELab{}},
		{26, "Black", colorConverter.SRGB{0, 0, 0}, colorConverter.CIELab{}},
	}

	for i, c := range bwColors {
		bwColors[i].LAB = colorConverter.SRGBToCIELab(c.RGB)
	}

	return bwColors
}
