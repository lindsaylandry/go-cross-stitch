# Go Cross-Stitch Pattern Generator
This is a project that will take an image and convert it to a png and pdf of DMC thread colors and instructions.

Other palettes include Anchor cross-stitch threads, simple LEGO colors, and greyscale.

## How To Use

### Build From Scratch
To build the binary, run the following:

```bash
go build
```

If this is your first time building, first run the following to download dependencies:

```bash
go mod tidy
```

### General Usage
Once the binary is compiled, use as follows:
```bash
./go-cross-stitch test_images/full-moon.png
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

### Configuration

Configure this app with file `configs/config.yaml`

Refer to the table below for accepted values.

| Key | Type | Default | Definition |
| - | - | - | - |
| `csv` | `string` | `""` | csv filename |
| `dither` | `bool` | `false` | implement dithering |
| `greyscale` | `bool` | `false` | convert image to greyscale |
| `palette` | `string` | `dmc` | color palette to use (OPTIONS: dmc, anchor, lego, bw) |
| `quantize` | `struct` | `-` | settings to quantize image (see below) |
| `rgb` | `bool` | `true` | use rgb color space |
| `width` | `int` | `300` | resize image width (0 means do not resize) |

#### quantize

Color quantization reduces the number of colors to match. See https://en.wikipedia.org/wiki/Color_quantization for more details.

If disabled, the app will attempt to match all colors to the specified color palette.

Enable this setting to reduce the number of colors in your instructions.

| Key | Type | Default | Definition |
| - | - | - | - |
| enabled | boolean | true | whether to enable color quantization |
| n | int | 7 | number of bisects for color quantization (2^n total) |

## Example Image Conversions

### Render all test images
```bash
make examples
```

### Mars (reds)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/mars.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/mars-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/mars-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Earth (blues and greens)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/earth-americas.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/earth-americas-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/earth-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Moon (greyscale)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/full-moon.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/full-moon-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/full-moon-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Full Color Spectrum
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/colors.jpg" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/colors-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/main/test_images/colors-dmc-lab.png" height="200" style="image-rendering: pixelated;">

## References
Color distance formulas: https://en.wikipedia.org/wiki/Color_difference

Color quantization: https://en.wikipedia.org/wiki/Color_quantization

CIELab color space: https://en.wikipedia.org/wiki/CIELAB_color_space


