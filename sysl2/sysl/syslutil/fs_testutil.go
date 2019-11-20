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

// AssertFsHasExactly asserts that fs contains the given files and only
// those. All paths must start with '/'.
func AssertFsHasExactly(t *testing.T, fs afero.Fs, paths ...string) bool {
	expected := make([]string, 0, len(paths))
	for _, p := range paths {
		expected = append(expected, MustAbsolute(t, p))
	}
	sort.Strings(expected)

	actual := []string{}
	root := MustAbsolute(t, string(os.PathSeparator))
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
}

func HandleCRLF(text []byte) []byte {
	if runtime.GOOS == "windows" {
		re := regexp.MustCompile("(\r\n|\r)")
		text = re.ReplaceAll(text, []byte("\n"))
	}
	return text
}
