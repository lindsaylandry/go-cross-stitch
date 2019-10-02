package convert

import (
	"bytes"
	"fmt"
	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io/ioutil"
)

func (c *Converter) writePDFFromHTML() (string, error) {
	path := c.getPath("html")
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	pdfg.Dpi.Set(300)
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
