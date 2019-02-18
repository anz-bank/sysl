package main

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.WarnLevel)
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
	assert.Equal(t, 4, len(typeRequest), "unexpected type attribute count")
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
}

func TestEvalIntegerAdd(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.NotNil(t, mod, "Module not loaded")
	txApp := mod.Apps["TransformApp"]

	assert.NotNil(t, txApp.Views["add"], "View not loaded")
	assert.Equal(t, 2, len(txApp.Views["add"].Param), "Params not correct")
	s := Scope{}
	s.AddInt("lhs", 1)
	s.AddInt("rhs", 1)
	out := Eval(txApp, &s, txApp.Views["add"].Expr)
	assert.Equal(t, int64(2), out.GetMap().Items["out1"].GetI(), "unexpected value")
	assert.Equal(t, int64(6), out.GetMap().Items["out2"].GetMap().Items["out3"].GetI(), "unexpected value")
}

func TestEvalUnionSet(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.NotNil(t, mod, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "UnionSet"

	assert.NotNil(t, txApp.Views[viewName], "View not loaded")
	assert.Equal(t, 1, len(txApp.Views[viewName].Param), "Params not correct")
	s := Scope{}
	s.AddInt("lhs", 1)
	out := Eval(txApp, &s, txApp.Views[viewName].Expr)
	assert.Equal(t, 2, len(out.GetMap().Items["out"].GetSet().Value), "unexpected value")
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
	out := Eval(txApp, &s, txApp.Views[viewName].Expr)
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
	out := Eval(txApp, &s, txApp.Views[viewName].Expr)
	assert.Equal(t, "String", out.GetMap().Items["out"].GetS(), "unexpected value")

	s["t"] = s["app"].GetMap().Items["types"].GetMap().Items["Response"].GetMap().Items["fields"].GetMap().Items["names"]
	out = Eval(txApp, &s, txApp.Views[viewName].Expr)
	assert.Equal(t, "List<Request>", out.GetMap().Items["out"].GetS(), "unexpected value")
}

func TestEvalGetAppAttributes(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "GetAppAttributes", &s)
	assert.Equal(t, "com.example.gen", out.GetMap().Items["out"].GetS(), "unexpected value")

	m := out.GetMap().Items["package"]
	assert.Equal(t, "com.example.gen", m.GetMap().Items["packageName"].GetS(), "packageName unexpected value")

	m = out.GetMap().Items["import"]
	assert.Equal(t, "Package1", m.GetList().Value[0].GetMap().Items["importPath"].GetS(), "import unexpected value")
	assert.Equal(t, "Package2", m.GetList().Value[1].GetMap().Items["importPath"].GetS(), "import unexpected value")

	m = out.GetMap().Items["definition"]
	assert.Equal(t, 2, len(m.GetSet().Value), "definition length is incorrect")
	assert.Equal(t, "RequestImpl", m.GetSet().Value[0].GetMap().Items["className"].GetS(), "Request unexpected value")
	assert.Equal(t, "ResponseImpl", m.GetSet().Value[1].GetMap().Items["className"].GetS(), "Response unexpected value")

	classBody := m.GetSet().Value[0].GetMap().Items["classBody"]
	assert.Equal(t, 2, len(classBody.GetSet().Value), "classBody unexpected value")

	requestId := classBody.GetSet().Value[0].GetMap().Items
	assert.Equal(t, 3, len(requestId), "requestId unexpected count")
	assert.Equal(t, "public", requestId["access"].GetS(), "access unexpected typename")
	assert.Equal(t, "*int", requestId["returnType"].GetS(), "returnType unexpected typename")
	assert.Equal(t, "getid", requestId["methodName"].GetS(), "returnType unexpected typename")

	responseId := classBody.GetSet().Value[1].GetMap().Items
	assert.Equal(t, 3, len(responseId), "requestId unexpected count")

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
	root := endpoints["GET /todos"].GetMap().Items
	assert.Equal(t, "GET /todos", root["name"].GetS(), "unexpected endpoint name")
	assert.Equal(t, "GET", root["method"].GetS(), "unexpected endpoint name")
	assert.Equal(t, "/todos", root["path"].GetS(), "unexpected endpoint name")
	assert.Equal(t, true, root["is_rest"].GetB(), "unexpected endpoint kind")
	assert.Equal(t, false, root["is_pubsub"].GetB(), "unexpected is_pubsub value")
	assert.Equal(t, "rest", root["attrs"].GetMap().Items["patterns"].GetList().Value[0].GetS(), "unexpected endpoint attrs")
}

func TestEvalStringOps(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "TodoApp"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "StringOps", &s)
	assert.NotNil(t, out.GetMap())
	items := out.GetMap().Items

	// Check if all functions have been tested
	assert.Equal(t, len(items), len(GoFuncMap))

	for name := range GoFuncMap {
		assert.NotNil(t, items[name])
	}

	assert.True(t, items["Contains"].GetB())
	assert.Equal(t, int64(3), items["Count"].GetI())
	assert.Equal(t, 2, len(items["Fields"].GetList().Value))
	assert.True(t, items["HasPrefix"].GetB())
	assert.True(t, items["HasSuffix"].GetB())
	assert.Equal(t, "Hello_World", items["Join"].GetS())
	assert.Equal(t, int64(12), items["LastIndex"].GetI())
	assert.Equal(t, "Hello_World", items["Replace"].GetS())
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
}

func TestIncorrectArgsToGoFunc(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "TodoApp"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "IncorrectArgsToGoFunc", &s)
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
	out := EvalView(mod, "TransformApp", "FlattenEndpointParams", &s)
	assert.NotNil(t, out.GetMap().Items["names"].GetSet())
	l := out.GetMap().Items["names"].GetSet().Value
	assert.Equal(t, 4, len(l))
	assert.Equal(t, "GET", l[0].GetS())
	assert.Equal(t, "GETid", l[1].GetS())
	assert.Equal(t, "GETidstatus", l[2].GetS())
	assert.Equal(t, "POST", l[3].GetS())
}

func TestEvalLinks(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := Scope{}
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "Links", &s)
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
