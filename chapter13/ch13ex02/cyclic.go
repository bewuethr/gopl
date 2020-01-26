// Package ch13ex02 provides a function that checks if its argument is cylic.
package ch13ex02

import (
	"reflect"
	"unsafe"
)

type typePtr struct {
	ptr unsafe.Pointer
	t   reflect.Type
}

// IsCyclic reports if its argument is cylic.
func IsCyclic(x interface{}) bool {
	seen := make(map[typePtr]bool)
	return isCyclic(reflect.ValueOf(x), seen)
}

func isCyclic(x reflect.Value, seen map[typePtr]bool) bool {
	if !x.IsValid() {
		return false
	}

	// cycle check
	if x.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		tp := typePtr{xptr, x.Type()}
		if seen[tp] {
			return true // already seen
		}
		seen[tp] = true
	}

	switch x.Kind() {
	case reflect.Ptr, reflect.Interface:
		return isCyclic(x.Elem(), seen)

	case reflect.Array, reflect.Slice:
		for i := 0; i < x.Len(); i++ {
			if isCyclic(x.Index(i), seen) {
				return true
			}
		}
		return false

	case reflect.Struct:
		for i, n := 0, x.NumField(); i < n; i++ {
			if isCyclic(x.Field(i), seen) {
				return true
			}
		}
		return false

	case reflect.Map:
		for _, k := range x.MapKeys() {
			if isCyclic(x.MapIndex(k), seen) {
				return true
			}
		}
		return false

	default: // all other types
		return false
	}
}
