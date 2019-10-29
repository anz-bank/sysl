package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/testutil"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func testMain2(t *testing.T, args []string, golden string) {
	logger, _ := test.NewNullLogger()
	_, fs := testutil.WriteToMemOverlayFs("/")
	output := "out.textpb"
	rc := main2(append([]string{"sysl", "pb", "-o", output}, args...), fs, logger, main3)
	assert.Zero(t, rc)

	actual, err := afero.ReadFile(fs, output)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(golden)
	assert.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func TestMain2TextPB(t *testing.T) {
	t.Parallel()

	testMain2(t, []string{"tests/args.sysl"}, "tests/args.sysl.golden.textpb")
}

func TestMain2JSON(t *testing.T) {
	t.Parallel()

	testMain2(t, []string{"--mode", "json", "tests/args.sysl"}, "tests/args.sysl.golden.json")
}

func testMain2Stdout(t *testing.T, args []string, golden string) {
	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(append([]string{"sysl", "pb", "-o", " - "}, args...), fs, logger, main3)
	assert.Zero(t, rc)

	_, err := ioutil.ReadFile(golden)
	require.NoError(t, err)

	_, err = os.Stat("-")
	assert.True(t, os.IsNotExist(err))

	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2TextPBStdout(t *testing.T) {
	t.Parallel()

	testMain2Stdout(t, []string{"tests/args.sysl"}, "tests/args.sysl.golden.textpb")
}

func TestMain2JSONStdout(t *testing.T) {
	t.Parallel()

	testMain2Stdout(t, []string{"--mode", "json", "tests/args.sysl"}, "tests/args.sysl.golden.json")
}

func TestMain2BadMode(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"pb",
			"-o", " - ",
			"--mode", "BAD",
			"tests/args.sysl",
		},
		fs, logger, main3,
	)
	assert.NotZero(t, rc)
	testutil.AssertFsHasExactly(t, memFs)

	_, err := os.Stat("-")
	assert.True(t, os.IsNotExist(err))
}

func TestMain2BadLog(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"pb",
			"-o", "-",
			"--log", "BAD",
			"tests/args.sysl",
		},
		fs, logger, main3,
	)
	assert.NotZero(t, rc)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2SeqdiagWithMissingFile(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sd",
			"-o", "%(epname).png",
			"tests/MISSING.sysl",
			"-a", "Project :: Sequences",
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, rc)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2SeqdiagWithNonsensicalOutput(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	out := "/out.zzz"
	_, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", out,
			"-b", "Server <- DB=call to database",
			"-b", "Server <- Login=call to database",
			"tests/call.sysl",
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
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", out,
			"-b", "Server <- DB=call to database",
			"-b", "Server <- Login=call to database",
			"tests/call.sysl",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, rc)
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "blackbox 'Server <- DB' passed on commandline not hit\n", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithReadOnlyFs(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	out := "/out2.png"
	_, fs := testutil.WriteToMemOverlayFs("/")
	fs = afero.NewReadOnlyFs(fs)
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", out,
			"-b", "Server <- DB=call to database",
			"-b", "Server <- Login=call to database",
			"tests/call.sysl",
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
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	ret := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", "tests/call2.png",
			"-b", "Server <- DB",
			"-b", "Server <- Login",
			"tests/call.sysl",
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, ret)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "expected KEY=VALUE got 'Server <- DB'", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithBlackboxSysl(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-o", "%(epname).png",
			"tests/blackbox.sysl",
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
	testutil.AssertFsHasExactly(t, memFs, "/SEQ-One.png", "/SEQ-Two.png")
}

func TestMain2WithBlackboxSyslEmptyEndpoints(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-o", "%(epname).png",
			"tests/blackbox.sysl",
			"-a", "Project :: Integrations",
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, rc)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "no call statements to build sequence diagram for endpoint PROJECT-E2E", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
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
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-g",
			"-o", "tests/call3.png",
			"tests/call.sysl",
		},
		fs, logger, main3,
	)
	assert.NotEqual(t, 0, rc)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "expected argument for flag '-g'", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithGroupingParamsCommandline(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	out := filepath.Clean("/out3.png")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-g", "owner",
			"-o", out,
			"tests/call.sysl",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, rc)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGroupingParamsSysl(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	rc := main2(
		[]string{
			"sysl",
			"sd",
			"-g", "location",
			"-o", "%(epname).png",
			"tests/groupby.sysl",
			"-a", "Project :: Sequences",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, rc)
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "Ignoring groupby passed from command line", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs, "/SEQ-One.png", "/SEQ-Two.png")
}

