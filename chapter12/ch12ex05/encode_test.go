package ch12ex05

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	one := 1
	var myIface interface{} = []int{1, 2, 3}
	var tests = []interface{}{
		"abc",                              // string
		1,                                  // int
		1.1,                                // float
		true,                               // bool
		false,                              // bool
		&one,                               // pointer
		myIface,                            // interface
		[]int{4, 5, 6},                     // int slice
		[...]int{7, 8, 9},                  // array
		[]string{"a", "b", "c"},            // string slice
		struct{ s string }{s: "abc"},       // struct
		map[string]int{"abc": 1, "def": 2}, // map with string keys
		map[int]int{1: 1, 2: 2},            // map with int keys
		[]struct {
			b bool
			n int
			s string
		}{
			{b: true, n: 1, s: "a"},
			{b: false, n: 2, s: "b"},
		}, // slice of structs
	}

	for _, before := range tests {
		// Encode with my encoder
		b, err := Marshal(before)
		if err != nil {
			t.Errorf("encoding %+v: got %q error, wanted nil", before, err)
			continue
		}

		// Decode with stlib JSON decoder
		var after interface{}
		err = json.Unmarshal(b, &after)
		if err != nil {
			t.Errorf("test %s, decoding: got %q error, wanted nil", b, err)
			continue
		}

		// Poor man's equality test: look yourself
		fmt.Printf("%+v --> %+v\n", before, after)
	}
}

func TestMarshal_maps(t *testing.T) {
	var tests = []struct {
		toEncode interface{}
		want     []byte
		errors   bool
	}{
		{
			toEncode: map[string]int{
				"abc": 1,
				"def": 2,
			},
			want:   []byte(`{"abc":1,"def":2}`),
			errors: false,
		},
		{
			toEncode: map[int]int{
				1: 1,
				2: 2,
			},
			want:   []byte(`{"1":1,"2":2}`),
			errors: false,
		},
		{
			toEncode: map[float64]int{
				1.1: 1,
				2.2: 2,
			},
			want:   []byte{},
			errors: true,
		},
	}

	for i, test := range tests {
		got, err := Marshal(test.toEncode)
		if !test.errors && err != nil {
			t.Errorf("test %d: got %q error, wanted nil", i, err)
			continue
		}
		if test.errors && err == nil {
			t.Errorf("test %d: got nil error, wanted non-nil", i)
			continue
		}
		if string(got) != string(test.want) {
			t.Errorf("Marshal(%+v) = %s, want %s", test.toEncode, got, test.want)
		}
	}
}
