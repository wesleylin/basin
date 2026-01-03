package orderedmap_test

import (
	"encoding/json"
	"testing"

	"github.com/wesleylin/basin/orderedmap"
)

func TestJSON_Marshal(t *testing.T) {
	m := orderedmap.New[string, int]()
	m.Put("z", 1)
	m.Put("a", 2)
	m.Put("m", 3)

	// Test 1: Simple Order
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expected := `{"z":1,"a":2,"m":3}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}

	// Test 2: Order after Delete
	m.Delete("a")
	data, err = json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	expectedAfterDelete := `{"z":1,"m":3}`
	if string(data) != expectedAfterDelete {
		t.Errorf("after delete: expected %s, got %s", expectedAfterDelete, string(data))
	}
}

func TestJSON_Unmarshal(t *testing.T) {
	input := `{"apple":100,"banana":200,"cherry":300}`
	m := orderedmap.New[string, int]()

	err := json.Unmarshal([]byte(input), m)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify count
	if m.Len() != 3 {
		t.Errorf("expected length 3, got %d", m.Len())
	}

	// Verify order is preserved during unmarshal
	// Using your existing Stream/Collect logic to verify order
	var keys []string
	for k := range m.Keys() {
		keys = append(keys, k)
	}

	expectedKeys := []string{"apple", "banana", "cherry"}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("at index %d: expected key %s, got %s", i, expectedKeys[i], k)
		}
	}
}

func TestJSON_ComplexValues(t *testing.T) {
	type Stats struct {
		Age    int  `json:"age"`
		Active bool `json:"active"`
	}

	m := orderedmap.New[string, Stats]()
	m.Put("user1", Stats{Age: 30, Active: true})
	m.Put("user2", Stats{Age: 25, Active: false})

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Round trip
	m2 := orderedmap.New[string, Stats]()
	if err := json.Unmarshal(data, m2); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	v, _ := m2.Get("user1")
	if v.Age != 30 {
		t.Errorf("expected age 30, got %d", v.Age)
	}
}
