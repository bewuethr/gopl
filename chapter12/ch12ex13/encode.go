// Package ch12ex13 provides a means for converting Go objects to and from
// S-expressions.
package ch12ex13

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// Marshal encodes a Go value in S-expression form. It respects "sexpr" struct
// tags for renaming, omitting null values ("omitempty" option) and generally
// skipping a field (tag "-"), analogous to JSON struct tags.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encode writes to buf an S-expression representation of v.
func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return encode(buf, v.Elem())

	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')

	case reflect.Struct: // ((name value) ...)
		buf.WriteByte('(')

		isFirst := true
		for i := 0; i < v.NumField(); i++ {
			fieldInfo := v.Type().Field(i)
			name, option, hasComma := getOption(fieldInfo.Tag.Get("sexpr"))
			if name == "" {
				name = strings.ToLower(fieldInfo.Name)
			}
			if name == "-" && option == "" && !hasComma {
				continue
			}
			if option == "omitempty" && v.Field(i).IsZero() {
				continue
			}

			if !isFirst {
				buf.WriteByte(' ')
			}
			isFirst = false

			fmt.Fprintf(buf, "(%s ", name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Map: // ((key value) ...)
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	default: // float, complex, bool, chan, func, interface
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

// getOption extracts the validation option from the tag string if present. The
// returned boolean indicates if there was any comma at all, even if there were
// no options.
func getOption(tag string) (string, string, bool) {
	if i := strings.Index(tag, ","); i != -1 {
		fields := strings.Split(tag, ",")
		return fields[0], fields[1], true
	}
	return tag, "", false
}
