package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/pbutil"
	"github.com/anz-bank/sysl/pkg/sysl"
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
	t.Parallel()

	syslKingpin := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := syslKingpin.Command("foo", "")
	_ = cmd.Flag("bar", "").Default("foo").String()
	_ = cmd.Flag("other", "").Default("foo").String()

	EnsureFlagsNonEmpty(cmd, "bar")

	args := []string{"foo", "--bar", ""}
	selected, err := syslKingpin.Parse(args)
	assert.Equal(t, "foo", selected)
	assert.NoError(t, err)
}

func TestLoadModule_stdinBytes(t *testing.T) {
	t.Parallel()

	var mods []*sysl.Module
	r := buildRunner(t, modsCmd(&mods))

	m, err := parse.NewParser().ParseString(emptyAppSrc)
	require.NoError(t, err)
	stdin := new(bytes.Buffer)
	require.NoError(t, pbutil.GeneratePBBinaryMessage(stdin, m))

	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin, io.Discard))

	assert.Len(t, mods, 1)
	assert.NotEmpty(t, mods[0].Apps["Stdin :: App"])
}

func TestLoadModule_stdinSplitBytes(t *testing.T) {
	t.Parallel()

	var mods []*sysl.Module
	opt := pbutil.OutputOptions{Compact: false}
	r := buildRunner(t, modsCmd(&mods))

	m, err := parse.NewParser().ParseString(emptyAppSrc)
	require.NoError(t, err)
	stdin := new(bytes.Buffer)
	dummyFs := afero.NewMemMapFs()
	require.NoError(t, pbutil.OutputSplitApplications(m, "binarypb", opt, "/", "test.pb", dummyFs))

	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin, io.Discard))

	fileExists, err := afero.Exists(dummyFs, path.Clean("/Stdin/App/test.pb"))
	assert.NoError(t, err)
	assert.True(t, fileExists)
}

func TestLoadModule_stdinSplitJSON(t *testing.T) {
	t.Parallel()

	var mods []*sysl.Module
	opt := pbutil.OutputOptions{Compact: false}
	r := buildRunner(t, modsCmd(&mods))

	m, err := parse.NewParser().ParseString(emptyAppSrc)
	require.NoError(t, err)
	stdin := new(bytes.Buffer)
	dummyFs := afero.NewMemMapFs()
	require.NoError(t, pbutil.OutputSplitApplications(m, "json", opt, "/", "test.json", dummyFs))

	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin, io.Discard))

	fileExists, err := afero.Exists(dummyFs, path.Clean("/Stdin/App/test.json"))
	assert.NoError(t, err)
	assert.True(t, fileExists)
}

func TestLoadModule_stdinSplitTextPB(t *testing.T) {
	t.Parallel()

	var mods []*sysl.Module
	opt := pbutil.OutputOptions{Compact: false}
	r := buildRunner(t, modsCmd(&mods))

	m, err := parse.NewParser().ParseString(emptyAppSrc)
	require.NoError(t, err)
	stdin := new(bytes.Buffer)
	dummyFs := afero.NewMemMapFs()
	require.NoError(t, pbutil.OutputSplitApplications(m, "textpb", opt, "/", "test.textpb", dummyFs))

	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin, io.Discard))

	fileExists, err := afero.Exists(dummyFs, path.Clean("/Stdin/App/test.textpb"))
	assert.NoError(t, err)
	assert.True(t, fileExists)
}

func TestLoadModule_stdin(t *testing.T) {
	t.Parallel()

	var mods []*sysl.Module
	r := buildRunner(t, modsCmd(&mods))

	stdin := toStdin(t, stdinFile{Path: filepath.Join(cwd(), "test.sysl"), Content: emptyAppSrc})
	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin, io.Discard))

	assert.Len(t, mods, 1)
	assert.NotEmpty(t, mods[0].Apps["Stdin :: App"])
}

func TestLoadModule_stdinImport(t *testing.T) {
	t.Parallel()

	var mods []*sysl.Module
	r := buildRunner(t, modsCmd(&mods))

	stdin := toStdin(t, stdinFile{Path: "test.sysl", Content: importSrc})
	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin, io.Discard))

	assert.Len(t, mods, 1)
	assert.Len(t, mods[0].Apps, 3)
	for _, a := range []string{"Stdin :: App", "Namespace1 :: App1", "Namespace1 :: App2"} {
		assert.Contains(t, mods[0].Apps, a)
	}
}

func TestLoadModule_stdinPathImport(t *testing.T) {
	t.Parallel()

	var mods []*sysl.Module
	r := buildRunner(t, modsCmd(&mods))

	stdin := toStdin(t, stdinFile{Path: filepath.Join("..", "..", "tests", "test.sysl"), Content: "import simple.sysl"})
	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin, io.Discard))

	assert.Len(t, mods, 1)
	assert.Len(t, mods[0].Apps, 2)
	for _, a := range []string{"Namespace1 :: App1", "Namespace1 :: App2"} {
		assert.Contains(t, mods[0].Apps, a)
		for _, src := range mods[0].Apps[a].SourceContexts {
			assert.True(t, strings.HasSuffix(src.File, "tests/simple.sysl"))
		}
	}
}

func TestLoadModule_pathsAndStdin(t *testing.T) {
	t.Parallel()

	args := []string{"simple.sysl"}
	cmd := &testCmd{name: cmdName, maxSyslModule: 1, exec: func(args cmdutils.ExecuteArgs) {}}
	r := buildRunnerWithRoot(t, cmd, testDir, args...)

	stdin := toStdin(t, stdinFile{Path: "test.sysl", Content: emptyAppSrc})
	require.NoError(t, r.Run(cmdName, afero.NewOsFs(), logrus.StandardLogger(), stdin, io.Discard))
}

func TestEnsureFlagsNonEmpty_AllowsExcludes(t *testing.T) {
	t.Parallel()

	syslKingpin := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := syslKingpin.Command("foo", "")
	_ = cmd.Flag("bar", "").Default("foo").String()
	_ = cmd.Flag("other", "").Default("foo").String()

	EnsureFlagsNonEmpty(cmd, "bar")

	args := []string{"foo", "--bar", ""}
	selected, err := syslKingpin.Parse(args)
	assert.Equal(t, "foo", selected)
	assert.NoError(t, err)
}

func TestEnsureFlagsNonEmpty(t *testing.T) {
	t.Parallel()

	syslKingpin := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := syslKingpin.Command("foo", "")
	cmd.Flag("bar", "").Default("foo")

	EnsureFlagsNonEmpty(cmd)

	args := []string{"foo", "--bar", ""}
	_, err := syslKingpin.ParseContext(args)
	assert.Error(t, err)
}

type testCmd struct {
	name          string
	maxSyslModule int
	exec          func(cmdutils.ExecuteArgs)
	modules       []string
}

// modsCmd returns a command that will set the mods arg to the array of loaded modules.
func modsCmd(mods *[]*sysl.Module) *testCmd {
	return &testCmd{name: cmdName, maxSyslModule: 1, exec: func(args cmdutils.ExecuteArgs) {
		*mods = args.Modules
	}}
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
