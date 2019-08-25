package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/testutil"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMain2(t *testing.T, args []string, golden string) {
	_, fs := testutil.WriteToMemOverlayFs(".")
	output := "out.textpb"
	rc := main2(append([]string{"sysl", "pb", "-o", output}, args...), fs, main3)
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
	memFs, fs := testutil.WriteToMemOverlayFs(".")
	rc := main2(append([]string{"sysl", "pb", "-o", " - "}, args...), fs, main3)
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

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	rc := main2([]string{"sysl", "pb", "-o", " - ", "--mode", "BAD", "tests/args.sysl"}, fs, main3)
	assert.NotZero(t, rc)
	testutil.AssertFsHasExactly(t, memFs)

	_, err := os.Stat("-")
	assert.True(t, os.IsNotExist(err))
}

func TestMain2BadLog(t *testing.T) {
	t.Parallel()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	rc := main2([]string{"sysl", "pb", "-o", "-", "--log", "BAD", "tests/args.sysl"}, fs, main3)
	assert.NotZero(t, rc)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2SeqdiagWithMissingFile(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	rc := main2([]string{"sd", "-o", "%(epname).png", "tests/MISSING.sysl", "-a", "Project :: Sequences"}, fs, main3)
	assert.NotEqual(t, 0, rc)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2SeqdiagWithImpossibleOutput(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	out := "tests/call.zzz"
	rc := main2(
		[]string{
			"sd",
			"-s", "MobileApp <- Login",
			"-o", out,
			"-b", "Server <- DB=call to database",
			"-b", "Server <- Login=call to database",
			"tests/call.sysl",
		},
		fs,
		main3,
	)
	assert.NotEqual(t, 0, rc)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithBlackboxParams(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	out := "/out.png"
	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-o", out,
			"-b", "Server <- DB=call to database",
			"-b", "Server <- Login=call to database",
			"tests/call.sysl",
		},
		fs,
		main3,
	)
	assert.Equal(t, logrus.WarnLevel, testHook.LastEntry().Level)
	assert.Equal(t, "blackbox 'Server <- DB' passed on commandline not hit\n", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithBlackboxParamsFaultyArguments(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
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
		fs,
		main3,
	)
	assert.NotEqual(t, 0, ret)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "expected KEY=VALUE got 'Server <- DB'", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithBlackboxSysl(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2(
		[]string{
			"sysl",
			"sd",
			"-o", "%(epname).png",
			"tests/blackbox.sysl",
			"-a", "Project :: Sequences",
		},
		fs,
		main3,
	)
	assert.Equal(t, logrus.WarnLevel, testHook.LastEntry().Level)
	assert.Equal(t, "blackbox 'SomeApp <- AppEndpoint' not hit in app 'Project :: Sequences'\n",
		testHook.Entries[len(testHook.Entries)-1].Message)
	assert.Equal(t, "blackbox 'SomeApp <- BarEndpoint1' not hit in app 'Project :: Sequences :: SEQ-Two'\n",
		testHook.Entries[len(testHook.Entries)-2].Message)
	assert.Equal(t, "blackbox 'SomeApp <- BarEndpoint' not hit in app 'Project :: Sequences :: SEQ-One'\n",
		testHook.Entries[len(testHook.Entries)-3].Message)
	testutil.AssertFsHasExactly(t, memFs, "/SEQ-One.png", "/SEQ-Two.png")
}

func TestMain2WithBlackboxSyslEmptyEndpoints(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2(
		[]string{
			"sysl",
			"sd",
			"-o", "%(epname).png",
			"tests/blackbox.sysl",
			"-a", "Project :: Integrations",
		},
		fs,
		main3,
	)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "No call statements to build sequence diagram for endpoint PROJECT-E2E", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2Fatal(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	assert.Equal(t, 42, main2(nil, nil, func(_ []string, _ afero.Fs) error {
		return parse.Exitf(42, "Exit error")
	}))
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
}

func TestMain2WithGroupingParamsGroupParamAbsent(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-g",
			"-o", "tests/call3.png",
			"tests/call.sysl",
		},
		fs,
		main3,
	)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "expected argument for flag '-g'", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithGroupingParamsCommandline(t *testing.T) {
	t.Parallel()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	out := "/out.png"
	main2(
		[]string{
			"sysl",
			"sd",
			"-s", "MobileApp <- Login",
			"-g", "owner",
			"-o", out,
			"tests/call.sysl",
		},
		fs,
		main3,
	)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGroupingParamsSysl(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2(
		[]string{
			"sysl",
			"sd",
			"-g", "location",
			"-o", "%(epname).png",
			"tests/groupby.sysl",
			"-a", "Project :: Sequences",
		},
		fs,
		main3,
	)
	assert.Equal(t, logrus.WarnLevel, testHook.LastEntry().Level)
	assert.Equal(t, "Ignoring groupby passed from command line", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs, "/SEQ-One.png", "/SEQ-Two.png")
}

func TestMain2WithGenerateIntegrations(t *testing.T) {
	t.Parallel()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
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
		fs,
		main3,
	)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGenerateCode(t *testing.T) {
	t.Parallel()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	ret := main2(
		[]string{
			"sysl",
			"gen",
			"--root-model", ".",
			"--model", "tests/model.sysl",
			"--root-transform", ".",
			"--transform", "tests/test.gen_multiple_annotations.sysl",
			"--grammar", "tests/test.gen.g",
			"--start", "javaFile",
		},
		fs,
		main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs, "/Model.java")
}

func TestMain2WithTextPbMode(t *testing.T) {
	t.Parallel()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	out := "/out.textpb"
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "textpb",
			"-o", out,
			"tests/call.sysl",
		},
		fs,
		main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithJSONMode(t *testing.T) {
	t.Parallel()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	out := "/out.json"
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "json",
			"-o", out,
			"tests/call.sysl",
		},
		fs,
		main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithTextPbModeStdout(t *testing.T) {
	t.Parallel()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode", "textpb",
			"-o", " - ",
			"tests/call.sysl",
			"-v",
		},
		fs,
		main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithJSONModeStdout(t *testing.T) {
	t.Parallel()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	ret := main2(
		[]string{
			"sysl",
			"pb",
			"--mode",
			"json",
			"-o",
			" - ",
			"tests/call.sysl",
			"-v",
		},
		fs,
		main3,
	)
	assert.Equal(t, 0, ret)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptySdParams(t *testing.T) {
	testHook := test.NewGlobal()
	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2([]string{"sysl", "sd", "-g", " ", "-o", "", "tests/groupby.sysl", "-a", " "}, fs, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "'output' value passed is empty\n"+
		"'groupby' value passed is empty\n"+
		"'app' value passed is empty\n", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyPbParams(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2([]string{"sysl", "pb", "--root", "", "-o", " ", "--mode", "", "tests/call.sysl"}, fs, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "'root' value passed is empty\n"+
		"'output' value passed is empty\n"+
		"'mode' value passed is empty\n", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyGenParams(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2([]string{"sysl", "gen", "--root-model", " ", "--root-transform", "", "--model", " ", "--transform",
		"tests/test.gen_multiple_annotations.sysl", "--grammar", " ", "--start", "", "--outdir", " "}, fs, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "'root-model' value passed is empty\n"+
		"'root-transform' value passed is empty\n"+
		"'model' value passed is empty\n"+
		"'grammar' value passed is empty\n"+
		"'start' value passed is empty\n"+
		"'outdir' value passed is empty\n", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyIntsParams(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()

	memFs, fs := testutil.WriteToMemOverlayFs(".")
	main2([]string{"sysl", "ints", "--root", " ", "-o", "", "-j", " ", "indirect_1.sysl"}, fs, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "'root' value passed is empty\n"+
		"'output' value passed is empty\n"+
		"'project' value passed is empty\n", testHook.LastEntry().Message)
	testutil.AssertFsHasExactly(t, memFs)
}
