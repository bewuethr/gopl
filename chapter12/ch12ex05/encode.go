// Package ch12ex05 provides a JSON encoder.
package ch12ex05

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// Marshal encodes a Go value as JSON.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encode writes to buf a JSON representation of v.
func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("null")

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

	case reflect.Array, reflect.Slice: // [value,...]
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(']')

	case reflect.Struct: // {"name":value,...}
		buf.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "%q:", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case reflect.Map: // {"key":value,...}; keys must be string or int type
		buf.WriteByte('{')
		switch v.Type().Key().Kind() {
		case reflect.String:
			mapKeys := v.MapKeys()
			strMapKeys := make([]string, len(mapKeys))
			for i, key := range mapKeys {
				strMapKeys[i] = key.String()
			}
			sort.Strings(strMapKeys)
			for i, key := range strMapKeys {
				if i > 0 {
					buf.WriteByte(',')
				}
				fmt.Fprintf(buf, "%q:", key)
				if err := encode(buf, v.MapIndex(reflect.ValueOf(key))); err != nil {
					return err
				}
			}
			buf.WriteByte('}')

		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			mapKeys := v.MapKeys()
			strMapKeys := make([]struct {
				s string
				v reflect.Value
			}, len(mapKeys))
			for i, key := range mapKeys {
				strMapKeys[i].s = strconv.FormatInt(key.Int(), 10)
				strMapKeys[i].v = key
			}
			sort.Slice(strMapKeys, func(i, j int) bool { return strMapKeys[i].s < strMapKeys[j].s })
			for i, key := range strMapKeys {
				if i > 0 {
					buf.WriteByte(',')
				}
				fmt.Fprintf(buf, "%q:", key.s)
				if err := encode(buf, v.MapIndex(key.v)); err != nil {
					return err
				}
			}
			buf.WriteByte('}')

		default:
			return fmt.Errorf("unsupported key type: %s", v.Type().Key().Kind())
		}

	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", v.Float())

	case reflect.Interface: // ("type name" value)
		return encode(buf, v.Elem())

	default: // chan, complex, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
