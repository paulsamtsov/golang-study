.PHONY: fmt lint test build all

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

test:
	go test -race ./...

build:
	go build -o bin/ ./cmd/app/

all: fmt lint test build
