package stream

import (
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
