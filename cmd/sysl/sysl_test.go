package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	currentWorkingDirectory = "."
	syslRootMarker          = parse.SyslRootMarker
)

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
	rc := main2(append([]string{"sysl", "pb", "-o", output}, args...), fs, logger, os.Stdin, main3)
	assert.Zero(t, rc)

	actual, err := afero.ReadFile(fs, output)
	assert.NoError(t, err)

	expected, err := ioutil.ReadFile(golden)
	assert.NoError(t, err)

	expected = syslutil.HandleCRLF(expected)

	// In this test case, it will generate different value of file in sourceContext as current code which
	// loading Sysl module will evaluate absolute path of Sysl module, and the path relies on environment.
	// In order to make sure file value can be the same in all environments, add this code to reset file
	// value as "" in runtime.
	if strings.HasSuffix(golden, ".json") {
		reg := regexp.MustCompile(`"file": *"[^,\n]*"`)
		expectedStr := reg.ReplaceAllString(string(expected), `"file": ""`)
		actualStr := reg.ReplaceAllString(string(actual), `"file": ""`)
		assert.Equal(t, expectedStr, actualStr)
	} else if strings.HasSuffix(golden, ".textpb") {
		reg := regexp.MustCompile(`file: *"[^,\n]*"`)
		expectedStr := reg.ReplaceAllString(string(expected), `file: ""`)
		actualStr := reg.ReplaceAllString(string(actual), `file: ""`)
		// In protobuf text file, the space in case like `apps: {` is not fixed, it can be `apps: {` or
		// or `apps:  {`. So update it make sure it has only one space in this case.
		reg = regexp.MustCompile(`[^ ]+: +[^\n ]+`) // indent is 2 or more spaces
		splitReg := regexp.MustCompile(` +`)
		actualStr = reg.ReplaceAllStringFunc(actualStr, func(foundStr string) string {
			newStr := strings.Join(splitReg.Split(foundStr, -1), " ")
			return newStr
		})
		assert.Equal(t, expectedStr, actualStr)
	}
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
	testAllMain2WithRoot(t, args, inputFile, golden)
	testAllMain2WithoutRoot(t, args, inputFile, golden)
}

func testAllMain2WithoutRoot(t *testing.T, args []string, inputFile string, golden string) {
	// no root defined
	noRootFile := filepath.Join(testDir, inputFile)
	testMain2(t, append(args, noRootFile), filepath.Join(testDir, golden))

	// root marker
	testMain2WithSyslRootMarker(t, append(args, noRootFile), filepath.Join(testDir, golden))
}

func testAllMain2WithRoot(t *testing.T, args []string, inputFile string, golden string) {
	// root flag defined
	rootFile := "tests/" + inputFile

	args = append([]string{"--root", projDir}, args...)
	golden = filepath.Join(testDir, golden)
	golden = filepath.Join(filepath.Dir(golden), filepath.Base(golden))
	testMain2(t, append(args, rootFile), golden)

	testMain2WithSyslRootMarker(t, append(args, rootFile), golden)
}

func TestMain2_args_textpb(t *testing.T) {
	t.Parallel()

	testAllMain2(t, []string{}, "args.sysl", "args.sysl.golden.textpb")
}

func TestMain2_args_json(t *testing.T) {
	t.Parallel()

	testAllMain2(t, []string{"--mode", "json"}, "args.sysl", "args.sysl.golden.json")
}

func TestMain2_args_textpb_stdout(t *testing.T) {
	t.Parallel()

	testMain2Stdout(t, []string{filepath.Join(testDir, "args.sysl")}, filepath.Join(testDir, "args.sysl.golden.textpb"))
}

func TestMain2_args_json_stdout(t *testing.T) {
	t.Parallel()

	testMain2Stdout(t,
		[]string{"--mode", "json", filepath.Join(testDir, "args.sysl")},
		filepath.Join(testDir, "args.sysl.golden.json"),
	)
}

func TestMain2_type_merge1(t *testing.T) {
	t.Parallel()

	testAllMain2(t, []string{}, "type_merge1.sysl", "type_merge1.sysl.golden.textpb")
}

func TestMain2_type_merge2(t *testing.T) {
	t.Parallel()

	testAllMain2(t, []string{}, "type_merge2.sysl", "type_merge2.sysl.golden.textpb")
}

func TestMain2_namespace_merge(t *testing.T) {
	t.Parallel()

	testAllMain2(t, []string{}, "namespace_merge.sysl", "namespace_merge.sysl.golden.textpb")
}

