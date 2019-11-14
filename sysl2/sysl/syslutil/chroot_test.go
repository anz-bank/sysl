//nolint:funlen
package syslutil

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	placeholder          = "new"
	modeToUse            = os.ModePerm
	createIfDoesNotExist = os.O_RDWR | os.O_CREATE | os.O_EXCL
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

func (ts *testStructure) checkMockFsResult(
	t *testing.T, res interface{}, err error, mockFs *MockFs) {
	require.Equal(t, ts.expectedErr, err)
	if ts.expectedErr == nil {
		mockFs.AssertExpectations(t)
	} else {
		require.Empty(t, res)
	}
}

func (ts *testStructure) checkMockErrorOnly(
	t *testing.T, err error, mockFs *MockFs) {
	require.Equal(t, ts.expectedErr, err)
	if ts.expectedErr == nil {
		mockFs.AssertExpectations(t)
	}
}

func (ts *testStructure) checkResult(t *testing.T, res interface{}, err error) {
	require.Equal(t, ts.expectedErr, err)
	if ts.expectedErr == nil {
		require.NotEmpty(t, res)
	} else {
		require.Empty(t, res)
	}
}

func (ts *testStructure) getTestFs(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	ts.buildFolderTest(t, fs)
	return NewChrootFs(fs, ts.root)
}

func (ts *testStructure) addPlaceholderToModule() string {
	// adding a placeholder string to the module name to create so that it creates
	// a different file with the same filepath as the module
	return ts.module + placeholder
}

func (ts *testStructure) turnModuleIntoDir() string {
	// turns the test structure module into a directory by removing the extension
	// so that it creates a directory that retains the original module's path
	newDir := strings.TrimRightFunc(ts.module, func(r rune) bool {
		return r != '.'
	})
	return strings.TrimRight(newDir, ".") + placeholder
}

func mustJoin(t *testing.T, fs *ChrootFs, module string) string {
	joined, err := fs.join(module)
	require.NoError(t, err)
	return joined
}

func testAllCases(t *testing.T, instructions func(testStructure) func(t *testing.T)) {
	tests := generateCases()

	for _, test := range tests {
		t.Run(test.name, instructions(test))
	}
}

func TestChrootFsCreate(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			res, err := chrootFs.Create(ts.addPlaceholderToModule())
			ts.checkResult(tt, res, err)
		}
	})
}

func TestChrootFsMkdir(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			err := chrootFs.Mkdir(ts.turnModuleIntoDir(), os.ModeTemporary)
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}

func TestChrootFsMkdirAll(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			err := chrootFs.MkdirAll(ts.turnModuleIntoDir(), os.ModeTemporary)
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}
func TestChrootFsOpen(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			res, err := chrootFs.Open(ts.module)
			ts.checkResult(tt, res, err)
		}
	})
}

func TestChrootFsOpenFile(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			res, err := chrootFs.OpenFile(ts.module, os.O_RDONLY, modeToUse)
			ts.checkResult(tt, res, err)

			// open files that does not exist
			res, err = chrootFs.OpenFile(ts.addPlaceholderToModule(),
				createIfDoesNotExist, modeToUse)
			ts.checkResult(tt, res, err)
		}
	})
}

func TestChrootFsRemove(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			err := chrootFs.Remove(ts.module)
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}

func TestChrootFsRemoveAll(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			err := chrootFs.RemoveAll(filepath.Dir(ts.module))
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}

func TestChrootFsRename(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			err := chrootFs.Rename(ts.module, ts.addPlaceholderToModule())
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}

func TestChrootFsStat(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			res, err := chrootFs.Stat(ts.module)
			ts.checkResult(tt, res, err)
		}
	})
}

func TestChrootFsChmod(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			err := chrootFs.Chmod(ts.module, modeToUse)
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}

func TestChrootFsChtimes(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			chrootFs := ts.getTestFs(tt)
			now := time.Now()
			err := chrootFs.Chtimes(ts.module, now, now)
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}

func TestMockChrootFsCreate(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			mockedFs.On("Create", mustJoin(tt, chrootFs, ts.module)).Return(struct{ afero.File }{}, nil)
			res, err := chrootFs.Create(ts.module)
			ts.checkMockFsResult(tt, res, err, mockedFs)
		}
	})
}

func TestMockChrootFsMkdir(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			newDir := ts.turnModuleIntoDir()
			mockedFs.On("Mkdir", mustJoin(tt, chrootFs, newDir), modeToUse).Return(nil)
			err := chrootFs.Mkdir(newDir, modeToUse)
			ts.checkMockErrorOnly(tt, err, mockedFs)
		}
	})
}

