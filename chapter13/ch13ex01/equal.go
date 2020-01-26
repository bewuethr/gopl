// Package ch13ex01 provides a deep equivalence relation that considers numbers
// of any type equal if they differ by less than one part in a billion.
package ch13ex01

import (
	"fmt"
	"math/cmplx"
	"reflect"
	"unsafe"
)

type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}

const epsilon = 1 / 1e9

// Equal reports whether x and y are deeply equal, considering numbers of any
// type to be equal if they differ by less than one billionth.
//
// Map keys are always compared with ==, not deeply.  (This matters for keys
// containing pointers or interfaces.)
func Equal(x, y interface{}) bool {
	seen := make(map[comparison]bool)
	return equal(reflect.ValueOf(x), reflect.ValueOf(y), seen)
}

func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}
	if !compatible(x.Type(), y.Type()) {
		fmt.Println("not compatible")
		return false
	}

	// cycle check
	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true // identical references
		}
		c := comparison{xptr, yptr, x.Type()}
		if seen[c] {
			return true // already seen
		}
		seen[c] = true
	}

	switch x.Kind() {
	case reflect.Bool:
		return x.Bool() == y.Bool()

	case reflect.String:
		return x.String() == y.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Float32, reflect.Float64, reflect.Uint64, reflect.Uintptr,
		reflect.Complex64, reflect.Complex128:
		return numEqual(x, y)

	case reflect.Chan, reflect.UnsafePointer, reflect.Func:
		return x.Pointer() == y.Pointer()

	case reflect.Ptr, reflect.Interface:
		return equal(x.Elem(), y.Elem(), seen)

	case reflect.Array, reflect.Slice:
		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i), seen) {
				return false
			}
		}
		return true

	case reflect.Struct:
		for i, n := 0, x.NumField(); i < n; i++ {
			if !equal(x.Field(i), y.Field(i), seen) {
				return false
			}
		}
		return true

	case reflect.Map:
		if x.Len() != y.Len() {
			return false
		}
		for _, k := range x.MapKeys() {
			if !equal(x.MapIndex(k), y.MapIndex(k), seen) {
				return false
			}
		}
		return true
	}
	panic("unreachable")
}

func compatible(x, y reflect.Type) bool {
	fmt.Printf("x: %v, y: %v\n", x.Kind().String(), y.Kind().String())
	if isNumeric(x.Kind()) && isNumeric(y.Kind()) {
		return true
	}
	return x.Kind() == y.Kind()
}

func isNumeric(k reflect.Kind) bool {
	return k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 ||
		k == reflect.Int32 || k == reflect.Int64 || k == reflect.Uint ||
		k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 ||
		k == reflect.Uint64 || k == reflect.Uintptr || k == reflect.Float32 ||
		k == reflect.Float64 || k == reflect.Complex64 || k == reflect.Complex128
}

func numEqual(x, y reflect.Value) bool {
	if isComplex(x) || isComplex(y) {
		return complexEqual(x, y)
	}

	floatX, floatY := getFloat(x), getFloat(y)
	if floatX == 0.0 {
		return abs(floatY) < epsilon
	}
	return abs((floatX-floatY)/floatX) < epsilon
}

func isComplex(x reflect.Value) bool {
	k := x.Kind()
	return k == reflect.Complex64 || k == reflect.Complex128
}

func complexEqual(x, y reflect.Value) bool {
	cX, cY := getComplex(x), getComplex(y)
	if cmplx.Abs(cX) == 0.0 {
		return cmplx.Abs(cY) < epsilon
	}
	return cmplx.Abs(cX-cY)/cmplx.Abs(cX) < epsilon
}

func getFloat(x reflect.Value) float64 {
	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(x.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		return float64(x.Uint())
	case reflect.Float32, reflect.Float64:
		return x.Float()
	default:
		panic(fmt.Sprintf("unexpected type %v in getFloat", x.Kind().String()))
	}
}

func getComplex(x reflect.Value) complex128 {
	if isComplex(x) {
		return x.Complex()
	}
	return complex(getFloat(x), 0.0)
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
