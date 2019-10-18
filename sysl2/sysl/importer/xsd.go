package importer

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"aqwari.net/xml/xsd"
	"github.com/sirupsen/logrus"
)

func LoadXSDText(args OutputData, text string, logger *logrus.Logger) (out string, err error) {
	xsd.StandardSchema = [][]byte{} // Ignore all the standard schemas,
	specs, err := xsd.Parse([]byte(text))
	if err != nil {
		return "", err
	}

	types := TypeList{}
	for _, schema := range specs {
		schemaTypes := loadSchemaTypes(schema, logger)
		types.Add(schemaTypes.Items()...)
	}
	types.Sort()

	info := SyslInfo{
		OutputData:  args,
		Description: "",
		Title:       "",
	}

	result := &bytes.Buffer{}
	w := newWriter(result, logger)
	if err := w.Write(info, types, ""); err != nil {
		return "", err
	}

	return result.String(), nil
}

func makeNamespacedType(name xml.Name, target Type) Type {
	return &ImportedBuiltInAlias{
		name:   fmt.Sprintf("%s:%s", name.Space, name.Local),
		Target: target,
	}
}

func loadSchemaTypes(schema xsd.Schema, logger *logrus.Logger) TypeList {
	types := TypeList{}

	var xsdToSyslMappings = map[string]string{
		"NMTOKEN":      StringTypeName,
		"integer":      "int",
		"time":         StringTypeName,
		"boolean":      "bool",
		StringTypeName: StringTypeName,
		"date":         "date",
	}
	for swaggerName, syslName := range xsdToSyslMappings {
		from := xml.Name{Space: "http://www.w3.org/2001/XMLSchema", Local: swaggerName}
		to := &SyslBuiltIn{syslName}
		types.Add(makeNamespacedType(from, to))
	}

	keys := make([]xml.Name, 0, len(schema.Types))
	for key := range schema.Types {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.Compare(
			fmt.Sprintf("%s:%s", keys[i].Space, keys[i].Local),
			fmt.Sprintf("%s:%s", keys[j].Space, keys[j].Local)) < 0
	})

	for _, name := range keys {
		data := schema.Types[name]
		if name.Local == "_self" {
			rootType := data.(*xsd.ComplexType)
			data = rootType.Elements[0].Type
			t := findType(data, &types)
			if t == nil {
				t = makeType(name, data, &types, logger)
				if name.Space != "" {
					types.Add(makeNamespacedType(name, t))
				}
			}
			x := t.(*StandardType)
			if !contains("~xml_root", x.Attributes) {
				x.Attributes = append(x.Attributes, "~xml_root")
			}
		} else if t := findType(data, &types); t == nil {
			types.Add(makeType(name, data, &types, logger))
		}
	}

	return types
}

func findType(t xsd.Type, knownTypes *TypeList) Type {
	if res, found := knownTypes.Find(fmt.Sprintf("%s:%s", xsd.XMLName(t).Space, xsd.XMLName(t).Local)); found {
		return res
	}
	if res, found := knownTypes.Find(xsd.XMLName(t).Local); found {
		return res
	}
	return nil
}

func makeType(_ xml.Name, from xsd.Type, knownTypes *TypeList, logger *logrus.Logger) Type {
	switch t := from.(type) {
	case *xsd.ComplexType:
		return makeComplexType(t, knownTypes, logger)
	case *xsd.SimpleType:
		return makeSimpleType(t, knownTypes, logger)
	case xsd.Builtin:
		return makeXsdBuiltinType(t, knownTypes)
	}
	return nil
}

func makeComplexType(from *xsd.ComplexType, knownTypes *TypeList, logger *logrus.Logger) Type {
	createChildItem := func(name xml.Name, data xsd.Type, isAttr, optional, plural bool) Field {
		childType := findType(data, knownTypes)
		if childType == nil {
			childType = makeType(name, data, knownTypes, logger)
			knownTypes.Add(childType)
		}
		f := Field{
			Name: name.Local,
			Type: childType,
		}
		if isAttr {
			f.Attributes = []string{"~xml_attribute"}
		}
		f.Optional = optional

		if from, ok := data.(*xsd.SimpleType); ok && from.Restriction.Min > 0 {
			spec := sizeSpec{
				Min: int(from.Restriction.Min),
				Max: int(from.Restriction.Max),
			}
			if spec.Max > 0 {
				spec.MaxType = MaxSpecified
			}
			f.SizeSpec = &spec
		}

		if plural {
			f.Type = &Array{Items: f.Type}
		}
		return f
	}

	item := &StandardType{
		name: from.Name.Local,
	}

	for _, child := range from.Elements {
		c := createChildItem(child.Name, child.Type, false, child.Optional, child.Plural)
		if c.SizeSpec == nil {
			c.SizeSpec = makeSizeSpecFromAttrs(child.Attr)
		}
		item.Properties = append(item.Properties, c)
	}
	for _, child := range from.Attributes {
		c := createChildItem(child.Name, child.Type, true, child.Optional, child.Plural)
		if c.SizeSpec == nil {
			c.SizeSpec = makeSizeSpecFromAttrs(child.Attr)
		}
		item.Properties = append(item.Properties, c)
	}

	return item
}

func makeSizeSpecFromAttrs(attrs []xml.Attr) *sizeSpec {
	ss := sizeSpec{}
	changed := false
	for _, attr := range attrs {
		switch attr.Name.Local {
		case "maxOccurs":
			if attr.Value == "unbounded" {
				ss.MaxType = OpenEnded
				changed = true
			} else if val, err := strconv.ParseInt(attr.Value, 0, 32); err == nil {
				ss.MaxType = MaxSpecified
				ss.Max = int(val)
				changed = true
			}
		case "minOccurs":
			if val, err := strconv.ParseInt(attr.Value, 0, 32); err == nil {
				ss.Min = int(val)
				changed = changed || val > 0
			}
		}
	}
	if changed {
		return &ss
	}
	return nil
}

func makeSimpleType(from *xsd.SimpleType, knownTypes *TypeList, logger *logrus.Logger) Type {
	item := &Alias{
		name:   from.Name.Local,
		Target: makeType(from.Name, from.Base, knownTypes, logger),
	}

	return item
}

func makeXsdBuiltinType(from xsd.Builtin, knownTypes *TypeList) Type {
	typeStr := StringTypeName
	switch from {
	case xsd.Integer, xsd.Int:
		typeStr = "int"
	}
	t, _ := knownTypes.Find(typeStr)
	return t
}
