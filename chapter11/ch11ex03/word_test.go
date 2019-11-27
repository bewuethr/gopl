package ch11ex03

import (
	"math/rand"
	"testing"
	"time"
	"unicode"
)

// randomNonPalindrome returns a string between 2 and 24 runes long that is
// guaranteed not to be a palindrome. The runes are guaranteed to be letters.
func randomNonPalindrome(rng *rand.Rand) string {
	n := rng.Intn(23) + 2 // random length between 2 and 24
	runes := make([]rune, n)

	// Set first and last rune to different values
	runes[0] = randomRune(rng)
	r := randomRune(rng)
	for r == runes[0] {
		r = randomRune(rng)
	}
	runes[n-1] = r

	// Fill rest of slice
	for i := 1; i < n-1; i++ {
		runes[i] = randomRune(rng)
	}
	return string(runes)
}

// randomRune returns a random letter rune.
func randomRune(rng *rand.Rand) rune {
	r := rune(rng.Intn(0x1000))
	for !unicode.IsLetter(r) {
		r = rune(rng.Intn(0x1000))
	}
	return r
}

func TestRandomNonPalindromes(t *testing.T) {
	// Initialize a pseudo-random number generator.
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))

	for i := 0; i < 1000; i++ {
		p := randomNonPalindrome(rng)
		if IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = true", p)
		}
	}
}
