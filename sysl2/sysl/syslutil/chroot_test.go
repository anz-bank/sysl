package syslutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockFunctionCall struct {
	mock.Mock
}
type testStructure struct {
	name, root, module string
	folders, files     []string
	expectedErr        error
}

func generateCases() []testStructure {
	return []testStructure{
		{
			name:   "import above root",
			root:   "/a",
			module: "../bad.sysl",
			folders: []string{
				"/a",
			},
			files: []string{
				"/bad.sysl",
			},
			expectedErr: errors.New("permission denied, file outside root"),
		},
		{
			name:   "successful import",
			root:   "/a",
			module: "/good.sysl",
			folders: []string{
				"/a",
			},
			files: []string{
				"/a/good.sysl",
			},
			expectedErr: nil,
		},
		{
			name:   "deeply nested structure with system root as root",
			root:   "/",
			module: "../../bar.sysl",
			folders: []string{
				"/a/b/c/d/e",
			},
			files: []string{
				"/bar.sysl",
				"/a/b/c/d/e/foo.sysl",
			},
			expectedErr: nil,
		},
		{
			name:   "deeply nested structure with directory as a root",
			root:   "/a/b",
			module: "../../bar.sysl",
			folders: []string{
				"/a/b/c/d/e",
			},
			files: []string{
				"/bar.sysl",
				"/a/b/c/d/e/foo.sysl",
			},
			expectedErr: errors.New("permission denied, file outside root"),
		},
		{
			name:   "relative import",
			root:   ".",
			module: "/good.sysl",
			folders: []string{
				".",
			},
			files: []string{
				"./good.sysl",
			},
			expectedErr: nil,
		},
		{
			name:   "failed relative import",
			root:   "./a/b",
			module: "../bad.sysl",
			folders: []string{
				"./a/b",
			},
			files: []string{
				"./a/bad.sysl",
			},
			expectedErr: errors.New("permission denied, file outside root"),
		},
		{
			name:   "strange file name",
			root:   "../a",
			module: "/..strange.sysl",
			folders: []string{
				"../a/",
			},
			files: []string{
				"../a/..strange.sysl",
			},
			expectedErr: nil,
		},
		{
			name:   "strange path",
			root:   "../a",
			module: "/b/../..strange.sysl",
			folders: []string{
				"../a/b",
			},
			files: []string{
				"../a/..strange.sysl",
			},
			expectedErr: nil,
		},
	}
}

func (ts *testStructure) buildFolderTest(t *testing.T, fs afero.Fs) {
	for _, folder := range ts.folders {
		folder, err := filepath.Abs(folder)
		require.NoError(t, err)

		err = fs.MkdirAll(folder, os.ModeTemporary)
		require.NoError(t, err)
	}

	for _, file := range ts.files {
		file, err := filepath.Abs(file)
		require.NoError(t, err)

		_, err = fs.Create(file)
		require.NoError(t, err)
	}
}

func (ts *testStructure) checkResultWithData(
	t *testing.T, res interface{}, err error, mockFs *MockFs) {
	require.Equal(t, ts.expectedErr, err)
	if ts.expectedErr == nil {
		require.NotEmpty(t, res)
		mockFs.AssertExpectations(t)
	} else {
		require.Empty(t, res)
	}
}

func (ts *testStructure) getSetUpMemMapFs(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	ts.buildFolderTest(t, fs)
	return fs
}

func (ts *testStructure) getSetUpMockFs(t *testing.T) (*MockFs, *ChrootFs) {
	memFs := ts.getSetUpMemMapFs(t)
	fs := NewMockFs(t, memFs)
	mockChrootFs := NewChrootFs(fs, ts.root)
	return fs, mockChrootFs
}

func (ts *testStructure) getJoined(t *testing.T, fs *ChrootFs) string {
	joined, err := fs.join(ts.module)
	require.NoError(t, err)
	return joined
}

func TestOpen(t *testing.T) {
	tests := generateCases()

	for _, test := range tests {
		t.Run(test.name, func(ts testStructure) func(t *testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				mockedFs, chrootfs := ts.getSetUpMockFs(tt)
				mockedFs.On("Open", ts.getJoined(tt, chrootfs)).Return(mock.AnythingOfType("afero.File"), nil)
				res, err := chrootfs.Open(ts.module)
				require.Equal(tt, ts.expectedErr, err)
				ts.checkResultWithData(tt, res, err, mockedFs)
			}
		}(test))
	}
}

func TestOpenFile(t *testing.T) {
	tests := generateCases()

	for _, test := range tests {
		t.Run(test.name, func(ts testStructure) func(t *testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				mockedFs, chrootfs := ts.getSetUpMockFs(tt)
				mockedFs.On(
					"OpenFile",
					ts.getJoined(tt, chrootfs),
					os.O_RDONLY,
					os.ModePerm).Return(mock.AnythingOfType("afero.File"), nil)
				res, err := chrootfs.OpenFile(ts.module, os.O_RDONLY, os.ModePerm)
				ts.checkResultWithData(t, res, err, mockedFs)
			}
		}(test))
	}
}
func TestStat(t *testing.T) {
	tests := generateCases()

	for _, test := range tests {
		t.Run(test.name, func(ts testStructure) func(t *testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				mockedFs, chrootfs := ts.getSetUpMockFs(tt)
				mockedFs.On("Stat", ts.getJoined(tt, chrootfs)).Return(mock.AnythingOfType("os.FileInfo"), nil)
				res, err := chrootfs.Stat(ts.module)
				ts.checkResultWithData(t, res, err, mockedFs)
			}
		}(test))
	}
}

func TestWrapCall(t *testing.T) {
	tests := generateCases()

	for _, test := range tests {
		t.Run(test.name, func(ts testStructure) func(t *testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				mockedFunction := new(mockFunctionCall)
				chrootfs := NewChrootFs(ts.getSetUpMemMapFs(tt), ts.root)

				joined := ts.getJoined(tt, chrootfs)

				mockedFunction.On("mockCall", joined).Return(nil)
				err := chrootfs.wrapCall(ts.module, mockedFunction.mockCall)
				require.Equal(t, ts.expectedErr, err)
			}
		}(test))
	}
}

func TestWrapCallWithData(t *testing.T) {
	tests := generateCases()

	for _, test := range tests {
		t.Run(test.name, func(ts testStructure) func(t *testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				mockedFunction := new(mockFunctionCall)
				chrootfs := NewChrootFs(ts.getSetUpMemMapFs(tt), ts.root)

				joined := ts.getJoined(tt, chrootfs)

				mockedFunction.On("mockCallWithData", joined).Return(mock.Anything, nil)
				res, err := chrootfs.wrapCallWithData(ts.module, mockedFunction.mockCallWithData)

				require.Equal(t, ts.expectedErr, err)
				if ts.expectedErr == nil {
					require.NotEmpty(t, res)
					mockedFunction.AssertExpectations(t)
				} else {
					require.Empty(t, res)
				}
			}
		}(test))
	}
}

func (m *mockFunctionCall) mockCall(path string) error {
	return m.Called(path).Error(0)
}

func (m *mockFunctionCall) mockCallWithData(path string) (interface{}, error) {
	res := m.Called(path)
	return res.Get(0), res.Error(1)
}
