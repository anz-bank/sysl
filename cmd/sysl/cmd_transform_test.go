package main

import (
	"bytes"
	"context"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/pkg/ctxfs"
	"github.com/arr-ai/arrai/pkg/ctxrootcache"
	"github.com/arr-ai/arrai/pkg/importcache"
	"github.com/arr-ai/arrai/syntax"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSyslPath          = "../../tests/simple.sysl"
	transformScript       = `\input (output: input.models(0).rel.app => .appName(1))`
	transformImportScript = `\input (output: input.models(0).rel.app => //{./util}.joinAppName(.appName))`
	transformOutput       = "{'App1', 'App2'}\n"
	transformImportOutput = "{'Namespace1 :: App1', 'Namespace1 :: App2'}\n"
	win                   = "windows"
)

func TestTransformInlineScript(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == win {
		t.Skip("This test fails on Windows with the error: " +
			"Rel: can't make \\input (output: input.models(0).rel.app => .appName(1)) relative to .")
	}

	output := runSyslWithOutput(t, ".sysl", nil,
		"transform", testSyslPath, "--script", transformScript)
	assert.Equal(t, transformOutput, output)
}

func TestTransformInlineScriptWithStdinModules(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == win {
		t.Skip("This test fails on Windows with the error: " +
			"Rel: can't make \\input (output: input.models(0).rel.app => .appName(1)) relative to .")
	}

	src, err := os.ReadFile(testSyslPath)
	require.NoError(t, err)
	stdin := toStdin(t, stdinFile{Path: testSyslPath, Content: string(src)})

	output := runSyslWithOutput(t, ".sysl", stdin,
		"transform", "--script", transformScript)
	assert.Equal(t, transformOutput, output)
}

func TestTransformTextFile(t *testing.T) {
	t.Parallel()

	scriptPath := writeTempFile(t, "transform.arrai", []byte(transformScript))
	output := runSyslWithOutput(t, ".sysl", nil,
		"transform", testSyslPath, "--script", scriptPath)
	assert.Equal(t, transformOutput, output)
}

func TestTransformTextFile_moduleStdin(t *testing.T) {
	t.Parallel()

	src, err := os.ReadFile(testSyslPath)
	require.NoError(t, err)
	scriptPath := writeTempFile(t, "transform.arrai", []byte(transformScript))
	stdin := toStdin(t, stdinFile{Path: testSyslPath, Content: string(src)})
	output := runSyslWithOutput(t, ".sysl", stdin, "transform", "--script", scriptPath)
	assert.Equal(t, transformOutput, output)
}

func TestTransformBundleFile(t *testing.T) {
	t.Parallel()

	// FIXME: windows CI test times out in this test.
	if runtime.GOOS == win {
		return
	}

	bundlePath := writeTempFile(t, "transform.arraiz", createBundle(t, transformScript))
	output := runSyslWithOutput(t, ".sysl", nil,
		"transform", testSyslPath, "--script", bundlePath)
	assert.Equal(t, transformOutput, output)
}

func TestTransformTextStdin(t *testing.T) {
	t.Parallel()

	output := runSyslWithOutput(t, ".sysl", strings.NewReader(transformScript),
		"transform", testSyslPath, "--script=-")
	assert.Equal(t, transformOutput, output)
}

func TestTransformTextStdin_relImport(t *testing.T) {
	t.Skip("relative imports not yet supported by stdin transform")
	t.Parallel()

	scriptPath := filepath.Join("..", "..", "pkg", "arrai", "foo.arrai")
	stdin := toStdin(t, stdinFile{Path: scriptPath, Content: transformImportScript})
	output := runSyslWithOutput(t, ".sysl", stdin,
		"transform", testSyslPath, "--script=-")
	assert.Equal(t, transformImportOutput, output)
}

func TestTransformBundleStdin(t *testing.T) {
	t.Parallel()

	output := runSyslWithOutput(t, ".sysl", bytes.NewReader(createBundle(t, transformScript)),
		"transform", testSyslPath, "--script=-")
	assert.Equal(t, transformOutput, output)
}

const simple = `Namespace1 :: App1 [~appTag]:
    !table Table1 [~tableTag]:
        Column2 <: string [~column2Tag]

    !table View1 [~view]:
        Column1 <: int [~vColumn1Tag]

    !type WebService1 [~typeTag]:
        Field1 <: datetime [~field1TagA, ~field1TagB]
        Field2 <: bool [~field2Tag]


Namespace1 :: App2:
    !table Table1:
        Column1 <: float [~pk]
        Column2 <: sequence of string
`

// Round trip tests the fidelity of round-tripping Sysl into a relational model and back via reconstruct,
// and expects to get the exact same Sysl as it was given.
func TestTransformRoundTripSimple(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == win {
		t.Skip("This test fails on Windows with the error: " +
			"Rel: can't make \\input (output: input.models(0).rel.app => .appName(1)) relative to .")
	}

	stdin := toStdin(t, stdinFile{Path: "simple.sysl", Content: simple})
	output := runSyslWithOutput(t, ".sysl", stdin,
		"transform", "--script",
		`\i (output: //{/pkg/arrai/reconstruct.arrai}(i.models(0).rel)('cmd')('sysl')('simple.sysl'))`)
	assert.Equal(t, simple, output)
}

func TestTransformRoundTripAll(t *testing.T) {
	t.Skip("TODO: Fix reconstruct to work with all.sysl")
	t.Parallel()

	if runtime.GOOS == win {
		t.Skip("This test fails on Windows with the error: " +
			"Rel: can't make \\input (output: input.models(0).rel.app => .appName(1)) relative to .")
	}

	runSyslWithExpectedOutput(t, "pkg/arrai/all.sysl",
		"transform", "../../pkg/arrai/all.sysl", "--script",
		`\i (output: //{/pkg/arrai/reconstruct.arrai}(i.models(0).rel)('cmd')('sysl')('all.sysl'))`)
}

func writeTempFile(t *testing.T, filename string, data []byte) string {
	filePath := path.Join(t.TempDir(), filename)
	err := os.WriteFile(filePath, data, 0600)
	require.NoError(t, err)
	return filePath
}

func createBundle(t *testing.T, script string) []byte {
	ctx := ctxfs.SourceFsOnto(context.Background(), afero.NewOsFs())
	ctx = ctxrootcache.WithRootCache(ctx)
	ctx, err := syntax.SetupBundle(ctx, "", []byte(script))
	require.NoError(t, err)

	_, err = syntax.Compile(importcache.WithNewImportCache(ctx), "", script)
	require.NoError(t, err)

	b := &bytes.Buffer{}
	err = syntax.OutputArraiz(ctx, b)
	require.NoError(t, err)

	return b.Bytes()
}
