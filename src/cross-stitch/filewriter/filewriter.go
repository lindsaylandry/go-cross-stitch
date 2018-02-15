package filewriter

import (
  "os"
  "image"
  "image/png"
  "html/template"
  "path/filepath"
  "strings"

  "cross-stitch/palette"
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

func WriteHTML(img image.Image, legend map[palette.Thread]int, symbols[][]int, path string) (error){
  // array of utf-8 decimal codes
  a := make([]int, 9983-9728)
  for i := range a {
    a[i] = 9728 + i
  }

  // struct to send to html
  type AA struct {
    Img image.Image
    Symbols [][]int
  }

  // funcs to use in html template
  fmap := template.FuncMap{
    "forloop": forloop,
  }

  // Write new image to png file
  absPath, err := filepath.Abs(path)
  absSplit := strings.Split(absPath, ".")

  // Write HTML instructions
  htmlPath := absSplit[0] + "-dmc.html"
  htmlFile, err := os.Create(htmlPath)
  if err != nil { return err }
  t := template.Must(template.New("instructions").Funcs(fmap).ParseFiles("../../templates/instructions.html"))
  if err := t.ExecuteTemplate(htmlFile, "instructions.html", AA{Img: img, Symbols: symbols}); err != nil { return err }

  return nil
}

func WritePNG(img image.Image, path string) (error) {
  // Write new image to png file
  absPath, err := filepath.Abs(path)
  absSplit := strings.Split(absPath, ".")

  newPath := absSplit[0] + "-dmc.png"
  place, err := os.Create(newPath)
  if err != nil { return err }
  defer place.Close()

  err = png.Encode(place, img)
  return err
}
