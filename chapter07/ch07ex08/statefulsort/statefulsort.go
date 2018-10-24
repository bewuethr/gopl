// Package statefulsort provides a sort interface to sort elements with a
// mutable order in which element fields are taken into account for sorting.
package statefulsort

import "fmt"

type sortStruct struct {
	name    string
	cmpFunc func(x, y interface{}) int
}

// StatefulSort implements sort.Interface and keeps track of the order in which
// the fields are looked at when comparing two elements.
type StatefulSort struct {
	elements    []interface{}
	sortStructs []*sortStruct
	reverse     bool
}

// NewStatefulSort returns a stateful sort for the slice of elements e, which
// must implement the Sortable interface. names is a slice of names
// corresponding to the fields of each element, and cmpFuncs is a slice of
// functions to compare two elements. The number of functions must be the same
// as the number of names, and each function compares based on the field
// corresponding to the name with the same index.
//
// A comparison function returns -1 if x sorts before y, 0 if they sort the
// same, and 1 if y sorts before x, just like C's strcmp.
func NewStatefulSort(e []interface{}, names []string, cmpFuncs []func(x, y interface{}) int) StatefulSort {
	if len(names) != len(cmpFuncs) {
		panic("different number of names and sort functions")
	}
	sortStructs := []*sortStruct{}
	for i := range names {
		sortStructs = append(sortStructs, &sortStruct{
			names[i],
			cmpFuncs[i],
		})
	}

	return StatefulSort{e, sortStructs, false}
}

// Elements returns the elements of a stateful sort.
func (s StatefulSort) Elements() []interface{} {
	return s.elements
}

// Len returns the number of elements in s.
func (s StatefulSort) Len() int { return len(s.elements) }

// Swap swaps the elements at indices i and j.
func (s StatefulSort) Swap(i, j int) { s.elements[i], s.elements[j] = s.elements[j], s.elements[i] }

// Less compares two eleements based on the current ordering of the comparison
// functions.
func (s StatefulSort) Less(i, j int) bool {
	for _, sStr := range s.sortStructs {
		switch sStr.cmpFunc(s.elements[i], s.elements[j]) {
		case -1:
			return !s.reverse
		case 1:
			return s.reverse
		}
	}
	return s.reverse
}

// SetPrimary moves the sort function corresponding to name n to the front of
// the sort function slice. If it already is at the front, the sort order is
// reversed by flipping reverse.
func (s *StatefulSort) SetPrimary(n string) error {
	idx := -1
	for i, v := range s.sortStructs {
		if v.name == n {
			idx = i
			break
		}
	}
	if idx == -1 {
		// Name not found
		return fmt.Errorf("sortStruct with name %v not found", n)
	}

	if idx == 0 {
		// Reverse sort order
		s.reverse = !s.reverse
		return nil
	}

	// Rearrange to move new primary to front
	newSortStructs := []*sortStruct{s.sortStructs[idx]}
	newSortStructs = append(newSortStructs, s.sortStructs[:idx]...)
	newSortStructs = append(newSortStructs, s.sortStructs[idx+1:]...)
	s.sortStructs = newSortStructs
	return nil
}
