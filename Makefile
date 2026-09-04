VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo unknown)
LDFLAGS := -X github.com/antoniorodr/folkctl/cmd.version=$(VERSION)
BINARY  := folkctl

.PHONY: build run clean test

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

test:
	go test ./...
