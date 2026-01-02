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

func TestMapAliases(t *testing.T) {
	m := New[string, int]()
	m.Set("x", 100)

	// Test All()
	for k, v := range m.All() {
		if k != "x" || v != 100 {
			t.Fail()
		}
	}
}

func TestMapKeys(t *testing.T) {
	m := New[string, int]()
	m.Set("x", 100)
	m.Set("y", 200)

	// Test Keys()
	for k := range m.Keys() {
		if k != "x" && k != "y" {
			fmt.Println("Unexpected key:", k)
			t.Fail()
		}
	}
}

func TestDeleteOrder(t *testing.T) {
	m := New[string, int]()
	m.Set("apple", 1)
	m.Set("banana", 2)
	m.Set("cherry", 3)

	m.Delete("banana") // Remove the middle item

	// banana should be gone
	if _, exists := m.Get("banana"); exists {
		t.Error("banana should have been deleted")
	}

	// apple and cherry should still be there, in that order
	expected := []string{"apple", "cherry"}
	i := 0
	for k := range m.Keys() {
		if k != expected[i] {
			t.Errorf("Order broken! Expected %s at index %d, got %s", expected[i], i, k)
		}
		i++
	}
}

func TestDeleteNonExistent(t *testing.T) {
	m := New[string, int]()
	m.Set("apple", 1)

	// Delete non-existent key
	m.Delete("banana") // Should not panic or error

	// apple should still be there
	val, exists := m.Get("apple")
	if !exists || val != 1 {
		t.Errorf("Expected apple to be 1, got %v", val)
	}
}
