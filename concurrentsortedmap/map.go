package concurrentsortedmap

import (
	"cmp"
	"sync"

	"github.com/wesleylin/basin/internal/hash"
	"github.com/wesleylin/basin/sortedmap"
)

const shardCount = 256
const shardMask = shardCount - 1

type Map[K cmp.Ordered, V any] struct {
	// shardMask is (shardCount - 1). Used for fast bitwise indexing.
	shardMask uint64
	shards    [shardCount]*shard[K, V]
}

type shard[K cmp.Ordered, V any] struct {
	// We embed a RWMutex so each shard can be locked independently
	sync.RWMutex
	data sortedmap.SortedMap[K, V]
}

func New[K cmp.Ordered, V any]() *Map[K, V] {
	cm := &Map[K, V]{
		shardMask: shardCount - 1,
	}

	for i := 0; i < shardCount; i++ {
		cm.shards[i] = &shard[K, V]{
			data: *sortedmap.New[K, V](),
		}
	}
	return cm
}

func (m *Map[K, V]) getShard(key K) *shard[K, V] {
	h := hash.Maphash(key)
	return m.shards[h&uint64(shardMask)]
}

// Get retrieves a value from the correct shard
func (m *Map[K, V]) Get(key K) (V, bool) {
	s := m.getShard(key)
	s.RLock()
	defer s.RUnlock()

	res, ok := s.data.Get(key)
	return res, ok
}

// Put inserts or updates a value
func (m *Map[K, V]) Put(key K, value V) bool {
	s := m.getShard(key)
	s.Lock()
	defer s.Unlock()

	_, ok := s.data.Put(key, value)
	return ok
}

// Delete removes a key while maintaining thread safety.
func (m *Map[K, V]) Delete(key K) {
	s := m.getShard(key)
	s.Lock()
	defer s.Unlock()

	s.data.Delete(key)
}
