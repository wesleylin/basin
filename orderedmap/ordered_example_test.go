package orderedmap

import (
	"testing"
)

type Animal struct {
	Name string
	Type string
}

// In Go 1.24, we use [Animal any] for the alias if we want it to be generic,
// but here we are pinning it to the specific 'Animal' struct.
type ZooMap = Map[string, Animal]

func TestOrderMapExample(t *testing.T) {
	// Initialize the map
	zoo := New[string, Animal]()

	// Using the Fluent API we designed:
	// 1. We use 'Set' for fluent Map insertion.
	// 2. Struct fields need strings in quotes.
	// 3. Keys (strings) must be passed separately from the Value (Animal).
	zoo.Set("kyle", Animal{"Kyle", "Kangaroo"}).
		Set("sam", Animal{"Sam", "Tiger"})

	// Iterate and verify order
	expectedOrder := []string{"kyle", "sam"}
	index := 0
	for key, animal := range zoo.All() {
		if key != expectedOrder[index] {
			t.Errorf("Expected key %s at index %d, got %s", expectedOrder[index], index, key)
		}
		index++
		_ = animal // In a real test, we might verify animal fields too
	}

	if index != len(expectedOrder) {
		t.Errorf("Expected %d items, got %d", len(expectedOrder), index)
	}
}

func TestOrderMapExample2(t *testing.T) {
	// Initialize the map
	zoo := New[string, Animal]()

	// can chain most calls calls
	zoo = zoo.Set("kyle", Animal{"Kyle", "Kangaroo"}).
		Set("sam", Animal{"Sam", "Tiger"})

	zooStream := zoo.Stream2()

	zooStream = zooStream.Filter(func(k string, a Animal) bool {
		return a.Type == "Tiger" || a.Type == "Lion"
	})

	// zooStream = stream.Map

	// zooStream = zooStream.Map(func(k string, a Animal) (string, Animal) {
	// 	a.Name = "Big " + a.Name
	// 	return k, a
	// }
}
