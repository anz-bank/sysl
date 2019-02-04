package main

import (
	"github.com/anz-bank/sysl/src/proto"
)

// Scope ...
type Scope map[string]*sysl.Value

// MakeValueI64 ...
func MakeValueBool(val bool) *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_B{
			B: val,
		},
	}
}

// MakeValueI64 ...
func MakeValueI64(val int64) *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_I{
			I: val,
		},
	}
}

// MakeValueString ...
func MakeValueString(val string) *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_S{
			S: val,
		},
	}
}

// MakeValueMap ...
func MakeValueMap() *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_Map_{
			Map: &sysl.Value_Map{
				Items: map[string]*sysl.Value{},
			},
		},
	}
}

// MakeValueList ...
func MakeValueList() *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_List_{
			List: &sysl.Value_List{
				Value: []*sysl.Value{},
			},
		},
	}
}

// MakeValueSet ...
func MakeValueSet() *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_Set{
			Set: &sysl.Value_List{
				Value: []*sysl.Value{},
			},
		},
	}
}

// IsCollectionType does obj holds multiple instances of the same thing?
func IsCollectionType(obj *sysl.Value) bool {
	switch obj.Value.(type) {
	case *sysl.Value_List_:
		return true
	case *sysl.Value_Set:
		return true
	}
	return false
}

func addItemToValueMap(m *sysl.Value_Map, name string, val *sysl.Value) {
	m.Items[name] = val
}

func appendListToValueList(m *sysl.Value_List, val *sysl.Value_List) {
	for _, v := range val.Value {
		appendItemToValueList(m, v)
	}
}

func appendItemToValueList(m *sysl.Value_List, val *sysl.Value) {
	m.Value = append(m.Value, val)
}

func attributeToValue(attr *sysl.Attribute) *sysl.Value {
	switch x := attr.Attribute.(type) {
	case *sysl.Attribute_S:
		return MakeValueString(x.S)
	case *sysl.Attribute_A:
		l := MakeValueList()
		arr := []*sysl.Value{}
		for _, elt := range x.A.Elt {
			arr = append(arr, attributeToValue(elt))
		}
		l.GetList().Value = arr
		return l
	}
	return nil
}

func getTypeDetail(t *sysl.Type) (string, string) {
	var typeName, typeDetail string
	switch x := t.Type.(type) {
	case *sysl.Type_Primitive_:
		typeName = "primitive"
		typeDetail = sysl.Type_Primitive_name[int32(x.Primitive)]
	case *sysl.Type_TypeRef:
		typeName = "type_ref"
		typeDetail = x.TypeRef.Ref.Path[0]
	case *sysl.Type_Sequence:
		typeName = "sequence"
		_, d := getTypeDetail(x.Sequence)
		typeDetail = d
	case *sysl.Type_Tuple_:
		typeName = "tuple"
	case *sysl.Type_Relation_:
		typeName = "relation"
	}
	return typeName, typeDetail
}

func fieldsToValueMap(fields map[string]*sysl.Type) *sysl.Value {
	fieldMap := MakeValueMap()

	for key, t := range fields {
		m := MakeValueMap()
		typeName, typeDetail := getTypeDetail(t)
		addItemToValueMap(m.GetMap(), "type", MakeValueString(typeName))
		addItemToValueMap(m.GetMap(), typeName, MakeValueString(typeDetail))
		if typeName == "sequence" {
			seqMap := MakeValueMap()
			addItemToValueMap(m.GetMap(), typeName, seqMap)

			typeName, typeDetail = getTypeDetail(t.GetSequence())
			addItemToValueMap(seqMap.GetMap(), "type", MakeValueString(typeName))
			addItemToValueMap(seqMap.GetMap(), typeName, MakeValueString(typeDetail))
		}

		addItemToValueMap(fieldMap.GetMap(), key, m)
		addItemToValueMap(m.GetMap(), "name", MakeValueString(key))
		addItemToValueMap(m.GetMap(), "optional", MakeValueBool(t.Opt))
		addItemToValueMap(m.GetMap(), "docstring", MakeValueString(t.Docstring))
	}
	return fieldMap
}

func typeToValue(t *sysl.Type) *sysl.Value {
	m := MakeValueMap()
	addItemToValueMap(m.GetMap(), "docstring", MakeValueString(t.Docstring))
	addItemToValueMap(m.GetMap(), "attrs", attrsToValueMap(t.Attrs))
	var typeName string
	switch x := t.Type.(type) {
	case *sysl.Type_Tuple_:
		typeName = "tuple"
		addItemToValueMap(m.GetMap(), "fields", fieldsToValueMap(x.Tuple.AttrDefs))
	case *sysl.Type_Relation_:
		typeName = "relation"
		addItemToValueMap(m.GetMap(), "fields", fieldsToValueMap(x.Relation.AttrDefs))
	case *sysl.Type_OneOf_:
		unionSet := MakeValueSet()
		typeName = "union"
		for _, embeddedType := range x.OneOf.Type {
			appendItemToValueList(unionSet.GetSet(), MakeValueString(embeddedType.GetTypeRef().Ref.Path[0]))
		}
		addItemToValueMap(m.GetMap(), "fields", unionSet)
	default:
		panic("typeToValue: unsupported type")
	}
	addItemToValueMap(m.GetMap(), "type", MakeValueString(typeName))
	return m
}

func attrsToValueMap(attrs map[string]*sysl.Attribute) *sysl.Value {
	m := MakeValueMap()
	for key, attr := range attrs {
		addItemToValueMap(m.GetMap(), key, attributeToValue(attr))
	}
	return m
}

func typesToValueMap(types map[string]*sysl.Type) *sysl.Value {
	m := MakeValueMap()
	for key, t := range types {
		switch t.Type.(type) {
		case *sysl.Type_OneOf_:
			continue
		default:
			addItemToValueMap(m.GetMap(), key, typeToValue(t))
		}
	}
	return m
}

func unionToValueMap(types map[string]*sysl.Type) *sysl.Value {
	m := MakeValueMap()
	for key, t := range types {
		switch t.Type.(type) {
		case *sysl.Type_OneOf_:
			addItemToValueMap(m.GetMap(), key, typeToValue(t))
		}
	}
	return m
}

// AddInt add int to scope
func (s *Scope) AddInt(name string, val int64) {
	(*s)[name] = MakeValueI64(val)
}

// AddString add string to scope
func (s *Scope) AddString(name string, val string) {
	(*s)[name] = MakeValueString(val)
}

// AddApp add sysl.App to scope
func (s *Scope) AddApp(name string, app *sysl.Application) {
	m := MakeValueMap()
	(*s)[name] = m
	addItemToValueMap(m.GetMap(), "name", MakeValueString(getAppName(app.Name)))
	addItemToValueMap(m.GetMap(), "attrs", attrsToValueMap(app.Attrs))
	addItemToValueMap(m.GetMap(), "types", typesToValueMap(app.Types))
	addItemToValueMap(m.GetMap(), "union", unionToValueMap(app.Types))
}
