package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	log "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
)

const currentWorkingDirectory = "."

type folderTestStructure struct {
	name,
	module,
	root,
	expectedRoot,
	rootMarkerPath string
	structure folderStructure
}

type folderStructure struct {
	folders, files []string
}

func TestEnsureFlagsNonEmpty_AllowsExcludes(t *testing.T) {
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := sysl.Command("foo", "")
	_ = cmd.Flag("bar", "").Default("foo").String()
	_ = cmd.Flag("other", "").Default("foo").String()

	EnsureFlagsNonEmpty(cmd, "bar")

	args := []string{"foo", "--bar", ""}
	selected, err := sysl.Parse(args)
	assert.Equal(t, "foo", selected)
	assert.NoError(t, err)
}

func TestEnsureFlagsNonEmpty(t *testing.T) {
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	cmd := sysl.Command("foo", "")
	cmd.Flag("bar", "").Default("foo")

	EnsureFlagsNonEmpty(cmd)

	args := []string{"foo", "--bar", ""}
	_, err := sysl.ParseContext(args)
	assert.Error(t, err)
}

func TestSetProjectRoot(t *testing.T) {
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
		},
		{
			name:         "Successful test: finding a root marker in the same directory as the module",
			root:         "",
			module:       successfulTest.files[1],
			structure:    successfulTest,
			expectedRoot: syslutil.MustAbsolute(t, "SuccessfulTest"),
		},
		{
			name:         "Successful test: finding the closest root marker",
			root:         "",
			module:       successfulTest.files[2],
			structure:    successfulTest,
			expectedRoot: syslutil.MustAbsolute(t, "SuccessfulTest/path/to/another"),
		},
		{
			name: "Root flag is defined and root marker does not exist",
			root: "DefinedRootAndSyslRootUndefinedTest/path/",
			module: syslutil.MustRelative(t, "DefinedRootAndSyslRootUndefinedTest/path/",
				definedRootNoMarker.files[0]),
			structure:    definedRootNoMarker,
			expectedRoot: "DefinedRootAndSyslRootUndefinedTest/path/",
		},
		{
			name:         "Defined relative root",
			root:         currentWorkingDirectory,
			module:       filepath.Clean(definedRootNoMarker.files[0]),
			structure:    definedRootNoMarker,
			expectedRoot: currentWorkingDirectory,
		},
		{
			root:         systemRoot,
			name:         "Defined absolute path root",
			module:       syslutil.MustAbsolute(t, definedRootNoMarker.files[0]),
			structure:    definedRootNoMarker,
			expectedRoot: systemRoot,
		},
		{
			name:         "Defined relative root with absolute module path rooted at root",
			root:         currentWorkingDirectory,
			module:       filepath.Join(systemRoot, filepath.Clean(definedRootNoMarker.files[0])),
			structure:    definedRootNoMarker,
			expectedRoot: currentWorkingDirectory,
		},
		{
			name:           "Defined root flag and root",
			root:           currentWorkingDirectory,
			module:         syslutil.MustRelative(t, currentWorkingDirectory, definedRootFlagAndMarkerFound.files[0]),
			structure:      definedRootFlagAndMarkerFound,
			expectedRoot:   currentWorkingDirectory,
			rootMarkerPath: syslutil.MustAbsolute(t, "./DefinedRootAndSyslRootDefinedTest/path/"),
		},
		{
			name:           "Defined root flag and root marker with absolute path module rooted at root",
			root:           "./DefinedRootAndSyslRootDefinedTest/",
			module:         "/path/to/module/test.sysl",
			structure:      definedRootFlagAndMarkerFound,
			expectedRoot:   "./DefinedRootAndSyslRootDefinedTest/",
			rootMarkerPath: syslutil.MustAbsolute(t, "./DefinedRootAndSyslRootDefinedTest/path/"),
		},
		{
			name:         "Root is not defined",
			root:         "",
			module:       undefinedRoot.files[0],
			structure:    undefinedRoot,
			expectedRoot: filepath.Dir(undefinedRoot.files[0]),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			logger, _ := log.NewNullLogger()
			fs := afero.NewMemMapFs()
			syslutil.BuildFolderTest(t, fs, test.structure.folders, test.structure.files)

			r := &cmdRunner{Root: test.root, module: test.module}

			require.NoError(t, r.setProjectRoot(fs, logger))
			require.Equal(t, test.expectedRoot, r.Root)
			require.Equal(t, test.getExpectedModule(t), r.module)
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
