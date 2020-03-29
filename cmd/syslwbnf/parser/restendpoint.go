package parser

import (
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/arr-ai/wbnf/ast"
)

func (p *epScope) buildHttpPath(prefix string, node HttpPathNode,
	urlParams []*sysl.Endpoint_RestParams_QueryParam) *sysl.Endpoint_RestParams {
	var parts []string
	if prefix != "" {
		parts = append(parts, prefix)
	}

	for _, part := range node.AllPart() {
		var sb strings.Builder
		if varPart := part.OneHttpPathVarWithType(); varPart != nil {
			name := varPart.OneVar().OneName().String()

			fmt.Fprintf(&sb, "{%s}", name)
			qp := &sysl.Endpoint_RestParams_QueryParam{
				Name: name,
				Type: &sysl.Type{},
			}
			t := varPart.OneType()
			if ndt := t.OneNativeDataTypes(); ndt != nil {
				buildNativeDataType(ndt, nil, qp.Type)
			} else if ref := t.OneReference(); ref != nil {
				p.buildReference(ref, qp.Type)
			} else if name := t.OneName(); name != nil {
				qp.Type.Type = &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{Ref: &sysl.Scope{Path: []string{name.String()}}},
				}
			}
			urlParams = append(urlParams, qp)
		} else {
			for _, p := range ast.Leafs(part.OneHttpPathPart().Node) {
				_, _ = sb.WriteString(p.Scanner().String())
			}
		}
		parts = append(parts, sb.String())
	}

	path := strings.Join(parts, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return &sysl.Endpoint_RestParams{
		Path:     path,
		UrlParam: urlParams,
	}
}

func (p *epScope) buildRestEndpoint(node *RestEndpointNode, parent *sysl.Endpoint) {
	ep := &sysl.Endpoint{
		Attrs:      AttributeMap{},
		RestParams: &sysl.Endpoint_RestParams{},
	}
	var prefix string
	if parent != nil {
		ep = cloneForChild(parent)
		prefix = ep.RestParams.Path
	}
	if attrs := node.OneAttribs(); attrs != nil {
		ep.Attrs = p.buildAttributes(*attrs)
	}
	ep.RestParams = p.buildHttpPath(prefix, *node.OneHttpPath(), ep.RestParams.UrlParam)
	addPattern(ep.Attrs, "rest")

	for _, restep := range node.AllRestEndpoint() {
		p.buildRestEndpoint(&restep, ep)
	}

	for _, def := range node.AllMethodDef() {
		p.buildRestMethodDef(def, cloneForChild(ep))
	}
}

func (p *epScope) buildRestMethodDef(node MethodDefNode, dest *sysl.Endpoint) {
	method := node.OneMethod().OneToken()
	dest.RestParams.Method = sysl.Endpoint_RestParams_Method(sysl.Endpoint_RestParams_Method_value[method])

	dest.Name = fmt.Sprintf("%s %s", dest.RestParams.Method, dest.RestParams.Path)
	dest.SourceContext = buildSourceContext(node.Node)

	for _, stmt := range node.AllStmt() {
		dest.Stmt = append(dest.Stmt, p.buildStatement(stmt))
	}
	if attrs := node.OneAttribs(); attrs != nil {
		mergeAttrs(p.buildAttributes(*attrs), dest.Attrs)
	}

	if params := node.OneParams(); params != nil {
		for _, param := range params.AllP() {
			val := &sysl.Param{
				Type: &sysl.Type{},
			}
			if field := param.OneField(); field != nil {
				val.Name = field.OneName().String()
				if ft := field.OneFieldType(); ft != nil {
					if col := ft.OneCollectionType(); col != nil {
						// todo
					} else {
						p.buildTypeSpec(*ft.OneTypeSpec(), val.Type)
					}
				}
			} else {
				p.buildReference(param.OneReference(), val.Type)
			}
			dest.Param = append(dest.Param, val)
		}
	}

	p.add(dest)
}

func cloneForChild(in *sysl.Endpoint) *sysl.Endpoint {
	return &sysl.Endpoint{
		LongName:  in.LongName,
		Docstring: in.Docstring,
		Attrs:     in.Attrs, // map so copy is safe
		Flag:      append([]string{}, in.Flag...),
		Source:    in.Source,
		Param:     append([]*sysl.Param{}, in.Param...),
		RestParams: &sysl.Endpoint_RestParams{
			Path:     in.RestParams.Path,
			UrlParam: append([]*sysl.Endpoint_RestParams_QueryParam{}, in.RestParams.UrlParam...),
		},
		SourceContext: nil,
	}
}
