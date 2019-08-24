package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/testutil"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testMain2(t *testing.T, args []string, golden string) {
	if output := testutil.TempFilename(t, "", "sysl-TestTextPBNilModule-*.textpb"); output != "" {
		rc := main2(append([]string{"sysl", "pb", "-o", output}, args...), main3)
		assert.Zero(t, rc)

		actual, err := ioutil.ReadFile(output)
		assert.NoError(t, err)

		expected, err := ioutil.ReadFile(golden)
		assert.NoError(t, err)

		assert.Equal(t, string(expected), string(actual))
	}
}

func TestMain2TextPB(t *testing.T) {
	testMain2(t, []string{"tests/args.sysl"}, "tests/args.sysl.golden.textpb")
}

func TestMain2JSON(t *testing.T) {
	testMain2(t, []string{"--mode", "json", "tests/args.sysl"}, "tests/args.sysl.golden.json")
}

func testMain2Stdout(t *testing.T, args []string, golden string) {
	rc := main2(append([]string{"sysl", "pb", "-o", " - "}, args...), main3)
	assert.Zero(t, rc)

	_, err := ioutil.ReadFile(golden)
	require.NoError(t, err)

	_, err = os.Stat("-")
	assert.True(t, os.IsNotExist(err))
}

func TestMain2TextPBStdout(t *testing.T) {
	testMain2Stdout(t, []string{"tests/args.sysl"}, "tests/args.sysl.golden.textpb")
}

func TestMain2JSONStdout(t *testing.T) {
	testMain2Stdout(t, []string{"--mode", "json", "tests/args.sysl"}, "tests/args.sysl.golden.json")
}

func TestMain2BadMode(t *testing.T) {
	rc := main2([]string{"sysl", "pb", "-o", " - ", "--mode", "BAD", "tests/args.sysl"}, main3)
	assert.NotZero(t, rc)

	_, err := os.Stat("-")
	assert.True(t, os.IsNotExist(err))
}

func TestMain2BadLog(t *testing.T) {
	rc := main2([]string{"sysl", "pb", "-o", "-", "--log", "BAD", "tests/args.sysl"}, main3)
	assert.NotZero(t, rc)

	_, err := os.Stat("-")
	assert.True(t, os.IsNotExist(err))
}

func TestMain2SeqdiagWithMissingFile(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	rc := main2([]string{"sd", "-o", "%(epname).png", "tests/MISSING.sysl", "-a", "Project :: Sequences"}, main3)
	assert.NotEqual(t, 0, rc)
}

func TestMain2SeqdiagWithImpossibleOutput(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	rc := main2([]string{"sd", "-s", "MobileApp <- Login", "-o", "tests/call.zzz", "-b", "Server <- DB=call to database",
		"-b", "Server <- Login=call to database", "tests/call.sysl"}, main3)
	assert.NotEqual(t, 0, rc)
}

func TestMain2WithBlackboxParams(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	out := "tests/call1.png"
	defer os.Remove(out)
	main2([]string{"sysl", "sd", "-s", "MobileApp <- Login", "-o", out, "-b", "Server <- DB=call to database",
		"-b", "Server <- Login=call to database", "tests/call.sysl"}, main3)
	assert.Equal(t, logrus.WarnLevel, testHook.LastEntry().Level)
	assert.Equal(t, "blackbox 'Server <- DB' passed on commandline not hit\n", testHook.LastEntry().Message)
}

func TestMain2WithBlackboxParamsFaultyArguments(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	ret := main2([]string{"sysl", "sd", "-s", "MobileApp <- Login", "-o", "tests/call2.png", "-b", "Server <- DB",
		"-b", "Server <- Login", "tests/call.sysl"}, main3)
	assert.NotEqual(t, 0, ret)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "expected KEY=VALUE got 'Server <- DB'", testHook.LastEntry().Message)
}

func TestMain2WithBlackboxSysl(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	main2([]string{"sysl", "sd", "-o", "%(epname).png", "tests/blackbox.sysl", "-a", "Project :: Sequences"}, main3)
	assert.Equal(t, logrus.WarnLevel, testHook.LastEntry().Level)
	assert.Equal(t, "blackbox 'SomeApp <- AppEndpoint' not hit in app 'Project :: Sequences'\n",
		testHook.Entries[len(testHook.Entries)-1].Message)
	assert.Equal(t, "blackbox 'SomeApp <- BarEndpoint1' not hit in app 'Project :: Sequences :: SEQ-Two'\n",
		testHook.Entries[len(testHook.Entries)-2].Message)
	assert.Equal(t, "blackbox 'SomeApp <- BarEndpoint' not hit in app 'Project :: Sequences :: SEQ-One'\n",
		testHook.Entries[len(testHook.Entries)-3].Message)
}

func TestMain2WithBlackboxSyslEmptyEndpoints(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	main2([]string{"sysl", "sd", "-o", "%(epname).png", "tests/blackbox.sysl", "-a", "Project :: Integrations"}, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "No call statements to build sequence diagram for endpoint PROJECT-E2E", testHook.LastEntry().Message)
}

