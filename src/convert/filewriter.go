package convert

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image/png"
	"os"
	"strings"

	"github.com/lindsaylandry/go-cross-stitch/src/palette"
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

func mod(i, j int) bool  { return i%j == 0 }
func plus(a, b int) int  { return a + b }
func minus(a, b int) int { return a - b }
func mult(a, b int) int  { return a * b }
func div(a, b int) int   { return a / b }

func (c *Converter) WriteHTML() (string, error) {
	// encode image to base64 string
	var buff bytes.Buffer
	png.Encode(&buff, c.newImage.image)
	imgString := base64.StdEncoding.EncodeToString(buff.Bytes())

	type Table struct {
		Symbols        [][]palette.Symbol
		Xstart, Ystart int
		Xend, Yend     int
	}

	// struct to send to html
	type AA struct {
		Img    string
		Legend []Legend
		Tables []Table
	}

	// funcs to use in html template
	fmap := template.FuncMap{
		"forloop": forloop,
		"mod":     mod,
		"plus":    plus,
		"minus":   minus,
		"mult":    mult,
		"div":     div,
	}

	htmlPath := c.getPath("html")
	htmlFile, err := os.Create(htmlPath)
	if err != nil {
		return "", err
	}

	aa := AA{
		Img:    imgString,
		Legend: c.newImage.legend,
	}

	xchunk := 50
	ychunk := 60

	xlen := len(c.newImage.symbols[0])
	ylen := len(c.newImage.symbols)

	xnum := xlen / xchunk
	if xlen%xchunk != 0 {
		xnum += 1
	}

	ynum := ylen / ychunk
	if ylen%ychunk != 0 {
		ynum += 1
	}

	for y := 0; y < ynum; y++ {
		for x := 0; x < xnum; x++ {
			table := Table{}
			table.Xstart = x * xchunk
			if (x+1)*xchunk <= xlen {
				table.Xend = (x+1)*xchunk - 1
			} else {
				table.Xend = xlen - 1
			}

			table.Ystart = y * ychunk
			if (y+1)*ychunk <= ylen {
				table.Yend = (y+1)*ychunk - 1
			} else {
				table.Yend = ylen - 1
			}

			fmt.Printf("%d %d %d %d\n", table.Xstart, table.Xend, table.Ystart, table.Yend)

			s := make([][]palette.Symbol, table.Yend-table.Ystart+1)
			for i := table.Ystart; i <= table.Yend; i++ {
				s2 := make([]palette.Symbol, table.Xend-table.Xstart+1)
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

func (c *Converter) WritePNG() (string, error) {
	// Write new image to png file
	newPath := c.getPath("png")
	place, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer place.Close()

	err = png.Encode(place, c.newImage.image)
	return newPath, err
}

func (c *Converter) getPath(extension string) string {
	// Write new image to png file
	split := strings.Split(c.path, ".")

	var newPath string
	if c.rgb {
		newPath = split[0] + "-dmc-rgb." + extension
	} else {
		newPath = split[0] + "-dmc-lab." + extension
	}

	return newPath
}
