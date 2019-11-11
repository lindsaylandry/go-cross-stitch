package convert

import (
	"image/color"
	//"fmt"
)

func (c *Converter) floydSteinbergDither() {
	bounds := c.image.Bounds()
	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		for x := bounds.Min.X; x < bounds.Dx(); x++ {
			xx := x
			if y % 2 != 0 {
				xx = bounds.Dx() - 1 - x
			}

			r32, g32, b32, _ := c.newImage.image.At(xx, y).RGBA()
			r1, g1, b1 := float64(uint8(r32)), float64(uint8(g32)), float64(uint8(b32))

			// x, y
			minIndex := c.setNewPixel(xx, y)
			if _, ok := c.newImage.count[c.pc[minIndex]]; ok {
				c.newImage.count[c.pc[minIndex]] += 1
			} else {
				c.newImage.count[c.pc[minIndex]] = 1
			}

			r232, g232, b232, _ := c.newImage.image.At(xx, y).RGBA()
			r2, g2, b2 := float64(uint8(r232)), float64(uint8(g232)), float64(uint8(b232))

			qR := r1 - r2
			qG := g1 - g2
			qB := b1 - b2

			// Odd
			if y % 2 != 0 {
				// x-1, y (7/16)
				if xx-1 > bounds.Min.X {
					c.setPixelError(xx-1, y, qR, qG, qB, 7.0/16.0)
				}

				// x+1, y+1 (3/16)
				if xx+1 < bounds.Dx() && y+1 < bounds.Dy() {
					c.setPixelError(xx+1, y+1, qR, qG, qB, 3.0/16.0)
				}

				// x, y+1 (5/16)
				if y+1 < bounds.Dy() {
					c.setPixelError(xx, y+1, qR, qG, qB, 5.0/16.0)
				}

				// x-1, y+1 (1/16)
				if xx-1 > bounds.Min.X && y+1 < bounds.Dy() {
					c.setPixelError(xx-1, y+1, qR, qG, qB, 1.0/16.0)
				}
			// Even
			} else {
				// x+1, y (7/16)
				if xx+1 < bounds.Dx() {
					c.setPixelError(xx+1, y, qR, qG, qB, 7.0/16.0)
				}

				// x-1, y+1 (3/16)
				if xx-1 > bounds.Min.X && y+1 < bounds.Dy() {
					c.setPixelError(xx-1, y+1, qR, qG, qB, 3.0/16.0)
				}

				// x, y+1 (5/16)
				if y+1 < bounds.Dy() {
					c.setPixelError(xx, y+1, qR, qG, qB, 5.0/16.0)
				}

				// x+1, y+1 (1/16)
				if xx+1 < bounds.Dx() && y+1 < bounds.Dy() {
					c.setPixelError(xx+1, y+1, qR, qG, qB, 1.0/16.0)
				}
			}
		}
	}
}

func (c *Converter) setPixelError(x, y int, qR, qG, qB, diffusion float64) {
	rr32, gg32, bb32, aa := c.newImage.image.At(x, y).RGBA()
	rr, gg, bb := uint8(rr32), uint8(gg32), uint8(bb32)
	rd := uint8(float64(rr) + qR*diffusion)
	gd := uint8(float64(gg) + qG*diffusion)
	bd := uint8(float64(bb) + qB*diffusion)
	c.newImage.image.Set(x, y, color.RGBA{rd, gd, bd, uint8(aa)})
}
