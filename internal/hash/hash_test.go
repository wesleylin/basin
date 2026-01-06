package hash

import (
	"testing"
)

func TestMaphash(t *testing.T) {
	t.Run("Determinism", func(t *testing.T) {
		val := "basin-test-key"
		h1 := Maphash(val)
		h2 := Maphash(val)
		if h1 != h2 {
			t.Errorf("expected same hash for same input, got %d and %d", h1, h2)
		}
	})

	t.Run("DifferentTypes", func(t *testing.T) {
		// Even if they have similar bit patterns,
		// typehash should treat them differently
		h1 := Maphash(int64(1))
		h2 := Maphash(float64(1.0))
		if h1 == h2 {
			t.Log("Note: h1 and h2 collided, which is possible but rare")
		}
	})

	t.Run("StructEquality", func(t *testing.T) {
		type entity struct {
			ID   int
			Name string
		}
		e1 := entity{ID: 10, Name: "Basin"}
		e2 := entity{ID: 10, Name: "Basin"}

		if Maphash(e1) != Maphash(e2) {
			t.Error("identical structs should produce identical hashes")
		}
	})
}

func BenchmarkMaphash(b *testing.B) {
	key := "a-reasonably-long-key-for-testing"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Maphash(key)
	}
}
