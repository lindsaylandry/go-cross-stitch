package convert

import (
  "os"
  "image"
  "image/png"
  "image/jpeg"
  "image/color"
  "encoding/csv"
  "strconv"
  "math"

  //"fmt"
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
  file, err := os.Open("../../palette/black-white-grey.csv")

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
      r32,g32,b32,a := img.At(x, y).RGBA()
      r, g, b := float64(uint8(r32)), float64(uint8(g32)), float64(uint8(b32))
      //fmt.Println(r, " ", g, " ", b)

      min := [2]float64 {0, math.MaxFloat64}
      for c := 0; c < len(record); c++ {
        rr, _ := strconv.Atoi(record[c][2])
        rg, _ := strconv.Atoi(record[c][3])
        rb, _ := strconv.Atoi(record[c][4])
        dist := math.Pow((float64(rr) - r), 2) + math.Pow((float64(rg) - g), 2) + math.Pow((float64(rb) - b), 2)
        if dist < min[1] {
          min[1] = dist
          min[0] = float64(c)
        }
      }

      dmcImg.Set(x, y, color.RGBA{uint8(strconv.Atoi(record[int(min[0])][2])),uint8(strconv.Atoi(record[int(min[0])][3])),uint8(strconv.Atoi(record[int(min[0])][4])),uint8(a)})
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
