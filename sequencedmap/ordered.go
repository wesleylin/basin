package sequencedmap

import (
	"iter"

	"github.com/wesleylin/basin/stream"
)

type entry[K comparable, V any] struct {
	key     K
	value   V
	deleted bool
}

type Map[K comparable, V any] struct {
	table        map[K]int
	slots        []entry[K, V]
	deletedCount int
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		table: make(map[K]int),
		// pre allocate small amount of space
		slots: make([]entry[K, V], 0, 8),
	}
}

func NewWithCapacity[K comparable, V any](cap int) *Map[K, V] {
	return &Map[K, V]{
		table: make(map[K]int, cap),
		slots: make([]entry[K, V], 0, cap),
	}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	idx, exists := m.table[key]
	if !exists || m.slots[idx].deleted {
		var zero V
		return zero, false
	}
	return m.slots[idx].value, true
}

// Put sets the value for a key in the map.
// Brand new insertions are placed at the end in order
// If the key already exists, its value is updated but its order remains unchanged.
// If the key was deleted previously it will be moved to the end as a new insertion.
func (m *Map[K, V]) Put(key K, val V) {
	// check if key exists
	if idx, exists := m.table[key]; exists {
		m.slots[idx].value = val
		return
	}

	m.table[key] = len(m.slots)
	m.slots = append(m.slots, entry[K, V]{
		key:   key,
		value: val,
	})
}

// Delete removes a key-value pair from the map. Returns true if the key was present.
func (m *Map[K, V]) Delete(key K) bool {
	idx, exists := m.table[key]
	if !exists {
		return false
	}

	// 1. Remove from lookup table
	delete(m.table, key)

	// 2. Clear the entry in the slice immediately.
	// Overwrite with a zero-value entry{deleted: true},
	// we release references to the Key and Value so the GC can free them
	// without waiting for a Compact().
	m.slots[idx] = entry[K, V]{deleted: true}
	m.deletedCount++

	// 3. Amortized compaction
	// We only pay the O(n) cost when the "waste" is significant.
	// 1024 is a sweet spot to avoid thrashing on small maps.
	if m.deletedCount > 1024 && m.deletedCount*2 > len(m.slots) {
		m.Compact()
	}

	return true
}

// compact remaining
func (m *Map[K, V]) Compact() {
	if m.deletedCount == 0 {
		return
	}

	// 'j' represents the next position for a live element
	j := 0
	for i := 0; i < len(m.slots); i++ {
		if m.slots[i].deleted {
			continue
		}

		// Move element forward if there's a gap
		if i != j {
			m.slots[j] = m.slots[i]
			// Update the table with the new index
			m.table[m.slots[j].key] = j
		}
		j++
	}

	// 40GB Safety Step:
	// We must zero out the "tail" of the slice that we just moved away from.
	// If we don't, the old pointers will stay in memory until the slice grows!
	for k := j; k < len(m.slots); k++ {
		m.slots[k] = entry[K, V]{}
	}

	// Reslice to the new live length
	m.slots = m.slots[:j]
	m.deletedCount = 0
}

// All returns an iterator for all key-value pairs in order.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, e := range m.slots {
			if e.deleted {
				continue
			}
			if !yield(e.key, e.value) {
				return
			}
		}
	}
}

// Keys returns an iterator for the keys in order.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, e := range m.slots {
			if !e.deleted {
				if !yield(e.key) {
					return
				}
			}
		}
	}
}

// Values returns an iterator for the values in order.
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, e := range m.slots {
			if !e.deleted {
				if !yield(e.value) {
					return
				}
			}
		}
	}
}

// Stream returns a Basin Stream2 which yields keys and values in insertion order.
func (m *Map[K, V]) Stream2() stream.Stream2[K, V] {
	var err error

	// Create the iterator logic
	seq := func(yield func(K, V) bool) {
		// Replace this loop with however your map actually iterates.
		// Example if you use a slice of entries internally:
		for _, entry := range m.slots {
			if !yield(entry.key, entry.value) {
				return
			}
		}
	}

	return stream.New2(seq, &err)
}

// Convenience methods
func (m *Map[K, V]) Len() int {
	return len(m.table)
}

func (m *Map[K, V]) Has(key K) bool {
	_, exists := m.table[key]
	return exists
}

func (m *Map[K, V]) Clear() {
	clear(m.table) // Built-in 'clear' (Go 1.21+) empties the map but keeps memory
	m.slots = m.slots[:0]
	m.deletedCount = 0
}
