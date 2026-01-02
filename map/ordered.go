package ordered

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
		slots: make([]entry[K, V], 0),
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
