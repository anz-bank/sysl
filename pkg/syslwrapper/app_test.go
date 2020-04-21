package syslutil

import (
	"encoding/json"
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func readSyslModule(filename string) (*sysl.Module, error) {
	mod, err := parse.NewParser().Parse(filename,
		syslutil.NewChrootFs(afero.NewOsFs(), "."))
	return mod, err
}

func TestMapRest(t *testing.T) {
	mod, err := readSyslModule("./tests/rest.sysl")
	assert.NoError(t, err)
	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()
	mapper.resolveTypes()
	simpleApps, err := mapper.Map()
	assert.NoError(t, err)
	printStr, _ := json.MarshalIndent(simpleApps, "", " ")
	t.Log(string(printStr))
	assert.Equal(t, "string", simpleApps["SampleRestApp"].Endpoints["POST /login/{CustomerID}"].Params["CustomerID"].Type.Type)
	assert.Equal(t, "string", simpleApps["SampleRestApp"].Endpoints["POST /login/{CustomerID}"].Params["newPost"].Type.Type)
	assert.Equal(t, "string", simpleApps["SampleRestApp"].Endpoints["POST /login/{CustomerID}"].Response["default"].Type.Type)
	assert.Equal(t, "string", simpleApps["SampleRestApp"].Endpoints["POST /post"].Params["newPost"].Type.Type)
	assert.Equal(t, "string", simpleApps["SampleRestApp"].Endpoints["GET /post"].Params["PostId"].Type.Type)
}
func TestMap(t *testing.T) {
	mod, err := readSyslModule("./tests/types.sysl")
	assert.NoError(t, err)
	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()
	mapper.resolveTypes()
	simpleApps, err := mapper.Map()
	assert.NoError(t, err)
	assert.Equal(t, "", simpleApps["Server"].Types["Response"].Properties["balance"].Type)
	assert.Equal(t, "tuple", simpleApps["Server"].Types["Response"].Properties["query"].Type)
	assert.Equal(t, "int", simpleApps["Server"].Types["Response"].Properties["query"].Properties["amount"].Type)
	assert.Equal(t, "tuple", simpleApps["MobileApp"].Endpoints["Login"].Params["input"].Type.Type)
	assert.Equal(t, "tuple", simpleApps["MobileApp"].Endpoints["Login"].Params["input"].Type.Type)
	assert.Equal(t, "string", simpleApps["MobileApp"].Endpoints["Login"].Params["input"].Type.Properties["query"].Type)
}
func TestResolveTypesWithSyslFile(t *testing.T) {
	mod, err := readSyslModule("./tests/types.sysl")
	assert.NoError(t, err)
	mapper := MakeAppMapper(mod)
	typeIndex := mapper.IndexTypes()
	mapper.resolveTypes()

	expectedResult := MakeTuple(map[string]*sysl.Type{
		"query":   MakePrimitive("int"),
		"balance": MakePrimitive("empty"),
	})

	assert.Equal(t, expectedResult.GetAttrs()["query"], typeIndex["Server:Response"].GetAttrs()["query"])
	assert.Equal(t, expectedResult.GetAttrs()["balance"], typeIndex["Server:Response"].GetAttrs()["balance"])
}

func TestMapTypeRef(t *testing.T) {
	type1 := MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"})
	type2 := MakePrimitive("string")
	var app1 = MakeApp("app1", []*sysl.Param{}, map[string]*sysl.Type{"typeName1": type1})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": &app1,
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()

	mappedType := mapper.MapType(type1)
	assert.Equal(t, "ref", mappedType.Type)
	assert.Equal(t, "app2:request", mappedType.Reference)
}
func TestResolveNonExistentType(t *testing.T) {
	type1 := MakeTypeRef("app1", []string{"login"}, "app2", []string{"nonexist"})
	param1 := MakeParam("Login", type1)
	var app1 = MakeApp("app1", []*sysl.Param{param1}, map[string]*sysl.Type{"list": nil})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": nil})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": &app1,
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()

	mapper.resolveTypes()
	assert.Equal(t, nil, mapper.Types["app2:request"])
	assert.Equal(t, nil, mapper.Types["app1:list"])
}
func TestResolveTypesNil(t *testing.T) {
	type1 := MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"})
	var app1 = MakeApp("app1", []*sysl.Param{}, map[string]*sysl.Type{"list": type1})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": nil})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": &app1,
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	typeIndex := mapper.IndexTypes()
	mapper.resolveTypes()

	assert.Equal(t, nil, typeIndex["app2:request"])
	assert.Equal(t, &sysl.Type{}, typeIndex["app1:list"])
}

func TestResolveTypeOneOf(t *testing.T) {
	type1 := MakeOneOf([]*sysl.Type{MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"})})
	type2 := MakePrimitive("string")
	param1 := MakeParam("Login", type1)
	var app1 = MakeApp("app1", []*sysl.Param{param1}, map[string]*sysl.Type{})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": &app1,
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	typeIndex := mapper.IndexTypes()

	syslType := mapper.resolveType(type1)
	assert.Equal(t, type2, typeIndex["app2"+":"+"request"])
	assert.Equal(t, MakeOneOf([]*sysl.Type{MakePrimitive("string")}), type1)
	assert.Equal(t, MakeOneOf([]*sysl.Type{MakePrimitive("string")}), syslType)
}

func TestResolveTypeMap(t *testing.T) {
	type1 := MakeMap(MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"}), MakePrimitive("string"))
	type2 := MakePrimitive("string")
	param1 := MakeParam("Login", type1)
	var app1 = MakeApp("app1", []*sysl.Param{param1}, map[string]*sysl.Type{})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": &app1,
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	typeIndex := mapper.IndexTypes()

	syslType := mapper.resolveType(type1)
	assert.Equal(t, type2, typeIndex["app2:request"])
	assert.Equal(t, MakeMap(MakePrimitive("string"), MakePrimitive("string")), type1)
	assert.Equal(t, MakeMap(MakePrimitive("string"), MakePrimitive("string")), syslType)
}
func TestResolveTypeList(t *testing.T) {
	type1 := MakeList(MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"}))
	type2 := MakePrimitive("string")
	param1 := MakeParam("Login", type1)
	var app1 = MakeApp("app1", []*sysl.Param{param1}, map[string]*sysl.Type{})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": &app1,
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	typeIndex := mapper.IndexTypes()

	syslType := mapper.resolveType(type1)
	assert.Equal(t, type2, typeIndex["app2:request"])
	assert.Equal(t, MakeList(MakePrimitive("string")), type1)
	assert.Equal(t, MakeList(MakePrimitive("string")), syslType)
}

func TestResolveTypeTypeRef(t *testing.T) {
	type1 := MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"})
	type2 := MakePrimitive("string")
	param1 := MakeParam("Login", type1)
	var app1 = MakeApp("app1", []*sysl.Param{param1}, map[string]*sysl.Type{})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": &app1,
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	typeIndex := mapper.IndexTypes()

	syslType := mapper.resolveType(type1)
	assert.Equal(t, type2, typeIndex["app2:request"])
	assert.Equal(t, MakePrimitive("string"), syslType)
}

func TestTypesFromRef(t *testing.T) {
	type1 := MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"})
	type2 := MakePrimitive("string")
	param1 := MakeParam("Login", type1)
	param2 := MakeParam("Auth", MakePrimitive("string"))
	var app1 = MakeApp("app1", []*sysl.Param{param1}, map[string]*sysl.Type{})
	var app2 = MakeApp("app2", []*sysl.Param{param2}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": &app1,
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()
	syslType, err := mapper.MapSyslType(mod.Apps["app1"].Endpoints["testEndpoint"].Param[0].Type)
	if err != nil {
		t.Error(err)
	}
	typeIndex := mapper.IndexTypes()
	assert.Equal(t, type2, typeIndex["app2:request"])
	assert.Equal(t, type2, syslType)
}

func TestTypeConversionPrimative(t *testing.T) {
	type2 := MakePrimitive("string")
	param2 := MakeParam("Auth", MakePrimitive("string"))
	var app2 = MakeApp("app2", []*sysl.Param{param2}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": &app2,
		},
	}

	mapper := MakeAppMapper(mod)
	convertedType1 := mapper.MapType(type2)
	assert.Equal(t, &Type{
		Type: "string",
	}, convertedType1)
}
func TestTypeConversionList(t *testing.T) {
	type2 := MakeList(MakePrimitive("string"))
	param2 := MakeParam("Auth", MakePrimitive("string"))
	var app2 = MakeApp("app2", []*sysl.Param{param2}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": &app2,
		},
	}

	expectedResult := &Type{
		Type: "list",
		Items: []*Type{
			{
				Type: "string",
			},
		},
	}

	mapper := MakeAppMapper(mod)
	convertedType1 := mapper.MapType(type2)
	assert.Equal(t, expectedResult, convertedType1)
}

func TestTypeConversionMap(t *testing.T) {
	type2 := MakeMap(MakePrimitive("string"), MakePrimitive("string"))
	param2 := MakeParam("Auth", MakePrimitive("string"))
	var app2 = MakeApp("app2", []*sysl.Param{param2}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": &app2,
		},
	}

	expectedResult := &Type{
		Type: "map",
		Items: []*Type{
			{
				Type: "string",
			}, {
				Type: "string",
			},
		},
	}
	mapper := MakeAppMapper(mod)
	convertedType1 := mapper.MapType(type2)
	assert.Equal(t, expectedResult, convertedType1)
}

func TestTypeConversionEnum(t *testing.T) {
	enumerate := map[string]int64{
		"apple":  1,
		"orange": 2,
	}
	mapper := MakeAppMapper(&sysl.Module{})
	typeToConvert := MakeEnum(enumerate)
	expectedResult := &Type{
		Type: "enum",
		Enum: enumerate,
	}
	convertedType1 := mapper.MapType(typeToConvert)
	assert.Equal(t, expectedResult, convertedType1)
}
func TestTypeConversionNoType(t *testing.T) {
	expectedResult := &Type{
		Type: "notype",
	}
	mapper := MakeAppMapper(&sysl.Module{})
	convertedType1 := mapper.MapType(MakeNoType())
	assert.Equal(t, expectedResult, convertedType1)
}

func TestConvertPrimitive(t *testing.T) {
	mapper := MakeAppMapper(&sysl.Module{})
	result := mapper.convertPrimitive("Primitive:STRING Source Context")
	assert.Equal(t, "string", result)
}
