package concurrentsortedmap

import (
	"iter"

	"github.com/wesleylin/basin/heap"
)

// All returns an iterator over all key-value pairs in sorted order.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		h := heap.New[K, V]()

		// 1. Initialize: Open an iterator for every shard and grab the first item
		for i := 0; i < shardCount; i++ {
			s := m.shards[i]
			s.RLock()
			// We need to keep the lock or a snapshot?
			// For a simple implementation, we'll collect the shard's iterator.
			next, stop := iter.Pull2(s.data.All())
			defer stop()

			if k, v, ok := next(); ok {
				heap.Push(h, shardItem[K, V]{k: k, v: v, next: next})
			}
			s.RUnlock()
		}

		// 2. Merge: Always pull the smallest key from the heap
		for h.Len() > 0 {
			item := heap.Pop(h).(shardItem[K, V])

			// Yield the smallest current value to the user
			if !yield(item.k, item.v) {
				return
			}

			// Pull the next item from the SAME shard that just produced a value
			if k, v, ok := item.next(); ok {
				heap.Push(h, shardItem[K, V]{k: k, v: v, next: item.next})
			}
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
