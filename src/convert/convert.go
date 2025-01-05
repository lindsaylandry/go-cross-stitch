package convert

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"sort"

	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
	"github.com/lindsaylandry/go-cross-stitch/src/palette"
)

type Legend struct {
	Color  palette.Thread
	Count  int
	Symbol rune
}

type ColorSymbol struct {
	Symbol palette.SymbolRune
	Color  palette.Thread
	Text   string
}

type NewData struct {
	Image   *image.RGBA
	Count   map[palette.Thread]int
	Legend  []Legend
	Symbols [][]ColorSymbol
	Path    string
	Extra   string
	Scheme  string
}

type Converter struct {
	image     image.Image
	newData   NewData
	limit     int
	rgb       bool
	pc        []palette.Thread
	dither    bool
	greyscale bool
	colorgrid bool
}

type Flags struct {
	Num int
	RGB bool
	All bool
	Palette string
	Dither bool
	Greyscale bool
	Pixel bool
	Color bool
	CSV string
}

func NewConverter(filename string, flags Flags) (*Converter, error) {
	c := Converter{}

	c.newData.Path = filename

	if err := c.getImage(); err != nil {
		return &c, err
	}

	bounds := c.image.Bounds()
	c.newData.Image = image.NewRGBA(bounds)

	c.newData.Symbols = make([][]ColorSymbol, bounds.Dy()-bounds.Min.Y)
	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		c.newData.Symbols[y] = make([]ColorSymbol, bounds.Dx()-bounds.Min.X)
	}

	c.limit = flags.Num
	c.rgb = flags.RGB
	c.dither = flags.Dither
	c.greyscale = flags.Greyscale
	c.colorgrid = flags.Color

	if c.rgb {
		c.newData.Extra = "-" + flags.Palette + "-rgb"
	} else {
		c.newData.Extra = "-" + flags.Palette + "-lab"
	}

	csvFile := flags.CSV
	if csvFile == "" {
		csvFile = flags.Palette
	}
		
	pc, err := palette.ReadCSV(csvFile)
	
	if err != nil {
		return &c, err
	}
	c.pc = pc

	if flags.Palette == "lego" {
		c.newData.Scheme = "LEGO"
	} else if flags.Palette == "dmc" || flags.Palette == "anchor" {
		if flags.Palette == "dmc" {
			c.newData.Scheme = "DMC"
		} else {
			c.newData.Scheme = "Anchor"
		}

		if !flags.All {
			if !flags.Pixel {
				//most colors rgb
				c.pc = c.convertPalette(c.pixel())
			} else {
				//best colors rgb
				c.pc = c.convertPalette(c.colorQuant())
			}
		}
	} else if flags.Palette == "bw" {
		c.newData.Scheme = "Black&White"
	} else {
		log.Fatalf("ERROR: -color not recognized")
	}

	c.newData.Count = make(map[palette.Thread]int)

	return &c, nil
}

func (c *Converter) Greyscale() {
	bounds := c.newData.Image.Bounds()
	for x := bounds.Min.X; x < bounds.Dx(); x++ {
		for y := bounds.Min.Y; y < bounds.Dy(); y++ {
			r32, g32, b32, a := c.newData.Image.At(x, y).RGBA()
			r, g, b := uint8(r32), uint8(g32), uint8(b32)

			gg := colorConverter.Greyscale(r, g, b)
			c.newData.Image.Set(x, y, color.RGBA{gg, gg, gg, uint8(a)})
		}
	}
}

func (c *Converter) Convert() (NewData, error) {
	err := c.convertImage()
	return c.newData, err
}

func (c *Converter) getImage() error {
	data, err := ioutil.ReadFile(c.newData.Path)
	if err != nil {
		return err
	}

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)

	c.image, _, err = image.Decode(bytes.NewReader(data))
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
				dist = colorConverter.RGBDistance(float64(c.pc[x].RGB.R), float64(c.pc[x].RGB.G), float64(c.pc[x].RGB.B), float64(colors[i].R), float64(colors[i].G), float64(colors[i].B))
			} else {
				cie := colorConverter.SRGBToCIELab(colors[i])
				dist = colorConverter.LABDistance(c.pc[x].LAB.L, c.pc[x].LAB.A, c.pc[x].LAB.B, cie.L, cie.A, cie.B)
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
				c.newData.Count[k] += v
			}
		}
	}

	symbols := palette.GetSymbolRunes()

	for i, v := range c.pc {
		l := Legend{v, c.newData.Count[v], symbols[i].Code}
		c.newData.Legend = append(c.newData.Legend, l)
	}

	sort.Slice(c.newData.Legend, func(i, j int) bool { return c.newData.Legend[i].Color.ID < c.newData.Legend[j].Color.ID })
	return nil
}

