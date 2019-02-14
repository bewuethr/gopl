// Ch05ex14 uses breadthFirst to traverse other structures.
package main

import "fmt"

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
}

func main() {
	breadthFirst(getDeps, []string{
		"calculus",
		"networks",
		"programming languages",
		"algorithms",
		"compilers",
		"databases",
	})
}

// breadthFirst calls f for each item in the worklist.  Any items returned by f
// are added to the worklist.  f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

// getDeps returns the prerequisites of item.
func getDeps(item string) []string {
	fmt.Println(item)
	return prereqs[item]
}
