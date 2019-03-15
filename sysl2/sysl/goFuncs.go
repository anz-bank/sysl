package main

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type goFunc struct {
	val  reflect.Value
	args []*sysl.Type
	ret  []*sysl.Type
}

type goFuncMap map[string]goFunc

var kindToPrimitiveType = map[reflect.Kind]sysl.Type_Primitive{
	reflect.Bool:   sysl.Type_BOOL,
	reflect.Int:    sysl.Type_INT,
	reflect.Int16:  sysl.Type_INT,
	reflect.Int32:  sysl.Type_INT,
	reflect.Int64:  sysl.Type_INT,
	reflect.String: sysl.Type_STRING,
}

var valueTypeToPrimitiveType = map[valueType]sysl.Type_Primitive{
	VALUE_BOOL:   sysl.Type_BOOL,
	VALUE_INT:    sysl.Type_INT,
	VALUE_FLOAT:  sysl.Type_FLOAT,
	VALUE_STRING: sysl.Type_STRING,
}

var stringType = &sysl.Type{
	Type: &sysl.Type_Primitive_{
		Primitive: sysl.Type_STRING,
	},
}

var intType = &sysl.Type{
	Type: &sysl.Type_Primitive_{
		Primitive: sysl.Type_INT,
	},
}

var listStringType = &sysl.Type{
	Type: &sysl.Type_List_{
		List: &sysl.Type_List{
			Type: &sysl.Type{
				Type: &sysl.Type_Primitive_{
					Primitive: sysl.Type_STRING,
				},
			},
		},
	},
}

var boolType = &sysl.Type{
	Type: &sysl.Type_Primitive_{
		Primitive: sysl.Type_BOOL,
	},
}

// GoFuncMap is Map of Function Names to goFunc{}
var GoFuncMap = map[string]goFunc{
	"Contains":      {reflect.ValueOf(strings.Contains), []*sysl.Type{stringType, stringType}, []*sysl.Type{boolType}},
	"Count":         {reflect.ValueOf(strings.Count), []*sysl.Type{stringType, stringType}, []*sysl.Type{intType}},
	"Fields":        {reflect.ValueOf(strings.Fields), []*sysl.Type{stringType}, []*sysl.Type{listStringType}},
	"FindAllString": {reflect.ValueOf(FindAllString), []*sysl.Type{stringType, stringType, intType}, []*sysl.Type{listStringType}},
	"HasPrefix":     {reflect.ValueOf(strings.HasPrefix), []*sysl.Type{stringType, stringType}, []*sysl.Type{boolType}},
	"HasSuffix":     {reflect.ValueOf(strings.HasSuffix), []*sysl.Type{stringType, stringType}, []*sysl.Type{boolType}},
	"Join":          {reflect.ValueOf(strings.Join), []*sysl.Type{listStringType, stringType}, []*sysl.Type{stringType}},
	"LastIndex":     {reflect.ValueOf(strings.LastIndex), []*sysl.Type{stringType, stringType}, []*sysl.Type{intType}},
	"MatchString":   {reflect.ValueOf(MatchString), []*sysl.Type{stringType, stringType}, []*sysl.Type{boolType}},
	"Replace":       {reflect.ValueOf(strings.Replace), []*sysl.Type{stringType, stringType, stringType, intType}, []*sysl.Type{stringType}},
	"Split":         {reflect.ValueOf(strings.Split), []*sysl.Type{stringType, stringType}, []*sysl.Type{listStringType}},
	"Title":         {reflect.ValueOf(strings.Title), []*sysl.Type{stringType}, []*sysl.Type{stringType}},
	"ToLower":       {reflect.ValueOf(strings.ToLower), []*sysl.Type{stringType}, []*sysl.Type{stringType}},
	"ToTitle":       {reflect.ValueOf(strings.ToTitle), []*sysl.Type{stringType}, []*sysl.Type{stringType}},
	"ToUpper":       {reflect.ValueOf(strings.ToUpper), []*sysl.Type{stringType}, []*sysl.Type{stringType}},
	"Trim":          {reflect.ValueOf(strings.Trim), []*sysl.Type{stringType, stringType}, []*sysl.Type{stringType}},
	"TrimLeft":      {reflect.ValueOf(strings.TrimLeft), []*sysl.Type{stringType, stringType}, []*sysl.Type{stringType}},
	"TrimPrefix":    {reflect.ValueOf(strings.TrimPrefix), []*sysl.Type{stringType, stringType}, []*sysl.Type{stringType}},
	"TrimRight":     {reflect.ValueOf(strings.TrimRight), []*sysl.Type{stringType, stringType}, []*sysl.Type{stringType}},
	"TrimSpace":     {reflect.ValueOf(strings.TrimSpace), []*sysl.Type{stringType}, []*sysl.Type{stringType}},
	"TrimSuffix":    {reflect.ValueOf(strings.TrimSuffix), []*sysl.Type{stringType, stringType}, []*sysl.Type{stringType}},
}

