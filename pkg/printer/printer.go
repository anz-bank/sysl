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
	ENDPOINTINDENT = 8
	MAXLINE = 80
)

// Printer prints sysl data structures out to source code
type Printer struct {
	io.Writer
}

// NewPrinter returns a printer that can be used to print out sysl source code from data structures
func NewPrinter(buf io.Writer) *Printer {
	return &Printer{Writer: buf}
}

// PrintModule Prints a whole module
func (p *Printer) PrintModule(mod *sysl.Module) {
	for _, key := range alphabeticalApplications(mod.Apps) {
		p.PrintApplication(mod.Apps[key])
	}
}

// PrintApplication prints applications:
// App:
func (p *Printer) PrintApplication(a *sysl.Application) {
	fmt.Fprintf(p.Writer, "\n%s", strings.Join(a.Name.GetPart(), ""))
	p.PrintPatterns(a.GetAttrs())
	fmt.Fprint(p.Writer, ":\n")
	for _, key := range alphabeticalAttributes(a.Attrs) {
		if key == "patterns"{
			continue
		}
		p.PrintAttrs(key, a.Attrs[key], APPLICATIONINDENT)
	}
	for _, key := range alphabeticalTypes(a.Types) {
		p.PrintTypeDecl(key, a.Types[key])
	}
	for _, key := range alphabeticalEndpoints(a.Endpoints) {
		p.PrintEndpoint(a.Endpoints[key])
	}
}

// PrintTypeDecl prints Type decelerations:
// !type Foo:
//     this <: string
func (p *Printer) PrintTypeDecl(key string, t *sysl.Type) {
	switch t.Type.(type) {
	case *sysl.Type_Enum_:
		fmt.Fprintf(p.Writer, "    !enum %s:\n", key)
		for _, key := range alphabeticalAttributes(t.GetAttrs()) {
			if key != "patterns"{
				p.PrintAttrs(key, t.GetAttrs()[key], ENDPOINTINDENT)
			}
		}
		enumFields := t.Type.(*sysl.Type_Enum_).Enum.Items
		for _, key := range alphabeticalInts(enumFields) {
			fmt.Fprintf(p.Writer, "        %s: %d\n", key, enumFields[key])
		}

	default:
		fmt.Fprintf(p.Writer, "    !type %s", key)
		p.PrintPatterns(t.GetAttrs())
		fmt.Fprint(p.Writer, ":\n")

		tuple := t.GetTuple()
		for _, key := range alphabeticalAttributes(t.GetAttrs()) {
			if key != "patterns"{
				p.PrintAttrs(key, t.GetAttrs()[key], ENDPOINTINDENT)
			}
		}
		if tuple == nil || tuple.AttrDefs == nil || len(tuple.AttrDefs) == 0{
			fmt.Fprintf(p.Writer, "        ...\n")
			return
		}
		for _, key := range alphabeticalTypes(tuple.AttrDefs) {
			typeClass, typeIdent := syslutil.GetTypeDetail(tuple.AttrDefs[key])
			switch typeClass{
			case "primitive":
				typeIdent = strings.ToLower(typeIdent)
			case "sequence":
				if foo := tuple.AttrDefs[key].GetSequence(); foo != nil{
					typeClass, typeIdent = syslutil.GetTypeDetail(foo)
					if typeClass == "primitive"{
						typeIdent = strings.ToLower(typeIdent)
					}
				}
				typeIdent = "sequence of "+  typeIdent
			}
			fmt.Fprintf(p.Writer, "        %s <: %s\n", key, typeIdent)
		}

	}

}

// Prints patterns in square brackets: [~foo, ~bar]
func (p *Printer)PrintPatterns(attrs map[string]*sysl.Attribute){
	if attrs == nil{
		return
	}
	patterns:= GetPatterns(attrs)
	if len(patterns)>0{
		fmt.Fprint(p.Writer, "[")
		for i, pattern := range patterns{
			fmt.Fprintf(p.Writer, "~%s", pattern)
			if i != len(patterns)-1{
				fmt.Fprintf(p.Writer, ", ")
			}
		}
		fmt.Fprint(p.Writer, "]")
	}
}

func GetPatterns(attrs map[string]*sysl.Attribute)[]string{
	var ret = []string{}
	patterns, has := attrs["patterns"]
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
func (p *Printer) PrintEndpoint(e *sysl.Endpoint) {
	fmt.Fprintf(p.Writer, "    %s", e.Name)

	if len(e.Param) != 0 {
		p.PrintParam(e.Param)
	}
	p.PrintPatterns(e.Attrs)
	fmt.Fprintf(p.Writer, ":\n")
	if len(e.Stmt) == 0 {
		fmt.Fprint(p.Writer, "        ...\n")
	}
	for key, attr := range e.Attrs{
		if key == "patterns"{
			continue
		}
		p.PrintAttrs(key, attr, ENDPOINTINDENT)
	}
	for _, stmnt := range e.Stmt {
		p.PrintStatement(stmnt)
	}
}

// PrintParam prints Parameters:
// Endpoint(This <: ParamHere):
func (p *Printer) PrintParam(params []*sysl.Param) {
	ans := "("
	for i, param := range params {
		ans += param.Name + " <: " + p.ParamType(param)
		if i != len(params)-1 {
			ans += ","
		}
	}
	ans += ")"
	fmt.Fprint(p.Writer, ans)
}

// PrintAttrs prints different statements:
// return string
// My <- call
// lookup db
func (p *Printer) PrintStatement(s *sysl.Statement) {
	if call := s.GetCall(); call != nil {
		p.PrintCall(call)
	}
	if action := s.GetAction(); action != nil {
		p.PrintAction(action)
	}
	if ret := s.GetRet(); ret != nil {
		p.PrintReturn(ret)
	}
}

// PrintReturn prints return statements:
// return foo <: type
func (p *Printer) PrintReturn(r *sysl.Return) {
	fmt.Fprintf(p.Writer, "        return %s\n", r.Payload)
}

// PrintAction prints actions:
// lookup data
func (p *Printer) PrintAction(a *sysl.Action) {
	fmt.Fprintf(p.Writer, "        %s\n", a.GetAction())
}

// PrintAttrs prints Attributes:
// @owner="server"
func (p *Printer) PrintAttrs(key string, a *sysl.Attribute, indentNum int) {
	multiLine := strings.Split(a.GetS(), "\n")
	indent := strings.Repeat(" ", indentNum)
	if len(multiLine)==1 && len (multiLine[0]) < MAXLINE{
			fmt.Fprintf(p.Writer, "%s@%s = \"%s\"\n", indent, key, multiLine[0])
			return

	}

	fmt.Fprintf(p.Writer, "%s@%s =:\n", indent, key)
	for _, line := range multiLine{
		for i := 0; i < len(line); i = i + MAXLINE{
			endindex := i + MAXLINE
			if lineLen := len(line); endindex >= lineLen {
				endindex = lineLen
			}
			fmt.Fprintf(p.Writer, "%s    |%s\n", indent, line[i:endindex])

		}
	}
}

// ParamType prints:
// foo(this <: <ParamType>):
func (p *Printer) ParamType(param *sysl.Param) string {
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
func (p *Printer) PrintCall(c *sysl.Call) {
	fmt.Fprintf(p.Writer, "        %s <- %s\n", strings.Join(c.Target.GetPart(), ""), c.GetEndpoint())
}
