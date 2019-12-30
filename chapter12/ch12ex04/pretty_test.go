package ch12ex04

import (
	"fmt"
	"testing"
)

type Movie struct {
	Title, Subtitle  string
	Year             int
	Color, LongTitle bool
	Actor            map[string]string
	Oscars           []string
	Sequel           *string
	BudgetMillions   float64
	ComplexNum       complex128
	InterfacePtr     *interface{}
}

func TestMarshalIndent(t *testing.T) {
	var iface interface{} = []int{1, 2, 3}
	strangelove := Movie{
		Title:     "Dr. Strangelove",
		Subtitle:  "How I Learned to Stop Worrying and Love the Bomb",
		Year:      1964,
		Color:     false,
		LongTitle: true,
		Actor: map[string]string{
			"Dr. Strangelove":            "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
		BudgetMillions: 1.8,
		ComplexNum:     1 + 2i,
		InterfacePtr:   &iface,
	}

	mySexpr, err := MarshalIndent(strangelove)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	fmt.Println(string(mySexpr))
}
