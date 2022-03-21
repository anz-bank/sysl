package eval

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/pkg/errors"
)

type unaryFunc func(*sysl.Value) *sysl.Value

//nolint:gochecknoglobals
var unaryFunctions = map[sysl.Expr_UnExpr_Op]unaryFunc{
	sysl.Expr_UnExpr_NEG:    unaryNeg,
	sysl.Expr_UnExpr_SINGLE: unarySingle,
	sysl.Expr_UnExpr_STRING: UnaryString,
}

func evalUnaryFunc(op sysl.Expr_UnExpr_Op, arg *sysl.Value) *sysl.Value {
	if x, has := unaryFunctions[op]; has {
		return x(arg)
	}
	panic(errors.Errorf("evalUnaryFunc: operation %v not supported", op))
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

func unarySingle(list *sysl.Value) *sysl.Value {
	if list == nil {
		panic(errors.Errorf("unarySingle received nil parameter"))
	}

	var v []*sysl.Value

	if _, ok := list.Value.(*sysl.Value_List_); ok {
		v = list.GetList().Value
	} else if _, ok := list.Value.(*sysl.Value_Set); ok {
		v = list.GetSet().Value
	} else {
		panic(errors.Errorf("unarySingle expecting List or Set, got %T", list))
	}

	if len(v) != 1 {
		panic(errors.Errorf("unarySingle expecting array length 1, got %v", len(v)))
	}
	return v[0]
}

func UnaryString(arg *sysl.Value) *sysl.Value {
	listfn := func(items []*sysl.Value) *sysl.Value {
		var parts []string
		for _, item := range items {
			parts = append(parts, UnaryString(item).GetS())
		}
		return MakeValueString(fmt.Sprintf("[%s]", strings.Join(parts, ", ")))
	}
	if arg == nil {
		arg = &sysl.Value{Value: &sysl.Value_Null_{}}
	}
	switch x := arg.Value.(type) {
	case *sysl.Value_S:
		return arg
	case *sysl.Value_I:
		return MakeValueString(fmt.Sprintf("%d", x.I))
	case *sysl.Value_B:
		return MakeValueString(map[bool]string{true: "true", false: "false"}[x.B])
	case *sysl.Value_List_:
		return listfn(x.List.GetValue())
	case *sysl.Value_Set:
		return listfn(x.Set.GetValue())
	case *sysl.Value_Map_:
		var parts []string
		for key, val := range x.Map.Items {
			parts = append(parts, fmt.Sprintf("%s: %s", key, UnaryString(val).GetS()))
		}
		sort.Strings(parts)
		return MakeValueString(fmt.Sprintf("{%s}", strings.Join(parts, ", ")))
	}
	return MakeValueString(arg.String())
}
