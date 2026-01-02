package set

import (
	"testing"
)

func BenchmarkBasinSet_Add(b *testing.B) {
	s := New[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Insert(i % 1000)
	}
}

func BenchmarkGoMap_Add(b *testing.B) {
	m := make(map[int]struct{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[i%1000] = struct{}{}
	}
}

func BenchmarkBasinSet_Iterate(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Insert(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range s.All() {
			// iterate
		}
	}
}

func BenchmarkGoMap_Iterate(b *testing.B) {
	m := make(map[int]struct{})
	for i := 0; i < 1000; i++ {
		m[i] = struct{}{}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range m {
			// iterate
		}
	}
}
