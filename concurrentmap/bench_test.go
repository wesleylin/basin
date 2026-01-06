package concurrentmap

import (
	"sync"
	"testing"
)

// BenchmarkParallelPut measures write performance across multiple goroutines.
func BenchmarkParallelPut(b *testing.B) {
	m := New[int, int]()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Put(i, i)
			i++
		}
	})
}

// BenchmarkParallelGet measures read performance.
// This should be extremely fast due to RLock.
func BenchmarkParallelGet(b *testing.B) {
	m := New[int, int]()
	// Pre-fill the map
	for i := 0; i < 10000; i++ {
		m.Put(i, i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = m.Get(i % 10000)
			i++
		}
	})
}

// BenchmarkContention simulates a real-world mix of reads and writes.
func BenchmarkContention(b *testing.B) {
	m := New[int, int]()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%4 == 0 {
				m.Put(i, i)
			} else {
				_, _ = m.Get(i)
			}
			i++
		}
	})
}

// Comparison Benchmark: Native Map with a single Mutex
// This shows why your Basin ConcurrentMap is superior.
func BenchmarkSingleMutexMap(b *testing.B) {
	var mu sync.RWMutex
	m := make(map[int]int)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%4 == 0 {
				mu.Lock()
				m[i] = i
				mu.Unlock()
			} else {
				mu.RLock()
				_ = m[i]
				mu.RUnlock()
			}
			i++
		}
	})
}
