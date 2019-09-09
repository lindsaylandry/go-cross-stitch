package convert

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"

	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
	"github.com/lindsaylandry/go-cross-stitch/src/palette"
)

type Legend struct {
	Color  palette.Thread
	Count  int
	Symbol int
}

type Converter struct {
	image    image.Image
	newImage struct {
		image   *image.RGBA
		p       int
		count   map[palette.Thread]int
		legend  []Legend
		symbols [][]palette.Symbol
	}
	path      string
	symbols   []palette.Symbol
	limit     int
	rgb       bool
	pc        []palette.Thread
	dither    bool
	greyscale bool
}

func (c *Converter) getImage() error {
	data, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	c.image, _, err = image.Decode(bytes.NewReader(data))
	return err
}

func NewConverter(filename string, num int, rgb, all bool, pal string, dit, gre bool) (*Converter, error) {
	c := Converter{}

	c.newImage.p = 4
	c.path = filename

	if err := c.getImage(); err != nil {
		return &c, err
	}

	bounds := c.image.Bounds()
	c.newImage.image = image.NewRGBA(bounds)

	c.newImage.symbols = make([][]palette.Symbol, bounds.Dy()-bounds.Min.Y)
	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		c.newImage.symbols[y] = make([]palette.Symbol, bounds.Dx()-bounds.Min.X)
		for x := bounds.Min.X; x < bounds.Dx(); x++ {
			pixel := c.image.At(x, y)
			c.newImage.image.Set(x, y, pixel)
		}
	}

	c.symbols = palette.GetSymbols()
	c.limit = num
	c.rgb = rgb
	c.dither = dit
	c.greyscale = gre

	if pal == "lego" {
		c.pc = palette.GetLEGOColors()
	} else if pal == "dmc" {
		c.pc = palette.GetDMCColors()

		if !all {
			//best colors rgb
			bcrgb := c.colorQuant()
			// Convert best-colors to thread palette
			c.pc = c.convertPalette(bcrgb)
		}
	} else {
		c.pc = palette.GetBWColors()
	}

	c.newImage.count = make(map[palette.Thread]int)

	return &c, nil
}

func (c *Converter) Greyscale() {
	bounds := c.newImage.image.Bounds()
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			r32, g32, b32, a := c.newImage.image.At(x, y).RGBA()
			r, g, b := uint8(r32), uint8(g32), uint8(b32)

			gg := colorConverter.Greyscale(r, g, b)
			c.newImage.image.Set(x, y, color.RGBA{gg, gg, gg, uint8(a)})
		}
	}
}

func (c *Converter) Convert() error {
	if err := c.convertImage(); err != nil {
		return err
	}

	err := c.WriteFiles()

	return err
}

// convert best-color palette to match available threads
func (c *Converter) convertPalette(colors []colorConverter.SRGB) []palette.Thread {
	dict := make(map[palette.Thread]int)
	var legend []palette.Thread
	for i := 0; i < len(colors); i++ {
		minLen := math.MaxFloat64
		minIndex := 0

		for x := 0; x < len(c.pc); x++ {
			var dist float64
			if c.rgb {
				dist = rgbDistance(float64(c.pc[x].RGB.R), float64(c.pc[x].RGB.G), float64(c.pc[x].RGB.B), float64(colors[i].R), float64(colors[i].G), float64(colors[i].B))
			} else {
				cie := colorConverter.SRGBToCIELab(colors[i])
				dist = labDistance(c.pc[x].LAB.L, c.pc[x].LAB.A, c.pc[x].LAB.B, cie.L, cie.A, cie.B)
			}
			if dist < minLen {
				minLen = dist
				minIndex = x
			}
		}

		if _, ok := dict[c.pc[minIndex]]; !ok {
			dict[c.pc[minIndex]] = 0
			legend = append(legend, c.pc[minIndex])
		}
	}

	return legend
}

func (c *Converter) convertImage() error {
	if c.dither {
		c.floydSteinbergDither()
	} else {

		bounds := c.image.Bounds()

		n := 4
		countChan := make(chan map[palette.Thread]int, n*n)

		// goroutines
		for m := 0; m < n; m++ {
			ylow := bounds.Min.Y + m*bounds.Dy()/n
			yhigh := (m + 1) * bounds.Dy() / n

			for p := 0; p < n; p++ {
				xlow := bounds.Min.X + p*bounds.Dx()/n
				xhigh := (p + 1) * bounds.Dx() / n

				go func() {
					count := c.convertImageChunk(xlow, xhigh, ylow, yhigh)
					countChan <- count
				}()
			}
		}

		for i := 0; i < n*n; i++ {
			count := <-countChan
			for k, v := range count {
				if _, ok := c.newImage.count[k]; ok {
					c.newImage.count[k] += v
				} else {
					c.newImage.count[k] = v
				}
			}
		}
	}

	for i, v := range c.pc {
		l := Legend{v, c.newImage.count[v], c.symbols[i].Code}
		c.newImage.legend = append(c.newImage.legend, l)
	}
	quickSortLegend(c.newImage.legend)
	return nil
}

