# Go Cross-Stitch Pattern Generator
This is a project that will take an image and convert it to a png and pdf of DMC thread colors and instructions.

Other palettes include Anchor cross-stitch threads, simple LEGO colors, and greyscale.

## How To Use

### Build From Scratch
To build the binary, run the following:

```bash
go build
```

### General Usage
Once the binary is compiled, use as follows:
```bash
./go-cross-stitch -n 10 test_images/full-moon.png
```
This will make four files in test_images:
```
full-moon-dmc-rgb.png
full-moon-dmc-rgb-A4.pdf
full-moon-dmc-rgb-A2.pdf
full-moon-dmc-rgb-A1.pdf
```
the png is the image converted to cross-stitch DMC thread colors.
the PDF is the instructions to stitch the pattern, with the DMC image included.

Run the help command to see all flags available:
```bash
./go-cross-stitch --help
```

### Render all test images
```bash
make examples
```

## References
Color distance formulas: https://en.wikipedia.org/wiki/Color_difference

Color quantization: https://en.wikipedia.org/wiki/Color_quantization

CIELab color space: https://en.wikipedia.org/wiki/CIELAB_color_space

## Example Image Conversions

### Mars (reds)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/mars.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/mars-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/mars-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Earth (blues and greens)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/earth.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/earth-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/earth-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Moon (greyscale)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/full-moon.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/full-moon-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/full-moon-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Full Color Spectrum
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/colors.jpg" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/colors-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/colors-dmc-lab.png" height="200" style="image-rendering: pixelated;">
