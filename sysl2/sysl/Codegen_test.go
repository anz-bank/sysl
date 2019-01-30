package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCode(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", "tests/test.gen.sysl", "tests/test.gen.g", "javaFile")
	root := output[0].output
	assert.True(t, len(output) == 1, "unexpected length of output")
	assert.True(t, len(root) == 1, "unexpected length of root")
	n1 := root[0].(Node)
	assert.True(t, len(n1) == 4, "unexpected length of javaFile")
	package1 := n1[0].(Node)
	comment1 := n1[1].(Node)
	import1 := n1[2].(Node)
	definition1 := n1[3].(string)
	assert.True(t, len(package1) == 1, "unexpected length of package")
	assert.True(t, len(comment1) == 2, "unexpected length of comment")
	assert.True(t, len(import1) == 2, "unexpected length of import")
	assert.True(t, definition1 == "some_value", "unexpected value of definition")

	package2 := package1[0].(Node)
	assert.True(t, len(package2) == 3, "unexpected length of package2")
	assert.True(t, package2[1].(string) == "com.example.gen", "unexpected length of package2")

	comment := []string{"comment1", "comment2"}
	for i := range comment {
		comment_0 := comment1[i].(Node)
		assert.True(t, len(comment_0) == 1, "unexpected length of comment2")
		comment_0_0 := comment_0[0].(string)
		assert.True(t, comment_0_0 == comment[i], "unexpected length of comment_i")
	}

	imports := []string{"import1", "import2"}
	for i := range imports {
		import_0 := import1[i].(Node)
		assert.True(t, len(import_0) == 1, "unexpected length of import2")
		import_0_0 := import_0[0].(Node)
		assert.True(t, len(import_0_0) == 3, "unexpected length of import2")
		assert.True(t, import_0_0[1].(string) == imports[i], "unexpected length of import_i")
	}
}

func TestGenerateCodeNoComment(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", "tests/test.gen_no_comment.sysl", "tests/test.gen.g", "javaFile")
	assert.True(t, len(output) == 1, "unexpected length of output")

	root := output[0].output
	assert.True(t, len(root) == 1, "unexpected length of root")
	n1 := root[0].(Node)
	assert.True(t, len(n1) == 3, "unexpected length of javaFile")
	package1 := n1[0].(Node)
	import1 := n1[1].(Node)
	definition1 := n1[2].(string)
	assert.True(t, len(package1) == 1, "unexpected length of package")
	assert.True(t, len(import1) == 2, "unexpected length of comment")
	assert.True(t, definition1 == "some_value", "unexpected value of definition")

	package2 := package1[0].(Node)
	assert.True(t, len(package2) == 3, "unexpected length of package2")
	assert.True(t, package2[1].(string) == "com.example.gen", "unexpected length of package2")

	imports := []string{"import1", "import2"}
	for i := range imports {
		import_0 := import1[i].(Node)
		assert.True(t, len(import_0) == 1, "unexpected length of import2")
		import_0_0 := import_0[0].(Node)
		assert.True(t, len(import_0_0) == 3, "unexpected length of import2")
		assert.True(t, import_0_0[1].(string) == imports[i], "unexpected length of import_i")
	}
}

func TestGenerateCodeNoPackage(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", "tests/test.gen_no_package.sysl", "tests/test.gen.g", "javaFile")
	root := output[0].output
	assert.True(t, root == nil, "unexpected length of root")
}

func TestGenerateCodeMultipleAnnotations(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", "tests/test.gen_multiple_annotations.sysl", "tests/test.gen.g", "javaFile")
	root := output[0].output
	assert.True(t, root == nil, "unexpected length of root")
}

func TestGenerateCodePerType(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", "tests/multiple_file.gen.sysl", "tests/test.gen.g", "javaFile")
	assert.True(t, len(output) == 1, "unexpected length of output")
	assert.True(t, output[0].filename == "Request.java", "unexpected length of output")

	root := output[0].output
	assert.True(t, len(root) == 1, "unexpected length of javaFile")

	requestRoot := root[0].(Node)
	assert.True(t, len(requestRoot) == 4, "unexpected length of requestRoot")

	package1 := requestRoot[0].(Node)
	comment1 := requestRoot[1].(Node)
	import1 := requestRoot[2].(Node)
	definition1 := requestRoot[3].(string)
	assert.True(t, len(package1) == 1, "unexpected length of package")
	assert.True(t, len(comment1) == 2, "unexpected length of comment")
	assert.True(t, len(import1) == 2, "unexpected length of import")
	assert.True(t, definition1 == "Request", "unexpected value of definition")
}

func TestSerialize(t *testing.T) {
	output := GenerateCode(".", "tests/model.sysl", "tests/test.gen.sysl", "tests/test.gen.g", "javaFile")
	out := new(bytes.Buffer)
	Serialize(out, " ", output[0].output)
	golden := "package com.example.gen \n comment1 comment2 import import1 \n import import2 \n some_value "
	val := out.String()
	assert.True(t, val == golden, "unexpected value of out string")
}
