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

func TestMapStructValues(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	m := New[string, Person]()
	alice := Person{Name: "Alice", Age: 30}
	bob := Person{Name: "Bob", Age: 25}

	m.Set("alice", alice)
	m.Set("bob", bob)

	// basic gets
	v, ok := m.Get("alice")
	if !ok || v.Name != "Alice" || v.Age != 30 {
		t.Errorf("expected alice to be %+v, got %+v (ok=%v)", alice, v, ok)
	}

	// update struct
	alice.Age = 31
	m.Set("alice", alice)
	v, _ = m.Get("alice")
	if v.Age != 31 {
		t.Errorf("expected alice age to be 31 after update, got %d", v.Age)
	}

	// order via All()
	expectOrder := []string{"alice", "bob"}
	i := 0
	for k, _ := range m.All() {
		if k != expectOrder[i] {
			t.Errorf("All() order broken: expected %s at index %d, got %s", expectOrder[i], i, k)
		}
		i++
	}
	if i != len(expectOrder) {
		t.Errorf("All() yielded %d items, expected %d", i, len(expectOrder))
	}
}

func TestCompactTriggeredAndOrderPreserved(t *testing.T) {
	m := New[string, int]()

	// insert 10 items
	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("k%d", i)
		m.Set(k, i)
	}

	// delete 6 items to trigger compact (deletedCount*2 > len(slots))
	for i := 0; i < 6; i++ {
		m.Delete(fmt.Sprintf("k%d", i))
	}

	// remaining should be k6..k9 in that order
	expected := []string{"k6", "k7", "k8", "k9"}
	idx := 0
	for k := range m.Keys() {
		if idx >= len(expected) {
			t.Errorf("unexpected extra key: %s", k)
			break
		}
		if k != expected[idx] {
			t.Errorf("Order broken after compact: expected %s at index %d, got %s", expected[idx], idx, k)
		}
		idx++
	}
	if idx != len(expected) {
		t.Errorf("expected %d keys after deletes, got %d", len(expected), idx)
	}

	// verify values still accessible
	for i, k := range expected {
		v, ok := m.Get(k)
		if !ok || v != i+6 {
			t.Errorf("expected %s -> %d, got %d (ok=%v)", k, i+6, v, ok)
		}
	}
}

func TestDeleteThenReinsertAppearsAtEnd(t *testing.T) {
	m := New[string, int]()
	// insert four items
	keys := []string{"one", "two", "three", "four"}
	for i, k := range keys {
		m.Set(k, i)
	}

	// remove the middle element "two"
	m.Delete("two")

	// re-insert "two" with a new value
	m.Set("two", 20)

	// expected order: one, three, four, two
	expected := []string{"one", "three", "four", "two"}
	i := 0
	for k := range m.Keys() {
		if i >= len(expected) {
			t.Errorf("unexpected extra key: %s", k)
			break
		}
		if k != expected[i] {
			t.Errorf("Order broken after reinsert: expected %s at index %d, got %s", expected[i], i, k)
		}
		i++
	}
	if i != len(expected) {
		t.Errorf("expected %d keys after reinsert, got %d", len(expected), i)
	}

	// verify the reinserted value is the new one
	if v, ok := m.Get("two"); !ok || v != 20 {
		t.Errorf("expected two=20 after reinsert, got %v (ok=%v)", v, ok)
	}
}

func TestNewWithCapacityBasic(t *testing.T) {
	m := NewWithCapacity[string, int](16)
	m.Set("one", 1)
	m.Set("two", 2)

	if v, ok := m.Get("one"); !ok || v != 1 {
		t.Errorf("expected one=1, got %v (ok=%v)", v, ok)
	}
	if v, ok := m.Get("two"); !ok || v != 2 {
		t.Errorf("expected two=2, got %v (ok=%v)", v, ok)
	}

	// delete and ensure other remains
	m.Delete("one")
	if _, ok := m.Get("one"); ok {
		t.Errorf("expected one to be deleted")
	}
	if v, ok := m.Get("two"); !ok || v != 2 {
		t.Errorf("expected two=2 after delete, got %v (ok=%v)", v, ok)
	}
}

func TestMapConvenience(t *testing.T) {
	m := New[string, int]()

	// 1. Test Len and Has on empty map
	if m.Len() != 0 {
		t.Errorf("Expected Len 0, got %d", m.Len())
	}
	if m.Has("apple") {
		t.Error("Empty map should not have 'apple'")
	}

	// 2. Test Len and Has after additions
	m.Set("apple", 1)
	m.Set("banana", 2)
	if m.Len() != 2 {
		t.Errorf("Expected Len 2, got %d", m.Len())
	}
	if !m.Has("apple") {
		t.Error("Map should have 'apple'")
	}

	// 3. Test Len after Delete (The most important part)
	m.Delete("apple")
	if m.Len() != 1 {
		t.Errorf("Expected Len 1 after delete, got %d", m.Len())
	}
	if m.Has("apple") {
		t.Error("Map should not have 'apple' after delete")
	}

	// 4. Test Clear
	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Expected Len 0 after Clear, got %d", m.Len())
	}
	if len(m.slots) != 0 {
		t.Errorf("Slots should be length 0 after Clear, got %d", len(m.slots))
	}

	// Verify we can still add things after Clear
	m.Set("cherry", 3)
	if m.Len() != 1 || !m.Has("cherry") {
		t.Error("Map should work correctly after being cleared")
	}
}
