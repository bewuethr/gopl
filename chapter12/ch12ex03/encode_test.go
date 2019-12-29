package ch12ex03

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	anInt := 8
	var anInterface interface{} = []int{1, 2, 3}
	aFunc := func() {}

	var tests = []struct {
		v       interface{}
		want    []byte
		wantErr error
	}{
		{8, []byte("8"), nil},                                     // int
		{-8, []byte("-8"), nil},                                   // negative int
		{uint(8), []byte("8"), nil},                               // uint
		{"abc", []byte(`"abc"`), nil},                             // string
		{&anInt, []byte("8"), nil},                                // pointer to int
		{[]int{0, 1, 2}, []byte("(0 1 2)"), nil},                  // slice
		{[...]int{0, 1, 2}, []byte("(0 1 2)"), nil},               // array
		{[]int{}, []byte("()"), nil},                              // empty slice
		{struct{ m, n int }{5, 6}, []byte("((m 5) (n 6))"), nil},  // struct
		{true, []byte("t"), nil},                                  // true
		{false, []byte("nil"), nil},                               // false
		{3.14159, []byte("3.14159"), nil},                         // floating point
		{1.1 + 2.1i, []byte("#C(1.1 2.1)"), nil},                  // complex number
		{&anInterface, []byte(`("[]int" (1 2 3))`), nil},          // interface
		{aFunc, []byte{}, fmt.Errorf("unsupported type: func()")}, // unsupported type

	}

	for _, test := range tests {
		got, gotErr := Marshal(test.v)
		if gotErr != nil && gotErr.Error() != test.wantErr.Error() {
			t.Errorf("Marshal(%v) = (_, %v); want (_, %v)", test.v, gotErr, test.wantErr)
			continue
		}

		if !bytes.Equal(got, test.want) {
			t.Errorf("Marshal(%v) = %s, want %s", test.v, got, test.want)
		}
	}
}
