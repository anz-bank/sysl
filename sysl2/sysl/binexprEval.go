package main

import (
	"fmt"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
)

type evalValueFunc func(*sysl.Value, *sysl.Value) *sysl.Value
type evalExprFunc func(txApp *sysl.Application, assign *Scope, list *sysl.Value, scopeVar string, rhs *sysl.Expr) *sysl.Value

// EvalStrategy interface to evaluate binary expr
type EvalStrategy interface {
	eval(*sysl.Application, *Scope, *sysl.Expr_BinExpr) *sysl.Value
}

// DefaultBinExprStrategy is to evaluate lhs and rhs expr's first and then pass it to fn
type DefaultBinExprStrategy struct{}

// Eval_LHS_Over_RhsStrategy binds rhs expression over each element of the LHS.
// Assumes lhs is a collection
type Eval_LHS_Over_RhsStrategy struct{}

var functionEvalStrategy = map[sysl.Expr_BinExpr_Op]EvalStrategy{
	sysl.Expr_BinExpr_EQ:      DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_ADD:     DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_SUB:     DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_MUL:     DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_MOD:     DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_DIV:     DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_IN:      DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_BITOR:   DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_GT:      DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_LT:      DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_GE:      DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_LE:      DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_NE:      DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_AND:     DefaultBinExprStrategy{},
	sysl.Expr_BinExpr_FLATTEN: Eval_LHS_Over_RhsStrategy{},
	sysl.Expr_BinExpr_WHERE:   Eval_LHS_Over_RhsStrategy{},
}

// key = op, lhs & rhs types
var valueFunctions = map[string]evalValueFunc{
	makeKey(sysl.Expr_BinExpr_ADD, VALUE_INT, VALUE_INT):       addInt64,
	makeKey(sysl.Expr_BinExpr_ADD, VALUE_STRING, VALUE_STRING): addString,
	makeKey(sysl.Expr_BinExpr_AND, VALUE_BOOL, VALUE_BOOL):     andBool,
	makeKey(sysl.Expr_BinExpr_BITOR, VALUE_LIST, VALUE_LIST):   concatList,
	makeKey(sysl.Expr_BinExpr_BITOR, VALUE_SET, VALUE_SET):     setUnion,
	makeKey(sysl.Expr_BinExpr_DIV, VALUE_INT, VALUE_INT):       divInt64,
	makeKey(sysl.Expr_BinExpr_EQ, VALUE_BOOL, VALUE_BOOL):      cmpBool,
	makeKey(sysl.Expr_BinExpr_EQ, VALUE_INT, VALUE_INT):        cmpInt,
	makeKey(sysl.Expr_BinExpr_EQ, VALUE_INT, VALUE_NULL):       cmpNullFalse,
	makeKey(sysl.Expr_BinExpr_EQ, VALUE_NULL, VALUE_NULL):      cmpNullTrue,
	makeKey(sysl.Expr_BinExpr_EQ, VALUE_NULL, VALUE_INT):       cmpNullFalse,
	makeKey(sysl.Expr_BinExpr_EQ, VALUE_NULL, VALUE_STRING):    cmpNullFalse,
	makeKey(sysl.Expr_BinExpr_EQ, VALUE_STRING, VALUE_NULL):    cmpNullFalse,
	makeKey(sysl.Expr_BinExpr_EQ, VALUE_STRING, VALUE_STRING):  cmpString,
	makeKey(sysl.Expr_BinExpr_GE, VALUE_INT, VALUE_INT):        geInt64,
	makeKey(sysl.Expr_BinExpr_GT, VALUE_INT, VALUE_INT):        gtInt64,
	makeKey(sysl.Expr_BinExpr_IN, VALUE_STRING, VALUE_LIST):    stringInList,
	makeKey(sysl.Expr_BinExpr_IN, VALUE_STRING, VALUE_NULL):    stringInNull,
	makeKey(sysl.Expr_BinExpr_IN, VALUE_STRING, VALUE_SET):     stringInSet,
	makeKey(sysl.Expr_BinExpr_LE, VALUE_INT, VALUE_INT):        leInt64,
	makeKey(sysl.Expr_BinExpr_LT, VALUE_INT, VALUE_INT):        ltInt64,
	makeKey(sysl.Expr_BinExpr_MOD, VALUE_INT, VALUE_INT):       modInt64,
	makeKey(sysl.Expr_BinExpr_MUL, VALUE_INT, VALUE_INT):       mulInt64,
	makeKey(sysl.Expr_BinExpr_NE, VALUE_BOOL, VALUE_BOOL):      not(cmpBool),
	makeKey(sysl.Expr_BinExpr_NE, VALUE_INT, VALUE_INT):        not(cmpInt),
	makeKey(sysl.Expr_BinExpr_NE, VALUE_NULL, VALUE_NULL):      cmpNullFalse,
	makeKey(sysl.Expr_BinExpr_NE, VALUE_NULL, VALUE_STRING):    cmpNullTrue,
	makeKey(sysl.Expr_BinExpr_NE, VALUE_STRING, VALUE_NULL):    cmpNullTrue,
	makeKey(sysl.Expr_BinExpr_NE, VALUE_STRING, VALUE_STRING):  not(cmpString),
	makeKey(sysl.Expr_BinExpr_SUB, VALUE_INT, VALUE_INT):       subInt64,
}

