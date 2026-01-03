package stream

import (
	"errors"
	"fmt"
	"slices"
	"testing"
)

func TestStream(t *testing.T) {
	t.Run("Basic Filter and Take", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		// Logic: Take even numbers, but only the first 3
		got, err := FromSeq(slices.Values(input)).
			Filter(func(n int) bool { return n%2 == 0 }).
			Take(3).
			Collect()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := []int{2, 4, 6}
		if !slices.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Lazy Evaluation", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		processedCount := 0

		// Track how many items actually pass through the filter
		s := FromSeq(slices.Values(input)).
			Take(2).
			Filter(func(n int) bool {
				processedCount++
				fmt.Println(processedCount)
				return true
			})

		// We only want 2 items, so the filter should NOT run for 3, 4, or 5
		_, _ = s.Collect()

		if processedCount != 2 {
			t.Errorf("expected to process 2 items, but processed %d", processedCount)
		}
	})

	t.Run("Empty Stream", func(t *testing.T) {
		var input []int
		got, err := FromSeq(slices.Values(input)).
			Filter(func(n int) bool { return true }).
			Collect()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if len(got) != 0 {
			t.Errorf("expected empty slice, got %v", got)
		}
	})

	t.Run("Take More Than Available", func(t *testing.T) {
		input := []int{1, 2}
		got, err := FromSeq(slices.Values(input)).
			Take(10).
			Collect()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !slices.Equal(got, input) {
			t.Errorf("got %v, want %v", got, input)
		}
	})
}

func TestCount(t *testing.T) {
	t.Run("counts all items in a simple slice", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}

		// Create a stream from a standard slice
		s := FromSeq(slices.Values(items))

		count, err := s.Count()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 5 {
			t.Errorf("expected count 5, got %d", count)
		}
	})

	t.Run("counts correctly after filtering", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5, 6}

		s := FromSeq(slices.Values(items))

		// Chain a filter for even numbers
		count, err := s.Filter(func(n int) bool { return n%2 == 0 }).Count()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 3 {
			t.Errorf("expected count 3, got %d", count)
		}
	})

	t.Run("returns error when stream is poisoned", func(t *testing.T) {
		sentinelErr := errors.New("database connection lost")
		var streamErr error

		// Create a stream that fails on the 3rd item
		s := New(func(yield func(int) bool) {
			for i := 1; i <= 5; i++ {
				if i == 3 {
					streamErr = sentinelErr // Trip the Live Wire
					return                  // Stop iterating
				}
				if !yield(i) {
					return
				}
			}
		}, &streamErr)

		count, err := s.Count()

		if !errors.Is(err, sentinelErr) {
			t.Errorf("expected error %v, got %v", sentinelErr, err)
		}
		if count != 0 {
			t.Errorf("expected count 0 on error, got %d", count)
		}
	})
}