func TestMain2_file_merge_textpb(t *testing.T) {
	t.Parallel()

	testAllMain2(t, []string{}, "file_merge.sysl", "file_merge.sysl.golden.textpb")
}

func TestMain2_file_merge_json(t *testing.T) {
	t.Parallel()

	testAllMain2(t, []string{"--mode", "json"}, "file_merge.sysl", "file_merge.sysl.golden.json")
}

func testMain2Stdout(t *testing.T, args []string, golden string) {
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	rc := main2(append([]string{"sysl", "pb", "-o", " - "}, args...), fs, logger, os.Stdin, main3)
	assert.Zero(t, rc)

	_, err := ioutil.ReadFile(golden)
	require.NoError(t, err)

	_, err = os.Stat("-")
	assert.True(t, os.IsNotExist(err))

	syslutil.AssertFsHasExactly(t, memFs)
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
		fs, logger, os.Stdin, main3,
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
		fs, logger, os.Stdin, main3,
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
		fs, logger, os.Stdin, main3,
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
		fs, logger, os.Stdin, main3,
	)
	assert.NotEqual(t, 0, rc)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	reg := regexp.MustCompile(`extension must be .[a-z]+,?|or not "zzz"`)
	assert.True(t, reg.MatchString(hook.LastEntry().Message))
}

func TestMain2WithBlackboxParams(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

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
		fs, logger, os.Stdin, main3,
	)
	assert.Equal(t, 0, rc)
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "blackbox 'Server <- DB' passed on commandline not hit\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithReadOnlyFs(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

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
		fs, logger, os.Stdin, main3,
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
		fs, logger, os.Stdin, main3,
	)
	assert.NotEqual(t, 0, ret)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "expected KEY=VALUE got 'Server <- DB'", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithBlackboxSysl(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

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
		fs, logger, os.Stdin, main3,
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
		fs, logger, os.Stdin, main3,
	)
	assert.NotEqual(t, 0, rc)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "no call statements to build sequence diagram for endpoint PROJECT-E2E", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2Fatal(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	assert.Equal(t, 42, main2(nil, nil, logger, os.Stdin,
		func(_ []string, _ afero.Fs, _ *logrus.Logger, stdio io.Reader) error {
			return syslutil.Exitf(42, "Exit error")
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
		fs, logger, os.Stdin, main3,
	)
	assert.NotEqual(t, 0, rc)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "expected argument for flag '-g'", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithGroupingParamsCommandline(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

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
		fs, logger, os.Stdin, main3,
	)
	assert.Equal(t, 0, rc)
	syslutil.AssertFsHasExactly(t, memFs, out)
}

func TestMain2WithGroupingParamsSysl(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

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
		fs, logger, os.Stdin, main3,
	)
	assert.Equal(t, 0, rc)
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "Ignoring groupby passed from command line", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs, "/SEQ-One.png", "/SEQ-Two.png")
}

func TestMain2WithGenerateIntegrations(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

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
		fs, logger, os.Stdin, main3,
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
			"--grammar", "test.gen.g",
			"--app-name", "Model",
			"--start", "javaFile",
			"--dep-path", "example.com/abc/asx/lmno/",
			"model.sysl",
		},
		fs, logger, os.Stdin, main3,
	)
	assert.Equal(t, 0, ret)
	out, err := filepath.Abs(filepath.Join(".", "Model.java"))
	require.NoError(t, err)
	syslutil.AssertFsHasExactly(t, memFs, out)
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
		fs, logger, os.Stdin, main3,
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
		fs, logger, os.Stdin, main3,
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
		fs, logger, os.Stdin, main3,
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
		fs, logger, os.Stdin, main3,
	)
	assert.Equal(t, 0, ret)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptySdParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "sd", "-g", " ", "-o", "", "tests/groupby.sysl", "-a", " "}, fs, logger, os.Stdin, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "'output' value passed is empty\n"+
		"'app' value passed is empty\n"+
		"'groupby' value passed is empty\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithValidateSyslNamespaces(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")
	ret := main2([]string{"sysl", "validate", filepath.Join(testDir, "apps_namespaces.sysl")}, fs, logger, os.Stdin, main3)
	assert.Equal(t, 0, ret)
	assert.False(t, len(hook.AllEntries()) > 0)
}

