package concurrentsequencedmap

import (
	"sync"

	orderedmap "github.com/wesleylin/basin/sequencedmap"
)

type Map[K comparable, V any] struct {
	shards   []*shard[K, V]
	sequence uint64 // The Global Atomic Clock
}

type entry[V any] struct {
	value V
	seq   uint64 // The "timestamp" for ordering
}

type shard[K comparable, V any] struct {
	sync.RWMutex
	data orderedmap.Map[K, V] // Now stores an entry struct instead of raw V
}
