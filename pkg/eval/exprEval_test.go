package eval

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/anz-bank/sysl/pkg/parse"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "../../tests/"

const (
	modelAppName = "Model"
	todoAppName  = "TodoApp"
)

func newExprEval(txApp *sysl.Application) *exprEval {
	log, _ := test.NewNullLogger()
	return &exprEval{
		txApp:     txApp,
		exprStack: exprStack{},
		logger:    log,
	}
}

func TestEvalStrategySetup(t *testing.T) {
	t.Parallel()

	for key := range valueFunctions {
		idx := strings.Index(key, "_Value")
		op := key[:idx]
		_, has := functionEvalStrategy[sysl.Expr_BinExpr_Op(sysl.Expr_BinExpr_Op_value[op])]
		assert.True(t, has, "%#v", key)
	}

	for key := range exprFunctions {
		idx := strings.Index(key, "_Value")
		op := key[:idx]
		_, has := functionEvalStrategy[sysl.Expr_BinExpr_Op(sysl.Expr_BinExpr_Op_value[op])]
		assert.True(t, has, "%#v", key)
	}
}

func TestScopeAddApp(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)
	s := Scope{}
	s.AddApp("app", mod.Apps[modelAppName])
	app := s["app"].GetMap().Items
	assert.Equal(t, modelAppName, app["name"].GetS())
	types := app["types"].GetMap().Items
	assert.Len(t, types, 2)
	typeRequest := types["Request"].GetMap().Items
	assert.Len(t, typeRequest, 6)
	assert.Equal(t, "tuple", typeRequest["type"].GetS())
	fields := typeRequest["fields"].GetMap().Items
	assert.Len(t, fields, 2)
	idField := fields["id"].GetMap().Items
	assert.Len(t, idField, 6)
	assert.Equal(t, "primitive", idField["type"].GetS())
	assert.Equal(t, "INT", idField["primitive"].GetS())

	union := app["union"].GetMap().Items
	unionMessage := union["Message"].GetMap().Items
	assert.Equal(t, "union", unionMessage["type"].GetS())
	assert.Len(t, unionMessage["fields"].GetSet().Value, 2)
	assert.Equal(t, "Request", unionMessage["fields"].GetSet().Value[0].GetS())
	assert.Equal(t, "Response", unionMessage["fields"].GetSet().Value[1].GetS())

	alias := app["alias"].GetMap().Items
	assert.Len(t, alias, 4)
	aliasError := alias["Error"].GetMap().Items
	assert.Equal(t, "primitive", aliasError["type"].GetS())
	assert.Equal(t, "STRING", aliasError["primitive"].GetS())

	aliasObject := alias["Object"].GetMap().Items
	assert.Equal(t, "type_ref", aliasObject["type"].GetS())
	assert.Equal(t, "Ignored", aliasObject["type_ref"].GetS())

	aliasTerms := alias["Terms"].GetMap().Items
	assert.Equal(t, "sequence", aliasTerms["type"].GetS())
	aliasSeqType := aliasTerms["sequence"].GetMap().Items
	assert.Equal(t, "Term", aliasSeqType["type_ref"].GetS())

	aliasAccounts := alias["Accounts"].GetMap().Items
	assert.Equal(t, "set", aliasAccounts["type"].GetS())
	assert.Equal(t, "Term", aliasAccounts["set"].GetS())
}

func TestEvalIntegerMath(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)
	txApp := mod.Apps["TransformApp"]
	viewName := "math"
	assert.NotNil(t, txApp.Views[viewName])
	assert.Len(t, txApp.Views[viewName].Param, 2)
	s := Scope{}
	s.AddInt("lhs", 1)
	s.AddInt("rhs", 1)
	out := Eval(newExprEval(txApp), s, txApp.Views[viewName].Expr).GetMap().Items
	assert.Equal(t, int64(2), out["out1"].GetI())

	assert.Equal(t, int64(6), out["out2"].GetMap().Items["out3"].GetI())

	assert.Equal(t, int64(0), out["out3"].GetI())
	assert.Equal(t, int64(1), out["out4"].GetI())
	assert.Equal(t, int64(1), out["out5"].GetI())
	assert.Equal(t, int64(0), out["out6"].GetI())
}

