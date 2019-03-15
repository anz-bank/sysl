package main

import (
	"github.com/anz-bank/sysl/src/proto"
)

// Scope holds the value of the variables during the execution of a transform
type Scope map[string]*sysl.Value

// MakeValueBool returns sysl.Value of type Value_B (bool)
func MakeValueBool(val bool) *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_B{
			B: val,
		},
	}
}

// MakeValueI64 returns sysl.Value of type Value_I (int64)
func MakeValueI64(val int64) *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_I{
			I: val,
		},
	}
}

// MakeValueString returns sysl.Value of type Value_S (string)
func MakeValueString(val string) *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_S{
			S: val,
		},
	}
}

// MakeValueMap returns sysl.Value of type Value_Map_ (map[string]*sysl.Value)
func MakeValueMap() *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_Map_{
			Map: &sysl.Value_Map{
				Items: map[string]*sysl.Value{},
			},
		},
	}
}

// MakeValueList returns sysl.Value of type Value_List_ ([]*sysl.Value)
func MakeValueList() *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_List_{
			List: &sysl.Value_List{
				Value: []*sysl.Value{},
			},
		},
	}
}

// MakeValueSet returns sysl.Value of type Value_Set ([]*sysl.Value)
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
	case *sysl.Value_List_, *sysl.Value_Set:
		return true
	}
	return false
}

// expects m.Value to be of type Value_Map
func addItemToValueMap(m *sysl.Value, name string, val *sysl.Value) {
	m.GetMap().Items[name] = val
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
		if x.TypeRef.Ref != nil && len(x.TypeRef.Ref.Path) == 1 {
			typeDetail = x.TypeRef.Ref.Path[0]
		} else {
			typeDetail = x.TypeRef.Ref.Appname.Part[0]
		}

	case *sysl.Type_Sequence:
		typeName = "sequence"
		_, d := getTypeDetail(x.Sequence)
		typeDetail = d
	case *sysl.Type_Set:
		typeName = "set"
		_, d := getTypeDetail(x.Set)
		typeDetail = d
	case *sysl.Type_List_:
		typeName = "list"
		_, d := getTypeDetail(x.List.Type)
		typeDetail = d
	case *sysl.Type_Tuple_:
		typeName = "tuple"
	case *sysl.Type_Relation_:
		typeName = "relation"
	case *sysl.Type_OneOf_:
		typeName = "union"
	}
	return typeName, typeDetail
}

func fieldsToValueMap(fields map[string]*sysl.Type) *sysl.Value {
	fieldMap := MakeValueMap()

	for key, t := range fields {
		m := typeToValue(t)
		addItemToValueMap(m, "name", MakeValueString(key))
		addItemToValueMap(fieldMap, key, m)
	}
	return fieldMap
}

func typeToValue(t *sysl.Type) *sysl.Value {
	m := MakeValueMap()
	addItemToValueMap(m, "docstring", MakeValueString(t.Docstring))
	addItemToValueMap(m, "attrs", attrsToValueMap(t.Attrs))
	typeName, typeDetail := getTypeDetail(t)
	addItemToValueMap(m, "type", MakeValueString(typeName))
	addItemToValueMap(m, typeName, MakeValueString(typeDetail))
	addItemToValueMap(m, "optional", MakeValueBool(t.Opt))

	switch x := t.Type.(type) {
	case *sysl.Type_Tuple_:
		addItemToValueMap(m, "fields", fieldsToValueMap(x.Tuple.AttrDefs))
	case *sysl.Type_Relation_:
		addItemToValueMap(m, "fields", fieldsToValueMap(x.Relation.AttrDefs))
	case *sysl.Type_OneOf_:
		unionSet := MakeValueSet()
		for _, embeddedType := range x.OneOf.Type {
			appendItemToValueList(unionSet.GetSet(), MakeValueString(embeddedType.GetTypeRef().Ref.Path[0]))
		}
		addItemToValueMap(m, "fields", unionSet)
	case *sysl.Type_Sequence:
		seqMap := MakeValueMap()
		addItemToValueMap(m, typeName, seqMap)
		typeName, typeDetail = getTypeDetail(t.GetSequence())
		addItemToValueMap(seqMap, "type", MakeValueString(typeName))
		addItemToValueMap(seqMap, typeName, MakeValueString(typeDetail))
		addItemToValueMap(seqMap, "optional", MakeValueBool(false))
	}
	return m
}

func attrsToValueMap(attrs map[string]*sysl.Attribute) *sysl.Value {
	m := MakeValueMap()
	for key, attr := range attrs {
		addItemToValueMap(m, key, attributeToValue(attr))
	}
	return m
}

func typesToValueMap(types map[string]*sysl.Type) *sysl.Value {
	m := MakeValueMap()
	for key, t := range types {
		switch t.Type.(type) {
		case *sysl.Type_Tuple_, *sysl.Type_Relation_:
			addItemToValueMap(m, key, typeToValue(t))
		}
	}
	return m
}

func unionToValueMap(types map[string]*sysl.Type) *sysl.Value {
	m := MakeValueMap()
	for key, t := range types {
		switch t.Type.(type) {
		case *sysl.Type_OneOf_:
			addItemToValueMap(m, key, typeToValue(t))
		}
	}
	return m
}

