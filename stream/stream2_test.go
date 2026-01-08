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

func TestStream2_Reduce(t *testing.T) {
	t.Run("Find Pair with Max Key", func(t *testing.T) {
		// Stream of ID -> Score
		s := stream.FromMap(map[int]string{
			10: "low",
			50: "high",
			20: "mid",
		})

		maxK, maxV, err := s.Reduce(func(k1 int, v1 string, k2 int, v2 string) (int, string) {
			if k1 > k2 {
				return k1, v1
			}
			return k2, v2
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if maxK != 50 || maxV != "high" {
			t.Errorf("expected 50:high, got %d:%s", maxK, maxV)
		}
	})

	t.Run("Empty Stream2 Error", func(t *testing.T) {
		s := stream.FromMap(map[string]int{}) // Empty
		_, _, err := s.Reduce(func(k1 string, v1 int, k2 string, v2 int) (string, int) {
			return k1, v1 + v2
		})

		if err == nil || err.Error() != "cannot reduce empty stream" {
			t.Errorf("expected empty stream error, got %v", err)
		}
	})

	t.Run("Error Propagation", func(t *testing.T) {
		var errSource = fmt.Errorf("shard timeout")
		// Simulate a stream that yields one item then fails
		s := stream.New2(func(yield func(int, int) bool) {
			if !yield(1, 1) {
				return
			}
			// In a real scenario, the source would check errSource or
			// the next yield would fail.
		}, &errSource)

		_, _, err := s.Reduce(func(k1, v1, k2, v2 int) (int, int) {
			return k1 + k2, v1 + v2
		})

		if err != errSource {
			t.Errorf("expected %v, got %v", errSource, err)
		}
	})
}
