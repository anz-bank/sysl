package eval

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	parser "github.com/anz-bank/sysl/src/sysl/grammar"
	"github.com/anz-bank/sysl/src/sysl/parse"

	sysl "github.com/anz-bank/sysl/src/proto_old"
)

type DebugFunc func(scope *Scope, app *sysl.Application, expr *sysl.Expr) error

type repl struct {
	input       *bufio.Scanner
	output      io.Writer
	breakpoints []*sysl.SourceContext
	step        bool
}

func (r *repl) run(scope *Scope, app *sysl.Application, expr *sysl.Expr) error {
	if expr == nil || r.step || r.isBreakpoint(expr) {
		return r.handleInput(scope, app, expr)
	}
	return nil
}

func (r *repl) handleInput(scope *Scope, app *sysl.Application, expr *sysl.Expr) error {
	for {
		writeOutput("sysl> ", r.output)
		if !r.input.Scan() {
			r.step = false
			return fmt.Errorf("input closed")
		}
		text := r.input.Text()
		switch len(text) {
		case 0:
			printUsage(r.output)
		case 1:
			switch text[0] {
			case 's':
				writeOutput("(step)\n", r.output)
				r.step = true
				return nil
			case 'c':
				writeOutput("(continue)\n", r.output)
				r.step = false
				return nil
			case 'd':
				writeOutput("(dump scope)\n", r.output)
				dumpScope(r.output, *scope)
			case 'q':
				writeOutput("(quit)\n", r.output)
				os.Exit(0)
			case '?', 'h':
				printUsage(r.output)
			default:
				writeOutput(parseExpression(text, scope), r.output)
			}
		default:
			writeOutput(parseExpression(text, scope), r.output)
		}
	}
}
func (r *repl) isBreakpoint(expr *sysl.Expr) bool {
	if expr.SourceContext == nil {
		return false
	}
	sc := expr.SourceContext
	for _, bsc := range r.breakpoints {
		if strings.HasSuffix(sc.File, bsc.File) &&
			sc.Start.Line == bsc.Start.Line &&
			(sc.Start.Col == bsc.Start.Col || bsc.Start.Col < 0) {
			return true
		}
	}
	return false
}

// Breakpoint format:
// filename:line[:col]
func makeBreakpoint(from string) *sysl.SourceContext {
	parts := strings.Split(from, ":")
	if len(parts) >= 2 {
		line, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return nil
		}
		res := &sysl.SourceContext{
			File: parts[0],
			Start: &sysl.SourceContext_Location{
				Line: int32(line),
				Col:  -1,
			},
		}
		if len(parts) >= 3 {
			col, err := strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return nil
			}
			res.Start.Col = int32(col)
		}
		return res
	}
	return nil
}

func printUsage(out io.Writer) {
	text := `
?, h		This help text
d		Dump the current scope
c		Continue execution until the next breakpoint
s		Step execution
EXPRESSION	Print the value of the supplied EXPRESSION
`
	writeOutput(text, out)
}

func dumpScope(out io.Writer, scope Scope) {
	names := make([]string, 0, len(scope))
	for k := range scope {
		names = append(names, k)
	}
	sort.Strings(names)
	lines := make([]string, 0, len(scope))
	for _, k := range names {
		lines = append(lines, fmt.Sprintf("  %s:  %s", k, scope[k].String()))
	}

	writeOutput(strings.Join(lines, "\n")+"\n", out)
}

func parseExpression(text string, scope *Scope) string {
	errorListener := parse.SyslParserErrorListener{}
	lexer := parser.NewSyslLexer(antlr.NewInputStream(text))
	lexer.SetMode(parser.SyslLexerVIEW_TRANSFORM)
	defer parser.DeleteLexerState(lexer)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	listener := parse.NewTreeShapeListener()

	p := parser.NewSyslParser(stream)
	p.AddErrorListener(&errorListener)
	p.AddParseListener(listener)

	var res *sysl.Expr

	ee := exprEval{
		txApp:     nil,
		exprStack: exprStack{},
		logger:    logrus.StandardLogger(),
	}

	// Need to figure out what sort of expression this is going to be
	lookAhead := []int{stream.LA(1), stream.LA(2)}
	if lookAhead[0] == parser.SyslLexerE_LET || lookAhead[1] == parser.SyslLexerE_EQ {
		// Need to push a transform onto the expression stack
		transform := &sysl.Expr_Transform{}
		listener.PushExpr(&sysl.Expr{Expr: &sysl.Expr_Transform_{Transform: transform}})
		if lookAhead[0] == parser.SyslLexerE_LET {
			p.Expr_let_statement()
		} else {
			p.Expr_simple_assign()
		}
		xformres := evalTransformStmts(&ee, *scope, transform)

		for k, v := range xformres.GetMap().Items {
			(*scope)[k] = v
		}
		return ""
	}
	p.Expr()
	res = listener.TopExpr()

	// There should be a single Expression in the listener,
	// Take it and force it to a string for output

	expr := &sysl.Expr{
		Expr: &sysl.Expr_Unexpr{
			Unexpr: &sysl.Expr_UnExpr{
				Op:  sysl.Expr_UnExpr_STRING,
				Arg: res,
			},
		},
	}

	val, err := ee.eval(*scope, expr)
	if err != nil {
		return err.Error()
	}
	return val.GetS() + "\n"
}

func NewREPL(input io.Reader, output io.Writer) DebugFunc {
	r := repl{
		input:       bufio.NewScanner(input),
		output:      output,
		breakpoints: []*sysl.SourceContext{},
		step:        true,
	}
	return r.run
}

func writeOutput(text string, to io.Writer) {
	if _, err := io.WriteString(to, text); err != nil {
		logrus.Errorf("Failed to write to output. %s", err.Error())
	}
}
