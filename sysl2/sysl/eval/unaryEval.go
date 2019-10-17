package eval

import (
	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
)

type unaryFunc func(*sysl.Value) *sysl.Value

//nolint:gochecknoglobals
var unaryFunctions = map[sysl.Expr_UnExpr_Op]unaryFunc{
	sysl.Expr_UnExpr_NEG: unaryNeg,
}

func evalUnaryFunc(op sysl.Expr_UnExpr_Op, arg *sysl.Value) *sysl.Value {
	if x, has := unaryFunctions[op]; has {
		return x(arg)
	}
	panic(errors.Errorf("evalUnaryFunc: Operation %v not supported\n", op))
}

func unaryNeg(arg *sysl.Value) *sysl.Value {
	if x, ok := arg.Value.(*sysl.Value_I); ok {
		return MakeValueI64(-x.I)
	}
	panic(errors.Errorf("unaryNeg for %v not supported", arg.Value))
}
