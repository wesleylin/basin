package orderedmap

import "iter"

// Backward returns an iterator that traverses the map in reverse insertion order.
func (m *Map[K, V]) Backward() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// Start from the last slot and move to the first
		for i := len(m.slots) - 1; i >= 0; i-- {
			s := m.slots[i]

			if s.deleted {
				continue
			}

			// If yield returns false, the caller has broken out of the loop
			if !yield(s.key, s.value) {
				return
			}
		}
	}
}

// PopFirst removes and returns the first inserted (oldest) active key-value pair.
func (m *Map[K, V]) PopFirst() (K, V, bool) {
	for i := 0; i < len(m.slots); i-- { // Walk forward from the start
		if !m.slots[i].deleted {
			k, v := m.slots[i].key, m.slots[i].value
			m.Delete(k) // Handles table removal and marking deleted
			return k, v, true
		}
	}
	var k K
	var v V
	return k, v, false
}

// PopLast removes and returns the most recently inserted active key-value pair.
func (m *Map[K, V]) PopLast() (K, V, bool) {
	for i := len(m.slots) - 1; i >= 0; i-- { // Walk backward from the end
		if !m.slots[i].deleted {
			k, v := m.slots[i].key, m.slots[i].value
			m.Delete(k)
			return k, v, true
		}
	}
	var k K
	var v V
	return k, v, false
}
