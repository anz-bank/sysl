package golang

import (
	"fmt"
	"go/token"
	"regexp"
	"strings"
	"unicode"

	"github.com/mitchellh/go-wordwrap"

	"github.com/go-errors/errors"
)

var idRE = regexp.MustCompile("^[\\pL_][\\pL_\\pN]*$")

var underscoreRE = regexp.MustCompile("_([a-z])")

// https://github.com/golang/lint/blob/8f45f776aaf18cebc8d65861cc70c33c60471952/lint.go#L771
var commonInitialismsRE = regexp.MustCompile("(" +
	"Acl|Api|Asc|Cpu|Css|Dns|Eof|Guid|Html|Http|Htt|Id|Ip|Json|Lhs|Qps|Ram|" +
	"Rhs|Rpc|Sla|Smtp|Sql|Ssh|Tcp|Tls|Ttl|Udp|Ui|Uid|Uuid|Uri|Url|Utf8|Vm|" +
	"Xml|Xmpp|Xsrf|Xss" +
	")([A-Z]|$)")

func underscoresToCapsCase(s string, first func(rune) rune) string {
	noUnderscores := underscoreRE.ReplaceAllStringFunc(
		string(first(rune(s[0])))+s[1:],
		func(match string) string {
			return strings.ToUpper(match[1:])
		},
	)
	return commonInitialismsRE.ReplaceAllStringFunc(noUnderscores, strings.ToUpper)
}

// Public creates a public ident with caps-case.
func Public(name string) string {
	return underscoresToCapsCase(name, unicode.ToUpper)
}

// Private creates a private ident with caps-case but lowercase first letter.
func Private(name string) string {
	name = underscoresToCapsCase(name, unicode.ToLower)
	if token.Lookup(name).IsKeyword() {
		return name + "_"
	}
	return name
}

// ArrayN creates a `[n]elt` ArrayType.
func ArrayN(n int, elt Expr) *ArrayType {
	return &ArrayType{Len: Int(n), Elt: elt}
}

// ArrayEllipsis creates a `[...]elt` ArrayType.
func ArrayEllipsis(elt Expr) *ArrayType {
	return &ArrayType{Len: &Ellipsis{}, Elt: elt}
}

// Assert creates a TypeAssertExpr.
func Assert(x, typ Expr) *TypeAssertExpr {
	return &TypeAssertExpr{X: x, Type: typ}
}

// AssertType creates a TypeAssertExpr.
func AssertType(x Expr) *TypeAssertExpr {
	return &TypeAssertExpr{X: x}
}

// Assign creates an = AssignStmt.
func Assign(lhs []Expr, rhs ...Expr) *AssignStmt {
	if len(rhs) == 0 {
		panic("Missing rhs")
	}
	return &AssignStmt{LHS: lhs, Tok: *T("="), RHS: rhs}
}

// Binary creates a BinaryExpr.
func Binary(x Expr, op string, y Expr) *BinaryExpr {
	return &BinaryExpr{X: x, Op: *T(op), Y: y}
}

// Block creates a BlockStmt.
func Block(stmt ...Stmt) *BlockStmt {
	return &BlockStmt{List: stmt}
}

// Break creates a `break` BranchStmt.
func Break() *BranchStmt {
	return &BranchStmt{Tok: *T("break")}
}

// BreakTo creates a `break label` BranchStmt.
func BreakTo(label string) *BranchStmt {
	return &BranchStmt{Tok: *T("break"), Label: I(label)}
}

// Call creates a CallStmt.
func Call(fun Expr, args ...Expr) *CallExpr {
	return &CallExpr{Fun: fun, Args: args}
}

// CallVararg creates an `arg...` CallStmt.
func CallVararg(fun Expr, args ...Expr) *CallExpr {
	return &CallExpr{Fun: fun, Args: args, Ellipsis: *T("...")}
}

