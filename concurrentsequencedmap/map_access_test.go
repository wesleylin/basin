package concurrentsequencedmap_test

import (
	"testing"

	"github.com/wesleylin/basin/concurrentsequencedmap"
)

func TestMap_KeysAndValues(t *testing.T) {
	m := concurrentsequencedmap.New[string, int]()

	// Use a small set of predictable data
	data := []struct {
		k string
		v int
	}{
		{"apple", 100},
		{"banana", 200},
		{"cherry", 300},
	}

	for _, d := range data {
		m.Put(d.k, d.v)
	}

	// 1. Test Keys()
	t.Run("Keys", func(t *testing.T) {
		i := 0
		for k := range m.Keys() {
			if k != data[i].k {
				t.Errorf("Keys() order mismatch at index %d: expected %s, got %s", i, data[i].k, k)
			}
			i++
		}
		if i != len(data) {
			t.Errorf("Keys() yielded %d items, expected %d", i, len(data))
		}
	})

	// 2. Test Values()
	t.Run("Values", func(t *testing.T) {
		i := 0
		for v := range m.Values() {
			if v != data[i].v {
				t.Errorf("Values() order mismatch at index %d: expected %d, got %d", i, data[i].v, v)
			}
			i++
		}
		if i != len(data) {
			t.Errorf("Values() yielded %d items, expected %d", i, len(data))
		}
	})
}
