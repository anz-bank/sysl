package exporter

import (
	"fmt"
	"strings"

	proto "github.com/anz-bank/sysl/src/proto"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"
)

type protoType struct {
	Format string
	Type   string
}

type TypeExporter struct {
	primitiveTypesMap map[proto.Type_Primitive]protoType
	log               *logrus.Logger
}

func makeTypeExporter(logger *logrus.Logger) *TypeExporter {
	return &TypeExporter{
		primitiveTypesMap: map[proto.Type_Primitive]protoType{
			proto.Type_NO_Primitive: {Format: "", Type: "object"},
			proto.Type_BOOL:         {Format: "", Type: "boolean"},
			proto.Type_INT:          {Format: "integer", Type: "number"},
			proto.Type_FLOAT:        {Format: "double", Type: "number"},
			proto.Type_DECIMAL:      {Format: "double", Type: "number"},
			proto.Type_STRING:       {Format: "string", Type: "string"},
			proto.Type_BYTES:        {Format: "string", Type: "string"},
			proto.Type_STRING_8:     {Format: "string", Type: "string"},
			proto.Type_DATE:         {Format: "string", Type: "string"},
			proto.Type_DATETIME:     {Format: "string", Type: "string"},
			proto.Type_XML:          {Format: "string", Type: "string"},
		},
		log: logger,
	}
}

func (t *TypeExporter) populateTypes(syslTypes map[string]*proto.Type, swaggerTypes spec.Definitions) error {
	for typeName, dataType := range syslTypes {
		typeSchema := spec.Schema{}
		if t.isComposite(dataType) {
			t.parseComposite(dataType, &typeSchema)
			swaggerTypes[typeName] = typeSchema
			continue
		}
		valueMap, err := t.findSwaggerType(dataType)
		if err != nil {
			t.log.Warnf("Type matching failed %s", err)
			return err
		}
		typeSchema.Format = valueMap.Format
		typeSchema.Type = append(typeSchema.Type, valueMap.Type)
		memberTypes := map[string]*proto.Type{}
		if typeSchema.Properties == nil {
			typeSchema.Properties = map[string]spec.Schema{}
		}
		if valueMap.Format == "tuple" {
			memberTypes = dataType.GetTuple().GetAttrDefs()
		} else if valueMap.Format == "relation" {
			memberTypes = dataType.GetRelation().GetAttrDefs()
		}
		for attK, attV := range memberTypes {
			elementSchema := spec.Schema{}
			if t.isComposite(attV) {
				t.parseComposite(attV, &elementSchema)
				swaggerTypes[attK] = elementSchema
			} else {
				attrValueMap, err := t.findSwaggerType(attV)
				if err != nil {
					t.log.Warnf("Type matching failed %s", err)
					return err
				}
				elementSchema.Format = attrValueMap.Format
				elementSchema.Type = []string{attrValueMap.Type}
			}
			typeSchema.Properties[attK] = elementSchema
		}
		swaggerTypes[typeName] = typeSchema
	}
	return nil
}

func (t *TypeExporter) findSwaggerType(syslType *proto.Type) (protoType, error) {
	switch s := syslType.Type.(type) {
	case *proto.Type_Primitive_:
		return protoType{
			Format: t.primitiveTypesMap[syslType.GetPrimitive()].Format,
			Type:   t.primitiveTypesMap[syslType.GetPrimitive()].Type,
		}, nil
	case *proto.Type_Enum_:
		return protoType{Format: "integer", Type: "number"}, nil
	case *proto.Type_Tuple_:
		return protoType{Format: "tuple", Type: "object"}, nil
	case *proto.Type_Relation_:
		return protoType{Format: "relation", Type: "object"}, nil
	case *proto.Type_TypeRef:
		if s.TypeRef.GetRef().GetAppname() == nil {
			return protoType{Format: s.TypeRef.GetRef().Path[0], Type: "object"}, nil
		}
		return protoType{Format: s.TypeRef.GetRef().GetAppname().Part[0], Type: "object"}, nil
	default:
		return protoType{}, fmt.Errorf("none of the Swagger Types match")
	}
}

func (t *TypeExporter) isComposite(containerType *proto.Type) bool {
	if _, isSet := containerType.Type.(*proto.Type_Set); isSet {
		return true
	} else if _, isSeq := containerType.Type.(*proto.Type_Sequence); isSeq {
		return true
	}
	return false
}

func (t *TypeExporter) isCompositeString(containerType string) bool {
	typeTokens := strings.Split(strings.ToLower(containerType), " ")
	switch typeTokens[0] {
	case "set", "sequence", "map":
		return true
	default:
		return false
	}
}

func (t *TypeExporter) parseComposite(containerType *proto.Type, schema *spec.Schema) {
	schema.Items = &spec.SchemaOrArray{}
	schema.Items.Schema = &spec.Schema{}
	if cType, isSet := containerType.Type.(*proto.Type_Set); isSet {
		schema.Type = []string{"array"}
		retMap, err := t.findSwaggerType(cType.Set)
		if err != nil {
			t.log.Warnf("Type matching failed %s", err)
			return
		}
		schema.Items.Schema.Format = retMap.Format
		schema.Items.Schema.Type = append(schema.Items.Schema.Type, retMap.Type)
	} else if cType, isSeq := containerType.Type.(*proto.Type_Sequence); isSeq {
		schema.Type = []string{"array"}
		retMap, err := t.findSwaggerType(cType.Sequence)
		if err != nil {
			t.log.Warnf("Type matching failed %s", err)
			return
		}
		schema.Items.Schema.Format = retMap.Format
		schema.Items.Schema.Type = append(schema.Items.Schema.Type, retMap.Type)
	}
}
