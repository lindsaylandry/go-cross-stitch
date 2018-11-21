# Go Cross-Stitch Pattern Generator
Do you love GoLang? Do you also love cross-stitching? Congratulations, you are a very rare cross-section of people!

This is a project that will take an image and convert it to a png and html of DMC thread colors and instructions.

## Example Image Conversions

### Mars (reds)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/mars.jpg" height="250" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/mars-dmc-rgb.png" height="250" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/mars-dmc-lab.png" height="250" style="image-rendering: pixelated;">

### Earth (blues and greens)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth200.jpg" height="250" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth200-dmc-rgb.png" height="250"style=" image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth200-dmc-lab.png" height="250" style="image-rendering: pixelated;">

### Pluto (greyscale/reds)
| Original | RGB Distance | CIELab Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/pluto150.jpg" height="250" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/pluto150-dmc-rgb.png" height="250" style="image-rendering: pixelated;"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/pluto150-dmc-lab.png" height="250" style="image-rendering: pixelated;">

## Build From Scratch
```make build```

## Usage
Once the binary is compiled, use as follows:
```
./cross-stitch -n 10 ~/Pictures/FullMoon150px.jpg
```
This will make two files in ~/Pictures
```
~/Pictures/FullMoon150px-dmc.png
~/Pictures/FullMoon150px-dmc.html
```
the png is the image converted to cross-stitch DMC thread colors.
the HTML is the instructions to stitch the pattern, with the DMC image included.

TODO: Currently the html instructions are not printer friendly. Will work on this.
