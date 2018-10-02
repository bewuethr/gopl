// Package treesort provides insertion sort using an unbalanced binary tree.
package treesort

import "fmt"

type tree struct {
	value       int
	left, right *tree
}

// Sort sorts values in place.
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)

	// Exercise 7.3: let's print the sorted tree
	fmt.Println(root)
}

// appendValues appends the elements of t to values in order and returns the
// resulting slice.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		// Equivalent to return &tree{value: value}.
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

// Exercise 7.3: String method for *tree
func (t *tree) String() string {
	var arr []int
	arr = traverse(t, arr)
	return fmt.Sprintf("%v", arr)
}

func traverse(t *tree, arr []int) []int {
	if t == nil {
		return arr
	}

	arr = traverse(t.left, arr)
	arr = append(arr, t.value)
	arr = traverse(t.right, arr)
	return arr
}
