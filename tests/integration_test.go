package tests

import (
	"testing"

	orderedmap "github.com/wesleylin/basin/sequencedmap"
)

func TestOrderedMapStreaming(t *testing.T) {
	m := orderedmap.New[string, int]()
	m.Put("first", 1)
	m.Put("second", 2)
	m.Put("third", 3)

	// Stream from the OrderedMap, filter, and collect
	results, err := m.Stream2().
		Filter(func(k string, v int) bool {
			return v > 1
		}).
		Collect()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify order is preserved: "second" must come before "third"
	expected := []struct {
		k string
		v int
	}{
		{"second", 2},
		{"third", 3},
	}

	if len(results) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(results))
	}

	for i, exp := range expected {
		if results[i].Key != exp.k || results[i].Value != exp.v {
			t.Errorf("at index %d: expected %s=%d, got %s=%d",
				i, exp.k, exp.v, results[i].Key, results[i].Value)
		}
	}
}
