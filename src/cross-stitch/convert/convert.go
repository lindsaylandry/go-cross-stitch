package convert

import (
  "os"
  "image"
  "image/png"
  "image/jpeg"
  "image/color"
  "math"

  //"fmt"
  "cross-stitch/palette"
)

func Open(filename string) (image.Image, error) {
  file, err := os.Open(filename)
  if err != nil {
     return nil, err
  }
  defer file.Close()

  image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
  image.RegisterFormat("jpg", "jpg", jpeg.Decode, jpeg.DecodeConfig)

  img, _, err := image.Decode(file)
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
  if err != nil {
    return greyImg, err
  }
  defer place.Close()

  err = png.Encode(place, greyImg)

  return greyImg, err
}

func DMC(img image.Image) (image.Image, error) {
  t, err := palette.DMCPalette()
  if err != nil { panic(err) }

  bounds := img.Bounds()
  dmcImg := image.NewRGBA(bounds)

  for x := bounds.Min.X; x < bounds.Dx(); x++ {
    for y := bounds.Min.Y; y < bounds.Dy(); y++ {
      // Euclidean distance
      r32,g32,b32,a := img.At(x, y).RGBA()
      r, g, b := float64(uint8(r32)), float64(uint8(g32)), float64(uint8(b32))
      //fmt.Println(r, " ", g, " ", b)

      minLen := math.MaxFloat64
      minIndex := 0
      for c := 0; c < len(t); c++ {
        dist := math.Pow((float64(t[c].R) - r), 2) + math.Pow((float64(t[c].G) - g), 2) + math.Pow((float64(t[c].B) - b), 2)
        if dist < minLen {
          minLen = dist
          minIndex = c
        }
      }

      dmcImg.Set(x, y, color.RGBA{t[minIndex].R, t[minIndex].G, t[minIndex].B, uint8(a)})
    }
  }

  place, err := os.Create("dmcoutput.png")
  if err != nil {
    return dmcImg, err
  }
  defer place.Close()

  err = png.Encode(place, dmcImg)

  return dmcImg, nil
}
