package eval

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/parse"
	sysl "github.com/anz-bank/sysl/pkg/proto_old"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func evalTransformStmts(ee *exprEval, assign Scope, tform *sysl.Expr_Transform) *sysl.Value {
	result := MakeValueMap()

	for _, s := range tform.Stmt {
		switch ss := s.Stmt.(type) {
		case *sysl.Expr_Transform_Stmt_Let:
			logrus.Debugf("Evaluating var %s", ss.Let.Name)
			res := Eval(ee, assign, ss.Let.Expr)
			assign[ss.Let.Name] = res
			logrus.Tracef("Eval Result %s =: %v\n", ss.Let.Name, res)
		case *sysl.Expr_Transform_Stmt_Assign_:
			logrus.Debugf("Evaluating %s", ss.Assign.Name)
			res := Eval(ee, assign, ss.Assign.Expr)
			logrus.Tracef("Eval Result %s =:\n\t\t %v:\n", ss.Assign.Name, res)
			AddItemToValueMap(result, ss.Assign.Name, res)
			// TODO: case *sysl.Expr_Transform_Stmt_Inject:
		}
	}
	if text, has := assign[parse.TemplateImpliedResult]; has {
		return text
	}
	return result
}

func setAppender(collection []*sysl.Value, newVal *sysl.Value) []*sysl.Value {
	found := false
	for _, result := range collection {
		if proto.Equal(newVal, result) {
			found = true
			break
		}
	}

	if !found {
		collection = append(collection, newVal)
	}
	return collection
}

func listAppender(collection []*sysl.Value, newVal *sysl.Value) []*sysl.Value {
	return append(collection, newVal)
}

func evalTransformUsingAppender(
	ee *exprEval,
	x *sysl.Expr_Transform,
	assign Scope,
	v []*sysl.Value,
	appender func(collection []*sysl.Value, newVal *sysl.Value) []*sysl.Value,
) []*sysl.Value {
	listResult := []*sysl.Value{}
	scopeVar := x.Scopevar

	for _, svar := range v {
		assign[scopeVar] = svar
		res := evalTransformStmts(ee, assign, x)
		listResult = appender(listResult, res)
	}
	delete(assign, scopeVar)
	logrus.Tracef("Transform Result (As List/Set): %v", listResult)
	return listResult
}

func evalTransformUsingValueSet(
	ee *exprEval,
	x *sysl.Expr_Transform,
	assign Scope,
	v []*sysl.Value,
) []*sysl.Value {
	return evalTransformUsingAppender(ee, x, assign, v, setAppender)
}

func evalTransformUsingValueList(
	ee *exprEval,
	x *sysl.Expr_Transform,
	assign Scope,
	v []*sysl.Value,
) []*sysl.Value {
	return evalTransformUsingAppender(ee, x, assign, v, listAppender)
}

// Eval expr
// TODO: Add type checks
func Eval(ee *exprEval, assign Scope, e *sysl.Expr) *sysl.Value {
	val, _ := ee.eval(assign, e) //nolint:errcheck
	return val
}

func (ee *exprEval) evalTransform(assign Scope, x *sysl.Expr_Transform_, e *sysl.Expr) *sysl.Value {
	logrus.Tracef("Evaluating Transform:\tRet Type: %v", e.Type)
	arg := x.Transform.Arg
	if arg.GetName() == "." {
		logrus.Warn("Expr Arg is empty")
		return nil
	}
	argValue := Eval(ee, assign, arg)
	dotValue, hasDot := assign["."]
	defer func() {
		if hasDot {
			assign["."] = dotValue
		}
	}()

	switch argValue.Value.(type) {
	case *sysl.Value_Set, *sysl.Value_List_:
		switch e.Type.Type.(type) {
		case *sysl.Type_Set:
			logrus.Debugf("Evaluation Argvalue as a set: %d times\n", len(GetValueSlice(argValue)))
			setResult := MakeValueSet()
			setResult.GetSet().Value = evalTransformUsingValueSet(ee, x.Transform, assign, GetValueSlice(argValue))
			return setResult
		default:
			logrus.Debugf("Evaluation Argvalue as a list: %d times\n", len(GetValueSlice(argValue)))
			listResult := MakeValueList()
			listResult.GetList().Value = evalTransformUsingValueList(ee, x.Transform, assign, GetValueSlice(argValue))
			return listResult
		}
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
				res := evalTransformStmts(ee, assign, x.Transform)
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
		res := evalTransformStmts(ee, assign, x.Transform)
		delete(assign, scopeVar)
		logrus.Tracef("Transform Result: %v", res)
		return res
	}
}

