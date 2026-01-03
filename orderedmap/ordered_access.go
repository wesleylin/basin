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
