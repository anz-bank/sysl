//nolint:lll
package syslwrapper

import (
	"encoding/json"
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestMapRest(t *testing.T) {
	mod, err := parse.NewParser().ParseFs("./tests/rest.sysl", afero.NewOsFs())
	assert.NoError(t, err)
	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()
	mapper.ResolveTypes()
	simpleApps, err := mapper.Map()
	assert.NoError(t, err)
	prettyPrint(t, simpleApps)
	assert.Equal(t, "string", simpleApps["SampleRestApp"].Endpoints["POST /login/{CustomerID}"].Params["CustomerID"].Type.Type)
	assert.Equal(t, "string", simpleApps["SampleRestApp"].Endpoints["POST /login/{CustomerID}"].Params["newPost"].Type.Type)
	assert.Equal(t, "ref", simpleApps["SampleRestApp"].Endpoints["POST /post"].Params["PostID"].Type.Type)
	assert.Equal(t, "string", simpleApps["SampleRestApp"].Endpoints["GET /post"].Params["PostId"].Type.Type)
}
func TestMap(t *testing.T) {
	mod, err := parse.NewParser().ParseFs("./tests/types.sysl", afero.NewOsFs())
	assert.NoError(t, err)
	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()
	mapper.ResolveTypes()
	simpleApps, err := mapper.Map()
	assert.NoError(t, err)
	prettyPrint(t, simpleApps)
	assert.Equal(t, "", simpleApps["Server"].Types["Response"].Properties["balance"].Type)
	assert.Equal(t, "tuple", simpleApps["Server"].Types["Response"].Properties["query"].Type)
	assert.Equal(t, "int", simpleApps["Server"].Types["Response"].Properties["query"].Properties["amount"].Type)
	assert.Equal(t, "ref", simpleApps["MobileApp"].Endpoints["Login"].Params["input"].Type.Type)
	assert.Equal(t, "Server:Request", simpleApps["MobileApp"].Endpoints["Login"].Params["input"].Type.Reference)
}
func TestResolveTypesWithSyslFile(t *testing.T) {
	mod, err := parse.NewParser().ParseFs("./tests/types.sysl", afero.NewOsFs())
	assert.NoError(t, err)
	mapper := MakeAppMapper(mod)
	typeIndex := mapper.IndexTypes()
	mapper.ResolveTypes()

	expectedResult := MakeTuple(map[string]*sysl.Type{
		"query":   MakePrimitive("int"),
		"balance": MakePrimitive("empty"),
	})

	assert.Equal(t, expectedResult.GetAttrs()["query"], typeIndex["Server:Response"].GetAttrs()["query"])
	assert.Equal(t, expectedResult.GetAttrs()["balance"], typeIndex["Server:Response"].GetAttrs()["balance"])
}

func TestConvertMapType(t *testing.T) {
	mod, err := parse.NewParser().ParseFs("./tests/map.sysl", afero.NewOsFs())
	assert.NoError(t, err)
	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()
	simpleTypes := mapper.ConvertTypes()

	expectedResult := &Type{
		Type: "map",
		Properties: map[string]*Type{
			"item_id": {
				Type: "string",
			},
			"quantity": {
				Type: "int",
			},
			"message": {
				Type:      "ref",
				Reference: "MapType:Message",
			},
		},
		PrimaryKey: "item_id",
	}

	assert.Equal(t, expectedResult, simpleTypes["MapType:InventoryResponse"])
}
func TestMapPetStoreToSimpleTypes(t *testing.T) {
	mod, err := parse.NewParser().ParseFs("../../demo/petshop/petshop.sysl", afero.NewOsFs())
	assert.NoError(t, err)
	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()
	mapper.ConvertTypes()
	expected := &Type{
		Type:      "ref",
		Reference: "PetShopModel:Breed",
	}
	assert.Equal(t, "relation", mapper.SimpleTypes["PetShopModel:Pet"].Type)
	assert.Equal(t, expected, mapper.SimpleTypes["PetShopModel:Pet"].Properties["breedId"])
}

func TestMapTypeRef(t *testing.T) {
	type1 := MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"})
	type2 := MakePrimitive("string")
	var app1 = MakeApp("app1", []*sysl.Param{}, map[string]*sysl.Type{"typeName1": type1})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": app1,
			"app2": app2,
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
			"app1": app1,
			"app2": app2,
		},
	}

	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()

	mapper.ResolveTypes()

	assert.Equal(t, nil, mapper.Types["app2:request"].GetType())
	assert.Equal(t, nil, mapper.Types["app1:list"].GetType())
}
func TestResolveTypesNil(t *testing.T) {
	type1 := MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"})
	var app1 = MakeApp("app1", []*sysl.Param{}, map[string]*sysl.Type{"list": type1})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": nil})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": app1,
			"app2": app2,
		},
	}

	mapper := MakeAppMapper(mod)
	typeIndex := mapper.IndexTypes()
	mapper.ResolveTypes()
	prettyPrint(t, typeIndex["app2:request"])
	assert.Equal(t, nil, typeIndex["app2:request"].GetType())
	assert.Equal(t, nil, typeIndex["app1:list"].GetType())
}

