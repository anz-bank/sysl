package syslutil

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func makeAppName(name string) *sysl.AppName {
	return &sysl.AppName{
		Part: strings.Split(name, " "),
	}
}

// Creates a type reference from one context to another.
// contextApp is the name of the app where the typeref is used
// refApp is the app where the
func makeTypeRef(contextApp string, contextTypePath []string, refApp string, refTypePath []string) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_TypeRef{
			TypeRef: &sysl.ScopedRef{
				Context: &sysl.Scope{
					Appname: makeAppName(contextApp),
					Path:    contextTypePath,
				},
				Ref: &sysl.Scope{
					Appname: makeAppName(refApp),
					Path:    refTypePath,
				},
			},
		},
	}
}

func makePrimitive(primType string) *sysl.Type {
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

func makeEnum(enum map[string]int64) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Enum_{
			Enum: &sysl.Type_Enum{
				Items: enum,
			},
		},
	}
}

func makeTuple(tuple map[string]*sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Tuple_{
			Tuple: &sysl.Type_Tuple{
				AttrDefs: tuple,
			},
		},
	}
}

func makeList(listType *sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_List_{
			List: &sysl.Type_List{
				Type: listType,
			},
		},
	}
}

func makeMap(keyType *sysl.Type, valueType *sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Map_{
			Map: &sysl.Type_Map{
				Key:   keyType,
				Value: valueType,
			},
		},
	}
}

func makeOneOf(oneOfType []*sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_OneOf_{
			OneOf: &sysl.Type_OneOf{
				Type: oneOfType,
			},
		},
	}
}

//TODO Relatino, set, sequence, notype
// This is complext. TODO
func makeRelation(oneOfType []*sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_Relation_{
			Relation: &sysl.Type_Relation{},
		},
	}
}

func makeSet(oneOfType []*sysl.Type) *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_OneOf_{
			OneOf: &sysl.Type_OneOf{
				Type: oneOfType,
			},
		},
	}
}

func makeNoType() *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_NoType_{
			NoType: &sysl.Type_NoType{},
		},
	}
}

func makeType(name string, value interface{}, t string) *sysl.Type {
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
				Set: makeType("app", "", "int"),
			},
		}
	case "sequence":
		resolvedType = &sysl.Type{
			Type: &sysl.Type_Sequence{
				Sequence: makeType("app", "", "int"),
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

func makeParam(name string, paramType *sysl.Type) *sysl.Param {
	var param = sysl.Param{
		Name: name,
		Type: paramType,
	}
	return &param
}

func makeApp(name string, params []*sysl.Param, types map[string]*sysl.Type) sysl.Application {
	var appName = makeAppName(name)

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
	return app
}
