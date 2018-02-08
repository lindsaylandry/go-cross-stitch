package main

import (
  "fmt"
  "os"

  "cross-stitch/dither"
)

func main() {
  if len(os.Args) < 2 {
    fmt.Println("Hello world!")
    os.Exit(0)
  }
  args := os.Args[1:]

  img, err := dither.Open(args[0])
  if err != nil {
    panic(err)
  }
  
  grey, err := dither.Greyscale(img, "output.png")
  if err != nil {
    panic(err)
  }
  fmt.Println(grey)
}