func evalCall(ee *exprEval, assign Scope, x *sysl.Expr_Call_) *sysl.Value {
	if callTransform, has := ee.txApp.Views[x.Call.Func]; has {
		logrus.Debugf("Calling View %s", x.Call.Func)
		params := callTransform.Param
		if len(params) != len(x.Call.Arg) {
			logrus.Warnf("Skipping Calling func(%s), args mismatch, %d args passed, %d required\n",
				x.Call.Func, len(x.Call.Arg), len(params))
			return nil
		}
		callScope := make(Scope)

		for i, argExpr := range x.Call.Arg {
			// TODO: Add type checks
			callScope[params[i].Name] = Eval(ee, assign, argExpr)
		}
		if callTransform.Expr.Type == nil {
			callTransform.Expr.Type = callTransform.RetType
		}
		return Eval(ee, callScope, callTransform.Expr)
	} else if strings.HasPrefix(x.Call.Func, ".") {
		switch x.Call.Func[1:] {
		case "count":
			argExpr := Eval(ee, assign, x.Call.Arg[0])
			switch t := argExpr.Value.(type) {
			case *sysl.Value_List_:
				return MakeValueI64(int64(len(t.List.Value)))
			case *sysl.Value_Set:
				return MakeValueI64(int64(len(t.Set.Value)))
			case *sysl.Value_Map_:
				return MakeValueI64(int64(len(t.Map.Items)))
			default:
				panic(errors.Errorf("Unexpected arg type: %v", x.Call.Arg))
			}
		default:
			panic(errors.Errorf("Unimplemented function: %s", x.Call.Func))
		}
	}
	list := MakeValueList()
	for _, argExpr := range x.Call.Arg {
		AppendItemToValueList(list.GetList(), Eval(ee, assign, argExpr))
	}
	return evalGoFunc(x.Call.Func, list)
}

// isInternalMap checks if this Value is created inplicitly inside evalTransform()
func isInternalMap(val *sysl.Value_Map) bool {
	founds := 0
	for key := range val.Items {
		switch key {
		case "key", "value":
			founds++
		default:
			return false
		}
	}
	return founds == 2
}

func evalName(ee *exprEval, assign Scope, x *sysl.Expr_Name) *sysl.Value {
	val, has := assign[x.Name]
	if !has {
		if x.Name == parse.TemplateImpliedResult {
			val = MakeValueString("")
			assign[x.Name] = val
		} else {
			ee.LogEntry().Errorf("Key: %s does not exist in scope", x.Name)
		}
	}
	return val
}

func evalGetAttr(ee *exprEval, assign Scope, x *sysl.Expr_GetAttr_) *sysl.Value {
	entry := ee.LogEntry()
	arg := Eval(ee, assign, x.GetAttr.Arg)
	if arg.GetMap() == nil {
		panic(errors.Errorf("%v", arg))
	}

	if isInternalMap(arg.GetMap()) {
		switch key := x.GetAttr.Attr; key {
		case "value":
			ee.LogEntry().Debugf("Unnecessary use of .value")
			fallthrough
		case "key":
			return arg.GetMap().Items[key]
		}
		arg = arg.GetMap().Items["value"]
	}
	val, has := arg.GetMap().Items[x.GetAttr.Attr]
	entry.Tracef("GetAttribute: %v result: %v: ", has, val)
	if !has {
		entry.Warnf("Failed to get key: %s. Map has following keys:", x.GetAttr.Attr)
		for key := range arg.GetMap().Items {
			entry.Warnf("\t%s", key)
		}
		return &sysl.Value{
			Value: &sysl.Value_Null_{
				Null: &sysl.Value_Null{},
			},
		}
	}
	return val
}

func evalIfelse(ee *exprEval, assign Scope, x *sysl.Expr_Ifelse) *sysl.Value {
	cond := Eval(ee, assign, x.Ifelse.Cond)
	if cond.GetB() {
		return Eval(ee, assign, x.Ifelse.IfTrue)
	}
	return Eval(ee, assign, x.Ifelse.IfFalse)
}

