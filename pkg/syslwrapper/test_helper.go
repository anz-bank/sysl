package syslwrapper

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func MakeAppName(name string) *sysl.AppName {
	return &sysl.AppName{
		Part: strings.Split(name, " "),
	}
}

// Creates a type reference from one context to another.
// contextApp is the name of the app where the typeref is used
// refApp is the app where the
func MakeTypeRef(contextApp string, contextTypePath []string, refApp string, refTypePath []string) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_TypeRef{
			TypeRef: &sysl.ScopedRef{
				Context: &sysl.Scope{
					Appname: MakeAppName(contextApp),
					Path:    contextTypePath,
				},
				Ref: &sysl.Scope{
					Appname: MakeAppName(refApp),
					Path:    refTypePath,
				},
			},
		},
	}
}

func MakePrimitive(primType string) *sysl.Type {
	var resolvedType *sysl.Type
	switch primType {
	case "noprimitive":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_NO_Primitive,
			},
		}
	case "empty":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_EMPTY,
			},
		}
	case "any":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_ANY,
			},
		}
	case "bool":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_BOOL,
			},
		}
	case "int":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_INT,
			},
		}
	case "float":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_FLOAT,
			},
		}
	case "decimal":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_DECIMAL,
			},
		}
	case "string":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_STRING,
			},
		}
	case "bytes":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_BYTES,
			},
		}
	case "string8":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_STRING_8,
			},
		}
	case "date":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_DATE,
			},
		}
	case "datetime":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_DATETIME,
			},
		}
	case "xml":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_XML,
			},
		}
	case "uuid":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Primitive_{
				Primitive: sysl.Type_UUID,
			},
		}
	}
	return resolvedType
}

func MakeEnum(enum map[int64]string) *sysl.Type {
	swappedEnum := make(map[string]int64)
	for index, str := range enum {
		swappedEnum[str] = index
	}
	return &sysl.Type{
		Type: &sysl.Type_Enum_{
			Enum: &sysl.Type_Enum{
				Items: swappedEnum,
			},
		},
	}
}

func MakeTuple(tuple map[string]*sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Tuple_{
			Tuple: &sysl.Type_Tuple{
				AttrDefs: tuple,
			},
		},
	}
}

func MakeList(listType *sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_List_{
			List: &sysl.Type_List{
				Type: listType,
			},
		},
	}
}

func MakeMap(keyType *sysl.Type, valueType *sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Map_{
			Map: &sysl.Type_Map{
				Key:   keyType,
				Value: valueType,
			},
		},
	}
}

func MakeOneOf(oneOfType []*sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_OneOf_{
			OneOf: &sysl.Type_OneOf{
				Type: oneOfType,
			},
		},
	}
}

//TODO relation, set, sequence, notype
func MakeRelation(types map[string]*sysl.Type, primaryKey string, keys []string) *sysl.Type {
	var relationKeys []*sysl.Type_Relation_Key
	for _, val := range keys {
		relationKeys = append(relationKeys, &sysl.Type_Relation_Key{
			AttrName: []string{val},
		})
	}
	relation := &sysl.Type{
		Type: &sysl.Type_Relation_{
			Relation: &sysl.Type_Relation{
				AttrDefs: types,
				PrimaryKey: &sysl.Type_Relation_Key{
					AttrName: []string{primaryKey},
				},
				Key: relationKeys,
			},
		},
	}

	return relation
}

func MakeSet(setType *sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Set{
			Set: setType,
		},
	}
}

func MakeSequence(seqType *sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Sequence{
			Sequence: seqType,
		},
	}
}

func MakeNoType() *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_NoType_{
			NoType: &sysl.Type_NoType{},
		},
	}
}

func MakeType(name string, value interface{}, t string) *sysl.Type {
	var resolvedType *sysl.Type

	switch t {
	case "relation":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Relation_{
				Relation: &sysl.Type_Relation{},
			},
		}
	case "set":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Set{
				Set: MakeType("app", "", "int"),
			},
		}
	case "sequence":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Sequence{
				Sequence: MakeType("app", "", "int"),
			},
		}
	case "notype":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_NoType_{
				NoType: &sysl.Type_NoType{},
			},
		}
	}
	return resolvedType
}

func MakeReturnStatement(payload string) *sysl.Statement {
	return &sysl.Statement{
		Stmt: &sysl.Statement_Ret{
			Ret: &sysl.Return{
				Payload: payload,
			},
		},
	}
}

func MakeParam(name string, paramType *sysl.Type) *sysl.Param {
	var param = sysl.Param{
		Name: name,
		Type: paramType,
	}
	return &param
}

func MakeApp(name string, params []*sysl.Param, types map[string]*sysl.Type) *sysl.Application {
	var appName = MakeAppName(name)

	var ep = sysl.Endpoint{
		Param: params,
	}
	var endpoints = map[string]*sysl.Endpoint{"testEndpoint": &ep}
	var app = sysl.Application{
		Name:      appName,
		LongName:  name,
		Endpoints: endpoints,
		Types:     types,
	}
	return &app
}
