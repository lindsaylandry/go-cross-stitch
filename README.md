# Go Cross-Stitch Pattern Generator
Do you love GoLang? Do you also love cross-stitching? Congratulations, you are a very rare cross-section of people!

This is a project that will take an image and convert it to a png and html of DMC thread colors and instructions.

## How To Use

### Build From Scratch
To build the binary, run the following:

```make build```

### General Usage
Once the binary is compiled, use as follows:
```
./bin/cross-stitch -n 10 examples/test_images/FullMoon150px.jpg
```
This will make two files in examples/test_images:
```
FullMoon150px-dmc-lab.png
FullMoon150px-dmc.html
```
the png is the image converted to cross-stitch DMC thread colors.
the HTML is the instructions to stitch the pattern, with the DMC image included.

### Flags
```
./bin/cross-stitch --help
Usage of ./bin/cross-stitch:
  -all
    	use all thread colors available (currently broken)
  -color string
    	color palette to use (OPTIONS: dmc, lego, bw) (default "dmc")
  -d	implement dithering
  -g	convert image to greyscale
  -n int
    	number of colors to attempt to match (2^n) (default 6)
  -rgb
    	use rgb color space
```

### Creating a PDF from HTML
For now, conversion from HTML to PDF is a manual process.

To convert to PDF, do the following:
* Open the generated HTML in Firefox. 
* Click `File => Print`
* Under `Appearance` check `Print Background Colors`
* Click the `PDF` dropdown and click `Save as PDF`

NOTE: it will take a while for the Firefox browser to convert the HTML to PDF. Do not close the page until the PDF is fully rendered.

TODO: auto-create PDF from HTML; find out what libraries Firefox is using.

### Render all test images
```make examples```

## References
Color distance formulas: https://en.wikipedia.org/wiki/Color_difference

Color quantization: https://en.wikipedia.org/wiki/Color_quantization

CIELab color space: https://en.wikipedia.org/wiki/CIELAB_color_space

## Example Image Conversions

### Mars (reds)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/mars.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/mars-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/mars-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Earth (blues and greens)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Moon (greyscale)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/full-moon.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/full-moon-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/full-moon-dmc-lab.png" height="200" style="image-rendering: pixelated;">

### Full Color Spectrum
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/colors.jpg" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/colors-dmc-rgb.png" height="200" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/colors-dmc-lab.png" height="200" style="image-rendering: pixelated;">


