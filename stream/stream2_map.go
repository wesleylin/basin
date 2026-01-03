package stream

import "iter"

// Map2 transforms K, V into new types NK, NV.
func Map2[K, V, NK, NV any](s Stream2[K, V], fn func(K, V) (NK, NV)) Stream2[NK, NV] {
	return Stream2[NK, NV]{
		err: s.err,
		seq: func(yield func(NK, NV) bool) {
			for k, v := range s.seq {
				if !yield(fn(k, v)) {
					return
				}
			}
		},
	}
}

// Map2Err allows transformation with error handling.
// If an error occurs, it updates the shared error pointer and halts.
func Map2Err[K, V any](s Stream2[K, V], fn func(K, V) (K, V, error)) Stream2[K, V] {
	return Stream2[K, V]{
		err: s.err,
		seq: func(yield func(K, V) bool) {
			for k, v := range s.seq {
				nk, nv, err := fn(k, v)
				if err != nil {
					if s.err != nil {
						*s.err = err
					}
					return
				}
				if !yield(nk, nv) {
					return
				}
			}
		},
	}
}

// FlatMap2 transforms a single K, V pair into a sequence of NK, NV pairs,
// then flattens them into the main stream.
func FlatMap2[K, V, NK, NV any](s Stream2[K, V], fn func(K, V) iter.Seq2[NK, NV]) Stream2[NK, NV] {
	return Stream2[NK, NV]{
		err: s.err,
		seq: func(yield func(NK, NV) bool) {
			for k, v := range s.seq {
				// Circuit Breaker: check if a previous step errored out
				if s.err != nil && *s.err != nil {
					return
				}

				subSeq := fn(k, v)
				for nk, nv := range subSeq {
					if !yield(nk, nv) {
						return
					}
				}
			}
		},
	}
}
