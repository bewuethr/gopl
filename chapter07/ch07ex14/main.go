// Ch07ex14 implements a new type min that implements the Expr interface.
package main

import (
	"fmt"

	"github.com/bewuethr/gopl/chapter07/ch07ex14/eval"
)

type testCase struct {
	env     eval.Env
	exprStr string
}

func main() {
	for _, tc := range []testCase{
		{
			exprStr: "min(,)",
		},
		{
			exprStr: "min()",
		},
		{
			exprStr: "min(1, 2, 3)",
		},
		{
			env:     eval.Env{"a": 3, "b": 2, "c": 1.1},
			exprStr: "min(a, b, c)",
		},
		{
			env:     eval.Env{"a": 3, "b": 2, "c": 1.1},
			exprStr: "min(min(a, b, c), 1.0)",
		},

		{
			env: eval.Env{"x": 1.1, "r": 0.85, "y": -0.25},
			exprStr: "min( sin(-x) * pow(1.5, -r), " +
				"pow(2, sin(y)) * pow(2, sin(x)) / 12, " +
				"sin(x*y/10) / 10 )",
		},
	} {
		expr, err := eval.Parse(tc.exprStr)
		if err != nil {
			fmt.Printf("\nparsing error for %s: %s\n", tc.exprStr, err)
			continue
		}
		vars := make(map[eval.Var]bool)
		if err = expr.Check(vars); err != nil {
			fmt.Printf("\nexpression %s fails check: %s\n", tc.exprStr, err)
			continue
		}
		fmt.Printf("\nInput:\t%s\n", tc.exprStr)
		fmt.Printf("Env:\t%+v\n", tc.env)
		fmt.Printf("Result:\t%.6g\n\n", expr.Eval(tc.env))
	}
}
