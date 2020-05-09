// package printer prints out sysl datamodels back to source code using the Printer struct.
// Source code does not have complete fidelity, and elements will be printed out in alphabetical order.
package printer

import (
	"fmt"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"io"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
)

const (
	APPLICATIONINDENT = 4
	ENDPOINTINDENT    = 8
	MAXLINE           = 80
	patterns          = "patterns"
)

func p(w io.Writer, I ...interface{}) {
	fmt.Fprint(w, I...)
}

// Module Prints a whole module
func Module(w io.Writer, mod *sysl.Module) {
	for _, key := range alphabeticalApplications(mod.Apps) {
		Application(w, mod.Apps[key])
	}
}

// Application prints applications:
// App:
func Application(w io.Writer, a *sysl.Application) {
	p(w, "\n", strings.Join(a.Name.GetPart(), ""))
	Patterns(w, a.GetAttrs())
	p(w, ":\n")
	for _, key := range alphabeticalAttributes(a.Attrs) {
		if key == patterns {
			continue
		}
		Attrs(w, key, a.Attrs[key], APPLICATIONINDENT)
	}
	for _, key := range alphabeticalTypes(a.Types) {
		TypeDecl(w, key, a.Types[key])
	}
	for _, key := range alphabeticalEndpoints(a.Endpoints) {
		Endpoint(w, a.Endpoints[key])
	}
}

// TypeDecl prints Type decelerations:
// !type Foo:
//     this <: string
func TypeDecl(w io.Writer, key string, t *sysl.Type) {
	switch t.Type.(type) {
	case *sysl.Type_Enum_:
		EnumDecl(w, key, t)
		return
	default:
		NonEnumDecl(w, key, t)
	}
}

func EnumDecl(w io.Writer, key string, t *sysl.Type) {
	switch v := t.Type.(type) {
	case *sysl.Type_Enum_:
		p(w, "    !enum ", key, ":\n")
		for _, key := range alphabeticalAttributes(t.GetAttrs()) {
			if key != patterns {
				Attrs(w, key, t.GetAttrs()[key], ENDPOINTINDENT)
			}
		}
		ef := v.Enum.Items
		for _, key := range alphabeticalInts(ef) {
			p(w, "        ", key, ": ", ef[key], "\n")
		}
	}
}

func NonEnumDecl(w io.Writer, key string, t *sysl.Type) {
	p(w, "    !type ", key)
	Patterns(w, t.GetAttrs())
	p(w, ":\n")
	tuple := t.GetTuple()
	for _, key := range alphabeticalAttributes(t.GetAttrs()) {
		if key != patterns {
			Attrs(w, key, t.GetAttrs()[key], ENDPOINTINDENT)
		}
	}
	if tuple == nil || tuple.AttrDefs == nil || len(tuple.AttrDefs) == 0 {
		p(w, "        ...\n")
		return
	}
	for _, key := range alphabeticalTypes(tuple.AttrDefs) {
		typeClass, typeIdent := syslutil.GetTypeDetail(tuple.AttrDefs[key])
		switch typeClass {
		case "primitive":
			typeIdent = strings.ToLower(typeIdent)
		case "sequence":
			if foo := tuple.AttrDefs[key].GetSequence(); foo != nil {
				typeClass, typeIdent = syslutil.GetTypeDetail(foo)
				if typeClass == "primitive" {
					typeIdent = strings.ToLower(typeIdent)
				}
			}
			typeIdent = "sequence of " + typeIdent
		}
		p(w, "        ", key, " <: ", typeIdent, "\n")
	}
}

// Prints patterns in square brackets: [~foo, ~bar]
func Patterns(w io.Writer, attrs map[string]*sysl.Attribute) {
	if attrs == nil {
		return
	}
	patterns := GetPatterns(attrs)
	if len(patterns) == 0 {
		return
	}
	p(w, "[")
	for i, pattern := range patterns {
		p(w, "~", pattern)
		if i != len(patterns)-1 {
			p(w, ", ")
		}
	}
	p(w, "]")
}

func GetPatterns(attrs map[string]*sysl.Attribute) []string {
	var ret = []string{}
	patterns, has := attrs[patterns]
	if !has {
		return nil
	}
	x := patterns.GetA()
	if x == nil {
		return nil
	}
	for _, y := range x.Elt {
		ret = append(ret, y.GetS())
	}
	return ret
}

// Endpoint prints endpoints:
// Endpoint:
func Endpoint(w io.Writer, e *sysl.Endpoint) {
	p(w, "    ", e.Name)

	if len(e.Param) != 0 {
		Param(w, e.Param)
	}
	Patterns(w, e.Attrs)
	p(w, ":\n")
	if len(e.Stmt) == 0 {
		p(w, "        ...\n")
	}
	for key, attr := range e.Attrs {
		if key == patterns {
			continue
		}
		Attrs(w, key, attr, ENDPOINTINDENT)
	}
	for _, stmnt := range e.Stmt {
		Statement(w, stmnt)
	}
}

// Param prints Parameters:
// Endpoint(This <: ParamHere):
func Param(w io.Writer, params []*sysl.Param) {
	p(w, "(")
	for i, param := range params {
		p(w, param.Name, " <: ", Type(param))
		if i != len(params)-1 {
			p(w, ",")
		}
	}
	p(w, ")")
}

// Attrs prints different statements:
// return string
// My <- call
// lookup db
func Statement(w io.Writer, s *sysl.Statement) {
	if call := s.GetCall(); call != nil {
		Call(w, call)
	}
	if action := s.GetAction(); action != nil {
		Action(w, action)
	}
	if ret := s.GetRet(); ret != nil {
		Return(w, ret)
	}
}

// Return prints return statements:
// return foo <: type
func Return(w io.Writer, r *sysl.Return) {
	p(w, "        return ", r.Payload, "\n")
}

// Action prints actions:
// lookup data
func Action(w io.Writer, a *sysl.Action) {
	p(w, "        ", a.GetAction(), "\n")
}

// Attrs prints Attributes:
// @owner="server"
func Attrs(w io.Writer, key string, a *sysl.Attribute, indentNum int) {
	lines := strings.Split(a.GetS(), "\n")
	indent := strings.Repeat(" ", indentNum)
	if len(lines) == 1 && len(lines[0]) < MAXLINE {
		p(w, indent, "@", key, ` = "`, lines[0], `"`, "\n")
		return
	}
	p(w, indent, "@", key, " =:\n")
	for _, line := range lines {
		for i := 0; i < len(line); i += MAXLINE {
			endIndex := i + MAXLINE
			if lineLen := len(line); endIndex >= lineLen {
				endIndex = lineLen
			}
			p(w, indent, "    |", line[i:endIndex], "\n")
		}
	}
}

// Type prints:
// foo(this <: <Type>):
func Type(param *sysl.Param) string {
	if param.Type == nil {
		return ""
	}
	// Ref.Appname.Part is the type name if the type is in the same package and the application name of where it's at
	// stored if its in another application and then Ref.Path is the type name
	if a := param.Type.GetTypeRef(); a != nil {
		ans := append(a.Ref.Appname.Part, a.Ref.Path...)
		return strings.Join(ans, ".")
	}
	return strings.Join(param.Type.GetTypeRef().Ref.Appname.Part, "")
}

// Call prints:
// AnApp <- AnEndpoint
func Call(w io.Writer, c *sysl.Call) {
	app := strings.Join(c.Target.GetPart(), "")
	ep := c.GetEndpoint()
	p(w, "        ", app, " <- ", ep, "\n")
}
