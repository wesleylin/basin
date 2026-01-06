package orderedmap_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	orderedmap "github.com/wesleylin/basin/sequencedmap"
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

func TestJSON_IntKeys(t *testing.T) {
	m := orderedmap.New[int, string]()
	input := `{"10": "ten", "20": "twenty"}`

	if err := json.Unmarshal([]byte(input), m); err != nil {
		t.Fatalf("Failed to unmarshal int keys: %v", err)
	}

	val, _ := m.Get(10)
	if val != "ten" {
		t.Errorf("Expected 'ten', got %v", val)
	}
}

func TestJSON_Array(t *testing.T) {
	t.Run("Array of ints as values", func(t *testing.T) {

		m := orderedmap.New[string, []int]()
		input := `{"numbers": [1, 2, 3]}`

		if err := json.Unmarshal([]byte(input), m); err != nil {
			t.Fatalf("Failed to unmarshal array: %v", err)
		}

		val, _ := m.Get("numbers")
		if len(val) != 3 || val[0] != 1 || val[1] != 2 || val[2] != 3 {
			t.Errorf("Expected [1,2,3], got %v", val)
		}
	})

	t.Run("Array of maps as values", func(t *testing.T) {

		m := orderedmap.New[string, []map[string]int]()
		input := `{"items": [{"a":1}, {"b":2}]}`

		if err := json.Unmarshal([]byte(input), m); err != nil {
			t.Fatalf("Failed to unmarshal array of maps: %v", err)
		}

		val, _ := m.Get("items")
		if len(val) != 2 || val[0]["a"] != 1 || val[1]["b"] != 2 {
			t.Errorf("Expected [{\"a\":1}, {\"b\":2}], got %v", val)
		}
	})

	t.Run("Array as one value", func(t *testing.T) {

		m := orderedmap.New[string, any]()
		input := `{"title": "test", "items": [{"a":1}, {"b":2}]}`

		if err := json.Unmarshal([]byte(input), m); err != nil {
			t.Fatalf("Failed to unmarshal array of maps: %v", err)
		}

		titleRaw, _ := m.Get("title")
		title, ok := titleRaw.(string)
		if !ok || title != "test" {
			t.Errorf("Expected title 'test', got %v", titleRaw)
		}

		val, _ := m.Get("items")
		items, ok := val.([]interface{})
		if !ok {
			t.Fatalf("Expected items to be []interface{}, got %T", val)
		}

		firstItem, ok := items[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected first item to be map[string]interface{}, got %T", items[0])
		}

		valA, ok := firstItem["a"].(float64) // JSON numbers are float64
		if !ok || valA != 1.0 {
			t.Fatalf("Expected first item 'a' to be float64, got %T", firstItem["a"])
		}
	})
}

func TestJSON_LargeNestedFile(t *testing.T) {
	filename := "complex_data.json"

	// 1. Setup: Generate the file
	err := GenerateLargeJSON(filename, 100)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer os.Remove(filename) // Clean up after test

	// 2. Load the file
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	// 3. Unmarshal into OrderedMap
	// Using map[string]any for the values to handle the complex nested nature
	m := orderedmap.New[string, map[string]any]()
	err = json.Unmarshal(content, m)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	// 4. Verify properties
	if m.Len() != 100 {
		t.Errorf("expected 100 users, got %d", m.Len())
	}

	// Check a specific nested field deep in the map
	user, ok := m.Get("user_0050")
	if !ok {
		t.Fatal("could not find user_0050")
	}

	profile := user["profile"].(map[string]any)
	socials := profile["socials"].(map[string]any)

	if socials["github"] != "https://github.com/basin" {
		t.Errorf("expected basin github, got %v", socials["github"])
	}

	// 5. Verify Stream Integrity
	// Ensure the first key is user_0000 and last is user_0099
	// (Standard maps would shuffle these)
	count := 0
	var firstKey, lastKey string
	for k := range m.Keys() {
		if count == 0 {
			firstKey = k
		}
		lastKey = k
		count++
	}

	if firstKey != "user_0000" || lastKey != "user_0099" {
		t.Errorf("Order lost during JSON unmarshal! First: %s, Last: %s", firstKey, lastKey)
	}
}

// GenerateLargeJSON("complex_data.json", 100)

func GenerateLargeJSON(filename string, count int) error {
	data := make(map[string]interface{})
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("user_%04d", i)
		data[key] = map[string]interface{}{
			"id":       i,
			"active":   i%2 == 0,
			"username": fmt.Sprintf("dev_node_%d", i),
			"profile": map[string]interface{}{
				"bio": "Software Engineer at Basin Corp",
				"socials": map[string]interface{}{
					"github":  "https://github.com/basin",
					"twitter": "@basin_io",
				},
				"stats": map[string]int{
					"commits": 100 * i,
					"stars":   5 * i,
				},
			},
			"tags": []string{"golang", "ordered-map", "stream", "performance"},
			"config": map[string]string{
				"theme": "dark",
				"font":  "JetBrains Mono",
			},
			// Adding more fields to hit the 30+ requirement...
			"meta_1": "v1", "meta_2": "v2", "meta_3": "v3",
			"meta_4": "v4", "meta_5": "v5", "meta_6": "v6",
			"meta_7": "v7", "meta_8": "v8", "meta_9": "v9",
			"meta_10": "v10",
		}
	}

	bytes, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile(filename, bytes, 0644)
}
