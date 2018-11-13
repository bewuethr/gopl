// ch07ex13 tests an implementation of the Expr interface with an additional
// pretty printer method Print().
package main

import (
	"fmt"

	"gopl/chapter07/ch07ex13/eval"
)

func main() {
	exprStr := "pow(2,sin(y))*pow(+2,sin(x))/-12"
	env := eval.Env{"x": 25, "y": 20}
	expr, err := eval.Parse(exprStr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Input:\t%s\n", exprStr)
	fmt.Printf("Env:\t%+v\n", env)
	fmt.Printf("Pretty:\t%s\n", expr.Print())
	fmt.Printf("Result:\t%.6g\n\n", expr.Eval(env))

	fmt.Println("Parsing expr.Print())...")
	expr, err = eval.Parse(expr.Print())
	fmt.Printf("Pretty:\t%s\n", expr.Print())
	fmt.Printf("Result:\t%.6g\n\n", expr.Eval(env))
}
