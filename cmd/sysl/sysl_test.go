package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const currentWorkingDirectory = "."

type folderTestStructure struct {
	name,
	module,
	root,
	expectedRoot,
	rootMarkerPath string
	rootFound bool
	structure folderStructure
}

type folderStructure struct {
	folders, files []string
}

func assertLogEntry(
	t *testing.T,
	entry *logrus.Entry,
	level logrus.Level,
	format string,
	args ...interface{},
) bool {
	okLevel := assert.Equal(t, level, entry.Level)
	okMessage := assert.Equal(t, fmt.Sprintf(format, args...), entry.Message)
	return okLevel && okMessage
}

func runMain2(t *testing.T, fs afero.Fs, args []string, golden string) {
	logger, _ := test.NewNullLogger()
	output := "out.textpb"
	rc := main2(append([]string{"sysl", "pb", "-o", output}, args...), fs, logger, main3)
	assert.Zero(t, rc)

	actual, err := afero.ReadFile(fs, output)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(golden)
	assert.NoError(t, err)

	expected = syslutil.HandleCRLF(expected)
	assert.Equal(t, string(expected), string(actual))
}

func testMain2WithSyslRootMarker(t *testing.T, args []string, golden string) {
	_, fs := syslutil.WriteToMemOverlayFs("/")
	dir := syslutil.MustAbsolute(t, fmt.Sprintf(testDir+"%s", syslRootMarker))
	require.NoError(t, fs.MkdirAll(dir, os.ModeDir))
	runMain2(t, fs, args, golden)
}

func testMain2(t *testing.T, args []string, golden string) {
	_, fs := syslutil.WriteToMemOverlayFs("/")
	runMain2(t, fs, args, golden)
}

func testAllMain2(t *testing.T, args []string, inputFile string, golden string) {
	// no root defined
	noRootFile := filepath.Join(testDir, inputFile)
	testMain2(t, append(args, noRootFile), filepath.Join(testDir, golden))

	// root marker
	testMain2WithSyslRootMarker(t, append(args, noRootFile), filepath.Join(testDir, golden))

	// root flag defined
	rootFile := "tests/" + inputFile

	args = append([]string{"--root", projDir}, args...)
	golden = filepath.Join(testDir, golden)
	out := "rooted" + filepath.Base(golden)
	golden = filepath.Join(filepath.Dir(golden), out)
	testMain2(t, append(args, rootFile), golden)

	testMain2WithSyslRootMarker(t, append(args, rootFile), golden)
}

func TestMain2TextPB(t *testing.T) {
	t.Parallel()
	// tests/ is used because the projDir in args is expecting to have a directory called tests
	testAllMain2(t, []string{}, "args.sysl", "args.sysl.golden.textpb")
}

func TestMain2JSON(t *testing.T) {
	t.Parallel()

	testAllMain2(t, []string{"--mode", "json"}, "args.sysl", "args.sysl.golden.json")
}

func testMain2Stdout(t *testing.T, args []string, golden string) {
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(append([]string{"sysl", "pb", "-o", " - "}, args...), fs, logger, main3)
	assert.Zero(t, rc)

	_, err := ioutil.ReadFile(golden)
	require.NoError(t, err)

	_, err = os.Stat("-")
	assert.True(t, os.IsNotExist(err))

	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2TextPBStdout(t *testing.T) {
	t.Parallel()

	testMain2Stdout(t, []string{filepath.Join(testDir, "args.sysl")}, filepath.Join(testDir, "args.sysl.golden.textpb"))
}

func TestMain2JSONStdout(t *testing.T) {
	t.Parallel()

	testMain2Stdout(t, []string{"--mode", "json",
		filepath.Join(testDir, "args.sysl")}, filepath.Join(testDir, "args.sysl.golden.json"))
}

func TestMain2BadMode(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"pb",
			"-o", " - ",
			"--mode", "BAD",
			filepath.Join(testDir, "args.sysl"),
		},
		fs, logger, main3,
	)
	assert.NotZero(t, rc)
	syslutil.AssertFsHasExactly(t, memFs)

	_, err := os.Stat("-")
	assert.True(t, os.IsNotExist(err))
}

