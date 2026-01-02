
// Drop is a fluent version of Pop, removing and discarding the top element of the heap.
func (h *Heap[T, P]) Drop() *Heap[T, P] {
	h.Pop()
	return h
}