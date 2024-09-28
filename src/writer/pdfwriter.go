package writer

import (
	"github.com/jung-kurt/gofpdf"

	"strconv"
	"fmt"
)

type Grid struct {
	Xstart, Ystart int
	Xend, Yend     int
}

func (w *Writer) writePDF(imgPath string, paperSize string) (string, error) {
	// TODO: multiplier for paper sizes
	var mult float64
	var maxChunkX, maxChunkY int

	switch paperSize {
  case "A4":
  	mult = 1.0
	  maxChunkX = 70
    maxChunkY = 100
  case "A2":
		mult = 2.0
    maxChunkX = 200
    maxChunkY = 200
  case "A1":
		mult = 2.5
    maxChunkX = 200
    maxChunkY = 400
  default:
    maxChunkX = 100
    maxChunkY = 100
  }

	// Setup pdf
	pdf := gofpdf.New("P", "mm", paperSize, "")
	pdf.AddUTF8Font("aaa", "", "fonts/arial-unicode.ttf")
	//pdf.AddUTF8Font("aaa", "", "fonts/NotoSansSymbols.ttf")
	pdf.SetAutoPageBreak(true, 1.5)

	bounds := w.data.Image.Bounds()

	// Title
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 32.0*mult)
	pdf.CellFormat(100, 30.0, w.title, "", 1, "LM", false, 0, "")

	// Image
	pdf.Image(imgPath, 10, 20, 190.0*mult, 0, true, "", 0, "")

	// Info
	aida := 14
	aidaColor := "Black"

	widthInches := float64(bounds.Max.X) / float64(aida)
	heightInches := float64(bounds.Max.Y) / float64(aida)

	pdf.SetFont("Arial", "B", 18*mult)
	pdf.CellFormat(100.0*mult, 20.0*mult, "INFO", "", 1, "LM", false, 0, "")

	pdf.SetFont("Arial", "B", 12*mult)
	pdf.CellFormat(90.0*mult, 5.5*mult, "Fabric:", "", 0, "LM", false, 0, "")
	pdf.SetFont("Arial", "", 12*mult)
	pdf.CellFormat(100.0*mult, 5.5*mult, fmt.Sprintf("%dct Aida %s", aida, aidaColor), "", 1, "RM", false, 0, "")

	pdf.SetFont("Arial", "B", 12*mult)
	pdf.CellFormat(90.0*mult, 5.5*mult, "Size:", "", 0, "LM", false, 0, "")
	pdf.SetFont("Arial", "", 12*mult)
	pdf.CellFormat(100.0*mult, 5.5*mult, fmt.Sprintf("%.1fx%.1fin", widthInches, heightInches), "", 1, "RM", false, 0, "")

	pdf.SetFont("Arial", "B", 12*mult)
	pdf.CellFormat(90.0*mult, 5.5*mult, "Color Scheme:", "", 0, "LM", false, 0, "")
	pdf.SetFont("Arial", "", 12*mult)
	pdf.CellFormat(100.0*mult, 5.5*mult, w.data.Scheme, "", 1, "RM", false, 0, "")

	pdf.SetFont("Arial", "B", 12*mult)
	pdf.CellFormat(90.0*mult, 5.5*mult, "Number of Colors:", "", 0, "LM", false, 0, "")
	pdf.SetFont("Arial", "", 12*mult)
	pdf.CellFormat(100.0*mult, 5.5*mult, strconv.Itoa(len(w.data.Legend)), "", 1, "RM", false, 0, "")

	// TODO: figure out how far down to put the info box
	ratio := 190*mult * float64(bounds.Max.Y) / float64(bounds.Max.X)
	pdf.Rect(10, 10+ratio+20*mult+25*mult, 190*mult, 30*mult, "D")

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

	// Find Middle Triangles
	midX := float64(bounds.Max.X) / 2.0
	midY := float64(bounds.Max.Y) / 2.0

	// Create Cells (bw)
	for _, grid := range grids {
		pdf.AddPage()
		w.CreateGrid(pdf, grid, midX, midY, false)
		setGridLines(pdf, grid)
	}

	// Create Cells (color)
	for _, grid := range grids {
		pdf.AddPage()
		w.CreateGrid(pdf, grid, midX, midY, true)
		setGridLines(pdf, grid)
	}

	// Write PDF file
	path := w.getPath("pdf", "-" + paperSize)
	err := pdf.OutputFileAndClose(path)

	return path, err
}

