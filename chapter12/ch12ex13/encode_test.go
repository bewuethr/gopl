package ch12ex13

import (
	"bytes"
	"testing"
)

type MyStruct1 struct {
	LongKey    string `sexpr:"lk"`
	LongKeyInt int    `sexpr:"lki"`
}

type MyStruct2 struct {
	NoEncode string `sexpr:"-"`
	DashName string `sexpr:"-,"`
}

type MyStruct3 struct {
	OmitYes string `sexpr:"oy,omitempty"`
	OmitNo  string `sexpr:"on"`
}

func TestMarshal(t *testing.T) {

	var tests = []struct {
		in   interface{}
		want []byte
	}{
		{
			in:   MyStruct1{LongKey: "foo", LongKeyInt: 123},
			want: []byte(`((lk "foo") (lki 123))`),
		},
		{
			in:   MyStruct2{NoEncode: "foo", DashName: "bar"},
			want: []byte(`((- "bar"))`),
		},
		{
			in:   MyStruct3{OmitYes: "foo", OmitNo: "bar"},
			want: []byte(`((oy "foo") (on "bar"))`),
		},
		{
			in:   MyStruct3{OmitYes: "", OmitNo: "bar"},
			want: []byte(`((on "bar"))`),
		},
	}

	for _, test := range tests {
		got, err := Marshal(test.in)
		if err != nil {
			t.Errorf("got error %v, expected nil", err)
			continue
		}

		if !bytes.Equal(got, test.want) {
			t.Errorf("Marshal(%+v) = %s, want %s", test.in, got, test.want)
		}
	}
}
