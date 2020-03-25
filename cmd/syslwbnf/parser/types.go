package parser

import "github.com/anz-bank/sysl/pkg/sysl"

func buildType(node TypeDeclNode) (string, *sysl.Type) {
	if x := node.OneTable(); x != nil {
		return buildTable(*x)
	}

	return "", nil
}

func buildTable(node TableNode) (string, *sysl.Type) {
	name := node.OneName().String()
	val := &sysl.Type{}
	val.SourceContext = buildSourceContext(node)

	if attrs := node.OneAttribs(); attrs != nil {
		val.Attrs = buildAttributes(*attrs)
	}
	tuple := &sysl.Type_Tuple{
		AttrDefs: map[string]*sysl.Type{},
	}
	for _, r := range node.AllTableRow() {
		row := &sysl.Type{}
		row.SourceContext = buildSourceContext(r)

		if attrs := r.OneAttribs(); attrs != nil {
			row.Attrs = buildAttributes(*attrs)
		}
		if r.OneOptional() != "" {
			row.Opt = true
		}
		tuple.AttrDefs[r.OneName().String()] = row

		if collection := r.OneType().OneCollectionType(); collection != nil {

		} else if type_spec := r.OneType().OneTypeSpec(); type_spec != nil {
			if ndt := type_spec.OneNativeDataTypes(); ndt != nil {
				mappings := map[string]sysl.Type_Primitive{
					"int32": sysl.Type_INT, "int64": sysl.Type_INT, "int": sysl.Type_INT,
					"float": sysl.Type_FLOAT, "decimal": sysl.Type_DECIMAL, "bool": sysl.Type_BOOL,
					"date": sysl.Type_DATE, "datetime": sysl.Type_DATETIME,
					"string": sysl.Type_STRING,
				}
				row.Type = &sysl.Type_Primitive_{Primitive: mappings[ndt.OneToken()]}
			}
		}
	}
	val.Type = &sysl.Type_Tuple_{Tuple: tuple}
	return name, val
}
