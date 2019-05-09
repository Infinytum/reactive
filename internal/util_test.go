package internal

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	a := []interface{}{"1", "2", "3", "4"}
	b := []interface{}{"4", "3", "2", "1"}

	c := [][]interface{}{a, b}
	d := [][]interface{}{b, a}
	Reverse(c)

	if !reflect.DeepEqual(c, d) {
		t.Error("c does not match reversed array")
	}
}
