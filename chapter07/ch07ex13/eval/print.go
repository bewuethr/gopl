package eval

import (
	"bytes"
	"fmt"
)

// Print prints variable v.
func (v Var) Print() string {
	return string(v)
}

func (l literal) Print() string {
	return fmt.Sprintf("%.6g", l)
}

func (u unary) Print() string {
	return fmt.Sprintf("%c%s", u.op, u.x.Print())
}

func (b binary) Print() string {
	return fmt.Sprintf("%s %c %s", b.x.Print(), b.op, b.y.Print())
}

func (c call) Print() string {
	if _, ok := numParams[c.fn]; !ok {
		panic(fmt.Sprintf("unknown function %q", c.fn))
	}
	bb := bytes.NewBufferString(fmt.Sprintf("%s( ", c.fn))
	for i, arg := range c.args {
		if i > 0 {
			bb.WriteString(", ")
		}
		bb.WriteString(arg.Print())
	}
	bb.WriteString(" )")
	return bb.String()
}
