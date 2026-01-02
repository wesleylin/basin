package set

import (
	"math/rand"
	"testing"
)

func TestSet_Chaos(t *testing.T) {
	s := New[int]()
	truth := make(map[int]struct{})

	// We'll perform 100,000 operations
	for i := 0; i < 100000; i++ {
		val := rand.Intn(1000) // Keep the range small to force collisions/deletes
		op := rand.Intn(3)

		switch op {
		case 0: // ADD
			s.Add(val)
			truth[val] = struct{}{}
		case 1: // DELETE
			s.Delete(val)
			delete(truth, val)
		case 2: // CHECK
			_, existsInTruth := truth[val]
			if s.Has(val) != existsInTruth {
				t.Fatalf("Consistency error at op %d: Has(%d) mismatch", i, val)
			}
		}

		// Periodically verify size and iteration
		if i%1000 == 0 {
			if s.Len() != len(truth) {
				t.Fatalf("Size mismatch: Basin %d, Truth %d", s.Len(), len(truth))
			}
		}
	}

	// Final check: Do the iterators produce the same items?
	count := 0
	s.All()(func(v int) bool {
		if _, ok := truth[v]; !ok {
			t.Errorf("Iterator produced value %d not in truth map", v)
		}
		count++
		return true
	})

	if count != len(truth) {
		t.Errorf("Iterator count mismatch: Got %d, Want %d", count, len(truth))
	}
}
