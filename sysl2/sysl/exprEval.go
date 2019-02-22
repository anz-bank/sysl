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
			addItemToValueMap(result, ss.Assign.Name, res)
			// TODO: case *sysl.Expr_Transform_Stmt_Inject:
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
			logrus.Println("Expr Arg is empty")
			return nil
		}
		argValue := Eval(txApp, assign, arg)
		dotValue, hasDot := (*assign)["."]
		defer func() {
			if hasDot {
				(*assign)["."] = dotValue
			}
		}()

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
				setResult := MakeValueSet()
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
					addItemToValueMap(a, "key", MakeValueString(key))
					addItemToValueMap(a, "value", item)
					(*assign)[scopeVar] = a
					res := evalTransformStmts(txApp, assign, x.Transform)
					appendItemToValueList(setResult.GetSet(), res)
				}
				delete(*assign, scopeVar)
				return setResult
			}
			logrus.Printf("Argvalue: %v", argValue)
			(*assign)[scopeVar] = argValue
			res := evalTransformStmts(txApp, assign, x.Transform)
			delete(*assign, scopeVar)
			logrus.Printf("Transform Result: %v", res)
			return res
		}
	case *sysl.Expr_Binexpr:
		return evalBinExpr(txApp, assign, x.Binexpr)
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
		}
		list := MakeValueList()
		for _, argExpr := range x.Call.Arg {
			appendItemToValueList(list.GetList(), Eval(txApp, assign, argExpr))
		}
		return evalGoFunc(x.Call.Func, list)
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
			for _, s := range x.Set.Expr {
				appendItemToValueList(setResult.GetSet(), Eval(txApp, assign, s))
			}
			return setResult
		}
	case *sysl.Expr_List_:
		{
			listResult := MakeValueList()
			for _, s := range x.List.Expr {
				appendItemToValueList(listResult.GetList(), Eval(txApp, assign, s))
			}
			return listResult
		}
	case *sysl.Expr_Unexpr:
		return evalUnaryFunc(x.Unexpr.Op, Eval(txApp, assign, x.Unexpr.Arg))
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
