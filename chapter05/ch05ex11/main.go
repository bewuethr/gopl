// Ch05ex11 extends topoSort to report cycles.
package main

import (
	"fmt"
	"os"
	"sort"
)

// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
	"linear algebra":        {"calculus"}, // circular dependency
}

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) []string {
	var order []string
	// if a key exists with value false, it counts as temporarily marked (used
	// to detect loops); after a node has been appended to order, it is marked
	// with true.
	seen := make(map[string]bool)
	var visitAll func(items []string)

	visitAll = func(items []string) {
		for _, item := range items {
			s, ok := seen[item]
			if ok && !s {
				// We've found a loop
				fmt.Fprintf(os.Stderr, "Found a loop involving %v\n", item)
				os.Exit(1)
			}
			if !s {
				seen[item] = false // temporary mark
				visitAll(m[item])
				order = append(order, item)
				seen[item] = true // definitive mark
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return order
}
