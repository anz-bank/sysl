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

func buildFolderTest(folders, files []string) (fs afero.Fs, err error) {
	fs = afero.NewMemMapFs()
	var folder, file string

	for _, folder = range folders {
		folder, err = filepath.Abs(folder)
		if err != nil {
			return
		}

		err = fs.MkdirAll(folder, os.ModeTemporary)
		if err != nil {
			return
		}
	}

	for _, file = range files {
		file, err = filepath.Abs(file)
		if err != nil {
			return
		}

		_, err = fs.Create(file)
		if err != nil {
			return
		}
	}

	return
}

func TestOpen(t *testing.T) {
	tests := []testStructure{
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

	for _, test := range tests {
		t.Run(test.name, func(ts testStructure) func(t *testing.T) {
			return func(tt *testing.T) {
				fs, err := buildFolderTest(ts.folders, ts.files)
				tt.Parallel()
				require.NoError(tt, err)
				chrootfs := NewChrootFs(fs, ts.root)
				_, err = chrootfs.Open(ts.module)
				require.Equal(tt, ts.expectedErr, err)
			}
		}(test))
	}
}
