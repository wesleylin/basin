package stream

import (
	"fmt"
	"iter"
	"slices"
	"strconv"
	"testing"
)

func TestMapFunctions(t *testing.T) {
	t.Run("Map transforms types correctly", func(t *testing.T) {
		input := []int{1, 2, 3}
		s := FromSlice(input)

		// Transform int -> string
		mappedStream := Map(s, func(n int) string {
			return fmt.Sprintf("val:%d", n)
		})

		results, err := mappedStream.Collect()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []string{"val:1", "val:2", "val:3"}
		if !slices.Equal(results, expected) {
			t.Errorf("expected %v, got %v", expected, results)
		}
	})

	t.Run("MapErr trips the live wire on failure", func(t *testing.T) {
		input := []string{"1", "2", "not_a_number", "4"}
		s := FromSlice(input)

		// Attempt to parse strings to ints
		// Note: Using a wrapper because your MapErr currently requires func(T) (T, error)
		parsedStream := MapErr(s, func(val string) (string, error) {
			if _, err := strconv.Atoi(val); err != nil {
				return "", err // This trips the wire
			}
			return val, nil
		})

		results, err := parsedStream.Collect()

		// 1. We should get an error
		if err == nil {
			t.Fatal("expected error from MapErr, got nil")
		}

		// 2. We are are using custom collect, we return nil/empty data on error
		if len(results) != 0 {
			t.Errorf("expected 0 items on error (idiomatic Go), got %d", len(results))
		}
	})

	t.Run("FlatMap flattens and respects short-circuiting", func(t *testing.T) {
		input := []int{1, 2}
		s := FromSlice(input)

		// Each N becomes [N, N]
		fm := FlatMap(s, func(n int) iter.Seq[int] {
			return func(yield func(int) bool) {
				if !yield(n) {
					return
				}
				if !yield(n) {
					return
				}
			}
		})

		// Use Take(3) to ensure it stops mid-flattening
		results, err := fm.Take(3).Collect()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []int{1, 1, 2} // The second '2' is truncated by Take(3)
		if !slices.Equal(results, expected) {
			t.Errorf("expected %v, got %v", expected, results)
		}
	})
}
