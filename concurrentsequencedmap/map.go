package concurrentsequencedmap

import (
	"sync"
	"sync/atomic"

	"github.com/wesleylin/basin/internal/hash"
	"github.com/wesleylin/basin/sequencedmap"
)

const ShardCount = 256
const shardMask = ShardCount - 1

type globalEntry[V any] struct {
	value V
	seq   uint64
}

type Map[K comparable, V any] struct {
	shards   [ShardCount]*shard[K, V]
	sequence uint64 // The Global Atomic Clock
}

type entry[V any] struct {
	value V
	seq   uint64 // The "timestamp" for ordering
}

type shard[K comparable, V any] struct {
	sync.RWMutex
	// Now stores point to entry struct instead of raw V
	data *sequencedmap.Map[K, globalEntry[V]]
}

func New[K comparable, V any]() *Map[K, V] {
	m := &Map[K, V]{}
	for i := 0; i < ShardCount; i++ {
		// initialize shards with a sequencedmap (insert order map)
		m.shards[i] = &shard[K, V]{
			data: sequencedmap.New[K, globalEntry[V]](),
		}
	}
	return m
}

func (m *Map[K, V]) getShard(key K) *shard[K, V] {
	h := hash.Maphash(key)
	return m.shards[h&uint64(shardMask)]
}

// Put adds a value and assigns the next global sequence order
func (m *Map[K, V]) Put(key K, value V) {
	// 1. Grab a global ticket
	seq := atomic.AddUint64(&m.sequence, 1)

	s := m.getShard(key)
	s.Lock()
	defer s.Unlock()

	// 2. Store with the ticket
	s.data.Put(key, globalEntry[V]{value: value, seq: seq})
}

// Get retrieves a value from the correct shard
func (m *Map[K, V]) Get(key K) (V, bool) {
	s := m.getShard(key)
	s.RLock()
	defer s.RUnlock()

	res, ok := s.data.Get(key)
	return res.value, ok
}

// Delete removes a key while maintaining thread safety.
func (m *Map[K, V]) Delete(key K) {
	s := m.getShard(key)
	s.Lock()
	defer s.Unlock()

	s.data.Delete(key)
}

// Len returns the total number of elements across all shards.
func (m *Map[K, V]) Len() int {
	var total int
	for i := 0; i < ShardCount; i++ {
		m.shards[i].RLock()
		total += m.shards[i].data.Len()
		m.shards[i].RUnlock()
	}
	return total
}
