package convert

func quickSort(c []int) {
  if len(c) <= 1 { return }

  p := len(c)/2
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
  if piv+1 <= len(c)-1 {quickSort(c[piv+1:len(c)])}
}

func quickSortLegend(c []Legend) {
  if len(c) <= 1 { return }

  p := len(c)/2
  // switch pivot with end
  tmp := c[p]
  c[p] = c[len(c)-1]
  c[len(c)-1] = tmp

  piv := 0
  for i := 0; i < len(c); i++ {
    // switch c with pivot point if less than pivot
    if c[i].Thread.ID < c[len(c)-1].Thread.ID {
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
  if piv+1 <= len(c)-1 {quickSortLegend(c[piv+1:len(c)])}
}

func quickSortColors(c [][]uint8, index int) {
  if len(c) <= 1 { return }

  p := len(c)/2
  // switch pivot with end
  tmp := c[p]
  c[p] = c[len(c)-1]
  c[len(c)-1] = tmp

  piv := 0
  for i := 0; i < len(c); i++ {
    // switch c with pivot point if less than pivot
    if c[i][index] < c[len(c)-1][index] {
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

  quickSortColors(c[0:piv], index)
  if piv+1 <= len(c)-1 {quickSortColors(c[piv+1:len(c)], index)}
}
