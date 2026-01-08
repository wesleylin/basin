package stream

import (
	"fmt"
	"testing"
)

func TestTerminalExtensions(t *testing.T) {
	t.Run("Last: Successful", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3, 4, 5})
		val, err := s.Last()
		if err != nil || val != 5 {
			t.Errorf("expected 5, got %v", val)
		}
	})

	t.Run("Last: Empty Stream", func(t *testing.T) {
		s := FromSlice([]int{})
		_, err := s.Last()
		if err == nil {
			t.Error("expected error for empty stream Last()")
		}
	})

	t.Run("ToMap: Successful", func(t *testing.T) {
		entries := []Entry[string, int]{
			{"apple", 1},
			{"banana", 2},
		}
		s := FromSlice(entries)
		res, err := ToMap(s)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res["apple"] != 1 || res["banana"] != 2 {
			t.Errorf("map content mismatch: %v", res)
		}
	})

	t.Run("GroupBy: Categorization", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		people := []Person{
			{"Alice", 25},
			{"Bob", 30},
			{"Charlie", 25},
		}
		s := FromSlice(people)

		// Group by Age
		grouped, err := GroupBy(s, func(p Person) int { return p.Age })
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(grouped[25]) != 2 || len(grouped[30]) != 1 {
			t.Errorf("grouping failed: %v", grouped)
		}
	})

	t.Run("Terminal: Error Handling", func(t *testing.T) {
		var errSource = fmt.Errorf("disk error")
		s := New(func(yield func(int) bool) {
			yield(1)
			// error happens here
		}, &errSource)

		_, err := s.Last()
		if err != errSource {
			t.Errorf("expected %v, got %v", errSource, err)
		}
	})
}
