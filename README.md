# Go Cross Stitch Pattern Generator
Do you love GoLang? Do you also love cross-stitching? Congratulations, you are a very rare cross-section of people!

This is a project that will take an image and convert it to a png and html of DMC thread colors and instructions.
![alt text][moon] 
![alt_text][moon-dmc]

[moon]: https://github.com/lindsaylandry/go-cross-stitch/examples/test_images/FullMoon150px.jpg
[moon-dmc]: https://github.com/lindsaylandry/go-cross-stitch/examples/test_images/FullMoon150px-dmc.png

## Build From Scratch
Currently go get does not work. Docker is your best bet.

### Using Docker
change GOOS and GOARCH to your preferred architecture
```
docker run --rm -i -t -v `pwd`:/go -e GOOS=darwin -e GOARCH=amd64 golang:1.9-alpine go install cross-stitch
```

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
