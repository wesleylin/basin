package set

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

// Has checks if an item is in the set
func (s *Set[K]) Has(key K) bool {
	_, exists := s.table[key]
	return exists
}
