package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lindsaylandry/go-cross-stitch/src/convert"
	"github.com/lindsaylandry/go-cross-stitch/src/writer"
)

var flags convert.Flags

func main() {
	rootCmd := &cobra.Command{
    Use:   "cross-stitch",
    Short: "Generate cross-stitch pattern",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				panic(errors.New("No input image provided"))
			}

			return CrossStitch(args[0])
    },
  }

	rootCmd.PersistentFlags().IntVarP(&flags.Num, "number", "n", 10, "number of colors to attempt to match (2^n)")
	rootCmd.PersistentFlags().BoolVarP(&flags.RGB, "rgb", "r", true, "use rgb color space")
	rootCmd.PersistentFlags().BoolVarP(&flags.All, "all", "a", false, "use all thread colors available")
	rootCmd.PersistentFlags().StringVarP(&flags.Palette, "pal", "p", "dmc", "color palette to use (OPTIONS: dmc, anchor, lego, bw)")
	rootCmd.PersistentFlags().BoolVarP(&flags.Dither, "dither", "d", false, "implement dithering")
	rootCmd.PersistentFlags().BoolVarP(&flags.Greyscale, "greyscale", "g", false, "convert image to greyscale")
	rootCmd.PersistentFlags().BoolVarP(&flags.Pixel, "px", "x", true, "quantize pixellated image")
	rootCmd.PersistentFlags().BoolVarP(&flags.Color, "colorgrid", "c", true, "include color grid instructions")
	rootCmd.PersistentFlags().StringVarP(&flags.CSV, "csv", "s", "", "csv filename (optional)")

	err := rootCmd.Execute()
  if err != nil {
    panic(err)
  }
}

func CrossStitch(filename string) error {
	fmt.Println(flags.CSV)

	c, err := convert.NewConverter(filename, flags)
	if err != nil {
		return err
	}

	d, err := c.Convert()
	if err != nil {
		return err
	}

	w := writer.NewWriter(d)
	
	return w.WriteFiles()
}