func TestMain2BadLog(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"pb",
			"-o", "-",
			"--log", "BAD",
			filepath.Join(testDir, "args.sysl"),
		},
		fs, logger, main3,
	)
	assert.NotZero(t, rc)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2SeqdiagWithMissingFile(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sd",
			"-o", "%(epname).png",
			filepath.Join(testDir, "MISSING.sysl"),
			"-a", "Project :: Sequences",
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, rc)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2SeqdiagWithNonsensicalOutput(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	out := "/out.zzz"
	_, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", out,
			"-b", "Server <- DB=call to database",
			"-b", "Server <- Login=call to database",
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, rc)
	assertLogEntry(t, hook.LastEntry(), logrus.ErrorLevel,
		`extension must be svg, png or uml, not "zzz"`)
}

func TestMain2WithBlackboxParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	out := filepath.Clean("/out1.png")
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", out,
			"-b", "Server <- DB=call to database",
			"-b", "Server <- Login=call to database",
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, rc)
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "blackbox 'Server <- DB' passed on commandline not hit\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithReadOnlyFs(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	out := "/out2.png"
	_, fs := syslutil.WriteToMemOverlayFs("/")
	fs = afero.NewReadOnlyFs(fs)
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", out,
			"-b", "Server <- DB=call to database",
			"-b", "Server <- Login=call to database",
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, rc)
	assertLogEntry(t, hook.LastEntry(), logrus.ErrorLevel,
		"writing %q: operation not permitted", out)
}

func TestMain2WithBlackboxParamsFaultyArguments(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	ret := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", "call2.png",
			"-b", "Server <- DB",
			"-b", "Server <- Login",
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, ret)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "expected KEY=VALUE got 'Server <- DB'", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithBlackboxSysl(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-o", "%(epname).png",
			filepath.Join(testDir, "blackbox.sysl"),
			"-a", "Project :: Sequences",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, rc)
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "blackbox 'SomeApp <- AppEndpoint' not hit in app 'Project :: Sequences'\n",
		hook.Entries[len(hook.Entries)-1].Message)
	assert.Equal(t, "blackbox 'SomeApp <- BarEndpoint1' not hit in app 'Project :: Sequences :: SEQ-Two'\n",
		hook.Entries[len(hook.Entries)-2].Message)
	assert.Equal(t, "blackbox 'SomeApp <- BarEndpoint' not hit in app 'Project :: Sequences :: SEQ-One'\n",
		hook.Entries[len(hook.Entries)-3].Message)
	syslutil.AssertFsHasExactly(t, memFs, "/SEQ-One.png", "/SEQ-Two.png")
}

func TestMain2WithBlackboxSyslEmptyEndpoints(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-o", "%(epname).png",
			filepath.Join(testDir, "blackbox.sysl"),
			"-a", "Project :: Integrations",
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, rc)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "no call statements to build sequence diagram for endpoint PROJECT-E2E", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2Fatal(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	assert.Equal(t, 42, main2(nil, nil, logger, func(_ []string, _ afero.Fs, _ *logrus.Logger) error {
		return parse.Exitf(42, "Exit error")
	}))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
}

func TestMain2WithGroupingParamsGroupParamAbsent(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-g",
			"-o", "tests/call3.png",
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, rc)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "expected argument for flag '-g'", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithGroupingParamsCommandline(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	out := filepath.Clean("/out3.png")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-g", "owner",
			"-o", out,
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, rc)
	syslutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGroupingParamsSysl(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-g", "location",
			"-o", "%(epname).png",
			filepath.Join(testDir, "groupby.sysl"),
			"-a", "Project :: Sequences",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, rc)
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "Ignoring groupby passed from command line", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs, "/SEQ-One.png", "/SEQ-Two.png")
}

