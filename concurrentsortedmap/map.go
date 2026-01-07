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

// shard has a pointer receiver because it must lock the mutex
func (s *shard[K, V]) Set(key K, value V) {
	s.Lock()
	defer s.Unlock()
	s.data.Set(key, value)
}

// Even Get needs a pointer receiver because it uses PathHints (which change)
func (s *shard[K, V]) Get(key K) (V, bool) {
	s.RLock()
	defer s.RUnlock()
	return s.data.Get(key)
}
