package concurrentsortedmap_test

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/wesleylin/basin/concurrentsortedmap"
)

func TestMap_All(t *testing.T) {
	// Initialize the map with 256 shards (assuming your New func does this)
	m := concurrentsortedmap.New[int, string]()

	// 1. Generate a large set of random data
	count := 10000
	expected := make([]int, count)
	for i := 0; i < count; i++ {
		val := rand.Intn(1000000)
		expected[i] = val
		m.Put(val, fmt.Sprintf("val-%d", val))
	}

	// 2. Create the "Ground Truth" by sorting our input slice
	slices.Sort(expected)
	// Remove duplicates since Put might overwrite
	expected = slices.Compact(expected)

	// 3. Collect results from the All() iterator
	var results []int
	for k, _ := range m.All() {
		results = append(results, k)
	}

	// 4. Validate results
	if len(results) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(results))
	}

	for i := 0; i < len(results); i++ {
		if results[i] != expected[i] {
			t.Errorf("At index %d: expected %d, got %d", i, expected[i], results[i])
			break // Break early to avoid flooding logs
		}
	}
}

func TestMap_All_EarlyExit(t *testing.T) {
	m := concurrentsortedmap.New[int, string]()
	for i := 0; i < 100; i++ {
		m.Put(i, "data")
	}

	count := 0
	for _, _ = range m.All() {
		count++
		if count == 10 {
			break // Exit early
		}
	}

	if count != 10 {
		t.Errorf("Expected to have iterated 10 times, got %d", count)
	}
	// If this test passes without hanging or leaking,
	// your 'defer stop()' logic is solid.
}
