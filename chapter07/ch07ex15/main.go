// Ch07ex15 reads one expression from standard input, prompts for variable
// values and evaluates the expression in the resulting environment.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/bewuethr/gopl/chapter07/ch07ex13/eval"
)

func main() {
	fmt.Println("Enter expression:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "problem scanning input")
		os.Exit(1)
	}

	exprStr := scanner.Text()

	expr, err := eval.Parse(exprStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing error: %v\n", err)
		os.Exit(1)
	}

	vars := make(map[eval.Var]bool)

	if err = expr.Check(vars); err != nil {
		fmt.Fprintf(os.Stderr, "semantic error: %v\n", err)
		os.Exit(1)
	}

	env, err := promptForValues(vars)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting values: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v = %v\n", expr.Print(), expr.Eval(env))
}

func promptForValues(vars map[eval.Var]bool) (eval.Env, error) {
	env := eval.Env{}
	scanner := bufio.NewScanner(os.Stdin)
	for v := range vars {
		fmt.Printf("Value for %v: ", v)
		scanner.Scan()
		strVal := scanner.Text()
		val, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return nil, errors.New("can't convert input to float")
		}
		env[v] = val
	}

	return env, nil
}
