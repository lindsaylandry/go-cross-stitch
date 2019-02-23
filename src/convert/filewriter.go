package convert

import (
	"bytes"
	"encoding/base64"
	"html/template"
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

func mod(i, j int) bool  { return i%j == 0 }
func plus(a, b int) int  { return a + b }
func minus(a, b int) int { return a - b }

func (c *Converter) WriteHTML() (string, error) {
	// encode image to base64 string
	var buff bytes.Buffer
	png.Encode(&buff, c.newImage.image)
	imgString := base64.StdEncoding.EncodeToString(buff.Bytes())

	// struct to send to html
	type AA struct {
		Img        string
		Symbols    [][]int
		Legend     []Legend
		Xlen, Ylen int
	}

	// funcs to use in html template
	fmap := template.FuncMap{
		"forloop": forloop,
		"mod":     mod,
		"plus":    plus,
		"minus":   minus,
	}

	// Write new image to png file
	split := strings.Split(c.path, ".")

	// Write HTML instructions
	var htmlPath string
	if c.rgb {
		htmlPath = split[0] + "-dmc-rgb.html"
	} else {
		htmlPath = split[0] + "-dmc-lab.html"
	}
	htmlFile, err := os.Create(htmlPath)
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("instructions").Funcs(fmap).ParseFiles("templates/instructions.html"))
	err = t.ExecuteTemplate(htmlFile, "instructions.html", AA{Img: imgString, Symbols: c.newImage.symbols, Legend: c.newImage.legend, Xlen: len(c.newImage.symbols[0]), Ylen: len(c.newImage.symbols)})

	return htmlPath, err
}

func (c *Converter) WritePNG() (string, error) {
	// Write new image to png file
	split := strings.Split(c.path, ".")

	var newPath string
	if c.rgb {
		newPath = split[0] + "-dmc-rgb.png"
	} else {
		newPath = split[0] + "-dmc-lab.png"
	}
	place, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer place.Close()

	err = png.Encode(place, c.newImage.image)
	return newPath, err
}
