package writer

import (
	"github.com/jung-kurt/gofpdf"

	"strconv"
	//"unicode/utf8"
	"fmt"
)

type Grid struct {
	Xstart, Ystart int
	Xend, Yend     int
}

func (w *Writer) writePDF(imgPath string) (string, error) {
	// Setup pdf
	path := w.getPath("pdf")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("aaa", "", "fonts/arial-unicode.ttf")
	//pdf.AddUTF8Font("aaa", "", "fonts/NotoSansSymbols.ttf")
	pdf.SetAutoPageBreak(true, 1.5)

	bounds := w.data.Image.Bounds()

	// Title
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 32)
	pdf.CellFormat(100, 30.0, w.title, "", 1, "LM", false, 0, "")

	// Image
	pdf.Image(imgPath, 10, 20, 190, 0, true, "", 0, "")

	// Info
	aida := 14
	aidaColor := "Black"

	widthInches := float64(bounds.Max.X) / float64(aida)
	heightInches := float64(bounds.Max.Y) / float64(aida)

	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(100.0, 20.0, "INFO", "", 1, "LM", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(90.0, 5.5, "Fabric:", "", 0, "LM", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(100.0, 5.5, fmt.Sprintf("%dct Aida %s", aida, aidaColor), "", 1, "RM", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(90.0, 5.5, "Size:", "", 0, "LM", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(100.0, 5.5, fmt.Sprintf("%.1fx%.1fin", widthInches, heightInches), "", 1, "RM", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(90.0, 5.5, "Color Scheme:", "", 0, "LM", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(100.0, 5.5, w.data.Scheme, "", 1, "RM", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(90.0, 5.5, "Number of Colors:", "", 0, "LM", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(100.0, 5.5, strconv.Itoa(len(w.data.Legend)), "", 1, "RM", false, 0, "")

	// TODO: figure out how far down to put the info box
	ratio := 190 * float64(bounds.Max.Y) / float64(bounds.Max.X)
	pdf.Rect(10, 10+ratio+20+25, 190, 30, "D")

	// Legend
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(100.0, 20.0, "Legend", "", 1, "LM", false, 0, "")
	pdf.SetFont("Arial", "B", 8)
	// header cells
	pdf.CellFormat(10.0, 4.5, "Color", "", 0, "CM", false, 0, "")
	pdf.CellFormat(15.0, 4.5, "Symbol", "", 0, "CM", false, 0, "")
	pdf.CellFormat(15.0, 4.5, "Color ID", "", 0, "CM", false, 0, "")
	pdf.CellFormat(15.0, 4.5, "Stitches", "", 0, "CM", false, 0, "")
	pdf.CellFormat(35.0, 4.5, "Color Description", "", 1, "LM", false, 0, "")

	// body cells
	pdf.SetFont("aaa", "", 8)
	for i := 0; i < len(w.data.Legend); i++ {
		fill := false
		pdf.SetFillColor(int(w.data.Legend[i].Color.RGB.R), int(w.data.Legend[i].Color.RGB.G), int(w.data.Legend[i].Color.RGB.B))
		pdf.CellFormat(10.0, 4.5, "", "", 0, "CM", true, 0, "")

		if i%2 == 0 {
			pdf.SetFillColor(200, 200, 200)
			fill = true
		}
		pdf.CellFormat(15.0, 4.5, string(w.data.Legend[i].Symbol), "", 0, "CM", fill, 0, "")
		pdf.CellFormat(15.0, 4.5, w.data.Legend[i].Color.StringID, "", 0, "RM", fill, 0, "")
		pdf.CellFormat(15.0, 4.5, strconv.Itoa(w.data.Legend[i].Count), "", 0, "RM", fill, 0, "")
		pdf.CellFormat(35.0, 4.5, w.data.Legend[i].Color.Name, "", 1, "LM", fill, 0, "")

		if i == 0 {
			pdf.SetLineWidth(0.4)
			pdf.Line(10.0, 10.0+20.0+4.5, 10.0+10.0+15.0+15.0+15.0+35.0, 10.0+20.0+4.5)
		}
	}

	// TODO: 70x100 chunks
	maxChunkX := 70
	maxChunkY := 100

	xnum := bounds.Max.X / maxChunkX
	if bounds.Max.X%maxChunkX != 0 {
		xnum += 1
	}

	ynum := bounds.Max.Y / maxChunkY
	if bounds.Max.Y%maxChunkY != 0 {
		ynum += 1
	}

	grids := []Grid{}

	for y := 0; y < ynum; y++ {
		for x := 0; x < xnum; x++ {
			grid := Grid{}
			grid.Xstart = x * maxChunkX
			if (x+1)*maxChunkX <= bounds.Max.X {
				grid.Xend = (x + 1) * maxChunkX
			} else {
				grid.Xend = bounds.Max.X
			}

			grid.Ystart = y * maxChunkY
			if (y+1)*maxChunkY <= bounds.Max.Y {
				grid.Yend = (y + 1) * maxChunkY
			} else {
				grid.Yend = bounds.Max.Y
			}

			grids = append(grids, grid)
		}
	}

	// Create Cells (bw)
	for _, grid := range grids {
		pdf.AddPage()
		w.CreateGrid(pdf, grid, false)
		setGridLines(pdf, grid)
	}

	// Create Cells (color)
	for _, grid := range grids {
		pdf.AddPage()
		w.CreateGrid(pdf, grid, true)
		setGridLines(pdf, grid)
	}

	err := pdf.OutputFileAndClose(path)

	return path, err
}

func (w *Writer) CreateGrid(pdf *gofpdf.Fpdf, grid Grid, color bool) {
	pdf.SetLineWidth(0.1)
	for y := grid.Ystart; y <= grid.Yend; y++ {
		for x := grid.Xstart; x <= grid.Xend; x++ {
			pdf.SetTextColor(0, 0, 0)
			ln := 0
			if x == grid.Xend {
				ln = 1
			}
			if y == grid.Ystart {
				pdf.SetFillColor(200, 200, 200)
				xLabel := ""
				if x%10 == 0 {
					xLabel = strconv.Itoa(x)
				}
				pdf.SetFont("Arial", "B", 5)
				pdf.CellFormat(2.5, 2.5, xLabel, "1", ln, "CM", true, 0, "")
			} else if x == grid.Xstart {
				pdf.SetFillColor(200, 200, 200)
				yLabel := ""
				if y%10 == 0 {
					yLabel = strconv.Itoa(y)
				}
				pdf.SetFont("Arial", "B", 5)
				pdf.CellFormat(2.5, 2.5, yLabel, "1", ln, "CM", true, 0, "")
			} else {
				fill := false
				if color {
					// TODO: set text color
					brightness := (int(w.data.Symbols[y-1][x-1].Color.RGB.R) + int(w.data.Symbols[y-1][x-1].Color.RGB.G) + int(w.data.Symbols[y-1][x-1].Color.RGB.B)) / 3
					if brightness < 100 {
						pdf.SetTextColor(255, 255, 255)
					}
					pdf.SetFillColor(int(w.data.Symbols[y-1][x-1].Color.RGB.R), int(w.data.Symbols[y-1][x-1].Color.RGB.G), int(w.data.Symbols[y-1][x-1].Color.RGB.B))
					fill = true
				}
				pdf.SetFont("aaa", "", 6)
				pdf.CellFormat(2.5, 2.5, string(w.data.Symbols[y-1][x-1].Symbol.Code), "1", ln, "CM", fill, 0, "")
			}
		}
	}
}

func setGridLines(pdf *gofpdf.Fpdf, grid Grid) {
	pdf.SetLineWidth(0.3)
	xRange := grid.Xend - grid.Xstart
	yRange := grid.Yend - grid.Ystart
	endX := 10.0 + 2.5*float64(xRange+1)
	endY := 10.0 + 2.5*float64(yRange+1)
	for x := 0; x <= xRange; x += 10 {
		if x == 0 {
			pdf.Line(10.0, 10.0, 10.0, endY)
			pdf.Line(10.0+2.5, 10.0, 10.0+2.5, endY)
		} else {
			pdf.Line(10.0+2.5+2.5*float64(x), 10.0, 10.0+2.5+2.5*float64(x), endY)
		}
	}

	for y := 0; y <= yRange; y += 10 {
		if y == 0 {
			pdf.Line(10.0, 10.0, endX, 10.0)
			pdf.Line(10.0, 10.0+2.5, endX, 10.0+2.5)
		} else {
			pdf.Line(10.0, 10.0+2.5+2.5*float64(y), endX, 10.0+2.5+2.5*float64(y))
		}
	}
	// R and B borders
	pdf.Line(endX, 10.0, endX, endY)
	pdf.Line(10.0, endY, endX, endY)
}