// key = op, outer container_type, inner container type
var exprFunctions = map[string]evalExprFunc{
	makeKey(sysl.Expr_BinExpr_FLATTEN, VALUE_LIST, VALUE_NO_ARG): flattenListList, // empty list
	makeKey(sysl.Expr_BinExpr_FLATTEN, VALUE_LIST, VALUE_LIST):   flattenListList,
	makeKey(sysl.Expr_BinExpr_FLATTEN, VALUE_LIST, VALUE_SET):    flattenListSet,
	makeKey(sysl.Expr_BinExpr_FLATTEN, VALUE_LIST, VALUE_MAP):    flattenListMap,
	makeKey(sysl.Expr_BinExpr_FLATTEN, VALUE_SET, VALUE_MAP):     flattenSetMap,
	makeKey(sysl.Expr_BinExpr_FLATTEN, VALUE_SET, VALUE_SET):     flattenSetSet,
	makeKey(sysl.Expr_BinExpr_WHERE, VALUE_LIST, VALUE_NO_ARG):   whereList,
	makeKey(sysl.Expr_BinExpr_WHERE, VALUE_LIST, VALUE_MAP):      whereList,
	makeKey(sysl.Expr_BinExpr_WHERE, VALUE_MAP, VALUE_NO_ARG):    whereMap,
	makeKey(sysl.Expr_BinExpr_WHERE, VALUE_SET, VALUE_NO_ARG):    whereSet,
	makeKey(sysl.Expr_BinExpr_WHERE, VALUE_SET, VALUE_MAP):       whereSet,
	makeKey(sysl.Expr_BinExpr_WHERE, VALUE_SET, VALUE_INT):       whereSet,
	makeKey(sysl.Expr_BinExpr_WHERE, VALUE_SET, VALUE_FLOAT):     whereSet,
	makeKey(sysl.Expr_BinExpr_WHERE, VALUE_SET, VALUE_STRING):    whereSet,
}

func makeKey(op sysl.Expr_BinExpr_Op, lhs, rhs valueType) string {
	return fmt.Sprintf("%s_%s_%s", op, lhs.String(), rhs.String())
}

func (op DefaultBinExprStrategy) eval(txApp *sysl.Application, assign *Scope, binexpr *sysl.Expr_BinExpr) *sysl.Value {
	lhs_v := Eval(txApp, assign, binexpr.Lhs)
	rhs_v := Eval(txApp, assign, binexpr.Rhs)
	key := makeKey(binexpr.Op, getValueType(lhs_v), getValueType(rhs_v))

	if f, has := valueFunctions[key]; has {
		return f(lhs_v, rhs_v)
	}
	panic(errors.Errorf("Unsupported operation:DefaultBinExprStrategy: %s", key))
}

func (op Eval_LHS_Over_RhsStrategy) eval(txApp *sysl.Application, assign *Scope, binexpr *sysl.Expr_BinExpr) *sysl.Value {
	lhs_v := Eval(txApp, assign, binexpr.Lhs)
	vType := getValueType(lhs_v)
	itemType := getContainedType(lhs_v)
	key := makeKey(binexpr.Op, vType, itemType)

	if f, has := exprFunctions[key]; has {
		return f(txApp, assign, lhs_v, binexpr.Scopevar, binexpr.Rhs)
	}
	panic(errors.Errorf("Unsupported operation: Cannot Execute %s", key))
}

func evalBinExpr(txApp *sysl.Application, assign *Scope, binexpr *sysl.Expr_BinExpr) *sysl.Value {
	if strategy, has := functionEvalStrategy[binexpr.Op]; has {
		return strategy.eval(txApp, assign, binexpr)
	}
	panic(errors.Errorf("Unsupported operation: %s", binexpr.Op))
}
