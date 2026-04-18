.PHONY: build run test lint clean

BINARY_NAME=git-resume
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

run:
	go run . $(ARGS)

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)
	rm -f data/*.db

# Development
dev:
	go run . analyze --dry-run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Create required directories
setup:
	mkdir -p $(BUILD_DIR) data
