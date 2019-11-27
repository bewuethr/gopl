package ch11ex04

import (
	"math/rand"
	"strings"
	"testing"
	"time"
	"unicode"
)

// randomNoisyPalindrome returns a string that, after removal of non-letters and
// normalization is a palindrome.
func randomNoisyPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25) // random length up to 24
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000)) // random rune up to '\u0999'
		runes[i] = r
		r1 := r
		if unicode.IsLetter(r) {
			if unicode.IsLower(r) {
				r = unicode.ToUpper(r)
				if unicode.ToLower(r) != r1 {
					r = r1
				}
			} else if unicode.IsUpper(r) {
				r = unicode.ToLower(r)
				if unicode.ToUpper(r) != r1 {
					r = r1
				}
			}
		}
		runes[n-1-i] = r
	}

	// Iterate over runes, and with a chance of 25%, insert a random space or
	// interpunction character.
	const noise = ".,-=()[]{}!?:;\"'/+|& \t"
	var b strings.Builder
	for _, r := range runes {
		b.WriteRune(r)
		if rng.Intn(4) == 0 {
			b.WriteByte(noise[rng.Intn(len(noise))])
		}
	}

	return b.String()
}

func TestRandomNoisyPalindromes(t *testing.T) {
	// Initialize a pseudo-random number generator.
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomNoisyPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}
