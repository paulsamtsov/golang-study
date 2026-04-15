// Package stats provides concurrent statistics counters.
package stats

import "sync"

// FIXED: Added mutex to prevent race condition
// BUGGY was: var GlobalStats = make(map[string]int) with concurrent writes
var (
	// GlobalStats tracks counts of processed items by type
	GlobalStats = make(map[string]int)

	// mu protects all concurrent access to GlobalStats map
	mu = &sync.RWMutex{}
)

// IncrementProcessed safely increments the counter for a given image type.
// FIXED: Wrapped with mutex to prevent data races
// BUGGY was: GlobalStats[imageType]++ without synchronization
func IncrementProcessed(imageType string) {
	mu.Lock()
	defer mu.Unlock()
	GlobalStats[imageType]++
}

// GetStats safely reads all current statistics.
func GetStats() map[string]int {
	mu.RLock()
	defer mu.RUnlock()

	// Create a copy to avoid external modifications
	result := make(map[string]int)
	for k, v := range GlobalStats {
		result[k] = v
	}
	return result
}
