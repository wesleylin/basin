package concurrentsortedmap

import (
	"iter"

	"github.com/wesleylin/basin/heap"
)

// All returns an iterator over all key-value pairs in sorted order.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// shardCursor acts as a local buffer for each shard's stream
		type shardCursor struct {
			curr V                   // The value we've already pulled and are "holding"
			next func() (K, V, bool) // The function to pull the next pair from this shard
		}

		// The heap is keyed by K (priority) and stores our shardCursor
		h := heap.New[K, shardCursor]()

		// 1. INITIALIZATION PHASE
		// Open pull iterators for all 256 shards and grab the first item from each.
		for i := 0; i < shardCount; i++ {
			s := m.shards[i]

			s.RLock()
			// Create a pull iterator from the shard's push iterator.
			// This captures a point-in-time view of the shard's B-Tree.
			pull, stop := iter.Pull2(s.data.All())

			// Crucial: stop() must be called to release B-Tree resources.
			// These defers will execute when the All() function returns.
			defer stop()

			if k, v, ok := pull(); ok {
				// We push the first key as the priority, and the rest in the cursor.
				h.Push(shardCursor{curr: v, next: pull}, k)
			}
			s.RUnlock()
		}

		// 2. MERGE PHASE
		// Always pop the globally smallest key across all 256 shards.
		for h.Len() > 0 {
			// k: the priority (key)
			// cursor: the struct containing the current value and the puller
			k, cursor, _ := h.Pop()

			// Yield the smallest current pair to the user's for-range loop.
			// If yield returns false, the user has 'broken' out of the loop.
			if !yield(k, cursor.curr) {
				return
			}

			// Refill: Get the next item from the shard that just provided 'k'.
			if nextK, nextV, ok := cursor.next(); ok {
				// Update the cursor with the NEW value we just pulled...
				cursor.curr = nextV

				// ...and push it back into the heap with its NEW key as priority.
				h.Push(cursor, nextK)
			}
			// If !ok, the shard is exhausted and we simply don't push it back.
		}
	}
}

// Keys returns an iterator over all keys in sorted order.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator over all values in sorted order of their keys.
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}