func TestMain2WithEmptyPbParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "pb", "-o", " ", "--mode", "", filepath.Join(testDir, "call.sysl")},
		fs, logger, os.Stdin, main3)
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
		"--grammar", " ", "--start", "", "--outdir", " "}, fs, logger, os.Stdin, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t,
		"'grammar' value passed is empty\n"+
			"'start' value passed is empty\n"+
			"'outdir' value passed is empty\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithEmptyIntsParams(t *testing.T) {
	t.Parallel()

	logger, hook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "ints", "-o", "", "-j", " ", "indirect_1.sysl"}, fs, logger, os.Stdin, main3)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t,
		"'output' value passed is empty\n"+
			"'project' value passed is empty\n", hook.LastEntry().Message)
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithDataMultipleFiles(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png",
		filepath.Join(testDir, "data.sysl"), "-j", "Project"}, fs, logger, os.Stdin, main3)
	syslutil.AssertFsHasExactly(t, memFs,
		"/Cross-Model.png", "/Relational-Model.png", "/Object-Model.png", "/Primitive-Alias-Model.png")
}

func TestMain2WithDataSingleFile(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "data.png",
		filepath.Join(testDir, "data.sysl"), "-j", "Project"}, fs, logger, os.Stdin, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/data.png")
}

func TestMain2WithDataNoProject(t *testing.T) {
	t.Parallel()
	logger, testHook := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png",
		filepath.Join(testDir, "data.sysl")}, fs, logger, os.Stdin, main3)
	assert.Equal(t, logrus.ErrorLevel, testHook.LastEntry().Level)
	assert.Equal(t, "project not found in sysl", testHook.LastEntry().Message)
	testHook.Reset()
	syslutil.AssertFsHasExactly(t, memFs)
}

func TestMain2WithDataFilter(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png", "-f",
		"Object-Model.png", filepath.Join(testDir, "data.sysl"), "-j",
		"Project"}, fs, logger, os.Stdin, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/Object-Model.png")
}

func TestMain2WithDataMultipleRelationships(t *testing.T) {
	t.Parallel()
	checkPlantUML(t)

	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "data", "-o", "%(epname).png",
		filepath.Join(testDir, "datareferences.sysl"), "-j", "Project"},
		fs, logger, os.Stdin, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/Relational-Model.png", "/Object-Model.png")
}

func TestMain2WithBinaryInfoCmd(t *testing.T) {
	t.Parallel()
	runSysl(t, 0, nil, "info")
}

func TestSwaggerExportCurrentDir(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "export", "-o", "SIMPLE_SWAGGER_EXAMPLE.yaml", "-a", "testapp",
		syslDir + "exporter/test-data/openapi2/SIMPLE_SWAGGER_EXAMPLE.sysl"}, fs, logger, os.Stdin, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/SIMPLE_SWAGGER_EXAMPLE.yaml")
}

func TestSwaggerExportTargetDir(t *testing.T) {
	t.Parallel()
	runSyslWithOutput(t, ".yaml", nil,
		"export", "-a", "testapp", path.Join(syslDir, "exporter/test-data/openapi2/SIMPLE_SWAGGER_EXAMPLE.sysl"))
}

func TestSwaggerExportJson(t *testing.T) {
	t.Parallel()
	runSyslWithOutput(t, ".json", nil,
		"export", "-a", "testapp", path.Join(syslDir, "exporter/test-data/openapi2/SIMPLE_SWAGGER_EXAMPLE.sysl"))
}

func TestSwaggerExportInvalid(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	_, fs := syslutil.WriteToMemOverlayFs("/")
	errInt := main2([]string{"sysl", "export", "-o", "SIMPLE_SWAGGER_EXAMPLE1.blah", "-a", "testapp",
		path.Join(syslDir, "exporter/test-data/openapi2/SIMPLE_SWAGGER_EXAMPLE.sysl")}, fs, logger, os.Stdin, main3)
	assert.True(t, errInt == 1)
}

func TestSwaggerAppExportNoDir(t *testing.T) {
	t.Parallel()
	outputDir := path.Join(t.TempDir(), "dirYetToBeCreated")
	runSysl(t, 0, nil, "export", "-o", path.Join(outputDir, "%(appname).yaml"),
		path.Join(syslDir, "exporter/test-data/openapi2/multiple/SIMPLE_SWAGGER_EXAMPLE_MULTIPLE.sysl"))
	assert.FileExists(t, path.Join(outputDir, "single.yaml"))
	assert.FileExists(t, path.Join(outputDir, "multiple.yaml"))
}

