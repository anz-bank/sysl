package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

const cmdName = "foo"

const emptyAppSrc = `
Stdin::App:
	...`

const importSrc = `
import ../../tests/simple.sysl

Stdin::App:
	...`

func TestLoadModule_path(t *testing.T) {
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

func TestLoadModule_stdin(t *testing.T) {
	var mods []*sysl.Module
	cmd := &testCmd{name: cmdName, maxSyslModule: 1, exec: func(args cmdutils.ExecuteArgs) {
		mods = args.Modules
	}}
	r := buildRunner(t, cmd)

	stdin := toStdin(t, stdinFile{Path: filepath.Join(cwd(), "test.sysl"), Content: emptyAppSrc})
	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin))

	assert.Len(t, mods, 1)
	assert.NotEmpty(t, mods[0].Apps["Stdin :: App"])
}

func TestLoadModule_stdinImport(t *testing.T) {
	var mods []*sysl.Module
	cmd := &testCmd{name: cmdName, maxSyslModule: 1, exec: func(args cmdutils.ExecuteArgs) {
		mods = args.Modules
	}}
	r := buildRunner(t, cmd)

	stdin := toStdin(t, stdinFile{Path: filepath.Join(cwd(), "test.sysl"), Content: importSrc})
	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin))

	assert.Len(t, mods, 1)
	assert.Len(t, mods[0].Apps, 3)
	for _, a := range []string{"Stdin :: App", "Namespace1 :: App1", "Namespace1 :: App2"} {
		assert.Contains(t, mods[0].Apps, a)
	}
}

func TestLoadModule_stdinPathImport(t *testing.T) {
	var mods []*sysl.Module
	cmd := &testCmd{name: cmdName, maxSyslModule: 1, exec: func(args cmdutils.ExecuteArgs) {
		mods = args.Modules
	}}
	r := buildRunner(t, cmd)

	stdin := toStdin(t, stdinFile{Path: filepath.Join("..", "..", "tests", "test.sysl"), Content: "import simple.sysl"})
	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin))

	assert.Len(t, mods, 1)
	assert.Len(t, mods[0].Apps, 2)
	for _, a := range []string{"Namespace1 :: App1", "Namespace1 :: App2"} {
		assert.Contains(t, mods[0].Apps, a)
	}
}

func TestLoadModule_pathsAndStdin(t *testing.T) {
	args := []string{"simple.sysl"}
	cmd := &testCmd{name: cmdName, maxSyslModule: 1, exec: func(args cmdutils.ExecuteArgs) {}}
	r := buildRunnerWithRoot(t, cmd, testDir, args...)

	stdin := toStdin(t, stdinFile{Path: filepath.Join("test.sysl"), Content: emptyAppSrc})
	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin))
}

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

type testCmd struct {
	name          string
	maxSyslModule int
	exec          func(cmdutils.ExecuteArgs)
	modules       []string
}

func (c *testCmd) Execute(args cmdutils.ExecuteArgs) error {
	c.exec(args)
	return nil
}

func (c *testCmd) Name() string {
	return c.name
}

func (c *testCmd) MaxSyslModule() int {
	return c.maxSyslModule
}

func (c *testCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	return app.Command(c.name, "")
}

func buildRunner(t *testing.T, cmd *testCmd, args ...string) *cmdRunner {
	return buildRunnerWithRoot(t, cmd, "", args...)
}

func buildRunnerWithRoot(t *testing.T, cmd *testCmd, root string, args ...string) *cmdRunner {
	app := kingpin.New("sysl", "System Modelling Language Toolkit")

	r := &cmdRunner{
		modules:  cmd.modules,
		Root:     root,
		commands: map[string]cmdutils.Command{cmd.name: cmd},
	}
	r.ConfigureCmd(app, cmd)

	_, err := app.Parse(append([]string{cmd.name}, args...))
	require.NoError(t, err)
	return r
}

func toStdin(t *testing.T, files ...stdinFile) io.Reader {
	bs, err := json.Marshal(files)
	require.NoError(t, err)
	return bytes.NewReader(bs)
}

func cwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
