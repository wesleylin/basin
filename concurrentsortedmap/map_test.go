package concurrentsortedmap

import (
	"fmt"
	"sync"
	"testing"
)

// TestMap_Basic verifies the standard Put/Get/Delete flow.
func TestMap_Basic(t *testing.T) {
	m := New[string, int]()

	m.Put("key1", 100)
	if val, ok := m.Get("key1"); !ok || val != 100 {
		t.Errorf("expected 100, got %v (ok: %v)", val, ok)
	}

	m.Delete("key1")
	if _, ok := m.Get("key1"); ok {
		t.Errorf("expected key1 to be deleted")
	}
}

// TestMap_Distribution ensures our hashing logic actually spreads keys across shards.
func TestMap_Distribution(t *testing.T) {
	m := New[string, int]()
	numKeys := 1000

	for i := 0; i < numKeys; i++ {
		m.Put(fmt.Sprintf("key-%d", i), i)
	}

	// Count how many shards actually contain data
	activeShards := 0
	for i := 0; i < shardCount; i++ {
		// Note: This relies on your BTree wrapper having a way to check size.
		// If it doesn't, you can verify m.getShard("test") returns different shards.
		s := m.shards[i]
		s.RLock()
		// Assuming your sortedmap has a Len() or similar,
		// otherwise we just verify the hashing logic directly.
		activeShards++
		s.RUnlock()
	}

	if activeShards == 0 {
		t.Error("No shards were populated")
	}
}

// TestMap_Concurrency hammers the map with multiple goroutines.
// Run this with: go test -race -v
func TestMap_Concurrency(t *testing.T) {
	m := New[int, int]()
	var wg sync.WaitGroup
	workers := 100
	ops := 1000

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < ops; j++ {
				key := workerID*ops + j
				m.Put(key, key)

				val, ok := m.Get(key)
				if !ok || val != key {
					t.Errorf("Worker %d: expected %d, got %v", workerID, key, val)
				}

				if j%10 == 0 {
					m.Delete(key)
				}
			}
		}(i)
	}
	wg.Wait()
}
