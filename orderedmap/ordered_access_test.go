package orderedmap_test

import (
	"testing"

	"github.com/wesleylin/basin/orderedmap"
)

func TestMap_Backward(t *testing.T) {
	m := orderedmap.New[string, int]()

	// 1. Setup initial state
	m.Put("A", 1)
	m.Put("B", 2)
	m.Put("C", 3)

	// 2. Delete the middle item
	m.Delete("B")

	// 3. Update an item (should stay in its original spot)
	m.Put("A", 10)

	// 4. Re-insert a deleted item (should move to the end)
	m.Put("B", 20)

	// Expected Forward: A, C, B
	// Expected Backward: B, C, A

	expected := []string{"B", "C", "A"}
	var actual []string

	// Use the new Backward iterator
	for k := range m.Backward() {
		actual = append(actual, k)
	}

	if len(actual) != len(expected) {
		t.Fatalf("Expected %d keys, got %d", len(expected), len(actual))
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("At index %d: expected %s, got %s", i, expected[i], actual[i])
		}
	}
}