func aliasToValueMap(types map[string]*sysl.Type) *sysl.Value {
	m := MakeValueMap()
	for key, t := range types {
		switch t.Type.(type) {
		case *sysl.Type_OneOf_, *sysl.Type_Tuple_, *sysl.Type_Relation_:
		default:
			addItemToValueMap(m, key, typeToValue(t))
		}
	}
	return m
}

func queryParamsToValue(qp *sysl.Endpoint_RestParams_QueryParam) *sysl.Value {
	m := MakeValueMap()
	addItemToValueMap(m, "name", MakeValueString(qp.Name))
	typeName, typeDetail := getTypeDetail(qp.Type)
	addItemToValueMap(m, "type", MakeValueString(typeName))
	addItemToValueMap(m, "optional", MakeValueBool(qp.Type.GetOpt()))
	addItemToValueMap(m, typeName, MakeValueString(typeDetail))
	return m
}

func paramToValue(qp *sysl.Param) *sysl.Value {
	m := MakeValueMap()
	addItemToValueMap(m, "name", MakeValueString(qp.Name))
	typeName, typeDetail := getTypeDetail(qp.Type)
	addItemToValueMap(m, "type", MakeValueString(typeName))
	addItemToValueMap(m, typeName, MakeValueString(typeDetail))
	addItemToValueMap(m, "attrs", attrsToValueMap(qp.Type.Attrs))
	addItemToValueMap(m, "optional", MakeValueBool(qp.Type.GetOpt()))
	return m
}

func stmtToValue(s *sysl.Statement) *sysl.Value {
	m := MakeValueMap()
	var stmt_type string
	switch x := s.Stmt.(type) {
	case *sysl.Statement_Action:
		stmt_type = "action"
		addItemToValueMap(m, "action", MakeValueString(x.Action.Action))
	case *sysl.Statement_Call:
		stmt_type = "call"
		addItemToValueMap(m, "endpoint", MakeValueString(x.Call.Endpoint))
		addItemToValueMap(m, "target", MakeValueString(getAppName(x.Call.Target)))
	case *sysl.Statement_Ret:
		stmt_type = "return"
		addItemToValueMap(m, "payload", MakeValueString(x.Ret.Payload))
	}
	addItemToValueMap(m, "type", MakeValueString(stmt_type))
	return m
}

func endpointToValue(e *sysl.Endpoint) *sysl.Value {
	m := MakeValueMap()
	addItemToValueMap(m, "name", MakeValueString(e.Name))
	addItemToValueMap(m, "longname", MakeValueString(e.LongName))
	addItemToValueMap(m, "docstring", MakeValueString(e.Docstring))
	addItemToValueMap(m, "attrs", attrsToValueMap(e.Attrs))
	addItemToValueMap(m, "is_rest", MakeValueBool(e.RestParams != nil))
	addItemToValueMap(m, "is_pubsub", MakeValueBool(e.IsPubsub))

	if e.RestParams != nil {
		addItemToValueMap(m, "method", MakeValueString(sysl.Endpoint_RestParams_Method_name[int32(e.RestParams.Method)]))
		addItemToValueMap(m, "path", MakeValueString(e.RestParams.Path))

		queryvars := MakeValueList()
		for _, query_param := range e.RestParams.QueryParam {
			appendItemToValueList(queryvars.GetList(), queryParamsToValue(query_param))
		}
		addItemToValueMap(m, "queryvars", queryvars)

		pathvars := MakeValueList()
		for _, query_param := range e.RestParams.UrlParam {
			appendItemToValueList(pathvars.GetList(), queryParamsToValue(query_param))
		}
		addItemToValueMap(m, "pathvars", pathvars)
	}

	params := MakeValueList()
	for _, param := range e.Param {
		appendItemToValueList(params.GetList(), paramToValue(param))
	}
	addItemToValueMap(m, "params", params)

	stmtsList := MakeValueList()
	for _, stmt := range e.Stmt {
		if stmt.GetRet() != nil {
			addItemToValueMap(m, "ret", stmtToValue(stmt))
		} else {
			appendItemToValueList(stmtsList.GetList(), stmtToValue(stmt))
		}
	}
	addItemToValueMap(m, "stmts", stmtsList)

	return m
}

func endpointsToValueMap(endpoints map[string]*sysl.Endpoint) *sysl.Value {
	m := MakeValueMap()
	for key, e := range endpoints {
		addItemToValueMap(m, key, endpointToValue(e))
	}
	return m
}

// AddInt add int to scope
func (s Scope) AddInt(name string, val int64) {
	s[name] = MakeValueI64(val)
}

// AddString add string to scope
func (s Scope) AddString(name string, val string) {
	s[name] = MakeValueString(val)
}

// AddApp add sysl.App to scope
func (s Scope) AddApp(name string, app *sysl.Application) {
	m := MakeValueMap()
	s[name] = m
	addItemToValueMap(m, "name", MakeValueString(getAppName(app.Name)))
	addItemToValueMap(m, "attrs", attrsToValueMap(app.Attrs))
	addItemToValueMap(m, "types", typesToValueMap(app.Types))
	addItemToValueMap(m, "union", unionToValueMap(app.Types))
	addItemToValueMap(m, "alias", aliasToValueMap(app.Types))
	addItemToValueMap(m, "endpoints", endpointsToValueMap(app.Endpoints))
}
