# Lab 3: Image Metadata Processor - Diagnostics and Refactoring

[![CI](https://github.com/paulsamtsov/golang-study/actions/workflows/ci.yml/badge.svg)](https://github.com/paulsamtsov/golang-study/actions/workflows/ci.yml)

## Overview

This lab demonstrates detection and fixing of three critical bugs in a high-performance Go service:

1. **Memory Leak** (30 pts) - Unbounded cache growth causing OOM
2. **Race Condition** (20 pts) - Concurrent map access without synchronization
3. **CPU Bottleneck** (30 pts) - Regex recompilation on every call
4. **Remote Debugging** (10 pts) - dlv integration for production debugging
5. **Reporting** (10 pts) - Analytical findings and optimization results

## Requirements

- Go SDK (1.26.1+)
- golangci-lint
- dlv (for remote debugging)
- pprof (built-in with Go)

## Building and Running

### Build the service
```bash
make build
./bin/service
```

### Run without building
```bash
make run
```

### Run tests
```bash
make test       # All tests with race detector
make race -v    # Verbose race condition detection
```

### Performance Analysis

#### Benchmarks (Before/After Optimization)
```bash
make bench
```

Output shows optimization results:
- **Buggy**: 2519 ns/op (regex compiled every call)
- **Fixed**: 257.7 ns/op (pre-compiled regex)
- **Speedup**: ~9.8x faster! ✓

#### Profiling with pprof

**Start the service:**
```bash
make profile
```

Service starts with pprof server on `http://localhost:6060/debug/pprof/`

**Heap Profile (Memory Leak Detection):**
```bash
make heap
```

This downloads heap profile and lets you analyze:
```bash
go tool pprof heap.prof
(pprof) top
(pprof) list processImage    # Shows memory allocations
```

**CPU Profile (Performance Analysis):**
```bash
make cpu
```

Then analyze:
```bash
go tool pprof cpu.prof
(pprof) top       # Shows functions consuming most CPU
(pprof) list processImage
(pprof) web       # Generate flame graph
```

### Remote Debugging with dlv

**Terminal 1: Start debugger in headless mode:**
```bash
make debug
# Output: listening on 127.0.0.1:2345, grpc enabled
```

**Terminal 2: Connect from VS Code:**

Create `.vscode/launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to dlv",
      "type": "go",
      "mode": "remote",
      "request": "attach",
      "port": 2345,
      "host": "127.0.0.1"
    }
  ]
}
```

Then press F5 to connect debugger.

**Terminal 2: Connect from GoLand:**
- Run → Edit Configurations → Go Remote
- Host: 127.0.0.1
- Port: 2345

**Terminal 2: Connect via dlv CLI:**
```bash
dlv connect 127.0.0.1:2345
(dlv) break processImage
(dlv) continue
(dlv) print data    # Examine variables
```

## Project Structure

```
lab3-detector/
├── cmd/service/
│   └── main.go              # Entry point with pprof :6060
├── internal/processor/
│   ├── metadata.go          # Fixed: pre-compiled regex + bounded cache
│   └── metadata_bench_test.go   # Benchmarks showing 9.8x speedup
├── internal/stats/
│   ├── counter.go           # Fixed: thread-safe with sync.RWMutex
│   └── counter_test.go      # Race condition tests
├── .github/workflows/
│   └── ci.yml              # GitHub Actions with benchmarks
├── Makefile
├── go.mod
└── README.md
```

## Key Bugs Fixed

### 1. Memory Leak: LeakCache
**Problem:**
```go
// BUGGY: Grows unbounded, 10KB per frame
var LeakCache = make(map[string][]byte)
// ... accumulates forever with no cleanup
```

**Solution:**
```go
// FIXED: Bounded cache with eviction policy
const maxCacheSize = 100
if len(leakCache) < maxCacheSize {
    leakCache[key] = data
}
```

### 2. Race Condition: GlobalStats
**Problem:**
```go
// BUGGY: Concurrent writes without synchronization
var GlobalStats = make(map[string]int)
func IncrementProcessed(imageType string) {
    GlobalStats[imageType]++  // RACE!
}
```

**Solution:**
```go
// FIXED: Protected by RWMutex
var mu sync.RWMutex
func IncrementProcessed(imageType string) {
    mu.Lock()
    defer mu.Unlock()
    GlobalStats[imageType]++
}
```

### 3. CPU Bottleneck: Regex Compilation
**Problem:**
```go
// BUGGY: Compiles regex on every call (2519 ns/op)
matched, _ := regexp.MatchString(`^image_worker\d+_\d+$`, data)
```

**Solution:**
```go
// FIXED: Pre-compile once (257.7 ns/op = 9.8x faster!)
var imagePattern = regexp.MustCompile(`^image_worker\d+_\d+$`)
if imagePattern.MatchString(data) {
    // ...
}
```

## Benchmark Results

```
BenchmarkProcessImageBuggy-12      455836    2519 ns/op    7216 B/op  80 allocs/op
BenchmarkProcessImageFixed-12     4786227     257.7 ns/op    56 B/op   2 allocs/op
                                                ^9.8x faster
```

## Diagnostics Tools

| Tool | Purpose | Command |
|------|---------|---------|
| go test -race | Detect race conditions | `make race` |
| pprof | Memory & CPU profiling | `make heap` / `make cpu` |
| go test -bench | Performance comparison | `make bench` |
| dlv | Remote debugging | `make debug` |

## CI/CD Pipeline

GitHub Actions automatically:
1. Runs all tests with race detector
2. Executes benchmarks (shows before/after)
3. Runs linter checks
4. Builds binary

Status visible in PR checks and Actions tab.

## Verification Checklist

- [x] Memory leak fixed (bounded cache + mutex)
- [x] Race condition fixed (sync.RWMutex)
- [x] CPU optimized (9.8x speedup from pre-compiled regex)
- [x] All tests pass: 0 race conditions detected
- [x] Benchmarks show significant improvement
- [x] Remote debugging setup (dlv headless mode)
- [x] CI/CD pipeline configured

## References

- [Go Profiling](https://pkg.go.dev/runtime/pprof)
- [Race Detector](https://golang.org/doc/articles/race_detector)
- [pprof Documentation](https://github.com/google/pprof)
- [dlv Debugger](https://github.com/go-delve/delve)
