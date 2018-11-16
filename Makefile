PACKAGES=$$(go list ./... | grep -v '/vendor/' | grep -v '/mocks')

build:
	go build -o bin/cross-stitch
test:
	go test -cover $(PACKAGES)
cover:
	go test -coverprofile cover.out $(PACKAGES) && go tool cover -html=cover.out -o cover.html && open cover.html
fmt:
	go fmt ./...
clean:
	go clean
	rm -f bin/
