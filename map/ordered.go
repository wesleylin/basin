package ordered

import (
	"iter"
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

func (m *Map[K, V]) Set(key K, val V) {
	// check if key exists
	if idx, exists := m.table[key]; exists {
		m.slots[idx].value = val
		m.slots[idx].deleted = false
		return
	}

	m.table[key] = len(m.slots)
	m.slots = append(m.slots, entry[K, V]{
		key:   key,
		value: val,
	})
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	idx, exists := m.table[key]
	if !exists || m.slots[idx].deleted {
		var zero V
		return zero, false
	}
	return m.slots[idx].value, true
}

// Delete removes a key-value pair from the map. Returns true if the key was present.
func (m *Map[K, V]) Delete(key K) bool {
	idx, exists := m.table[key]

	if !exists {
		return false
	}

	delete(m.table, key)
	m.slots[idx].deleted = true
	m.deletedCount++

	// GC optimization, clear data so garbage collector can reclaim memory
	var zeroK K
	var zeroV V
	m.slots[idx].key = zeroK
	m.slots[idx].value = zeroV

	// compact if more than half are deleted
	if m.deletedCount*2 > len(m.slots) {
		m.compact()
	}
	return true
}

// compact remaining
func (m *Map[K, V]) compact() {
	newSlots := make([]entry[K, V], 0, len(m.slots)-m.deletedCount)

	for _, e := range m.slots {
		if !e.deleted {
			// Update the table with items new position in new slice
			m.table[e.key] = len(newSlots)
			newSlots = append(newSlots, e)
		}
	}

	m.slots = newSlots
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