func TestMockChrootFsMkdirAll(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			newDir := ts.turnModuleIntoDir()
			mockedFs.On("MkdirAll", mustJoin(tt, chrootFs, newDir), modeToUse).Return(nil)
			err := chrootFs.MkdirAll(newDir, modeToUse)
			ts.checkMockErrorOnly(tt, err, mockedFs)
		}
	})
}

func TestMockChrootFsOpen(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			mockedFs.On("Open", mustJoin(tt, chrootFs, ts.module)).Return(struct{ afero.File }{}, nil)
			res, err := chrootFs.Open(ts.module)
			ts.checkMockFsResult(tt, res, err, mockedFs)
		}
	})
}

func TestMockChrootFsOpenFile(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			mockedFs.On(
				"OpenFile",
				mustJoin(tt, chrootFs, ts.module),
				createIfDoesNotExist,
				modeToUse).Return(struct{ afero.File }{}, nil)

			// only create files if it does not exist
			res, err := chrootFs.OpenFile(ts.module, createIfDoesNotExist, modeToUse)
			ts.checkMockFsResult(tt, res, err, mockedFs)
		}
	})
}

func TestMockChrootFsRemove(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			mockedFs.On("Remove", mustJoin(tt, chrootFs, ts.module)).Return(nil)
			err := chrootFs.Remove(ts.module)
			ts.checkMockErrorOnly(tt, err, mockedFs)
		}
	})
}

func TestMockChrootFsRename(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			newName := ts.addPlaceholderToModule()
			mockedFs.On(
				"Rename",
				mustJoin(tt, chrootFs, ts.module),
				mustJoin(tt, chrootFs, newName)).Return(nil)
			err := chrootFs.Rename(ts.module, newName)
			ts.checkMockErrorOnly(tt, err, mockedFs)
		}
	})
}

func TestMockChrootFsStat(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			mockedFs.On("Stat", mustJoin(tt, chrootFs, ts.module)).Return(struct{ os.FileInfo }{}, nil)
			res, err := chrootFs.Stat(ts.module)
			ts.checkMockFsResult(tt, res, err, mockedFs)
		}
	})
}

func TestMockChrootFsChmod(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			mockedFs.On("Chmod", mustJoin(tt, chrootFs, ts.module), modeToUse).Return(nil)
			err := chrootFs.Chmod(ts.module, modeToUse)
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}

func TestMockChrootFsChtimes(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFs, chrootFs := NewMockChrootFs(ts.root)
			now := time.Now()
			mockedFs.On("Chtimes", mustJoin(tt, chrootFs, ts.module), now, now).Return(nil)
			err := chrootFs.Chtimes(ts.module, now, now)
			require.Equal(tt, ts.expectedErr, err)
		}
	})
}

func TestWrapCall(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFunction := new(mockFunctionCall)
			chrootFs := NewChrootFs(ts.getTestFs(tt), ts.root)

			joined := mustJoin(tt, chrootFs, ts.module)

			mockedFunction.On("mockCall", joined).Return(nil)
			err := chrootFs.wrapCall(ts.module, mockedFunction.mockCall)
			require.Equal(tt, ts.expectedErr, err)
			if ts.expectedErr == nil {
				mockedFunction.AssertExpectations(t)
			}
		}
	})
}

func TestWrapCallWithData(t *testing.T) {
	testAllCases(t, func(ts testStructure) func(t *testing.T) {
		return func(tt *testing.T) {
			tt.Parallel()

			mockedFunction := new(mockFunctionCall)
			chrootFs := NewChrootFs(ts.getTestFs(tt), ts.root)

			joined := mustJoin(tt, chrootFs, ts.module)

			mockedFunction.On("mockCallWithData", joined).Return(struct{}{}, nil)
			res, err := chrootFs.wrapCallWithData(ts.module, mockedFunction.mockCallWithData)

			require.Equal(t, ts.expectedErr, err)
			if ts.expectedErr == nil {
				mockedFunction.AssertExpectations(t)
			} else {
				require.Empty(t, res)
			}
		}
	})
}

func (m *mockFunctionCall) mockCall(path string) error {
	return m.Called(path).Error(0)
}

func (m *mockFunctionCall) mockCallWithData(path string) (interface{}, error) {
	res := m.Called(path)
	return res.Get(0), res.Error(1)
}
