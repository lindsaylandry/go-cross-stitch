package writer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/lindsaylandry/go-cross-stitch/src/convert"
)

type Writer struct {
	data  convert.NewData
	title string
}

func NewWriter(d convert.NewData) *Writer {
	w := Writer{}
	w.data = d
	w.title = getTitle(d.Path)

	return &w
}

func (w *Writer) WriteFiles() error {
	// write new image file
	imgPath, _, imgErr := w.writePNG()
	if imgErr != nil {
		return imgErr
	}
	fmt.Printf("Wrote new PNG to %s\n", imgPath)

	// write PDF instructions
	pdfPath, pdfErr := w.writePDF(imgPath)
	if pdfErr != nil {
		return pdfErr
	}
	fmt.Printf("Wrote PDF to %s\n", pdfPath)

	return nil
}

func (w *Writer) writePNG() (string, *image.RGBA, error) {
	p := 10
	// Make each pixel 3x3
	bounds := w.data.Image.Bounds()
	bounds.Max.X = bounds.Max.X * p
	bounds.Max.Y = bounds.Max.Y * p
	img := image.NewRGBA(bounds)

	newPath := w.getPath("png")
	place, err := os.Create(newPath)
	if err != nil {
		return "", img, err
	}
	defer place.Close()

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			pixel := w.data.Image.At(x, y)
			for xx := 0; xx < p; xx++ {
				for yy := 0; yy < p; yy++ {
					if xx == p-1 || yy == p-1 {
						px := color.Gray16{0}
						img.Set(x*p+xx, y*p+yy, px)
					} else {
						img.Set(x*p+xx, y*p+yy, pixel)
					}
				}
			}
		}
	}

	err = png.Encode(place, img)
	return newPath, img, err
}

func getTitle(filename string) string {
	fn := strings.SplitAfter(filename, "/")
	n := strings.Split(fn[len(fn)-1], ".")
	name := strings.ReplaceAll(n[0], "-", " ")
	name2 := strings.ReplaceAll(name, "_", " ")

	return strings.ToUpper(name2)
}

func (w *Writer) getPath(extension string) string {
	// Write new image to png file
	split := strings.Split(w.data.Path, ".")

	newPath := split[0] + w.data.Extra + "." + extension

	return newPath
}
