.PHONY: build install clean test

# Build the binary
build:
	go build -o bin/shai ./cmd/shai

# Install the binary
install:
	go install ./cmd/shai

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Get dependencies
deps:
	go mod tidy
