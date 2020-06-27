PACKAGES=$$(go list ./... | grep -v '/vendor/' | grep -v '/mocks')

build:
	go build -o bin/cross-stitch
test:
	go test -cover $(PACKAGES)
cover:
	go test -coverprofile cover.out $(PACKAGES) && go tool cover -html=cover.out -o cover.html && open cover.html
fmt:
	go fmt ./...
deps:
	dep ensure -v
clean:
	go clean
	rm -f bin/
examples: build
	./bin/cross-stitch -rgb=true -n=10 examples/test_images/full-moon.png
	./bin/cross-stitch -rgb=false -n=10 examples/test_images/full-moon.png
	./bin/cross-stitch -rgb=false -n=10 examples/test_images/colors.jpg
	./bin/cross-stitch -rgb=true -n=10 examples/test_images/colors.jpg
	./bin/cross-stitch -rgb=true -n=10 examples/test_images/mars.png
	./bin/cross-stitch -rgb=false -n=10 examples/test_images/mars.png
	./bin/cross-stitch -rgb=true -n=10 examples/test_images/earth.png
	./bin/cross-stitch -rgb=false -n=10 examples/test_images/earth.png
