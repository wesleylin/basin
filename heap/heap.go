package heap

import (
	"cmp"
	"iter"
)

type entry[T any, P cmp.Ordered] struct {
	value    T
	priority P
}

type Heap[T any, P cmp.Ordered] struct {
	data []entry[T, P]
	min  bool
}

// New returns a Min-Heap (smallest priority at the top)
func New[T any, P cmp.Ordered]() *Heap[T, P] {
	return &Heap[T, P]{min: true}
}

// NewMax returns a Max-Heap (largest priority at the top)
func NewMax[T any, P cmp.Ordered]() *Heap[T, P] {
	return &Heap[T, P]{min: false}
}

func (h *Heap[T, P]) Len() int { return len(h.data) }

func (h *Heap[T, P]) Insert(val T, priority P) {
	h.data = append(h.data, entry[T, P]{val, priority})
	h.up(len(h.data) - 1)
}

func (h *Heap[T, P]) Pop() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}

	n := len(h.data) - 1
	h.swap(0, n)
	h.down(0, n)

	item := h.data[n]
	h.data = h.data[:n]
	return item.value, true
}

func (h *Heap[T, P]) Peek() (T, bool) {
	if len(h.data) == 0 {
		var zero T
		return zero, false
	}
	return h.data[0].value, true
}

// Drain removes and yields all elements from the heap in priority order.
// Usage: for v := range h.Drain() { ... } will pop all elements.
func (h *Heap[T, P]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			val, ok := h.Pop()
			if !ok || !yield(val) {
				return
			}
		}
	}
}

// --- Internal Heap Math ---

func (h *Heap[T, P]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.less(j, i) {
			break
		}
		h.swap(i, j)
		j = i
	}
}

func (h *Heap[T, P]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.less(j2, j1) {
			j = j2 // right child
		}
		if !h.less(j, i) {
			break
		}
		h.swap(i, j)
		i = j
	}
	return i > i0
}

func (h *Heap[T, P]) less(i, j int) bool {
	if h.min {
		return h.data[i].priority < h.data[j].priority
	}
	return h.data[i].priority > h.data[j].priority
}

func (h *Heap[T, P]) swap(i, j int) {
	h.data[i], h.data[j] = h.data[j], h.data[i]
}
