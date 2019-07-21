package pbutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	if filename := testutil.TempFilename(t, "", "sysl-TestJSONPB-*.json"); filename != "" {
		require.NoError(t, JSONPB(testModule, filename))
		output, err := ioutil.ReadFile(filename)
		require.NoError(t, err)
		assert.Equal(t, testModuleJSONPB, string(output))
	}
}

func TestJSONPBNilModule(t *testing.T) {
	if tf := testutil.NewTempFile(t, "", "sysl-TestJSONPB-*.json"); tf != nil {
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
	if filename := testutil.TempFilename(t, "", "sysl-TestJSONPB-*.json"); filename != "" {
		require.NoError(t, TextPB(testModule, filename))
		output, err := ioutil.ReadFile(filename)
		require.NoError(t, err)
		assert.Equal(t, testModuleTextPB, string(output))
	}
}

func TestTextPBNilModule(t *testing.T) {
	if tf := testutil.NewTempFile(t, "", "sysl-TestTextPBNilModule-*.textpb"); tf != nil {
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
