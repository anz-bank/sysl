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
	err := runner.Configure(sysl)
	assert.NoError(t, err)

	tests := []map[string]interface{}{
		{
			"root":           "",
			"module":         "tests/root_finder_tests/SuccessfulTest/path/to/module/test.sysl",
			"relativeOutput": "tests/root_finder_tests/SuccessfulTest",
			"error":          func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			"root":           "",
			"module":         "tests/root_finder_tests/SuccessfulTest/test2.sysl",
			"relativeOutput": "tests/root_finder_tests/SuccessfulTest",
			"error":          func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			"root":           "",
			"module":         "tests/root_finder_tests/SuccessfulTest/path/to/another/module/test3.sysl",
			"relativeOutput": "tests/root_finder_tests/SuccessfulTest/path/to/another",
			"error":          func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			"root":           "./tests/root_finder_tests/DefinedRootAndSyslRootUndefinedTest/",
			"module":         "path/to/module/test.sysl",
			"relativeOutput": "./tests/root_finder_tests/DefinedRootAndSyslRootUndefinedTest/",
			"error":          func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			"root":           ".",
			"module":         "path/to/module/test.sysl",
			"relativeOutput": ".",
			"error":          func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			"root":           "",
			"module":         "tests/root_finder_tests/UndefinedRootAndUndefinedSyslRoot/test.sysl",
			"relativeOutput": "",
			"error":          func(t *testing.T, err error) { assert.EqualError(t, err, "root is not defined") },
		},
	}

	for i, test := range tests {
		runner.Root = test["root"].(string)
		runner.module = test["module"].(string)
		realAbsOutput, err := filepath.Abs(test["root"].(string))
		assert.NoError(t, err)

		err = runner.rootHandler(fs, logger)

		t.Logf("Test #%d, root %s", i, realAbsOutput)
		assert.Equal(t, test["relativeOutput"], runner.Root)

		// testing the appropriate error assertion
		test["error"].(func(t *testing.T, err error))(t, err)
	}
}
