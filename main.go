package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lindsaylandry/go-cross-stitch/src/convert"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No input image provided")
		os.Exit(0)
	}

	num := flag.Int("n", 10, "number of colors to attempt to match (2^n)")
	rgb := flag.Bool("rgb", false, "use rgb color space (default false)")
	all := flag.Bool("all", false, "use all thread colors available (currently broken)")
	pal := flag.String("color", "dmc", "color palette to use (OPTIONS: dmc, anchor, lego, bw)")
	dit := flag.Bool("d", false, "implement dithering (default false)")
	gre := flag.Bool("g", false, "convert image to greyscale (default false)")
	pix := flag.Bool("px", false, "quantize pixellated image (default false)")
	col := flag.Bool("c", true, "include color grid instructions")
	flag.Parse()

	c, err := convert.NewConverter(flag.Args()[0], *num, *rgb, *all, *pal, *dit, *gre, *pix, *col)
	if err != nil {
		panic(err)
	}

	if err := c.Convert(); err != nil {
		panic(err)
	}
}
