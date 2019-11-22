package intset

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"testing"
)

type refSet map[int]bool

func (s refSet) String() string {
	var v []int
	for x := range s {
		v = append(v, x)
	}
	sort.Ints(v)

	var buf bytes.Buffer
	buf.WriteByte('{')
	for _, x := range v {
		if buf.Len() > len("{") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, strconv.Itoa(x))
	}
	buf.WriteByte('}')
	return buf.String()
}

func (s refSet) unionWith(t refSet) {
	for k := range t {
		s[k] = true
	}
}

func TestIntSet(t *testing.T) {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)

	refX := refSet{1: true, 144: true, 9: true}

	if x.String() != refX.String() {
		t.Errorf("Add(1,144,9) = %v, expected %v", &x, refX)
	}

	y.Add(9)
	y.Add(42)

	refY := refSet{9: true, 42: true}
	if y.String() != refY.String() {
		t.Errorf("Add(9,42) = %v, expected %v", &y, refY)
	}

	x.UnionWith(&y)

	refX.unionWith(refY)

	if x.String() != refX.String() {
		t.Errorf("x.UnionWith(&y) = %v, expected %v", &x, refX)
	}

	for _, n := range []int{9, 23} {
		if x.Has(n) != refX[n] {
			t.Errorf("x.(%d) = %v, expected %v", n, x.Has(n), refX[n])
		}
	}
}
