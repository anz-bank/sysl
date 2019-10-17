package exporter

import (
	"fmt"
	"strings"

	proto "github.com/anz-bank/sysl/src/proto"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

type TypeExporter struct {
	primitiveTypesMap map[proto.Type_Primitive]map[string]string
	log               *logrus.Logger
}

func makeTypeExporter(logger *logrus.Logger) *TypeExporter {
	return &TypeExporter{
		primitiveTypesMap: map[proto.Type_Primitive]map[string]string{
			proto.Type_NO_Primitive: {"format": "", "type": "object"},
			proto.Type_BOOL:         {"format": "", "type": "boolean"},
			proto.Type_INT:          {"format": "integer", "type": "number"},
			proto.Type_FLOAT:        {"format": "double", "type": "number"},
			proto.Type_DECIMAL:      {"format": "double", "type": "number"},
			proto.Type_STRING:       {"format": "string", "type": "string"},
			proto.Type_BYTES:        {"format": "string", "type": "string"},
			proto.Type_STRING_8:     {"format": "string", "type": "string"},
			proto.Type_DATE:         {"format": "string", "type": "string"},
			proto.Type_DATETIME:     {"format": "string", "type": "string"},
			proto.Type_XML:          {"format": "string", "type": "string"},
		},
		log: logger,
	}
}

func (v *TypeExporter) exportTypes(syslTypes map[string]*proto.Type, swaggerTypes spec.Definitions) error {
	for typeName, typeType := range syslTypes {
		typeSchema := spec.Schema{}
		if v.isComposite(typeType) {
			v.parseComposite(typeType, &typeSchema)
			swaggerTypes[typeName] = typeSchema
			continue
		}
		valueMap, err := v.findSwaggerType(typeType)
		if err != nil {
			v.log.Warnf("Type matching failed %s", err)
			return err
		}
		var ok bool
		typeSchema.Format, ok = valueMap["format"]
		typeSchema.Type = append(typeSchema.Type, valueMap["type"])
		if !ok {
			return fmt.Errorf("malformed sysl missing type info")
		}
		memberTypes := map[string]*proto.Type{}
		if typeSchema.Properties == nil {
			typeSchema.Properties = map[string]spec.Schema{}
		}
		if valueMap["format"] == "tuple" {
			memberTypes = typeType.GetTuple().GetAttrDefs()
		} else if valueMap["format"] == "relation" {
			memberTypes = typeType.GetRelation().GetAttrDefs()
		}
		var attrValueMap map[string]string
		for attK, attV := range memberTypes {
			elementSchema := spec.Schema{}
			if v.isComposite(attV) {
				v.parseComposite(attV, &elementSchema)
				swaggerTypes[attK] = elementSchema
			} else {
				attrValueMap, err = v.findSwaggerType(attV)
				if err != nil {
					return err
				}
			}
			elementSchema.Format = attrValueMap["Format"]
			elementSchema.Type = []string{attrValueMap["type"]}
			typeSchema.Properties[attK] = elementSchema
		}
		swaggerTypes[typeName] = typeSchema
	}
	return nil
}

func (v *TypeExporter) findSwaggerType(syslType *proto.Type) (map[string]string, error) {
	switch s := syslType.Type.(type) {
	case *proto.Type_Primitive_:
		return v.primitiveTypesMap[syslType.GetPrimitive()], nil
	case *proto.Type_Enum_:
		return map[string]string{"format": "integer", "type": "number"}, nil
	case *proto.Type_Tuple_:
		return map[string]string{"format": "tuple", "type": "object"}, nil
	case *proto.Type_Relation_:
		return map[string]string{"format": "relation", "type": "object"}, nil
	case *proto.Type_TypeRef:
		if s.TypeRef.GetRef().GetAppname() == nil {
			return map[string]string{"format": s.TypeRef.GetRef().Path[0], "type": "object"}, nil
		}
		return map[string]string{"format": s.TypeRef.GetRef().GetAppname().Part[0], "type": "object"}, nil
	default:
		return nil, fmt.Errorf("none of the Swagger Types match")
	}
}

func (v *TypeExporter) isComposite(containerType *proto.Type) bool {
	if _, isSet := containerType.Type.(*proto.Type_Set); isSet {
		return true
	} else if _, isSeq := containerType.Type.(*proto.Type_Sequence); isSeq {
		return true
	}
	return false
}

func (v *TypeExporter) isCompositeString(containerType string) bool {
	typeTokens := strings.Split(strings.ToLower(containerType), " ")
	switch typeTokens[0] {
	case "set", "sequence", "map":
		return true
	default:
		return false
	}
}

func (v *TypeExporter) parseComposite(containerType *proto.Type, schema *spec.Schema) {
	schema.Items = &spec.SchemaOrArray{}
	schema.Items.Schema = &spec.Schema{}
	if cType, isSet := containerType.Type.(*proto.Type_Set); isSet {
		schema.Type = []string{"array"}
		retMap, err := v.findSwaggerType(cType.Set)
		if err != nil {
			v.log.Warnf("Type matching failed %s", err)
			return
		}
		schema.Items.Schema.Format = retMap["format"]
		schema.Items.Schema.Type = append(schema.Items.Schema.Type, retMap["type"])
	} else if cType, isSeq := containerType.Type.(*proto.Type_Sequence); isSeq {
		schema.Type = []string{"array"}
		retMap, err := v.findSwaggerType(cType.Sequence)
		if err != nil {
			v.log.Warnf("Type matching failed %s", err)
			return
		}
		schema.Items.Schema.Format = retMap["format"]
		schema.Items.Schema.Type = append(schema.Items.Schema.Type, retMap["type"])
	}
}
