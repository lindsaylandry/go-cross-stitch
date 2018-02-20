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

  "fmt"
  "cross-stitch/palette"
  "cross-stitch/filewriter"
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

func DMC(path string, limit int) (error) {
  t, err := palette.DMCPalette()
  if err != nil { panic(err) }

  img, err := open(path)
  if err != nil { panic(err) }

  //bcrgb := colorQuant(img)

  // Convert best-colors to thread palette
  //bt := convertPalette(bcrgb, t)
  //fmt.Println(bt)

  dmcImg, legend, symbols := convertImage(img, t)

  _, nerr := filewriter.WritePNG(dmcImg, path)
  if nerr != nil { panic(nerr) }
  err = filewriter.WriteHTML(dmcImg, legend, symbols, path)
  if err != nil { panic(err) }

  return nil
}

// convert best-color palette to match available threads
func convertPalette(colors [][]uint8, t []palette.Thread) (map[palette.Thread]int) {
  legend := make(map[palette.Thread]int)
  for i := 0; i < len(colors); i++ {
    minLen := math.MaxFloat64
    minIndex := 0

    for c := 0; c < len(t); c++ {
      dist := 2*math.Pow(float64(t[c].R - colors[i][0]), 2) + 4*math.Pow(float64(t[c].G - colors[i][1]), 2) + 3*math.Pow(float64(t[c].B - colors[i][2]), 2) + float64(t[c].R + colors[i][0])/2*(math.Pow(float64(t[c].R - colors[i][0]), 2)-math.Pow(float64(t[c].B - colors[i][2]), 2))/256
      if dist < minLen {
        minLen = dist
        minIndex = c
      }
    }

    if _, ok := legend[t[minIndex]]; !ok {
      legend[t[minIndex]] = 0
    }
  }

  return legend
}

func convertImage(img image.Image, t []palette.Thread) (image.Image, map[palette.Thread]int, [][]int) {
  legend := make(map[palette.Thread]int)
  bounds := img.Bounds()
  newImg := image.NewRGBA(bounds)

  for y := bounds.Min.Y; y < bounds.Dy(); y++ {
    for x := bounds.Min.X; x < bounds.Dx(); x++ {
      pixel := img.At(x, y)
      newImg.Set(x, y, pixel)
    }
  }

  symbols := make([][]int, bounds.Dy() - bounds.Min.Y)
  for y := bounds.Min.Y; y < bounds.Dy(); y++ {
    symbols[y] = make([]int, bounds.Dx() - bounds.Min.X)
    for x := bounds.Min.X; x < bounds.Dx(); x++ {
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

      symbols[y][x] = t[minIndex].UTFDec
      newImg.Set(x, y, color.RGBA{t[minIndex].R, t[minIndex].G, t[minIndex].B, uint8(a)})
    }
  }

  return newImg, legend, symbols
}

// start with 16 colors
func colorQuant (img image.Image) ([][]uint8) {
  bounds := img.Bounds()
  newImg := image.NewRGBA(bounds)

  for y := bounds.Min.Y; y < bounds.Dy(); y++ {
    for x := bounds.Min.X; x < bounds.Dx(); x++ {
      pixel := img.At(x, y)
      newImg.Set(x, y, pixel)
    }
  }

  allcolors := make([][]uint8, (bounds.Dy() - bounds.Min.Y)*(bounds.Dx() - bounds.Min.X))
  for y := bounds.Min.Y; y < bounds.Dy(); y++ {
    for x := bounds.Min.X; x < bounds.Dx(); x++ {
      r32,g32,b32,a := newImg.At(x, y).RGBA()
      allcolors[y*(bounds.Dx() - bounds.Min.X)+x] = []uint8{uint8(r32), uint8(g32), uint8(b32), uint8(a)}
    }
  }

  // 1 2 4 8 16
  slices := []int {0, len(allcolors)-1}

  for i := 0; i < 4; i++ {
    for j := 0; j < len(slices)-1; j++ {
      // get a slice of allcolors
      s := allcolors[slices[j]:slices[j+1]]

      colorRanges := [][]uint8{{math.MaxUint8, 0}, {math.MaxUint8, 0},{math.MaxUint8, 0}}
      for c := 0; c < len(s); c++ {
        //R
        if s[c][0] < colorRanges[0][0] { colorRanges[0][0] = s[c][0] }
        if s[c][0] > colorRanges[0][1] { colorRanges[0][1] = s[c][0] }
        //G
        if s[c][1] < colorRanges[1][0] { colorRanges[1][0] = s[c][1] }
        if s[c][1] > colorRanges[1][1] { colorRanges[1][1] = s[c][1] }
        //B
        if s[c][2] < colorRanges[2][0] { colorRanges[2][0] = s[c][2] }
        if s[c][2] > colorRanges[2][1] { colorRanges[2][1] = s[c][2] }
      }
      var index int
      xr := colorRanges[0][1] - colorRanges[0][0]
      yr := colorRanges[1][1] - colorRanges[1][0]
      zr := colorRanges[2][1] - colorRanges[2][0]

      if xr > yr && xr > zr { 
        index = 0
      } else if yr > xr && yr > zr { 
        index = 1 
      } else { index = 2 }
      quickSortColors(s, index)
    }

    // insert 2^n more slice indexes
    max := len(slices) - 1
    for k := 0; k < max; k++ {
      slices = append(slices, (slices[k+1]+slices[k])/2)
    }
    quickSort(slices)
  }

  // Average all sliced colors and insert into bestcolors
  bestcolors := make([][]uint8, len(slices)-1)
  for i := 0; i < len(slices)-1; i++ {
    s := allcolors[slices[i]:slices[i+1]]
    avgR := float64(0)
    avgG := float64(0)
    avgB := float64(0)
    for c := 0; c < len(s); c++ {
      avgR = avgR + math.Pow(float64(s[c][0]), 2)
      avgG = avgG + math.Pow(float64(s[c][1]), 2)
      avgB = avgB + math.Pow(float64(s[c][2]), 2)
    }

    bestcolors[i] = []uint8 {uint8(math.Sqrt(avgR/float64(len(s)))), uint8(math.Sqrt(avgG/float64(len(s)))), uint8(math.Sqrt(avgB/float64(len(s))))}
  }

  fmt.Println(bestcolors)
  return bestcolors
}
