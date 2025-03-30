package convert

import (
	"image/color"
	"log/slog"
)

type errorPix struct {
	r, g, b float64
}

// Problem: pixel error can exceed 0 and 255. Need to account for big errors.
func (c *Converter) floydSteinbergDither() {
	slog.Info("Dithering image...")
	bounds := c.image.Bounds()

	// create an error pixel matrix
	boundX := bounds.Dx() - bounds.Min.X
	boundY := bounds.Dy() - bounds.Min.Y

	errs := make([][]errorPix, boundY)
	for i := range errs {
		errs[i] = make([]errorPix, boundX)
	}

	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		for x := bounds.Min.X; x < bounds.Dx(); x++ {
			r32, g32, b32, _ := c.newData.Image.At(x, y).RGBA()
			errs[y][x] = errorPix{float64(uint8(r32)), float64(uint8(g32)), float64(uint8(b32))}
		}
	}

	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		for x := bounds.Min.X; x < bounds.Dx(); x++ {
			xx := x
			if y%2 != 0 {
				xx = bounds.Dx() - 1 - x
			}

			// clip to 0-255, set pixel at this site
			px := errs[y][xx]
			r1, g1, b1 := px.r, px.g, px.b

			if r1 < 0 {
				r1 = 0
			} else if r1 > 255 {
				r1 = 255
			}

			if g1 < 0 {
				g1 = 0
			} else if g1 > 255 {
				g1 = 255
			}

			if b1 < 0 {
				b1 = 0
			} else if b1 > 255 {
				b1 = 255
			}

			a := 255

			c.newData.Image.Set(xx, y, color.RGBA{uint8(r1), uint8(g1), uint8(b1), uint8(a)})

			// convert pixel to nearest palette color
			minIndex := c.setNewPixel(xx, y)
			c.newData.Count[c.pc[minIndex]] += 1

			// get new pixel values
			r32, g32, b32, _ := c.newData.Image.At(xx, y).RGBA()
			r2, g2, b2 := float64(uint8(r32)), float64(uint8(g32)), float64(uint8(b32))

			// quantization error
			qR := r1 - r2
			qG := g1 - g2
			qB := b1 - b2

			// x, y+1 (5/16)
			if y+1 < bounds.Dy() {
				setPixelError(&errs[y+1][xx], qR, qG, qB, 5.0/16.0)
			}

			// Odd: Backward
			if y%2 != 0 {
				// x-1, y (7/16)
				if xx-1 > bounds.Min.X {
					setPixelError(&errs[y][xx-1], qR, qG, qB, 7.0/16.0)
				}

				// x+1, y+1 (3/16)
				if xx+1 < bounds.Dx() && y+1 < bounds.Dy() {
					setPixelError(&errs[y+1][xx+1], qR, qG, qB, 3.0/16.0)
				}

				// x-1, y+1 (1/16)
				if xx-1 > bounds.Min.X && y+1 < bounds.Dy() {
					setPixelError(&errs[y+1][xx-1], qR, qG, qB, 1.0/16.0)
				}
				// Even: Forward
			} else {
				// x+1, y (7/16)
				if xx+1 < bounds.Dx() {
					setPixelError(&errs[y][xx+1], qR, qG, qB, 7.0/16.0)
				}

				// x-1, y+1 (3/16)
				if xx-1 > bounds.Min.X && y+1 < bounds.Dy() {
					setPixelError(&errs[y+1][xx-1], qR, qG, qB, 3.0/16.0)
				}

				// x+1, y+1 (1/16)
				if xx+1 < bounds.Dx() && y+1 < bounds.Dy() {
					setPixelError(&errs[y+1][xx+1], qR, qG, qB, 1.0/16.0)
				}
			}
		}
	}

	slog.Info("Done")
}

func setPixelError(e *errorPix, qR, qG, qB, diffusion float64) {
	e.r += qR * diffusion
	e.g += qG * diffusion
	e.b += qB * diffusion
}
