package convert

import (
	"github.com/lindsaylandry/go-cross-stitch/src/colorConverter"
)

func quickSort(c []int) {
	if len(c) <= 1 {
		return
	}

	p := len(c) / 2
	// switch pivot with end
	tmp := c[p]
	c[p] = c[len(c)-1]
	c[len(c)-1] = tmp

	piv := 0
	for i := 0; i < len(c); i++ {
		// switch c with pivot point if less than pivot
		if c[i] < c[len(c)-1] {
			tmp = c[i]
			c[i] = c[piv]
			c[piv] = tmp

			piv++
		}
	}

	// put pivot point back where it was
	tmp = c[piv]
	c[piv] = c[len(c)-1]
	c[len(c)-1] = tmp

	quickSort(c[0:piv])
	if piv+1 <= len(c)-1 {
		quickSort(c[piv+1:])
	}
}

func quickSortLegend(c []Legend) {
	if len(c) <= 1 {
		return
	}

	p := len(c) / 2
	// switch pivot with end
	tmp := c[p]
	c[p] = c[len(c)-1]
	c[len(c)-1] = tmp

	piv := 0
	for i := 0; i < len(c); i++ {
		// switch c with pivot point if less than pivot
		if c[i].Color.ID < c[len(c)-1].Color.ID {
			tmp = c[i]
			c[i] = c[piv]
			c[piv] = tmp

			piv++
		}
	}

	// put pivot point back where it was
	tmp = c[piv]
	c[piv] = c[len(c)-1]
	c[len(c)-1] = tmp

	quickSortLegend(c[0:piv])
	if piv+1 <= len(c)-1 {
		quickSortLegend(c[piv+1:])
	}
}

func quickSortRed(c []colorConverter.SRGB) {
	if len(c) <= 1 {
		return
	}

	p := len(c) / 2
	// switch pivot with end
	tmp := c[p]
	c[p] = c[len(c)-1]
	c[len(c)-1] = tmp

	piv := 0
	for i := 0; i < len(c); i++ {
		// switch c with pivot point if less than pivot
		if c[i].R < c[len(c)-1].R {
			tmp = c[i]
			c[i] = c[piv]
			c[piv] = tmp

			piv++
		}
	}

	// put pivot point back where it was
	tmp = c[piv]
	c[piv] = c[len(c)-1]
	c[len(c)-1] = tmp

	quickSortRed(c[0:piv])
	if piv+1 <= len(c)-1 {
		quickSortRed(c[piv+1:])
	}
}

func quickSortGreen(c []colorConverter.SRGB) {
	if len(c) <= 1 {
		return
	}

	p := len(c) / 2
	// switch pivot with end
	tmp := c[p]
	c[p] = c[len(c)-1]
	c[len(c)-1] = tmp

	piv := 0
	for i := 0; i < len(c); i++ {
		// switch c with pivot point if less than pivot
		if c[i].G < c[len(c)-1].G {
			tmp = c[i]
			c[i] = c[piv]
			c[piv] = tmp

			piv++
		}
	}

	// put pivot point back where it was
	tmp = c[piv]
	c[piv] = c[len(c)-1]
	c[len(c)-1] = tmp

	quickSortGreen(c[0:piv])
	if piv+1 <= len(c)-1 {
		quickSortGreen(c[piv+1:])
	}
}

func quickSortBlue(c []colorConverter.SRGB) {
	if len(c) <= 1 {
		return
	}

	p := len(c) / 2
	// switch pivot with end
	tmp := c[p]
	c[p] = c[len(c)-1]
	c[len(c)-1] = tmp

	piv := 0
	for i := 0; i < len(c); i++ {
		// switch c with pivot point if less than pivot
		if c[i].B < c[len(c)-1].B {
			tmp = c[i]
			c[i] = c[piv]
			c[piv] = tmp

			piv++
		}
	}

	// put pivot point back where it was
	tmp = c[piv]
	c[piv] = c[len(c)-1]
	c[len(c)-1] = tmp

	quickSortBlue(c[0:piv])
	if piv+1 <= len(c)-1 {
		quickSortBlue(c[piv+1:])
	}
}
