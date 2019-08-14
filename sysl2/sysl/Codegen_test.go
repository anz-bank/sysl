package main

import (
	"bytes"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/eval"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCode(t *testing.T) {
	output := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	root := output[0].output
	assert.Equal(t, 1, len(output), "unexpected length of output")
	assert.Equal(t, 1, len(root), "unexpected length of root")
	n1 := root[0].(Node)
	assert.Equal(t, 4, len(n1), "unexpected length of javaFile")
	package1 := n1[0].(Node)
	comment1 := n1[1].(Node)
	import1 := n1[2].(Node)
	definition1 := n1[3].(string)
	assert.Equal(t, 1, len(package1), "unexpected length of package")
	assert.Equal(t, 2, len(comment1), "unexpected length of comment")
	assert.Equal(t, 2, len(import1), "unexpected length of import")
	assert.Equal(t, "some_value", definition1, "unexpected value of definition")

	package2 := package1[0].(Node)
	assert.Equal(t, 3, len(package2), "unexpected length of package2")
	assert.Equal(t, "com.example.gen", package2[1].(string), "unexpected length of package2")

	for i, comment := range []string{"comment1", "comment2"} {
		comment0 := comment1[i].(Node)
		assert.Equal(t, 1, len(comment0), "unexpected length of comment2")
		comment0_0 := comment0[0].(string)
		assert.Equal(t, comment, comment0_0, "unexpected length of comment_i")
	}

	for i, imports := range []string{"import1", "import2"} {
		import0 := import1[i].(Node)
		assert.Equal(t, 1, len(import0), "unexpected length of import2")
		import0_0 := import0[0].(Node)
		assert.Equal(t, 3, len(import0_0), "unexpected length of import2")
		assert.Equal(t, imports, import0_0[1].(string), "unexpected length of import_i")
	}
}

func TestGenerateCodeNoComment(t *testing.T) {
	output := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen_no_comment.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	assert.Equal(t, 1, len(output), "unexpected length of output")
	root := output[0].output
	assert.Equal(t, 1, len(root), "unexpected length of root")
	n1 := root[0].(Node)
	assert.Equal(t, 3, len(n1), "unexpected length of javaFile")
	package1 := n1[0].(Node)
	import1 := n1[1].(Node)
	definition1 := n1[2].(string)
	assert.Equal(t, 1, len(package1), "unexpected length of package")
	assert.Equal(t, 2, len(import1), "unexpected length of comment")
	assert.Equal(t, "some_value", definition1, "unexpected value of definition")

	package2 := package1[0].(Node)
	assert.Equal(t, 3, len(package2), "unexpected length of package2")
	assert.Equal(t, "com.example.gen", package2[1].(string), "unexpected length of package2")

	for i, imports := range []string{"import1", "import2"} {
		import0 := import1[i].(Node)
		assert.Equal(t, 1, len(import0), "unexpected length of import2")
		import0_0 := import0[0].(Node)
		assert.Equal(t, 3, len(import0_0), "unexpected length of import2")
		assert.Equal(t, imports, import0_0[1].(string), "unexpected length of import_i")
	}
}

func TestGenerateCodeNoPackage(t *testing.T) {
	output := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen_no_package.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	root := output[0].output
	assert.Nil(t, root, "unexpected root")
}

func TestGenerateCodeMultipleAnnotations(t *testing.T) {
	output := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen_multiple_annotations.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	root := output[0].output
	assert.Nil(t, root, "unexpected root")
}

func TestGenerateCodePerType(t *testing.T) {
	output := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/multiple_file.gen.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	assert.Equal(t, 1, len(output), "unexpected length of output")
	assert.Equal(t, "Request.java", output[0].filename, "unexpected length of output")

	root := output[0].output
	assert.Equal(t, 1, len(root), "unexpected length of javaFile")

	requestRoot := root[0].(Node)
	assert.Equal(t, 4, len(requestRoot), "unexpected length of requestRoot")

	package1 := requestRoot[0].(Node)
	comment1 := requestRoot[1].(Node)
	import1 := requestRoot[2].(Node)
	definition1 := requestRoot[3].(string)
	assert.Equal(t, 1, len(package1), "unexpected length of package")
	assert.Equal(t, 2, len(comment1), "unexpected length of comment")
	assert.Equal(t, 2, len(import1), "unexpected length of import")
	assert.Equal(t, "Request", definition1, "unexpected value of definition")
}

func TestSerialize(t *testing.T) {
	output := GenerateCodeWithParams(".", "tests/model.sysl", ".", "tests/test.gen.sysl",
		"tests/test.gen.g", "javaFile", "warn", false)
	out := new(bytes.Buffer)
	require.NoError(t, Serialize(out, " ", output[0].output))
	golden := "package com.example.gen \n comment1 comment2 import import1 \n import import2 \n some_value "
	assert.Equal(t, golden, out.String(), "unexpected value of out string")
}

func TestOutputForPureTokenOnlyRule(t *testing.T) {
	g := readGrammar("tests/token_only_rule.g", "gen", "pureToken")
	obj := eval.MakeValueMap()
	m := eval.MakeValueMap()
	eval.AddItemToValueMap(m, "text", eval.MakeValueString("hello"))
	eval.AddItemToValueMap(obj, "header", eval.MakeValueBool(true))
	eval.AddItemToValueMap(obj, "tail", eval.MakeValueBool(true))
	eval.AddItemToValueMap(obj, "body", m)
	output := processRule(g, obj, "pureToken")
	assert.NotNil(t, output)

	root := output[0].(Node)
	assert.NotNil(t, root)
	assert.Equal(t, 3, len(root))

	header := root[0].(Node)
	assert.Equal(t, 1, len(header))
	assert.Equal(t, "head", header[0].(Node)[0].(string))

	body := root[1].(Node)
	assert.Equal(t, 1, len(body))
	assert.Equal(t, "hello", body[0].(Node)[1].(string))

	tail := root[2].(Node)
	assert.Equal(t, 1, len(tail))
	assert.Equal(t, "tail", tail[0].(Node)[0].(string))
}

func GenerateCodeWithParams(rootModel, model, rootTransform, transform, grammar, start string, loglevel string,
	isVerbose bool,
) []*CodeGenOutput {
	cmdContextParamCodegen := &CmdContextParamCodegen{
		rootModel:     &rootModel,
		model:         &model,
		rootTransform: &rootTransform,
		transform:     &transform,
		grammar:       &grammar,
		start:         &start,
		loglevel:      &loglevel,
		isVerbose:     &isVerbose,
	}
	return GenerateCode(cmdContextParamCodegen)
}
