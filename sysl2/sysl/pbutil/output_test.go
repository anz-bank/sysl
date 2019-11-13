package pbutil

import (
	"bytes"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testModule() *sysl.Module {
	return &sysl.Module{
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
}

func testModuleJSONPB() string {
	return `{
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
}

func testModuleTextPB() string {
	return `apps: <
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
}

func TestJSONPB(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	filename := "out.pb.json"
	require.NoError(t, JSONPB(testModule(), filename, fs))
	output, err := afero.ReadFile(fs, filename)
	require.NoError(t, err)
	assert.Equal(t, testModuleJSONPB(), string(output))
}

func TestJSONPBNilModule(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	filename := "out.pb.json"
	require.Error(t, JSONPB(nil, filename, fs))
	syslutil.AssertFsHasExactly(t, fs)
}

func TestFJSONPB(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	require.NoError(t, FJSONPB(&output, testModule()))
	assert.Equal(t, testModuleJSONPB(), output.String())
}

func TestFJSONPBNilModule(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	require.Error(t, FJSONPB(&output, nil))
	assert.Equal(t, "", output.String())
}

func TestTextPB(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	filename := "/out.textpb"
	require.NoError(t, TextPB(testModule(), filename, fs))
	output, err := afero.ReadFile(fs, filename)
	require.NoError(t, err)
	assert.Equal(t, testModuleTextPB(), string(output))
}

func TestTextPBNilModule(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	filename := "/out.textpb"
	require.Error(t, TextPB(nil, filename, fs))
	syslutil.AssertFsHasExactly(t, fs)
}

func TestFTextPB(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	require.NoError(t, FTextPB(&output, testModule()))
	assert.Equal(t, testModuleTextPB(), output.String())
}

func TestFTextPBNilModule(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	require.Error(t, FTextPB(&output, nil))
	assert.Equal(t, "", output.String())
}
