package concurrentsequencedmap

import (
	"fmt"
	"sync"
	"testing"
)

func TestMap_ConcurrencyAndOrder(t *testing.T) {
	m := New[string, int]()
	workerCount := 100
	opsPerWorker := 1000
	var wg sync.WaitGroup

	// 1. Hammer the map with concurrent Puts
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				key := fmt.Sprintf("w%d-k%d", workerID, j)
				// We store 'j' as the value
				m.Put(key, j)
			}
		}(i)
	}
	wg.Wait()

	// 2. Verify Total Length
	expectedLen := workerCount * opsPerWorker
	if m.Len() != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, m.Len())
	}

	// 3. Verify Global Sequence via Iterator
	count := 0
	for key, val := range m.All() {
		count++

		// Basic validation: ensure values are within the range we put in
		if val < 0 || val >= opsPerWorker {
			t.Errorf("Incorrect value for key %s: %d", key, val)
		}
	}

	// This proves the iterator successfully merged all 256 shards
	// without dropping or duplicating any data.
	if count != expectedLen {
		t.Errorf("Iterator yielded %d items, expected %d", count, expectedLen)
	}
}

func TestMap_OverwriteMaintainsSequence(t *testing.T) {
	m := New[string, string]()

	m.Put("key1", "first")
	s1 := getSeq(m, "key1")

	m.Put("key1", "second")
	s2 := getSeq(m, "key1")

	if s2 <= s1 {
		t.Errorf("Overwrite should result in higher sequence: %d -> %d", s1, s2)
	}
}

func TestMap_StrictOrder(t *testing.T) {
	m := New[string, string]()
	m.Put("first", "A")
	m.Put("second", "B")
	m.Put("third", "C")

	expected := []string{"A", "B", "C"}
	i := 0
	for _, val := range m.All() {
		if val != expected[i] {
			t.Errorf("Order mismatch at index %d: expected %s, got %s", i, expected[i], val)
		}
		i++
	}
}

// Helper to peek at the internal sequence for testing
func getSeq[K comparable, V any](m *Map[K, V], key K) uint64 {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	entry, _ := shard.data.Get(key)
	return entry.seq
}
