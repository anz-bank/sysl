package syslutil

import (
	"os"
	"path/filepath"
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
		expected = append(expected, getAbsolute(t, p))
	}
	sort.Strings(expected)

	actual := []string{}
	root, err := filepath.Abs(string(os.PathSeparator))
	require.NoError(t, err)
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

func getAbsolute(t *testing.T, path string) string {
	fullPath, err := filepath.Abs(path)
	require.NoError(t, err)
	return fullPath
}
