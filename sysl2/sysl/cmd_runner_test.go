package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestEnsureFlagsNonEmpty_AllowsExcludes(t *testing.T) {
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := sysl.Command("foo", "")
	_ = cmd.Flag("bar", "").Default("foo").String()
	_ = cmd.Flag("other", "").Default("foo").String()

	EnsureFlagsNonEmpty(cmd, "bar")

	args := []string{"foo", "--bar", ""}
	selected, err := sysl.Parse(args)
	assert.Equal(t, "foo", selected)
	assert.NoError(t, err)
}

func TestEnsureFlagsNonEmpty(t *testing.T) {
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := sysl.Command("foo", "")
	cmd.Flag("bar", "").Default("foo")

	EnsureFlagsNonEmpty(cmd)

	args := []string{"foo", "--bar", ""}
	_, err := sysl.ParseContext(args)
	assert.Error(t, err)
}

func TestRootHandler(t *testing.T) {
	fs := afero.NewOsFs()
	logger := logrus.StandardLogger()
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	runner := cmdRunner{}
	runner.Configure(sysl)

	getAbsolute := func(path string) string {
		abs, _ := filepath.Abs(path)
		return abs
	}

	tests := []map[string]interface{}{
		{
			"root":   "",
			"module": "tests/root_finder_tests/NoSyslModuleTest/fake.sysl",
			"err":    func(err error) { assert.EqualError(t, err, "Sysl module not found") },
			"output": "",
		},
		{
			"root":   "",
			"module": "tests/root_finder_tests/RootFlagAndSyslRootUndefinedTest/test.sysl",
			"err":    func(err error) { assert.EqualError(t, err, "Project root is undefined") },
			"output": "",
		},
		{
			"root":   "",
			"module": "tests/root_finder_tests/SuccessfulTest/path/to/module/test.sysl",
			"err":    func(err error) { assert.NoError(t, err) },
			"output": getAbsolute("tests/root_finder_tests/SuccessfulTest"),
		},
		{
			"root":   "",
			"module": "tests/root_finder_tests/SuccessfulTest/test2.sysl",
			"err":    func(err error) { assert.NoError(t, err) },
			"output": getAbsolute("tests/root_finder_tests/SuccessfulTest"),
		},
		{
			"root":   "",
			"module": "tests/root_finder_tests/SuccessfulTest/path/to/another/module/test3.sysl",
			"err":    func(err error) { assert.NoError(t, err) },
			"output": getAbsolute("tests/root_finder_tests/SuccessfulTest/path/to/another"),
		},
		{
			"root":   getAbsolute("tests/root_finder_tests/DefinedRootAndSyslRootUndefinedTest/"),
			"module": "tests/root_finder_tests/DefinedRootAndSyslRootUndefinedTest/path/to/module/test.sysl",
			"err":    func(err error) { assert.NoError(t, err) },
			"output": getAbsolute("tests/root_finder_tests/DefinedRootAndSyslRootUndefinedTest/"),
		},
	}

	for i, test := range tests {
		t.Logf("Test %d\n", i)
		runner.Root = test["root"].(string)
		runner.module = test["module"].(string)
		absolutePath, _ := filepath.Abs(runner.module)
		t.Logf("%s\n", absolutePath)
		err := runner.rootHandler(fs, logger)

		t.Log(err)
		assert.Equal(t, test["output"], runner.Root)

		// Calling the appropriate error assertion
		test["err"].(func(err error))(err)

	}
}