func (c *Converter) convertImageChunk(xlow, xhigh, ylow, yhigh int) map[palette.Thread]int {
	count := make(map[palette.Thread]int)
	for y := ylow; y < yhigh; y++ {
		for x := xlow; x < xhigh; x++ {
			minIndex := c.setNewPixel(x, y)
			count[c.pc[minIndex]] += 1
		}
	}
	return count
}

func (c *Converter) setNewPixel(x, y int) int {
	r32, g32, b32, a := c.image.At(x, y).RGBA()
	r, g, b := uint8(r32), uint8(g32), uint8(b32)

	symbols := palette.GetSymbolRunes()

	minLen := math.MaxFloat64
	minIndex := 0
	for i := 0; i < len(c.pc); i++ {
		var dist float64
		if c.rgb || colorConverter.Greyscale(r, g, b) < 100 {
			dist = colorConverter.RGBDistance(float64(r), float64(g), float64(b), float64(c.pc[i].RGB.R), float64(c.pc[i].RGB.G), float64(c.pc[i].RGB.B))
		} else {
			cie := colorConverter.SRGBToCIELab(colorConverter.SRGB{R: r, G: g, B: b})
			dist = colorConverter.LABDistance(c.pc[i].LAB.L, c.pc[i].LAB.A, c.pc[i].LAB.B, cie.L, cie.A, cie.B)
		}
		if dist < minLen {
			minLen = dist
			minIndex = i
		}
	}

	c.newData.Symbols[y][x].Symbol = symbols[minIndex]
	c.newData.Symbols[y][x].Color = c.pc[minIndex]
	c.newData.Symbols[y][x].Text = getTextColor(c.pc[minIndex])
	c.newData.Image.Set(x, y, color.RGBA{uint8(c.pc[minIndex].RGB.R), uint8(c.pc[minIndex].RGB.G), uint8(c.pc[minIndex].RGB.B), uint8(a)})

	return minIndex
}

func getTextColor(cc palette.Thread) string {
	gg := colorConverter.Greyscale(cc.RGB.R, cc.RGB.G, cc.RGB.B)

	if gg < 100 {
		return "white"
	}
	return "black"
}

// start with 32 colors
func (c *Converter) colorQuant() []colorConverter.SRGB {
	bounds := c.image.Bounds()

	allcolors := make([]colorConverter.SRGB, (bounds.Dy()-bounds.Min.Y)*(bounds.Dx()-bounds.Min.X))
	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		for x := bounds.Min.X; x < bounds.Dx(); x++ {
			r32, g32, b32, _ := c.image.At(x, y).RGBA()
			allcolors[y*(bounds.Dx()-bounds.Min.X)+x] = colorConverter.SRGB{R: uint8(r32), G: uint8(g32), B: uint8(b32)}
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
				sort.Slice(s, func(i, j int) bool { return s[i].R < s[j].R })
			} else if yr > xr && yr > zr {
				sort.Slice(s, func(i, j int) bool { return s[i].G < s[j].G })
			} else {
				sort.Slice(s, func(i, j int) bool { return s[i].B < s[j].B })
			}
		}

		// insert 2^n more slice indexes
		max := len(slices) - 1
		for k := 0; k < max; k++ {
			slices = append(slices, (slices[k+1]+slices[k])/2)
		}
		sort.Ints(slices)
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

		bestColors[i] = colorConverter.SRGB{R: uint8(math.Sqrt(avgR / float64(len(s)))), G: uint8(math.Sqrt(avgG / float64(len(s)))), B: uint8(math.Sqrt(avgB / float64(len(s))))}
	}

	return bestColors
}
