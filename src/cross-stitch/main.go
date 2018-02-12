package main

import (
  "fmt"
  "os"
  "flag"

  "cross-stitch/convert"
)

//go:generate go run scripts/setupColors.go

func main() {
  if len(os.Args) < 2 {
    fmt.Println("No input image provided")
    os.Exit(0)
  }
  args := os.Args[1:]

  num := flag.Int("ncolor", 500, "number of colors to use")
  flag.Parse()

  _, err := convert.DMC(args[0], *num)
  if err != nil { panic(err) }
}
