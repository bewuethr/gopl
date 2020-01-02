// Package ch12ex07 provides an S-expression encoder with a streaming API.
package ch12ex07

import (
	"fmt"
	"io"
	"reflect"
)

// Encoder writes S-expression values to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the S-expression encoding of v to the stream, followed by a
// newline character.
func (enc *Encoder) Encode(v interface{}) error {
	if err := enc.encode(reflect.ValueOf(v)); err != nil {
		return err
	}
	fmt.Fprint(enc.w, "\n")
	return nil
}

// encode writes an S-expression representation of v to the underlying writer.
func (enc *Encoder) encode(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Fprint(enc.w, "nil")

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(enc.w, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(enc.w, "%d", v.Uint())

	case reflect.String:
		fmt.Fprintf(enc.w, "%q", v.String())

	case reflect.Ptr:
		return enc.encode(v.Elem())

	case reflect.Array, reflect.Slice: // (value ...)
		fmt.Fprint(enc.w, "(")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				fmt.Fprint(enc.w, " ")
			}
			if err := enc.encode(v.Index(i)); err != nil {
				return err
			}
		}
		fmt.Fprint(enc.w, ")")

	case reflect.Struct: // ((name value) ...)
		fmt.Fprint(enc.w, "(")
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				fmt.Fprint(enc.w, " ")
			}
			fmt.Fprintf(enc.w, "(%s ", v.Type().Field(i).Name)
			if err := enc.encode(v.Field(i)); err != nil {
				return err
			}
			fmt.Fprint(enc.w, ")")
		}
		fmt.Fprint(enc.w, ")")

	case reflect.Map: // ((key value) ...)
		fmt.Fprint(enc.w, "(")
		for i, key := range v.MapKeys() {
			if i > 0 {
				fmt.Fprint(enc.w, " ")
			}
			fmt.Fprint(enc.w, "(")
			if err := enc.encode(key); err != nil {
				return err
			}
			fmt.Fprint(enc.w, " ")
			if err := enc.encode(v.MapIndex(key)); err != nil {
				return err
			}
			fmt.Fprint(enc.w, ")")
		}
		fmt.Fprint(enc.w, ")")

	case reflect.Bool:
		if v.Bool() {
			fmt.Fprint(enc.w, "t")
		} else {
			fmt.Fprint(enc.w, "nil")
		}

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(enc.w, "%g", v.Float())

	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(enc.w, "#C(%g %g)", real(v.Complex()), imag(v.Complex()))

	case reflect.Interface: // ("type name" value)
		fmt.Fprint(enc.w, "(")
		fmt.Fprintf(enc.w, `"%v" `, v.Elem().Type())
		if err := enc.encode(v.Elem()); err != nil {
			return err
		}
		fmt.Fprint(enc.w, ")")

	default: // chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
