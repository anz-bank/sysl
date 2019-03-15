package main

import (
	"strings"
	"testing"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.WarnLevel)
}

func TestEvalStrategySetup(t *testing.T) {
	for key := range valueFunctions {
		idx := strings.Index(key, "_VALUE_")
		op := key[:idx]
		_, has := functionEvalStrategy[sysl.Expr_BinExpr_Op(sysl.Expr_BinExpr_Op_value[op])]
		assert.Truef(t, has, "Op %s exists in functionEvalStrategy", op[1])
	}

	for key := range exprFunctions {
		idx := strings.Index(key, "_VALUE_")
		op := key[:idx]
		_, has := functionEvalStrategy[sysl.Expr_BinExpr_Op(sysl.Expr_BinExpr_Op_value[op])]
		assert.Truef(t, has, "Op %s exists in functionEvalStrategy", op[1])
	}
}

func TestScopeAddApp(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	app := s["app"].GetMap().Items
	assert.Equal(t, appName, app["name"].GetS(), "unexpected app name")
	types := app["types"].GetMap().Items
	assert.Equal(t, 2, len(types), "unexpected types count")
	typeRequest := types["Request"].GetMap().Items
	assert.Equal(t, 6, len(typeRequest), "unexpected type attribute count")
	assert.Equal(t, "tuple", typeRequest["type"].GetS(), "unexpected typename")
	fields := typeRequest["fields"].GetMap().Items
	assert.Equal(t, 2, len(fields), "unexpected field count")
	idField := fields["id"].GetMap().Items
	assert.Equal(t, 6, len(idField), "unexpected id Field count")
	assert.Equal(t, "primitive", idField["type"].GetS(), "unexpected id field type")
	assert.Equal(t, "INT", idField["primitive"].GetS(), "unexpected id field type name")

	union := app["union"].GetMap().Items
	unionMessage := union["Message"].GetMap().Items
	assert.Equal(t, "union", unionMessage["type"].GetS(), "unexpected id Field count")
	assert.Equal(t, 2, len(unionMessage["fields"].GetSet().Value), "unexpected id Field count")
	assert.Equal(t, "Request", unionMessage["fields"].GetSet().Value[0].GetS(), "unexpected id Field count")
	assert.Equal(t, "Response", unionMessage["fields"].GetSet().Value[1].GetS(), "unexpected id Field count")

	alias := app["alias"].GetMap().Items
	assert.Equal(t, 4, len(alias))
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
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.NotNil(t, mod, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "math"
	assert.NotNil(t, txApp.Views[viewName], "View not loaded")
	assert.Equal(t, 2, len(txApp.Views[viewName].Param), "Params not correct")
	s := Scope{}
	s.AddInt("lhs", 1)
	s.AddInt("rhs", 1)
	out := Eval(txApp, s, txApp.Views[viewName].Expr).GetMap().Items
	assert.Equal(t, int64(2), out["out1"].GetI(), "unexpected value")

	assert.Equal(t, int64(6), out["out2"].GetMap().Items["out3"].GetI(), "unexpected value")

	assert.Equal(t, int64(0), out["out3"].GetI(), "unexpected value")
	assert.Equal(t, int64(1), out["out4"].GetI(), "unexpected value")
	assert.Equal(t, int64(1), out["out5"].GetI(), "unexpected value")
	assert.Equal(t, int64(0), out["out6"].GetI(), "unexpected value")
}

func TestEvalCompare(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.NotNil(t, mod, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "compare"
	assert.NotNil(t, txApp.Views[viewName], "View not loaded")
	s := Scope{}
	s.AddInt("lhs", 1)
	s.AddInt("rhs", 1)
	out := Eval(txApp, s, txApp.Views[viewName].Expr).GetMap().Items
	assert.Equal(t, 6, len(out))
	assert.True(t, out["eq"].GetB())
	assert.False(t, out["gt"].GetB())
	assert.False(t, out["lt"].GetB())
	assert.True(t, out["ge"].GetB())
	assert.True(t, out["le"].GetB())
	assert.False(t, out["ne"].GetB())
}

func TestEvalListSetOps(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.NotNil(t, mod, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "ListSetOps"

	assert.NotNil(t, txApp.Views[viewName], "View not loaded")
	assert.Equal(t, 1, len(txApp.Views[viewName].Param), "Params not correct")
	s := Scope{}
	s.AddInt("lhs", 1)
	out := Eval(txApp, s, txApp.Views[viewName].Expr)
	strs := out.GetMap().Items["strs"].GetSet()
	assert.NotNil(t, strs)
	assert.Equal(t, 2, len(strs.Value))
	assert.Equal(t, "lhs", strs.Value[0].GetS())
	assert.Equal(t, "rhs", strs.Value[1].GetS())

	assert.Equal(t, int64(2), out.GetMap().Items["count1"].GetI())
	assert.Equal(t, int64(2), out.GetMap().Items["count2"].GetI())
	assert.Equal(t, 3, len(out.GetMap().Items["list"].GetList().Value))

	numbers := out.GetMap().Items["numbers"].GetSet()
	assert.NotNil(t, numbers)
	assert.Equal(t, 2, len(numbers.Value))
	numbers2 := out.GetMap().Items["numbers2"].GetSet()
	assert.NotNil(t, numbers2)
	assert.Equal(t, 0, len(numbers2.Value))
}

func TestEvalIsKeyword(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.NotNil(t, mod, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "IsKeyword"

	assert.NotNil(t, txApp.Views[viewName], "View not loaded")
	assert.Equal(t, 1, len(txApp.Views[viewName].Param), "Params not correct")
	s := Scope{}
	s.AddString("word", "defer")
	out := Eval(txApp, s, txApp.Views[viewName].Expr)
	assert.True(t, out.GetMap().Items["out"].GetB(), "unexpected value")
}

func TestEvalIfElseAlt(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.NotNil(t, mod, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "JavaType"

	assert.NotNil(t, txApp.Views[viewName], "View not loaded")
	assert.Equal(t, 1, len(txApp.Views[viewName].Param), "Params not correct")
	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	s["t"] = s["app"].GetMap().Items["types"].GetMap().Items["Request"].GetMap().Items["fields"].GetMap().Items["payload"]
	out := Eval(txApp, s, txApp.Views[viewName].Expr)
	assert.Equal(t, "String", out.GetMap().Items["out"].GetS(), "unexpected value")

	s["t"] = s["app"].GetMap().Items["types"].GetMap().Items["Response"].GetMap().Items["fields"].GetMap().Items["names"]
	out = Eval(txApp, s, txApp.Views[viewName].Expr)
	assert.Equal(t, "List<Request>", out.GetMap().Items["out"].GetS(), "unexpected value")
}

func TestEvalGetAppAttributes(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "GetAppAttributes", s)
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
	assert.Equal(t, 2, len(defSet), "definition length is incorrect")

	requestBody := defSet[0].GetMap().Items
	assert.Equal(t, "RequestImpl", requestBody["className"].GetS())

	// check members of Request
	requestClassBody := requestBody["classBody"].GetSet().Value
	assert.Equal(t, 4, len(requestClassBody))

	getRequestId := requestClassBody[0].GetMap().Items
	assert.Equal(t, 3, len(getRequestId))
	assert.Equal(t, "public", getRequestId["access"].GetS())
	assert.Equal(t, "*int", getRequestId["returnType"].GetS())
	assert.Equal(t, "getid", getRequestId["methodName"].GetS())

	getRequestPayload := requestClassBody[1].GetMap().Items
	assert.Equal(t, 3, len(getRequestPayload))
	assert.Equal(t, "String", getRequestPayload["returnType"].GetS())
	assert.Equal(t, "getpayload", getRequestPayload["methodName"].GetS())

	setRequestId := requestClassBody[2].GetMap().Items
	assert.Equal(t, 2, len(setRequestId))
	assert.Equal(t, "setid", setRequestId["methodName"].GetS())

	setRequestPayload := requestClassBody[3].GetMap().Items
	assert.Equal(t, 2, len(setRequestPayload))
	assert.Equal(t, "setpayload", setRequestPayload["methodName"].GetS())

	// check members of Response
	responseBody := defSet[1].GetMap().Items
	assert.Equal(t, "ResponseImpl", responseBody["className"].GetS(), "Response unexpected value")

}

func TestEvalNullCheckAppAttrs(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "NullCheckAppAttrs", s)

	assert.False(t, out.GetMap().Items["NotHasAttrName"].GetB())
	assert.True(t, out.GetMap().Items["NotHasAttrFoo"].GetB())
	assert.True(t, out.GetMap().Items["hasAttrName"].GetB())
	assert.False(t, out.GetMap().Items["hasAttrFoo"].GetB())
}

func TestScopeAddRestApp(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	s := Scope{}
	appName := "TodoApp"
	s.AddApp("app", mod.Apps[appName])
	app := s["app"].GetMap().Items
	assert.Equal(t, appName, app["name"].GetS(), "unexpected app name")
	endpoints := app["endpoints"].GetMap().Items
	assert.Equal(t, 4, len(endpoints), "unexpected endpoint count")
	rootTodos := endpoints["GET /todos"].GetMap().Items
	assert.Equal(t, "GET /todos", rootTodos["name"].GetS(), "unexpected endpoint name")
	assert.Equal(t, "GET", rootTodos["method"].GetS())
	assert.Equal(t, "todos", rootTodos["ret"].GetMap().Items["payload"].GetS())
	assert.Equal(t, 0, len(rootTodos["pathvars"].GetList().Value))
	assert.Equal(t, "/todos", rootTodos["path"].GetS())
	assert.Equal(t, true, rootTodos["is_rest"].GetB())
	assert.Equal(t, false, rootTodos["is_pubsub"].GetB())
	assert.Equal(t, "rest", rootTodos["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

	postTodo := endpoints["POST /todos"].GetMap().Items
	assert.Equal(t, "POST /todos", postTodo["name"].GetS(), "unexpected endpoint name")
	paramList := postTodo["params"].GetList().Value
	assert.Equal(t, 2, len(paramList))
	paramItem0 := paramList[0].GetMap().Items
	assert.Equal(t, "newTodo", paramItem0["name"].GetS())
	assert.Equal(t, "body", paramItem0["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

	paramItem1 := paramList[1].GetMap().Items
	assert.Equal(t, "accept", paramItem1["name"].GetS())
	assert.Equal(t, "header", paramItem1["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

	todosById := endpoints["GET /todos/{id}"].GetMap().Items
	assert.Equal(t, "GET /todos/{id}", todosById["name"].GetS(), "unexpected endpoint name")
	assert.Equal(t, "GET", todosById["method"].GetS())
	assert.Equal(t, "todo", todosById["ret"].GetMap().Items["payload"].GetS())
	assert.Equal(t, 1, len(todosById["pathvars"].GetList().Value))
	assert.Equal(t, "/todos/{id}", todosById["path"].GetS())
	assert.Equal(t, true, todosById["is_rest"].GetB())
	assert.Equal(t, false, todosById["is_pubsub"].GetB())
	assert.Equal(t, "rest", todosById["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

	todosByIdStatus := endpoints["GET /todos/{id}/{status}"].GetMap().Items
	assert.Equal(t, "GET /todos/{id}/{status}", todosByIdStatus["name"].GetS(), "unexpected endpoint name")
	assert.Equal(t, "GET", todosByIdStatus["method"].GetS())
	assert.Equal(t, "todoWithStatus", todosByIdStatus["ret"].GetMap().Items["payload"].GetS())
	assert.Equal(t, 2, len(todosByIdStatus["pathvars"].GetList().Value))
	assert.Equal(t, "/todos/{id}/{status}", todosByIdStatus["path"].GetS())
	assert.Equal(t, true, todosByIdStatus["is_rest"].GetB())
	assert.Equal(t, false, todosByIdStatus["is_pubsub"].GetB())
	assert.Equal(t, "rest", todosByIdStatus["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS())

}

func TestEvalStringOps(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "TodoApp"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "StringOps", s)
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
	assert.Equal(t, 2, len(items["Fields"].GetList().Value))
	assert.True(t, items["HasPrefix"].GetB())
	assert.True(t, items["HasSuffix"].GetB())
	assert.Equal(t, "Hello_World", items["Join"].GetS())
	assert.Equal(t, int64(12), items["LastIndex"].GetI())
	assert.Equal(t, "Hello_World", items["Replace"].GetS())
	assert.Equal(t, 3, len(items["Split"].GetList().Value))
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
	assert.Equal(t, 2, len(items["FindAllString"].GetList().Value))
	assert.Equal(t, 2, len(items["tabs"].GetList().Value))

}

func TestIncorrectArgsToGoFunc(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "TodoApp"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "IncorrectArgsToGoFunc", s)
	assert.NotNil(t, out.GetMap())
	items := out.GetMap().Items
	contains, has := items["Contains"]
	assert.True(t, has)
	assert.Nil(t, contains)
	wrongNumberOfArgs, has_Args := items["WrongNumberOfArgs"]
	assert.True(t, has_Args)
	assert.Nil(t, wrongNumberOfArgs)
}

func TestEvalFlatten(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "TodoApp"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "Flatten", s)
	assert.NotNil(t, out.GetMap().Items["names"].GetSet())
	l := out.GetMap().Items["names"].GetSet().Value
	assert.Equal(t, 4, len(l))
	assert.Equal(t, "GET", l[0].GetS())
	assert.Equal(t, "GETid", l[1].GetS())
	assert.Equal(t, "GETidstatus", l[2].GetS())
	assert.Equal(t, "POST", l[3].GetS())

	numbers1 := out.GetMap().Items["listOfNumbers1"].GetList().Value
	assert.Equal(t, 6, len(numbers1))
	assert.Equal(t, int64(3), numbers1[0].GetI())

	numbers2 := out.GetMap().Items["listOfNumbers2"].GetList().Value
	assert.Equal(t, 6, len(numbers2))
	assert.Equal(t, int64(2), numbers2[0].GetI())

	numbers3 := out.GetMap().Items["setOfNumbers1"].GetSet().Value
	assert.Equal(t, 6, len(numbers3))
	assert.Equal(t, int64(2), numbers3[0].GetI())
}

func TestEvalWhere(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "Where", s)

	numbers1 := out.GetMap().Items["greaterThanOne"].GetSet().Value
	assert.Equal(t, 2, len(numbers1))

	RequestFromList := out.GetMap().Items["RequestFromList"].GetList().Value
	assert.Equal(t, 1, len(RequestFromList))

	strOne := out.GetMap().Items["strOne"].GetSet().Value
	assert.Equal(t, 1, len(strOne))

	request := out.GetMap().Items["Request"].GetSet().Value
	assert.Equal(t, 1, len(request))

	listofNames := out.GetMap().Items["ListofNames"].GetList().Value
	assert.Equal(t, 2, len(listofNames))

	NotObjectAliases := out.GetMap().Items["NotObjectAliases"].GetMap().Items
	assert.Equal(t, 3, len(NotObjectAliases))
}

func TestEvalLinks(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "Links", s)
	assert.NotNil(t, out.GetMap().Items["links"].GetSet())
	l := out.GetMap().Items["links"].GetSet().Value
	assert.Equal(t, 5, len(l))
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
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "TestDotScope", s).GetMap().Items
	assert.Equal(t, 3, len(out))
}

func TestListOfTypeNames(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "ListOfTypeNames", s)
	l := out.GetList()
	assert.NotNil(t, l)
	assert.Equal(t, 2, len(l.Value))
}
