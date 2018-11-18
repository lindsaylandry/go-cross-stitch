package main

import (
  "fmt"
  "os"
  "flag"

  "github.com/lindsaylandry/go-cross-stitch/src/convert"
)

//go:generate go run scripts/setupColors.go

func main() {
  if len(os.Args) < 2 {
    fmt.Println("No input image provided")
    os.Exit(0)
  }
 
  num := flag.Int("n", 6, "number of colors to use (2^n)")
  flag.Parse()

	c, err := convert.NewConverter(flag.Args()[0], *num)
	if err != nil { panic(err) }

  if err := c.DMC(); err != nil { panic(err) }
}
