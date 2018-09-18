// ch05ex15 implements variadic max and min functions.
package main

import "fmt"

var (
	errMaxNoValue = fmt.Errorf("max: no values provided")
	errMinNoValue = fmt.Errorf("min: no values provided")
)

func main() {
	cases1 := []struct {
		in        []int
		maxExp    int
		maxErrExp error
		minExp    int
		minErrExp error
	}{
		{[]int{1, 2, 3}, 3, nil, 1, nil},
		{[]int{3, 2, 1}, 3, nil, 1, nil},
		{[]int{2, 2}, 2, nil, 2, nil},
		{[]int{1}, 1, nil, 1, nil},
		{[]int{}, 0, errMaxNoValue, 0, errMinNoValue},
	}

	for _, c := range cases1 {
		max, err1 := max(c.in...)
		min, err2 := min(c.in...)
		if max != c.maxExp || min != c.minExp || err1 != c.maxErrExp || err2 != c.minErrExp {
			fmt.Printf("Expected\tActual\n%v\t%v\n%v\t%v\n%v\t%v\n%v\t%v\n",
				c.maxExp, max, c.minExp, min, err1, c.maxErrExp, err2, c.minErrExp)
		}
	}

	cases2 := []struct {
		in     []int
		maxExp int
		minExp int
	}{
		{[]int{1, 2, 3}, 3, 1},
		{[]int{3, 2, 1}, 3, 1},
		{[]int{2, 2}, 2, 2},
		{[]int{1}, 1, 1},
	}

	for _, c := range cases2 {
		max := max2(c.in[0], c.in[1:]...)
		min := min2(c.in[0], c.in[1:]...)
		if max != c.maxExp || min != c.minExp {
			fmt.Printf("Expected\tActual\n%v\t%v\n%v\t%v\n",
				c.maxExp, max, c.minExp, min)
		}
	}
}

// max returns the maximum argument; if there are no arguments, the error is
// non-nil.
func max(vals ...int) (int, error) {
	if len(vals) == 0 {
		return 0, errMaxNoValue
	}
	m := vals[0]
	for _, v := range vals[1:] {
		if v > m {
			m = v
		}
	}
	return m, nil
}

// min returns the minimum argument; if there are no arguments, the error is
// non-nil.
func min(vals ...int) (int, error) {
	if len(vals) == 0 {
		return 0, errMinNoValue
	}
	m := vals[0]
	for _, v := range vals[1:] {
		if v < m {
			m = v
		}
	}
	return m, nil
}

// max2 return the maximum argument and makes sure at least one argument is
// provided.
func max2(first int, vals ...int) int {
	m := first
	for _, v := range vals {
		if v > m {
			m = v
		}
	}
	return m
}

// min2 return the minimum argument and makes sure at least one argument is
// provided.
func min2(first int, vals ...int) int {
	m := first
	for _, v := range vals {
		if v < m {
			m = v
		}
	}
	return m
}
