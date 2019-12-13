package intset

import (
	"testing"
)

func TestIntSet(t *testing.T) {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)

	refX := RefSet{1: true, 144: true, 9: true}

	if x.String() != refX.String() {
		t.Errorf("Add(1,144,9) = %v, expected %v", &x, refX)
	}

	y.Add(9)
	y.Add(42)

	refY := RefSet{9: true, 42: true}
	if y.String() != refY.String() {
		t.Errorf("Add(9,42) = %v, expected %v", &y, refY)
	}

	x.UnionWith(&y)

	refX.UnionWith(refY)

	if x.String() != refX.String() {
		t.Errorf("x.UnionWith(&y) = %v, expected %v", &x, refX)
	}

	for _, n := range []int{9, 23} {
		if x.Has(n) != refX[n] {
			t.Errorf("x.(%d) = %v, expected %v", n, x.Has(n), refX[n])
		}
	}
}