func (c *Converter) convertImageChunk(xlow, xhigh, ylow, yhigh int) map[palette.Thread]int {
	count := make(map[palette.Thread]int)
	for y := ylow; y < yhigh; y++ {
		for x := xlow; x < xhigh; x++ {
			minIndex := c.setNewPixel(x, y)

			if _, ok := count[c.pc[minIndex]]; ok {
				count[c.pc[minIndex]] += 1
			} else {
				count[c.pc[minIndex]] = 1
			}
		}
	}
	return count
}

func (c *Converter) setNewPixel(x, y int) int {
	r32, g32, b32, a := c.newImage.image.At(x, y).RGBA()
	r, g, b := uint8(r32), uint8(g32), uint8(b32)

	minLen := math.MaxFloat64
	minIndex := 0
	for i := 0; i < len(c.pc); i++ {
		var dist float64
		if c.rgb || colorConverter.Greyscale(r, g, b) < 100 {
			dist = rgbDistance(float64(r), float64(g), float64(b), float64(c.pc[i].RGB.R), float64(c.pc[i].RGB.G), float64(c.pc[i].RGB.B))
		} else {
			cie := colorConverter.SRGBToCIELab(colorConverter.SRGB{r, g, b})
			dist = labDistance(c.pc[i].LAB.L, c.pc[i].LAB.A, c.pc[i].LAB.B, cie.L, cie.A, cie.B)
		}
		if dist < minLen {
			minLen = dist
			minIndex = i
		}
	}

	c.newImage.symbols[y][x] = c.symbols[minIndex]
	c.newImage.image.Set(x, y, color.RGBA{uint8(c.pc[minIndex].RGB.R), uint8(c.pc[minIndex].RGB.G), uint8(c.pc[minIndex].RGB.B), uint8(a)})

	return minIndex
}

// start with 32 colors
func (c *Converter) colorQuant() []colorConverter.SRGB {
	bounds := c.image.Bounds()

	allcolors := make([]colorConverter.SRGB, (bounds.Dy()-bounds.Min.Y)*(bounds.Dx()-bounds.Min.X))
	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		for x := bounds.Min.X; x < bounds.Dx(); x++ {
			r32, g32, b32, _ := c.newImage.image.At(x, y).RGBA()
			allcolors[y*(bounds.Dx()-bounds.Min.X)+x] = colorConverter.SRGB{uint8(r32), uint8(g32), uint8(b32)}
		}
	}

	// 1 2 4 8 16
	slices := []int{0, len(allcolors) - 1}

	for i := 0; i < c.limit; i++ {
		for j := 0; j < len(slices)-1; j++ {
			// get a slice of allcolors
			s := allcolors[slices[j]:slices[j+1]]

			colorRanges := [][]uint8{{math.MaxUint8, 0}, {math.MaxUint8, 0}, {math.MaxUint8, 0}}
			for c := 0; c < len(s); c++ {
				//R
				if s[c].R < colorRanges[0][0] {
					colorRanges[0][0] = s[c].R
				}
				if s[c].R > colorRanges[0][1] {
					colorRanges[0][1] = s[c].R
				}
				//G
				if s[c].G < colorRanges[1][0] {
					colorRanges[1][0] = s[c].G
				}
				if s[c].G > colorRanges[1][1] {
					colorRanges[1][1] = s[c].G
				}
				//B
				if s[c].B < colorRanges[2][0] {
					colorRanges[2][0] = s[c].B
				}
				if s[c].B > colorRanges[2][1] {
					colorRanges[2][1] = s[c].B
				}
			}
			xr := colorRanges[0][1] - colorRanges[0][0]
			yr := colorRanges[1][1] - colorRanges[1][0]
			zr := colorRanges[2][1] - colorRanges[2][0]

			if xr == 0 && yr == 0 && zr == 0 {
				continue
			}

			// Sort channel that has greatest variance
			if xr > yr && xr > zr {
				quickSortRed(s)
			} else if yr > xr && yr > zr {
				quickSortGreen(s)
			} else {
				quickSortBlue(s)
			}
		}

		// insert 2^n more slice indexes
		max := len(slices) - 1
		for k := 0; k < max; k++ {
			slices = append(slices, (slices[k+1]+slices[k])/2)
		}
		quickSort(slices)
	}

	// Average all sliced colors and insert into bestcolors
	bestColors := make([]colorConverter.SRGB, len(slices)-1)
	for i := 0; i < len(slices)-1; i++ {
		s := allcolors[slices[i]:slices[i+1]]
		var avgR, avgG, avgB float64
		for c := 0; c < len(s); c++ {
			avgR = avgR + math.Pow(float64(s[c].R), 2)
			avgG = avgG + math.Pow(float64(s[c].G), 2)
			avgB = avgB + math.Pow(float64(s[c].B), 2)
		}

		bestColors[i] = colorConverter.SRGB{uint8(math.Sqrt(avgR / float64(len(s)))), uint8(math.Sqrt(avgG / float64(len(s)))), uint8(math.Sqrt(avgB / float64(len(s))))}
	}

	return bestColors
}