func TestMain2WithGenerateIntegrations(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	out := "/indirect_1.png"
	main2(
		[]string{
			"sysl",
			"ints",
			"--root", "./tests/",
			"-o", out,
			"-j", "Project",
			"indirect_1.sysl",
		},
		fs, logger, main3,
	)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGenerateCode(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	root := "."
	ret := main2(
		[]string{
			"sysl",
			"gen",
			"--root", root,
			"--root-transform", ".",
			"--transform", "tests/test.gen_multiple_annotations.sysl",
			"--grammar", "tests/test.gen.g",
			"--app-name", "Model",
			"--start", "javaFile",
			"tests/model.sysl",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	out, err := filepath.Abs(filepath.Join(root, "Model.java"))
	require.NoError(t, err)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGenerateCodeReadOnlyFs(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	_, fs := testutil.WriteToMemOverlayFs("/")
	fs = afero.NewReadOnlyFs(fs)
	ret := main2(
		[]string{
			"sysl",
			"gen",
			"--root", ".",
			"--root-transform", ".",
			"--transform", "tests/test.gen_multiple_annotations.sysl",
			"--grammar", "tests/test.gen.g",
			"--app-name", "Model",
			"--start", "javaFile",
			"tests/model.sysl",
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
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	out := "/out.textpb"
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "textpb",
			"-o", out,
			"tests/call.sysl",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithJSONMode(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	out := filepath.Clean("/out.json")
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "json",
			"-o", out,
			"tests/call.sysl",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithTextPbModeStdout(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "textpb",
			"-o", " - ",
			"tests/call.sysl",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithJSONModeStdout(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "json",
			"-o", " - ",
			"tests/call.sysl",
		},
		fs, logger, main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptySdParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "sd", "-g", " ", "-o", "", "tests/groupby.sysl", "-a", " "}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "'output' value passed is empty\n"+
		"'app' value passed is empty\n"+
		"'groupby' value passed is empty\n", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyPbParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "pb", "-o", " ", "--mode", "", "tests/call.sysl"}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t,
		"'output' value passed is empty\n'mode' value passed is empty\n", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyGenParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "gen", "--transform",
		"tests/test.gen_multiple_annotations.sysl", "--grammar", " ", "--start", "", "--outdir", " "}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t,
		"'grammar' value passed is empty\n"+
			"'start' value passed is empty\n"+
			"'outdir' value passed is empty\n", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyIntsParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "ints", "-o", "", "-j", " ", "indirect_1.sysl"}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t,
		"'output' value passed is empty\n"+
			"'project' value passed is empty\n", hook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithDataMultipleFiles(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png", "tests/data.sysl", "-j", "Project"}, fs, logger, main3)
	testutil.AssertFsHasExactly(t, memFs, "/Relational-Model.png", "/Object-Model.png")
}

func TestMain2WithDataSingleFile(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "data.png", "tests/data.sysl", "-j", "Project"}, fs, logger, main3)
	testutil.AssertFsHasExactly(t, memFs, "/data.png")
}

func TestMain2WithDataNoProject(t *testing.T) {
	t.Parallel()
	logger, testHook := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png", "tests/data.sysl"}, fs, logger, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "project not found in sysl", testHook.LastEntry().Message)
	testHook.Reset()
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithDataFilter(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png", "-f", "Object-Model.png", "tests/data.sysl", "-j",
		"Project"}, fs, logger, main3)
	testutil.AssertFsHasExactly(t, memFs, "/Object-Model.png")
}

func TestMain2WithDataMultipleRelationships(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png", "tests/datareferences.sysl", "-j", "Project"},
		fs, logger, main3)
	testutil.AssertFsHasExactly(t, memFs, "/Relational-Model.png", "/Object-Model.png")
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
	memFs, fs := testutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "export", "-o", "SIMPLE_SWAGGER_EXAMPLE.yaml", "-a", "testapp",
		"exporter/test-data/SIMPLE_SWAGGER_EXAMPLE.sysl"}, fs, logger, main3)
	testutil.AssertFsHasExactly(t, memFs, "/SIMPLE_SWAGGER_EXAMPLE.yaml")
}

func TestSwaggerExportTargetDir(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	tmp1, err := ioutil.TempDir("", "tmp1")
	assert.Nil(t, err)
	main2([]string{"sysl", "export", "-o", tmp1 + "/SIMPLE_SWAGGER_EXAMPLE1.yaml", "-a", "testapp",
		"exporter/test-data/SIMPLE_SWAGGER_EXAMPLE.sysl"}, afero.NewOsFs(), logger, main3)
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
		"-a", "testapp", "exporter/test-data/SIMPLE_SWAGGER_EXAMPLE.sysl"}, afero.NewOsFs(), logger, main3)
	_, err = ioutil.ReadFile(tmp2 + "/SIMPLE_SWAGGER_EXAMPLE2.json")
	assert.Nil(t, err)
	os.RemoveAll(tmp2)
}

func TestSwaggerExportInvalid(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	_, fs := testutil.WriteToMemOverlayFs("/")
	errInt := main2([]string{"sysl", "export", "-o", "SIMPLE_SWAGGER_EXAMPLE1.blah", "-a", "testapp",
		"exporter/test-data/SIMPLE_SWAGGER_EXAMPLE.sysl"}, fs, logger, main3)
	assert.True(t, errInt == 1)
}

func TestSwaggerAppExportNoDir(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	main2([]string{"sysl", "export", "-o", "out/%(appname).yaml",
		"exporter/test-data/multiple/SIMPLE_SWAGGER_EXAMPLE_MULTIPLE.sysl"}, afero.NewOsFs(), logger, main3)
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
		"exporter/test-data/multiple/SIMPLE_SWAGGER_EXAMPLE_MULTIPLE.sysl"}, afero.NewOsFs(), logger, main3)
	for _, file := range []string{tmp3 + "/single.yaml", tmp3 + "/multiple.yaml"} {
		_, err := ioutil.ReadFile(file)
		assert.Nil(t, err)
	}
	os.RemoveAll(tmp3)
}
