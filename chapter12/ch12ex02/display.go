// Package ch12ex02 provides a means to display structured data. The display
// function is modified to limit recursion.
package ch12ex02

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// maxDepth is the default maximum recursion depth; it can be modified by
// clients.
var maxDepth = 3

// SetMaxDepth allows clients to set the maximum recursion depth.
func SetMaxDepth(d int) error {
	if d < 1 {
		return errors.New("maximum depth must be >= 1")
	}
	maxDepth = d
	return nil
}

// Display is a recursive value printer.
func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T), max depth %d:\n", name, x, maxDepth)
	display(name, reflect.ValueOf(x), 0)
}

// formatAtom formats a value without inspecting its internal structure.
func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	// ...floating-point and complex cases omitted for brevity...
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr,
		reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return v.Type().String() + " value"
	}
}

func display(path string, v reflect.Value, depth int) {
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if depth < maxDepth {
				display(fmt.Sprintf("%s[%d]", path, i), v.Index(i), depth+1)
			} else {
				fmt.Printf("%s = %s\n", path, formatAtom(v))
			}
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			if depth < maxDepth {
				display(fieldPath, v.Field(i), depth+1)
			} else {
				fmt.Printf("%s = %s\n", fieldPath, formatAtom(v))
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if depth < maxDepth {
				display(fmt.Sprintf("%s[%s]", path,
					formatAtom(key)), v.MapIndex(key), depth+1)
			} else {
				fmt.Printf("%s = %s\n", path, formatAtom(v))
			}
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			if depth < maxDepth {
				display(fmt.Sprintf("(*%s)", path), v.Elem(), depth+1)
			} else {
				fmt.Printf("%s = %s\n", path, formatAtom(v))
			}
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			if depth < maxDepth {
				display(path+".value", v.Elem(), depth+1)
			} else {
				fmt.Printf("%s = %s\n", path, formatAtom(v))
			}
		}
	default: // basic types, channels, funcs
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}
