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
