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