func TestMain2WithGenerateIntegrations(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	out := "/indirect_1.png"
	main2(
		[]string{
			"sysl",
			"ints",
			"--root", testDir,
			"-o", out,
			"-j", "Project",
			"indirect_1.sysl",
		},
		fs, logger, main3,
	)
	syslutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGenerateCode(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	ret := main2(
		[]string{
			"sysl",
			"gen",
			"--root", testDir,
			"--root-transform", testDir,
			"--transform", "test.gen_multiple_annotations.sysl",
			"--grammar", filepath.Join(testDir, "test.gen.g"),
			"--app-name", "Model",
			"--start", "javaFile",
			"--dep-path", "github.service.anz/anzx/bff/gen/",
			"model.sysl",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	out, err := filepath.Abs(filepath.Join(".", "Model.java"))
	require.NoError(t, err)
	syslutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGenerateCodeReadOnlyFs(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")
	fs = afero.NewReadOnlyFs(fs)
	ret := main2(
		[]string{
			"sysl",
			"gen",
			"--root", testDir,
			"--root-transform", testDir,
			"--transform", "test.gen_multiple_annotations.sysl",
			"--grammar", filepath.Join(testDir, "test.gen.g"),
			"--app-name", "Model",
			"--start", "javaFile",
			"--dep-path", "github.service.anz/anzx/bff/gen/",
			"model.sysl",
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, ret)
	assertLogEntry(t, hook.LastEntry(), logrus.ErrorLevel,
		`unable to create "Model.java": operation not permitted`)
}

func TestMain2WithTextPbMode(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	out := "/out.textpb"
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "textpb",
			"-o", out,
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	syslutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithJSONMode(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	out := filepath.Clean("/out.json")
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "json",
			"-o", out,
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	syslutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithTextPbModeStdout(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "textpb",
			"-o", " - ",
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithJSONModeStdout(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "json",
			"-o", " - ",
			filepath.Join(testDir, "call.sysl"),
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptySdParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "sd", "-g", " ", "-o", "", "tests/groupby.sysl", "-a", " "}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "'output' value passed is empty\n"+
		"'app' value passed is empty\n"+
		"'groupby' value passed is empty\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyPbParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "pb", "-o", " ", "--mode", "", filepath.Join(testDir, "call.sysl")}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t,
		"'output' value passed is empty\n'mode' value passed is empty\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyGenParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "gen", "--transform",
		filepath.Join(testDir, "test.gen_multiple_annotations.sysl"),
		"--grammar", " ", "--start", "", "--outdir", " "}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t,
		"'grammar' value passed is empty\n"+
			"'start' value passed is empty\n"+
			"'outdir' value passed is empty\n"+
			"'dep-path' value passed is empty\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyIntsParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "ints", "-o", "", "-j", " ", "indirect_1.sysl"}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t,
		"'output' value passed is empty\n"+
			"'project' value passed is empty\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithDataMultipleFiles(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png",
		filepath.Join(testDir, "data.sysl"), "-j", "Project"}, fs, logger, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/Relational-Model.png", "/Object-Model.png")
}

func TestMain2WithDataSingleFile(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "data.png",
		filepath.Join(testDir, "data.sysl"), "-j", "Project"}, fs, logger, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/data.png")
}

func TestMain2WithDataNoProject(t *testing.T) {
	t.Parallel()
	logger, testHook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png",
		filepath.Join(testDir, "data.sysl")}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "project not found in sysl", testHook.LastEntry().Message)
	testHook.Reset()
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithDataFilter(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png", "-f",
		"Object-Model.png", filepath.Join(testDir, "data.sysl"), "-j",
		"Project"}, fs, logger, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/Object-Model.png")
}

func TestMain2WithDataMultipleRelationships(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png",
		filepath.Join(testDir, "datareferences.sysl"), "-j", "Project"},
		fs, logger, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/Relational-Model.png", "/Object-Model.png")
}

func TestMain2WithBinaryInfoCmd(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	exitCode := main2([]string{"sysl", "info"}, nil, logger, main3)
	assert.Equal(t, 0, exitCode)
}

func TestSwaggerExportCurrentDir(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "export", "-o", "SIMPLE_SWAGGER_EXAMPLE.yaml", "-a", "testapp",
		syslDir + "exporter/test-data/SIMPLE_SWAGGER_EXAMPLE.sysl"}, fs, logger, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/SIMPLE_SWAGGER_EXAMPLE.yaml")
}

func TestSwaggerExportTargetDir(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	tmp1, err := ioutil.TempDir("", "tmp1")
	assert.Nil(t, err)
	main2([]string{"sysl", "export", "-o", tmp1 + "/SIMPLE_SWAGGER_EXAMPLE1.yaml", "-a", "testapp",
		syslDir + "exporter/test-data/SIMPLE_SWAGGER_EXAMPLE.sysl"}, afero.NewOsFs(), logger, main3)
	_, err = ioutil.ReadFile(tmp1 + "/SIMPLE_SWAGGER_EXAMPLE1.yaml")
	assert.Nil(t, err)
	os.RemoveAll(tmp1)
}