func (w *Writer) CreateGrid(pdf *gofpdf.Fpdf, grid Grid, midX, midY float64, color bool) {
	pdf.SetLineWidth(0.1)
	for y := grid.Ystart; y <= grid.Yend; y++ {
		for x := grid.Xstart; x <= grid.Xend; x++ {
			pdf.SetTextColor(0, 0, 0)
			ln := 0
			if x == grid.Xend {
				ln = 1
			}
			// x-axis labels
			if y == grid.Ystart {
				fill := false
				xLabel := ""
				if x%10 == 0 {
					xLabel = strconv.Itoa(x)
				}
				pdf.SetFont("Arial", "B", 5)
				pdf.CellFormat(2.5, 2.5, xLabel, "", ln, "CM", fill, 0, "")
			// y-axis labels
			} else if x == grid.Xstart {
				fill := false
				yLabel := ""
				if y%10 == 0 {
					yLabel = strconv.Itoa(y)
				}
				pdf.SetFont("Arial", "B", 5)
				pdf.CellFormat(2.5, 2.5, yLabel, "", ln, "CM", fill, 0, "")
			} else {
				fill := false
				if color {
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

	if midX >= float64(grid.Xstart) && midX <= float64(grid.Xend) {
		x := midX - float64(grid.Xstart)
    pt1 := gofpdf.PointType{X: float64(x + 1)*2.5 - 1, Y: 10}
		pt2 := gofpdf.PointType{X: float64(x + 1)*2.5, Y: 10+2.5}
		pt3 := gofpdf.PointType{X: float64(x + 1)*2.5 + 1, Y: 10}
		pts := make([]gofpdf.PointType, 0, 3)
		pts = append(pts,	pt1)
		pts = append(pts, pt2)
		pts = append(pts, pt3)
		
		pdf.SetFillColor(0, 0, 0)
		pdf.Polygon(pts, "FD")
  }

	if midY >= float64(grid.Ystart) && midY <= float64(grid.Yend) {
    y := midY - float64(grid.Ystart)
    pt1 := gofpdf.PointType{X: 10, Y: float64(y + 1)*2.5 - 1}
    pt2 := gofpdf.PointType{X: 10+2.5, Y: float64(y + 1)*2.5}
    pt3 := gofpdf.PointType{X: 10, Y: float64(y + 1)*2.5 + 1}
    pts := make([]gofpdf.PointType, 0, 3)
    pts = append(pts, pt1)
    pts = append(pts, pt2)
    pts = append(pts, pt3)
 
    pdf.SetFillColor(0, 0, 0)
    pdf.Polygon(pts, "FD")
  }
}

func setGridLines(pdf *gofpdf.Fpdf, grid Grid) {
	pdf.SetLineWidth(0.3)
	xRange := grid.Xend - grid.Xstart
	yRange := grid.Yend - grid.Ystart
	endX := 10.0 + 2.5*float64(xRange+1)
	endY := 10.0 + 2.5*float64(yRange+1)
	for x := 0; x <= xRange; x += 10 {
		pdf.Line(10.0+2.5*float64(x+1), 10.0+2.5, 10.0+2.5*float64(x+1), endY)
	}

	for y := 0; y <= yRange; y += 10 {
		pdf.Line(10.0+2.5, 10.0+2.5*float64(y+1), endX, 10.0+2.5*float64(y+1))
	}
	// R and B borders
	pdf.Line(endX, 10.0+2.5, endX, endY)
	pdf.Line(10.0+2.5, endY, endX, endY)
}
