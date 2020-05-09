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

// PrintModule Prints a whole module
func PrintModule(w io.Writer, mod *sysl.Module) {
	for _, key := range alphabeticalApplications(mod.Apps) {
		PrintApplication(w, mod.Apps[key])
	}
}

// PrintApplication prints applications:
// App:
func PrintApplication(w io.Writer, a *sysl.Application) {
	fmt.Fprintf(w, "\n%s", strings.Join(a.Name.GetPart(), ""))
	PrintPatterns(w, a.GetAttrs())
	fmt.Fprint(w, ":\n")
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
		fmt.Fprintf(w, "    !enum %s:\n", key)
		for _, key := range alphabeticalAttributes(t.GetAttrs()) {
			if key != patterns {
				PrintAttrs(w, key, t.GetAttrs()[key], ENDPOINTINDENT)
			}
		}
		enumFields := v.Enum.Items
		for _, key := range alphabeticalInts(enumFields) {
			fmt.Fprintf(w, "        %s: %d\n", key, enumFields[key])
		}

	default:
		fmt.Fprintf(w, "    !type %s", key)
		PrintPatterns(w, t.GetAttrs())
		fmt.Fprint(w, ":\n")

		tuple := t.GetTuple()
		for _, key := range alphabeticalAttributes(t.GetAttrs()) {
			if key != patterns {
				PrintAttrs(w, key, t.GetAttrs()[key], ENDPOINTINDENT)
			}
		}
		if tuple == nil || tuple.AttrDefs == nil || len(tuple.AttrDefs) == 0 {
			fmt.Fprintf(w, "        ...\n")
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
			fmt.Fprintf(w, "        %s <: %s\n", key, typeIdent)
		}
	}
}

// Prints patterns in square brackets: [~foo, ~bar]
func PrintPatterns(w io.Writer, attrs map[string]*sysl.Attribute) {
	if attrs == nil {
		return
	}
	patterns := GetPatterns(attrs)
	if len(patterns) > 0 {
		fmt.Fprint(w, "[")
		for i, pattern := range patterns {
			fmt.Fprintf(w, "~%s", pattern)
			if i != len(patterns)-1 {
				fmt.Fprintf(w, ", ")
			}
		}
		fmt.Fprint(w, "]")
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
	fmt.Fprintf(w, "    %s", e.Name)

	if len(e.Param) != 0 {
		PrintParam(w, e.Param)
	}
	PrintPatterns(w, e.Attrs)
	fmt.Fprintf(w, ":\n")
	if len(e.Stmt) == 0 {
		fmt.Fprint(w, "        ...\n")
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
	ans := "("
	for i, param := range params {
		ans += param.Name + " <: " + ParamType(param)
		if i != len(params)-1 {
			ans += ","
		}
	}
	ans += ")"
	fmt.Fprint(w, ans)
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
	fmt.Fprintf(w, "        return %s\n", r.Payload)
}

// PrintAction prints actions:
// lookup data
func PrintAction(w io.Writer, a *sysl.Action) {
	fmt.Fprintf(w, "        %s\n", a.GetAction())
}

// PrintAttrs prints Attributes:
// @owner="server"
func PrintAttrs(w io.Writer, key string, a *sysl.Attribute, indentNum int) {
	multiLine := strings.Split(a.GetS(), "\n")
	indent := strings.Repeat(" ", indentNum)
	if len(multiLine) == 1 && len(multiLine[0]) < MAXLINE {
		fmt.Fprintf(w, "%s@%s = \"%s\"\n", indent, key, multiLine[0])
		return
	}
	fmt.Fprintf(w, "%s@%s =:\n", indent, key)
	for _, line := range multiLine {
		for i := 0; i < len(line); i += MAXLINE {
			endindex := i + MAXLINE
			if lineLen := len(line); endindex >= lineLen {
				endindex = lineLen
			}
			fmt.Fprintf(w, "%s    |%s\n", indent, line[i:endindex])
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
	fmt.Fprintf(w, "        %s <- %s\n", strings.Join(c.Target.GetPart(), ""), c.GetEndpoint())
}
