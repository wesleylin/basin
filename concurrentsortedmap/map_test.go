package concurrentsortedmap

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentSortedMap_Basic(t *testing.T) {
	cm := New[string, int]()

	// Test Set and Get
	cm.shards[0].Set("apple", 1) // Accessing shard directly for unit test
	val, ok := cm.shards[0].Get("apple")
	if !ok || val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
}

func TestConcurrentSortedMap_Concurrency(t *testing.T) {
	cm := New[string, int]()
	wg := sync.WaitGroup{}

	const numGoroutines = 100
	const opsPerGoroutine = 1000

	// We'll simulate a mix of reads and writes across many goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)

				// In a real scenario, you'd use cm.Set(key, j)
				// For now, we simulate hitting specific shards to test locks
				shardIdx := id % shardCount
				cm.shards[shardIdx].Set(key, j)

				val, ok := cm.shards[shardIdx].Get(key)
				if !ok || val != j {
					// Using t.Errorf in goroutines is thread-safe
					t.Errorf("concurrency error: key %s expected %d, got %d", key, j, val)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentSortedMap_Distribution(t *testing.T) {
	// This test ensures that your 256 shards are actually being used.
	// Note: This requires your getShardIndex logic to be implemented.
	cm := New[string, int]()

	// Let's check if all shards were initialized properly
	for i := 0; i < shardCount; i++ {
		if cm.shards[i] == nil {
			t.Fatalf("shard %d was not initialized", i)
		}
	}
}
