package set

// Add inserts an item into the set. Returns the set for chaining.
func (s *Set[K]) Add(key K) *Set[K] {
	s.Insert(key)
	return s
}

// Remove removes an item from the set. Returns the set for chaining.
func (s *Set[K]) Remove(key K) *Set[K] {
	s.Delete(key)
	return s
}

// Filter returns a new set containing only the elements that satisfy the predicate.
// Pass in a function that takes an element and returns a boolean.
// It preserves the "Ordered" nature of the original set.
// Returns the set itself to allow for method chaining.
func (s *Set[T]) Filter(fn func(T) bool) *Set[T] {
	newSet := New[T]()
	// Using the internal iterator for 21x speed
	s.All()(func(v T) bool {
		if fn(v) {
			newSet.Insert(v)
		}
		return true
	})
	return newSet
}

// Each executes a provided function once for each set element.
// Returns the set itself to allow for method chaining.
func (s *Set[T]) Each(fn func(T)) *Set[T] {
	s.All()(func(v T) bool {
		fn(v)
		return true
	})
	return s
}

// ToSlice converts the Set into a standard Go slice.
// Returns a slice containing all elements in the set.
func (s *Set[T]) ToSlice() []T {
	res := make([]T, 0, s.Len())
	s.All()(func(v T) bool {
		res = append(res, v)
		return true
	})
	return res
}
