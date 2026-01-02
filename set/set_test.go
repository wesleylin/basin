package set

import (
	"fmt"
	"testing"
)

func TestSetBasic(t *testing.T) {
	s := New[string]()
	fmt.Println("Created new set:", s)
}
