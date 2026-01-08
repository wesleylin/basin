package stream

import (
	"cmp"
	"errors"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Sum calculates the total of all elements in the stream.
func Sum[T Number](s Stream[T]) (T, error) {
	return Fold(s, T(0), func(acc T, v T) T {
		return acc + v
	})
}

// Max returns the largest element in the stream.
// Returns an error if the stream is empty.
func Max[T cmp.Ordered](s Stream[T]) (T, error) {
	return s.Reduce(func(a, b T) T {
		if a > b {
			return a
		}
		return b
	})
}

// Min returns the smallest element in the stream.
// Returns an error if the stream is empty.
func Min[T cmp.Ordered](s Stream[T]) (T, error) {
	return s.Reduce(func(a, b T) T {
		if a < b {
			return a
		}
		return b
	})
}

// Average calculates the arithmetic mean.
func Average[T Number](s Stream[T]) (float64, error) {
	var sum float64
	var count int

	err := s.ForEach(func(v T) {
		sum += float64(v)
		count++
	})

	if err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("cannot calculate average of empty stream")
	}

	return sum / float64(count), nil
}
