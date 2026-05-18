.PHONY: build run clean test install deps

BINARY=popeye
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/popeye

run: build
	./$(BUILD_DIR)/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)
	go clean

test:
	go test -v ./...

install: build
	cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/

deps:
	go mod tidy
	go mod download

lint:
	golangci-lint run

fmt:
	go fmt ./...

dev:
	go run ./cmd/popeye
