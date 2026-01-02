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

func TestUnion(t *testing.T) {
	s1 := New[string]()
	s1.Add("apple")
	s1.Add("banana")

	s2 := New[string]()
	s2.Add("banana") // Duplicate
	s2.Add("cherry")
	s2.Add("date")

	// The expected order is:
	// 1. Everything from s1 ("apple", "banana")
	// 2. Everything from s2 NOT in s1 ("cherry", "date")
	expected := []string{"apple", "banana", "cherry", "date"}

	count := 0
	for item := range Union(s1, s2) {
		if count >= len(expected) {
			t.Errorf("Union yielded more items than expected")
			break
		}
		if item != expected[count] {
			t.Errorf("At index %d: expected %s, got %s", count, expected[count], item)
		}
		count++
	}

	if count != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), count)
	}
}

func TestUnionEmpty(t *testing.T) {
	s1 := New[int]()
	s2 := New[int]()
	s2.Add(1)

	// Union of empty and non-empty
	found := false
	for item := range Union(s1, s2) {
		if item == 1 {
			found = true
		}
	}
	if !found {
		t.Error("Union with empty set failed to yield items from non-empty set")
	}
}
