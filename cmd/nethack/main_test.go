package main

import "testing"

func TestMake2DArrayOfChar(t *testing.T) {
	a := make2DArrayOfChar(5, 6)

	if 5 != len(a) {
		t.Errorf("wrong length of x dim: %d", len(a))
	}

	if 6 != len(a[0]) {
		t.Errorf("wrong length of y dim: %d", len(a[0]))
	}
}