func TestSwaggerAppExportDirExists(t *testing.T) {
	t.Parallel()
	outputDir := t.TempDir()
	runSysl(t, 0, nil, "export", "-o", path.Join(outputDir, "%(appname).yaml"),
		path.Join(syslDir, "exporter/test-data/openapi2/multiple/SIMPLE_SWAGGER_EXAMPLE_MULTIPLE.sysl"))
	assert.FileExists(t, path.Join(outputDir, "single.yaml"))
	assert.FileExists(t, path.Join(outputDir, "multiple.yaml"))
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
			fs := syslutil.NewChrootFs(afero.NewMemMapFs(), "/")
			syslutil.BuildFolderTest(t, fs, ts.structure.folders, ts.structure.files)

			config := loader.NewProjectConfiguration()
			err := config.ConfigureProject(ts.root, ts.module, fs, logger)

			require.Equal(t, ts.rootFound, config.RootIsFound)
			require.NoError(t, err)
			require.Equal(t, ts.expectedRoot, config.Root)
			require.Equal(t, ts.getExpectedModule(t), config.Module)
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

/*
func TestCodegenGrammarImportDefOut(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "import", "-i",
		syslDir + "importer/tests-grammar/simplerules.gen.g", "-a", "go"}, fs, logger, os.Stdin, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/output.sysl")
}

func TestCodegenGrammarImport(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "import", "-i", syslDir + "importer/tests-grammar/unions.gen.g",
		"-o", "out.sysl", "-a", "go"}, fs, logger, os.Stdin, main3)
	syslutil.AssertFsHasExactly(t, memFs, "/out.sysl")
}
*/
func TestTemplating(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	memFs, fs := syslutil.WriteToMemOverlayFs("/")
	main2([]string{"sysl", "tmpl", "--root", "../../demo/codegen/AuthorisationAPI",
		"--root-template", "../../demo/codegen",
		"--template", "AuthorisationAPI/grpc.sysl", "--app-name", "AuthorisationAPI", "--start", "start",
		"--outdir", "../../demo/codegen/AuthorisationAPI/", "authorisation"}, fs, logger, os.Stdin, main3)
	outputFilename := "../../demo/codegen/AuthorisationAPI/AuthorisationAPI.proto"
	syslutil.AssertFsHasExactly(t, memFs, outputFilename)
	expected, err := ioutil.ReadFile(outputFilename)
	assert.NoError(t, err)
	expected = syslutil.HandleCRLF(expected)
	actual, err := afero.ReadFile(memFs, syslutil.MustAbsolute(t, outputFilename))
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(actual))
}

func TestSpannerSQLImport(t *testing.T) {
	t.Parallel()
	runSyslWithOutput(t, "", nil,
		"import", "--input", "../../pkg/importer/sql/tests/spanner/spanner.sql", "--app-name", "customeraccounts",
		"--package", "retail", "--format", "spannerSQL")
}

func TestSpannerSQLImportWithoutPackage(t *testing.T) {
	t.Parallel()
	runSyslWithOutput(t, "", nil,
		"import", "--input", "../../pkg/importer/sql/tests/spanner/spanner.sql", "--app-name", "customeraccounts",
		"--format", "spannerSQL")
}

func TestSpannerSQLImportDefOut(t *testing.T) {
	t.Parallel()
	runSyslWithOutput(t, "", nil,
		"import", "--input", "../../pkg/importer/sql/tests/spanner/spanner.sql", "--app-name", "customeraccounts",
		"--format", "spannerSQL", "--package", "retail")
}

func TestJsonSchemaImport(t *testing.T) {
	t.Parallel()
	runSyslWithExpectedOutput(t, "transforms/importers/jsonschema/expected.sysl",
		"import",
		"--input", "../../transforms/importers/jsonschema/input.json",
		"--app-name", "TestNamespace::TestApp",
		"--format", "JSONSchema")
}

func TestSpannerDirImport(t *testing.T) {
	t.Parallel()
	runSyslWithExpectedOutput(t, "pkg/importer/sql/tests/spanner/migrations/migrations.sysl",
		"import",
		"--input", "../../pkg/importer/sql/tests/spanner/migrations",
		"--app-name", "TestApp",
		"--format", "spannerSQLdir")
}

// runSyslWithExpectedOutput runs the Sysl command line tool with the specified arguments and adds an --output switch
// with a file directed at a temporary output directory. The content of the file are verified to be identical to the
// file specified in the expectedPathFromRepoRoot parameter (no need for "../../" in its path).
func runSyslWithExpectedOutput(t *testing.T, expectedPathFromRepoRoot string, args ...string) {
	expectedBytes, err := ioutil.ReadFile(path.Join("..", "..", expectedPathFromRepoRoot))
	require.NoError(t, err)
	expected := string(expectedBytes)

	actual := runSyslWithOutput(t, path.Ext(expectedPathFromRepoRoot), nil, args...)
	assert.Equal(t, expected, actual)
}

