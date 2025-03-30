package convert

import (
	"fmt"
	"log/slog"
	"image"

	"golang.org/x/image/draw"
)

func resize(src image.Image, width int) image.Image {
	// TODO: do math to get image dimensions
	bounds := src.Bounds()

	dstHeight := int(float64(bounds.Max.Y) / float64(bounds.Max.X) * float64(width))

	slog.Info(fmt.Sprintf("Resizing image to %dx%d pixels... ", width, dstHeight))

	dst := image.NewRGBA(image.Rect(0, 0, width, dstHeight))
	draw.CatmullRom.Scale(dst, image.Rect(0, 0, width, dstHeight), src, image.Rect(0, 0, bounds.Max.X, bounds.Max.Y), draw.Over, nil)

	slog.Info("Done")

	return dst
}
