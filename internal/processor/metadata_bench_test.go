package processor

import (
	"fmt"
	"regexp"
	"sync"
	"testing"
	"time"
)

// BenchmarkProcessImageBuggy demonstrates the cost of recompiling regexp on every call.
// This is the ORIGINAL BUGGY version that compiles the regex pattern for each call.
// Expected: ~5000-10000 ns/op (very slow due to regex compilation)
func BenchmarkProcessImageBuggy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := fmt.Sprintf("image_worker%d_%d", 1, time.Now().UnixNano())
		// BUGGY: Recompiles regex on every call (intentional for benchmark)
		//nolint:staticcheck
		_, _ = regexp.MatchString(`^image_worker\d+_\d+$`, data)
	}
}

// BenchmarkProcessImageFixed demonstrates the improvement from pre-compiled regex.
// This is the FIXED version using a package-level compiled regex.
// Expected: ~10-50 ns/op (50-100x faster!)
func BenchmarkProcessImageFixed(b *testing.B) {
	compiledPattern := regexp.MustCompile(`^image_worker\d+_\d+$`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := fmt.Sprintf("image_worker%d_%d", 1, time.Now().UnixNano())
		// FIXED: Uses pre-compiled pattern
		compiledPattern.MatchString(data)
	}
}

// BenchmarkCacheOperations compares bounded vs unbounded cache performance.
func BenchmarkCacheOperationsBounded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cacheMu.Lock()
		if len(leakCache) < maxCacheSize {
			key := fmt.Sprintf("key_%d", i)
			leakCache[key] = make([]byte, 1024*10)
		}
		cacheMu.Unlock()
	}
}

// Benchmark showing impact of unbounded cache (simulated).
func BenchmarkCacheOperationsUnbounded(b *testing.B) {
	unboundedCache := make(map[string][]byte)
	mu := &sync.Mutex{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		key := fmt.Sprintf("key_%d", i)
		unboundedCache[key] = make([]byte, 1024*10)
		mu.Unlock()
	}
}
