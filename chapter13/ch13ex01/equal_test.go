package ch13ex01

import "testing"

func TestEqual(t *testing.T) {
	var tests = []struct {
		a, b interface{}
		want bool
	}{
		{1, 1, true},
		{1, 2, false},
		{1, 1.0, true},
		{1, 1 + 0i, true},
		{1e10, 1e10 - 1, true},
		{1e9, 1e9 - 2, false},
		{[]int{1}, []float32{1.0}, true},
		{[]float32{1.0, 2.0, 3.0}, []complex64{1 + 0i, 2 + 0i, 3 + 0i}, true},
		{1e10 + 0i, 1e10 + 1i, true},
		{1e10 + 0i, 1e10 + 10i, false},
		{
			struct {
				i int
				f float64
			}{100, 32.0},
			struct {
				c complex64
				u uintptr
			}{100 + 0i, 32},
			true,
		},
	}

	for _, test := range tests {
		if got := Equal(test.a, test.b); got != test.want {
			t.Errorf("Equal(%v, %v) == %v, want %v", test.a, test.b, got, test.want)
		}
	}
}
