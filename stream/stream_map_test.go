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

func TestFold(t *testing.T) {
	t.Run("Sum Lengths of Strings", func(t *testing.T) {
		// Stream of strings -> reduced to an int (total characters)
		s := FromSlice([]string{"apple", "pear", "kiwi"})

		// Fold(stream, initialValue, accumulator)
		totalLen, err := Fold(s, 0, func(acc int, s string) int {
			return acc + len(s)
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if totalLen != 13 { // 5 + 4 + 4
			t.Errorf("expected 13, got %d", totalLen)
		}
	})

	t.Run("Build a Map from Stream", func(t *testing.T) {
		// This is a common pattern for "collecting" data into a lookup table
		type User struct {
			ID   int
			Name string
		}
		users := []User{{1, "Alice"}, {2, "Bob"}}
		s := FromSlice(users)

		userMap, err := Fold(s, make(map[int]string), func(acc map[int]string, u User) map[int]string {
			acc[u.ID] = u.Name
			return acc
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if userMap[1] != "Alice" || userMap[2] != "Bob" {
			t.Errorf("map was not built correctly: %v", userMap)
		}
	})

	t.Run("Fold with Upstream Error", func(t *testing.T) {
		var errSource = fmt.Errorf("connection lost")
		s := New(func(yield func(int) bool) {
			if !yield(1) {
				return
			}
			if !yield(2) {
				return
			}
		}, &errSource)

		// The fold should return the initial value and the error
		result, err := Fold(s, 100, func(acc int, v int) int {
			return acc + v
		})

		if err != errSource {
			t.Errorf("expected %v, got %v", errSource, err)
		}
		// In an error scenario, we return the seed value
		if result != 100 {
			t.Errorf("expected seed 100, got %d", result)
		}
	})
}