func evalSet(ee *exprEval, assign Scope, x *sysl.Expr_Set) *sysl.Value {
	{
		setResult := MakeValueSet()
		for _, s := range x.Set.Expr {
			AppendItemToValueList(setResult.GetSet(), Eval(ee, assign, s))
		}
		return setResult
	}
}

func evalList(ee *exprEval, assign Scope, x *sysl.Expr_List_) *sysl.Value {
	{
		listResult := MakeValueList()
		for _, s := range x.List.Expr {
			AppendItemToValueList(listResult.GetList(), Eval(ee, assign, s))
		}
		return listResult
	}
}

func EvaluateApp(app *sysl.Application, view *sysl.View, s Scope) *sysl.Value {
	ee := exprEval{
		txApp:     app,
		exprStack: exprStack{},
		logger:    logrus.StandardLogger(),
	}
	val, err := ee.eval(s, view.Expr)
	if err != nil {
		logrus.Panic(err.Error())
	}
	return val
}

// EvaluateView evaluate the view using the Scope
func EvaluateView(mod *sysl.Module, appName, viewName string, s Scope) *sysl.Value {
	txApp := mod.Apps[appName]
	view := txApp.Views[viewName]
	if view.Expr.Type == nil {
		view.Expr.Type = view.RetType
	}

	return EvaluateApp(txApp, view, s)
}

type exprEval struct {
	txApp     *sysl.Application
	exprStack exprStack
	logger    *logrus.Logger
}

func logentry(logger *logrus.Logger, expr *sysl.Expr) *logrus.Entry {
	if expr.SourceContext == nil {
		return logger.WithFields(logrus.Fields{})
	}
	return logger.WithFields(logrus.Fields{
		"filename": expr.SourceContext.File,
		"line":     expr.SourceContext.Start.Line,
		"col":      expr.SourceContext.Start.Col,
	})
}

func (ee *exprEval) LogEntry() *logrus.Entry {
	return logentry(ee.logger, ee.exprStack.Peek().e)
}

func (ee *exprEval) handlePanic() {
	if r := recover(); r != nil {
		ee.logger.Errorf("Evaluation Failed: ")
		stack := ee.exprStack.s
		for i := len(stack) - 1; i >= 0; i-- {
			logentry(ee.logger, stack[i].e).Errorf("... %s", getExprText(stack[i].e))
		}
		os.Exit(1)
	}
}

func (ee *exprEval) eval(scope Scope, expr *sysl.Expr) (val *sysl.Value, err error) {
	entry := logentry(ee.logger, expr)

	entry.Tracef("Entering: %s", getExprText(expr))

	ee.exprStack.Push(scope, expr)
	defer ee.exprStack.Pop()
	defer ee.handlePanic()

	switch e := expr.Expr.(type) {
	case *sysl.Expr_Transform_:
		val = ee.evalTransform(scope, e, expr)
	case *sysl.Expr_Binexpr:
		val = evalBinExpr(ee, scope, e.Binexpr)
	case *sysl.Expr_Call_:
		val = evalCall(ee, scope, e)
	case *sysl.Expr_Name:
		val = evalName(ee, scope, e)
	case *sysl.Expr_GetAttr_:
		val = evalGetAttr(ee, scope, e)
	case *sysl.Expr_Ifelse:
		val = evalIfelse(ee, scope, e)
	case *sysl.Expr_Literal:
		val = e.Literal
	case *sysl.Expr_Set:
		val = evalSet(ee, scope, e)
	case *sysl.Expr_List_:
		val = evalList(ee, scope, e)
	case *sysl.Expr_Unexpr:
		val, err = ee.eval(scope, e.Unexpr.Arg)
		if err != nil {
			return nil, err
		}
		val = evalUnaryFunc(e.Unexpr.Op, val)
	default:
		return nil, fmt.Errorf("unhandled sysl.Expr type '%s' @ %s:%d:%d",
			reflect.TypeOf(expr.Expr).String(),
			expr.SourceContext.File, expr.SourceContext.Start.Line, expr.SourceContext.Start.Col)
	}
	entry.Tracef("Result: %s", val.String())
	return val, err
}
