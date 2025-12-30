package palette

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

type Color struct {
	ID       int
	StringID string                `csv:"id"`
	Name     string                `csv:"name"`
	R        uint8                 `csv:"r"`
	G        uint8                 `csv:"g"`
	B        uint8                 `csv:"b"`
	RGB      colorConverter.SRGB   `csv:"-"`
	LAB      colorConverter.CIELab `csv:"-"`
}

func ReadCSV(filename string) ([]Color, error) {
	dmcColors := []Color{}

	path := fmt.Sprintf("palette/%s.csv", filename)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := gocsv.UnmarshalFile(file, &dmcColors); err != nil {
		return nil, err
	}

	maxID := 10000
	for i, c := range dmcColors {
		dmcColors[i].RGB = colorConverter.SRGB{R: c.R, G: c.G, B: c.B}
		dmcColors[i].LAB = colorConverter.SRGBToCIELab(dmcColors[i].RGB)

		dmcColors[i].ID, err = strconv.Atoi(dmcColors[i].StringID)
		if err != nil {
			dmcColors[i].ID = maxID
			maxID += 1
		}
	}

	return dmcColors, nil
}

func ConvertOriginal(colors []colorConverter.SRGB) []Color {
  var legend []Color
	for _, c := range(colors) {
		col := Color{
			R: c.R,
			G: c.G,
			B: c.B,
			RGB: c,
		}

		legend = append(legend, col)
  }

  return legend
}
