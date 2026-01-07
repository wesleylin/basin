package concurrentsequencedmap

import "iter"

// All returns a globally ordered iterator across all 256 shards
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// 1. Snapshot the heads of all shards
		// To avoid holding locks for 40GB of iteration, we collect the ordered data
		shardIters := make([]func(func(K, globalEntry[V]) bool), ShardCount)
		for i := 0; i < ShardCount; i++ {
			shardIters[i] = m.shards[i].data.All()
		}

		// 2. Setup "Pull" iterators (Go 1.23 feature)
		// This lets us manually advance each shard's iterator
		next, stop := make([]func() (K, globalEntry[V], bool), ShardCount), make([]func(), ShardCount)
		for i := 0; i < ShardCount; i++ {
			next[i], stop[i] = iter.Pull2(shardIters[i])
			defer stop[i]()
		}

		// 3. Current "Head" of each shard
		keys := make([]K, ShardCount)
		entries := make([]globalEntry[V], ShardCount)
		active := make([]bool, ShardCount)

		for i := 0; i < ShardCount; i++ {
			keys[i], entries[i], active[i] = next[i]()
		}

		// 4. Merge Loop: Always pick the smallest global sequence ID
		for {
			bestShard := -1
			var minSeq uint64 = ^uint64(0)

			for i := 0; i < ShardCount; i++ {
				if active[i] && entries[i].seq < minSeq {
					minSeq = entries[i].seq
					bestShard = i
				}
			}

			if bestShard == -1 {
				break // All shards exhausted
			}

			// Yield the globally next value
			if !yield(keys[bestShard], entries[bestShard].value) {
				return
			}

			// Advance the shard we just used
			keys[bestShard], entries[bestShard], active[bestShard] = next[bestShard]()
		}
	}
}

// Keys returns an iterator for the map's keys in insertion order.
func (m *Map[K, V]) Keys() iter.Seq[K] {
	// TODO: possibly optimize to not call All to remove retrieving Value() as well
	return func(yield func(K) bool) {
		for k, _ := range m.All() {
			if !yield(k) {
				return
			}
		}
	}
}

// Values returns an iterator for the map's values in insertion order.
func (m *Map[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m.All() {
			if !yield(v) {
				return
			}
		}
	}
}
