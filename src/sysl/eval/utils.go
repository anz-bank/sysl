package eval

import (
	"fmt"
	"reflect"

	sysl "github.com/anz-bank/sysl/src/proto_old"
)

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

func (e *exprStack) Peek() *exprData {
	return e.s[len(e.s)-1]
}

func getExprText(expr *sysl.Expr) string {
	switch e := expr.Expr.(type) {
	case *sysl.Expr_Name:
		return fmt.Sprintf("Name -> %s", e.Name)
	case *sysl.Expr_Call_:
		return fmt.Sprintf("Call -> %s()", e.Call.Func)
	case *sysl.Expr_GetAttr_:
		return fmt.Sprintf("GetAttr -> %s", e.GetAttr.Attr)
	case *sysl.Expr_Binexpr:
		return fmt.Sprintf("BinaryOp -> %s", e.Binexpr.Op.String())
	case *sysl.Expr_Unexpr:
		return fmt.Sprintf("UnaryOp -> %s", e.Unexpr.Op.String())
	case *sysl.Expr_Literal:
		return fmt.Sprintf("Literal -> %s", e.Literal.String())
	default:
		return fmt.Sprintf("Eval -> %s", reflect.TypeOf(expr.Expr).String())
	}
}
