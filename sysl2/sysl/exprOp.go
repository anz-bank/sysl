package main

import (
	"fmt"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
)

type evalFunc func(*sysl.Value, *sysl.Value) *sysl.Value

// OpMap is map of string to Function of type BinFunc
type OpMap map[string]evalFunc

// Consts for ValueTypes
const (
	VALUE_NO_ARG     = -1
	VALUE_BOOL   int = iota
	VALUE_INT
	VALUE_FLOAT
	VALUE_STRING
	VALUE_STRING_DECIMAL
	VALUE_LIST
	VALUE_MAP
	VALUE_SET
)

// FuncMap ...
var FuncMap OpMap

// OpArgs represents valid combination of an operator, lhs and rhs value types
type OpArgs struct {
	op       string // String representation of the function name, See Expr_BinExpr_Op_value or Expr_RelExpr_Op_name.
	lhs, rhs int    // VALUE_BOOL etc
	f        evalFunc
}

func addInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() + rhs.GetI())
}
func subInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() - rhs.GetI())
}
func mulInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() * rhs.GetI())
}
func divInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() / rhs.GetI())
}
func modInt64(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueI64(lhs.GetI() % rhs.GetI())
}

func addString(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueString(lhs.GetS() + rhs.GetS())
}

func cmpString(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetS() == rhs.GetS())
}

func cmpBool(lhs, rhs *sysl.Value) *sysl.Value {
	return MakeValueBool(lhs.GetB() == rhs.GetB())
}

func setUnion(lhs, rhs *sysl.Value) *sysl.Value {
	res := MakeValueSet()
	appendListToValueList(res.GetSet(), lhs.GetSet())
	appendListToValueList(res.GetSet(), rhs.GetSet())
	logrus.Printf("Union set: lhs %d, rhs %d res %d\n", len(lhs.GetSet().Value), len(rhs.GetSet().Value), len(res.GetSet().Value))
	return res
}

func stringInSet(lhs, rhs *sysl.Value) *sysl.Value {
	str := lhs.GetS()
	for _, v := range rhs.GetSet().Value {
		if str == v.GetS() {
			return MakeValueBool(true)
		}
	}
	return MakeValueBool(false)
}

func stringInList(lhs, rhs *sysl.Value) *sysl.Value {
	str := lhs.GetS()
	for _, v := range rhs.GetList().Value {
		if str == v.GetS() {
			return MakeValueBool(true)
		}
	}
	return MakeValueBool(false)
}

func getValueType(v *sysl.Value) int {
	if v == nil {
		return VALUE_NO_ARG
	}
	switch v.Value.(type) {
	case *sysl.Value_B:
		return VALUE_BOOL
	case *sysl.Value_I:
		return VALUE_INT
	case *sysl.Value_D:
		return VALUE_FLOAT
	case *sysl.Value_S:
		return VALUE_STRING
	case *sysl.Value_Decimal:
		return VALUE_STRING_DECIMAL
	case *sysl.Value_Set:
		return VALUE_SET
	case *sysl.Value_List_:
		return VALUE_SET
	case *sysl.Value_Map_:
		return VALUE_MAP
	default:
		panic("exprOp: getValueType: unhandled type")
	}
}

func makeKey(op_name string, lhs, rhs int) string {
	return fmt.Sprintf("%s_%d_%d", op_name, lhs, rhs)
}

func addToFuncMap(op_name string, lhs, rhs int, f evalFunc) {
	if _, has := FuncMap[makeKey(op_name, lhs, rhs)]; has {
		panic("adding duplicate entries in the func map")
	}
	FuncMap[makeKey(op_name, lhs, rhs)] = f
}

func evalBinExpr(op sysl.Expr_BinExpr_Op, lhs, rhs *sysl.Value) *sysl.Value {
	key := makeKey(sysl.Expr_BinExpr_Op_name[int32(op)], getValueType(lhs), getValueType(rhs))
	if fn, has := FuncMap[key]; has {
		return fn(lhs, rhs)
	}
	panic("unsupported operation: " + key)
}

func addBinaryOps() {
	var ops = []OpArgs{
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_EQ)], VALUE_BOOL, VALUE_BOOL, cmpBool},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_ADD)], VALUE_INT, VALUE_INT, addInt64},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_SUB)], VALUE_INT, VALUE_INT, subInt64},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_MUL)], VALUE_INT, VALUE_INT, mulInt64},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_DIV)], VALUE_INT, VALUE_INT, divInt64},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_MOD)], VALUE_INT, VALUE_INT, modInt64},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_ADD)], VALUE_STRING, VALUE_STRING, addString},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_EQ)], VALUE_STRING, VALUE_STRING, cmpString},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_IN)], VALUE_STRING, VALUE_SET, stringInSet},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_IN)], VALUE_STRING, VALUE_LIST, stringInList},
		OpArgs{sysl.Expr_BinExpr_Op_name[int32(sysl.Expr_BinExpr_BITOR)], VALUE_SET, VALUE_SET, setUnion},
	}
	for _, op := range ops {
		addToFuncMap(op.op, op.lhs, op.rhs, op.f)
	}
}

func init() {
	FuncMap = make(OpMap)
	addBinaryOps()
}
