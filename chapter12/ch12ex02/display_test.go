package ch12ex02

import "testing"

func TestDisplay(t *testing.T) {
	// a pointer that points to itself
	type P *P
	var p P
	p = &p
	Display("p", p)
	// Output:
	// Display p (ch12ex02.P), max depth 3:
	// (*(*(*p))) = ch12ex02.P 0xc00000e028

	// a map that contains itself
	type M map[string]M
	m := make(M)
	m[""] = m
	Display("m", m)
	// Output:
	// Display m (ch12ex02.M), max depth 3:
	// m[""][""][""] = ch12ex02.M 0xc000078270

	// a slice that contains itself
	type S []S
	s := make(S, 1)
	s[0] = s
	Display("s", s)
	// Output:
	// Display s (ch12ex02.S), max depth 3:
	// s[0][0][0] = ch12ex02.S 0xc00000c120

	// a linked list that eats its own tail
	type Cycle struct {
		Value int
		Tail  *Cycle
	}
	var c Cycle
	c = Cycle{42, &c}
	Display("c", c)
	// Output:
	// Display c (ch12ex02.Cycle), max depth 3:
	// c.Value = 42
	// (*c.Tail).Value = 42
	// (*c.Tail).Tail = *ch12ex02.Cycle 0xc0000105f0
}
