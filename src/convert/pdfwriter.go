package convert

import (
	"github.com/jung-kurt/gofpdf"

	"strconv"
	//"unicode/utf8"
)

func (c *Converter) writePDF(imgPath string) (string, error) {
	// Setup pdf
	path := c.getPath("pdf")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("aaa", "", "fonts/arial-unicode.ttf")
	pdf.SetAutoPageBreak(true, 1.5)

	bounds := c.newImage.image.Bounds()

	// Image
	pdf.AddPage()
	pdf.Image(imgPath, 10, 10, 190, 0, false, "", 0, "")

	// Legend
	pdf.AddPage()
	pdf.SetFont("aaa", "", 32)
	pdf.CellFormat(100.0, 20.0, "LEGEND", "", 1, "LM", false, 0, "")
	pdf.SetFont("aaa", "", 8)
	// header cells
	pdf.CellFormat(10.0, 5.0, "Color", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(20.0, 5.0, "Symbol", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(20.0, 5.0, "Count", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(20.0, 5.0, "Color ID", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(30.0, 5.0, "Name", "1", 1, "CM", false, 0, "")

	// body cells
	for i := 0; i < len(c.newImage.legend); i++ {
		pdf.SetFillColor(int(c.newImage.legend[i].Color.RGB.R), int(c.newImage.legend[i].Color.RGB.G), int(c.newImage.legend[i].Color.RGB.B))
		pdf.CellFormat(10.0, 5.0, "", "1", 0, "CM", true, 0, "")
		pdf.CellFormat(20.0, 5.0, string(c.newImage.legend[i].Symbol), "1", 0, "CM", false, 0, "")
		pdf.CellFormat(20.0, 5.0, c.newImage.legend[i].Color.StringID, "1", 0, "RM", false, 0, "")
		pdf.CellFormat(20.0, 5.0, strconv.Itoa(c.newImage.legend[i].Count), "1", 0, "RM", false, 0, "")
		pdf.CellFormat(30.0, 5.0, c.newImage.legend[i].Color.Name, "1", 1, "LM", false, 0, "")
	}

	// Create Cells (color)
	pdf.AddPage()
	pdf.SetFont("aaa", "", 6)
	for y := 0; y < bounds.Max.Y; y++ {
		for x := 1; x <= bounds.Max.X; x++ {
			ln := 0
			if x == bounds.Max.X {
				ln = 1
			}
			pdf.SetFillColor(int(c.newImage.symbols[x-1][y].Color.RGB.R), int(c.newImage.symbols[x-1][y].Color.RGB.G), int(c.newImage.symbols[x-1][y].Color.RGB.B))
			pdf.CellFormat(2.6, 2.6, string(c.newImage.symbols[x-1][y].Symbol.Code), "1", ln, "CM", true, 0, "")
		}
	}

	// Create Cells (bw)
	pdf.AddPage()

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 1; x <= bounds.Max.X; x++ {
			ln := 0
			if x == bounds.Max.X {
				ln = 1
			}
			pdf.CellFormat(2.6, 2.6, string(c.newImage.symbols[x-1][y].Symbol.Code), "1", ln, "CM", false, 0, "")
		}
	}

	err := pdf.OutputFileAndClose(path)

	return path, err
}
