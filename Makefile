.PHONY: fmt lint test build all race bench profile debug help

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

test:
	go test -race ./...

build:
	go build -o bin/service ./cmd/service/

run:
	go run ./cmd/service/

all: fmt lint test build

# Lab 3 specific targets
bench:
	go test -bench=. -benchmem -count=3 ./internal/processor/

race:
	go test -race -v ./...

profile:
	go run ./cmd/service/ &

heap:
	@echo "Downloading heap profile from running service..."
	curl -s http://localhost:6060/debug/pprof/heap > heap.prof
	@echo "Heap profile saved to heap.prof"
	@echo "Analyze with: go tool pprof heap.prof"

cpu:
	@echo "Downloading CPU profile from running service..."
	curl -s http://localhost:6060/debug/pprof/profile?seconds=5 > cpu.prof
	@echo "CPU profile saved to cpu.prof"
	@echo "Analyze with: go tool pprof cpu.prof"

debug:
	dlv debug ./cmd/service/ --headless --listen=:2345 --api-version=2

help:
	@echo "Available targets:"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  test       - Run tests with race detector"
	@echo "  build      - Build binary"
	@echo "  run        - Run service"
	@echo "  all        - fmt + lint + test + build"
	@echo ""
	@echo "Lab 3 targets:"
	@echo "  bench      - Run benchmarks (before/after optimization)"
	@echo "  race       - Verbose race condition detection"
	@echo "  profile    - Start service for profiling"
	@echo "  heap       - Download heap profile"
	@echo "  cpu        - Download CPU profile"
	@echo "  debug      - Start dlv debugger in headless mode"
