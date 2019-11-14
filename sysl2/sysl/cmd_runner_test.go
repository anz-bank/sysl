package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	currentWorkingDirectory          = "."
	noRootMarkerLog         logCases = iota
	rootMarkerExistsLog
	noRootLog
)

type logCases int
type folderTestStructure struct {
	name, module, root, expectedRoot, rootMarkerPath string
	expectedLog                                      logCases
	rootMarkerExists                                 bool
	structure                                        folderStructure
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

func TestGetProjectRoot(t *testing.T) {
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
			name:             "Successful test: finding a root marker",
			root:             "",
			module:           successfulTest.files[0],
			structure:        successfulTest,
			expectedRoot:     syslutil.MustAbsolute(t, "SuccessfulTest"),
			rootMarkerExists: true,
		},
		{
			name:             "Successful test: finding a root marker in the same directory as the module",
			root:             "",
			module:           successfulTest.files[1],
			structure:        successfulTest,
			expectedRoot:     syslutil.MustAbsolute(t, "SuccessfulTest"),
			rootMarkerExists: true,
		},
		{
			name:             "Successful test: finding the closest root marker",
			root:             "",
			module:           successfulTest.files[2],
			structure:        successfulTest,
			expectedRoot:     syslutil.MustAbsolute(t, "SuccessfulTest/path/to/another"),
			rootMarkerExists: true,
		},
		{
			name: "Root flag is defined and root marker does not exist",
			root: "DefinedRootAndSyslRootUndefinedTest/path/",
			module: syslutil.MustRelative(t, "DefinedRootAndSyslRootUndefinedTest/path/",
				definedRootNoMarker.files[0]),
			structure:        definedRootNoMarker,
			expectedRoot:     "DefinedRootAndSyslRootUndefinedTest/path/",
			expectedLog:      noRootMarkerLog,
			rootMarkerExists: false,
		},
		{
			name:             "Defined relative root",
			root:             currentWorkingDirectory,
			module:           filepath.Clean(definedRootNoMarker.files[0]),
			structure:        definedRootNoMarker,
			expectedRoot:     currentWorkingDirectory,
			expectedLog:      noRootMarkerLog,
			rootMarkerExists: false,
		},
		{
			root:             systemRoot,
			name:             "Defined absolute path root",
			module:           syslutil.MustAbsolute(t, definedRootNoMarker.files[0]),
			structure:        definedRootNoMarker,
			expectedRoot:     systemRoot,
			expectedLog:      noRootMarkerLog,
			rootMarkerExists: false,
		},
		{
			name:             "Defined relative root with absolute module path rooted at root",
			root:             currentWorkingDirectory,
			module:           filepath.Join(systemRoot, filepath.Clean(definedRootNoMarker.files[0])),
			structure:        definedRootNoMarker,
			expectedRoot:     currentWorkingDirectory,
			expectedLog:      noRootMarkerLog,
			rootMarkerExists: false,
		},
		{
			name:             "Defined root flag and root",
			root:             currentWorkingDirectory,
			module:           syslutil.MustRelative(t, currentWorkingDirectory, definedRootFlagAndMarkerFound.files[0]),
			structure:        definedRootFlagAndMarkerFound,
			expectedRoot:     currentWorkingDirectory,
			expectedLog:      rootMarkerExistsLog,
			rootMarkerPath:   syslutil.MustAbsolute(t, "./DefinedRootAndSyslRootDefinedTest/path/"),
			rootMarkerExists: false,
		},
		{
			name:             "Defined root flag and root marker with absolute path module rooted at root",
			root:             "./DefinedRootAndSyslRootDefinedTest/",
			module:           "/path/to/module/test.sysl",
			structure:        definedRootFlagAndMarkerFound,
			expectedRoot:     "./DefinedRootAndSyslRootDefinedTest/",
			expectedLog:      rootMarkerExistsLog,
			rootMarkerPath:   syslutil.MustAbsolute(t, "./DefinedRootAndSyslRootDefinedTest/path/"),
			rootMarkerExists: false,
		},
		{
			name:             "Root is not defined",
			root:             "",
			module:           undefinedRoot.files[0],
			structure:        undefinedRoot,
			expectedRoot:     filepath.Dir(undefinedRoot.files[0]),
			expectedLog:      noRootLog,
			rootMarkerExists: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(ts folderTestStructure) func(t *testing.T) {
			return func(tt *testing.T) {
				tt.Parallel()

				logger, hook := log.NewNullLogger()
				fs := afero.NewMemMapFs()
				syslutil.BuildFolderTest(tt, fs, ts.structure.folders, ts.structure.files)

				r := &cmdRunner{Root: ts.root, module: ts.module}

				require.NoError(tt, r.getProjectRoot(fs, logger))
				require.Equal(tt, ts.expectedRoot, r.Root)
				require.Equal(tt, ts.getExpectedModule(tt), r.module)

				if !ts.rootMarkerExists {
					require.Equal(tt, 1, len(hook.Entries))
					require.Equal(tt, logrus.WarnLevel, hook.LastEntry().Level)
					require.Equal(tt, ts.getExpectedLog(), hook.LastEntry().Message)
				} else {
					require.Equal(tt, 0, len(hook.Entries))
				}
			}
		}(test))
	}
}

func (ts folderTestStructure) getExpectedModule(t *testing.T) string {
	// if root is defined, expected root and root param is the same and module is not changed
	if ts.expectedRoot == ts.root {
		return ts.module
	}
	return syslutil.MustRelative(t, ts.expectedRoot, ts.module)
}

func (ts folderTestStructure) getExpectedLog() string {
	switch ts.expectedLog {
	case noRootMarkerLog:
		return fmt.Sprintf(noRootMarkerWarning, syslRootMarker, ts.expectedRoot)
	case rootMarkerExistsLog:
		return fmt.Sprintf(rootMarkerExistsWarning, syslRootMarker, ts.rootMarkerPath, ts.expectedRoot)
	case noRootLog:
		return fmt.Sprintf(noRootWarning, syslRootMarker, ts.expectedRoot)
	}
	return ""
}
