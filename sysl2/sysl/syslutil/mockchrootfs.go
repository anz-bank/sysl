package syslutil

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockFs struct {
	fs afero.Fs
	mock.Mock
}

var _ afero.Fs = &MockFs{}

// NewMockFs creates a wrapper fs to record function calls with the
// mock library, it requires a memory map fs
func NewMockFs(t *testing.T, memFs afero.Fs) *MockFs {
	require.IsType(t, &afero.MemMapFs{}, memFs)
	return &MockFs{fs: memFs}
}

func (m *MockFs) Create(name string) (afero.File, error) {
	m.Called(name)
	return m.fs.Create(name)
}

func (m *MockFs) Mkdir(name string, perm os.FileMode) error {
	m.Called(name, perm)
	return m.fs.Mkdir(name, perm)
}

func (m *MockFs) MkdirAll(name string, perm os.FileMode) error {
	m.Called(name, perm)
	return m.fs.MkdirAll(name, perm)
}

func (m *MockFs) Open(name string) (afero.File, error) {
	m.Called(name)
	return m.fs.Open(name)
}

func (m *MockFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	m.Called(name, flag, perm)
	return m.fs.OpenFile(name, flag, perm)
}

func (m *MockFs) Remove(name string) error {
	m.Called(name)
	return m.fs.Remove(name)
}

func (m *MockFs) RemoveAll(path string) error {
	m.Called(path)
	return m.fs.RemoveAll(path)
}

func (m *MockFs) Rename(oldname, newname string) error {
	m.Called(oldname, newname)
	return m.fs.Rename(oldname, newname)
}

func (m *MockFs) Stat(name string) (os.FileInfo, error) {
	m.Called(name)
	return m.fs.Stat(name)
}

func (m *MockFs) Name() string {
	m.Called()
	return "MockFs"
}

func (m *MockFs) Chmod(name string, mode os.FileMode) error {
	m.Called(name, mode)
	return m.fs.Chmod(name, mode)
}

func (m *MockFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	m.Called(name, atime, mtime)
	return m.fs.Chtimes(name, atime, mtime)
}
