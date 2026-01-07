package sortedmap

import (
	"fmt"
	"testing"
)

func TestSortedMap(t *testing.T) {
	sm := New[string, int]()

	// 1. Test Set and Get
	t.Run("SetAndGet", func(t *testing.T) {
		sm.Put("apple", 10)
		sm.Put("banana", 20)

		val, ok := sm.Get("apple")
		if !ok || val != 10 {
			t.Errorf("expected 10, got %d", val)
		}
	})

	// 2. Test Overwrite
	t.Run("Overwrite", func(t *testing.T) {
		sm.Put("apple", 15)
		val, _ := sm.Get("apple")
		if val != 15 {
			t.Errorf("expected overwritten value 15, got %d", val)
		}
	})

	// 3. Test Delete
	t.Run("Delete", func(t *testing.T) {
		sm.Delete("banana")
		_, ok := sm.Get("banana")
		if ok {
			t.Error("expected banana to be deleted")
		}
	})

	// 4. Test Iteration (Go 1.23 range over func)
	t.Run("Iteration", func(t *testing.T) {
		// Clear and reset for clean test
		sm = New[string, int]()
		data := map[string]int{"c": 3, "a": 1, "b": 2}
		for k, v := range data {
			sm.Put(k, v)
		}

		expectedOrder := []string{"a", "b", "c"}
		i := 0
		// This uses the All() iterator you wrote
		for k, v := range sm.All() {
			if k != expectedOrder[i] {
				t.Errorf("at index %d, expected key %s, got %s", i, expectedOrder[i], k)
			}
			if v != data[k] {
				t.Errorf("value mismatch for key %s", k)
			}
			i++
		}
	})

	// 5. Test Range Scan
	t.Run("RangeScan", func(t *testing.T) {
		sm = New[string, int]()
		for i := 0; i < 10; i++ {
			sm.Put(fmt.Sprintf("key-%02d", i), i)
		}

		// Scan from key-03 to key-06 (exclusive of end)
		count := 0
		for k, _ := range sm.Range("key-03", "key-07") {
			count++
			fmt.Println("Range found:", k)
		}

		if count != 4 {
			t.Errorf("expected 4 items in range, got %d", count)
		}
	})
}
