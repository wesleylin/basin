package tests

import (
	"testing"

	"github.com/wesleylin/basin/concurrentsequencedmap"
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

func TestConcurrentMapStreamingIntegration(t *testing.T) {
	// Initialize Basin's Concurrent Map
	m := concurrentsequencedmap.New[string, int]()

	// 1. Load data out of order (Shards will handle this)
	m.Put("A", 10)
	m.Put("B", 20)
	m.Put("C", 30)
	m.Put("D", 40)

	// 2. Build a pipeline: Stream -> Filter -> Collect
	// We want to verify that even after filtering, the Basin global order is kept.
	results, err := m.Stream2().
		Filter(func(k string, v int) bool {
			// Filter out A and D
			return v > 10 && v < 40
		}).
		Collect()

	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}

	// 3. Verify Integration
	// We expect "B" then "C" because of the global insertion sequence.
	expected := []struct {
		k string
		v int
	}{
		{"B", 20},
		{"C", 30},
	}

	if len(results) != len(expected) {
		t.Fatalf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, exp := range expected {
		if results[i].Key != exp.k || results[i].Value != exp.v {
			t.Errorf("Sequence broken at index %d: expected %s=%d, got %s=%d",
				i, exp.k, exp.v, results[i].Key, results[i].Value)
		}
	}
}
