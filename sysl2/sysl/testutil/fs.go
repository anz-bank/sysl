package testutil

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertFsHasExactly asserts that fs contains the given files and only
// those. All paths must start with '/'.
func AssertFsHasExactly(t *testing.T, fs afero.Fs, paths ...string) bool {
	expected := make([]string, 0, len(paths))
	for _, p := range paths {
		expected = append(expected, filepath.Clean(p))
	}
	sort.Strings(expected)

	actual := []string{}
	require.NoError(t, afero.Walk(fs, "/", func(path string, info os.FileInfo, err error) error {
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

func WriteToMemOverlayFs(t *testing.T, osRoot string) (memFs, fs afero.Fs) {
	memFs = CreateTestChrootFs(t, afero.NewMemMapFs(), "/")
	fs = afero.NewCopyOnWriteFs(CreateTestChrootFs(t, afero.NewOsFs(), osRoot), memFs)
	return
}

func CreateTestChrootFs(t *testing.T, fs afero.Fs, root string) (chrootfs afero.Fs) {
	var err error
	chrootfs, err = syslutil.NewChrootFs(fs, root)
	require.NoError(t, err)
	return chrootfs
}
