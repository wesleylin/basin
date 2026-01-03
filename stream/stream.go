package stream

import (
	"iter"
	"slices"
)

type Stream[T any] struct {
	seq iter.Seq[T]
	err *error
}

// Filter creates a lazy iterator that only yields matching items.
func (s Stream[T]) Filter(fn func(T) bool) Stream[T] {
	return Stream[T]{seq: func(yield func(T) bool) {
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
	return Stream[T]{seq: func(yield func(T) bool) {
		count := 0
		for v := range s.seq {
			if count >= n || !yield(v) {
				return
			}
			count++
		}
	}}
}

// Collect gathers all items into a slice and returns any error encountered.
func (s Stream[T]) Collect() ([]T, error) {
	items := slices.Collect(s.seq)

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
