package main

import (
	"github.com/anz-bank/sysl/src/proto"
)

// Scope ...
type Scope map[string]*sysl.Value

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

func fieldsToValueMap(fields map[string]*sysl.Type) *sysl.Value {
	var typeName string
	fieldMap := MakeValueMap()

	for key, t := range fields {
		m := MakeValueMap()
		addItemToValueMap(fieldMap.GetMap(), key, m)
		addItemToValueMap(m.GetMap(), "name", MakeValueString(key))
		addItemToValueMap(m.GetMap(), "docstring", MakeValueString(t.Docstring))
		switch x := t.Type.(type) {
		case *sysl.Type_Primitive_:
			typeName = "primitive"
			addItemToValueMap(m.GetMap(), "primitive", MakeValueString(sysl.Type_Primitive_name[int32(x.Primitive)]))
		case *sysl.Type_TypeRef:
			typeName = "type_ref"
		}
		addItemToValueMap(m.GetMap(), "type", MakeValueString(typeName))
	}
	return fieldMap
}

func typeToValue(t *sysl.Type) *sysl.Value {
	m := MakeValueMap()
	addItemToValueMap(m.GetMap(), "docstring", MakeValueString(t.Docstring))
	addItemToValueMap(m.GetMap(), "attrs", attrsToValueMap(t.Attrs))
	var typeName string
	var attrDefs map[string]*sysl.Type
	switch x := t.Type.(type) {
	case *sysl.Type_Tuple_:
		typeName = "tuple"
		attrDefs = x.Tuple.AttrDefs
	case *sysl.Type_Relation_:
		typeName = "relation"
		attrDefs = x.Relation.AttrDefs
	}
	addItemToValueMap(m.GetMap(), "fields", fieldsToValueMap(attrDefs))
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
		addItemToValueMap(m.GetMap(), key, typeToValue(t))
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
}
