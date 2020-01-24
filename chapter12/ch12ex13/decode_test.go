package ch12ex13

import (
	"reflect"
	"testing"
)

func TestUnmarshalTags(t *testing.T) {
	type MyType struct {
		LongKey    string `sexpr:"lk"`
		LongKeyInt int    `sexpr:"lki"`
	}
	in := []byte(`((lk "abc") (lki 123))`)
	want := MyType{LongKey: "abc", LongKeyInt: 123}

	got := MyType{}
	if err := Unmarshal(in, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal(%v) = %+v, want %+v", in, got, want)
	}
}

func TestUnmarshalNoTags(t *testing.T) {
	type MyType struct {
		StrKey string
		IntKey int
	}
	in := []byte(`((StrKey "abc") (IntKey 123))`)
	want := MyType{StrKey: "abc", IntKey: 123}

	got := MyType{}
	if err := Unmarshal(in, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal(%v) = %+v, want %+v", in, got, want)
	}
}

func TestUnmarshalSliceValue(t *testing.T) {
	type MyType struct {
		SliceKey []string `sexpr:"sk"`
	}
	in := []byte(`((sk ("abc" "def" "ghi")))`)
	want := MyType{SliceKey: []string{"abc", "def", "ghi"}}

	got := MyType{}
	if err := Unmarshal(in, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal(%v) = %+v, want %+v", in, got, want)
	}
}