func TestResolveTypeOneOf(t *testing.T) {
	type1 := MakeOneOf([]*sysl.Type{MakeTypeRef("app1", []string{"login"}, "app2", []string{"request"})})
	type2 := MakePrimitive("string")
	param1 := MakeParam("Login", type1)
	var app1 = MakeApp("app1", []*sysl.Param{param1}, map[string]*sysl.Type{})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app1": app1,
			"app2": app2,
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
			"app1": app1,
			"app2": app2,
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
			"app1": app1,
			"app2": app2,
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
			"app1": app1,
			"app2": app2,
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
			"app1": app1,
			"app2": app2,
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
			"app2": app2,
		},
	}

	mapper := MakeAppMapper(mod)
	convertedType1 := mapper.MapType(type2)
	assert.Equal(t, &Type{
		Type: "string",
	}, convertedType1)
}

func TestTypeConversionRelationNoPrimaryKey(t *testing.T) {
	relation := &sysl.Type{
		Type: &sysl.Type_Relation_{
			Relation: &sysl.Type_Relation{
				AttrDefs: map[string]*sysl.Type{
					"id":   MakePrimitive("string"),
					"name": MakePrimitive("string"),
				},
			},
		},
	}
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": relation})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": app2,
		},
	}

	expectedResult := &Type{
		Type: "relation",
		Properties: map[string]*Type{
			"id": {
				Type: "string",
			},
			"name": {
				Type: "string",
			},
		},
	}

	mapper := MakeAppMapper(mod)
	convertedType1 := mapper.MapType(relation)
	assert.Equal(t, expectedResult, convertedType1)
}

func TestTypeConversionRelation(t *testing.T) {
	type2 := MakeRelation(map[string]*sysl.Type{
		"id":   MakePrimitive("string"),
		"name": MakePrimitive("string"),
	}, "id", []string{})
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": app2,
		},
	}

	expectedResult := &Type{
		Type:       "relation",
		PrimaryKey: "id",
		Properties: map[string]*Type{
			"id": {
				Type: "string",
			},
			"name": {
				Type: "string",
			},
		},
	}

	mapper := MakeAppMapper(mod)
	convertedType1 := mapper.MapType(type2)
	assert.Equal(t, expectedResult, convertedType1)
}
func TestTypeConversionSet(t *testing.T) {
	type2 := MakeSet(MakePrimitive("string"))
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": app2,
		},
	}

	expectedResult := &Type{
		Type: "set",
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
func TestTypeConversionList(t *testing.T) {
	type2 := MakeList(MakePrimitive("string"))
	param2 := MakeParam("Auth", MakePrimitive("string"))
	var app2 = MakeApp("app2", []*sysl.Param{param2}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": app2,
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

func TestMapReturnStatements(t *testing.T) {
	type2 := MakeMap(MakePrimitive("string"), MakePrimitive("string"))
	var app2 = MakeApp("app2", []*sysl.Param{}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": app2,
		},
	}

	mapper := MakeAppMapper(mod)
	mapper.IndexTypes()
	statements := []*sysl.Statement{
		MakeReturnStatement("string"),
		MakeReturnStatement("default <: string"),
		MakeReturnStatement("301 <: set of app2.request"),
		MakeReturnStatement("401 <: request"),
		MakeReturnStatement("404 <: sequence of request"),
		MakeReturnStatement("500 <: sequence of app2.request"),
	}
	result := mapper.mapResponse(statements, "app2")
	expectedResult := map[string]*Parameter{
		"200": {
			Name: "200",
			Type: &Type{
				Type: "string",
			},
		},
		"301": {
			Name: "301",
			Type: &Type{
				Items: []*Type{
					{
						Type:      "ref",
						Reference: "app2:request",
					},
				},
				Type: "set",
			},
		},
		"default": {
			Name: "default",
			Type: &Type{
				Type: "string",
			},
		},
		"401": {
			Name: "401",
			Type: &Type{
				Reference: "app2:request",
				Type:      "ref",
			},
		},
		"404": {
			Name: "404",
			Type: &Type{
				Items: []*Type{
					{
						Type:      "ref",
						Reference: "app2:request",
					},
				},
				Type: "list",
			},
		},
		"500": {
			Name: "500",
			Type: &Type{
				Items: []*Type{
					{
						Type:      "ref",
						Reference: "app2:request",
					},
				},
				Type: "list",
			},
		},
	}

	assert.Equal(t, expectedResult, result)
}

func TestTypeConversionMap(t *testing.T) {
	type2 := MakeMap(MakePrimitive("string"), MakePrimitive("string"))
	param2 := MakeParam("Auth", MakePrimitive("string"))
	var app2 = MakeApp("app2", []*sysl.Param{param2}, map[string]*sysl.Type{"request": type2})
	var mod = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"app2": app2,
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
	enumerate := map[int64]string{
		1: "apple",
		2: "orange",
	}
	mapper := MakeAppMapper(&sysl.Module{})
	typeToConvert := MakeEnum(enumerate)
	convertedType1 := mapper.MapType(typeToConvert)
	assert.Equal(t, "enum", convertedType1.Type)
	assert.Equal(t, "apple", convertedType1.Enum[1])
	assert.Equal(t, "orange", convertedType1.Enum[2])
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

func prettyPrint(t *testing.T, v interface{}) {
	json, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		t.Log(t, err)
	}
	t.Log(t, string(json))
}
