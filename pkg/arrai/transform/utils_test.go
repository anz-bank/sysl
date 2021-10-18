package transform

import (
	"context"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/arr-ai/arrai/pkg/bundle"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestTransformWithArraiz(t *testing.T) {
	type cases struct {
		files                         map[string]string
		script, param, expected, name string
	}
	testCases := []cases{
		{
			files: map[string]string{
				"/simple.arrai": `\_ 1`,
			},
			script:   "/simple.arrai",
			param:    `0`,
			expected: `1`,
			name:     "simple",
		},
		{
			files: map[string]string{
				"/param.arrai": `\param param`,
			},
			script:   "/param.arrai",
			param:    `0`,
			expected: `0`,
			name:     "param",
		},
		{
			files: map[string]string{
				"/import.arrai": `\param //{./file} + param`,
				"/file.arrai":   `1`,
			},
			script:   "/import.arrai",
			param:    `1`,
			expected: `2`,
			name:     "import",
		},
	}
	for _, c := range testCases {
		t.Run(
			c.name,
			func(t *testing.T) {
				t.Parallel()
				assertValueOfTransformedArraiz(t, c.files, c.script, c.param, c.expected)
			},
		)
	}
}

func TestTransformWithArraizErrors(t *testing.T) {
	t.Parallel()

	assertErrorOfTransformedArraiz(t,
		map[string]string{
			"/number.arrai": `1`,
		},
		"/number.arrai",
		`0`,
		errNotClosure.Error(),
	)
}

func assertValueOfTransformedArraiz(
	t *testing.T,
	files map[string]string,
	mainScript, param, expected string,
) {
	paramVal, err := syntax.EvalWithScope(context.Background(), "", param, rel.EmptyScope)
	require.NoError(t, err)
	val, err := getTransformArraizResult(t, files, mainScript, paramVal)
	require.NoError(t, err)
	syntax.AssertEvalExprString(t, expected, val.String())
}

func assertErrorOfTransformedArraiz(
	t *testing.T,
	files map[string]string,
	mainScript, param, expectedErrMsg string,
) {
	paramVal, err := syntax.EvalWithScope(context.Background(), "", param, rel.EmptyScope)
	require.NoError(t, err)
	_, err = getTransformArraizResult(t, files, mainScript, paramVal)
	assert.EqualError(t, err, expectedErrMsg)
}

func getTransformArraizResult(
	t *testing.T,
	files map[string]string,
	mainScript string,
	param rel.Value,
) (rel.Value, error) {
	fs := afero.NewMemMapFs()
	scriptPath := "/test.arraiz"

	f, err := fs.Create(scriptPath)
	require.NoError(t, err)

	_, err = f.Write(bundle.MustCreateTestBundleFromMap(t, files, mainScript))
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)

	scriptBytes, err := afero.ReadFile(fs, scriptPath)
	if err != nil {
		return nil, err
	}

	return EvalWithParam(scriptBytes, scriptPath, param)
}
