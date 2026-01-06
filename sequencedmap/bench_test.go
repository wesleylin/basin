package sequencedmap_test

import (
	"fmt"
	"testing"

	orderedmap "github.com/wesleylin/basin/sequencedmap"
)

// --- PUT BENCHMARKS ---

func BenchmarkPut_Basin(b *testing.B) {
	for _, n := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("Size-%d", n), func(b *testing.B) {
			m := orderedmap.New[int, int]()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Put(i%n, i)
			}
		})
	}
}

func BenchmarkPut_StdMap(b *testing.B) {
	for _, n := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("Size-%d", n), func(b *testing.B) {
			m := make(map[int]int)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m[i%n] = i
			}
		})
	}
}

// --- GET BENCHMARKS ---

func BenchmarkGet_Basin(b *testing.B) {
	m := orderedmap.New[int, int]()
	size := 1000
	for i := 0; i < size; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Get(i % size)
	}
}

func BenchmarkGet_StdMap(b *testing.B) {
	m := make(map[int]int)
	size := 1000
	for i := 0; i < size; i++ {
		m[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[i%size]
	}
}

// --- ITERATION BENCHMARKS ---

func BenchmarkIterate_Basin(b *testing.B) {
	m := orderedmap.New[int, int]()
	size := 10000
	for i := 0; i < size; i++ {
		m.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range m.Keys() {
			// Basin walks a slice (CPU-friendly)
		}
	}
}

func BenchmarkIterate_StdMap(b *testing.B) {
	m := make(map[int]int)
	size := 10000
	for i := 0; i < size; i++ {
		m[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range m {
			// StdMap walks a hash table (Randomized/Heavier)
		}
	}
}

// BenchmarkMemoryPressure tests the memory pressure of Basin's ordered map
// when there are many deletions causing "holes" in the internal slice.
// It compares iteration performance before and after compacting the map.
func BenchmarkMemoryPressure(b *testing.B) {
	size := 100000
	m := orderedmap.New[int, int]()

	// 1. Fill the map
	for i := 0; i < size; i++ {
		m.Put(i, i)
	}

	// 2. Delete 99% of the items
	// This leaves "holes" in the slots slice
	for i := 0; i < size-1000; i++ {
		m.Delete(i)
	}

	b.Run("BeforeCompact", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Iterate over the map with 99,000 holes
			for range m.Keys() {
			}
		}
	})

	b.Run("AfterCompact", func(b *testing.B) {
		m.Compact() // Clean up the holes
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Iterate over the dense map (only 1,000 active items)
			for range m.Keys() {
			}
		}
	})
}

// Compare memory pressure of Basin vs standard map under heavy put/delete
func BenchmarkMemoryPressure_Comparison(b *testing.B) {
	const size = 100000

	b.Run("Basin_NoCompact", func(b *testing.B) {
		m := orderedmap.New[int, int]()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// We put and then delete.
			// In a standard map, memory stays flat.
			// In Basin, m.slots grows to the size of b.N!
			m.Put(i, i)
			m.Delete(i)
		}
	})

	b.Run("StdMap", func(b *testing.B) {
		m := make(map[int]int)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			m[i] = i
			delete(m, i)
		}
	})
}
