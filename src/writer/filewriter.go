package writer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log/slog"
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
	slog.Info("Writing new image file...")
	imgPath, _, imgErr := w.writePNG()
	if imgErr != nil {
		return imgErr
	}
	slog.Info(fmt.Sprintf("Wrote new PNG to %s\n", imgPath))

	paperSizes := [3]string{"A4", "A2", "A1"}

	// write PDF instructions
	slog.Info("Writing PDF instructions...")
	for _, p := range paperSizes {
		pdfPath, pdfErr := w.writePDF(imgPath, p)
		if pdfErr != nil {
			return pdfErr
		}
		slog.Info(fmt.Sprintf("Wrote PDF to %s\n", pdfPath))
	}

	return nil
}

func (w *Writer) writePNG() (string, *image.RGBA, error) {
	p := 12
	bounds := w.data.Image.Bounds()
	bounds.Max.X = bounds.Max.X * p
	bounds.Max.Y = bounds.Max.Y * p
	img := image.NewRGBA(bounds)

	newPath := w.getPath("png", "")
	place, err := os.Create(newPath)
	if err != nil {
		return "", img, err
	}
	defer place.Close()

	n := 4
	msgChan := make(chan string, n*n)

	for m := 0; m < n; m++ {
		ylow := bounds.Min.Y + m*bounds.Dy()/n
		yhigh := (m + 1) * bounds.Dy() / n

		for p := 0; p < n; p++ {
			xlow := bounds.Min.X + p*bounds.Dx()/n
			xhigh := (p + 1) * bounds.Dx() / n

			go func() {
				msg := w.writeImageChunk(img, xlow, xhigh, ylow, yhigh)
				msgChan <- msg
			}()
		}
	}

	for i := 0; i < n*n; i++ {
		<-msgChan
	}

	slog.Info("Done")

	err = png.Encode(place, img)
	return newPath, img, err
}

func (w *Writer) writeImageChunk(img *image.RGBA, xlow, xhigh, ylow, yhigh int) string {
	p := 12
	for x := xlow; x < xhigh; x++ {
		for y := ylow; y < yhigh; y++ {
			px := w.data.Image.At(x, y)
			for xx := 0; xx < p; xx++ {
				for yy := 0; yy < p; yy++ {
					if xx == p-1 || yy == p-1 {
						img.Set(x*p+xx, y*p+yy, color.Gray16{0})
					} else {
						img.Set(x*p+xx, y*p+yy, px)
					}
				}
			}
		}
	}
	return "done"
}

func getTitle(filename string) string {
	fn := strings.SplitAfter(filename, "/")
	n := strings.Split(fn[len(fn)-1], ".")
	name := strings.ReplaceAll(n[0], "-", " ")
	name2 := strings.ReplaceAll(name, "_", " ")

	return strings.ToUpper(name2)
}

func (w *Writer) getPath(extension string, extra string) string {
	split := strings.Split(w.data.Path, ".")
	newPath := split[0] + w.data.Extra + extra + "." + extension

	return newPath
}
