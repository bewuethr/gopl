package ch12ex09

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

// A Token is an interface holding one of the token types: Symbol, String, Int,
// StartList, or EndList.
type Token interface{}

type (
	// A Symbol is an unquoted name.
	Symbol string

	// A String is an S-expression string value.
	String string

	// An Int is an S-expression integer value.
	Int int

	// StartList represents the beginning of a list.
	StartList byte

	// EndList represents the end of a list.
	EndList byte
)

func (s StartList) String() string { return string(s) }
func (s EndList) String() string   { return string(s) }

// A Decoder represents an S-epxressions parser reading tokens from an input
// stream.
type Decoder struct {
	s *scanner.Scanner
}

// NewDecoder creates a new Deocoder reading from r.
func NewDecoder(r io.Reader) *Decoder {
	var s scanner.Scanner
	return &Decoder{s: s.Init(r)}
}

// Token returns the next S-expression token in the input stream. At the end of
// the input stream, Token returns nil, io.EOF.
func (d *Decoder) Token() (Token, error) {
	var t Token
	var err error
	switch r := d.s.Scan(); r {
	case scanner.EOF:
		err = io.EOF
	case scanner.Ident:
		t = Symbol(d.s.TokenText())
	case scanner.String:
		t = String(d.s.TokenText())
	case scanner.Int:
		var i int
		i, err = strconv.Atoi(d.s.TokenText())
		if err == nil {
			t = Int(i)
		}
	case '(':
		t = StartList('(')
	case ')':
		t = EndList(')')
	default:
		err = fmt.Errorf("unexpected token: %v", r)
	}
	return t, err
}
