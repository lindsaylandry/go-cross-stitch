# Go Cross Stitch Pattern Generator
Right now this code converts an input image to greyscale

# Build From Scratch

change GOOS and GOARCH to your preferred architecture
```
docker run --rm -i -t -v `pwd`:/go -e GOOS=darwin -e GOARCH=amd64 golang:1.9-alpine go install cross-stitch
```
