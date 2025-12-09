VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BINARY := tosconv
DIST := dist

.PHONY: build install test clean release

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/tosconv

install:
	go install $(LDFLAGS) ./cmd/tosconv

test:
	go test -v ./...

clean:
	rm -f $(BINARY)
	rm -rf $(DIST)

release: clean
	mkdir -p $(DIST)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-arm64 ./cmd/tosconv
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-amd64 ./cmd/tosconv
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-amd64 ./cmd/tosconv
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-arm64 ./cmd/tosconv
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-windows-amd64.exe ./cmd/tosconv