func TestSwaggerExportJson(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	tmp2, err := ioutil.TempDir("", "tmp2")
	assert.Nil(t, err)
	main2([]string{"sysl", "export", "-o", tmp2 + "/SIMPLE_SWAGGER_EXAMPLE2.json",
		"-a", "testapp", syslDir + "exporter/test-data/SIMPLE_SWAGGER_EXAMPLE.sysl"}, afero.NewOsFs(), logger, main3)
	_, err = ioutil.ReadFile(tmp2 + "/SIMPLE_SWAGGER_EXAMPLE2.json")
	assert.Nil(t, err)
	os.RemoveAll(tmp2)
}

func TestSwaggerExportInvalid(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")
	errInt := main2([]string{"sysl", "export", "-o", "SIMPLE_SWAGGER_EXAMPLE1.blah", "-a", "testapp",
		syslDir + "exporter/test-data/SIMPLE_SWAGGER_EXAMPLE.sysl"}, fs, logger, main3)
	assert.True(t, errInt == 1)
}

func TestSwaggerAppExportNoDir(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	main2([]string{"sysl", "export", "-o", "out/%(appname).yaml",
		syslDir + "exporter/test-data/multiple/SIMPLE_SWAGGER_EXAMPLE_MULTIPLE.sysl"}, afero.NewOsFs(), logger, main3)
	for _, file := range []string{"out/single.yaml", "out/multiple.yaml"} {
		_, err := ioutil.ReadFile(file)
		assert.Nil(t, err)
	}
	os.RemoveAll("out")
}

func TestSwaggerAppExportDirExists(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	tmp3, err := ioutil.TempDir("", "tmp3")
	assert.Nil(t, err)
	main2([]string{"sysl", "export", "-o", tmp3 + "/%(appname).yaml",
		syslDir + "exporter/test-data/multiple/SIMPLE_SWAGGER_EXAMPLE_MULTIPLE.sysl"}, afero.NewOsFs(), logger, main3)
	for _, file := range []string{tmp3 + "/single.yaml", tmp3 + "/multiple.yaml"} {
		_, err := ioutil.ReadFile(file)
		assert.Nil(t, err)
	}
	os.RemoveAll(tmp3)
}

