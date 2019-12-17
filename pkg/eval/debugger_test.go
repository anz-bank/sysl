package eval

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/sysl"

	"github.com/stretchr/testify/require"
)

func TestREPL_simpleCommands(t *testing.T) {
	type data struct {
		name     string
		command  string
		expected string
	}

	tests := []data{
		{"help", "?", "This help text"},
		{"step", "s", "(step)"},
		{"continue", "c", "(continue)"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := bytes.NewBufferString(tt.command + "\n")
			output := bytes.Buffer{}

			r := NewREPL(input, &output)
			r(&Scope{}, nil, nil) //nolint: errcheck
			require.Contains(t, output.String(), tt.expected)
		})
	}
}

func TestREPL_Expression(t *testing.T) {
	input := bytes.NewBufferString("let foo = 1\n f = foo\nf\n")
	output := bytes.Buffer{}

	scope := &Scope{}
	r := NewREPL(input, &output)
	r(scope, nil, nil) //nolint: errcheck

	result := output.String()
	require.Contains(t, result, "1\n")
}

func TestIsBreakpoint(t *testing.T) {
	type data struct {
		name     string
		bptext   string
		exprSC   string // will use makeBreakpoint() to get the SourceContext value
		expected bool
	}
	tests := []data{
		{"no source context", "foo.sysl:123", "", false},
		{"yes", "foo.sysl:123", "foo.sysl:123", true},
		{"yes-with-expr-col-set", "foo.sysl:123", "foo.sysl:123:32", true},
		{"no-with-expr-col-set", "foo.sysl:123:0", "foo.sysl:123:32", false},
		{"wrong-file", "foo.sysl:123", "bar.sysl:123:32", false},
		{"file-path-suffix", "foo.sysl:123", "this/is/a/path/foo.sysl:123:32", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := repl{
				breakpoints: []*sysl.SourceContext{
					makeBreakpoint(tt.bptext),
				},
			}
			expr := &sysl.Expr{SourceContext: makeBreakpoint(tt.exprSC)}
			require.Equal(t, tt.expected, r.isBreakpoint(expr))
		})
	}
}

func Test_makeBreakpoint(t *testing.T) {
	type data struct {
		name string
		text string
		// expects
		filename string
		line     int
		col      int
	}
	tests := []data{
		{
			"simple",
			"simple.sysl:100",
			"simple.sysl",
			100,
			-1,
		},
		{
			"with-col",
			"simple.sysl:100:30",
			"simple.sysl",
			100,
			30,
		},
		{
			"missing-line",
			"simple.sysl",
			"",
			0,
			0,
		},
		{
			"missing-line-conv",
			"simple.sysl:aaa",
			"",
			0,
			0,
		},
		{
			"missing-col-conv",
			"simple.sysl:123:aaa",
			"",
			0,
			0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			bp := makeBreakpoint(tt.text)
			if tt.filename == "" {
				require.Nil(t, bp)
			} else {
				require.NotNil(t, bp)
				require.Equal(t, tt.filename, bp.File)
				require.EqualValues(t, tt.line, bp.Start.Line)
				require.EqualValues(t, tt.col, bp.Start.Col)
			}
		})
	}
}

func TestREPL_DumpScope(t *testing.T) {
	input := bytes.NewBufferString("d\n")
	output := bytes.Buffer{}

	scope := &Scope{
		"a": MakeValueString("hello"),
		"b": MakeValueBool(false),
		"c": makeNullValue(),
	}
	r := NewREPL(input, &output)
	r(scope, nil, nil) //nolint: errcheck

	result := output.String()
	require.Contains(t, result, "(dump scope)\n")
	for k, v := range *scope {
		require.Contains(t, result, fmt.Sprintf("%s:  %s", k, unaryString(v).GetS()))
	}
}

func TestREPL_CallView(t *testing.T) {
	filename := testDir + "transform1.sysl"
	sysl, appname, err := parse.LoadAndGetDefaultApp(filename, afero.NewOsFs(), parse.NewParser())
	require.NoError(t, err)

	scope := &Scope{}
	scope.AddApp("app", sysl.Apps[appname])
	input := bytes.NewBufferString("l\nTfmFilenameInvalid2(app).foo\n")
	output := bytes.Buffer{}

	r := NewREPL(input, &output)
	r(scope, sysl.Apps[appname], nil) //nolint: errcheck
	result := output.String()
	require.Contains(t, result, "servicehandler.go\n")
}
