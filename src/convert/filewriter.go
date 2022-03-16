package convert

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

func forloop(start, end int) (stream chan int) {
	stream = make(chan int)
	go func() {
		for i := start; i <= end; i++ {
			stream <- i
		}
		close(stream)
	}()
	return
}

func mod(i, j int) bool         { return i%j == 0 }
func plus(a, b int) int         { return a + b }
func minus(a, b int) int        { return a - b }
func mult(a, b int) int         { return a * b }
func div(a, b int) float32      { return float32(a) / float32(b) }
func fmtFloat(a float32) string { return fmt.Sprintf("%.1f", a) }

func (c *Converter) WriteFiles() error {
	// write new image file
	path, img, err := c.writePNG()
	if err != nil {
		return err
	}
	fmt.Printf("Wrote new PNG to %s\n", path)

	// write HTML instructions
	path, err = c.writeHTML(img)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote instructions to %s\n", path)

	// write PDF instructions
	path, err = c.writePDF(img)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote PDF to %s\n", path)

	return nil
}

func (c *Converter) writeHTML(img *image.RGBA) (string, error) {
	// encode image to base64 string
	var buff bytes.Buffer
	png.Encode(&buff, img)
	imgString := base64.StdEncoding.EncodeToString(buff.Bytes())

	type Table struct {
		Symbols        [][]ColorSymbol
		Xstart, Ystart int
		Xend, Yend     int
	}

	// struct to send to html
	type AA struct {
		Img       string
		Legend    []Legend
		Tables    []Table
		Width     int
		Height    int
		Scheme    string
		Title     string
		ColorGrid bool
	}

	// funcs to use in html template
	fmap := template.FuncMap{
		"forloop":  forloop,
		"mod":      mod,
		"plus":     plus,
		"minus":    minus,
		"mult":     mult,
		"div":      div,
		"fmtFloat": fmtFloat,
	}

	htmlPath := c.getPath("html")
	htmlFile, err := os.Create(htmlPath)
	if err != nil {
		return "", err
	}

	aa := AA{
		Img:       imgString,
		Legend:    c.newImage.legend,
		Width:     len(c.newImage.symbols[0]),
		Height:    len(c.newImage.symbols),
		Scheme:    "DMC",
		Title:     c.title,
		ColorGrid: c.colorgrid,
	}

	xchunk := 50
	ychunk := 60

	xnum := aa.Width / xchunk
	if aa.Width%xchunk != 0 {
		xnum += 1
	}

	ynum := aa.Height / ychunk
	if aa.Height%ychunk != 0 {
		ynum += 1
	}

	for y := 0; y < ynum; y++ {
		for x := 0; x < xnum; x++ {
			table := Table{}
			table.Xstart = x * xchunk
			if (x+1)*xchunk <= aa.Width {
				table.Xend = (x+1)*xchunk - 1
			} else {
				table.Xend = aa.Width - 1
			}

			table.Ystart = y * ychunk
			if (y+1)*ychunk <= aa.Height {
				table.Yend = (y+1)*ychunk - 1
			} else {
				table.Yend = aa.Height - 1
			}

			s := make([][]ColorSymbol, table.Yend-table.Ystart+1)
			for i := table.Ystart; i <= table.Yend; i++ {
				s2 := make([]ColorSymbol, table.Xend-table.Xstart+1)
				for j := table.Xstart; j <= table.Xend; j++ {
					s2[j-table.Xstart] = c.newImage.symbols[i][j]
				}
				s[i-table.Ystart] = s2
			}

			table.Symbols = s

			aa.Tables = append(aa.Tables, table)
		}
	}

	t := template.Must(template.New("instructions").Funcs(fmap).ParseFiles("templates/instructions.html"))
	err = t.ExecuteTemplate(htmlFile, "instructions.html", aa)

	return htmlPath, err
}

func (c *Converter) writePNG() (string, *image.RGBA, error) {
	// Make each pixel 3x3
	bounds := c.newImage.image.Bounds()
	bounds.Max.X = bounds.Max.X * c.newImage.p
	bounds.Max.Y = bounds.Max.Y * c.newImage.p
	img := image.NewRGBA(bounds)

	newPath := c.getPath("png")
	place, err := os.Create(newPath)
	if err != nil {
		return "", img, err
	}
	defer place.Close()

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			pixel := c.newImage.image.At(x, y)
			for xx := 0; xx < c.newImage.p; xx++ {
				for yy := 0; yy < c.newImage.p; yy++ {
					if xx == c.newImage.p-1 || yy == c.newImage.p-1 {
						px := color.Gray16{0}
						img.Set(x*c.newImage.p+xx, y*c.newImage.p+yy, px)
					} else {
						img.Set(x*c.newImage.p+xx, y*c.newImage.p+yy, pixel)
					}
				}
			}
		}
	}

	err = png.Encode(place, img)
	return newPath, img, err
}

func (c *Converter) getPath(extension string) string {
	// Write new image to png file
	split := strings.Split(c.path, ".")

	newPath := split[0] + c.extra + "." + extension

	return newPath
}