// Case creates a `case x, y, z:` CaseClause.
func Case(list []Expr, stmts ...Stmt) *CaseClause {
	return &CaseClause{List: list, Body: stmts}
}

// Comments creates a CommentGroup from a sequence of comments.
func Comments(comments ...string) *CommentGroup {
	list := make([]Comment, 0, len(comments))
	for _, c := range comments {
		list = append(list, Comment{Token: *T(c)})
	}
	return &CommentGroup{List: list}
}

// BeforeGenDecl returns a copy of a GenDecl with Doc added.
func (doc *CommentGroup) BeforeGenDecl(gd GenDecl) *GenDecl {
	gd.Doc = doc
	return &gd
}

// BeforeFuncDecl returns a copy of a FuncDecl with Doc added.
func (doc *CommentGroup) BeforeFuncDecl(fd FuncDecl) *FuncDecl {
	fd.Doc = doc
	return &fd
}

// Const creates a `const` GenDecl.
func Const(values ...ValueSpec) *GenDecl {
	specs := make([]Spec, 0, len(values))
	for _, spec := range values {
		c := spec
		specs = append(specs, &c)
	}
	return &GenDecl{Tok: *T("const"), Specs: specs}
}

// Continue creates a `continue` BranchStmt.
func Continue() *BranchStmt {
	return &BranchStmt{Tok: *T("continue")}
}

// ContinueTo creates a `continue label` BranchStmt.
func ContinueTo(label string) *BranchStmt {
	return &BranchStmt{Tok: *T("continue"), Label: I(label)}
}

// Composite creates a CompositeLit.
func Composite(typ Expr, elts ...Expr) *CompositeLit {
	return &CompositeLit{Type: typ, Elts: elts}
}

// CopyChain copies an IfStmt and, recursively, n.Else if it is an IfStmt.
// Returns the copy and the last node in the chain.
func (n *IfStmt) CopyChain() (head, last *IfStmt) {
	if n == nil {
		return n, n
	}
	c := *n
	head, last = &c, &c
	for i, ok := last.Else.(*IfStmt); ok; i, ok = last.Else.(*IfStmt) {
		last.Else, last = i.CopyChain()
	}
	return head, last
}

// DefaultCase creates a `default:` CaseClause.
func DefaultCase(stmts ...Stmt) *CaseClause {
	return &CaseClause{Body: stmts}
}

// DefaultComm creates a `default:` CommClause.
func DefaultComm(stmts ...Stmt) *CommClause {
	return &CommClause{Body: stmts}
}

// Defer creates a DeferStmt.
func Defer(fun Expr, args ...Expr) *DeferStmt {
	return &DeferStmt{Call: *Call(fun, args...)}
}

// DeferVararg creates an `arg...` DeferStmt.
func DeferVararg(fun Expr, args ...Expr) *DeferStmt {
	return &DeferStmt{Call: *CallVararg(fun, args...)}
}

// Dec creates an `x--` IncDecStmt.
func Dec(x Expr) *IncDecStmt {
	return &IncDecStmt{X: x, Tok: *T("--")}
}

// Dot creates a SelectorExpr.
func Dot(x Expr, id string) *SelectorExpr {
	return &SelectorExpr{X: x, Sel: *I(id)}
}

// Fallthrough creates a `fallthrough` BranchStmt.
func Fallthrough() *BranchStmt {
	return &BranchStmt{Tok: *T("fallthrough")}
}

// NewField creates a Field.
func NewField(typ Expr, ids ...string) *Field {
	return &Field{Names: Idents(ids...), Type: typ}
}

// Float creates a BasicLit for a float64.
func Float(v float64) *BasicLit {
	return &BasicLit{Token{Text: fmt.Sprintf("%#v", v)}}
}

// Func creates a FuncLit.
func Func(params *FieldList, results *FieldList, stmts ...Stmt) *FuncLit {
	return &FuncLit{
		Type: *NewFuncType(params, results),
		Body: BlockStmt{List: stmts},
	}
}

