package main

import (
	"bytes"
	"context"
	"os"
	"path"
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
	transformScript = `\input (output: input.models(0).rel.app => .appName(1))`
	transformOutput = "{'App1', 'App2'}\n"
)

func TestTransformInlineScript(t *testing.T) {
	t.Parallel()

	// FIXME: windows CI test fails on
	// "Rel: can't make \input (output: input.models(0).rel.app => .appName(1)) relative to ."
	if runtime.GOOS == "windows" {
		return
	}

	output := runSyslWithOutput(t, ".sysl", nil,
		"transform", "../../tests/simple.sysl", "--script", transformScript)
	assert.Equal(t, transformOutput, output)
}

func TestTransformTextFile(t *testing.T) {
	t.Parallel()
	scriptPath := writeTempFile(t, "transform.arrai", []byte(transformScript))
	output := runSyslWithOutput(t, ".sysl", nil,
		"transform", "../../tests/simple.sysl", "--script", scriptPath)
	assert.Equal(t, transformOutput, output)
}

func TestTransformBundleFile(t *testing.T) {
	t.Parallel()

	// FIXME: windows CI test times out in this test.
	if runtime.GOOS == "windows" {
		return
	}

	bundlePath := writeTempFile(t, "transform.arraiz", createBundle(t, transformScript))
	output := runSyslWithOutput(t, ".sysl", nil,
		"transform", "../../tests/simple.sysl", "--script", bundlePath)
	assert.Equal(t, transformOutput, output)
}

func TestTransformTextStdin(t *testing.T) {
	t.Parallel()

	output := runSyslWithOutput(t, ".sysl", strings.NewReader(transformScript),
		"transform", "../../tests/simple.sysl", "--script=-")
	assert.Equal(t, transformOutput, output)
}

func TestTransformBundleStdin(t *testing.T) {
	t.Parallel()

	output := runSyslWithOutput(t, ".sysl", bytes.NewReader(createBundle(t, transformScript)),
		"transform", "../../tests/simple.sysl", "--script=-")
	assert.Equal(t, transformOutput, output)
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
