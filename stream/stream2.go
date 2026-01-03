package stream

import (
	"iter"
)

// Stream2 represents a sequence of Key-Value pairs with a "Live Wire" error pointer.
type Stream2[K, V any] struct {
	seq iter.Seq2[K, V]
	err *error
}

// Pair is a simple container for when users want to collect Stream2 into a slice.
type Pair[K, V any] struct {
	Key   K
	Value V
}

// --- Constructors ---

// New2 wraps a standard Go 1.23 K-V iterator into a Basin Stream2.
func New2[K, V any](seq iter.Seq2[K, V], errPtr *error) Stream2[K, V] {
	if errPtr == nil {
		var err error
		errPtr = &err
	}
	return Stream2[K, V]{
		seq: seq,
		err: errPtr,
	}
}

// FromMap creates a Stream2 from a standard Go map.
func FromMap[K comparable, V any](m map[K]V) Stream2[K, V] {
	var err error
	return Stream2[K, V]{
		err: &err,
		seq: func(yield func(K, V) bool) {
			for k, v := range m {
				if !yield(k, v) {
					return
				}
			}
		},
	}
}

// --- Bridges (Stream2 -> Stream) ---

// Keys returns a Stream containing only the keys.
func (s Stream2[K, V]) Keys() Stream[K] {
	return Stream[K]{
		err: s.err,
		seq: func(yield func(K) bool) {
			for k := range s.seq {
				if !yield(k) {
					return
				}
			}
		},
	}
}

// Values returns a Stream containing only the values.
func (s Stream2[K, V]) Values() Stream[V] {
	return Stream[V]{
		err: s.err,
		seq: func(yield func(V) bool) {
			for _, v := range s.seq {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// --- Filtering Functions ---

func (s Stream2[K, V]) Filter(fn func(K, V) bool) Stream2[K, V] {
	return Stream2[K, V]{
		err: s.err,
		seq: func(yield func(K, V) bool) {
			for k, v := range s.seq {
				if fn(k, v) {
					if !yield(k, v) {
						return
					}
				}
			}
		},
	}
}

func (s Stream2[K, V]) Take(n int) Stream2[K, V] {
	return Stream2[K, V]{
		err: s.err,
		seq: func(yield func(K, V) bool) {
			count := 0
			for k, v := range s.seq {
				if count >= n || !yield(k, v) {
					return
				}
				count++
			}
		},
	}
}

// --- Standalone Transformations (Non-Methods for Type Inference) ---

// MapValues transforms the values (V -> R) while keeping the keys the same.
func MapValues[K, V, R any](s Stream2[K, V], fn func(V) R) Stream2[K, R] {
	return Stream2[K, R]{
		err: s.err,
		seq: func(yield func(K, R) bool) {
			for k, v := range s.seq {
				if !yield(k, fn(v)) {
					return
				}
			}
		},
	}
}

// MapErr2 is a fallible transformation for both key and value.
// If fn returns an error, the "Live Wire" trips and the stream stops.
func MapErr2[K, V, NK, NV any](s Stream2[K, V], fn func(K, V) (NK, NV, error)) Stream2[NK, NV] {
	return Stream2[NK, NV]{
		err: s.err,
		seq: func(yield func(NK, NV) bool) {
			for k, v := range s.seq {
				nk, nv, err := fn(k, v)
				if err != nil {
					*s.err = err
					return
				}
				if !yield(nk, nv) {
					return
				}
			}
		},
	}
}

// --- Terminal Functions ---

func (s Stream2[K, V]) Collect() ([]Pair[K, V], error) {
	var results []Pair[K, V]
	for k, v := range s.seq {
		results = append(results, Pair[K, V]{Key: k, Value: v})
	}

	if s.err != nil && *s.err != nil {
		return nil, *s.err
	}

	return results, nil
}

func (s Stream2[K, V]) Count() (int, error) {
	n := 0
	for range s.seq {
		n++
	}
	if s.err != nil && *s.err != nil {
		return 0, *s.err
	}
	return n, nil
}

// check is a private helper to keep things DRY
func (s Stream2[K, V]) check() error {
	if s.err != nil && *s.err != nil {
		return *s.err
	}
	return nil
}
