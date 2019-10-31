package eval

import (
	"fmt"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
)

type unaryFunc func(*sysl.Value) *sysl.Value

//nolint:gochecknoglobals
var unaryFunctions = map[sysl.Expr_UnExpr_Op]unaryFunc{
	sysl.Expr_UnExpr_NEG:    unaryNeg,
	sysl.Expr_UnExpr_STRING: unaryString,
}

func evalUnaryFunc(op sysl.Expr_UnExpr_Op, arg *sysl.Value) *sysl.Value {
	if x, has := unaryFunctions[op]; has {
		return x(arg)
	}
	panic(errors.Errorf("evalUnaryFunc: Operation %v not supported\n", op))
}

func unaryNeg(arg *sysl.Value) *sysl.Value {
	switch x := arg.Value.(type) {
	case *sysl.Value_I:
		return MakeValueI64(-x.I)
	case *sysl.Value_B:
		return MakeValueBool(!x.B)
	}
	panic(errors.Errorf("unaryNeg for %v not supported", arg.Value))
}

func unaryString(arg *sysl.Value) *sysl.Value {
	switch x := arg.Value.(type) {
	case *sysl.Value_S:
		return arg
	case *sysl.Value_I:
		return MakeValueString(fmt.Sprintf("%d", x.I))
	case *sysl.Value_B:
		return MakeValueString(map[bool]string{true: "true", false: "false"}[x.B])
	}
	return MakeValueString(arg.String())
}