func valueToReflectValue(v *sysl.Value, t *sysl.Type) reflect.Value {
	switch x := t.Type.(type) {
	case *sysl.Type_Primitive_:
		switch x.Primitive {
		case sysl.Type_BOOL:
			return reflect.ValueOf(v.GetB())
		case sysl.Type_INT:
			return reflect.ValueOf(int(v.GetI()))
		case sysl.Type_STRING:
			return reflect.ValueOf(v.GetS())
		}
	case *sysl.Type_List_:
		listOf := x.List.Type
		var valueList []*sysl.Value
		if v.GetList() != nil {
			valueList = v.GetList().Value
		} else {
			valueList = v.GetSet().Value
		}

		if p := listOf.GetPrimitive(); p != sysl.Type_NO_Primitive {
			switch p {
			// Uncomment if required in future.
			// case sysl.Type_INT:
			// 	intSlice := []int{}
			// 	for _, listItem := range valueList {
			// 		intSlice = append(intSlice, int(listItem.GetI()))
			// 	}
			// 	return reflect.ValueOf(intSlice)
			case sysl.Type_STRING:
				stringSlice := []string{}
				for _, listItem := range valueList {
					stringSlice = append(stringSlice, listItem.GetS())
				}
				return reflect.ValueOf(stringSlice)
			}
		}
	}
	panic(errors.Errorf("valueToReflectValue: unsupported value type: %v", v))
}

func isValueExpectedType(v *sysl.Value, t *sysl.Type) bool {
	vType := getValueType(v)
	type1, has := valueTypeToPrimitiveType[vType]
	inType := t.GetPrimitive()
	// if both are primitive types
	if has && inType != sysl.Type_NO_Primitive {
		return inType == type1
	}

	if vType == VALUE_LIST && len(v.GetList().Value) == 0 {
		return true
	}

	if vType == VALUE_SET && len(v.GetSet().Value) == 0 {
		return true
	}

	if vType == VALUE_LIST && t.GetList() != nil {
		return isValueExpectedType(v.GetList().Value[0], t.GetList().Type)
	}

	if vType == VALUE_SET && t.GetSet() != nil {
		return isValueExpectedType(v.GetSet().Value[0], t.GetSet())
	}

	// Pass elements of sets as a slice to Go funcs.
	if vType == VALUE_SET && t.GetList() != nil {
		return isValueExpectedType(v.GetSet().Value[0], t.GetList().Type)
	}
	return false
}

func isReflectValueExpectedType(r reflect.Value, typ *sysl.Type) bool {
	kind := r.Kind()
	_, has := kindToPrimitiveType[kind]
	if has {
		return true
	}
	if kind == reflect.Slice {
		if typ.GetList() == nil || typ.GetSet() != nil {
			return false
		}
		_, has = kindToPrimitiveType[r.Index(0).Kind()]
		return has
	}
	return false
}

func reflectToValue(r reflect.Value, typ *sysl.Type) *sysl.Value {
	if isReflectValueExpectedType(r, typ) == false {
		logrus.Warnf("Got %s, Expected Value type: %v \n", r.Kind(), typ.Type)
	}

	switch r.Kind() {
	case reflect.Bool:
		return MakeValueBool(r.Bool())
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		if typ.GetPrimitive() != sysl.Type_INT {
			logrus.Warnf("Got int, Expected Value type: %v \n", typ.Type)
		}
		return MakeValueI64(r.Int())
	case reflect.String:
		if typ.GetPrimitive() != sysl.Type_STRING {
			logrus.Warnf("Got string, Expected Value type: %v \n", typ.Type)
		}
		return MakeValueString(r.String())
	case reflect.Slice:
		list := MakeValueList()
		primitiveType := typ.GetList().GetType()

		for i := 0; i < r.Len(); i++ {
			appendItemToValueList(list.GetList(), reflectToValue(r.Index(i), primitiveType))
		}
		return list
	}
	panic(errors.Errorf("reflectToValue: kind %s not supported\n", r.Kind().String()))
}

func evalGoFunc(name string, list *sysl.Value) *sysl.Value {
	if f, has := GoFuncMap[name]; has {
		var in []reflect.Value
		if len(list.GetList().Value) != len(f.args) {
			logrus.Errorf("Incorrect number of arg function %s\n", name)
			logrus.Errorf("Expected number of arg: (%d)\n", len(f.args))
			logrus.Errorf("Got num args: (%d)\n", len(list.GetList().Value))
			return nil
		}

		for i, l := range list.GetList().Value {
			if isValueExpectedType(l, f.args[i]) {
				in = append(in, valueToReflectValue(l, f.args[i]))
			} else {
				return nil
			}
		}
		result := f.val.Call(in)
		if len(f.ret) != len(result) {
			logrus.Errorf("Incorrect number of return types for function %s\n", name)
			logrus.Errorf("Expected number of return types: (%d)\n", len(f.ret))
			logrus.Errorf("Got result: (%d)\n", len(result))
		}
		return reflectToValue(result[0], f.ret[0])
	}
	return nil
}

// MatchString exposes regexp.MatchString to sysl transforms
func MatchString(pattern, word string) bool {
	return regexp.MustCompile(pattern).MatchString(word)

}

// FindAllString exposes regexp.FindAllString to sysl transforms
func FindAllString(pattern, word string, n int) []string {
	return regexp.MustCompile(pattern).FindAllString(word, n)
}
