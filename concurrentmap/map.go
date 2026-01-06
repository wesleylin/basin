package concurrentmap

import (
	"sync"

	"github.com/wesleylin/basin/internal/hash"
)

// ShardCount must be a power of 2 for the bitwise mask to work
const ShardCount = 256
const shardMask = ShardCount - 1

type Map[K comparable, V any] struct {
	shards []*shard[K, V]
}

type shard[K comparable, V any] struct {
	sync.RWMutex
	data map[K]V
}

func New[K comparable, V any]() *Map[K, V] {
	m := &Map[K, V]{
		shards: make([]*shard[K, V], ShardCount),
	}
	for i := 0; i < ShardCount; i++ {
		m.shards[i] = &shard[K, V]{
			data: make(map[K]V),
		}
	}
	return m
}

// getShard picks the correct bucket based on the key's hash
func (m *Map[K, V]) getShard(key K) *shard[K, V] {
	h := hash.Maphash(key)
	return m.shards[h&uint64(shardMask)]
}

// Set adds or updates a key-value pair.
func (m *Map[K, V]) Put(key K, value V) {
	s := m.getShard(key)
	s.Lock()
	s.data[key] = value
	s.Unlock()
}

// Get retrieves a value from the map.
func (m *Map[K, V]) Get(key K) (V, bool) {
	s := m.getShard(key)
	s.RLock()
	val, ok := s.data[key]
	s.RUnlock()
	return val, ok
}

// Delete removes a key from the map.
func (m *Map[K, V]) Delete(key K) {
	s := m.getShard(key)
	s.Lock()
	delete(s.data, key)
	s.Unlock()
}

// Pop deletes and returns the value (Atomic Get + Delete)
// use this instead of using both Get and Delete to ensure atomicity
func (m *Map[K, V]) Pop(key K) (V, bool) {
	s := m.getShard(key)
	s.Lock()
	defer s.Unlock()
	val, ok := s.data[key]
	if ok {
		delete(s.data, key)
	}
	return val, ok
}
