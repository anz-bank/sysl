package main

import (
	"github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
)

type valueType int

// const definitions for various sysl.Value types
const (
	VALUE_NO_ARG valueType = -1
	VALUE_BOOL   valueType = iota
	VALUE_INT
	VALUE_FLOAT
	VALUE_STRING
	VALUE_STRING_DECIMAL
	VALUE_LIST
	VALUE_MAP
	VALUE_SET
	VALUE_NULL
)

var valueTypeNames = map[valueType]string{
	VALUE_NO_ARG:         "VALUE_NO_ARG",
	VALUE_BOOL:           "VALUE_BOOL",
	VALUE_INT:            "VALUE_INT",
	VALUE_FLOAT:          "VALUE_FLOAT",
	VALUE_STRING:         "VALUE_STRING",
	VALUE_STRING_DECIMAL: "VALUE_STRING_DECIMAL",
	VALUE_LIST:           "VALUE_LIST",
	VALUE_MAP:            "VALUE_MAP",
	VALUE_SET:            "VALUE_SET",
	VALUE_NULL:           "VALUE_NULL",
}

func (v valueType) String() string {
	return valueTypeNames[v]
}

func getValueType(v *sysl.Value) valueType {
	if v == nil {
		return VALUE_NO_ARG
	}
	switch v.Value.(type) {
	case *sysl.Value_B:
		return VALUE_BOOL
	case *sysl.Value_I:
		return VALUE_INT
	case *sysl.Value_D:
		return VALUE_FLOAT
	case *sysl.Value_S:
		return VALUE_STRING
	case *sysl.Value_Decimal:
		return VALUE_STRING_DECIMAL
	case *sysl.Value_Set:
		return VALUE_SET
	case *sysl.Value_List_:
		return VALUE_LIST
	case *sysl.Value_Map_:
		return VALUE_MAP
	case *sysl.Value_Null_:
		return VALUE_NULL
	default:
		panic(errors.Errorf("exprOp: getValueType: unhandled type: %v", v))
	}
}

func getContainedType(container *sysl.Value) valueType {
	var list []*sysl.Value
	switch x := container.Value.(type) {
	case *sysl.Value_List_:
		list = x.List.Value
	case *sysl.Value_Set:
		list = x.Set.Value
	default:
		return VALUE_NO_ARG
	}

	if len(list) == 0 {
		return VALUE_NO_ARG
	}
	return getValueType(list[0])
}
