package golang

import (
	"fmt"
	"regexp"
)

var idRE = regexp.MustCompile(`^[\pL_][\pL_\pN]*$`)

// ArrayN creates a `[n]elt` ArrayType.
func ArrayN(n int, elt Expr) *ArrayType {
	return &ArrayType{Len: Int(n), Elt: elt}
}

// ArrayEllipsis creates a `[...]elt` ArrayType.
func ArrayEllipsis(elt Expr) *ArrayType {
	return &ArrayType{Len: &Ellipsis{}, Elt: elt}
}

// Assert creates a TypeAssertExpr.
func Assert(x, t Expr) *TypeAssertExpr {
	return &TypeAssertExpr{X: x, Type: t}
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
func Composite(t Expr, elts ...Expr) *CompositeLit {
	return &CompositeLit{Type: t, Elts: elts}
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

// Float creates a BasicLit for a float64.
func Float(v float64) *BasicLit {
	return &BasicLit{Token{Text: fmt.Sprintf("%#v", v)}}
}

// Func creates a FuncLit.
func Func(params FieldList, results *FieldList, stmts ...Stmt) *FuncLit {
	return &FuncLit{
		Type: FuncType{Params: params, Results: results},
		Body: BlockStmt{List: stmts},
	}
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

// KV creates a KeyValueExpr
func KV(key, value Expr) *KeyValueExpr {
	return &KeyValueExpr{Key: key, Value: value}
}

// Map creates a MapType.
func Map(key, value Expr) *MapType {
	return &MapType{Key: key, Value: value}
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
