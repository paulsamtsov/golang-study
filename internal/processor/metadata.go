// Package processor handles image metadata processing with worker pool.
package processor

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

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
	log.Info().Int("workers", count).Msg("Starting worker pool")

	for i := 0; i < count; i++ {
		go func(id int) {
			log.Info().Int("worker_id", id).Msg("Worker started")
			for {
				processImage(id)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Log stats periodically
	go logStats()

	select {} // block forever
}

// logStats logs processing statistics every 5 seconds.
func logStats() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cacheSize, cacheMax := GetCacheStats()
		allStats := stats.GetStats()
		total := stats.GetTotal()

		log.Info().
			Int("total_processed", total).
			Int("cache_size", cacheSize).
			Int("cache_max", cacheMax).
			Interface("by_type", allStats).
			Msg("Processing stats")
	}
}

// processImage processes a single image metadata.
// BUGGY issues (now FIXED):
// 1. Memory Leak: LeakCache grew unbounded (10KB per frame, no eviction)
// 2. CPU Bottleneck: regexp.MatchString compiled regex on every call
// 3. Race Condition: LeakCache written from multiple goroutines without mutex
func processImage(workerID int) {
	start := time.Now()

	// Create image data identifier
	data := fmt.Sprintf("image_worker%d_%d", workerID, time.Now().UnixNano())

	// FIXED: Use pre-compiled pattern instead of regexp.MatchString
	// This avoids O(n) regex compilation on every call
	if imagePattern.MatchString(data) {
		stats.IncrementProcessed("image")

		// FIXED: Bounded cache with mutex prevents memory leak and race conditions
		// BUGGY was: direct write without bounds: LeakCache[key] = make([]byte, 1024*10)
		cacheMu.Lock()
		cached := false
		if len(leakCache) < maxCacheSize {
			key := fmt.Sprintf("key_%d", time.Now().UnixNano())
			leakCache[key] = make([]byte, 1024*10) // 10KB per entry
			cached = true
		}
		cacheMu.Unlock()

		log.Debug().
			Int("worker_id", workerID).
			Dur("duration", time.Since(start)).
			Bool("cached", cached).
			Msg("Image processed")
	}
}

// GetCacheStats returns current cache statistics (for monitoring).
func GetCacheStats() (size int, maxSize int) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	return len(leakCache), maxCacheSize
}
