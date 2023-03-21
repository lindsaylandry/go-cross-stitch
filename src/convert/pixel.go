package convert

import (
	"fmt"
	"image/color"

	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

func (c *Converter) pixel() []colorConverter.SRGB {
	bounds := c.image.Bounds()

	colors := make(map[color.Color]int)
	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		for x := bounds.Min.X; x < bounds.Dx(); x++ {
			pixel := c.image.At(x, y)
			if _, ok := colors[pixel]; !ok {
				colors[pixel] = 1
			} else {
				colors[pixel] += 1
			}
		}
	}

	bestColors := []colorConverter.SRGB{}
	for k, v := range colors {
		fmt.Printf("%+v:%d\n", k, v)
		r, g, b, _ := k.RGBA()
		bestColors = append(bestColors, colorConverter.SRGB{uint8(r), uint8(g), uint8(b)})
	}
	fmt.Printf("Number Colors: %d\n", len(bestColors))

	return bestColors
}
