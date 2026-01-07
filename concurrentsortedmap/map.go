package concurrentsortedmap

import (
	"cmp"
	"sync"

	"github.com/wesleylin/basin/sortedmap"
)

const shardCount = 256

type ConcurrentSortedMap[K cmp.Ordered, V any] struct {
	// shardMask is (shardCount - 1). Used for fast bitwise indexing.
	shardMask uint64
	shards    [shardCount]*shard[K, V]
}

type shard[K cmp.Ordered, V any] struct {
	// We embed a RWMutex so each shard can be locked independently
	sync.RWMutex
	data sortedmap.SortedMap[K, V]
}

func New[K cmp.Ordered, V any]() *ConcurrentSortedMap[K, V] {
	cm := &ConcurrentSortedMap[K, V]{
		shardMask: shardCount - 1,
	}

	for i := 0; i < shardCount; i++ {
		cm.shards[i] = &shard[K, V]{
			data: *sortedmap.New[K, V](),
		}
	}
	return cm
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
