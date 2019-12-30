// Package ch12ex04 provides a pretty-printing encoder for S-expressions.
package ch12ex04

import (
	"bytes"
	"fmt"
	"reflect"
)

// MarshalIndent pretty-encodes a Go value in S-expression form.
func MarshalIndent(v interface{}) ([]byte, error) {
	p := printer{margin: margin, space: margin}
	if err := encode(&p, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return p.Bytes(), nil
}

const margin = 80

type tokenKind int

const (
	strToken tokenKind = iota
	openDelim
	closeDelim
	blankToken
)

type token struct {
	kind   tokenKind
	str    string
	length int
}

type printer struct {
	bytes.Buffer

	tokens []*token

	// used for printing
	pStack []int
	margin int
	space  int

	// used for scanning
	sStack    []*token
	rightotal int
}

func (p *printer) print(t *token) {
	switch t.kind {
	case strToken:
		p.WriteString(t.str)
		p.space -= len(t.str)
	case openDelim:
		p.pStack = append(p.pStack, p.space)
	case closeDelim:
		p.pStack = p.pStack[:len(p.pStack)-1]
	case blankToken:
		if t.length > p.space {
			p.space = p.pStack[len(p.pStack)-1] - 2
			p.WriteString(fmt.Sprintf("\n%*s", p.margin-p.space, ""))
		} else {
			p.WriteByte(' ')
			p.space--
		}
	}
}

func (p *printer) scanStrToken(s string) {
	t := &token{
		kind:   strToken,
		str:    s,
		length: len(s),
	}
	if len(p.sStack) == 0 {
		p.print(t)
	} else {
		p.tokens = append(p.tokens, t)
		p.rightotal += t.length
	}
}

func (p *printer) scanOpenDelim() {
	t := &token{kind: openDelim}
	if len(p.sStack) == 0 {
		p.rightotal = 1
	}
	t.length = -p.rightotal
	p.tokens = append(p.tokens, t)
	p.sStack = append(p.sStack, t)
	p.scanStrToken("(")
}

func (p *printer) scanCloseDelim() {
	t := &token{kind: closeDelim}
	p.scanStrToken(")")
	p.tokens = append(p.tokens, t)
	x := p.sStack[len(p.sStack)-1]
	p.sStack = p.sStack[:len(p.sStack)-1]
	x.length += p.rightotal
	if x.kind == blankToken {
		x = p.sStack[len(p.sStack)-1]
		p.sStack = p.sStack[:len(p.sStack)-1]
		x.length += p.rightotal
	}
	if len(p.sStack) == 0 {
		for _, token := range p.tokens {
			p.print(token)
		}
		p.tokens = []*token{}
	}
}

func (p *printer) scanBlankToken() {
	t := &token{kind: blankToken}
	last := len(p.sStack) - 1
	if p.sStack[last].kind == blankToken {
		x := p.sStack[last]
		p.sStack = p.sStack[:last]
		x.length += p.rightotal
	}
	t.length = -p.rightotal
	p.tokens = append(p.tokens, t)
	p.sStack = append(p.sStack, t)
	p.rightotal++
}

// encode writes an S-expression representation of v to pretty printer p.
func encode(p *printer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		p.scanStrToken("nil")

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		p.scanStrToken(fmt.Sprintf("%d", v.Int()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.scanStrToken(fmt.Sprintf("%d", v.Uint()))

	case reflect.String:
		p.scanStrToken(fmt.Sprintf("%q", v.String()))

	case reflect.Ptr:
		return encode(p, v.Elem())

	case reflect.Array, reflect.Slice: // (value ...)
		p.scanOpenDelim()
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				p.scanBlankToken()
			}
			if err := encode(p, v.Index(i)); err != nil {
				return err
			}
		}
		p.scanCloseDelim()

	case reflect.Struct: // ((name value) ...)
		p.scanOpenDelim()
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				p.scanBlankToken()
			}
			p.scanOpenDelim()
			p.scanStrToken(v.Type().Field(i).Name)
			p.scanBlankToken()
			if err := encode(p, v.Field(i)); err != nil {
				return err
			}
			p.scanCloseDelim()
		}
		p.scanCloseDelim()

	case reflect.Map: // ((key value) ...)
		p.scanOpenDelim()
		for i, key := range v.MapKeys() {
			if i > 0 {
				p.scanBlankToken()
			}
			p.scanOpenDelim()
			if err := encode(p, key); err != nil {
				return err
			}
			p.scanBlankToken()
			if err := encode(p, v.MapIndex(key)); err != nil {
				return err
			}
			p.scanCloseDelim()
		}
		p.scanCloseDelim()

	default: // float, complex, bool, chan, func, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
