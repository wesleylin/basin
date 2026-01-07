package sortedmap

import (
	"cmp"
	"iter"

	"github.com/tidwall/btree"
)

// kv is our internal container to keep keys and values together in the tree nodes.
type kv[K cmp.Ordered, V any] struct {
	key   K
	value V
}

// SortedMap wraps tidwall/btree to provide a clean, generic API for Basin shards.
type SortedMap[K cmp.Ordered, V any] struct {
	tree *btree.BTreeG[kv[K, V]]
	hint btree.PathHint // Optimizes sequential operations
}

// New initializes a new SortedMap using the tidwall/btree engine.
func New[K cmp.Ordered, V any]() *SortedMap[K, V] {
	return &SortedMap[K, V]{
		tree: btree.NewBTreeG(func(a, b kv[K, V]) bool {
			return a.key < b.key
		}),
	}
}

// Set inserts or updates a key-value pair.
func (m *SortedMap[K, V]) Put(key K, value V) {
	// Using SetHint makes sequential writes significantly faster.
	m.tree.SetHint(kv[K, V]{key, value}, &m.hint)
}

// Get retrieves a value. Returns the zero value and false if not found.
func (m *SortedMap[K, V]) Get(key K) (V, bool) {
	item, ok := m.tree.GetHint(kv[K, V]{key: key}, &m.hint)
	if !ok {
		var zero V
		return zero, false
	}
	return item.value, true
}

// Delete removes a key from the map.
func (m *SortedMap[K, V]) Delete(key K) bool {
	_, ok := m.tree.DeleteHint(kv[K, V]{key: key}, &m.hint)
	return ok
}

// All returns a Go 1.23 iterator for a full sorted scan.
func (m *SortedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.tree.Scan(func(item kv[K, V]) bool {
			return yield(item.key, item.value)
		})
	}
}

// Range allows for efficient sub-section scans (useful for time-series or prefix queries).
func (m *SortedMap[K, V]) Range(start, end K) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.tree.Ascend(kv[K, V]{key: start}, func(item kv[K, V]) bool {
			if item.key >= end {
				return false // Stop iterating
			}
			return yield(item.key, item.value)
		})
	}
}
