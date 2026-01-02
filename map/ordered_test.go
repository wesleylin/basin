package ordered

import (
	"fmt"
	"testing"
)

func TestMapBasic(t *testing.T) {
	m := New[string, int]()
	fmt.Println("Created new ordered map:", m)

	// 1. setting value
	m.Set("apple", 1)
	m.Set("banana", 2)

	// 2. getting value
	val, exists := m.Get("apple")
	if !exists || val != 1 {
		t.Errorf("Expected apple to be 1, got %v", val)
	}

	// 3. updating value
	m.Set("apple", 10)
	val, _ = m.Get("apple")
	if val != 10 {
		t.Errorf("Expected updated apple to be 10, got %v", val)
	}
}

func TestMapOrder(t *testing.T) {
	m := New[string, int]()

	// add items in order
	keys := []string{"first", "second", "third"}
	for i, k := range keys {
		m.Set(k, i)
	}

	// check same order
	for i, k := range keys {
		val, _ := m.Get(k)
		if val != i {
			t.Errorf("Order broken! Expected %s to have value %d", k, i)
		}
	}
}
