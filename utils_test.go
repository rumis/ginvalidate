package ginvalidate

import (
	"fmt"
	"testing"
)

func TestKeyFormat(t *testing.T) {

	k1 := "x[]"
	k2 := "y[0]"
	k3 := "y[1]"

	fmt.Println(FormatKey(k1))
	fmt.Println(FormatKey(k2))
	fmt.Println(FormatKey(k3))

}
