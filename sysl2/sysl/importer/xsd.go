package importer

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"aqwari.net/xml/xsd"
	"github.com/sirupsen/logrus"
)

func LoadXSDText(args OutputData, text string, targetNS string, logger *logrus.Logger) (out string, err error) {
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
		"NMTOKEN": "string",
		"integer": "int",
		"time":    "string",
		"boolean": "bool",
		"string":  "string",
		"date":    "date",
	}
	for swaggerName, syslName := range xsdToSyslMappings {
		from := xml.Name{Space: "http://www.w3.org/2001/XMLSchema", Local: swaggerName}
		to := &SyslBuiltIn{syslName}
		types.Add(makeNamespacedType(from, to))
	}

	for name, data := range schema.Types {
		if name.Local == "_self" {
			rootType := data.(*xsd.ComplexType)
			data = rootType.Elements[0].Type
			t := FindType(data, &types)
			if t == nil {
				t = makeType(name, data, &types, logger)
				types.Add(t)
				if name.Space != "" {
					types.Add(makeNamespacedType(name, t))
				}
			}
			x := t.(*StandardType)
			if !contains("~xml_root", x.Attributes) {
				x.Attributes = append(x.Attributes, "~xml_root")
			}
		} else if t := FindType(data, &types); t == nil {
			item := makeType(name, data, &types, logger)
			if name.Space != "" {
				item = makeNamespacedType(name, item)
			}
			types.Add(item)
		}
	}

	return types
}

func FindType(t xsd.Type, knownTypes *TypeList) Type {
	if res, found := knownTypes.Find(fmt.Sprintf("%s:%s", xsd.XMLName(t).Space, xsd.XMLName(t).Local)); found {
		return res
	}
	if res, found := knownTypes.Find(xsd.XMLName(t).Local); found {
		return res
	}
	return nil
}

func makeType(name xml.Name, from xsd.Type, knownTypes *TypeList, logger *logrus.Logger) Type {
	switch t := from.(type) {
	case *xsd.ComplexType:
		return makeComplexType(name, t, knownTypes, logger)
	case *xsd.SimpleType:
		return makeSimpleType(name, t, logger)
	case xsd.Builtin:
		return makeXsdBuiltinType(t, knownTypes)
	}
	return nil
}

func makeComplexType(_ xml.Name, from *xsd.ComplexType, knownTypes *TypeList, logger *logrus.Logger) Type {
	if isArray(from) {
		panic("Cant make a complex array yet")
	}

	createChildItem := func(name xml.Name, data xsd.Type) Field {
		childType := FindType(data, knownTypes)
		if childType == nil {
			childType = makeType(name, data, knownTypes, logger)
			knownTypes.Add(childType)
		}
		f := Field{
			Name: name.Local,
			Type: childType,
		}
		return f
	}

	configureChildItem := func(field *Field, isAttr, optional, plural bool) {
		if isAttr {
			field.Attributes = []string{"~xml_attribute"}
		}
		field.Optional = optional

		if plural {
			field.Type = &Array{Items: field.Type}
		}
	}

	item := &StandardType{
		name: from.Name.Local,
	}

	for _, child := range from.Elements {
		c := createChildItem(child.Name, child.Type)
		configureChildItem(&c, false, child.Optional, child.Plural)
		item.Properties = append(item.Properties, c)
	}
	for _, child := range from.Attributes {
		c := createChildItem(child.Name, child.Type)
		configureChildItem(&c, true, child.Optional, child.Plural)
		item.Properties = append(item.Properties, c)
	}

	return item
}

func makeSimpleType(name xml.Name, from *xsd.SimpleType, logger *logrus.Logger) Type {
	return nil
}

func makeXsdBuiltinType(from xsd.Builtin, knownTypes *TypeList) Type {
	if from == xsd.String {
		t, _ := knownTypes.Find("string")
		return t
	}
	return nil
}

func isArray(from *xsd.ComplexType) bool {
	if max := getAttrValue(from.Attributes, "maxOccurs"); max != nil {
		return true
	}
	return false
}

func getAttrValue(attrs []xsd.Attribute, which string) *string {
	for _, a := range attrs {
		if a.Name.Local == which || fmt.Sprintf("%s:%s", a.Name.Space, a.Name.Local) == which {
			return &a.Attr[0].Value
		}
	}
	return nil
}
