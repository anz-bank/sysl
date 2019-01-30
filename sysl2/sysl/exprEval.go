package main

import (
	"sort"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
)

func evalPrimitiveI(op sysl.Expr_BinExpr_Op, lhs int64, rhs int64) *sysl.Value {
	var result int64
	switch op {
	case sysl.Expr_BinExpr_ADD:
		result = lhs + rhs
	case sysl.Expr_BinExpr_SUB:
		result = lhs - rhs
	case sysl.Expr_BinExpr_MUL:
		result = lhs * rhs
	case sysl.Expr_BinExpr_DIV:
		result = lhs / rhs
	}
	return MakeValueI64(result)
}

func evalPrimitiveS(op sysl.Expr_BinExpr_Op, lhs string, rhs string) *sysl.Value {
	var result string
	switch op {
	case sysl.Expr_BinExpr_ADD:
		result = lhs + rhs
	default:
		panic("binary ops on strings not supported!")
	}
	return MakeValueString(result)
}

func evalSets(op sysl.Expr_BinExpr_Op, lhs *sysl.Value_List, rhs *sysl.Value_List) *sysl.Value {
	switch op {
	case sysl.Expr_BinExpr_ADD:
		// TODO: add checks for uniqueness
		res := MakeValueSet()
		appendListToValueList(res.GetSet(), lhs)
		appendListToValueList(res.GetSet(), rhs)
		return res
	default:
		logrus.Warningln("binary ops on sets not supported!")
	}
	return nil
}

func evalTransformStmts(txApp *sysl.Application, assign *Scope, tform *sysl.Expr_Transform) *sysl.Value {
	local := make(Scope)
	result := MakeValueMap()

	for _, s := range tform.Stmt {
		switch ss := s.Stmt.(type) {
		case *sysl.Expr_Transform_Stmt_Let:
			local[ss.Let.Name] = Eval(txApp, assign, ss.Let.Expr)
		case *sysl.Expr_Transform_Stmt_Assign_:
			addItemToValueMap(result.GetMap(), ss.Assign.Name, Eval(txApp, assign, ss.Assign.Expr))
			// case *sysl.Expr_Transform_Stmt_Inject:
		}
	}
	return result
}

func evalTransformUsingValueList(txApp *sysl.Application, x *sysl.Expr_Transform, assign *Scope, v []*sysl.Value) *sysl.Value {
	listResult := MakeValueList()
	scopeVar := x.Scopevar

	for _, svar := range v {
		(*assign)[scopeVar] = svar
		res := evalTransformStmts(txApp, assign, x)
		listResult.GetList().Value = append(listResult.GetList().Value, res)
	}
	delete(*assign, scopeVar)
	return listResult
}

// Eval expr
// TODO: Add type checks
func Eval(txApp *sysl.Application, assign *Scope, e *sysl.Expr) *sysl.Value {
	switch x := e.Expr.(type) {
	case *sysl.Expr_Transform_:
		arg := x.Transform.Arg
		if arg.GetName() == "." {
			// TODO: return error
			logrus.Println("Expr Arg is empty")
			return nil
		}
		argValue := Eval(txApp, assign, arg)

		switch argValue.Value.(type) {
		case *sysl.Value_Set:
			return evalTransformUsingValueList(txApp, x.Transform, assign, argValue.GetSet().Value)
		case *sysl.Value_List_:
			return evalTransformUsingValueList(txApp, x.Transform, assign, argValue.GetList().Value)
		default:
			// HACK: scopevar == '.', then we are not unpacking the map entries
			scopeVar := x.Transform.Scopevar
			if argValue.GetMap() != nil && scopeVar != "." {
				listResult := MakeValueList()
				// Sort keys, to get stable output
				var keys []string
				for key := range argValue.GetMap().Items {
					keys = append(keys, key)
				}
				sort.Strings(keys)
				items := argValue.GetMap().Items
				for _, key := range keys {
					item := items[key]
					a := MakeValueMap()
					addItemToValueMap(a.GetMap(), "key", MakeValueString(key))
					addItemToValueMap(a.GetMap(), "value", item)
					(*assign)[scopeVar] = a
					res := evalTransformStmts(txApp, assign, x.Transform)
					listResult.GetList().Value = append(listResult.GetList().Value, res)
				}
				delete(*assign, scopeVar)
				return listResult
			}

			(*assign)[scopeVar] = argValue
			res := evalTransformStmts(txApp, assign, x.Transform)
			delete(*assign, scopeVar)
			return res
		}
	case *sysl.Expr_Binexpr:
		lhs_v := Eval(txApp, assign, x.Binexpr.Lhs)
		rhs_v := Eval(txApp, assign, x.Binexpr.Rhs)
		switch lhs_v.Value.(type) {
		case *sysl.Value_I:
			{
				return evalPrimitiveI(x.Binexpr.Op, lhs_v.GetI(), rhs_v.GetI())
			}
		case *sysl.Value_S:
			{
				return evalPrimitiveS(x.Binexpr.Op, lhs_v.GetS(), rhs_v.GetS())
			}
		case *sysl.Value_Set:
			return evalSets(x.Binexpr.Op, lhs_v.GetSet(), rhs_v.GetSet())
		default:
			logrus.Warnf("Skipping: Binary Op: %d for lhs(%T), rhs(%T)", x.Binexpr.Op, lhs_v.Value, rhs_v.Value)
			return nil
		}
	case *sysl.Expr_Call_:
		if callTransform, has := txApp.Views[x.Call.Func]; has {

			params := callTransform.Param
			if len(params) != len(x.Call.Arg) {
				logrus.Warnf("Skipping Calling func(%s), args mismatch, %d args passed, %d required\n", x.Call.Func, len(x.Call.Arg), len(params))
				return nil
			}
			callScope := make(Scope)

			for i, argExpr := range x.Call.Arg {
				// TODO: Add type checks
				callScope[params[i].Name] = Eval(txApp, assign, argExpr)
			}
			return Eval(txApp, &callScope, callTransform.Expr)
		}
	case *sysl.Expr_Name:
		return (*assign)[x.Name]
	case *sysl.Expr_GetAttr_:
		logrus.Printf("Evaluating x: %v:\n", x)
		arg := Eval(txApp, assign, x.GetAttr.Arg)
		val := arg.GetMap().Items[x.GetAttr.Attr]
		logrus.Printf("result: %v: ", val)
		return val
	case *sysl.Expr_Literal:
		return x.Literal
	case *sysl.Expr_Set:
		{
			setResult := MakeValueSet()
			v := []*sysl.Value{}
			for _, s := range x.Set.Expr {
				v = append(v, Eval(txApp, assign, s))
			}
			setResult.GetSet().Value = v
			return setResult
		}
	default:
		logrus.Warnf("Skipping Expr of type %T\n", x)
	}
	return nil
}

// EvalView evaluate the view using the Scope
func EvalView(mod *sysl.Module, appName, viewName string, s *Scope) *sysl.Value {
	txApp := mod.Apps[appName]
	return Eval(txApp, s, txApp.Views[viewName].Expr)
}