func TestHandleProjectRoot(t *testing.T) {
	successfulTest := folderStructure{
		folders: []string{
			"./SuccessfulTest/path/to/module",
			fmt.Sprintf("./SuccessfulTest/%s", syslRootMarker),
			"./SuccessfulTest/path/to/another/module",
			fmt.Sprintf("./SuccessfulTest/path/to/another/%s", syslRootMarker),
		},
		files: []string{
			"./SuccessfulTest/path/to/module/test.sysl",
			"./SuccessfulTest/test2.sysl",
			"./SuccessfulTest/path/to/another/module/test3.sysl",
		},
	}

	definedRootNoMarker := folderStructure{
		folders: []string{
			"./DefinedRootAndSyslRootUndefinedTest/path/to/module/",
		},
		files: []string{
			"./DefinedRootAndSyslRootUndefinedTest/path/to/module/test.sysl",
		},
	}

	definedRootFlagAndMarkerFound := folderStructure{
		folders: []string{
			"./DefinedRootAndSyslRootDefinedTest/path/to/module/",
			fmt.Sprintf("./DefinedRootAndSyslRootDefinedTest/path/%s", syslRootMarker),
		},
		files: []string{
			"./DefinedRootAndSyslRootDefinedTest/path/to/module/test.sysl",
		},
	}

	undefinedRoot := folderStructure{
		folders: []string{
			"./UndefinedRootAndUndefinedSyslRoot/",
		},
		files: []string{
			"./UndefinedRootAndUndefinedSyslRoot/test.sysl",
		},
	}
	systemRoot := syslutil.MustAbsolute(t, string(os.PathSeparator))
	tests := []folderTestStructure{
		{
			name:         "Successful test: finding a root marker",
			root:         "",
			module:       successfulTest.files[0],
			structure:    successfulTest,
			expectedRoot: syslutil.MustAbsolute(t, "SuccessfulTest"),
			rootFound:    true,
		},
		{
			name:         "Successful test: finding a root marker in the same directory as the module",
			root:         "",
			module:       successfulTest.files[1],
			structure:    successfulTest,
			expectedRoot: syslutil.MustAbsolute(t, "SuccessfulTest"),
			rootFound:    true,
		},
		{
			name:         "Successful test: finding the closest root marker",
			root:         "",
			module:       successfulTest.files[2],
			structure:    successfulTest,
			expectedRoot: syslutil.MustAbsolute(t, "SuccessfulTest/path/to/another"),
			rootFound:    true,
		},
		{
			name: "Root flag is defined and root marker does not exist",
			root: "DefinedRootAndSyslRootUndefinedTest/path/",
			module: syslutil.MustRelative(t, "DefinedRootAndSyslRootUndefinedTest/path/",
				definedRootNoMarker.files[0]),
			structure:    definedRootNoMarker,
			expectedRoot: "DefinedRootAndSyslRootUndefinedTest/path/",
			rootFound:    true,
		},
		{
			name:         "Defined relative root",
			root:         currentWorkingDirectory,
			module:       filepath.Clean(definedRootNoMarker.files[0]),
			structure:    definedRootNoMarker,
			expectedRoot: currentWorkingDirectory,
			rootFound:    true,
		},
		{
			root:         systemRoot,
			name:         "Defined absolute path root",
			module:       syslutil.MustAbsolute(t, definedRootNoMarker.files[0]),
			structure:    definedRootNoMarker,
			expectedRoot: systemRoot,
			rootFound:    true,
		},
		{
			name:         "Defined relative root with absolute module path rooted at root",
			root:         currentWorkingDirectory,
			module:       filepath.Join(systemRoot, filepath.Clean(definedRootNoMarker.files[0])),
			structure:    definedRootNoMarker,
			expectedRoot: currentWorkingDirectory,
			rootFound:    true,
		},
		{
			name:           "Defined root flag and root",
			root:           currentWorkingDirectory,
			module:         syslutil.MustRelative(t, currentWorkingDirectory, definedRootFlagAndMarkerFound.files[0]),
			structure:      definedRootFlagAndMarkerFound,
			expectedRoot:   currentWorkingDirectory,
			rootMarkerPath: syslutil.MustAbsolute(t, "./DefinedRootAndSyslRootDefinedTest/path/"),
			rootFound:      true,
		},
		{
			name:           "Defined root flag and root marker with absolute path module rooted at root",
			root:           "./DefinedRootAndSyslRootDefinedTest/",
			module:         "/path/to/module/test.sysl",
			structure:      definedRootFlagAndMarkerFound,
			expectedRoot:   "./DefinedRootAndSyslRootDefinedTest/",
			rootMarkerPath: syslutil.MustAbsolute(t, "./DefinedRootAndSyslRootDefinedTest/path/"),
			rootFound:      true,
		},
		{
			name:         "Root is not defined",
			root:         "",
			module:       undefinedRoot.files[0],
			structure:    undefinedRoot,
			expectedRoot: filepath.Dir(undefinedRoot.files[0]),
			rootFound:    false,
		},
	}

	for _, ts := range tests {
		ts := ts
		t.Run(ts.name, func(t *testing.T) {
			t.Parallel()

			logger, _ := test.NewNullLogger()
			fs := afero.NewMemMapFs()
			syslutil.BuildFolderTest(t, fs, ts.structure.folders, ts.structure.files)

			config := newProjectConfiguration()
			err := config.configureProject(ts.root, ts.module, fs, logger)

			require.Equal(t, ts.rootFound, config.rootIsFound)
			require.NoError(t, err)
			require.Equal(t, ts.expectedRoot, config.root)
			require.Equal(t, ts.getExpectedModule(t), config.module)
		})
	}
}

func (ts folderTestStructure) getExpectedModule(t *testing.T) string {
	// if root is defined, expected root and root param is the same and module is not changed
	if ts.expectedRoot == ts.root {
		return ts.module
	}
	return syslutil.MustRelative(t, ts.expectedRoot, ts.module)
}

func TestCodegenGrammarImportDefOut(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "import", "-i",
		syslDir + "importer/tests-grammar/simplerules.gen.g", "-a", "go"}, fs, logger, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/output.sysl")
}

func TestCodegenGrammarImport(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "import", "-i", syslDir + "importer/tests-grammar/unions.gen.g",
		"-o", "out.sysl", "-a", "go"}, fs, logger, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/out.sysl")
}
