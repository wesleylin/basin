package stream

import (
	"iter"
	"slices"
)

type Stream[T any] struct {
	seq iter.Seq[T]
	err *error
}

// constructors

// New wraps a standard Go 1.23 iterator into a Basin Stream.
// It requires a pointer to an error so the stream remains "Error-Aware."
func New[T any](seq iter.Seq[T], errPtr *error) Stream[T] {
	if errPtr == nil {
		// Safety: If the user doesn't provide a pointer,
		// we give them a local one so the stream doesn't crash.
		var err error
		errPtr = &err
	}
	return Stream[T]{
		seq: seq,
		err: errPtr,
	}
}

func FromSeq[T any](seq iter.Seq[T]) Stream[T] {
	// err is instantiated here
	var err error
	return Stream[T]{seq: seq, err: &err}
}

// Seq returns the raw Go 1.23 iterator for use in for-range loops.
func (s Stream[T]) Seq() iter.Seq[T] {
	return s.seq
}

// FromSlice creates a Stream from a standard Go slice.
// Since slices are in-memory, the error pointer will stay nil
// unless a later operation (like MapErr) trips it.
func FromSlice[T any](items []T) Stream[T] {
	var err error
	return Stream[T]{
		err: &err,
		seq: func(yield func(T) bool) {
			for _, v := range items {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Filtering functions Filter, Take, Skip

// Filter creates a lazy iterator that only yields matching items.
func (s Stream[T]) Filter(fn func(T) bool) Stream[T] {
	return Stream[T]{
		err: s.err,
		seq: func(yield func(T) bool) {
			for v := range s.seq {
				if fn(v) {
					if !yield(v) {
						return
					}
				}
			}
		}}
}

// Take limits the number of items yielded.
func (s Stream[T]) Take(n int) Stream[T] {
	return Stream[T]{
		err: s.err,
		seq: func(yield func(T) bool) {
			count := 0
			for v := range s.seq {
				if count >= n || !yield(v) {
					return
				}
				count++
			}
		}}
}

func (s Stream[T]) Skip(n int) Stream[T] {
	return Stream[T]{
		err: s.err,
		seq: func(yield func(T) bool) {
			skipped := 0
			for v := range s.seq {
				// exit early if there's an error
				if s.err != nil && *s.err != nil {
					return
				}

				if skipped < n {
					skipped++
					continue
				}

				if !yield(v) {
					return
				}
			}
		}}
}

// Short circuting functions First, Any, All

// First returns the first element of the stream.
func (s Stream[T]) First() (T, error) {
	var zero T
	// range s.seq will execute the iterator
	for v := range s.seq {
		// Even for the first item, we check if the source failed
		if s.err != nil && *s.err != nil {
			return zero, *s.err
		}
		return v, nil
	}
	return zero, s.check()
}

// Any returns true if any element of the stream matches the predicate.
func (s Stream[T]) Any(fn func(T) bool) (bool, error) {
	for v := range s.seq {
		if s.err != nil && *s.err != nil {
			return false, *s.err
		}
		if fn(v) {
			// short circuit and return true
			return true, nil
		}
	}
	return false, s.check()
}

// All returns true if all elements of the stream match the predicate.
// note if the stream is empty, All returns true.
func (s Stream[T]) All(fn func(T) bool) (bool, error) {
	for v := range s.seq {
		if s.err != nil && *s.err != nil {
			return false, *s.err
		}
		if !fn(v) {
			return false, nil
		}
	}
	return true, s.check()
}

// Terminal functions Collect, Count, and ForEach

// Count counts the number of items in the stream.
func (s Stream[T]) Count() (int, error) {
	n := 0
	for range s.seq {
		n++
	}

	if s.err != nil && *s.err != nil {
		return 0, *s.err
	}

	return n, nil
}

// Collect gathers all items into a slice and returns any error encountered.
func (s Stream[T]) Collect() ([]T, error) {
	items := slices.Collect(s.seq)

	// the s.err is not potentially set until after Collect is called
	// so check s.err != nil after
	if s.err != nil && *s.err != nil {
		return nil, *s.err
	}

	return items, nil
}

// check is a private helper to keep things DRY
func (s Stream[T]) check() error {
	if s.err != nil && *s.err != nil {
		return *s.err
	}
	return nil
}
