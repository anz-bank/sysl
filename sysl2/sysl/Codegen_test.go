package main

import (
	"bytes"
	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/eval"
	"github.com/anz-bank/sysl/sysl2/sysl/testutil"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	output, err := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	require.NoError(t, err)
	root := output[0].output
	assert.Len(t, output, 1)
	assert.Len(t, root, 1)
	n1 := root[0].(Node)
	assert.Len(t, n1, 4)
	package1 := n1[0].(Node)
	comment1 := n1[1].(Node)
	import1 := n1[2].(Node)
	definition1 := n1[3].(string)
	assert.Len(t, package1, 1)
	assert.Len(t, comment1, 2)
	assert.Len(t, import1, 2)
	assert.Equal(t, "some_value", definition1)

	package2 := package1[0].(Node)
	assert.Len(t, package2, 3)
	assert.Equal(t, "com.example.gen", package2[1].(string))

	for i, comment := range []string{"comment1", "comment2"} {
		comment0 := comment1[i].(Node)
		assert.Len(t, comment0, 1)
		comment0_0 := comment0[0].(string)
		assert.Equal(t, comment, comment0_0)
	}

	for i, imports := range []string{"import1", "import2"} {
		import0 := import1[i].(Node)
		assert.Len(t, import0, 1)
		import0_0 := import0[0].(Node)
		assert.Len(t, import0_0, 3)
		assert.Equal(t, imports, import0_0[1].(string))
	}
}

func TestGenerateCodeNoComment(t *testing.T) {
	t.Parallel()

	output, err := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen_no_comment.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	require.NoError(t, err)
	assert.Len(t, output, 1)
	root := output[0].output
	assert.Len(t, root, 1)
	n1 := root[0].(Node)
	assert.Len(t, n1, 3)
	package1 := n1[0].(Node)
	import1 := n1[1].(Node)
	definition1 := n1[2].(string)
	assert.Len(t, package1, 1)
	assert.Len(t, import1, 2)
	assert.Equal(t, "some_value", definition1)

	package2 := package1[0].(Node)
	assert.Len(t, package2, 3)
	assert.Equal(t, "com.example.gen", package2[1].(string))

	for i, imports := range []string{"import1", "import2"} {
		import0 := import1[i].(Node)
		assert.Len(t, import0, 1)
		import0_0 := import0[0].(Node)
		assert.Len(t, import0_0, 3)
		assert.Equal(t, imports, import0_0[1].(string))
	}
}

func TestGenerateCodeNoPackage(t *testing.T) {
	t.Parallel()

	output, err := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen_no_package.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	require.NoError(t, err)
	root := output[0].output
	assert.Nil(t, root)
}

func TestGenerateCodeMultipleAnnotations(t *testing.T) {
	t.Parallel()

	output, err := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen_multiple_annotations.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	require.NoError(t, err)
	root := output[0].output
	assert.Nil(t, root)
}

func TestGenerateCodePerType(t *testing.T) {
	t.Parallel()

	output, err := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/multiple_file.gen.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	require.NoError(t, err)
	assert.Len(t, output, 1)
	assert.Equal(t, "Request.java", output[0].filename)

	root := output[0].output
	assert.Len(t, root, 1)

	requestRoot := root[0].(Node)
	assert.Len(t, requestRoot, 4)

	package1 := requestRoot[0].(Node)
	comment1 := requestRoot[1].(Node)
	import1 := requestRoot[2].(Node)
	definition1 := requestRoot[3].(string)
	assert.Len(t, package1, 1)
	assert.Len(t, comment1, 2)
	assert.Len(t, import1, 2)
	assert.Equal(t, "Request", definition1)
}

func TestSerialize(t *testing.T) {
	t.Parallel()

	output, err := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	require.NoError(t, err)
	out := new(bytes.Buffer)
	require.NoError(t, Serialize(out, " ", output[0].output))
	golden := "package com.example.gen \n comment1 comment2 import import1 \n import import2 \n some_value "
	assert.Equal(t, golden, out.String())
}

func TestOutputForPureTokenOnlyRule(t *testing.T) {
	t.Parallel()

	g, err := readGrammar("tests/token_only_rule.g", "gen", "pureToken")
	require.NoError(t, err)
	obj := eval.MakeValueMap()
	m := eval.MakeValueMap()
	eval.AddItemToValueMap(m, "text", eval.MakeValueString("hello"))
	eval.AddItemToValueMap(obj, "header", eval.MakeValueBool(true))
	eval.AddItemToValueMap(obj, "tail", eval.MakeValueBool(true))
	eval.AddItemToValueMap(obj, "body", m)
	logger, _ := test.NewNullLogger()
	output := processRule(g, obj, "pureToken", logger)
	assert.NotNil(t, output)

	root := output[0].(Node)
	assert.NotNil(t, root)
	assert.Len(t, root, 3)

	header := root[0].(Node)
	assert.Len(t, header, 1)
	assert.Equal(t, "head", header[0].(Node)[0].(string))

	body := root[1].(Node)
	assert.Len(t, body, 1)
	assert.Equal(t, "hello", body[0].(Node)[1].(string))

	tail := root[2].(Node)
	assert.Len(t, tail, 1)
	assert.Equal(t, "tail", tail[0].(Node)[0].(string))
}

func GenerateCodeWithParams(
	rootModel, model, rootTransform, transform, grammar, start, loglevel string,
	isVerbose bool,
) ([]*CodeGenOutput, error) {
	_, fs := testutil.WriteToMemOverlayFs(rootModel)
	return GenerateCodeWithParamsFs(
		rootModel, model, rootTransform, transform, grammar, start, loglevel,
		isVerbose, fs,
	)
}

func GenerateCodeWithParamsFs(
	rootModel, model, rootTransform, transform, grammar, start, loglevel string,
	isVerbose bool, fs afero.Fs,
) ([]*CodeGenOutput, error) {

	modelParser := parse.NewParser()
	mod, modelAppName, err := parse.LoadAndGetDefaultApp(model, syslutil.NewChrootFs(fs, rootModel), modelParser)
	if err != nil {
		return nil, err
	}
	cmdContextParamCodegen := &CmdContextParamCodegen{
		model:         mod,
		modelAppName:  modelAppName,
		rootTransform: &rootTransform,
		transform:     &transform,
		grammar:       &grammar,
		start:         &start,
		outDir:        nil,
	}
	logger, _ := test.NewNullLogger()
	return GenerateCode(cmdContextParamCodegen, fs, logger)
}