// runSyslWithOutput runs the Sysl command line tool with the specified arguments and adds an --output switch with a
// file directed at a temporary output directory with the specified file extension (if empty, '.out' is used). This file
// is then read and its contents are returned.
func runSyslWithOutput(t *testing.T, outFileExt string, stdin io.Reader, args ...string) string {
	if outFileExt == "" {
		outFileExt = ".out"
	}
	outputPath := filepath.Join(t.TempDir(), "output"+outFileExt)
	args = append(args, "--output", outputPath)
	runSysl(t, 0, stdin, args...)
	require.FileExists(t, outputPath)
	output, err := ioutil.ReadFile(outputPath)
	require.NoError(t, err)
	return string(output)
}

// runSysl runs the Sysl command line tool with the specified arguments (without needing 'sysl' as the first argument)
// and then ensures it completed with the specified return code.
func runSysl(t *testing.T, expectedRet int, stdin io.Reader, args ...string) {
	logger, hook := test.NewNullLogger()

	if args[0] != "sysl" {
		args = append([]string{"sysl"}, args...)
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	ret := main2(args, afero.NewOsFs(), logger, stdin, main3)

	lastEntry := hook.LastEntry()
	var lastMessage string
	if lastEntry != nil {
		lastMessage = lastEntry.Message
	}
	require.Equal(t, expectedRet, ret, lastMessage)
}

func TestMain3(t *testing.T) {
	logger, _ := test.NewNullLogger()
	fs := afero.NewOsFs()

	assert.Error(t, main3([]string{"sysl", "codegen"}, fs, logger, os.Stdin))

	assert.Error(t, main3([]string{"sysl", "codegen", "@tests/config.txt"}, fs, logger, os.Stdin))

	assert.Error(t, main3([]string{"sysl", "codegen", "@tests/config_new.txt"}, fs, logger, os.Stdin))

	assert.Error(t, main3([]string{"sysl", "codegen", "--grammar=go.gen.g", "--transform=go.gen.sysl", "model.sysl"},
		fs, logger, os.Stdin))
}

// Refers https://golang.org/pkg/testing/#hdr-Main
func TestMain(m *testing.M) {
	// Set mod.SyslModules = false to disable Sysl modules remote fetching for all test cases in main package.
	// Or test case executing will try to fetch remote sysl module and generate go.mod and go.sum files in
	// ./cmd/sysl folder, because this remote fetching uses go get. Generated go.mod and go.sum cause it can't
	// build Sysl binary file successfully in CI and local environments.
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestAvroImport(t *testing.T) {
	t.Parallel()
	logger, _ := test.NewNullLogger()
	outputPath := filepath.Join(t.TempDir(), "output.sysl")
	ret := main2([]string{"sysl", "import",
		"--input", filepath.Join(testDir, "/simple_avro.avsc"),
		"--app-name", "testapp",
		"--package", "test",
		"--output", outputPath}, afero.NewOsFs(), logger, os.Stdin, main3)
	assert.Equal(t, 0, ret)
	assert.FileExists(t, outputPath)
}

// Validate all test Sysl files with `sysl validate`.
func TestSyslSyntaxValidate(t *testing.T) {
	t.Parallel()

	dirs := []string{
		"../../pkg/importer/avro/tests",
		"../../pkg/importer/proto/tests",
		"../../pkg/importer/tests/grammar",
		"../../pkg/importer/tests/openapi3",
		"../../pkg/importer/tests/openapi2",
		"../../pkg/importer/tests/xsd",
		"../../tests"}

	for _, dir := range dirs {
		logger, _ := test.NewNullLogger()
		files, err := ioutil.ReadDir(dir)
		assert.NoError(t, err)
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".sysl") {
				t.Run(file.Name(), func(t *testing.T) {
					ret := main2([]string{"sysl", "validate",
						filepath.Join(dir, file.Name())}, afero.NewOsFs(), logger, os.Stdin, main3)
					assert.Equal(t, 0, ret)
				})
			}
		}
	}
}

// checkPlantUML causes a test to fail with a useful error message if a local PlantUML cannot be
// reached for tests that require it.
func checkPlantUML(t *testing.T) {
	plantURL := os.Getenv("SYSL_PLANTUML")
	res, err := http.Head(plantURL) //nolint:gosec
	require.NoError(t, err, "PlantUML is not running at %s", plantURL)
	_ = res.Body.Close()
}
