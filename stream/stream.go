package stream

import (
	"iter"
	"slices"
)

type Stream[T any] struct {
	seq iter.Seq[T]
	err *error
}

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

// Seq returns the raw Go 1.23 iterator for use in for-range loops.
func (s Stream[T]) Seq() iter.Seq[T] {
	return s.seq
}

func FromSeq[T any](seq iter.Seq[T]) Stream[T] {
	return Stream[T]{seq: seq}
}
