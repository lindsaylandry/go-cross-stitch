package convert

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

func (c *Converter) writePDF(imgPath string) (string, error) {
	// Setup pdf
	path := c.getPath("pdf")

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8Font("aaa", "", "fonts/arial-unicode.ttf")
	pdf.SetAutoPageBreak(false, 1.5)

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
	pdf.CellFormat(40.0, 5.0, "Name", "1", 1, "CM", true, 0, "")

	// body cells
	for i := 0; i < len(c.newImage.legend); i++ {
		pdf.SetFillColor(int(c.newImage.legend[i].Color.RGB.R), int(c.newImage.legend[i].Color.RGB.G), int(c.newImage.legend[i].Color.RGB.B))
		pdf.CellFormat(15.0, 5.0, "", "1", 0, "CM", true, 0, "")
		pdf.CellFormat(20.0, 5.0, string(c.newImage.legend[i].Symbol), "1", 0, "CM", false, 0, "")
		pdf.CellFormat(20.0, 5.0, c.newImage.legend[i].Color.StringID, "1", 0, "RM", false, 0, "")
		pdf.CellFormat(20.0, 5.0, strconv.Itoa(c.newImage.legend[i].Count), "1", 0, "RM", false, 0, "")
		pdf.CellFormat(40.0, 5.0, c.newImage.legend[i].Color.Name, "1", 1, "LM", false, 0, "")
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
				grid.Xend = (x+1)*maxChunkX - 1
			} else {
				grid.Xend = bounds.Max.X - 1
			}

			grid.Ystart = y * maxChunkY
			if (y+1)*maxChunkY <= bounds.Max.Y {
				grid.Yend = (y+1)*maxChunkY - 1
			} else {
				grid.Yend = bounds.Max.Y - 1
			}

			grids = append(grids, grid)
		}
	}

	fmt.Println(grids)

	// Create Cells (color)
	// TODO: white vs black font
	for _, grid := range grids {
		pdf.AddPage()
		c.CreateGrid(pdf, grid, true)
		setGridLines(pdf, grid)
	}

	// Create Cells (bw)
	// TOOO: line width on 10 spaces
	for _, grid := range grids {
		pdf.AddPage()
		c.CreateGrid(pdf, grid, false)
		setGridLines(pdf, grid)
	}

	err := pdf.OutputFileAndClose(path)

	return path, err
}

func (c *Converter) CreateGrid(pdf *gofpdf.Fpdf, grid Grid, color bool) {
	pdf.SetFont("aaa", "", 6)
	pdf.SetLineWidth(0.1)
	for y := grid.Ystart; y <= grid.Yend; y++ {
		for x := grid.Xstart; x <= grid.Xend; x++ {
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
				pdf.CellFormat(2.5, 2.5, xLabel, "1", ln, "CM", true, 0, "")
			} else if x == grid.Xstart {
				pdf.SetFillColor(200, 200, 200)
				yLabel := ""
				if y%10 == 0 {
					yLabel = strconv.Itoa(y)
				}
				pdf.CellFormat(2.5, 2.5, yLabel, "1", ln, "CM", true, 0, "")
			} else {
				fill := false
				if color == true {
					pdf.SetFillColor(int(c.newImage.symbols[y-1][x-1].Color.RGB.R), int(c.newImage.symbols[y-1][x-1].Color.RGB.G), int(c.newImage.symbols[y-1][x-1].Color.RGB.B))
					fill = true
				}
				pdf.CellFormat(2.5, 2.5, string(c.newImage.symbols[y-1][x-1].Symbol.Code), "1", ln, "CM", fill, 0, "")
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
