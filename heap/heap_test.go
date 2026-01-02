package heap

import (
	"slices"
	"testing"
)

func TestHeapBasic(t *testing.T) {
	h := New[string, int]() // Min-Heap

	h.Insert("low", 10)
	h.Insert("high", 1)
	h.Insert("mid", 5)

	if h.Len() != 3 {
		t.Errorf("Expected length 3, got %d", h.Len())
	}

	val, _ := h.Pop()
	if val != "high" {
		t.Errorf("Expected 'high' (1), got %s", val)
	}

	val, _, _ = h.Peek()
	if val != "mid" {
		t.Errorf("Peek expected 'mid' (5), got %s", val)
	}
}

func TestMaxHeap(t *testing.T) {
	h := NewMax[string, int]()

	h.Insert("low", 1)
	h.Insert("high", 10)
	h.Insert("mid", 5)

	val, _ := h.Pop()
	if val != "high" {
		t.Errorf("Max-Heap expected 'high' (10), got %s", val)
	}
}

func TestHeapDrain(t *testing.T) {
	h := New[int, int]()
	input := []int{5, 3, 8, 1}
	for _, v := range input {
		h.Insert(v, v)
	}

	var result []int
	for val := range h.Drain() {
		result = append(result, val)
	}

	expected := []int{1, 3, 5, 8}
	if !slices.Equal(result, expected) {
		t.Errorf("Drain order incorrect. Got %v, want %v", result, expected)
	}

	if h.Len() != 0 {
		t.Error("Heap should be empty after Drain")
	}
}

func TestUnstableNature(t *testing.T) {
	// This test documents that we do NOT guarantee order for equal priorities
	h := New[string, int]()

	// Adding three items with the same priority
	h.Insert("A", 1)
	h.Insert("B", 1)
	h.Insert("C", 1)

	var result []string
	for val := range h.Drain() {
		result = append(result, val)
	}

	// In a standard binary heap, "A" is usually first,
	// but the rest depends on the internal tree swaps.
	t.Logf("Unstable order result: %v", result)
}

func TestEmptyHeap(t *testing.T) {
	h := New[int, int]()
	_, ok := h.Pop()
	if ok {
		t.Error("Pop on empty heap should return ok=false")
	}
}
