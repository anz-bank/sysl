package main

import (
	"sort"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
)

func evalTransformStmts(txApp *sysl.Application, assign *Scope, tform *sysl.Expr_Transform) *sysl.Value {
	result := MakeValueMap()

	for _, s := range tform.Stmt {
		switch ss := s.Stmt.(type) {
		case *sysl.Expr_Transform_Stmt_Let:
			res := Eval(txApp, assign, ss.Let.Expr)
			(*assign)[ss.Let.Name] = res
			logrus.Printf("Eval %s: %v\n", ss.Let.Name, res)
		case *sysl.Expr_Transform_Stmt_Assign_:
			logrus.Printf("Eval %s:\n", ss.Assign.Name)
			res := Eval(txApp, assign, ss.Assign.Expr)
			logrus.Printf("Eval %s ==\n\t\t %v:\n", ss.Assign.Name, res)
			addItemToValueMap(result.GetMap(), ss.Assign.Name, res)
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
	logrus.Printf("Transform Result (As List/Set): %v", listResult)
	return listResult
}

// Eval expr
// TODO: Add type checks
func Eval(txApp *sysl.Application, assign *Scope, e *sysl.Expr) *sysl.Value {
	switch x := e.Expr.(type) {
	case *sysl.Expr_Transform_:
		logrus.Printf("Evaluating Transform:\tRet Type: %v\n", e.Type)
		logrus.Println("Evaluating Transform")
		arg := x.Transform.Arg
		if arg.GetName() == "." {
			// TODO: return error
			logrus.Println("Expr Arg is empty")
			return nil
		}
		argValue := Eval(txApp, assign, arg)

		switch argValue.Value.(type) {
		case *sysl.Value_Set:
			logrus.Printf("Evaluation Argvalue as a set: %d times\n", len(argValue.GetSet().Value))
			return evalTransformUsingValueList(txApp, x.Transform, assign, argValue.GetSet().Value)
		case *sysl.Value_List_:
			logrus.Printf("Evaluation Argvalue as a list: %d times\n", len(argValue.GetList().Value))
			return evalTransformUsingValueList(txApp, x.Transform, assign, argValue.GetList().Value)
		default:
			// HACK: scopevar == '.', then we are not unpacking the map entries
			scopeVar := x.Transform.Scopevar
			if argValue.GetMap() != nil && scopeVar != "." {
				// TODO: add check that return type is defined as 'set of ...'
				listResult := MakeValueSet()
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
					logrus.Printf("Evaluation Argvalue as a map: key=(%s), value=(%v)\n", key, item)
					addItemToValueMap(a.GetMap(), "key", MakeValueString(key))
					addItemToValueMap(a.GetMap(), "value", item)
					(*assign)[scopeVar] = a
					res := evalTransformStmts(txApp, assign, x.Transform)
					listResult.GetSet().Value = append(listResult.GetSet().Value, res)
				}
				delete(*assign, scopeVar)
				return listResult
			} else {
				logrus.Printf("Argvalue: %v", argValue)
			}

			(*assign)[scopeVar] = argValue
			res := evalTransformStmts(txApp, assign, x.Transform)
			delete(*assign, scopeVar)
			logrus.Printf("Transform Result: %v", res)
			return res
		}
	case *sysl.Expr_Binexpr:
		lhs_v := Eval(txApp, assign, x.Binexpr.Lhs)
		rhs_v := Eval(txApp, assign, x.Binexpr.Rhs)
		return evalBinExpr(x.Binexpr.Op, lhs_v, rhs_v)
	case *sysl.Expr_Call_:
		if callTransform, has := txApp.Views[x.Call.Func]; has {
			logrus.Printf("Calling %s\n", x.Call.Func)
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
		} else {
			// TODO: see if the func is in strings package
		}
	case *sysl.Expr_Name:
		return (*assign)[x.Name]
	case *sysl.Expr_GetAttr_:
		logrus.Printf("Evaluating x: %v:\n", x)
		arg := Eval(txApp, assign, x.GetAttr.Arg)
		val, has := arg.GetMap().Items[x.GetAttr.Attr]
		logrus.Printf("GetAttribute: %v result: %v: ", has, val)
		return val
	case *sysl.Expr_Ifelse:
		cond := Eval(txApp, assign, x.Ifelse.Cond)
		if cond.GetB() {
			return Eval(txApp, assign, x.Ifelse.IfTrue)
		}
		return Eval(txApp, assign, x.Ifelse.IfFalse)
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
