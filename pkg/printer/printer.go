// package printer prints out sysl datamodels back to source code using the Printer struct.
// Source code does not have complete fidelity, and elements will be printed out in alphabetical order.
package printer

import (
	"fmt"
	"io"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
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

// PrintModule Prints a whole module
func PrintModule(w io.Writer, mod *sysl.Module) {
	for _, key := range alphabeticalApplications(mod.Apps) {
		PrintApplication(w, mod.Apps[key])
	}
}

// PrintApplication prints applications:
// App:
func PrintApplication(w io.Writer, a *sysl.Application) {
	p(w, "\n", strings.Join(a.Name.GetPart(), ""))
	PrintPatterns(w, a.GetAttrs())
	p(w, ":\n")
	for _, key := range alphabeticalAttributes(a.Attrs) {
		if key == patterns {
			continue
		}
		PrintAttrs(w, key, a.Attrs[key], APPLICATIONINDENT)
	}
	for _, key := range alphabeticalTypes(a.Types) {
		PrintTypeDecl(w, key, a.Types[key])
	}
	for _, key := range alphabeticalEndpoints(a.Endpoints) {
		PrintEndpoint(w, a.Endpoints[key])
	}
}

// PrintTypeDecl prints Type decelerations:
// !type Foo:
//     this <: string
func PrintTypeDecl(w io.Writer, key string, t *sysl.Type) {
	switch v := t.Type.(type) {
	case *sysl.Type_Enum_:
		p(w, "    !enum ", key, ":\n")
		for _, key := range alphabeticalAttributes(t.GetAttrs()) {
			if key != patterns {
				PrintAttrs(w, key, t.GetAttrs()[key], ENDPOINTINDENT)
			}
		}
		enumFields := v.Enum.Items
		for _, key := range alphabeticalInts(enumFields) {
			p(w, "        ", key, ": ", enumFields[key], "\n")
		}
		return
	}
	p(w, "    !type ", key)
	PrintPatterns(w, t.GetAttrs())
	p(w, ":\n")

	tuple := t.GetTuple()
	for _, key := range alphabeticalAttributes(t.GetAttrs()) {
		if key != patterns {
			PrintAttrs(w, key, t.GetAttrs()[key], ENDPOINTINDENT)
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
		p(w, "        ", key, " <: ", typeIdent)
	}

}

// Prints patterns in square brackets: [~foo, ~bar]
func PrintPatterns(w io.Writer, attrs map[string]*sysl.Attribute) {
	if attrs == nil {
		return
	}
	patterns := GetPatterns(attrs)
	if len(patterns) > 0 {
		p(w, "[")
		for i, pattern := range patterns {
			p(w, "~", pattern)
			if i != len(patterns)-1 {
				p(w, ", ")
			}
		}
		p(w, "]")
	}
}

func GetPatterns(attrs map[string]*sysl.Attribute) []string {
	var ret = []string{}
	patterns, has := attrs[patterns]
	if has {
		if x := patterns.GetA(); x != nil {
			for _, y := range x.Elt {
				ret = append(ret, y.GetS())
			}
		}
	}
	return ret
}

// PrintEndpoint prints endpoints:
// Endpoint:
func PrintEndpoint(w io.Writer, e *sysl.Endpoint) {
	p(w, "    ", e.Name)

	if len(e.Param) != 0 {
		PrintParam(w, e.Param)
	}
	PrintPatterns(w, e.Attrs)
	p(w, ":\n")
	if len(e.Stmt) == 0 {
		p(w, "        ...\n")
	}
	for key, attr := range e.Attrs {
		if key == patterns {
			continue
		}
		PrintAttrs(w, key, attr, ENDPOINTINDENT)
	}
	for _, stmnt := range e.Stmt {
		PrintStatement(w, stmnt)
	}
}

// PrintParam prints Parameters:
// Endpoint(This <: ParamHere):
func PrintParam(w io.Writer, params []*sysl.Param) {
	p(w, "(")
	for i, param := range params {
		p(w, param.Name, " <: ", ParamType(param))
		if i != len(params)-1 {
			p(w, ",")
		}
	}
	p(w, ")")
}

// PrintAttrs prints different statements:
// return string
// My <- call
// lookup db
func PrintStatement(w io.Writer, s *sysl.Statement) {
	if call := s.GetCall(); call != nil {
		PrintCall(w, call)
	}
	if action := s.GetAction(); action != nil {
		PrintAction(w, action)
	}
	if ret := s.GetRet(); ret != nil {
		PrintReturn(w, ret)
	}
}

// PrintReturn prints return statements:
// return foo <: type
func PrintReturn(w io.Writer, r *sysl.Return) {
	p(w, "        return ", r.Payload, "\n")
}

// PrintAction prints actions:
// lookup data
func PrintAction(w io.Writer, a *sysl.Action) {
	p(w, "        ", a.GetAction(), "\n")
}

// PrintAttrs prints Attributes:
// @owner="server"
func PrintAttrs(w io.Writer, key string, a *sysl.Attribute, indentNum int) {
	multiLine := strings.Split(a.GetS(), "\n")
	indent := strings.Repeat(" ", indentNum)
	if len(multiLine) == 1 && len(multiLine[0]) < MAXLINE {
		p(w, indent, "@", key, `="`, multiLine[0], `"`, "\n")
		return
	}
	p(w, indent, "@", key, "=:\n")
	for _, line := range multiLine {
		for i := 0; i < len(line); i += MAXLINE {
			endindex := i + MAXLINE
			if lineLen := len(line); endindex >= lineLen {
				endindex = lineLen
			}
			p(w, indent, "    |", line[i:endindex], "\n")
		}
	}
}

// ParamType prints:
// foo(this <: <ParamType>):
func ParamType(param *sysl.Param) string {
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

// PrintCall prints:
// AnApp <- AnEndpoint
func PrintCall(w io.Writer, c *sysl.Call) {
	app := strings.Join(c.Target.GetPart(), "")
	ep := c.GetEndpoint()
	p(w, "        ", app, "<-", ep, "\n")
}
