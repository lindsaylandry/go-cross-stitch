# Go Cross-Stitch Pattern Generator
Do you love GoLang? Do you also love cross-stitching? Congratulations, you are a very rare cross-section of people!

This is a project that will take an image and convert it to a png and html of DMC thread colors and instructions.
![alt text][moon] 
![alt_text][moon-dmc]

## Example Image Conversions

### Moon (greyscale)

[moon]: https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/FullMoon150px.jpg
[moon-dmc]: https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/FullMoon150px-dmc.png

### Mars (reds)

### Earth (blues and greens)
| Original | RGB Distance | CIELAb Distance |
|:--:|:--:|:--:|
| <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth200.jpg" height="250"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth200-dmc-rgb.png" height="250"> | <img src="https://github.com/lindsaylandry/go-cross-stitch/blob/master/examples/test_images/earth200-dmc-lab.png" height="250">

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
