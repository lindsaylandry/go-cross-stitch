package convert

import (
	"bytes"
	"fmt"
	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/jung-kurt/gofpdf"
	"io/ioutil"
	"strconv"
)

func (c *Converter) WritePDF() (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 10)

	for _, y := range c.newImage.symbols {
		for _, x := range y {
			symbol := rune(x.Code)
			pdf.CellFormat(5.0, 5.0, strconv.QuoteRuneToGraphic(symbol), "1", 0, "CM", false, 0, "")
		}
		pdf.Ln(-1.0)
	}

	path := c.getPath("pdf")
	return path, pdf.OutputFileAndClose(path)
}

func (c *Converter) WritePDFFromHTML() (string, error) {
	path := c.getPath("html")
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return "", err
	}

	htmlfile, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	page2 := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlfile))
	pdfg.AddPage(page2)
	err = pdfg.Create()
	if err != nil {
		return "", err
	}

	newPath := c.getPath("pdf")
	err = pdfg.WriteFile(newPath)
	if err != nil {
		return "", err
	}
	fmt.Printf("PDF size %vkB\n", len(pdfg.Bytes())/1024)

	return path, nil
}
