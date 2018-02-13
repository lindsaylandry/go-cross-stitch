package convert

import (
  "os"
  "io/ioutil"
  "bytes"
  "image"
  "image/png"
  "image/jpeg"
  "image/color"
  "math"
  "path/filepath"
  "strings"
  "text/template"

  "fmt"
  "cross-stitch/palette"
)

func open(filename string) (image.Image, error) {
  data, err := ioutil.ReadFile(filename)
  if err != nil { return nil, err }

  image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
  image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

  img, _, err := image.Decode(bytes.NewReader(data))
  return img, err
}

func Greyscale(img image.Image, outputLoc string) (*image.Gray, error) {
  bounds := img.Bounds()
  greyImg := image.NewGray(bounds)

  for x := bounds.Min.X; x < bounds.Dx(); x++ {
    for y := bounds.Min.Y; y < bounds.Dy(); y++ {
      pix := img.At(x, y)
      greyImg.Set(x, y, pix)
    }
  }

  place, err := os.Create(outputLoc)
  if err != nil { return greyImg, err }
  defer place.Close()

  err = png.Encode(place, greyImg)

  return greyImg, err
}

func DMC(path string, limit int) (image.Image, error) {
  t, err := palette.DMCPalette()
  if err != nil { panic(err) }

  img, err := open(path)
  if err != nil { panic(err) }

  dmcImg, legend := convertColors(img, t)

  fmt.Println(legend)
  fmt.Println("Colors: ", len(legend))

  // array of utf-8 decimal codes
  a := make([]int, 9983-9728)
  for i := range a {
    a[i] = 9728 + i
  }

  type AA struct {
    A []int
  }

  // Write new image to png file
  absPath, err := filepath.Abs(path)
  absSplit := strings.Split(absPath, ".")
  newPath := absSplit[0] + "-dmc.png"
  place, err := os.Create(newPath)
  if err != nil { return dmcImg, err }
  defer place.Close()

  err = png.Encode(place, dmcImg)

  // Write HTML instructions
  htmlPath := absSplit[0] + "-dmc.html"
  htmlFile, err := os.Create(htmlPath)
  if err != nil { return dmcImg, err }
  templates := template.Must(template.ParseFiles("../../templates/instructions.html"))
  if err := templates.Execute(htmlFile, AA{A: a}); err != nil { panic(err) }

  return dmcImg, nil
}

func convertColors(img image.Image, t []palette.Thread) (image.Image, map[palette.Thread]int) {
  legend := make(map[palette.Thread]int)
  bounds := img.Bounds()
  newImg := image.NewRGBA(bounds)

  for x := bounds.Min.X; x < bounds.Dx(); x++ {
    for y := bounds.Min.Y; y < bounds.Dy(); y++ {
      pixel := img.At(x, y)
      newImg.Set(x, y, pixel)
    }
  }

  for x := bounds.Min.X; x < bounds.Dx(); x++ {
    for y := bounds.Min.Y; y < bounds.Dy(); y++ {
      // Euclidean distance
      r32,g32,b32,a := newImg.At(x, y).RGBA()
      r, g, b := float64(uint8(r32)), float64(uint8(g32)), float64(uint8(b32))

      minLen := math.MaxFloat64
      minIndex := 0
      for c := 0; c < len(t); c++ {
        dist := 2*math.Pow((float64(t[c].R) - r), 2) + 4*math.Pow((float64(t[c].G) - g), 2) + 3*math.Pow((float64(t[c].B) - b), 2) + (float64(t[c].R) + r)/2*(math.Pow((float64(t[c].R) - r), 2)-math.Pow((float64(t[c].B) - b), 2))/256
        if dist < minLen {
          minLen = dist
          minIndex = c
        }
      }

      if _, ok := legend[t[minIndex]]; ok {
        legend[t[minIndex]] += 1
      } else {
        legend[t[minIndex]] = 1
      }
      newImg.Set(x, y, color.RGBA{t[minIndex].R, t[minIndex].G, t[minIndex].B, uint8(a)})
    }
  }

  return newImg, legend
}
