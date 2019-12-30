package ch12ex04

import (
	"fmt"
	"testing"
)

type Movie struct {
	Title, Subtitle string
	Year            int
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func TestMarshalIndent(t *testing.T) {
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
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
	}

	mySexpr, err := MarshalIndent(strangelove)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	fmt.Println(string(mySexpr))
}
