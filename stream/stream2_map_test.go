package stream

import (
	"errors"
	"fmt"
	"iter"
	"reflect"
	"testing"
)

// Helper to turn a map into a Stream2 for testing
func streamFromMap[K comparable, V any](m map[K]V) Stream2[K, V] {
	var err error
	return Stream2[K, V]{
		err: &err,
		seq: func(yield func(K, V) bool) {
			for k, v := range m {
				if !yield(k, v) {
					return
				}
			}
		},
	}
}

func TestMap2(t *testing.T) {
	input := map[string]int{"one": 1, "two": 2}
	s := streamFromMap(input)

	// Transform: string->int (len), int->string (val)
	mapped := Map2(s, func(k string, v int) (int, string) {
		return len(k), fmt.Sprintf("val-%d", v)
	})

	results := make(map[int]string)
	for k, v := range mapped.seq {
		results[k] = v
	}

	// expected := map[int]string{3: "val-1", 3: "val-2"} // Both "one" and "two" have length 3
	// Note: in a real map "one" would overwrite "two" here, but the stream yields both.
	if len(results) != 1 || results[3] == "" {
		t.Errorf("Map2 failed, got %v", results)
	}
}

func TestMap2Err(t *testing.T) {
	t.Run("Success Path", func(t *testing.T) {
		var errPtr error
		s := Stream2[string, int]{
			err: &errPtr,
			seq: func(yield func(string, int) bool) {
				yield("a", 1)
			},
		}

		mapped := Map2Err(s, func(k string, v int) (string, int, error) {
			return k + "!", v * 10, nil
		})

		for k, v := range mapped.seq {
			if k != "a!" || v != 10 {
				t.Errorf("Unexpected values: %s, %d", k, v)
			}
		}
		if *mapped.err != nil {
			t.Errorf("Expected nil error, got %v", *mapped.err)
		}
	})

	t.Run("Error Path", func(t *testing.T) {
		var errPtr error
		s := Stream2[string, int]{
			err: &errPtr,
			seq: func(yield func(string, int) bool) {
				// Check the boolean! If Map2Err says "stop" (false), we must stop.
				if !yield("a", 1) {
					return
				}
				yield("b", 2)
			},
		}

		sentinelErr := errors.New("boom")
		mapped := Map2Err(s, func(k string, v int) (string, int, error) {
			if k == "a" {
				return "", 0, sentinelErr
			}
			return k, v, nil
		})

		// Trigger execution
		for range mapped.seq {
		}

		if *mapped.err != sentinelErr {
			t.Errorf("Expected error %v, got %v", sentinelErr, *mapped.err)
		}
	})
}

func TestFlatMap2(t *testing.T) {
	var errPtr error
	s := Stream2[string, int]{
		err: &errPtr,
		seq: func(yield func(string, int) bool) {
			yield("numbers", 2)
		},
	}

	// FlatMap transforms 1 entry into 2 entries
	flat := FlatMap2(s, func(k string, v int) iter.Seq2[string, int] {
		return func(yield func(string, int) bool) {
			for i := 1; i <= v; i++ {
				if !yield(fmt.Sprintf("%s-%d", k, i), i) {
					return
				}
			}
		}
	})

	var results []string
	for k, _ := range flat.seq {
		results = append(results, k)
	}

	expected := []string{"numbers-1", "numbers-2"}
	if !reflect.DeepEqual(results, expected) {
		t.Errorf("Expected %v, got %v", expected, results)
	}
}
