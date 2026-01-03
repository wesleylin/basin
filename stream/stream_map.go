package stream

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