// NewFuncType create a FuncType.
func NewFuncType(fieldLists ...*FieldList) *FuncType {
	var params, results *FieldList
	switch len(fieldLists) {
	case 2:
		results = fieldLists[1]
		fallthrough
	case 1:
		params = fieldLists[0]
	case 0:
	default:
		panic("Too many field lists")
	}
	if len(fieldLists) > 1 {
	}

	if params == nil {
		params = &FieldList{
			Opening: *T("("),
			Closing: *T(")"),
		}
	} else {
		paramsWithParens := *params
		paramsWithParens.Opening = *T("(")
		paramsWithParens.Closing = *T(")")
		params = &paramsWithParens
	}
	if results != nil {
		resultsWithParens := *results
		resultsWithParens.Opening = *T("(")
		resultsWithParens.Closing = *T(")")
		results = &resultsWithParens
	}

	return &FuncType{Params: *params, Results: results}
}

// WordWrappedComment creates an automatically word-wrapped CommentGroup from a
// string of text.
func WordWrappedComment(s string) *CommentGroup {
	lines := strings.Split(wordwrap.WrapString(s, 76), "\n")
	list := make([]Comment, 0, len(lines))
	for _, line := range lines {
		list = append(list, Comment{Token: *T("// " + line)})
	}
	return &CommentGroup{List: list}
}

// WithDoc adds a doc string to FuncDecl, automatically wrapping long strings.
func (fd FuncDecl) WithDoc(doc string) *FuncDecl {
	return WordWrappedComment(fd.Name.Name.Text + " " + doc).BeforeFuncDecl(fd)
}

// WithDocf adds a formatted string to FuncDecl, automatically wrapping long
// strings.
func (fd FuncDecl) WithDocf(format string, a ...interface{}) *FuncDecl {
	return fd.WithDoc(fmt.Sprintf(format, a...))
}

// WithDoc adds a doc string to GenDecl, automatically wrapping long strings.
func (gd GenDecl) WithDoc(doc string) *GenDecl {
	if len(gd.Specs) != 1 {
		panic(errors.Errorf("WithDoc: GenDecl must have exactly one spec"))
	}
	var tname string
	switch spec := gd.Specs[0].(type) {
	case *TypeSpec:
		tname = spec.Name.Name.Text
	case *ValueSpec:
		if len(spec.Names) != 1 {
			panic(errors.Errorf("WithDoc: ValueSpec must have exactly one name"))
		}
		tname = spec.Names[0].Name.Text
	default:
		panic(errors.Errorf("WithDoc: spec must be TypeSpec or ImportSpec"))
	}
	return WordWrappedComment(tname + " " + doc).BeforeGenDecl(gd)
}

// WithDocf adds a formatted string to GenDecl, automatically wrapping long
// strings.
func (gd GenDecl) WithDocf(format string, a ...interface{}) *GenDecl {
	return gd.WithDoc(fmt.Sprintf(format, a...))
}

// Goto creates a `goto label` BranchStmt.
func Goto(label string) *BranchStmt {
	return &BranchStmt{Tok: *T("goto"), Label: I(label)}
}

// I creates an Ident.
func I(id string) *Ident {
	if !idRE.MatchString(id) {
		panic(fmt.Sprintf("Not a valid ident: %#v", id))
	}
	return &Ident{Token{Text: id}}
}

// Idents creates an []Ident from ids.
func Idents(ids ...string) []Ident {
	idents := make([]Ident, 0, len(ids))
	for _, id := range ids {
		idents = append(idents, Ident{Token{Text: id}})
	}
	return idents
}

// If creates an IfStmt.
func If(init Stmt, cond Expr, stmts ...Stmt) *IfStmt {
	return &IfStmt{Init: init, Cond: cond, Body: *Block(stmts...)}
}

