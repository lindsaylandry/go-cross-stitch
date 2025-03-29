package convert

import (
	"fmt"
	"image"

	"golang.org/x/image/draw"
)

func resize(src image.Image, width int) (image.Image) {
	// TODO: do math to get image dimensions
	bounds := src.Bounds()

	dstHeight := int(float64(bounds.Max.X) / float64(bounds.Max.Y) * float64(width))

	fmt.Printf("Resizing image to %dx%d pixels... ", width, dstHeight)

	dst := image.NewRGBA(image.Rect(0, 0, width, dstHeight))
	draw.NearestNeighbor.Scale(dst, image.Rect(0, 0, width, dstHeight), src, image.Rect(0, 0, bounds.Max.X, bounds.Max.Y), draw.Over, nil)

	fmt.Println("Done")

	return dst
}
