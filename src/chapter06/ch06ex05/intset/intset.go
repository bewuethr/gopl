// Package intset provides a set of integers based on a bit vector, using a
// platform-specific word size.
package intset

import (
	"bytes"
	"fmt"
)

const wordSize = 32 << (^uint(0) >> 63)

// An IntSet is a set of small non-negative integers. Its zero value represents
// the empty set.
type IntSet struct {
	words []uint
}

// Len returns the number of elements in s.
func (s *IntSet) Len() int {
	var count int
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for i := 0; i < wordSize; i++ {
			if word&(1<<uint(i)) != 0 {
				count++
			}
		}
	}
	return count
}

// Remove removes x from s, if s contains x.
func (s *IntSet) Remove(x int) {
	if s.Has(x) {
		word, bit := x/wordSize, uint(x%wordSize)
		s.words[word] &^= 1 << bit
	}
}

// Clear removes all elements from s.
func (s *IntSet) Clear() {
	for i := range s.words {
		s.words[i] = 0
	}
}

// Copy returns a copy of s.
func (s *IntSet) Copy() *IntSet {
	var t IntSet
	for _, word := range s.words {
		t.words = append(t.words, word)
	}
	return &t
}

// AddAll adds a list of values to s.
func (s *IntSet) AddAll(vals ...int) {
	for _, v := range vals {
		s.Add(v)
	}
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/wordSize, uint(x%wordSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/wordSize, uint(x%wordSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// IntersectWith sets s to the intersection of s and t.
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if !t.Has(wordSize*i + j) {
					s.Remove(wordSize*i + j)
				}
			}
		}
	}
}

// DifferenceWith sets s to the set difference s - t.
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if t.Has(wordSize*i + j) {
					s.Remove(wordSize*i + j)
				}
			}
		}
	}
}

// SymmetricDifference sets s to the symmetric difference of s and t.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	s2 := s.Copy()
	s2.IntersectWith(t)
	t2 := t.Copy()

	s.DifferenceWith(t)
	t2.DifferenceWith(s2)
	s.UnionWith(t2)
}

// Elems returns an int slice contaning the elements of s.
func (s *IntSet) Elems() []int {
	var elems []int
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				elems = append(elems, i*wordSize+j)
			}
		}
	}
	return elems
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < wordSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", wordSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}
