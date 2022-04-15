package convert

import (
	"fmt"
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
	imgPath, _, imgErr := c.writePNG()
	if imgErr != nil {
		return imgErr
	}
	fmt.Printf("Wrote new PNG to %s\n", imgPath)

	// write PDF instructions
	pdfPath, pdfErr := c.writePDF(imgPath)
	if pdfErr != nil {
		return pdfErr
	}
	fmt.Printf("Wrote PDF to %s\n", pdfPath)

	return nil
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
