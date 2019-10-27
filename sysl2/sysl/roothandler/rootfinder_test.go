package roothandler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type folderStructure struct {
	folders, files []string
}

type rootHandlerTestStructure struct {
	root, module, foundRoot, name string
	folders, files                []string
	rootIsFound                   bool
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

func absPathRelativeToCurrentDirectory(path string, t *testing.T) string {

	currentDirectory, err := filepath.Abs(".")
	assert.NoError(t, err)

	absPath, err := filepath.Abs(path)
	assert.NoError(t, err)

	relPath, err := filepath.Rel(currentDirectory, absPath)
	assert.NoError(t, err)
	return relPath
}

func TestRootHandler(t *testing.T) {

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

	definedRootFlagUndefinedMarker := folderStructure{
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

	tests := []rootHandlerTestStructure{
		{
			root:        "",
			name:        "Successful test: finding a root marker",
			module:      successfulTest.files[0],
			foundRoot:   "SuccessfulTest",
			folders:     successfulTest.folders,
			files:       successfulTest.files,
			rootIsFound: true,
		},
		{
			root:        "",
			name:        "Successful test: finding a root marker in the same directory as the module",
			module:      successfulTest.files[1],
			foundRoot:   "SuccessfulTest",
			folders:     successfulTest.folders,
			files:       successfulTest.files,
			rootIsFound: true,
		},
		{
			root:        "",
			name:        "Successful test: finding the closest root marker",
			module:      successfulTest.files[2],
			foundRoot:   "SuccessfulTest/path/to/another",
			folders:     successfulTest.folders,
			files:       successfulTest.files,
			rootIsFound: true,
		},
		{
			root:        "DefinedRootAndSyslRootUndefinedTest/path/",
			name:        "Root flag is defined and root marker does not exist",
			module:      definedRootFlagUndefinedMarker.files[0],
			foundRoot:   "DefinedRootAndSyslRootUndefinedTest/path",
			folders:     definedRootFlagUndefinedMarker.folders,
			files:       definedRootFlagUndefinedMarker.files,
			rootIsFound: true,
		},
		{
			root:        ".",
			name:        "Root with relative output",
			module:      definedRootFlagUndefinedMarker.files[0],
			foundRoot:   ".",
			folders:     definedRootFlagUndefinedMarker.folders,
			files:       definedRootFlagUndefinedMarker.files,
			rootIsFound: true,
		},
		{
			root:        "/",
			name:        "Root with absolute output",
			module:      definedRootFlagUndefinedMarker.files[0],
			foundRoot:   absPathRelativeToCurrentDirectory("/", t),
			folders:     definedRootFlagUndefinedMarker.folders,
			files:       definedRootFlagUndefinedMarker.files,
			rootIsFound: true,
		},
		{
			root:        "~",
			name:        "Root flag and Root marker is defined",
			module:      definedRootFlagAndMarkerFound.files[0],
			foundRoot:   absPathRelativeToCurrentDirectory("~", t),
			folders:     definedRootFlagAndMarkerFound.folders,
			files:       definedRootFlagAndMarkerFound.files,
			rootIsFound: true,
		},
		{
			root:        "",
			name:        "Root is not defined",
			module:      undefinedRoot.files[0],
			foundRoot:   "UndefinedRootAndUndefinedSyslRoot",
			folders:     undefinedRoot.folders,
			files:       undefinedRoot.files,
			rootIsFound: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(test rootHandlerTestStructure) func(t *testing.T) {
			return func(t *testing.T) {
				t.Parallel()
				logger := logrus.StandardLogger()
				fs, err := buildFolderTest(test.folders, test.files)
				assert.NoError(t, err)

				rootHandler := NewRootHandler(test.root, test.module)
				err = rootHandler.HandleRoot(fs, logger)

				assert.NoError(t, err)
				assert.Equal(t, test.foundRoot, rootHandler.Root())
				assert.Equal(t, test.rootIsFound, rootHandler.RootIsFound())
			}
		}(test))
	}
}

func TestImportAllowed(t *testing.T) {
	tests := []struct {
		name, module string
		handledRoot  *rootStatus
		result       bool
		err          error
	}{
		{
			name:   "Successful test, regular relative import",
			module: "test/module",
			handledRoot: &rootStatus{
				root:          "../",
				rootIsDefined: true,
			},
			result: true,
			err:    nil,
		},
		{
			name:   "Successful test, regular absolute import",
			module: "/test/module",
			handledRoot: &rootStatus{
				root:          "../",
				rootIsDefined: true,
			},
			result: true,
			err:    nil,
		},
		{
			name:   "Relative import with ..",
			module: "../test/module",
			handledRoot: &rootStatus{
				root:          "../",
				rootIsDefined: true,
			},
			result: false,
			err:    errors.New("import does not allow \"..\""),
		},
		{
			name:   "Absolute import with ..",
			module: "/../test/module",
			handledRoot: &rootStatus{
				root:          "../",
				rootIsDefined: true,
			},
			result: false,
			err:    errors.New("import does not allow \"..\""),
		},
		{
			name:   "Strange case, \"..\" in the middle",
			module: "/test/../module",
			handledRoot: &rootStatus{
				root:          "",
				rootIsDefined: false,
			},
			result: true,
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result, err := test.handledRoot.ImportAllowed(test.module)
			assert.Equal(t, test.result, result)
			assert.Equal(t, test.err, err)
		})
	}
}
