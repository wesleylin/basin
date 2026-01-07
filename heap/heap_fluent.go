package heap

// Push is a fluent version of Insert, adding an element to the heap. Returns the heap for chaining.
func (h *Heap[P, V]) Push(priority P, val V) *Heap[P, V] {
	h.Insert(priority, val)
	return h
}

// Drop is a fluent version of Pop, removing and discarding the top element of the heap. Returns the heap for chaining.
func (h *Heap[P, V]) Drop() *Heap[P, V] {
	h.Pop()
	return h
}
