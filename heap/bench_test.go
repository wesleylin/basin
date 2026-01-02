package heap

import (
	"container/heap"
	"testing"
)

// --- Standard Library implementation for comparison ---

type stdItem struct {
	val      int
	priority int
}

type stdHeap []stdItem

func (h stdHeap) Len() int           { return len(h) }
func (h stdHeap) Less(i, j int) bool { return h[i].priority < h[j].priority }
func (h stdHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *stdHeap) Push(x any)        { *h = append(*h, x.(stdItem)) }
func (h *stdHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// --- Benchmarks ---

func BenchmarkBasinHeap_PushPop(b *testing.B) {
	h := New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Push(i, i%100)
		if h.Len() > 1000 {
			h.Pop()
		}
	}
}

func BenchmarkStdHeap_PushPop(b *testing.B) {
	h := &stdHeap{}
	heap.Init(h)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Push(h, stdItem{i, i % 100})
		if h.Len() > 1000 {
			heap.Pop(h)
		}
	}
}
