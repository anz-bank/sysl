package eval

import (
	"sort"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func evalTransformStmts(txApp *sysl.Application, assign Scope, tform *sysl.Expr_Transform) *sysl.Value {
	result := MakeValueMap()

	for _, s := range tform.Stmt {
		switch ss := s.Stmt.(type) {
		case *sysl.Expr_Transform_Stmt_Let:
			logrus.Debugf("Evaluating var %s", ss.Let.Name)
			res := Eval(txApp, assign, ss.Let.Expr)
			assign[ss.Let.Name] = res
			logrus.Tracef("Eval Result %s =: %v\n", ss.Let.Name, res)
		case *sysl.Expr_Transform_Stmt_Assign_:
			logrus.Debugf("Evaluating %s", ss.Assign.Name)
			res := Eval(txApp, assign, ss.Assign.Expr)
			logrus.Tracef("Eval Result %s =:\n\t\t %v:\n", ss.Assign.Name, res)
			AddItemToValueMap(result, ss.Assign.Name, res)
			// TODO: case *sysl.Expr_Transform_Stmt_Inject:
		}
	}
	return result
}

func evalTransformUsingValueList(
	txApp *sysl.Application,
	x *sysl.Expr_Transform,
	assign Scope,
	v []*sysl.Value,
) []*sysl.Value {
	listResult := []*sysl.Value{}
	scopeVar := x.Scopevar
	logrus.Debugf("Scopevar: %s", scopeVar)

	for _, svar := range v {
		assign[scopeVar] = svar
		res := evalTransformStmts(txApp, assign, x)
		listResult = append(listResult, res)
	}
	delete(assign, scopeVar)
	logrus.Tracef("Transform Result (As List/Set): %v", listResult)
	return listResult
}

// Eval expr
// TODO: Add type checks
func Eval(txApp *sysl.Application, assign Scope, e *sysl.Expr) *sysl.Value {
	switch x := e.Expr.(type) {
	case *sysl.Expr_Transform_:
		return evalTransform(txApp, assign, x, e)
	case *sysl.Expr_Binexpr:
		return evalBinExpr(txApp, assign, x.Binexpr)
	case *sysl.Expr_Call_:
		return evalCall(txApp, assign, x)
	case *sysl.Expr_Name:
		return evalName(assign, x)
	case *sysl.Expr_GetAttr_:
		return evalGetAttr(txApp, assign, x)
	case *sysl.Expr_Ifelse:
		return evalIfelse(txApp, assign, x)
	case *sysl.Expr_Literal:
		return x.Literal
	case *sysl.Expr_Set:
		return evalSet(txApp, assign, x)
	case *sysl.Expr_List_:
		return evalList(txApp, assign, x)
	case *sysl.Expr_Unexpr:
		return evalUnaryFunc(x.Unexpr.Op, Eval(txApp, assign, x.Unexpr.Arg))
	default:
		logrus.Warnf("Skipping Expr of type %T\n", x)
		return nil
	}
}

func evalTransform(txApp *sysl.Application, assign Scope, x *sysl.Expr_Transform_, e *sysl.Expr) *sysl.Value {
	logrus.Tracef("Evaluating Transform:\tRet Type: %v", e.Type)
	arg := x.Transform.Arg
	if arg.GetName() == "." {
		logrus.Warn("Expr Arg is empty")
		return nil
	}
	argValue := Eval(txApp, assign, arg)
	dotValue, hasDot := assign["."]
	defer func() {
		if hasDot {
			assign["."] = dotValue
		}
	}()

	switch argValue.Value.(type) {
	case *sysl.Value_Set:
		logrus.Debugf("Evaluation Argvalue as a set: %d times\n", len(argValue.GetSet().Value))
		setResult := MakeValueSet()
		setResult.GetSet().Value = evalTransformUsingValueList(txApp, x.Transform, assign, argValue.GetSet().Value)
		return setResult
	case *sysl.Value_List_:
		logrus.Debugf("Evaluation Argvalue as a list: %d times\n", len(argValue.GetList().Value))
		listResult := MakeValueList()
		listResult.GetList().Value = evalTransformUsingValueList(txApp, x.Transform, assign, argValue.GetList().Value)
		return listResult
	default:
		// HACK: scopevar == '.', then we are not unpacking the map entries
		scopeVar := x.Transform.Scopevar
		if argValue.GetMap() != nil && scopeVar != "." {
			// TODO: add check that return type is defined as 'set of ...'
			resultList := &sysl.Value_List{}
			// Sort keys, to get stable output
			var keys []string
			for key := range argValue.GetMap().Items {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			items := argValue.GetMap().Items
			for _, key := range keys {
				item := items[key]
				logrus.Debugf("Evaluation Argvalue %s", key)
				a := MakeValueMap()
				logrus.Tracef("Evaluation Argvalue as a map: key=(%s), value=(%v)\n", key, item)
				AddItemToValueMap(a, "key", MakeValueString(key))
				AddItemToValueMap(a, "value", item)
				assign[scopeVar] = a
				res := evalTransformStmts(txApp, assign, x.Transform)
				AppendItemToValueList(resultList, res)
			}
			delete(assign, scopeVar)
			if e.Type.GetSet() != nil {
				return &sysl.Value{
					Value: &sysl.Value_Set{
						Set: resultList,
					},
				}
			}
			return &sysl.Value{
				Value: &sysl.Value_List_{
					List: resultList,
				},
			}
		}
		logrus.Tracef("Argvalue: %v", argValue)
		assign[scopeVar] = argValue
		res := evalTransformStmts(txApp, assign, x.Transform)
		delete(assign, scopeVar)
		logrus.Tracef("Transform Result: %v", res)
		return res
	}
}

func evalCall(txApp *sysl.Application, assign Scope, x *sysl.Expr_Call_) *sysl.Value {
	if callTransform, has := txApp.Views[x.Call.Func]; has {
		logrus.Debugf("Calling View %s\n", x.Call.Func)
		params := callTransform.Param
		if len(params) != len(x.Call.Arg) {
			logrus.Warnf("Skipping Calling func(%s), args mismatch, %d args passed, %d required\n",
				x.Call.Func, len(x.Call.Arg), len(params))
			return nil
		}
		callScope := make(Scope)

		for i, argExpr := range x.Call.Arg {
			// TODO: Add type checks
			callScope[params[i].Name] = Eval(txApp, assign, argExpr)
		}
		if callTransform.Expr.Type == nil {
			callTransform.Expr.Type = callTransform.RetType
		}
		return Eval(txApp, callScope, callTransform.Expr)
	} else if strings.HasPrefix(x.Call.Func, ".") {
		switch x.Call.Func[1:] {
		case "count":
			argExpr := Eval(txApp, assign, x.Call.Arg[0])
			if argExpr.GetList() != nil {
				return MakeValueI64(int64(len(argExpr.GetList().Value)))
			} else if argExpr.GetSet() != nil {
				return MakeValueI64(int64(len(argExpr.GetSet().Value)))
			}
			panic(errors.Errorf("Unexpected arg type: %v", x.Call.Arg))
		default:
			panic(errors.Errorf("Unimplemented function: %s", x.Call.Func))
		}
	}
	list := MakeValueList()
	for _, argExpr := range x.Call.Arg {
		AppendItemToValueList(list.GetList(), Eval(txApp, assign, argExpr))
	}
	return evalGoFunc(x.Call.Func, list)
}

func evalName(assign Scope, x *sysl.Expr_Name) *sysl.Value {
	val, has := assign[x.Name]
	if !has {
		logrus.Errorf("Key: %s does not exist in scope", x.Name)
	}
	return val
}

func evalGetAttr(txApp *sysl.Application, assign Scope, x *sysl.Expr_GetAttr_) *sysl.Value {
	logrus.Tracef("Evaluating x: %v:", x)
	arg := Eval(txApp, assign, x.GetAttr.Arg)
	if arg.GetMap() == nil {
		panic(errors.Errorf("%v", arg))
	}
	val, has := arg.GetMap().Items[x.GetAttr.Attr]
	logrus.Tracef("GetAttribute: %v result: %v: ", has, val)
	if !has {
		logrus.Warnf("Failed to get key: %s\nMap has following keys:", x.GetAttr.Attr)
		for key := range arg.GetMap().Items {
			logrus.Warnf("\t%s", key)
		}
		return &sysl.Value{
			Value: &sysl.Value_Null_{
				Null: &sysl.Value_Null{},
			},
		}
	}
	return val
}

func evalIfelse(txApp *sysl.Application, assign Scope, x *sysl.Expr_Ifelse) *sysl.Value {
	cond := Eval(txApp, assign, x.Ifelse.Cond)
	if cond.GetB() {
		return Eval(txApp, assign, x.Ifelse.IfTrue)
	}
	return Eval(txApp, assign, x.Ifelse.IfFalse)
}

func evalSet(txApp *sysl.Application, assign Scope, x *sysl.Expr_Set) *sysl.Value {
	{
		setResult := MakeValueSet()
		for _, s := range x.Set.Expr {
			AppendItemToValueList(setResult.GetSet(), Eval(txApp, assign, s))
		}
		return setResult
	}
}

func evalList(txApp *sysl.Application, assign Scope, x *sysl.Expr_List_) *sysl.Value {
	{
		listResult := MakeValueList()
		for _, s := range x.List.Expr {
			AppendItemToValueList(listResult.GetList(), Eval(txApp, assign, s))
		}
		return listResult
	}
}

// EvaluateView evaluate the view using the Scope
func EvaluateView(mod *sysl.Module, appName, viewName string, s Scope) *sysl.Value {
	txApp := mod.Apps[appName]
	view := txApp.Views[viewName]
	if view.Expr.Type == nil {
		view.Expr.Type = view.RetType
	}
	return Eval(txApp, s, view.Expr)
}
