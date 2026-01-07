package concurrentsequencedmap

import (
	// "container/heap"
	"iter"

	"github.com/wesleylin/basin/heap"

	"github.com/wesleylin/basin/stream"
)

// All returns a Go 1.23 iterator that yields all key-value pairs in the
// map according to their global insertion order.
// Uses a small heap-Merge with micro-locks, guaranteed to be in order for the items
// that existed when All() was called
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// Define a local struct for the heap values
		type mergeRef struct {
			key      K
			value    V
			shardIdx int
		}

		h := heap.New[uint64, mergeRef]()

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
				h.Insert(e.seq, mergeRef{
					key:      k,
					value:    e.value,
					shardIdx: i,
				})
			}
		}

		// The Merge Loop:
		// We always pull the item with the lowest sequence ID across all shards.
		for h.Len() > 0 {
			// 1. Get the globally earliest item.
			item, _, _ := h.Peek()

			// 2. Yield to the user.
			// No locks are held here, allowing the user to process data
			// without blocking map writes.
			if !yield(item.key, item.value) {
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
				h.Replace(nextE.seq, mergeRef{
					key:      nextK,
					value:    nextE.value,
					shardIdx: idx,
				})
			} else {
				// shard exhausted
				h.Pop()
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
	return stream.FromSeq2(m.All())
}