func TestEvalCompare(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)
	txApp := mod.Apps["TransformApp"]
	viewName := "compare"
	assert.NotNil(t, txApp.Views[viewName])
	s := Scope{}
	s.AddInt("lhs", 1)
	s.AddInt("rhs", 1)
	out := Eval(newExprEval(txApp), s, txApp.Views[viewName].Expr).GetMap().Items
	assert.Len(t, out, 6)
	assert.True(t, out["eq"].GetB())
	assert.False(t, out["gt"].GetB())
	assert.False(t, out["lt"].GetB())
	assert.True(t, out["ge"].GetB())
	assert.True(t, out["le"].GetB())
	assert.False(t, out["ne"].GetB())
}

func TestEvalListSetOps(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)
	txApp := mod.Apps["TransformApp"]
	viewName := "ListSetOps"

	assert.NotNil(t, txApp.Views[viewName])
	assert.Len(t, txApp.Views[viewName].Param, 1)
	s := Scope{}
	s.AddInt("lhs", 1)
	out := Eval(newExprEval(txApp), s, txApp.Views[viewName].Expr)
	strs := out.GetMap().Items["strs"].GetSet()
	assert.NotNil(t, strs)
	assert.Len(t, strs.Value, 2)
	assert.Equal(t, "lhs", strs.Value[0].GetS())
	assert.Equal(t, "rhs", strs.Value[1].GetS())

	assert.Equal(t, int64(2), out.GetMap().Items["count1"].GetI())
	assert.Equal(t, int64(2), out.GetMap().Items["count2"].GetI())
	assert.Equal(t, int64(3), out.GetMap().Items["count3"].GetI())
	assert.Len(t, out.GetMap().Items["list"].GetList().Value, 3)

	numbers := out.GetMap().Items["numbers"].GetSet()
	assert.NotNil(t, numbers)
	assert.Len(t, numbers.Value, 2)
	numbers2 := out.GetMap().Items["numbers2"].GetSet()
	assert.NotNil(t, numbers2)
	assert.Len(t, numbers2.Value, 0)
}

func TestEvalIsKeyword(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)
	txApp := mod.Apps["TransformApp"]
	viewName := "IsKeyword"

	assert.NotNil(t, txApp.Views[viewName])
	assert.Len(t, txApp.Views[viewName].Param, 1)
	s := Scope{}
	s.AddString("word", "defer")
	out := Eval(newExprEval(txApp), s, txApp.Views[viewName].Expr)
	assert.True(t, out.GetMap().Items["out"].GetB())
}

func TestEvalIfElseAlt(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)
	txApp := mod.Apps["TransformApp"]
	viewName := "JavaType"

	assert.NotNil(t, txApp.Views[viewName])
	assert.Len(t, txApp.Views[viewName].Param, 1)
	s := Scope{}
	s.AddApp("app", mod.Apps[modelAppName])
	s["t"] = s["app"].GetMap().Items["types"].GetMap().Items["Request"].GetMap().Items["fields"].GetMap().Items["payload"]
	out := Eval(newExprEval(txApp), s, txApp.Views[viewName].Expr)
	assert.Equal(t, "String", out.GetMap().Items["out"].GetS())

	s["t"] = s["app"].GetMap().Items["types"].GetMap().Items["Response"].GetMap().Items["fields"].GetMap().Items["names"]
	out = Eval(newExprEval(txApp), s, txApp.Views[viewName].Expr)
	assert.Equal(t, "List<Request>", out.GetMap().Items["out"].GetS())
}

