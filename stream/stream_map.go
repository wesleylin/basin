package stream

import "iter"

// basin.Map(myStream, func(int) string)
func Map[T, R any](s Stream[T], fn func(T) R) Stream[R] {
	return Stream[R]{
		err: s.err,
		seq: func(yield func(R) bool) {
			for v := range s.seq {
				if !yield(fn(v)) {
					return
				}
			}
		}}
}

func MapErr[T any](s Stream[T], fn func(T) (T, error)) Stream[T] {
	return Stream[T]{
		err: s.err,
		seq: func(yield func(T) bool) {
			for v := range s.seq {
				mapped, err := fn(v)
				if err != nil {
					// if err is found, set err as return value and terminate
					*s.err = err
					return
				}
				if !yield(mapped) {
					return
				}
			}
		}}
}

// FlatMap transforms T into an iterator of R, then flattens them into a single Stream[R].
func FlatMap[T, R any](s Stream[T], fn func(T) iter.Seq[R]) Stream[R] {
	return Stream[R]{
		err: s.err, // Preserve error
		seq: func(yield func(R) bool) {
			for v := range s.seq {
				// Circuit Breaker
				if s.err != nil && *s.err != nil {
					return
				}

				subSeq := fn(v)

				//Flatten the sub-sequence into the main yield
				for subItem := range subSeq {
					if !yield(subItem) {
						return
					}
				}
			}
		},
	}
}

// Fold collapses a Stream[T] into a single value of type U.
// It requires an initial value (the "seed") and a function to accumulate results.
func Fold[T any, U any](s Stream[T], initial U, fn func(U, T) U) (U, error) {
	acc := initial
	for v := range s.seq {
		// Check for upstream errors before each step
		if s.err != nil && *s.err != nil {
			return initial, *s.err
		}
		acc = fn(acc, v)
	}
	// Final check for errors that might have occurred at the very end of the sequence
	return acc, s.check()
}
