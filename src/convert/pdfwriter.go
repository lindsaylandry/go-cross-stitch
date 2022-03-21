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

	// Title
	pdf.AddPage()
	pdf.SetFont("aaa", "", 32)
	pdf.CellFormat(100, 30.0, "TITLE", "", 1, "LM", false, 0, "")

	// Image
	pdf.Image(imgPath, 10, 20, 190, 0, true, "", 0, "")

	// Info
	pdf.SetFont("aaa", "", 8)
	pdf.CellFormat(190.0, 50.0, "INFO", "1", 1, "LM", false, 0, "")

	// Legend
	pdf.AddPage()
	pdf.SetFont("aaa", "", 32)
	pdf.CellFormat(100.0, 20.0, "LEGEND", "", 1, "LM", false, 0, "")
	pdf.SetFont("aaa", "", 8)
	// header cells
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(15.0, 5.0, "Color", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(20.0, 5.0, "Symbol", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(20.0, 5.0, "Color ID", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(20.0, 5.0, "Count", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(30.0, 5.0, "Name", "1", 1, "CM", true, 0, "")

	// body cells
	for i := 0; i < len(c.newImage.legend); i++ {
		pdf.SetFillColor(int(c.newImage.legend[i].Color.RGB.R), int(c.newImage.legend[i].Color.RGB.G), int(c.newImage.legend[i].Color.RGB.B))
		pdf.CellFormat(15.0, 5.0, "", "1", 0, "CM", true, 0, "")
		pdf.CellFormat(20.0, 5.0, string(c.newImage.legend[i].Symbol), "1", 0, "CM", false, 0, "")
		pdf.CellFormat(20.0, 5.0, c.newImage.legend[i].Color.StringID, "1", 0, "RM", false, 0, "")
		pdf.CellFormat(20.0, 5.0, strconv.Itoa(c.newImage.legend[i].Count), "1", 0, "RM", false, 0, "")
		pdf.CellFormat(30.0, 5.0, c.newImage.legend[i].Color.Name, "1", 1, "LM", false, 0, "")
	}

	// Create Cells (color)
	// TODO: white vs black font
	pdf.AddPage()
	pdf.SetFont("aaa", "", 6)
	for y := 0; y <= bounds.Max.Y; y++ {
		for x := 0; x <= bounds.Max.X; x++ {
			ln := 0
			if x == bounds.Max.X {
				ln = 1
			}
			if y == 0 {
				pdf.SetFillColor(200, 200, 200)
				xLabel := ""
				if x%10 == 0 {
					xLabel = strconv.Itoa(x)
				}
				pdf.CellFormat(2.5, 2.5, xLabel, "1", ln, "CM", true, 0, "")
			} else if x == 0 {
				pdf.SetFillColor(200, 200, 200)
				yLabel := ""
				if y%10 == 0 {
					yLabel = strconv.Itoa(y)
				}
				pdf.CellFormat(2.5, 2.5, yLabel, "1", ln, "CM", true, 0, "")
			} else {
				pdf.SetFillColor(int(c.newImage.symbols[y-1][x-1].Color.RGB.R), int(c.newImage.symbols[y-1][x-1].Color.RGB.G), int(c.newImage.symbols[y-1][x-1].Color.RGB.B))
				pdf.CellFormat(2.5, 2.5, string(c.newImage.symbols[y-1][x-1].Symbol.Code), "1", ln, "CM", true, 0, "")
			}
		}
	}

	// Create Cells (bw)
	// TOOO: line width on 10 spaces
	pdf.AddPage()

	for y := 0; y <= bounds.Max.Y; y++ {
		for x := 0; x <= bounds.Max.X; x++ {
			ln := 0
			if x == bounds.Max.X {
				ln = 1
			}
			if y == 0 {
				pdf.SetFillColor(200, 200, 200)
				xLabel := ""
				if x%10 == 0 {
					xLabel = strconv.Itoa(x)
				}
				pdf.CellFormat(2.5, 2.5, xLabel, "1", ln, "CM", true, 0, "")
			} else if x == 0 {
				pdf.SetFillColor(200, 200, 200)
				yLabel := ""
				if y%10 == 0 {
					yLabel = strconv.Itoa(y)
				}
				pdf.CellFormat(2.5, 2.5, yLabel, "1", ln, "CM", true, 0, "")
			} else {
				pdf.SetLineWidth(0.2)
				border := "1"
				if x%10 == 0 && y%10 == 0 {
					pdf.SetLineWidth(0.3)
					border = "RB"
				} else if x%10 == 0 && y%2 == 0 {
					pdf.SetLineWidth(0.3)
					border = "R"
				} else if x%10 == 1 && y%2 == 1 {
					pdf.SetLineWidth(0.3)
					border = "L"
				} else if y%10 == 0 && x%2 == 0 {
					pdf.SetLineWidth(0.3)
					border = "B"
				} else if y%10 == 1 && x%2 == 1 {
					pdf.SetLineWidth(0.3)
					border = "T"
				}
				pdf.CellFormat(2.5, 2.5, string(c.newImage.symbols[y-1][x-1].Symbol.Code), border, ln, "CM", false, 0, "")
			}
		}
	}

	err := pdf.OutputFileAndClose(path)

	return path, err
}
