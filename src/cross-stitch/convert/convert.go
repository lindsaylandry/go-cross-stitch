package convert

import (
  "os"
  "image"
  "image/png"
  "image/jpeg"
  "image/color"
  "encoding/csv"
  "strconv"
  "fmt"
)

type Dither struct {
  Filter [][] float32
}

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

  // TODO: Many color palettes in the future?
  //file, err := os.Open("palette/dmc-floss.csv")
  file, err := os.Open("palette/black-white-grey.csv")

  if err != nil {
     return nil, err
  }
  defer file.Close()

  // convert dmc data to csv hash
  reader := csv.NewReader(file)
  reader.Comma = ','
  record, err := reader.ReadAll()
  if err != nil {
    return nil, err
  }

  // Record: [["Floss#","Description","Red","Green","Blue"],...]

  bounds := img.Bounds()
  dmcImg := image.NewRGBA(bounds)

  for x := bounds.Min.X; x < bounds.Dx(); x++ {
    for y := bounds.Min.Y; y < bounds.Dy(); y++ {
      // Euclidean distance
      r,g,b,a := img.At(x, y).RGBA()

      dist := make([]int, len(record))
      for c := 0; c < len(record); c++ {
        rr, _ := strconv.Atoi(record[c][2])
        rg, _ := strconv.Atoi(record[c][3])
        rb, _ := strconv.Atoi(record[c][4])
        dist[c] = (rr - int(r))^2 + (rg - int(g))^2 + (rb - int(b))^2
      }

      fmt.Println(dist)

      dmcImg.Set(x, y, color.RGBA{uint8(r),uint8(g),uint8(b),uint8(a)})
    }
  }

  return dmcImg, nil
}
