// Package processor handles image metadata processing with worker pool.
package processor

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/paulsamtsov/lab3-detector/internal/stats"
)

// FIXED: Pre-compiled regexp to avoid recompilation overhead
// BUGGY was: regexp.MatchString(...) in every processImage call
var imagePattern = regexp.MustCompile(`^image_worker\d+_\d+$`)

// FIXED: Bounded cache with mutex to prevent memory leak and race conditions
// BUGGY was: var LeakCache = make(map[string][]byte) with no bounds
const maxCacheSize = 100

var (
	leakCache = make(map[string][]byte)
	cacheMu   = &sync.Mutex{}
)

// RunWorkerPool starts a pool of worker goroutines for image processing.
func RunWorkerPool(count int) {
	for i := 0; i < count; i++ {
		go func(id int) {
			for {
				processImage(id)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}
	select {} // block forever
}

// processImage processes a single image metadata.
// BUGGY issues (now FIXED):
// 1. Memory Leak: LeakCache grew unbounded (10KB per frame, no eviction)
// 2. CPU Bottleneck: regexp.MatchString compiled regex on every call
// 3. Race Condition: LeakCache written from multiple goroutines without mutex
func processImage(workerID int) {
	// Create image data identifier
	data := fmt.Sprintf("image_worker%d_%d", workerID, time.Now().UnixNano())

	// FIXED: Use pre-compiled pattern instead of regexp.MatchString
	// This avoids O(n) regex compilation on every call
	if imagePattern.MatchString(data) {
		stats.IncrementProcessed("image")

		// FIXED: Bounded cache with mutex prevents memory leak and race conditions
		// BUGGY was: direct write without bounds: LeakCache[key] = make([]byte, 1024*10)
		cacheMu.Lock()
		if len(leakCache) < maxCacheSize {
			key := fmt.Sprintf("key_%d", time.Now().UnixNano())
			leakCache[key] = make([]byte, 1024*10) // 10KB per entry
		}
		cacheMu.Unlock()
	}
}

// GetCacheStats returns current cache statistics (for monitoring).
func GetCacheStats() (size int, maxSize int) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	return len(leakCache), maxCacheSize
}