func TestMain2Fatal(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	assert.Equal(t, 42, main2(nil, func(_ []string) error {
		return parse.Exitf(42, "Exit error")
	}))
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
}

func TestMain2WithGroupingParamsGroupParamAbsent(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	main2([]string{"sysl", "sd", "-s", "MobileApp <- Login", "-g", "-o", "tests/call3.png", "tests/call.sysl"}, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "expected argument for flag '-g'", testHook.LastEntry().Message)
}

func TestMain2WithGroupingParamsCommandline(t *testing.T) {
	out := "tests/call4.png"
	defer os.Remove(out)
	main2([]string{"sysl", "sd", "-s", "MobileApp <- Login", "-g", "owner", "-o", out,
		"tests/call.sysl"}, main3)
	_, err := os.Stat(out)
	assert.NoError(t, err)
}

func TestMain2WithGroupingParamsSysl(t *testing.T) {
	testHook := test.NewGlobal()
	defer testHook.Reset()
	main2([]string{"sysl", "sd", "-g", "location", "-o", "%(epname).png", "tests/groupby.sysl", "-a",
		"Project :: Sequences"}, main3)
	for _, filename := range []string{"SEQ-One.png", "SEQ-Two.png"} {
		_, err := os.Stat(filename)
		assert.NoError(t, err)
		os.Remove(filename)
	}
	assert.Equal(t, logrus.WarnLevel, testHook.LastEntry().Level)
	assert.Equal(t, "Ignoring groupby passed from command line", testHook.LastEntry().Message)
}

func TestMain2WithGenerateIntegrations(t *testing.T) {
	out := "indirect_1.png"
	defer os.Remove(out)
	main2([]string{"sysl", "ints", "--root", "./tests/", "-o", out, "-j", "Project", "indirect_1.sysl"}, main3)
	_, err2 := os.Stat(out)
	assert.NoError(t, err2)
}

func TestMain2WithGenerateCode(t *testing.T) {
	defer os.Remove("Model.java")
	ret := main2([]string{"sysl", "gen", "--root-model", ".", "--model", "tests/model.sysl",
		"--root-transform", ".", "--transform", "tests/test.gen_multiple_annotations.sysl",
		"--grammar", "tests/test.gen.g", "--start", "javaFile"}, main3)
	assert.Equal(t, 0, ret)
}

func TestMain2WithTestPbJsonMode(t *testing.T) {
	out := "tests/callout1"
	defer os.Remove(out)
	ret := main2([]string{"sysl", "pb", "--mode", "textpb", "-o", out, "tests/call.sysl"}, main3)
	assert.Equal(t, 0, ret)
}

func TestMain2WithTestPbMode(t *testing.T) {
	out := "tests/callout2"
	defer os.Remove(out)
	ret := main2([]string{"sysl", "pb", "--mode", "json", "-o", out, "tests/call.sysl"}, main3)
	assert.Equal(t, 0, ret)
}

func TestMain2WithTestPbJsonConsole(t *testing.T) {
	ret := main2([]string{"sysl", "pb", "--mode", "textpb", "-o", " - ", "tests/call.sysl", "-v"}, main3)
	assert.Equal(t, 0, ret)
}

func TestMain2WithTestPbConsole(t *testing.T) {
	ret := main2([]string{"sysl", "pb", "--mode", "json", "-o", " - ", "tests/call.sysl", "-v"}, main3)
	assert.Equal(t, 0, ret)
}

func TestMain2WithEmptySdParams(t *testing.T) {
	testHook := test.NewGlobal()
	main2([]string{"sysl", "sd", "-g", " ", "-o", "", "tests/groupby.sysl", "-a", " "}, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "'output' value passed is empty\n"+
		"'groupby' value passed is empty\n"+
		"'app' value passed is empty\n", testHook.LastEntry().Message)
}

func TestMain2WithEmptyPbParams(t *testing.T) {
	testHook := test.NewGlobal()
	main2([]string{"sysl", "pb", "--root", "", "-o", " ", "--mode", "", "tests/call.sysl"}, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "'root' value passed is empty\n"+
		"'output' value passed is empty\n"+
		"'mode' value passed is empty\n", testHook.LastEntry().Message)
}

func TestMain2WithEmptyGenParams(t *testing.T) {
	testHook := test.NewGlobal()
	main2([]string{"sysl", "gen", "--root-model", " ", "--root-transform", "", "--model", " ", "--transform",
		"tests/test.gen_multiple_annotations.sysl", "--grammar", " ", "--start", "", "--outdir", " "}, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "'root-model' value passed is empty\n"+
		"'root-transform' value passed is empty\n"+
		"'model' value passed is empty\n"+
		"'grammar' value passed is empty\n"+
		"'start' value passed is empty\n"+
		"'outdir' value passed is empty\n", testHook.LastEntry().Message)
}

func TestMain2WithEmptyIntsParams(t *testing.T) {
	testHook := test.NewGlobal()
	main2([]string{"sysl", "ints", "--root", " ", "-o", "", "-j", " ", "indirect_1.sysl"}, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "'root' value passed is empty\n"+
		"'output' value passed is empty\n"+
		"'project' value passed is empty\n", testHook.LastEntry().Message)
}
