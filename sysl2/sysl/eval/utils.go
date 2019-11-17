package eval

import sysl "github.com/anz-bank/sysl/src/proto"

type exprData struct {
	e    *sysl.Expr
	args Scope
}

type exprStack struct {
	s []*exprData
}

func (e *exprStack) Push(scope Scope, expr *sysl.Expr) {
	e.s = append(e.s, &exprData{expr, scope})
}

func (e *exprStack) Pop() *exprData {
	top := e.s[len(e.s)-1]
	e.s = e.s[:len(e.s)-1]
	return top
}
