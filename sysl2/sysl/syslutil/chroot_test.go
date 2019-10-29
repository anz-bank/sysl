package syslutil

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

type testStructure struct {
	name, root, module string
	folders, files     []string
	expectedErr        error
}

func buildFolderTest(t *testing.T, folders, files []string) (fs afero.Fs) {
	fs = afero.NewMemMapFs()

	for _, folder := range folders {
		folder, err := filepath.Abs(folder)
		require.NoError(t, err)

		err = fs.MkdirAll(folder, os.ModeTemporary)
		require.NoError(t, err)
	}

	for _, file := range files {
		file, err := filepath.Abs(file)
		require.NoError(t, err)

		_, err = fs.Create(file)
		require.NoError(t, err)
	}

	return
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

func TestOpen(t *testing.T) {
	tests := generateCases()

	for _, test := range tests {
		t.Run(test.name, func(ts testStructure) func(t *testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()
				fs := buildFolderTest(tt, ts.folders, ts.files)
				chrootfs := NewChrootFs(fs, ts.root)

				_, err := chrootfs.Open(ts.module)
				require.Equal(tt, ts.expectedErr, err)
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
				fs := buildFolderTest(tt, ts.folders, ts.files)
				chrootfs := NewChrootFs(fs, ts.root)

				_, err := chrootfs.Stat(ts.module)
				require.Equal(tt, ts.expectedErr, err)
			}
		}(test))
	}
}
