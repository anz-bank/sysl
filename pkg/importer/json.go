package importer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
func LoadJSONText(args OutputData, jsonText string, logger *logrus.Logger) (out string, err error) {
	if len(jsonText) == 0 {
		return "", fmt.Errorf("processed json file is empty")
	}
	if len(args.Transform.Apps) > 1 {
		return "", fmt.Errorf("it can have only one app in file %s", args.TransformFile)
	}

	types := TypeList{}
	for _, app := range args.Transform.Apps {
		mTags, pTags := getTypeMeta(app)
		var data interface{}
		err := json.Unmarshal([]byte(strings.TrimSpace(jsonText)), &data)
		if err != nil {
			return "", fmt.Errorf("can't unmarshall json file")
		}

		if isArray(jsonText) {
			buildType(mTags, pTags, data, &types)
		} else {
			return "", fmt.Errorf("can't support current json")
		}
	}
	types.Sort()

	info := SyslInfo{
		OutputData:  args,
		Description: "",
		Title:       "",
	}
	result := &bytes.Buffer{}
	w := newWriter(result, logger)
	if err := w.Write(info, types); err != nil {
		return "", err
	}

	return result.String(), nil
}

func getTypeMeta(app *sysl.Application) (modelTags []string, propertyTags []string) {
	var mTags []string
	for _, attrDef := range app.GetTypes()["model"].GetTuple().GetAttrDefs() {
		mTags = append(mTags, attrDef.GetAttrs()[jsonTag].GetS())
	}

	var pTags []string
	for _, attrDef := range app.GetTypes()["property"].GetTuple().GetAttrDefs() {
		pTags = append(pTags, attrDef.GetAttrs()[jsonTag].GetS())
	}

	return mTags, pTags
}

func buildType(modelTags []string, propertyTags []string, data interface{}, typeList *TypeList) {
	mCount := 0
	mPath := "$[" + strconv.Itoa(mCount) + "]"
	for mName, mErr := jsonpath.Read(data, mPath+".name"); mErr == nil; mName, mErr = jsonpath.Read(data, mPath+".name") {
		model := StandardType{name: mName.(string)}
		fCount := 0
		fPath := mPath + ".fields[" + strconv.Itoa(fCount) + "]"
		for _, fsErr := jsonpath.Read(data, fPath); fsErr == nil; _, fsErr = jsonpath.Read(data, fPath) {
			for _, fTag := range propertyTags {
				if fieldTag, fErr := jsonpath.Read(data, fPath+"."+fTag); fErr == nil {
					model.Properties = append(model.Properties, Field{Name: fieldTag.(string),
						Type: &SyslBuiltIn{name: "string"}})
				}
			}
			fCount++
			fPath = mPath + ".fields[" + strconv.Itoa(fCount) + "]"
		}

		typeList.Add(&model)
		mCount++
		mPath = "$[" + strconv.Itoa(mCount) + "]"
	}
}

// Check JSON root element is array or map/object
func isArray(json string) (isArray bool) {
	bytes := []byte(strings.TrimSpace(json))
	array := (len(bytes) > 0 && bytes[0] == '[')
	return array
}
