package heap

import (
	"testing"
)

func TestHeapBasic(t *testing.T) {
	h := New[string, int]() // Min-Heap

	h.Push("low", 10)
	h.Push("high", 1)
	h.Push("mid", 5)

	if h.Len() != 3 {
		t.Errorf("Expected length 3, got %d", h.Len())
	}

	val, _ := h.Pop()
	if val != "high" {
		t.Errorf("Expected 'high' (1), got %s", val)
	}

	val, _ = h.Peek()
	if val != "mid" {
		t.Errorf("Peek expected 'mid' (5), got %s", val)
	}
}

func TestMaxHeap(t *testing.T) {
	h := NewMax[string, int]()

	h.Push("low", 1)
	h.Push("high", 10)
	h.Push("mid", 5)

	val, _ := h.Pop()
	if val != "high" {
		t.Errorf("Max-Heap expected 'high' (10), got %s", val)
	}
}

func TestEmptyHeap(t *testing.T) {
	h := New[int, int]()
	_, ok := h.Pop()
	if ok {
		t.Error("Pop on empty heap should return ok=false")
	}
}