func TestEvalGetAppAttributes(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[modelAppName])
	out := EvaluateView(mod, "TransformApp", "GetAppAttributes", s)
	assert.Equal(t, "com.example.gen", out.GetMap().Items["out"].GetS())
	assert.Nil(t, out.GetMap().Items["Nil"])
	assert.False(t, out.GetMap().Items["stringInNull"].GetB())
	assert.False(t, out.GetMap().Items["stringInList"].GetB())

	packageMap := out.GetMap().Items["package"].GetMap().Items
	assert.Equal(t, "com.example.gen", packageMap["packageName"].GetS())

	importSet := out.GetMap().Items["import"].GetSet().Value
	assert.Equal(t, "Package1", importSet[0].GetMap().Items["importPath"].GetS())
	assert.Equal(t, "Package2", importSet[1].GetMap().Items["importPath"].GetS())

	defSet := out.GetMap().Items["definition"].GetSet().Value
	assert.Len(t, defSet, 2)

	requestBody := defSet[0].GetMap().Items
	assert.Equal(t, "RequestImpl", requestBody["className"].GetS())

	// check members of Request
	requestClassBody := requestBody["classBody"].GetSet().Value
	assert.Len(t, requestClassBody, 4)

	getRequestID := requestClassBody[0].GetMap().Items
	assert.Len(t, getRequestID, 3)
	assert.Equal(t, "public", getRequestID["access"].GetS())
	assert.Equal(t, "*int", getRequestID["returnType"].GetS())
	assert.Equal(t, "getid", getRequestID["methodName"].GetS())

	getRequestPayload := requestClassBody[1].GetMap().Items
	assert.Len(t, getRequestPayload, 3)
	assert.Equal(t, "String", getRequestPayload["returnType"].GetS())
	assert.Equal(t, "getpayload", getRequestPayload["methodName"].GetS())

	setRequestID := requestClassBody[2].GetMap().Items
	assert.Len(t, setRequestID, 2)
	assert.Equal(t, "setid", setRequestID["methodName"].GetS())

	setRequestPayload := requestClassBody[3].GetMap().Items
	assert.Len(t, setRequestPayload, 2)
	assert.Equal(t, "setpayload", setRequestPayload["methodName"].GetS())

	// check members of Response
	responseBody := defSet[1].GetMap().Items
	assert.Equal(t, "ResponseImpl", responseBody["className"].GetS())
}

func TestEvalNullCheckAppAttrs(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[modelAppName])
	out := EvaluateView(mod, "TransformApp", "NullCheckAppAttrs", s)

	assert.False(t, out.GetMap().Items["NotHasAttrName"].GetB())
	assert.True(t, out.GetMap().Items["NotHasAttrFoo"].GetB())
	assert.True(t, out.GetMap().Items["hasAttrName"].GetB())
	assert.False(t, out.GetMap().Items["hasAttrFoo"].GetB())
}

