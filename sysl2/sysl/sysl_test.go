package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExit(t *testing.T) {
	format := "Exiting: %s"
	param := "Oopsies!"
	message := fmt.Sprintf(format, param)
	code := 42
	e := exitf(code, format, param)
	assert.Error(t, e)
	assert.Equal(t, message, e.message)
	assert.Equal(t, message, e.Error())
	assert.Equal(t, 42, e.code)
}

//nolint:gochecknoglobals
var (
	testModule = &sysl.Module{
		Apps: map[string]*sysl.Application{
			"Test": {
				Name: &sysl.AppName{
					Part: []string{"Test"},
				},
				Endpoints: map[string]*sysl.Endpoint{
					"GetInfo": {
						Name: "GetInfo",
						Stmt: []*sysl.Statement{
							{
								Stmt: &sysl.Statement_Action{
									Action: &sysl.Action{
										Action: "Do something",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	testModuleJSONPB = `{
 "apps": {
  "Test": {
   "name": {
    "part": [
     "Test"
    ]
   },
   "endpoints": {
    "GetInfo": {
     "name": "GetInfo",
     "stmt": [
      {
       "action": {
        "action": "Do something"
       }
      }
     ]
    }
   }
  }
 }
}`

	testModuleTextPB = `apps: <
  key: "Test"
  value: <
    name: <
      part: "Test"
    >
    endpoints: <
      key: "GetInfo"
      value: <
        name: "GetInfo"
        stmt: <
          action: <
            action: "Do something"
          >
        >
      >
    >
  >
>
`
)

func TestJSONPB(t *testing.T) {
	if filename := testTempFilename(t, "", "sysl-TestJSONPB-*.json"); filename != "" {
		require.NoError(t, JSONPB(testModule, filename))
		output, err := ioutil.ReadFile(filename)
		require.NoError(t, err)
		assert.Equal(t, testModuleJSONPB, string(output))
	}
}

func TestJSONPBNilModule(t *testing.T) {
	if tf := newTestTempFile(t, "", "sysl-TestJSONPB-*.json"); tf != nil {
		filename := tf.Name()
		tf.CloseAndRemove()
		require.Error(t, JSONPB(nil, filename))
		_, err := os.Stat(filename)
		assert.True(t, os.IsNotExist(err))
	}
}

func TestFJSONPB(t *testing.T) {
	var output bytes.Buffer
	require.NoError(t, FJSONPB(&output, testModule))
	assert.Equal(t, testModuleJSONPB, output.String())
}

func TestFJSONPBNilModule(t *testing.T) {
	var output bytes.Buffer
	require.Error(t, FJSONPB(&output, nil))
	assert.Equal(t, "", output.String())
}

func TestTextPB(t *testing.T) {
	if filename := testTempFilename(t, "", "sysl-TestJSONPB-*.json"); filename != "" {
		require.NoError(t, TextPB(testModule, filename))
		output, err := ioutil.ReadFile(filename)
		require.NoError(t, err)
		assert.Equal(t, testModuleTextPB, string(output))
	}
}

func TestTextPBNilModule(t *testing.T) {
	if tf := newTestTempFile(t, "", "sysl-TestTextPBNilModule-*.textpb"); tf != nil {
		filename := tf.Name()
		tf.CloseAndRemove()
		require.Error(t, TextPB(nil, filename))
		_, err := os.Stat(filename)
		assert.True(t, os.IsNotExist(err))
	}
}

func TestFTextPB(t *testing.T) {
	var output bytes.Buffer
	require.NoError(t, FTextPB(&output, testModule))
	assert.Equal(t, testModuleTextPB, output.String())
}

func TestFTextPBNilModule(t *testing.T) {
	var output bytes.Buffer
	require.Error(t, FTextPB(&output, nil))
	assert.Equal(t, "", output.String())
}

func testMain2(t *testing.T, args []string, golden string) {
	if output := testTempFilename(t, "", "sysl-TestTextPBNilModule-*.textpb"); output != "" {
		var stdout, stderr bytes.Buffer
		rc := main2(&stdout, &stderr, append([]string{"sysl", "-o", output}, args...), main3)
		if !assert.Zero(t, rc) {
			t.Error(stderr.String())
		}
		assert.True(t, stdout.Len() == 0)

		actual, err := ioutil.ReadFile(output)
		require.NoError(t, err)

		expected, err := ioutil.ReadFile(golden)
		require.NoError(t, err)

		assert.Equal(t, string(expected), string(actual))
	}
}

func TestMain2TextPB(t *testing.T) {
	testMain2(t, []string{"tests/args.sysl"}, "tests/args.sysl.golden.textpb")
}

func TestMain2JSON(t *testing.T) {
	testMain2(t, []string{"-mode", "json", "tests/args.sysl"}, "tests/args.sysl.golden.json")
}

func testMain2Stdout(t *testing.T, args []string, golden string) {
	var stdout, stderr bytes.Buffer
	rc := main2(&stdout, &stderr, append([]string{"sysl", "-o", "-"}, args...), main3)
	if !assert.Zero(t, rc) {
		t.Error(stderr.String())
	}

	expected, err := ioutil.ReadFile(golden)
	require.NoError(t, err)

	assert.Equal(t, string(expected), stdout.String())

	_, err = os.Stat("-")
	assert.True(t, os.IsNotExist(err), "Should not have created file '-'")
}

func TestMain2TextPBStdout(t *testing.T) {
	testMain2Stdout(t, []string{"tests/args.sysl"}, "tests/args.sysl.golden.textpb")
}

func TestMain2JSONStdout(t *testing.T) {
	testMain2Stdout(t, []string{"-mode", "json", "tests/args.sysl"}, "tests/args.sysl.golden.json")
}

func TestMain2BadMode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	rc := main2(&stdout, &stderr, []string{"sysl", "-o", "-", "-mode", "BAD", "tests/args.sysl"}, main3)
	assert.NotZero(t, rc)

	_, err := os.Stat("-")
	assert.True(t, os.IsNotExist(err), "Should not have created file '-'")
}

func TestMain2BadLog(t *testing.T) {
	var stdout, stderr bytes.Buffer
	rc := main2(&stdout, &stderr, []string{"sysl", "-o", "-", "-log", "BAD", "tests/args.sysl"}, main3)
	assert.NotZero(t, rc)

	_, err := os.Stat("-")
	assert.True(t, os.IsNotExist(err), "Should not have created file '-'")
}

func TestMain2WithBlackboxParams(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	wantStdout := "blackbox not hit: Server <- DB"

	main2(stdout, stderr, []string{"sd", "-s", "MobileApp <- Login", "-o", "tests/call.png", "-b",
		"Server <- DB=call to database", "-b", "Server <- Login=call to database", "tests/call.sysl"}, main3)

	if gotStdout := stdout.String(); strings.Compare(gotStdout, wantStdout) == 0 {
		t.Errorf("main2() = %v, want %v", gotStdout, wantStdout)
	}
}

func TestMain2WithBlackboxSysl(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	wantStdout := "blackbox 'SomeApp <- BarEndpoint' not hit in application 'Project :: Integrations :: PROJECT-E2E'"
	main2(stdout, stderr, []string{"sd", "-o", "%(epname).png", "tests/blackbox.sysl", "-a",
		"Project :: Integrations"}, main3)

	if gotStdout := stdout.String(); strings.Compare(gotStdout, wantStdout) == 0 {
		t.Errorf("main2() = %v, want %v", gotStdout, wantStdout)
	}
}

func TestMain2Fatal(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assert.Equal(t, 1, main2(&stdout, &stderr, nil, func(_, _ io.Writer, _ []string) error {
		return fmt.Errorf("Generic error")
	}))
	assert.Equal(t, "Generic error\n", stderr.String())
	stderr.Reset()
	assert.Equal(t, 42, main2(&stdout, &stderr, nil, func(_, _ io.Writer, _ []string) error {
		return exitf(42, "Exit error")
	}))
	assert.Equal(t, "Exit error\n", stderr.String())
}
