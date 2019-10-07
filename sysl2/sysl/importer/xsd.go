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

	var types TypeList
	for _, schema := range specs {
		if schema.TargetNS == targetNS {
			types = append(types, loadSchemaTypes(schema, logger)...)

			self, ok := types.Find("_self")

			fmt.Printf("%v %v", self, ok)
		}
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

func loadSchemaTypes(schema xsd.Schema, logger *logrus.Logger) TypeList {
	defs := TypeList{}

	for name, data := range schema.Types {
		if t := FindType(data, &defs); t == nil {
			defs = append(defs, makeType(name, data, &defs, logger))
		}
	}

	return defs
}

func FindType(t xsd.Type, knownTypes *TypeList) Type {
	res, found := knownTypes.Find(xsd.XMLName(t).Local)
	if found {
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
		return nil
	}

	item := &StandardType{
		name: from.Name.Local,
	}
	for _, child := range from.Elements {
		childType := FindType(child.Type, knownTypes)
		if childType == nil {
			childType = makeType(child.Name, child.Type, knownTypes, logger)
			*knownTypes = append(*knownTypes, childType)
		}
		f := Field{
			Name: child.Name.Local,
			Type: childType,
		}
		item.Properties = append(item.Properties, f)
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
