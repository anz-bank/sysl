package pbutil

import (
	"bytes"
	"testing"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/anz-bank/sysl/pkg/sysl"
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

func TestJSONPB(t *testing.T) {
	t.Parallel()

	unmarshalled := &sysl.Module{}
	fs := afero.NewMemMapFs()
	filename := "out.pb.json"
	require.NoError(t, JSONPB(testModule(), filename, fs))
	output, err := afero.ReadFile(fs, filename)
	require.NoError(t, err)
	require.NoError(t, protojson.Unmarshal(output, unmarshalled))
	assert.True(t, proto.Equal(unmarshalled, protoreflect.ProtoMessage(testModule())))
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

	unmarshalled := &sysl.Module{}
	var output bytes.Buffer
	require.NoError(t, FJSONPB(&output, testModule()))
	require.NoError(t, protojson.Unmarshal(output.Bytes(), unmarshalled))
	assert.True(t, proto.Equal(unmarshalled, protoreflect.ProtoMessage(testModule())))
}

func TestFJSONPBNilModule(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	require.Error(t, FJSONPB(&output, nil))
	assert.Equal(t, "", output.String())
}

func TestTextPB(t *testing.T) {
	t.Parallel()

	unmarshalled := &sysl.Module{}
	fs := afero.NewMemMapFs()
	filename := "/out.textpb"
	require.NoError(t, TextPB(testModule(), filename, fs))
	output, err := afero.ReadFile(fs, filename)
	require.NoError(t, err)
	require.NoError(t, prototext.Unmarshal(output, unmarshalled))
	assert.True(t, proto.Equal(unmarshalled, protoreflect.ProtoMessage(testModule())))
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

	unmarshalled := &sysl.Module{}
	var output bytes.Buffer
	require.NoError(t, FTextPB(&output, testModule()))
	require.NoError(t, prototext.Unmarshal(output.Bytes(), unmarshalled))
	assert.True(t, proto.Equal(unmarshalled, protoreflect.ProtoMessage(testModule())))
}

func TestFTextPBNilModule(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	require.Error(t, FTextPB(&output, nil))
	assert.Equal(t, "", output.String())
}

func TestGeneratePBBinaryMessage(t *testing.T) {
	t.Parallel()

	unmarshalled := &sysl.Module{}
	var output bytes.Buffer
	require.NoError(t, GeneratePBBinaryMessage(&output, testModule()))
	require.NoError(t, proto.Unmarshal(output.Bytes(), unmarshalled))
	assert.True(t, proto.Equal(unmarshalled, protoreflect.ProtoMessage(testModule())))
}

func TestGeneratePBBinaryMessageNilModule(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	require.Error(t, GeneratePBBinaryMessage(&output, nil))
	assert.Equal(t, "", output.String())
}

func TestGeneratePBBinaryMessageFile(t *testing.T) {
	t.Parallel()

	unmarshalled := &sysl.Module{}
	fs := afero.NewMemMapFs()
	filename := "/out.pb"
	require.NoError(t, GeneratePBBinaryMessageFile(testModule(), filename, fs))
	output, err := afero.ReadFile(fs, filename)
	require.NoError(t, err)
	require.NoError(t, proto.Unmarshal(output, unmarshalled))
	assert.True(t, proto.Equal(unmarshalled, protoreflect.ProtoMessage(testModule())))
}

func TestGeneratePBBinaryMessageFileNilModule(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	filename := "/out.pb"
	require.Error(t, GeneratePBBinaryMessageFile(nil, filename, fs))
	syslutil.AssertFsHasExactly(t, fs)
}
