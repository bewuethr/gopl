package ch13ex02

import "testing"

type link struct {
	value string
	tail  *link
}

func TestIsCyclic(t *testing.T) {
	a, b := &link{value: "a"}, &link{value: "b"}
	a.tail = a
	var aIntf interface{} = a
	var tests = []struct {
		a    interface{}
		want bool
	}{
		{1, false},
		{a, true},
		{[]int{1, 2, 3}, false},
		{[]link{*a}, true},
		{struct{ l *link }{l: a}, true},
		{*a, true},
		{aIntf, true},
		{map[int]link{1: *a}, true},
		{map[int]link{1: *b}, false},
	}

	for _, test := range tests {
		if got := IsCyclic(test.a); got != test.want {
			t.Errorf("IsCyclic(%v) == %v, want %v", test.a, got, test.want)
		}
	}
}
