// Ch07ex16 implements an expression evaluator as a web service. The request
// expects two query string parameters, "expr" containing the expression, and
// "env" containing a JSON object with the environment, for example
// ?expr=a%2Bb&env={"a"%3A1%2C"b"%3A2}
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bewuethr/gopl/chapter07/ch07ex13/eval"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	// Get expression and environment from query string parameters
	values := req.URL.Query()
	exprStr := values.Get("expr")
	if exprStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `missing "expr" query string parameter`)
		return
	}
	envStr := values.Get("env")
	if envStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `missing "env" query string parameter`)
		return
	}

	// Decode environment
	var env eval.Env
	err := json.Unmarshal([]byte(envStr), &env)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't decode environment JSON blob %v: %v\n", envStr, err)
		return
	}

	// Parse expression
	expr, err := eval.Parse(exprStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "can't parse expression %v: %v\n", exprStr, err)
		return
	}

	// Check expression
	vars := make(map[eval.Var]bool)
	if err = expr.Check(vars); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "semantic error: %v\n", err)
		return
	}

	// Evaluate and return
	fmt.Fprintf(w, "%v = %v\n", expr.Print(), expr.Eval(env))
}
