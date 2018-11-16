package convert

import (	
  "os"	
  "image"	
  "image/png"	
  "html/template"	
  "path/filepath"	
  "strings"	
  "bytes"	
  "encoding/base64"	
)

func forloop (start, end int) (stream chan int) {	
    stream = make(chan int)	
    go func() {	
        for i := start; i <= end; i++ {	
            stream <- i	
        }	
        close(stream)	
    }()	
    return	
}	
	
func mod (i, j int) bool { return i%j == 0 }	
func plus (a, b int) int {return a + b}	
func minus (a, b int) int {return a - b}	
	
func WriteHTML(img image.Image, legend []Legend, symbols [][]int, path string) (error){	
  // encode image to base64 string	
  var buff bytes.Buffer	
  png.Encode(&buff, img)	
  imgString := base64.StdEncoding.EncodeToString(buff.Bytes())	
	
  // struct to send to html	
  type AA struct {	
    Img string	
    Symbols [][]int	
    Legend []Legend	
    Xlen, Ylen int	
  }	
	
  // funcs to use in html template	
  fmap := template.FuncMap{	
    "forloop": forloop,	
    "mod": mod,	
    "plus": plus,	
    "minus": minus,	
  }	
	
  // Write new image to png file	
  absPath, err := filepath.Abs(path)	
  absSplit := strings.Split(absPath, ".")	
	
  // Write HTML instructions	
  htmlPath := absSplit[0] + "-dmc.html"	
  htmlFile, err := os.Create(htmlPath)	
  if err != nil { return err }	
  t := template.Must(template.New("instructions").Funcs(fmap).ParseFiles("../../templates/instructions.html"))	
  if err := t.ExecuteTemplate(htmlFile, "instructions.html", AA{Img: imgString, Symbols: symbols, Legend: legend, Xlen: len(symbols[0]), Ylen: len(symbols)}); err != nil { return err }	
	
  return nil	
}	
	
func WritePNG(img image.Image, path string) (string, error) {	
  // Write new image to png file	
  absPath, err := filepath.Abs(path)	
  absSplit := strings.Split(absPath, ".")	
	
  newPath := absSplit[0] + "-dmc.png"	
  place, err := os.Create(newPath)	
  if err != nil { return "", err }	
  defer place.Close()	
	
  err = png.Encode(place, img)	
  return newPath, err	
}
