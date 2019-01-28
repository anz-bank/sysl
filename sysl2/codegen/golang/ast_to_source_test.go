package golang

import (
	"fmt"
	"go/format"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testFormat(t *testing.T, expected string, n interface{}) {
	expected = strings.TrimSpace(expected)
	expected = strings.Replace(expected, "\n\t\t", "\n", -1)
	actual := fmt.Sprintf("%s", n)
	formatted, err := format.Source([]byte(actual))
	if assert.NoError(t, err, "unformatted: %s", actual) {
		trimmed := strings.TrimSpace(string(formatted))
		assert.Equal(t, expected, trimmed, "unformatted: %s", actual)
	}
}

func testFormatEmpty(t *testing.T, n interface{}) {
	assert.Zero(t, fmt.Sprintf("%s", n))
}

func testFormatPanics(t *testing.T, n interface{}) {
	actual := fmt.Sprintf("%s", n)
	assert.True(t, strings.Contains(actual, "%!s(PANIC="), "%s", actual)
}

func testFormatInfix(t *testing.T, prefix, expected, suffix string, n interface{}) bool {
	expected = strings.TrimSpace(expected)
	expected = strings.Replace(expected, "\n\t\t", "\n", -1)
	output := fmt.Sprintf("%s", n)
	formatted, err := format.Source([]byte(output))
	if !assert.NoError(t, err, "unformatted: %s", output) {
		return false
	}
	actual := strings.TrimSpace(string(formatted))
	if !assert.True(t, strings.HasPrefix(actual, prefix)) ||
		!assert.True(t, strings.HasSuffix(actual, suffix)) {
		t.Errorf("unformatted: %#v\nformatted: %#v", output, actual)
		return false
	}
	return assert.Equal(t,
		expected, strings.TrimSpace(actual[len(prefix):len(actual)-len(suffix)]),
		"unformatted: %s", output)
}

func testFormatType(t *testing.T, expected string, typ Expr) {
	testFormatInfix(t,
		`var _ `, expected, ``,
		&GenDecl{
			Tok:   *T("var"),
			Specs: []Spec{&TypeSpec{Name: *I("_"), Type: typ}},
		},
	)
}

func TestFormat(t *testing.T) {
	_, err := Format("", &File{Name: *I("go")})
	assert.Error(t, err)
	_, err = Format("", &File{Name: *I("goo")})
	assert.NoError(t, err)
}

func TestToken(t *testing.T) {
	testFormatEmpty(t, (*Token)(nil))
	testFormat(t, `foo`, T("foo"))
}

func TestArrayType(t *testing.T) {
	testFormatEmpty(t, (*ArrayType)(nil))
	testFormatType(t, `[]int`, SliceType(I("int")))
	testFormatType(t, `[10]string`, ArrayN(10, I("string")))
}

func TestAssignStmt(t *testing.T) {
	testFormatEmpty(t, (*AssignStmt)(nil))
	testFormat(t, `x = 42`, Assign([]Expr{I("x")}, Int(42)))
	testFormat(t, `x := 42`, Init([]string{"x"}, Int(42)))
	testFormat(t,
		`a, b := 42, "hello"`,
		Init([]string{"a", "b"}, Int(42), String("hello")),
	)
}

func TestBadDecl(t *testing.T) {
	testFormatEmpty(t, (*BadDecl)(nil))
	testFormatPanics(t, &BadDecl{})
}

func TestBadExpr(t *testing.T) {
	testFormatEmpty(t, (*BadExpr)(nil))
	testFormatPanics(t, &BadExpr{})
}

func TestBadStmt(t *testing.T) {
	testFormatEmpty(t, (*BadStmt)(nil))
	testFormatPanics(t, &BadStmt{})
}

func TestBasicLit(t *testing.T) {
	testFormatEmpty(t, (*BasicLit)(nil))
	testFormat(t, `123.45`, &BasicLit{*T("123.45")})
	testFormat(t, `"123.45"`, &BasicLit{*T(`"123.45"`)})
}

func TestBinaryExpr(t *testing.T) {
	testFormatEmpty(t, (*BinaryExpr)(nil))
	// https://golang.org/ref/spec#Operator_precedence
	for _, op := range []string{
		"*", "/", "%", "<<", ">>", "&", "&^",
		"+", "-", "*", "/",
		"==", "!=", "<", "<=", ">", ">=",
		"&&",
		"||",
	} {
		testFormat(t, fmt.Sprintf(`(1 %s 2)`, op), Binary(Int(1), op, Int(2)))
	}
}

func TestBlockStmt(t *testing.T) {
	testFormatEmpty(t, (*BlockStmt)(nil))
	testFormat(t,
		`{
		}`,
		&BlockStmt{},
	)
	testFormat(t,
		`{
			x
		}`,
		&BlockStmt{List: []Stmt{&ExprStmt{X: I("x")}}},
	)
	testFormat(t,
		`{
			x := 42
			return x, nil
		}`,
		&BlockStmt{
			List: []Stmt{Init([]string{"x"}, Int(42)), Return(I("x"), Nil())},
		},
	)
}

func TestBranchStmt(t *testing.T) {
	testFormatEmpty(t, (*BranchStmt)(nil))
	testFormat(t, `break`, Break())
	testFormat(t, `break fast`, BreakTo("fast"))
	testFormat(t, `continue`, Continue())
	testFormat(t, `continue dreaming`, ContinueTo("dreaming"))
	testFormat(t, `goto jail`, Goto("jail"))
	testFormat(t, `fallthrough`, Fallthrough())
}

func TestCallExpr(t *testing.T) {
	testFormatEmpty(t, (*CallExpr)(nil))
	testFormat(t, `f()`, &CallExpr{Fun: I("f")})
	testFormat(t, `sin(π)`, Call(I("sin"), I("π")))
	testFormat(t, `math.Atan2(y, x)`, Call(Dot(I("math"), "Atan2"), I("y"), I("x")))
	testFormat(t,
		`fmt.Sprintf("Hi %s and %s.", names...)`,
		CallVararg(Dot(I("fmt"), "Sprintf"), String("Hi %s and %s."), I("names")),
	)
}

func TestCaseClause(t *testing.T) {
	testFormatEmpty(t, (*CaseClause)(nil))

	testFormatCaseClause := func(expected string, n *CaseClause) {
		testFormatInfix(t,
			`switch {`, expected, `}`,
			&SwitchStmt{Body: BlockStmt{List: []Stmt{n}}},
		)
	}

	testFormatCaseClause(`default:`, DefaultCase())
	testFormatCaseClause(
		`default:
			{
			}`,
		DefaultCase(&BlockStmt{}),
	)
	testFormatCaseClause(
		`default:
			return π`,
		DefaultCase(Return(I("π"))),
	)
	testFormatCaseClause(
		`case (π < 4):
			log.Info("π seems about right!")
			return π`,
		Case([]Expr{Binary(I("π"), "<", Int(4))},
			&ExprStmt{
				X: Call(Dot(I("log"), "Info"), String("π seems about right!")),
			},
			Return(I("π")),
		),
	)
}

func TestChanType(t *testing.T) {
	testFormatEmpty(t, (*ChanType)(nil))
	value := Struct()
	testFormatType(t, `chan struct{}`, &ChanType{Value: value, Dir: ""})
	testFormatType(t, `chan<- struct{}`, &ChanType{Value: value, Dir: "SEND"})
	testFormatType(t, `<-chan struct{}`, &ChanType{Value: value, Dir: "RECV"})
	testFormatPanics(t, &ChanType{Value: value, Dir: "TRANSCEIVE"})
}

func TestCommClause(t *testing.T) {
	testFormatEmpty(t, (*CommClause)(nil))

	testFormatCommClause := func(n *CommClause, expected string) {
		testFormatInfix(t,
			`select {`, expected, `}`,
			&SelectStmt{Body: BlockStmt{List: []Stmt{n}}},
		)
	}

	testFormatCommClause(DefaultComm(), `default:`)
	testFormatCommClause(
		DefaultComm(&BlockStmt{}),
		`default:
			{
			}`,
	)
	testFormatCommClause(
		DefaultComm(Return(I("π"))),
		`default:
			return π`,
	)

	body := []Stmt{
		&ExprStmt{X: Call(Dot(I("log"), "Info"), String("They need π!"))},
		Return(I("π")),
	}
	bodySource := `
			log.Info("They need π!")
			return π`

	testFormatCommClause(
		SendComm(I("ch"), I("π"), body...),
		`case ch <- π:`+bodySource)
	testFormatCommClause(
		RecvAssignComm([]Expr{I("x")}, I("ch"), body...),
		`case x = <-ch:`+bodySource,
	)
	testFormatCommClause(
		RecvInitComm([]string{"x"}, Call(I("ch")), body...),
		`case x := <-ch():`+bodySource,
	)
}

func TestComment(t *testing.T) {
	testFormatEmpty(t, (*Comment)(nil))

	testFormat(t, `//`, &Comment{*T("//")})
	testFormat(t, `// Hmmmm`, &Comment{*T("// Hmmmm")})
	testFormat(t, `/* Hmmmm */`, &Comment{*T("/* Hmmmm */")})
}

func TestCommentGroup(t *testing.T) {
	testFormatEmpty(t, (*CommentGroup)(nil))

	testFormat(t, ``, &CommentGroup{})
	testFormat(t,
		`// Comment 1`,
		&CommentGroup{List: []Comment{{*T("// Comment 1")}}},
	)
	testFormat(t,
		`// Comment 1
		// Comment 2`,
		&CommentGroup{
			List: []Comment{
				{*T("// Comment 1")},
				{*T("// Comment 2")},
			},
		},
	)
}

func TestCompositeLit(t *testing.T) {
	testFormatEmpty(t, (*CompositeLit)(nil))
	testFormat(t, `Foo{X: 42}`, Composite(I("Foo"), KV(I("X"), Int(42))))
	testFormat(t,
		`Foo{X: 42, Y: 3.14159}`,
		Composite(I("Foo"), KV(I("X"), Int(42)), KV(I("Y"), Float(3.14159))),
	)
	testFormat(t,
		`[...]string{"hello", "world"}`,
		Composite(ArrayEllipsis(I("string")), String("hello"), String("world")),
	)
}

func TestDeclStmt(t *testing.T) {
	testFormatEmpty(t, (*DeclStmt)(nil))
	testFormat(t,
		`var x int`,
		&DeclStmt{
			&GenDecl{
				Tok:   *T("var"),
				Specs: []Spec{&TypeSpec{Name: *I("x"), Type: I("int")}},
			},
		},
	)
}

func TestDeferStmt(t *testing.T) {
	testFormatEmpty(t, (*DeferStmt)(nil))
	testFormat(t, `defer f(x)`, Defer(I("f"), I("x")))
	testFormat(t,
		`defer func(a, b int) {
			log.Infof("(%d, %d)", a, b)
		}(x, y)`,
		Defer(
			Func(
				ParenFields(Field{Names: Idents("a", "b"), Type: I("int")}),
				nil,
				&ExprStmt{
					X: Call(Dot(I("log"), "Infof"), String("(%d, %d)"), I("a"), I("b"))},
			),
			I("x"), I("y"),
		),
	)
}

func TestEllipsis(t *testing.T) {
	testFormatEmpty(t, (*Ellipsis)(nil))
	// Tested via ArrayType
}

func TestEmptyStmt(t *testing.T) {
	testFormatEmpty(t, (*EmptyStmt)(nil))
	testFormat(t, `;`, &EmptyStmt{})
}

func TestExprStmt(t *testing.T) {
	testFormatEmpty(t, (*ExprStmt)(nil))
	testFormat(t, `0`, &ExprStmt{X: Int(0)})
	testFormat(t,
		`[]interface{}{}`,
		&ExprStmt{X: Composite(&ArrayType{Elt: &InterfaceType{}})},
	)
}

func TestField(t *testing.T) {
	testFormatEmpty(t, (*Field)(nil))
}

func TestFieldList(t *testing.T) {
	testFormatEmpty(t, (*FieldList)(nil))
}

func TestFile(t *testing.T) {
	testFormatEmpty(t, (*File)(nil))
	testFormat(t, `package main`, &File{Name: *I("main")})
	testFormat(t,
		`package main

		var x = 0`,
		&File{
			Name: *I("main"),
			Decls: []Decl{
				Var(ValueSpec{Names: Idents("x"), Values: []Expr{Int(0)}}),
			},
		},
	)
	testFormat(t,
		`package main

		import "fmt"`,
		&File{
			Name: *I("main"),
			Imports: []ImportSpec{
				{Path: *String("fmt")},
			},
		},
	)
	testFormat(t,
		`package main

		import (
			"fmt"
			"math"
		)

		func main() {
			fmt.Printf("Hello, %g\n", math.PI)
		}`,
		&File{
			Name: *I("main"),
			Imports: []ImportSpec{
				{Path: *String("fmt")},
				{Path: *String("math")},
			},
			Decls: []Decl{
				&FuncDecl{
					Name: *I("main"),
					Type: FuncType{Params: *ParenFields()},
					Body: &BlockStmt{
						List: []Stmt{
							&ExprStmt{
								X: Call(
									Dot(I("fmt"), "Printf"),
									String("Hello, %g\n"),
									Dot(I("math"), "PI"),
								),
							},
						},
					},
				},
			},
		},
	)
	testFormat(t,
		`package foo

		// A ...
		type A struct{}

		// B ...
		type B struct{}`,
		&File{
			Name: *I("foo"),
			Decls: []Decl{
				&GenDecl{
					Doc: Comments("// A ..."),
					Tok: *T("type"),
					Specs: []Spec{
						&TypeSpec{Name: *I("A"), Type: Struct()},
					},
				},
				&GenDecl{
					Doc: Comments("// B ..."),
					Tok: *T("type"),
					Specs: []Spec{
						&TypeSpec{Name: *I("B"), Type: Struct()},
					},
				},
			},
		},
	)
}

func TestForStmt(t *testing.T) {
	testFormatEmpty(t, (*ForStmt)(nil))

	testFormat(t,
		`for {
		}`,
		&ForStmt{},
	)
	testFormat(t,
		`for {
			0
		}`,
		&ForStmt{Body: *Block(&ExprStmt{X: Int(0)})},
	)
}

func TestFuncDecl(t *testing.T) {
	testFormatEmpty(t, (*FuncDecl)(nil))

	testFormat(t,
		`func f(a, b int) {
			log.Infof("(%d, %d)", a, b)
		}`,
		&FuncDecl{
			Name: *I("f"),
			Type: FuncType{
				Params: *ParenFields(
					Field{Names: Idents("a", "b"), Type: I("int")},
				),
			},
			Body: Block(
				&ExprStmt{
					X: Call(Dot(I("log"), "Infof"),
						&BasicLit{Token{Text: `"(%d, %d)"`}}, I("a"), I("b"),
					),
				},
			),
		},
	)
	testFormat(t,
		`func f(a int, b ...int) {}`,
		&FuncDecl{
			Name: *I("f"),
			Type: FuncType{
				Params: *ParenFields(
					Field{Names: Idents("a"), Type: I("int")},
					Field{Names: Idents("b"), Type: &Ellipsis{Elt: I("int")}},
				),
			},
			Body: Block(),
		},
	)
}

func TestFuncLit(t *testing.T) {
	testFormatEmpty(t, (*FuncLit)(nil))
	testFormat(t,
		`[]interface{}{func(a, b int) {
			log.Infof("(%d, %d)", a, b)
		}}`,
		Composite(
			&ArrayType{Elt: &InterfaceType{}},
			Func(
				ParenFields(Field{Names: Idents("a", "b"), Type: I("int")}),
				nil,
				&ExprStmt{
					X: Call(Dot(I("log"), "Infof"),
						&BasicLit{Token{Text: `"(%d, %d)"`}}, I("a"), I("b"),
					),
				},
			),
		),
	)
}

func TestFuncType(t *testing.T) {
	testFormatEmpty(t, (*FuncType)(nil))
	// Tested via FuncDecl and FuncLit.
}

func TestGenDeclImport(t *testing.T) {
	testFormatEmpty(t, (*GenDecl)(nil))

	testFormat(t, `import "fmt"`, Import(ImportSpec{Path: *String("fmt")}))
	testFormat(t,
		`import (
			queries "database/sql"
			"fmt"
			"io"
			"math"
		)`,
		Import(
			ImportSpec{Name: I("queries"), Path: *String("database/sql")},
			ImportSpec{Path: *String("fmt")},
			ImportSpec{Path: *String("io")},
			ImportSpec{Path: *String("math")},
		),
	)
}

func TestGenDeclVar(t *testing.T) {
	testFormatEmpty(t, (*GenDecl)(nil))

	testFormat(t,
		`var x = 42`,
		Var(ValueSpec{Names: Idents("x"), Values: []Expr{Int(42)}}),
	)
	testFormat(t,
		`var (
			x              = 42
			r, g, b string = "red", "green", "blue"
		)`,
		Var(
			ValueSpec{Names: Idents("x"), Values: []Expr{Int(42)}},
			ValueSpec{
				Names:  Idents("r, g, b"),
				Type:   I("string"),
				Values: []Expr{String("red"), String("green"), String("blue")},
			},
		),
	)
}

func TestGenDeclConst(t *testing.T) {
	testFormatEmpty(t, (*GenDecl)(nil))

	testFormat(t,
		`const x = 42`,
		Const(ValueSpec{Names: Idents("x"), Values: []Expr{Int(42)}}),
	)
	testFormat(t,
		`const (
			x              = 42
			r, g, b string = "red", "green", "blue"
		)`,
		Const(
			ValueSpec{Names: Idents("x"), Values: []Expr{Int(42)}},
			ValueSpec{
				Names:  Idents("r, g, b"),
				Type:   I("string"),
				Values: []Expr{String("red"), String("green"), String("blue")},
			},
		),
	)
}

func TestGenDeclType(t *testing.T) {
	testFormatEmpty(t, (*GenDecl)(nil))

	testFormat(t,
		`type id uint64`,
		Types(TypeSpec{Name: *I("id"), Type: I("uint64")}),
	)

	coords := Idents("x", "y", "z", "w")

	testFormat(t,
		`type (
			vec2 struct {
				x, y float32
			}
			vec3 struct {
				x, y, z float32
			}
			vec4 struct {
				x, y, z, w float32
			}
		)`,
		Types(
			TypeSpec{
				Name: *I("vec2"),
				Type: Struct(Field{Names: coords[:2], Type: I("float32")}),
			},
			TypeSpec{
				Name: *I("vec3"),
				Type: Struct(Field{Names: coords[:3], Type: I("float32")}),
			},
			TypeSpec{
				Name: *I("vec4"),
				Type: Struct(Field{Names: coords[:4], Type: I("float32")}),
			},
		),
	)
}

func TestGoStmt(t *testing.T) {
	testFormatEmpty(t, (*GoStmt)(nil))
	testFormat(t, `go f()`, &GoStmt{Call: CallExpr{Fun: I("f")}})
}

func TestIdent(t *testing.T) {
	testFormatEmpty(t, (*Ident)(nil))
	testFormat(t, `x`, I("x"))
	testFormat(t, `fooBarBaz`, I("fooBarBaz"))
	testFormat(t, `π`, I("π"))
}

func TestIfStmt(t *testing.T) {
	testFormatEmpty(t, (*IfStmt)(nil))

	testFormat(t,
		`if ok {
			return 0
		}`,
		If(nil, I("ok"), Return(Int(0))),
	)
	testFormat(t,
		`if x, ok := f(); ok {
			return 0
		}`,
		If(Init([]string{"x", "ok"}, Call(I("f"))), I("ok"), Return(Int(0))),
	)
	testFormat(t,
		`if x, ok := f(); ok {
			return 0
		}`,
		If(Init([]string{"x", "ok"}, Call(I("f"))), I("ok"), Return(Int(0))),
	)
}

func TestIfElse(t *testing.T) {
	testFormat(t,
		`if x, ok := f(); ok {
			return 0
		} else {
			return 1
		}`,
		If(Init([]string{"x", "ok"}, Call(I("f"))), I("ok"), Return(Int(0))).
			WithElse(Return(Int(1))),
	)
	testFormat(t,
		`if x, ok := f(); ok {
			y := (1 / x)
			return y
		} else {
			fmt.Printf("Oops! (%d)", x)
			return 1
		}`,
		If(
			Init([]string{"x", "ok"}, Call(I("f"))),
			I("ok"),
			Init([]string{"y"}, Binary(Int(1), "/", I("x"))),
			Return(I("y")),
		).WithElse(
			&ExprStmt{
				X: Call(
					Dot(I("fmt"), "Printf"),
					String("Oops! (%d)"),
					I("x"),
				),
			},
			Return(Int(1)),
		),
	)
}

func TestIfElseIf(t *testing.T) {
	testFormat(t,
		`if x, ok := f(); ok {
			y := (1 / x)
			return y
		} else if panicOnError {
			panic("Ouch!")
		}`,
		If(
			Init([]string{"x", "ok"}, Call(I("f"))),
			I("ok"),
			Init([]string{"y"}, Binary(Int(1), "/", I("x"))),
			Return(I("y")),
		).WithElseIf(
			nil,
			I("panicOnError"),
			&ExprStmt{X: Call(I("panic"), String("Ouch!"))},
		),
	)
}

func TestIfElseIfElse(t *testing.T) {
	testFormat(t,
		`if x, ok := f(); ok {
			y := (1 / x)
			return y
		} else if panicOnError {
			panic("Ouch!")
		} else {
			fmt.Printf("Oops! (%d)", x)
			return 1
		}`,
		If(
			Init([]string{"x", "ok"}, Call(I("f"))),
			I("ok"),
			Init([]string{"y"}, Binary(Int(1), "/", I("x"))),
			Return(I("y")),
		).WithElseIf(
			nil,
			I("panicOnError"),
			&ExprStmt{X: Call(I("panic"), String("Ouch!"))},
		).WithElse(
			&ExprStmt{
				X: Call(
					Dot(I("fmt"), "Printf"),
					String("Oops! (%d)"),
					I("x"),
				),
			},
			Return(Int(1)),
		),
	)
}

func TestIfElseIfElseIf(t *testing.T) {
	testFormat(t,
		`if x, ok := f(); ok {
			y := (1 / x)
			return y, nil
		} else if panicOnError {
			panic("Ouch!")
		} else if defaultOnError {
			fmt.Printf("Oops! (%d)", x)
			return 1, nil
		}`,
		If(
			Init([]string{"x", "ok"}, Call(I("f"))),
			I("ok"),
			Init([]string{"y"}, Binary(Int(1), "/", I("x"))),
			Return(I("y"), Nil()),
		).WithElseIf(
			nil,
			I("panicOnError"),
			&ExprStmt{X: Call(I("panic"), String("Ouch!"))},
		).WithElseIf(
			nil,
			I("defaultOnError"),
			&ExprStmt{
				X: Call(Dot(I("fmt"), "Printf"), String("Oops! (%d)"), I("x")),
			},
			Return(Int(1), Nil()),
		),
	)
}

func TestIfElseIfElseIfElse(t *testing.T) {
	testFormat(t,
		`if x, ok := f(); ok {
			y := (1 / x)
			return y, nil
		} else if panicOnError {
			panic("Ouch!")
		} else if defaultOnError {
			fmt.Printf("Oops! (%d)", x)
			return 1, nil
		} else {
			return 0, fmt.Errorf("Oops! (%d)", x)
		}`,
		If(
			Init([]string{"x", "ok"}, Call(I("f"))),
			I("ok"),
			Init([]string{"y"}, Binary(Int(1), "/", I("x"))),
			Return(I("y"), Nil()),
		).WithElseIf(
			nil,
			I("panicOnError"),
			&ExprStmt{X: Call(I("panic"), String("Ouch!"))},
		).WithElseIf(
			nil,
			I("defaultOnError"),
			&ExprStmt{
				X: Call(Dot(I("fmt"), "Printf"), String("Oops! (%d)"), I("x")),
			},
			Return(Int(1), Nil()),
		).WithElse(
			Return(Int(0), Call(Dot(I("fmt"), "Errorf"), String("Oops! (%d)"), I("x"))),
		),
	)
}

func TestImportSpec(t *testing.T) {
	testFormatEmpty(t, (*ImportSpec)(nil))
	testFormat(t, ``, &ImportSpec{})
}

func TestIncDecStmt(t *testing.T) {
	testFormatEmpty(t, (*IncDecStmt)(nil))
	testFormat(t, `x++`, Inc(I("x")))
	testFormat(t, `x--`, Dec(I("x")))
}

func TestIndexExpr(t *testing.T) {
	testFormatEmpty(t, (*IndexExpr)(nil))
	testFormat(t, `x[y]`, Index(I("x"), I("y")))
}

func TestInterfaceType(t *testing.T) {
	testFormatEmpty(t, (*InterfaceType)(nil))
	testFormatType(t, `interface{}`, &InterfaceType{})
	testFormatType(t,
		`interface {
			Foo(a, b int) (int, string)
		}`,
		&InterfaceType{
			Methods: FieldList{
				List: []Field{
					{
						Names: Idents("Foo"),
						Type: &FuncType{
							Params: *ParenFields(
								Field{Names: Idents("a", "b"), Type: I("int")},
							),
							Results: ParenFields(
								Field{Type: I("int")},
								Field{Type: I("string")},
							),
						},
					},
				},
			},
		},
	)
}

func TestKeyValueExpr(t *testing.T) {
	testFormatEmpty(t, (*KeyValueExpr)(nil))
}

func TestLabeledStmt(t *testing.T) {
	testFormatEmpty(t, (*LabeledStmt)(nil))

	testFormat(t,
		`here:
		return x, y`,
		&LabeledStmt{Label: *I("here"), Stmt: Return(I("x"), I("y"))},
	)
	testFormat(t,
		`here:
		{
			return x, y
		}`,
		&LabeledStmt{Label: *I("here"), Stmt: Block(Return(I("x"), I("y")))},
	)
}

func TestMapType(t *testing.T) {
	testFormatEmpty(t, (*MapType)(nil))

	testFormatType(t, `map[int]string`, Map(I("int"), I("string")))
	testFormatType(t,
		`map[int]map[float]string`,
		Map(I("int"), Map(I("float"), I("string"))),
	)
}

func TestParenExpr(t *testing.T) {
	testFormatEmpty(t, (*ParenExpr)(nil))
	testFormat(t, `x`, &ParenExpr{X: I("x")})
	testFormat(t, `(x ^ y)`, &ParenExpr{X: Binary(I("x"), "^", I("y"))})
}

func TestRangeStmt(t *testing.T) {
	testFormatEmpty(t, (*RangeStmt)(nil))

	testFormat(t,
		`for range records {
			x++
		}`,
		Range(nil, nil, "", I("records"), &IncDecStmt{X: I("x"), Tok: *T("++")}),
	)
	testFormat(t,
		`for i := range input {
			output[i] = input[i]
		}`,
		Range(I("i"), nil, ":=", I("input"),
			Assign([]Expr{Index(I("output"), I("i"))}, Index(I("input"), I("i"))),
		),
	)
	testFormat(t,
		`for _, x := range input {
			output = append(output, f(x))
		}`,
		Range(I("_"), I("x"), ":=", I("input"),
			Assign(
				[]Expr{I("output")},
				Call(I("append"), I("output"), Call(I("f"), I("x"))),
			),
		),
	)
	testFormat(t,
		`for i, elt = range os.Args {
			if i > 0 {
				process(elt)
			}
		}`,
		Range(I("i"), I("elt"), "=", Dot(I("os"), "Args"),
			If(nil, Binary(I("i"), ">", Int(0)), &ExprStmt{X: Call(I("process"), I("elt"))}),
		),
	)
}

func TestReturnStmt(t *testing.T) {
	testFormatEmpty(t, (*ReturnStmt)(nil))
	testFormat(t, `return`, Return())
	testFormat(t, `return x`, Return(I("x")))
}

func TestSelectStmt(t *testing.T) {
	testFormatEmpty(t, (*SelectStmt)(nil))

	testFormat(t, `select {}`, Select())
	testFormat(t,
		`select {
		case ch <- value:
		}`,
		Select(SendComm(I("ch"), I("value"))),
	)
	testFormat(t,
		`select {
		case ch <- value:
			x++
		}`,
		Select(SendComm(I("ch"), I("value"), Inc(I("x")))),
	)
	testFormat(t,
		`select {
		case value := <-ch:
			return value
		}`,
		Select(RecvInitComm([]string{"value"}, I("ch"), Return(I("value")))),
	)
	testFormat(t,
		`select {
		case ch <- value:
			x++
		case value := <-ch:
			return value
		default:
			panic("I can't wait!")
		}`,
		Select(
			SendComm(I("ch"), I("value"), Inc(I("x"))),
			RecvInitComm([]string{"value"}, I("ch"), Return(I("value"))),
			DefaultComm(&ExprStmt{Call(I("panic"), String("I can't wait!"))}),
		),
	)
}

func TestSelectorExpr(t *testing.T) {
	testFormatEmpty(t, (*SelectorExpr)(nil))
	testFormat(t, `x.y`, Dot(I("x"), "y"))
}

func TestSendStmt(t *testing.T) {
	testFormatEmpty(t, (*SendStmt)(nil))
	testFormat(t, `ch <- x`, Send(I("ch"), I("x")))
}

func TestSliceExpr(t *testing.T) {
	testFormatEmpty(t, (*SliceExpr)(nil))

	x := I("x")
	low := I("low")
	high := I("high")
	max := I("max")

	testFormat(t, `x[:]`, Slice(x))
	testFormat(t, `x[low:]`, Slice(x, low))
	testFormat(t, `x[:high]`, Slice(x, nil, high))
	testFormat(t, `x[low:high]`, Slice(x, low, high))
	testFormat(t, `x[:high:max]`, Slice(x, nil, high, max))
	testFormat(t, `x[low:high:max]`, Slice(x, low, high, max))

	assert.Panics(t, func() { Slice(x, low, high, max, I("wat")) })
}

func TestStarExpr(t *testing.T) {
	testFormatEmpty(t, (*StarExpr)(nil))
	testFormat(t, `*x`, Star(I("x")))
	testFormat(t, `*(x + y)`, Star(Binary(I("x"), "+", I("y"))))
}

func TestStructType(t *testing.T) {
	testFormatEmpty(t, (*StructType)(nil))

	testFormatType(t, `struct{}`, Struct())
	testFormatType(t,
		`struct {
			X int
		}`,
		Struct(Field{Names: Idents("X"), Type: I("int")}),
	)
	testFormatType(t,
		`struct {
			X, Y int
			Z    int "json:\"id\""
		}`,
		Struct(
			Field{Names: Idents("X", "Y"), Type: I("int")},
			Field{Names: Idents("Z"), Type: I("int"), Tag: String(`json:"id"`)},
		),
	)
}

func TestSwitchStmt(t *testing.T) {
	testFormatEmpty(t, (*SwitchStmt)(nil))
	testFormat(t,
		`switch {
		}`,
		&SwitchStmt{},
	)
	testFormat(t,
		`switch {
		default:
			{
			}
		}`,
		&SwitchStmt{Body: *Block(DefaultCase(&BlockStmt{}))},
	)
	testFormat(t,
		`switch {
		default:
			{
			}
		}`,
		&SwitchStmt{Body: *Block(DefaultCase(&BlockStmt{}))},
	)
}

func TestTypeAssertExpr(t *testing.T) {
	testFormatEmpty(t, (*TypeAssertExpr)(nil))
	testFormat(t, `x.(y)`, Assert(I("x"), I("y")))
	testFormat(t, `x.(type)`, AssertType(I("x")))
}

func TestTypeSpec(t *testing.T) {
	testFormatEmpty(t, (*TypeSpec)(nil))
}

func TestTypeSwitchStmt(t *testing.T) {
	testFormatEmpty(t, (*TypeSwitchStmt)(nil))

	testFormat(t,
		`switch y.(type) {
		}`,
		TypeSwitch(nil, "", I("y")),
	)
	testFormat(t,
		`switch x := y.(type) {
		}`,
		TypeSwitch(nil, "x", I("y")),
	)
	testFormat(t,
		`switch y := 1; x := y.(type) {
		case int, float:
			return 0
		}`,
		TypeSwitch(
			Init([]string{"y"}, Int(1)), "x", I("y"),
			Case([]Expr{I("int"), I("float")}, Return(Int(0))),
		),
	)
}

func TestUnaryExpr(t *testing.T) {
	testFormatEmpty(t, (*UnaryExpr)(nil))
	// https://golang.org/ref/spec#Operators
	for _, op := range []string{"+", "-", "!", "^", "*", "&", "<-"} {
		testFormat(t, op+`x`, Unary(op, I("x")))
	}
}

func TestValueSpec(t *testing.T) {
	testFormatEmpty(t, (*ValueSpec)(nil))
}
