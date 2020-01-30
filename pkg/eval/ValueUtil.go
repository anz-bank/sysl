package eval

import (
	"strings"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	log "github.com/sirupsen/logrus"
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
func MakeValueList(values ...*sysl.Value) *sysl.Value {
	if values == nil {
		values = []*sysl.Value{}
	}
	return &sysl.Value{
		Value: &sysl.Value_List_{
			List: &sysl.Value_List{
				Value: values,
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

// GetValueSlice returns the slice from either a set or a list
func GetValueSlice(obj *sysl.Value) []*sysl.Value {
	switch obj.Value.(type) {
	case *sysl.Value_List_:
		return obj.GetList().Value
	case *sysl.Value_Set:
		return obj.GetSet().Value
	default:
		panic("GetValueSlice expecting a collection type")
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
func AddItemToValueMap(m *sysl.Value, name string, val *sysl.Value) {
	m.GetMap().Items[name] = val
}

func AppendItemToValueList(m *sysl.Value_List, val *sysl.Value) {
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
	fieldMap := MakeValueMap()

	for key, t := range fields {
		m := TypeToValue(t)
		AddItemToValueMap(m, "name", MakeValueString(key))
		AddItemToValueMap(fieldMap, key, m)
	}
	return fieldMap
}

func TypeToValue(t *sysl.Type) *sysl.Value {
	m := MakeValueMap()
	AddItemToValueMap(m, "docstring", MakeValueString(t.Docstring))
	AddItemToValueMap(m, "attrs", attrsToValueMap(t.Attrs))
	typeName, typeDetail := syslutil.GetTypeDetail(t)
	AddItemToValueMap(m, "type", MakeValueString(typeName))
	AddItemToValueMap(m, typeName, MakeValueString(typeDetail))
	AddItemToValueMap(m, "optional", MakeValueBool(t.Opt))

	switch x := t.Type.(type) {
	case *sysl.Type_Tuple_:
		AddItemToValueMap(m, "fields", fieldsToValueMap(x.Tuple.AttrDefs))
	case *sysl.Type_Relation_:
		AddItemToValueMap(m, "fields", fieldsToValueMap(x.Relation.AttrDefs))
	case *sysl.Type_OneOf_:
		unionSet := MakeValueSet()
		for _, embeddedType := range x.OneOf.Type {
			AppendItemToValueList(unionSet.GetSet(), MakeValueString(embeddedType.GetTypeRef().Ref.Path[0]))
		}
		AddItemToValueMap(m, "fields", unionSet)
	case *sysl.Type_Sequence:
		seqMap := MakeValueMap()
		AddItemToValueMap(m, typeName, seqMap)
		typeName, typeDetail = syslutil.GetTypeDetail(t.GetSequence())
		AddItemToValueMap(seqMap, "type", MakeValueString(typeName))
		AddItemToValueMap(seqMap, typeName, MakeValueString(typeDetail))
		AddItemToValueMap(seqMap, "optional", MakeValueBool(false))
	}
	return m
}

func attrsToValueMap(attrs map[string]*sysl.Attribute) *sysl.Value {
	m := MakeValueMap()
	for key, attr := range attrs {
		AddItemToValueMap(m, key, attributeToValue(attr))
	}
	return m
}

func typesToValueMap(types map[string]*sysl.Type) *sysl.Value {
	m := MakeValueMap()
	for key, t := range types {
		switch t.Type.(type) {
		case *sysl.Type_Tuple_, *sysl.Type_Relation_:
			AddItemToValueMap(m, key, TypeToValue(t))
		}
	}
	return m
}

func unionToValueMap(types map[string]*sysl.Type) *sysl.Value {
	m := MakeValueMap()
	for key, t := range types {
		if _, ok := t.Type.(*sysl.Type_OneOf_); ok {
			AddItemToValueMap(m, key, TypeToValue(t))
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
			AddItemToValueMap(m, key, TypeToValue(t))
		}
	}
	return m
}

func queryParamsToValue(qp *sysl.Endpoint_RestParams_QueryParam) *sysl.Value {
	m := MakeValueMap()
	AddItemToValueMap(m, "name", MakeValueString(qp.Name))
	typeName, typeDetail := syslutil.GetTypeDetail(qp.Type)
	AddItemToValueMap(m, "type", MakeValueString(typeName))
	AddItemToValueMap(m, "optional", MakeValueBool(qp.Type.GetOpt()))
	AddItemToValueMap(m, typeName, MakeValueString(typeDetail))
	return m
}

func paramToValue(qp *sysl.Param) *sysl.Value {
	m := MakeValueMap()
	AddItemToValueMap(m, "name", MakeValueString(qp.Name))
	typeName, typeDetail := syslutil.GetTypeDetail(qp.Type)
	AddItemToValueMap(m, "type", MakeValueString(typeName))
	AddItemToValueMap(m, typeName, MakeValueString(typeDetail))
	AddItemToValueMap(m, "attrs", attrsToValueMap(qp.Type.Attrs))
	AddItemToValueMap(m, "optional", MakeValueBool(qp.Type.GetOpt()))
	return m
}

func stmtToValue(s *sysl.Statement) *sysl.Value {
	m := MakeValueMap()
	var stmtType string
	switch x := s.Stmt.(type) {
	case *sysl.Statement_Action:
		stmtType = "action"
		AddItemToValueMap(m, "action", MakeValueString(x.Action.Action))
	case *sysl.Statement_Call:
		stmtType = "call"
		AddItemToValueMap(m, "endpoint", MakeValueString(x.Call.Endpoint))
		AddItemToValueMap(m, "target", MakeValueString(syslutil.GetAppName(x.Call.Target)))
	case *sysl.Statement_Ret:
		stmtType = "return"
		AddItemToValueMap(m, "payload", MakeValueString(x.Ret.Payload))
	}
	AddItemToValueMap(m, "type", MakeValueString(stmtType))
	return m
}

func endpointToValue(e *sysl.Endpoint) *sysl.Value {
	m := MakeValueMap()
	AddItemToValueMap(m, "name", MakeValueString(e.Name))
	AddItemToValueMap(m, "longname", MakeValueString(e.LongName))
	AddItemToValueMap(m, "docstring", MakeValueString(e.Docstring))
	AddItemToValueMap(m, "attrs", attrsToValueMap(e.Attrs))
	AddItemToValueMap(m, "is_rest", MakeValueBool(e.RestParams != nil))
	AddItemToValueMap(m, "is_pubsub", MakeValueBool(e.IsPubsub))
	retTypes := MakeValueMap()
	var retValues []string

	if e.RestParams != nil {
		AddItemToValueMap(m, "method", MakeValueString(sysl.Endpoint_RestParams_Method_name[int32(e.RestParams.Method)]))
		AddItemToValueMap(m, "path", MakeValueString(e.RestParams.Path))

		queryvars := MakeValueList()
		for _, queryParam := range e.RestParams.QueryParam {
			AppendItemToValueList(queryvars.GetList(), queryParamsToValue(queryParam))
		}
		AddItemToValueMap(m, "queryvars", queryvars)

		pathvars := MakeValueList()
		for _, queryParam := range e.RestParams.UrlParam {
			AppendItemToValueList(pathvars.GetList(), queryParamsToValue(queryParam))
		}
		AddItemToValueMap(m, "pathvars", pathvars)
	}

	params := MakeValueList()
	for _, param := range e.Param {
		if param.GetType() == nil {
			log.Warnf("empty param defined in endpoint: %s", e.Name)
			continue
		}
		AppendItemToValueList(params.GetList(), paramToValue(param))
	}
	AddItemToValueMap(m, "params", params)

	stmtsList := MakeValueList()
	for _, stmt := range e.Stmt {
		switch s := stmt.Stmt.(type) {
		case *sysl.Statement_Ret:
			retValues = strings.Split(s.Ret.GetPayload(), " <: ")
		case *sysl.Statement_Cond:
			retValues = strings.Split(s.Cond.GetStmt()[0].GetRet().GetPayload(), " <: ")
		case *sysl.Statement_Group:
			if s.Group.GetTitle() == "else" || strings.Contains(s.Group.GetTitle(), "else if") {
				retValues = strings.Split(s.Group.GetStmt()[0].GetRet().GetPayload(), " <: ")
			} else {
				log.Warnf("Unexpected statement %s found", s.Group.GetTitle())
				retValues = []string{stmtToValue(s.Group.Stmt[0]).GetS()}
			}
		default:
			AppendItemToValueList(stmtsList.GetList(), stmtToValue(stmt))
			continue
		}
		if len(retValues) > 1 {
			retTypes.GetMap().Items[retValues[0]] = MakeValueString(retValues[1])
		} else if len(retValues) == 1 {
			retTypes.GetMap().Items["payload"] = MakeValueString(retValues[0])
		}
	}
	AddItemToValueMap(m, "ret", retTypes)
	AddItemToValueMap(m, "stmts", stmtsList)

	return m
}

func endpointsToValueMap(endpoints map[string]*sysl.Endpoint) *sysl.Value {
	m := MakeValueMap()
	for key, e := range endpoints {
		AddItemToValueMap(m, key, endpointToValue(e))
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
	s[name] = addAppToValueMap(app)
}

// AddModule add module <: sysl.Module to scope
func (s Scope) AddModule(name string, module *sysl.Module) {
	apps := MakeValueMap()
	types := MakeValueMap()

	for n, a := range module.GetApps() {
		AddItemToValueMap(apps, n, addAppToValueMap(a))
	}

	m := MakeValueMap()
	AddItemToValueMap(m, "apps", apps)
	AddItemToValueMap(m, "types", types)

	s[name] = m
}

func (s Scope) ToValue() *sysl.Value {
	return &sysl.Value{
		Value: &sysl.Value_Map_{
			Map: &sysl.Value_Map{
				Items: s,
			},
		},
	}
}

func addAppToValueMap(app *sysl.Application) *sysl.Value {
	m := MakeValueMap()
	AddItemToValueMap(m, "name", MakeValueString(syslutil.GetAppName(app.Name)))
	AddItemToValueMap(m, "attrs", attrsToValueMap(app.Attrs))
	AddItemToValueMap(m, "types", typesToValueMap(app.Types))
	AddItemToValueMap(m, "union", unionToValueMap(app.Types))
	AddItemToValueMap(m, "alias", aliasToValueMap(app.Types))
	AddItemToValueMap(m, "endpoints", endpointsToValueMap(app.Endpoints))

	return m
}
