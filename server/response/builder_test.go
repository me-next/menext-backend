package builder_test

import (
	"fmt"
	"testing"
)

func getArr() interface{} {
	ret := []string{"a", "b"}

	return ret
}

func TestSimple(t *testing.T) {
	fmt.Println(getArr())
}
