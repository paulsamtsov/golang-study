package stats

import (
	"sync"
	"testing"
)

// TestIncrementProcessedNoRace verifies that concurrent increments are race-free.
// Run with: go test -race ./...
// BEFORE FIX: Would report race condition
// AFTER FIX: No race condition detected
func TestIncrementProcessedNoRace(t *testing.T) {
	// Reset global stats
	mu.Lock()
	GlobalStats = make(map[string]int)
	mu.Unlock()

	numGoroutines := 100
	incrementsPerGoroutine := 1000
	expectedTotal := numGoroutines * incrementsPerGoroutine

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				IncrementProcessed("jpeg")
			}
		}()
	}
	wg.Wait()

	// Verify correctness
	stats := GetStats()
	if stats["jpeg"] != expectedTotal {
		t.Errorf("Expected %d, got %d", expectedTotal, stats["jpeg"])
	}
}

// TestConcurrentReadWrite verifies thread-safety of read and write operations.
func TestConcurrentReadWrite(t *testing.T) {
	mu.Lock()
	GlobalStats = make(map[string]int)
	mu.Unlock()

	numWriters := 10
	numReaders := 10
	operations := 1000

	var wg sync.WaitGroup

	// Start writers
	for w := 0; w < numWriters; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < operations; i++ {
				IncrementProcessed("image")
			}
		}()
	}

	// Start readers
	for r := 0; r < numReaders; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < operations; i++ {
				GetStats()
			}
		}()
	}

	wg.Wait()

	// Verify final count
	stats := GetStats()
	expectedCount := numWriters * operations
	if stats["image"] != expectedCount {
		t.Errorf("Expected %d, got %d", expectedCount, stats["image"])
	}
}
