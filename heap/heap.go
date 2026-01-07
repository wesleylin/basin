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

// Insert adds a value with the given priority to the heap and restores the heap invariant.
func (h *Heap[T, P]) Insert(val T, priority P) {
	h.data = append(h.data, entry[T, P]{val, priority})
	h.up(len(h.data) - 1)
}

// Replace is a slightly faster way of doing a Pop() then immediate an Insert(T, P).
// Specifically a high-performance path for doing K-Way merges.
// It overwrites the root and bubbles it down, saving an 'up' pass.
func (h *Heap[T, P]) Replace(val T, priority P) {
	if len(h.data) == 0 {
		h.Insert(val, priority)
		return
	}
	h.data[0] = entry[T, P]{val, priority}
	h.down(0, len(h.data))
}

// Fix re-establishes heap order after an element at index i has changed its priority.
func (h *Heap[T, P]) Fix(i int) {
	if i < 0 || i >= len(h.data) {
		return
	}
	if !h.down(i, len(h.data)) {
		h.up(i)
	}
}

func (h *Heap[T, P]) Pop() (P, T, bool) {
	if len(h.data) == 0 {
		var zeroP P
		var zeroT T
		return zeroP, zeroT, false
	}

	n := len(h.data) - 1
	h.swap(0, n)
	h.down(0, n)

	item := h.data[n]

	// Memory Safety: Zero out the slot to prevent stale pointer leaks
	var zero T
	var zeroP P
	h.data[n] = entry[T, P]{zero, zeroP}

	h.data = h.data[:n]
	return item.priority, item.value, true
}

func (h *Heap[T, P]) Peek() (T, P, bool) {
	if len(h.data) == 0 {
		var zero T
		var zeroP P
		return zero, zeroP, false
	}
	return h.data[0].value, h.data[0].priority, true
}

// Drain removes and yields all elements from the heap in priority order.
// Usage: for v := range h.Drain() { ... } will pop all elements.
func (h *Heap[T, P]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			_, val, ok := h.Pop()
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
