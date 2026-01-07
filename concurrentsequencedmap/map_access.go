package concurrentsequencedmap

import (
	"container/heap"
	"iter"

	"github.com/wesleylin/basin/stream"
)

// mergeItem tracks the current head of one specific shard's iterator.
type mergeItem[K comparable, V any] struct {
	key      K
	entry    globalEntry[V]
	shardIdx int
}

// mergeHeap implements heap.Interface to provide an O(log N) min-priority queue
// based on the global sequence ID.
type mergeHeap[K comparable, V any] []mergeItem[K, V]

func (h mergeHeap[K, V]) Len() int           { return len(h) }
func (h mergeHeap[K, V]) Less(i, j int) bool { return h[i].entry.seq < h[j].entry.seq }
func (h mergeHeap[K, V]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *mergeHeap[K, V]) Push(x any) {
	*h = append(*h, x.(mergeItem[K, V]))
}

func (h *mergeHeap[K, V]) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// All returns a Go 1.23 iterator that yields all key-value pairs in the
// map according to their global insertion order.
// Uses a small heap-Merge with micro-locks, guaranteed to be in order for the items
// that existed when All() was called
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		h := &mergeHeap[K, V]{}
		heap.Init(h)

		// Create pull iterators for all 256 shards.
		// iter.Pull2 starts a background goroutine for each shard.
		nexts := make([]func() (K, globalEntry[V], bool), ShardCount)
		stops := make([]func(), ShardCount)

		for i := 0; i < ShardCount; i++ {
			s := m.shards[i]

			// We use a closure to ensure we are calling the shard's All() correctly.
			shardIter := s.data.All()
			nexts[i], stops[i] = iter.Pull2(shardIter)

			// Crucial: Ensure goroutines are cleaned up if iteration ends early.
			defer stops[i]()

			// Initial "prime" of the heap:
			// We lock the shard briefly just to grab the first element.
			s.RLock()
			k, e, ok := nexts[i]()
			s.RUnlock()

			if ok {
				heap.Push(h, mergeItem[K, V]{
					key:      k,
					entry:    e,
					shardIdx: i,
				})
			}
		}

		// The Merge Loop:
		// We always pull the item with the lowest sequence ID across all shards.
		for h.Len() > 0 {
			// 1. Get the globally earliest item.
			item := heap.Pop(h).(mergeItem[K, V])

			// 2. Yield to the user.
			// No locks are held here, allowing the user to process data
			// without blocking map writes.
			if !yield(item.key, item.entry.value) {
				return
			}

			// 3. Refill the heap from the shard we just took from.
			idx := item.shardIdx
			s := m.shards[idx]

			// MICRO-LOCK: Lock ONLY the shard we are advancing.
			s.RLock()
			nextK, nextE, ok := nexts[idx]()
			s.RUnlock()

			if ok {
				heap.Push(h, mergeItem[K, V]{
					key:      nextK,
					entry:    nextE,
					shardIdx: idx,
				})
			}
		}
	}
}

// Keys returns an iterator for the map's keys in insertion order.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	// TODO: possibly optimize to not call All to remove retrieving Value() as well
	return func(yield func(K) bool) {
		for k, _ := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator for the map's values in insertion order.
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}

// Stream returns a new Stream initialized with the globally ordered data from this map.
func (m *Map[K, V]) Stream2() stream.Stream2[K, V] {
	// We pass the global iterator directly into the stream constructor
	return stream.New2(m.All(), nil /*no error passed */)
}
