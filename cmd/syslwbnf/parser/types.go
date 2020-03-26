package parser

import (
	"strconv"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func (p *pscope) buildType(node TypeDeclNode) (string, *sysl.Type) {
	if x := node.OneTable(); x != nil {
		return p.buildTable(*x)
	}

	return "", nil
}

func (p *pscope) buildTable(node TableNode) (string, *sysl.Type) {
	name := node.OneName().String()
	val := &sysl.Type{}
	val.SourceContext = buildSourceContext(node.Node)

	if attrs := node.OneAttribs(); attrs != nil {
		val.Attrs = p.buildAttributes(*attrs)
	}
	attrs := map[string]*sysl.Type{}
	var pk []string
	for _, r := range node.AllTableRow() {
		row := &sysl.Type{}
		row.SourceContext = buildSourceContext(r.Node)

		if attrs := r.OneAttribs(); attrs != nil {
			row.Attrs = p.buildAttributes(*attrs)
			if hasPattern(row.Attrs, "pk") {
				pk = append(pk, r.OneName().String())
			}
		}
		if r.OneOptional() != "" {
			row.Opt = true
		}
		attrs[r.OneName().String()] = row

		if collection := r.OneType().OneCollectionType(); collection != nil {
			var nested sysl.Type
			p.buildTypeSpec(*collection.OneTypeSpec(), &nested)
			switch collection.OneCol().OneToken() {
			case "set":
				row.Type = &sysl.Type_Set{Set: &nested}
			case "sequence":
				row.Type = &sysl.Type_Sequence{Sequence: &nested}
			}
		} else if type_spec := r.OneType().OneTypeSpec(); type_spec != nil {
			p.buildTypeSpec(*type_spec, row)
		}
	}
	if node.OneMode().OneToken() == "!table" {
		rel := &sysl.Type_Relation{
			AttrDefs: attrs,
		}
		if len(pk) > 0 {
			rel.PrimaryKey = &sysl.Type_Relation_Key{AttrName: pk}
		}
		val.Type = &sysl.Type_Relation_{Relation: rel}
	} else {
		val.Type = &sysl.Type_Tuple_{Tuple: &sysl.Type_Tuple{
			AttrDefs: attrs,
		}}
	}
	return name, val
}

func (p *pscope) buildTypeSpec(node TypeSpecNode, dest *sysl.Type) {
	if ndt := node.OneNativeDataTypes(); ndt != nil {
		mappings := map[string]sysl.Type_Primitive{
			"int32": sysl.Type_INT, "int64": sysl.Type_INT, "int": sysl.Type_INT,
			"float": sysl.Type_FLOAT, "decimal": sysl.Type_DECIMAL, "bool": sysl.Type_BOOL,
			"date": sysl.Type_DATE, "datetime": sysl.Type_DATETIME,
			"string": sysl.Type_STRING,
		}
		pt := mappings[ndt.OneToken()]
		dest.Type = &sysl.Type_Primitive_{Primitive: pt}
		dest.Constraint = buildTypeConstraints(node.OneSizeSpec(), pt)
		dest.SourceContext = buildSourceContext(node.Node)
	} else if ref := node.OneReference(); ref != nil {
		scope := &sysl.Scope{}
		if pkg := ref.OnePkg(); pkg != nil {
			var path []string
			for _, name := range ref.OnePkg().AllName() {
				path = append(path, name.String())
			}
			scope.Appname = &sysl.AppName{Part: path}
		}
		scope.Path = appName(*ref.OneAppname()).Part
		dest.Type = &sysl.Type_TypeRef{TypeRef: &sysl.ScopedRef{
			Context: nil,
			Ref:     scope,
		}}
	}
}

func buildTypeConstraints(node *TypeSpecSizeSpecNode, primitive sysl.Type_Primitive) []*sysl.Type_Constraint {
	if node == nil || primitive == sysl.Type_NO_Primitive {
		return nil
	}
	var c []*sysl.Type_Constraint
	leading, err := strconv.Atoi(node.OneLeading())
	syslutil.PanicOnErrorf(err, "should not happen: unable to parse SizeSpec:leading")
	var trailing int
	var hasTrailing bool
	if val := node.OneTrailing(); val != "" {
		hasTrailing = true
		trailing, err = strconv.Atoi(val)
		syslutil.PanicOnErrorf(err, "should not happen: unable to parse SizeSpec:trailing")
	}
	if node.OneArray() != "" {
		switch primitive {
		case sysl.Type_DATE, sysl.Type_DATETIME, sysl.Type_INT, sysl.Type_DECIMAL, sysl.Type_STRING:
			ct := &sysl.Type_Constraint{
				Length: &sysl.Type_Constraint_Length{},
			}
			if primitive != sysl.Type_STRING {
				ct.Length.Min = int64(leading)
			}
			if hasTrailing {
				ct.Length.Max = int64(trailing)
			}
			c = append(c, ct)
		default:
			panic("should not be here")
		}
	} else {
		c = append(c, &sysl.Type_Constraint{
			Length: &sysl.Type_Constraint_Length{
				Max: int64(leading),
			},
		})
		switch primitive {
		case sysl.Type_DATE, sysl.Type_DATETIME, sysl.Type_INT, sysl.Type_STRING:
		case sysl.Type_DECIMAL:
			if hasTrailing {
				c[0].Precision = int32(leading)
				c[0].Scale = int32(trailing)
			}
		default:
			panic("should not be here")
		}
	}

	return c
}
