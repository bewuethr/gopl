package ch12ex10

import (
	"reflect"
	"testing"
)

func TestUnmarshalInt(t *testing.T) {
	var got int
	want := 8
	sexpr := []byte("8")

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if got != want {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalString(t *testing.T) {
	var got string
	want := "abc"
	sexpr := []byte(`"abc"`)

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if got != want {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalArray(t *testing.T) {
	var got [3]int
	want := [...]int{1, 2, 3}
	sexpr := []byte("(1 2 3)")

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if got != want {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalSlice(t *testing.T) {
	var got []string
	want := []string{"abc", "def", "ghi"}
	sexpr := []byte(`("abc" "def" "ghi")`)

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalStruct(t *testing.T) {
	var got struct {
		N int
		S string
	}
	want := struct {
		N int
		S string
	}{N: 6, S: "abc"}
	sexpr := []byte(`((N 6) (S "abc"))`)

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalMap(t *testing.T) {
	var got map[string]int
	want := map[string]int{
		"abc": 1,
		"def": 2,
	}
	sexpr := []byte(`(("abc" 1) ("def" 2))`)

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalFalse(t *testing.T) {
	got := true
	want := false
	sexpr := []byte("nil")

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if got != want {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalTrue(t *testing.T) {
	var got bool
	want := true
	sexpr := []byte("t")

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if got != want {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalFloat(t *testing.T) {
	var got float64
	want := 3.14159
	sexpr := []byte("3.14159")

	if err := Unmarshal(sexpr, &got); err != nil {
		t.Fatalf("got error %v, expected nil", err)
	}
	if got != want {
		t.Errorf("Unmarshal(%s) = %v, want %v", sexpr, got, want)
	}
}

func TestUnmarshalInterface(t *testing.T) {
	var tests = []struct {
		input []byte
		want  interface{}
	}{
		{[]byte(`(("int" 5))`), 5},
		{[]byte(`(("float64" 1.1))`), 1.1},
		{[]byte(`(("string" "abc"))`), "abc"},
		{[]byte(`(("[]int" (1 2 3)))`), []int{1, 2, 3}},
		{[]byte(`(("bool" t))`), true},
	}

	for _, test := range tests {
		var got interface{}
		if err := Unmarshal(test.input, &got); err != nil {
			t.Errorf("got error %v, expected nil", err)
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Unmarshal(%s) = %v, want %v", test.input, got, test.want)
		}
	}
}
