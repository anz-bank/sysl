package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, 5, len(idField), "unexpected id Field count")
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

func TestEvalAddSet(t *testing.T) {
	mod, _ := Parse("tests/eval_expr.sysl", "")
	assert.NotNil(t, mod, "Module not loaded")
	txApp := mod.Apps["TransformApp"]
	viewName := "addSet"

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
	assert.Equal(t, "List<String>", out.GetMap().Items["out"].GetS(), "unexpected value")
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
	assert.Equal(t, "package1", m.GetList().Value[0].GetMap().Items["importPath"].GetS(), "import unexpected value")
	assert.Equal(t, "package2", m.GetList().Value[1].GetMap().Items["importPath"].GetS(), "import unexpected value")

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
