package convert

import (
	"fmt"
	"image/color"
)

func (c *Converter) floydSteinbergDither() {
	bounds := c.image.Bounds()
	for y := bounds.Min.Y; y < bounds.Dy(); y++ {
		for x := bounds.Min.X; x < bounds.Dx(); x++ {
			r32, g32, b32, _ := c.newImage.image.At(x, y).RGBA()
			r1, g1, b1 := uint8(r32), uint8(g32), uint8(b32)

			// x, y
			minIndex := c.setNewPixel(x, y)
			if _, ok := c.newImage.count[c.pc[minIndex]]; ok {
				c.newImage.count[c.pc[minIndex]] += 1
			} else {
				c.newImage.count[c.pc[minIndex]] = 1
			}

			r232, g232, b232, _ := c.newImage.image.At(x, y).RGBA()
			r2, g2, b2 := uint8(r232), uint8(g232), uint8(b232)

			qR := float64(int8(r1) - int8(r2))
			qG := float64(int8(g1) - int8(g2))
			qB := float64(int8(b1) - int8(b2))

			fmt.Printf("%f %f %f\n", qR, qG, qB)

			// x+1, y (7/16)
			if x+1 < bounds.Dx() {
				c.setPixelError(x+1, y, qR, qG, qB, 7.0/16.0)
			}

			// x-1, y+1 (3/16)
			if x-1 > bounds.Min.X && y+1 < bounds.Dy() {
				c.setPixelError(x-1, y+1, qR, qG, qB, 3.0/16.0)
			}

			// x, y+1 (5/16)
			if y+1 < bounds.Dy() {
				c.setPixelError(x, y+1, qR, qG, qB, 5.0/16.0)
			}

			// x+1, y+1 (1/16)
			if x+1 < bounds.Dx() && y+1 < bounds.Dy() {
				c.setPixelError(x+1, y+1, qR, qG, qB, 1.0/16.0)
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
