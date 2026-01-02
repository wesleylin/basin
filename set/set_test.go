package set

import (
	"fmt"
	"testing"
)

func TestSetBasic(t *testing.T) {
	s := New[string]()
	fmt.Println("Created new set:", s)

	// 1. adding value
	added := s.Add("apple")
	if !added {
		t.Errorf("Expected apple to be newly added")
	}
	added = s.Add("banana")
	if !added {
		t.Errorf("Expected banana to be newly added")
	}

	// 2. checking existence
	exists := s.Has("apple")
	if !exists {
		t.Errorf("Expected apple to exist in the set")
	}

	// 3. adding duplicate
	added = s.Add("apple")
	if added {
		t.Errorf("Expected apple to not be newly added again")
	}
}

func TestSetDelete(t *testing.T) {
	s := New[int]()
	s.Add(10)
	s.Add(20)

	s.Delete(10)
	if s.Has(10) || s.Len() != 1 {
		t.Error("Delete failed to remove item or update length")
	}

	// Re-add deleted item
	if !s.Add(10) {
		t.Error("Should be able to re-add a deleted item")
	}
}

func TestSetIteration(t *testing.T) {
	s := New[int]()
	items := []int{100, 200, 300}
	for _, v := range items {
		s.Add(v)
	}

	i := 0
	for val := range s.All() {
		if val != items[i] {
			t.Errorf("Iteration order broken: expected %d, got %d", items[i], val)
		}
		i++
	}
}

func TestSetClear(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Clear()
	if s.Len() != 0 {
		t.Error("Clear failed")
	}
}
