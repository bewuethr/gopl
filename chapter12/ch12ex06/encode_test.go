package ch12ex06

import (
	"testing"
)

type type1 struct {
	s1, s2 string
}

type type2 struct {
	s1, s2   string
	n1, n2   int
	b1, b2   bool
	c1, c2   complex128
	f1, f2   float64
	t1, t2   type1
	sl1, sl2 []int
}

func TestMarshal(t *testing.T) {
	var tests = []struct {
		v    interface{}
		want string
	}{
		{
			type1{s1: "abc", s2: ""},
			`((s1 "abc"))`,
		},
		{
			type1{s1: "", s2: ""},
			`()`,
		},
		{
			type2{
				s1:  "abc",
				s2:  "",
				n1:  1,
				n2:  int(0),
				b1:  true,
				b2:  false,
				c1:  1.1 + 1.1i,
				c2:  0 + 0i,
				f1:  1.1,
				f2:  0.0,
				t1:  type1{s1: "abc", s2: "def"},
				t2:  type1{},
				sl1: []int{1, 2, 3},
				sl2: nil,
			},
			`((s1 "abc") (n1 1) (b1 t) (c1 #C(1.1 1.1)) (f1 1.1) (t1 ((s1 "abc") (s2 "def"))) (sl1 (1 2 3)))`,
		},
	}

	for _, test := range tests {
		got, err := Marshal(test.v)
		if err != nil {
			t.Errorf("got %v, want nil error", err)
			continue
		}
		if string(got) != test.want {
			t.Errorf("Marshal(%+v) = %s, want %s", test.v, got, test.want)
		}
	}
}
