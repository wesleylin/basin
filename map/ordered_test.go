package ordered

import (
	"fmt"
	"testing"
)

func TestMapBasic(t *testing.T) {
	m := New[string, int]()
	fmt.Println("Created new ordered map:", m)
}
