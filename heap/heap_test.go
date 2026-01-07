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

func TestHeapReplace(t *testing.T) {
	h := New[string, int]()
	h.Insert("medium", 50)
	h.Insert("low", 100)
	h.Insert("high", 10)

	// Replace "high" (10) with "ultra-low" (200)
	// The root was 10, now the new root should be 50
	h.Replace("ultra-low", 200)

	val, _ := h.Pop()
	if val != "medium" {
		t.Errorf("Expected medium (50) after replace, got %s", val)
	}

	val, _ = h.Pop()
	if val != "low" {
		t.Errorf("Expected low (100), got %s", val)
	}

	val, _ = h.Pop()
	if val != "ultra-low" {
		t.Errorf("Expected ultra-low (200), got %s", val)
	}
}

func TestHeapFix(t *testing.T) {
	h := New[string, int]()
	h.Insert("A", 10)
	h.Insert("B", 20)
	h.Insert("C", 30)

	// Manually sabotage the priority of the root (A)
	h.data[0].priority = 40
	// Fix it
	h.Fix(0)

	// New root should be B (20)
	val, _ := h.Pop()
	if val != "B" {
		t.Errorf("After Fix, expected B at root, got %s", val)
	}
}

func TestHeapMemorySafety(t *testing.T) {
	type complexObj struct {
		data []byte
	}
	h := New[*complexObj, int]()
	obj := &complexObj{data: make([]byte, 1024)}

	h.Insert(obj, 10)
	h.Pop()

	// Check if the underlying slice cleared the reference
	if h.data[:1][0].value != nil {
		t.Error("Pop did not zero out the underlying array element; potential memory leak")
	}
}
