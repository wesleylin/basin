package stream

import "fmt"

// Last consumes the stream and returns the very last element.
// Returns an error if the stream is empty or if an upstream error occurs.
func (s Stream[T]) Last() (T, error) {
	var last T
	var found bool

	err := s.ForEach(func(v T) {
		last = v
		found = true
	})

	if err != nil {
		return last, err
	}
	if !found {
		var zero T
		return zero, fmt.Errorf("cannot get last element of empty stream")
	}
	return last, nil
}

// ToMap takes a stream of Entry structs and collects them into a Go map.
func ToMap[K comparable, V any](s Stream[Entry[K, V]]) (map[K]V, error) {
	res := make(map[K]V)
	err := s.ForEach(func(e Entry[K, V]) {
		res[e.Key] = e.Value
	})
	return res, err
}

// GroupBy organizes elements into a map of slices based on a key extraction function.
func GroupBy[T any, K comparable](s Stream[T], keyFn func(T) K) (map[K][]T, error) {
	res := make(map[K][]T)
	err := s.ForEach(func(v T) {
		key := keyFn(v)
		res[key] = append(res[key], v)
	})
	return res, err
}
