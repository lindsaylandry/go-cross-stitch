# Go Cross Stitch Pattern Generator
Do you love GoLang? Do you also love cross-stitching? Congratulations, you are a very rare cross-section of people!

This is a project that will take an image and convert it to a png and html of DMC thread colors and instructions.
[moon] [moon-dmc]
[moon]: https://github.com/lindsaylandry/go-cross-stitch/examplse/test_images/FullMoon150px.jpg
[moon-dmc]: https://github.com/lindsaylandry/go-cross-stitch/examplse/test_images/FullMoon150px-dmc.png

## Build From Scratch
Currently go get does not work. Docker is your best bet.

### Using Docker
change GOOS and GOARCH to your preferred architecture
```
docker run --rm -i -t -v `pwd`:/go -e GOOS=darwin -e GOARCH=amd64 golang:1.9-alpine go install cross-stitch
```
