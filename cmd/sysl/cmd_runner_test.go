package main

import (
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

const emptyAppSrc = `
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
	cmdName := "foo"
	var mods []*sysl.Module
	cmd := &testCmd{name: cmdName, maxSyslModule: 1, exec: func(args cmdutils.ExecuteArgs) {
		mods = args.Modules
	}}
	r := buildRunner(t, cmd)

	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), strings.NewReader(emptyAppSrc)))

	assert.Len(t, mods, 1)
	assert.NotEmpty(t, mods[0].Apps["Stdin :: App"])
}

func TestLoadModule_pathsAndStdin(t *testing.T) {
	cmdName := "foo"
	args := []string{"simple.sysl"}
	var mods []*sysl.Module
	cmd := &testCmd{name: cmdName, maxSyslModule: 1, exec: func(args cmdutils.ExecuteArgs) {
		mods = args.Modules
	}}
	r := buildRunnerWithRoot(t, cmd, testDir, args...)

	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), strings.NewReader(emptyAppSrc)))

	assert.Len(t, mods, 1)
	// The "App" input on stdin is ignored since inputs were provided.
	assert.NotEmpty(t, mods[0].Apps["Namespace1 :: App1"])
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
	cmd := app.Command(c.name, "")
	if c.maxSyslModule > 0 {
		cmd.Arg("MODULE", "").StringsVar(&c.modules)
	}
	return cmd
}

func buildRunner(t *testing.T, cmd *testCmd, args ...string) *cmdRunner {
	return buildRunnerWithRoot(t, cmd, "", args...)
}

func buildRunnerWithRoot(t *testing.T, cmd *testCmd, root string, args ...string) *cmdRunner {
	app := kingpin.New("sysl", "System Modelling Language Toolkit")
	_ = cmd.Configure(app)
	_, err := app.Parse(append([]string{cmd.name}, args...))
	require.NoError(t, err)

	return &cmdRunner{
		modules:  cmd.modules,
		Root:     root,
		commands: map[string]cmdutils.Command{cmd.name: cmd},
	}
}
