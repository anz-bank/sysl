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

// Printer prints sysl data structures out to source code
type Printer struct {
	io.Writer
}

// NewPrinter returns a printer that can be used to print out sysl source code from data structures
func NewPrinter(buf io.Writer) *Printer {
	return &Printer{Writer: buf}
}

// PrintModule Prints a whole module, calling
func (p *Printer) PrintModule(mod *sysl.Module) {
	for _, key := range alphabeticalApplications(mod.Apps) {
		p.PrintApplication(mod.Apps[key])
	}
}

// PrintApplication prints applications:
// App:
func (p *Printer) PrintApplication(a *sysl.Application) {
	fmt.Fprintf(p.Writer, "%s:\n", strings.Join(a.Name.GetPart(), ""))
	for _, key := range alphabeticalAttributes(a.Attrs) {
		p.PrintAttrs(key, a.Attrs[key])
	}
	for _, key := range alphabeticalTypes(a.Types) {
		p.PrintTypeDecl(key, a.Types[key])
	}
	for _, key := range alphabeticalEndpoints(a.Endpoints) {
		p.PrintEndpoint(a.Endpoints[key])
	}
}

// PrintTypeDecl prints Type declerations:
// !type Foo:
//     this <: string
func (p *Printer) PrintTypeDecl(key string, t *sysl.Type) {
	fmt.Fprintf(p.Writer, "    !type %s:\n", key)
	if tuple := t.GetTuple(); tuple != nil {
		for _, key := range alphabeticalTypes(tuple.AttrDefs) {
			typeClass, typeIdent := syslutil.GetTypeDetail(tuple.AttrDefs[key])
			if typeClass == "primitive" {
				typeIdent = strings.ToLower(typeIdent)
			}
			fmt.Fprintf(p.Writer, "        %s <: %s\n", key, typeIdent)
		}
	}
}

// PrintEndpoint prints endpoints:
// Endpoint:
func (p *Printer) PrintEndpoint(e *sysl.Endpoint) {
	fmt.Fprintf(p.Writer, "    %s", e.Name)

	if len(e.Param) != 0 {
		p.PrintParam(e.Param)
	}
	fmt.Fprintf(p.Writer, ":\n")
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
func (p *Printer) PrintAttrs(key string, a *sysl.Attribute) {
	fmt.Fprintf(p.Writer, "    @%s=\"%s\"\n", key, a.GetS())
}

// ParamType prints:
// foo(this <: <ParamType>):
func (p *Printer) ParamType(param *sysl.Param) string {
	if param.Type == nil {
		return ""
	}
	if param.Type.GetTypeRef() == nil {
		return ""
	}
	return strings.Join(param.Type.GetTypeRef().Ref.Appname.Part, "")
}

// PrintCall prints:
// AnApp <- AnEndpoint
func (p *Printer) PrintCall(c *sysl.Call) {
	fmt.Fprintf(p.Writer, "        %s <- %s\n", strings.Join(c.Target.GetPart(), ""), c.GetEndpoint())
}
