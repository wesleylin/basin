package heap

// Push is a fluent version of Insert, adding an element to the heap. Returns the heap for chaining.
func (h *Heap[T, P]) Push(val T, priority P) *Heap[T, P] {
	h.Insert(priority, val)
	return h
}

// Drop is a fluent version of Pop, removing and discarding the top element of the heap. Returns the heap for chaining.
func (h *Heap[T, P]) Drop() *Heap[T, P] {
	h.Pop()
	return h
}
