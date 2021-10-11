package syslutil

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sep = string(os.PathSeparator)

// AssertFsHasExactly asserts that fs contains the given files and only
// those. All paths must start with '/'.
func AssertFsHasExactly(t *testing.T, fs afero.Fs, paths ...string) bool {
	expected := make([]string, 0, len(paths))
	for _, p := range paths {
		expected = append(expected, MustAbsolute(t, p))
	}
	sort.Strings(expected)

	actual := []string{}

	root := MustAbsolute(t, sep)

	// special case for MemFs in windows as walking from root results in infinite loop
	_, isMemFs := fs.(*afero.MemMapFs)
	isMemFsWindows := runtime.GOOS == windows && isMemFs
	if isMemFsWindows {
		root = sep
	}

	require.NoError(t, afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			actual = append(actual, path)
		}
		return nil
	}))
	sort.Strings(actual)

	// add back the volume name as expected always have absolute path
	if isMemFsWindows {
		for i := 0; i < len(actual); i++ {
			actual[i] = MustAbsolute(t, filepath.Join(sep, actual[i]))
		}
	}

	return assert.Equal(t, expected, actual)
}

func WriteToMemOverlayFs(osRoot string) (memFs, fs afero.Fs) {
	memFs = NewChrootFs(afero.NewMemMapFs(), "/")
	fs = afero.NewCopyOnWriteFs(NewChrootFs(afero.NewOsFs(), osRoot), memFs)
	return
}

func MustAbsolute(t *testing.T, path string) string {
	fullPath, err := filepath.Abs(path)
	require.NoError(t, err)
	return fullPath
}

func MustRelative(t *testing.T, base, target string) string {
	base = MustAbsolute(t, base)
	target = MustAbsolute(t, target)
	rel, err := filepath.Rel(base, target)
	require.NoError(t, err)
	return rel
}

func BuildFolderTest(t *testing.T, fs afero.Fs, folders, files []string) {
	for _, folder := range folders {
		require.NoError(t, fs.MkdirAll(trimVolumeName(MustAbsolute(t, folder)), os.ModeTemporary))
	}

	for _, file := range files {
		_, err := fs.Create(trimVolumeName(MustAbsolute(t, file)))
		require.NoError(t, err)
	}
}

func HandleCRLF(text []byte) []byte {
	if runtime.GOOS == windows {
		re := regexp.MustCompile("(\r\n|\r)")
		text = re.ReplaceAll(text, []byte("\n"))
	}
	return text
}
