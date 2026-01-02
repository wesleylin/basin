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

// Has checks if the set contains the given key.
func (s *Set[K]) Has(key K) bool {
	_, exists := s.table[key]
	return exists
}

// Insert inserts an item into the set. Returns true if it was newly added.
func (s *Set[K]) Insert(key K) bool {
	if _, exists := s.table[key]; exists {
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

// Delete removes an item from the set. Returns true if the item was present.
func (s *Set[K]) Delete(key K) bool {
	idx, exists := s.table[key]
	if !exists {
		return false
	}

	// 1. remove from table to be added later
	delete(s.table, key)
	// 2. mark tombstone in slots
	s.slots[idx].deleted = true

	// GC optimization: Zero out the key
	var zero K
	s.slots[idx].key = zero

	s.deletedCount++

	// periodic cleanup if more than half are deleted
	if s.deletedCount*2 > len(s.slots) {
		s.compact()
	}
	return true
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

// helper functions
func (s *Set[K]) Len() int {
	return len(s.table)
}

func (s *Set[K]) Clear() {
	clear(s.table)
	s.slots = s.slots[:0]
	s.deletedCount = 0
}

// set algebra operations

func Union[K comparable](a, b *Set[K]) iter.Seq[K] {
	return func(yield func(K) bool) {
		// Just streaming from the existing 'slots' of a
		for item := range a.All() {
			if !yield(item) {
				return
			}
		}
		// Just streaming from b, checking against a's hash table
		for item := range b.All() {
			if !a.Has(item) {
				if !yield(item) {
					return
				}
			}
		}
	}
}

func Intersect[K comparable](a, b *Set[K]) iter.Seq[K] {
	return func(yield func(K) bool) {
		// Optimization: iterate over the smaller set to minimize lookups
		target, check := a, b
		if b.Len() < a.Len() {
			target, check = b, a
		}

		for item := range target.All() {
			if check.Has(item) {
				if !yield(item) {
					return
				}
			}
		}
	}
}

func Difference[K comparable](a, b *Set[K]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for item := range a.All() {
			if !b.Has(item) {
				if !yield(item) {
					return
				}
			}
		}
	}
}
