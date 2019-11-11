package testutil

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/nofun97/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertFsHasExactly asserts that fs contains the given files and only
// those. All paths must start with '/'.
func AssertFsHasExactly(t *testing.T, fs afero.Fs, paths ...string) bool {
	expected := make([]string, 0, len(paths))
	for _, p := range paths {
		expected = append(expected, GetAbsolute(t, p))
	}
	sort.Strings(expected)

	actual := []string{}
	root, err := filepath.Abs(string(os.PathSeparator))
	require.NoError(t, err)
	require.NoError(t, afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
		t.Log("Walking: ", path)
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
	memFs = syslutil.NewChrootFs(afero.NewMemMapFs(), "/")
	fs = afero.NewCopyOnWriteFs(syslutil.NewChrootFs(afero.NewOsFs(), osRoot), memFs)
	return
}

func GetAbsolute(t *testing.T, path string) string {
	fullPath, err := filepath.Abs(path)
	require.NoError(t, err)
	return fullPath
}
