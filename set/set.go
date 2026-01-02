package set

import "iter"

// entry only needs the key and the tombstone
type entry[K comparable] struct {
	key     K
	deleted bool
}

type Set[K comparable] struct {
	table        map[K]int
	slots        []entry[K]
	deletedCount int
}

func New[K comparable]() *Set[K] {
	return &Set[K]{
		table: make(map[K]int),
		slots: make([]entry[K], 0, 8),
	}
}

func NewWithCapacity[K comparable](cap int) *Set[K] {
	return &Set[K]{
		table: make(map[K]int, cap),
		slots: make([]entry[K], 0, cap),
	}
}

// Add inserts an item into the set. Returns true if it was newly added.
func (s *Set[K]) Add(key K) bool {
	if idx, exists := s.table[key]; exists {
		if s.slots[idx].deleted {
			s.slots[idx].deleted = false
			s.deletedCount--
			return true
		}
		return false // Already existed and wasn't deleted
	}

	s.table[key] = len(s.slots)
	s.slots = append(s.slots, entry[K]{key: key})
	return true
}

func (s *Set[K]) All() iter.Seq[K] {
	return func(yield func(K) bool) {
		for _, e := range s.slots {
			if !e.deleted {
				if !yield(e.key) {
					return
				}
			}
		}
	}
}

func (s *Set[K]) Delete(key K) {
	idx, exists := s.table[key]
	if !exists {
		return
	}

	delete(s.table, key)
	s.slots[idx].deleted = true

	// GC optimization: Zero out the key
	var zero K
	s.slots[idx].key = zero

	s.deletedCount++

	if s.deletedCount*2 > len(s.slots) {
		s.compact()
	}
}

func (s *Set[K]) compact() {
	newSlots := make([]entry[K], 0, len(s.slots)-s.deletedCount)
	for _, e := range s.slots {
		if !e.deleted {
			s.table[e.key] = len(newSlots)
			newSlots = append(newSlots, e)
		}
	}
	s.slots = newSlots
	s.deletedCount = 0
}

func (s *Set[K]) Has(key K) bool {
	_, exists := s.table[key]
	return exists
}

func (s *Set[K]) Len() int {
	return len(s.table)
}

func (s *Set[K]) Clear() {
	clear(s.table)
	s.slots = s.slots[:0]
	s.deletedCount = 0
}
