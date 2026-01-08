package stream

import (
	"fmt"
	"testing"
)

func TestMath(t *testing.T) {
	t.Run("Sum integers", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3, 4, 5})
		val, err := Sum(s)
		if err != nil || val != 15 {
			t.Errorf("expected 15, got %v (err: %v)", val, err)
		}
	})

	t.Run("Max/Min with floats", func(t *testing.T) {
		s := FromSlice([]float64{10.5, 2.1, 55.0, 0.5})

		maxVal, _ := Max(s)
		if maxVal != 55.0 {
			t.Errorf("expected max 55.0, got %v", maxVal)
		}

		// Re-create stream for Min (since Max consumed the first one)
		s = FromSlice([]float64{10.5, 2.1, 55.0, 0.5})
		minVal, _ := Min(s)
		if minVal != 0.5 {
			t.Errorf("expected min 0.5, got %v", minVal)
		}
	})

	t.Run("Average success", func(t *testing.T) {
		s := FromSlice([]int{10, 20, 30})
		avg, err := Average(s)
		if err != nil || avg != 20.0 {
			t.Errorf("expected 20.0, got %v", avg)
		}
	})

	t.Run("Math Error Propagation", func(t *testing.T) {
		var errSource = fmt.Errorf("read failure")
		// Stream that fails mid-way
		s := New(func(yield func(int) bool) {
			yield(10)
			// error happens here
		}, &errSource)

		_, err := Sum(s)
		if err != errSource {
			t.Errorf("expected error %v, got %v", errSource, err)
		}
	})

	t.Run("Empty Stream Average", func(t *testing.T) {
		s := FromSlice([]int{})
		_, err := Average(s)
		if err == nil {
			t.Error("expected error for empty stream average, got nil")
		}
	})
}
