package palette

import (
	"os"

	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
	"github.com/gocarina/gocsv"
)

func GetDMCColors() ([]Thread, error) {
	// TODO: read CSV file
	dmcColors := []Thread{}

	file, err := os.Open("palette/dmc.csv")
	if err != nil {
    return nil, err
	}
	defer file.Close()

	if err := gocsv.UnmarshalFile(file, &dmcColors); err != nil {
		return nil, err
	}	

	for i, c := range dmcColors {
		c.RGB = colorConverter.SRGB{R: c.R, G: c.G, B: c.B}
		dmcColors[i].LAB = colorConverter.SRGBToCIELab(c.RGB)
	}

	return dmcColors, nil
}