// WithElseIf creates a copy of IfStmt with an Else stmt added.
func (n *IfStmt) WithElseIf(init Stmt, cond Expr, stmts ...Stmt) *IfStmt {
	head, last := n.CopyChain()
	if last.Else != nil {
		panic("Cannot chain else to else-tail")
	}
	last.Else = &IfStmt{Init: init, Cond: cond, Body: *Block(stmts...)}
	return head
}

// WithElse creates a copy of IfStmt with an Else stmt added.
func (n *IfStmt) WithElse(stmts ...Stmt) *IfStmt {
	head, last := n.CopyChain()
	if last.Else != nil {
		panic("Cannot chain else to else-tail")
	}
	last.Else = Block(stmts...)
	return head
}

// Import creates an `import` GenDecl.
func Import(imports ...ImportSpec) *GenDecl {
	specs := make([]Spec, 0, len(imports))
	for _, spec := range imports {
		c := spec
		specs = append(specs, &c)
	}
	return &GenDecl{Tok: *T("import"), Specs: specs}
}

// Inc creates an `x++` IncDecStmt.
func Inc(x Expr) *IncDecStmt {
	return &IncDecStmt{X: x, Tok: *T("++")}
}

// Index creates an IndexExpr.
func Index(a, b Expr) *IndexExpr {
	return &IndexExpr{X: a, Index: b}
}

// Init creates an := AssignStmt.
func Init(idents []string, rhs ...Expr) *AssignStmt {
	lhs := make([]Expr, 0, len(idents))
	for _, ident := range idents {
		lhs = append(lhs, I(ident))
	}
	return &AssignStmt{LHS: lhs, Tok: *T(":="), RHS: rhs}
}

// Int creates a BasicLit for an int.
func Int(v int) *BasicLit {
	return &BasicLit{Token{Text: fmt.Sprintf("%#v", v)}}
}

// Interface creates an InterfaceType.
func Interface(fields ...Field) *InterfaceType {
	return &InterfaceType{Methods: FieldList{List: fields}}
}

// KV creates a KeyValueExpr
func KV(key, value Expr) *KeyValueExpr {
	return &KeyValueExpr{Key: key, Value: value}
}

// Map creates a MapType.
func Map(key, value Expr) *MapType {
	return &MapType{Key: key, Value: value}
}

// Method creates a method FuncDecl.
func Method(recvName string, recv Expr, name string, ft FuncType, stmts ...Stmt) *FuncDecl {
	recvArgs := []string{}
	if recvName != "" {
		recvArgs = append(recvArgs, recvName)
	}
	return &FuncDecl{
		Recv: ParenFields(*NewField(recv, recvArgs...)),
		Name: *I(name),
		Type: ft,
		Body: Block(stmts...),
	}
}

// WithParams creates a copy of FuncType with Params containing fields.
func (ft FuncType) WithParams(fields ...Field) *FuncType {
	ft.Params = *ParenFields(fields...)
	return &ft
}

// WithResults creates a copy of FuncType with Result containing fields.
func (ft FuncType) WithResults(fields ...Field) *FuncType {
	ft.Results = ParenFields(fields...)
	return &ft
}

// Nil creates a `nil` Ident.
func Nil() *Ident {
	return I("nil")
}

// ParenFields creates a ()-delimited FieldList.
func ParenFields(fields ...Field) *FieldList {
	return &FieldList{
		Opening: *T("("),
		List:    fields,
		Closing: *T(")"),
	}
}

// Range creates a RangeStmt.
func Range(key, value Expr, tok string, x Expr, body ...Stmt) *RangeStmt {
	return &RangeStmt{Key: key, Value: value, Tok: *T(tok), X: x, Body: *Block(body...)}
}

// Recv creates a `<-ch` UnaryExpr.
func Recv(ch Expr) *UnaryExpr {
	return Unary("<-", ch)
}

