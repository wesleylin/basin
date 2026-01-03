package stream_test

import (
	"fmt"
	"testing"

	"github.com/wesleylin/basin/stream"
)

func TestStream2(t *testing.T) {
	t.Run("FromMap and Collect", func(t *testing.T) {
		input := map[string]int{"a": 1, "b": 2}
		s := stream.FromMap(input)

		results, err := s.Collect()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Since standard maps are unordered, we check length and content
		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}

		// Verify contents exist
		foundA := false
		for _, p := range results {
			if p.Key == "a" && p.Value == 1 {
				foundA = true
			}
		}
		if !foundA {
			t.Error("key 'a' not found in collected results")
		}
	})

	t.Run("MapValues transformation", func(t *testing.T) {
		input := map[string]int{"one": 1, "two": 2}
		s := stream.FromMap(input)

		// Transform int -> string values
		s2 := stream.MapValues(s, func(v int) string {
			return fmt.Sprintf("num-%d", v)
		})

		results, err := s2.Collect()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, p := range results {
			expected := fmt.Sprintf("num-%d", input[p.Key])
			if p.Value != expected {
				t.Errorf("for key %s: expected %s, got %s", p.Key, expected, p.Value)
			}
		}
	})

	t.Run("Keys and Values bridges", func(t *testing.T) {
		input := map[string]int{"a": 1}
		s := stream.FromMap(input)

		// Test Keys() bridge
		keys, err := s.Keys().Collect()
		if err != nil || len(keys) != 1 || keys[0] != "a" {
			t.Errorf("Keys() bridge failed: got %v, %v", keys, err)
		}

		// Test Values() bridge
		values, err := s.Values().Collect()
		if err != nil || len(values) != 1 || values[0] != 1 {
			t.Errorf("Values() bridge failed: got %v, %v", values, err)
		}
	})
}