func TestScopeAddRestApp(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)
	s := Scope{}
	s.AddApp("app", mod.Apps[todoAppName])
	app := s["app"].GetMap().Items
	assert.Equal(t, todoAppName, app["name"].GetS())
	endpoints := app["endpoints"].GetMap().Items
	assert.Len(t, endpoints, 4)
	rootTodos := endpoints["GET /todos"].GetMap().Items
	assert.Equal(t, "GET /todos", rootTodos["name"].GetS())
	assert.Equal(t, "GET", rootTodos["method"].GetS())
	assert.Equal(t, "NotFoundError", rootTodos["ret"].GetMap().Items["404"].GetS())
	assert.Equal(t, "ServerError", rootTodos["ret"].GetMap().Items["500"].GetS())
	assert.Equal(t, "todos", rootTodos["ret"].GetMap().Items["200"].GetS())
	assert.Len(t, rootTodos["pathvars"].GetList().Value, 0)
	assert.Equal(t, "/todos", rootTodos["path"].GetS())
	assert.True(t, rootTodos["is_rest"].GetB())
	assert.False(t, rootTodos["is_pubsub"].GetB())
	assert.Equal(t, "rest", rootTodos["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

	postTodo := endpoints["POST /todos"].GetMap().Items
	assert.Equal(t, "POST /todos", postTodo["name"].GetS())
	paramList := postTodo["params"].GetList().Value
	assert.Len(t, paramList, 2)
	paramItem0 := paramList[0].GetMap().Items
	assert.Equal(t, "newTodo", paramItem0["name"].GetS())
	assert.Equal(t, "body", paramItem0["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

	paramItem1 := paramList[1].GetMap().Items
	assert.Equal(t, "accept", paramItem1["name"].GetS())
	assert.Equal(t, "header", paramItem1["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

	todosByID := endpoints["GET /todos/{id}"].GetMap().Items
	assert.Equal(t, "GET /todos/{id}", todosByID["name"].GetS())
	assert.Equal(t, "GET", todosByID["method"].GetS())
	assert.Equal(t, "todo", todosByID["ret"].GetMap().Items["payload"].GetS())
	assert.Len(t, todosByID["pathvars"].GetList().Value, 1)
	assert.Equal(t, "/todos/{id}", todosByID["path"].GetS())
	assert.True(t, todosByID["is_rest"].GetB())
	assert.False(t, todosByID["is_pubsub"].GetB())
	assert.Equal(t, "rest", todosByID["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

	todosByIDStatus := endpoints["GET /todos/{id}/{status}"].GetMap().Items
	assert.Equal(t, "GET /todos/{id}/{status}", todosByIDStatus["name"].GetS())
	assert.Equal(t, "GET", todosByIDStatus["method"].GetS())
	assert.Equal(t, "todoWithStatus", todosByIDStatus["ret"].GetMap().Items["payload"].GetS())
	assert.Len(t, todosByIDStatus["pathvars"].GetList().Value, 2)
	assert.Equal(t, "/todos/{id}/{status}", todosByIDStatus["path"].GetS())
	assert.True(t, todosByIDStatus["is_rest"].GetB())
	assert.False(t, todosByIDStatus["is_pubsub"].GetB())
	assert.Equal(t, "rest", todosByIDStatus["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())
}

func TestEvalStringOps(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[todoAppName])
	out := EvaluateView(mod, "TransformApp", "StringOps", s)
	assert.NotNil(t, out.GetMap())
	items := out.GetMap().Items

	for name := range items {
		assert.NotNilf(t, items[name], "%s", name)
	}

	for name := range GoFuncMap {
		assert.NotNilf(t, items[name], "%s", name)
	}

	assert.True(t, items["Contains"].GetB())
	assert.Equal(t, int64(3), items["Count"].GetI())
	assert.Len(t, items["Fields"].GetList().Value, 2)
	assert.True(t, items["HasPrefix"].GetB())
	assert.True(t, items["HasSuffix"].GetB())
	assert.Equal(t, "Hello_World", items["Join"].GetS())
	assert.Equal(t, int64(12), items["LastIndex"].GetI())
	assert.Equal(t, "Hello_World", items["Replace"].GetS())
	assert.Len(t, items["Split"].GetList().Value, 3)
	assert.Equal(t, "Hello World!", items["Title"].GetS())
	assert.Equal(t, "hello world!", items["ToLower"].GetS())
	assert.Equal(t, "HELLO WORLD!", items["ToTitle"].GetS())
	assert.Equal(t, "HELLO WORLD!", items["ToUpper"].GetS())
	assert.Equal(t, "hello world!", items["Trim"].GetS())
	assert.Equal(t, "hello world! ", items["TrimLeft"].GetS())
	assert.Equal(t, " world! ", items["TrimPrefix"].GetS())
	assert.Equal(t, " hello world!", items["TrimRight"].GetS())
	assert.Equal(t, "hello world!", items["TrimSpace"].GetS())
	assert.Equal(t, " hello ", items["TrimSuffix"].GetS())
	assert.True(t, items["hasHello"].GetB())
	assert.True(t, items["MatchString"].GetB())
	assert.Len(t, items["FindAllString"].GetList().Value, 2)
	assert.Len(t, items["tabs"].GetList().Value, 2)
}

func TestIncorrectArgsToGoFunc(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[todoAppName])
	out := EvaluateView(mod, "TransformApp", "IncorrectArgsToGoFunc", s)
	assert.NotNil(t, out.GetMap())
	items := out.GetMap().Items
	contains, has := items["Contains"]
	assert.True(t, has)
	assert.Nil(t, contains)
	wrongNumberOfArgs, hasArgs := items["WrongNumberOfArgs"]
	assert.True(t, hasArgs)
	assert.Nil(t, wrongNumberOfArgs)
}

func TestEvalFlatten(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[todoAppName])
	out := EvaluateView(mod, "TransformApp", "Flatten", s)
	assert.NotNil(t, out.GetMap().Items["names"].GetSet())
	l := out.GetMap().Items["names"].GetSet().Value
	assert.Len(t, l, 4)
	assert.Equal(t, "GET", l[0].GetS())
	assert.Equal(t, "GETid", l[1].GetS())
	assert.Equal(t, "GETidstatus", l[2].GetS())
	assert.Equal(t, "POST", l[3].GetS())

	numbers1 := out.GetMap().Items["listOfNumbers1"].GetList().Value
	assert.Len(t, numbers1, 6)
	assert.Equal(t, int64(3), numbers1[0].GetI())

	numbers2 := out.GetMap().Items["listOfNumbers2"].GetList().Value
	assert.Len(t, numbers2, 6)
	assert.Equal(t, int64(2), numbers2[0].GetI())

	numbers3 := out.GetMap().Items["setOfNumbers1"].GetSet().Value
	assert.Len(t, numbers3, 6)
	assert.Equal(t, int64(2), numbers3[0].GetI())
}

func TestEvalWhere(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[modelAppName])
	out := EvaluateView(mod, "TransformApp", "Where", s)

	numbers1 := out.GetMap().Items["greaterThanOne"].GetSet().Value
	assert.Len(t, numbers1, 2)

	RequestFromList := out.GetMap().Items["RequestFromList"].GetList().Value
	assert.Len(t, RequestFromList, 1)

	strOne := out.GetMap().Items["strOne"].GetSet().Value
	assert.Len(t, strOne, 1)

	request := out.GetMap().Items["Request"].GetSet().Value
	assert.Len(t, request, 1)

	listofNames := out.GetMap().Items["ListofNames"].GetList().Value
	assert.Len(t, listofNames, 2)

	NotObjectAliases := out.GetMap().Items["NotObjectAliases"].GetMap().Items
	assert.Len(t, NotObjectAliases, 3)
}

func TestEvalLinks(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[modelAppName])
	out := EvaluateView(mod, "TransformApp", "Links", s)
	assert.NotNil(t, out.GetMap().Items["links"].GetSet())
	l := out.GetMap().Items["links"].GetSet().Value
	assert.Len(t, l, 5)
	assert.Equal(t, "id", l[0].GetMap().Items["Left"].GetS())
	assert.Equal(t, "int", l[0].GetMap().Items["Right"].GetS())

	assert.Equal(t, "payload", l[1].GetMap().Items["Left"].GetS())
	assert.Equal(t, "String", l[1].GetMap().Items["Right"].GetS())

	assert.Equal(t, "code", l[2].GetMap().Items["Left"].GetS())
	assert.Equal(t, "int", l[2].GetMap().Items["Right"].GetS())

	assert.Equal(t, "message", l[3].GetMap().Items["Left"].GetS())
	assert.Equal(t, "String", l[3].GetMap().Items["Right"].GetS())

	assert.Equal(t, "names", l[4].GetMap().Items["Left"].GetS())
	assert.Equal(t, "List<Request>", l[4].GetMap().Items["Right"].GetS())
}

func TestDotScope(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[modelAppName])
	out := EvaluateView(mod, "TransformApp", "TestDotScope", s).GetMap().Items
	assert.Len(t, out, 3)
}

func TestListOfTypeNames(t *testing.T) {
	t.Parallel()

	mod, err := parse.NewParser().Parse("eval_expr.sysl", syslutil.NewChrootFs(afero.NewOsFs(), testDir))
	require.NoError(t, err)
	require.NotNil(t, mod)

	s := Scope{}
	s.AddApp("app", mod.Apps[modelAppName])
	out := EvaluateView(mod, "TransformApp", "ListOfTypeNames", s)
	l := out.GetList()
	assert.NotNil(t, l)
	assert.Len(t, l.Value, 2)
}

func TestListAppender_Nil(t *testing.T) {
	result := listAppender(nil, nil)
	require.Nil(t, result[0])
}

func TestSetAppender_Nil(t *testing.T) {
	result := setAppender(nil, nil)
	require.Nil(t, result[0])
}

func TestListAppender_AllowDup(t *testing.T) {
	val := MakeValueString("val")
	list := MakeValueList(MakeValueString("val"))

	result := listAppender(list.GetList().Value, val)

	require.Len(t, result, 2)
	assert.Equal(t, "val", result[0].GetS())
	assert.Equal(t, "val", result[1].GetS())
}

func TestListAppender_BlockDup(t *testing.T) {
	val := MakeValueString("val")
	set := MakeValueSet()
	set.GetSet().Value = append(set.GetSet().Value, MakeValueString("val"))

	result := setAppender(set.GetSet().Value, val)

	require.Len(t, result, 1)
	assert.Equal(t, "val", result[0].GetS())
}

func Test_isInternalMap(t *testing.T) {
	type data struct {
		name     string
		keys     []string
		expected bool
	}
	tests := []data{
		{
			name:     "simple-yes",
			keys:     []string{"key", "value"},
			expected: true,
		},
		{
			name:     "simple-no",
			keys:     []string{"value"},
			expected: false,
		},
		{
			name:     "simple-no2",
			keys:     []string{"key"},
			expected: false,
		},
		{
			name:     "extra-items",
			keys:     []string{"key", "value", "foo"},
			expected: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := MakeValueMap().GetMap()
			for _, x := range tt.keys {
				m.Items[x] = MakeValueString(x)
			}
			assert.Equal(t, tt.expected, isInternalMap(m))
		})
	}
}
