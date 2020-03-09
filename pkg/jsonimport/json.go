package jsonimport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/anz-bank/sysl/pkg/eval"
	"github.com/anz-bank/sysl/pkg/importer"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/yalp/jsonpath"
)

const jsonTag = "json_tag"

/*
LoadJSONText produces sysl file with json file. Please note this json file is not a json schema format, it looks like:
[
	{
		"book":{"name":"Book1"}
	},
	{
		"book":{"name":"Book2"}
	}
]
or
{
	"books" : [
		{
			"book" : {"name":"Book1"}
		},
		{
			"book" : {"name":"Book2"}
		}
	]
}
So it needs transform to do a mapping.
*/
func LoadJSONText(args importer.OutputData, jsonText string, logger *logrus.Logger) (string, error) {
	if len(jsonText) == 0 {
		return "", fmt.Errorf("processed json file is empty")
	}
	if args.Transform == nil {
		return "", fmt.Errorf("transform is essential, but it is empty")
	}
	if len(args.Transform.Apps) > 1 {
		return "", fmt.Errorf("it can have only one app in file %s", args.TransformFile)
	}

	types := importer.TypeList{}
	for _, app := range args.Transform.Apps {
		mNameTag, pTags := getTypeMeta(app)
		var data interface{}
		err := json.Unmarshal([]byte(strings.TrimSpace(jsonText)), &data)
		if err != nil {
			return "", fmt.Errorf("can't unmarshall json file")
		}

		if isArray(jsonText) {
			types = loadModelTypes(mNameTag, data, logger)
			buildModelTypesDetail(app, mNameTag, pTags, data, &types, logger)
		} else {
			return "", fmt.Errorf("can't support current json")
		}
	}

	info := importer.SyslInfo{
		OutputData:  args,
		Description: "",
		Title:       "",
	}
	result := &bytes.Buffer{}
	w := importer.NewWriter(result, logger)
	if err := w.Write(info, types); err != nil {
		return "", err
	}

	return result.String(), nil
}

// Use JsonPath to get data in json file
func loadModelTypes(modelNameTag string, data interface{},
	logger *logrus.Logger) importer.TypeList {
	types := importer.TypeList{}

	mCount := 0
	mPath := "$[" + strconv.Itoa(mCount) + "]"
	mName, mErr := jsonpath.Read(data, mPath+"."+modelNameTag)
	for ; mErr == nil; mName, mErr = jsonpath.Read(data, mPath+".name") {
		model := importer.StandardType{}
		model.SetName(mName.(string))
		types.Add(&model)
		mCount++
		mPath = "$[" + strconv.Itoa(mCount) + "]"
	}

	types.Sort()
	return types
}

// Use JsonPath to get data in json file
func buildModelTypesDetail(app *sysl.Application, modelNameTag string, propertyTags map[string]string, data interface{},
	knownTypeList *importer.TypeList, logger *logrus.Logger) {
	mCount := 0
	mPath := "$[" + strconv.Itoa(mCount) + "]"
	mName, mErr := jsonpath.Read(data, mPath+"."+modelNameTag)
	for ; mErr == nil; mName, mErr = jsonpath.Read(data, mPath+".name") {
		modelType, existed := knownTypeList.Find(mName.(string))
		if !existed {
			logger.Warnf("found type %s which is not existing in json", mName.(string))
			continue
		}

		switch model := modelType.(type) {
		case *importer.StandardType:
			fCount := 0
			fPath := mPath + ".fields[" + strconv.Itoa(fCount) + "]"
			for _, fsErr := jsonpath.Read(data, fPath); fsErr == nil; _, fsErr = jsonpath.Read(data, fPath) {
				fieldName, fnErr := jsonpath.Read(data, fPath+"."+propertyTags["name"])
				fieldType, ftErr := jsonpath.Read(data, fPath+"."+propertyTags["type"])
				if fnErr == nil && ftErr == nil && fieldName != "" && fieldType != "" {
					jsonType := makeJSONType(app, fieldType.(string), fPath, data, knownTypeList)
					model.Properties = append(model.Properties, importer.Field{Name: fieldName.(string),
						Type: jsonType})
				} else {
					logger.Errorf("can't process field %s", fPath+"."+propertyTags["name"])
				}

				fCount++
				fPath = mPath + ".fields[" + strconv.Itoa(fCount) + "]"
			}
		default:
		}
		mCount++
		mPath = "$[" + strconv.Itoa(mCount) + "]"
	}
}

func getTypeMeta(app *sysl.Application) (modelNameTag string, propertyTags map[string]string) {
	mNameAttrDef := app.GetTypes()["model"].GetTuple().GetAttrDefs()["name"]
	mNameTag := mNameAttrDef.GetAttrs()[jsonTag].GetS()

	var pTags map[string]string = make(map[string]string)
	for key, attrDef := range app.GetTypes()["property"].GetTuple().GetAttrDefs() {
		pTags[key] = attrDef.GetAttrs()[jsonTag].GetS()
	}

	return mNameTag, pTags
}

func makeJSONType(app *sysl.Application, jsonDataType string, jsonPath string, data interface{},
	knownTypeList *importer.TypeList) importer.Type {
	if jsonDataType == "Foreign Key" {
		// Some foreign key can't find model, if it can't find refrenced model, return default builtin type
		dataType := makeJSONComplexType(jsonPath, data, knownTypeList)
		if dataType == nil {
			syslBuiltInType := importer.SyslBuiltIn{}
			syslBuiltInType.SetName("string")
			return &syslBuiltInType
		}
		return dataType
	}
	return makeJSONBuiltinType(app, jsonDataType)
}

func makeJSONComplexType(jsonPath string, data interface{}, knownTypeList *importer.TypeList) importer.Type {
	obj, oErr := jsonpath.Read(data, jsonPath+".relatedObject")
	_, fErr := jsonpath.Read(data, jsonPath+".referencedField")
	if oErr == nil && fErr == nil {
		modelType, existed := knownTypeList.Find(obj.(string))
		if existed {
			return modelType
		}
	}

	return nil
}

func makeJSONBuiltinType(app *sysl.Application, jsonDataType string) importer.Type {
	s := eval.Scope{}
	s.AddString("t", jsonDataType)
	out := evaluateView(app, "DataType", s)
	typeStr := out.GetMap().GetItems()["type"].GetS()

	syslBuiltInType := importer.SyslBuiltIn{}
	syslBuiltInType.SetName(typeStr)
	return &syslBuiltInType
}

// Check JSON root element is array or map/object
func isArray(json string) bool {
	bytes := []byte(strings.TrimSpace(json))
	array := (len(bytes) > 0 && bytes[0] == '[')
	return array
}

func evaluateView(app *sysl.Application, viewName string, s eval.Scope) *sysl.Value {
	view := app.Views[viewName]
	if view.Expr.Type == nil {
		view.Expr.Type = view.RetType
	}

	return eval.EvaluateApp(app, view, s)
}
