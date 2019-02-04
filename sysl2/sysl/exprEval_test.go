package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScopeAddApp(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	s := make(Scope)
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	app := s["app"].GetMap().Items
	assert.True(t, app["name"].GetS() == "Model", "unexpected app name")
	types := app["types"].GetMap().Items
	assert.True(t, len(types) == 2, "unexpected types count")
	typeRequest := types["Request"].GetMap().Items
	assert.True(t, len(typeRequest) == 4, "unexpected type attribute count")
	assert.True(t, typeRequest["type"].GetS() == "tuple", "unexpected typename")
	fields := typeRequest["fields"].GetMap().Items
	assert.True(t, len(fields) == 2, "unexpected field count")
	idField := fields["id"].GetMap().Items
	assert.True(t, len(idField) == 5, "unexpected id Field count")
	assert.True(t, idField["type"].GetS() == "primitive", "unexpected id field type")
	assert.True(t, idField["primitive"].GetS() == "INT", "unexpected id field type name")

	union := app["union"].GetMap().Items
	unionMessage := union["Message"].GetMap().Items
	assert.True(t, unionMessage["type"].GetS() == "union", "unexpected id Field count")
	assert.True(t, len(unionMessage["fields"].GetSet().Value) == 2, "unexpected id Field count")
	assert.True(t, unionMessage["fields"].GetSet().Value[0].GetS() == "Request", "unexpected id Field count")
	assert.True(t, unionMessage["fields"].GetSet().Value[1].GetS() == "Response", "unexpected id Field count")
}

func TestEvalIntegerAdd(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.True(t, mod != nil, "Module not loaded")
	txApp := mod.Apps["TransformApp"]

	assert.True(t, txApp.Views["add"] != nil, "View not loaded")
	assert.True(t, len(txApp.Views["add"].Param) == 2, "Params not correct")
	s := make(Scope)
	s.AddInt("lhs", 1)
	s.AddInt("rhs", 1)
	out := Eval(txApp, &s, txApp.Views["add"].Expr)
	assert.True(t, out.GetMap().Items["out1"].GetI() == 2, "unexpected value")
	assert.True(t, out.GetMap().Items["out2"].GetMap().Items["out3"].GetI() == 6, "unexpected value")
}

func TestEvalAddSet(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.True(t, mod != nil, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "addSet"

	assert.True(t, txApp.Views[viewName] != nil, "View not loaded")
	assert.True(t, len(txApp.Views[viewName].Param) == 1, "Params not correct")
	s := make(Scope)
	s.AddInt("lhs", 1)
	out := Eval(txApp, &s, txApp.Views[viewName].Expr)
	assert.True(t, len(out.GetMap().Items["out"].GetSet().Value) == 2, "unexpected value")
}

func TestEvalIsKeyword(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.True(t, mod != nil, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "IsKeyword"

	assert.True(t, txApp.Views[viewName] != nil, "View not loaded")
	assert.True(t, len(txApp.Views[viewName].Param) == 1, "Params not correct")
	s := make(Scope)
	s.AddString("word", "defer")
	out := Eval(txApp, &s, txApp.Views[viewName].Expr)
	assert.True(t, out.GetMap().Items["out"].GetB(), "unexpected value")
}

func TestEvalIfElseAlt(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.True(t, mod != nil, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "JavaType"

	assert.True(t, txApp.Views[viewName] != nil, "View not loaded")
	assert.True(t, len(txApp.Views[viewName].Param) == 1, "Params not correct")
	s := make(Scope)
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	s["t"] = s["app"].GetMap().Items["types"].GetMap().Items["Request"].GetMap().Items["fields"].GetMap().Items["payload"]
	out := Eval(txApp, &s, txApp.Views[viewName].Expr)
	assert.True(t, out.GetMap().Items["out"].GetS() == "String", "unexpected value")

	s["t"] = s["app"].GetMap().Items["types"].GetMap().Items["Response"].GetMap().Items["fields"].GetMap().Items["names"]
	out = Eval(txApp, &s, txApp.Views[viewName].Expr)
	assert.True(t, out.GetMap().Items["out"].GetS() == "List<String>", "unexpected value")
}

func TestEvalGetAppAttributes(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")

	s := make(Scope)
	appName := "Model"
	s.AddApp("app", mod.Apps[appName])
	out := EvalView(mod, "TransformApp", "GetAppAttributes", &s)
	assert.True(t, out.GetMap().Items["out"].GetS() == "com.example.gen", "unexpected value")

	m := out.GetMap().Items["package"]
	assert.True(t, m.GetMap().Items["packageName"].GetS() == "com.example.gen", "packageName unexpected value")

	m = out.GetMap().Items["import"]
	assert.True(t, m.GetList().Value[0].GetMap().Items["importPath"].GetS() == "package1", "import unexpected value")
	assert.True(t, m.GetList().Value[1].GetMap().Items["importPath"].GetS() == "package2", "import unexpected value")

	m = out.GetMap().Items["definition"]
	assert.True(t, len(m.GetSet().Value) == 2, "definition length is incorrect")
	assert.True(t, m.GetSet().Value[0].GetMap().Items["className"].GetS() == "RequestImpl", "Request unexpected value")
	assert.True(t, m.GetSet().Value[1].GetMap().Items["className"].GetS() == "ResponseImpl", "Response unexpected value")

	classBody := m.GetSet().Value[0].GetMap().Items["classBody"]
	assert.True(t, len(classBody.GetSet().Value) == 2, "classBody unexpected value")

	requestId := classBody.GetSet().Value[0].GetMap().Items
	assert.True(t, len(requestId) == 3, "requestId unexpected count")
	assert.True(t, requestId["access"].GetS() == "public", "access unexpected typename")
	assert.True(t, requestId["returnType"].GetS() == "*int", "returnType unexpected typename")
	assert.True(t, requestId["methodName"].GetS() == "getid", "returnType unexpected typename")

	responseId := classBody.GetSet().Value[1].GetMap().Items
	assert.True(t, len(responseId) == 3, "requestId unexpected count")

}
