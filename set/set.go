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
