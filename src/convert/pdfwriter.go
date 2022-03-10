package convert

import (
	"github.com/jung-kurt/gofpdf"

	"image"
	//"unicode/utf8"
)

func (c *Converter) writePDF(*image.RGBA) (string, error) {
	path := c.getPath("pdf")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("aaa", "", "fonts/arial-unicode.ttf")
	pdf.SetFont("aaa", "", 7)
	pdf.SetAutoPageBreak(true, 1.5)
	pdf.AddPage()

	// Create Cells
	maxX := 60
	maxY := 90

	code := '\u9312'

	for i := 0; i < maxY; i++ {
		for j := 1; j <= maxX; j++ {
			ln := 0
			if j == maxX {
				ln = 1
			}
			pdf.CellFormat(3.0, 3.0, string(code), "1", ln, "CM", false, 0, "")
		}
	}

	err := pdf.OutputFileAndClose(path)

	return path, err
}
