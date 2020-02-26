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
func (p *Printer) PrintApplication(A *sysl.Application) {
	fmt.Fprintf(p.Writer, "%s:\n", strings.Join(A.Name.GetPart(), ""))
	for _, key := range alphabeticalAttributes(A.Attrs) {
		p.PrintAttrs(key, A.Attrs[key])
	}
	for _, key := range alphabeticalTypes(A.Types) {
		p.PrintTypeDecl(key, A.Types[key])
	}
	for _, key := range alphabeticalEndpoints(A.Endpoints) {
		p.PrintEndpoint(A.Endpoints[key])
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
func (p *Printer) PrintEndpoint(E *sysl.Endpoint) {
	fmt.Fprintf(p.Writer, "    %s", E.Name)

	if len(E.Param) != 0 {
		p.PrintParam(E.Param)
	}
	fmt.Fprintf(p.Writer, ":\n")
	for _, stmnt := range E.Stmt {
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
func (p *Printer) PrintStatement(S *sysl.Statement) {
	if call := S.GetCall(); call != nil {
		p.PrintCall(call)
	}
	if action := S.GetAction(); action != nil {
		p.PrintAction(action)
	}
	if ret := S.GetRet(); ret != nil {
		p.PrintReturn(ret)
	}
}

// PrintReturn prints return statements:
// return foo <: type
func (p *Printer) PrintReturn(R *sysl.Return) {
	fmt.Fprintf(p.Writer, "        return %s\n", R.Payload)
}

// PrintAction prints actions:
// lookup data
func (p *Printer) PrintAction(A *sysl.Action) {
	fmt.Fprintf(p.Writer, "        %s\n", A.GetAction())
}

// PrintAttrs prints Attributes:
// @owner="server"
func (p *Printer) PrintAttrs(key string, A *sysl.Attribute) {
	fmt.Fprintf(p.Writer, "    @%s=\"%s\"\n", key, A.GetS())
}

// ParamType prints:
// foo(this <: <ParamType>):
func (p *Printer) ParamType(P *sysl.Param) string {
	if P.Type == nil {
		return ""
	}
	if P.Type.GetTypeRef() == nil {
		return ""
	}
	return strings.Join(P.Type.GetTypeRef().Ref.Appname.Part, "")
}

// PrintCall prints:
// AnApp <- AnEndpoint
func (p *Printer) PrintCall(c *sysl.Call) {
	fmt.Fprintf(p.Writer, "        %s <- %s\n", strings.Join(c.Target.GetPart(), ""), c.GetEndpoint())
}
