package concurrentmap

import (
	"sync"
	"testing"
)

func TestMap_Basic(t *testing.T) {
	m := New[string, int]()

	// Test Set and Get
	m.Put("apple", 1)
	val, ok := m.Get("apple")
	if !ok || val != 1 {
		t.Errorf("expected 1, got %v", val)
	}

	// Test Update
	m.Put("apple", 2)
	val, _ = m.Get("apple")
	if val != 2 {
		t.Errorf("expected updated value 2, got %v", val)
	}

	// Test Delete
	m.Delete("apple")
	_, ok = m.Get("apple")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestMap_Concurrency(t *testing.T) {
	m := New[int, int]()
	var wg sync.WaitGroup
	numIterations := 1000
	numGoroutines := 10

	// Concurrent Writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(gID int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				// Use different keys to test shard distribution
				key := gID*numIterations + j
				m.Put(key, j)
			}
		}(i)
	}
	wg.Wait()

	// Concurrent Reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(gID int) {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				key := gID*numIterations + j
				_, ok := m.Get(key)
				if !ok {
					t.Errorf("expected to find key %d", key)
				}
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkMap_Set(b *testing.B) {
	m := New[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Put(i, i)
			i++
		}
	})
}
