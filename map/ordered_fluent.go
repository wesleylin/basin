package ordered

// GetOr returns the value for a key, or the provided default if it doesn't exist.
func (m *Map[K, V]) GetOr(key K, defaultVal V) V {
	val, exists := m.Get(key)
	if exists {
		return val
	}
	return defaultVal
}

// Set sets the value for a key in the map. It returns the ordered map for chaining.
func (m *Map[K, V]) Set(key K, val V) *Map[K, V] {
	m.Put(key, val)
	return m
}