// RecvAssignComm creates a `case x, y = <-ch:` CommClause.
func RecvAssignComm(lhs []Expr, ch Expr, stmts ...Stmt) *CommClause {
	return &CommClause{Comm: Assign(lhs, Recv(ch)), Body: stmts}
}

// RecvInitComm creates a `case x, y := <-ch:` CommClause.
func RecvInitComm(lhs []string, ch Expr, stmts ...Stmt) *CommClause {
	return &CommClause{Comm: Init(lhs, Recv(ch)), Body: stmts}
}

// Return creates a ReturnStmt.
func Return(results ...Expr) *ReturnStmt {
	return &ReturnStmt{Results: results}
}

// Select creates a SelectStmt.
func Select(body ...Stmt) *SelectStmt {
	return &SelectStmt{Body: *Block(body...)}
}

// Send creates a SendStmt.
func Send(ch, value Expr) *SendStmt {
	return &SendStmt{Chan: ch, Value: value}
}

// SendComm creates a `case ch <- value:` CommClause.
func SendComm(ch, value Expr, stmts ...Stmt) *CommClause {
	return &CommClause{Comm: Send(ch, value), Body: stmts}
}

// SliceType creates a `[]elt` ArrayType.
func SliceType(elt Expr) *ArrayType {
	return &ArrayType{Elt: elt}
}

// Slice creates a SliceExpr.
func Slice(x Expr, args ...Expr) *SliceExpr {
	if len(args) > 3 {
		panic("Too many slice args")
	}
	args = append(args, nil, nil, nil)[:3]
	return &SliceExpr{X: x, Low: args[0], High: args[1], Max: args[2]}
}

// Star creates a StarExpr.
func Star(x Expr) *StarExpr {
	return &StarExpr{X: x}
}

// String creates a BasicLit for a string.
func String(v string) *BasicLit {
	return &BasicLit{Token{Text: fmt.Sprintf("%#v", v)}}
}

// Struct creates a StructType.
func Struct(fields ...Field) *StructType {
	return &StructType{Fields: FieldList{List: fields}}
}

// Switch creates a SwitchStmt.
func Switch(init Expr, body ...Stmt) *SwitchStmt {
	return &SwitchStmt{Body: *Block(body...)}
}

// T creates a Token.
func T(text string) *Token {
	return &Token{Text: text}
}

// TypeDecl creates a "type" GenDecl.
func TypeDecl(specs ...Spec) *GenDecl {
	return &GenDecl{
		Tok:   *T("type"),
		Specs: specs,
	}
}

// NewTypeSpec create a TypeSpec.
func NewTypeSpec(name string, typ Expr) *TypeSpec {
	return &TypeSpec{
		Name: *I(name),
		Type: typ,
	}
}

// Types creates a `type` GenDecl.
func Types(types ...TypeSpec) *GenDecl {
	specs := make([]Spec, 0, len(types))
	for _, spec := range types {
		c := spec
		specs = append(specs, &c)
	}
	return &GenDecl{Tok: *T("type"), Specs: specs}
}

// TypeSwitch creates a TypeSwitchStmt.
func TypeSwitch(init Stmt, x string, y Expr, body ...Stmt) *TypeSwitchStmt {
	var assign Stmt
	if x != "" {
		assign = Init([]string{x}, AssertType(y))
	} else {
		assign = &ExprStmt{AssertType(y)}
	}
	return &TypeSwitchStmt{Init: init, Assign: assign, Body: *Block(body...)}
}

// Unary creates a UnaryExpr.
func Unary(op string, x Expr) *UnaryExpr {
	return &UnaryExpr{Op: *T(op), X: x}
}

// Var creates a `var` GenDecl.
func Var(values ...ValueSpec) *GenDecl {
	specs := make([]Spec, 0, len(values))
	for _, spec := range values {
		c := spec
		specs = append(specs, &c)
	}
	return &GenDecl{Tok: *T("var"), Specs: specs}
}
