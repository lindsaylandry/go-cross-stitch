package convert

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log/slog"
	"math"
	"os"
	"sort"

	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
	"github.com/lindsaylandry/go-cross-stitch/src/config"
	"github.com/lindsaylandry/go-cross-stitch/src/palette"
)

type Legend struct {
	Color  palette.Color
	Count  int
	Symbol rune
}

type ColorSymbol struct {
	Symbol palette.SymbolRune
	Color  palette.Color
	Text   string
}

type NewData struct {
	Image   *image.RGBA
	Count   map[palette.Color]int
	Legend  []Legend
	Symbols [][]ColorSymbol
	Path    string
	Extra   string
	Scheme  string
	Type    config.Type
}

type Converter struct {
	image     image.Image
	newData   NewData
	rgb       bool
	pc        []palette.Color
	greyscale bool
}

func NewConverter(filename string, config *config.Config) (*Converter, error) {
	c := Converter{}

	c.newData.Path = filename

	if err := c.getImage(); err != nil {
		return &c, err
	}

	if config.Width > 0 {
		dst := resize(c.image, config.Width)
		c.image = dst
	}

	bounds := c.image.Bounds()

	c.newData.Image = image.NewRGBA(bounds)

	c.newData.Symbols = make([][]ColorSymbol, bounds.Dy()-bounds.Min.Y)
	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		c.newData.Symbols[y] = make([]ColorSymbol, bounds.Dx()-bounds.Min.X)
	}

	c.newData.Count = make(map[palette.Color]int)

	c.rgb = config.Rgb
	c.greyscale = config.Greyscale

	if c.rgb {
		c.newData.Extra = "-" + config.Palette + "-rgb"
	} else {
		c.newData.Extra = "-" + config.Palette + "-lab"
	}

	if config.Palette == "original" {
		if !config.Quantize.Enabled {
			slog.Error(fmt.Sprintf("Cannot use %s palette without enabling convertPalette. Check config settings and try again.", config.Palette))
		}
		c.newData.Scheme = "Quantize"

		c.pc = palette.ConvertOriginal(c.colorQuant(config.Quantize.N))
		return &c, nil
	}

	csvFile := config.CsvFile
	if csvFile == "" {
		csvFile = config.Palette
	}

	pc, err := palette.ReadCSV(csvFile, config.Excludes)

	if err != nil {
		return &c, err
	}
	c.pc = pc

	if config.Palette == "dmc" || config.Palette == "anchor" || config.Palette == "lego" {
		if config.Palette == "lego" {
			c.newData.Scheme = "LEGO"
			c.newData.Type = config.Lego
		} else if config.Palette == "dmc" {
			c.newData.Scheme = "DMC"
			c.newData.Type = config.DMC
		} else {
			c.newData.Scheme = "Anchor"
			c.newData.Type = config.Anchor
		}

		if config.Quantize.Enabled {
			// best colors
			c.pc = c.convertPalette(c.colorQuant(config.Quantize.N))
		} else {
			//most colors
			c.pc = c.convertPalette(c.pixel())
		}
		slog.Info(fmt.Sprintf("Number of Colors: %d\n", len(c.pc)))
	} else if config.Palette == "bw" {
		c.newData.Scheme = "Black&White"
	} else {
		return &c, errors.New("--color not recognized")
	}

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

func (c *Converter) Convert(dither bool) (NewData, error) {
	err := c.convertImage(dither)
	return c.newData, err
}

func (c *Converter) getImage() error {
	data, err := os.ReadFile(c.newData.Path)
	if err != nil {
		return err
	}

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)

	c.image, _, err = image.Decode(bytes.NewReader(data))
	return err
}

// convert image best-colors to match available palette
func (c *Converter) convertPalette(colors []colorConverter.SRGB) []palette.Color {
	dict := make(map[palette.Color]int)
	var legend []palette.Color
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

func (c *Converter) convertImage(dither bool) error {
	if dither {
		c.floydSteinbergDither()
	} else {
		bounds := c.image.Bounds()

		slog.Info(fmt.Sprintf("Converting image to %s palette...", c.newData.Scheme))
		// TODO: make variable amount of chunks based on image size
		n := 4
		countChan := make(chan map[palette.Color]int, n*n)

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
		if c.newData.Count[v] == 0 {
			continue
		}
		l := Legend{v, c.newData.Count[v], symbols[i].Code}
		c.newData.Legend = append(c.newData.Legend, l)
	}

	sort.Slice(c.newData.Legend, func(i, j int) bool { return c.newData.Legend[i].Color.ID < c.newData.Legend[j].Color.ID })

	slog.Info("Done")

	return nil
}

func (c *Converter) convertImageChunk(xlow, xhigh, ylow, yhigh int) map[palette.Color]int {
	count := make(map[palette.Color]int)
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

func getTextColor(cc palette.Color) string {
	gg := colorConverter.Greyscale(cc.RGB.R, cc.RGB.G, cc.RGB.B)

	if gg < 100 {
		return "white"
	}
	return "black"
}

// start with 32 colors
func (c *Converter) colorQuant(n int) []colorConverter.SRGB {
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

	for i := 0; i < n; i++ {
		for j := 0; j < len(slices)-1; j++ {
			// get a slice of allcolors
			s := allcolors[slices[j]:slices[j+1]]

			colorRanges := [][]uint8{{math.MaxUint8, 0}, {math.MaxUint8, 0}, {math.MaxUint8, 0}}
			for k := 0; k < len(s); k++ {
				//R
				if s[k].R < colorRanges[0][0] {
					colorRanges[0][0] = s[k].R
				}
				if s[k].R > colorRanges[0][1] {
					colorRanges[0][1] = s[k].R
				}
				//G
				if s[k].G < colorRanges[1][0] {
					colorRanges[1][0] = s[k].G
				}
				if s[k].G > colorRanges[1][1] {
					colorRanges[1][1] = s[k].G
				}
				//B
				if s[k].B < colorRanges[2][0] {
					colorRanges[2][0] = s[k].B
				}
				if s[k].B > colorRanges[2][1] {
					colorRanges[2][1] = s[k].B
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
		for k := 0; k < len(s); k++ {
			avgR = avgR + math.Pow(float64(s[k].R), 2)
			avgG = avgG + math.Pow(float64(s[k].G), 2)
			avgB = avgB + math.Pow(float64(s[k].B), 2)
		}

		bestColors[i] = colorConverter.SRGB{R: uint8(math.Sqrt(avgR / float64(len(s)))), G: uint8(math.Sqrt(avgG / float64(len(s)))), B: uint8(math.Sqrt(avgB / float64(len(s))))}
	}

	// TODO: filter out any duplicates

	slog.Debug("Best Colors RGB Values:")
	for _, b := range bestColors {
		slog.Debug(fmt.Sprintf("R: %v, G: %v, B: %v", b.R, b.G, b.B))
	}

	return bestColors
}
