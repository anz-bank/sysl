package syslutil

import (
	"log"
	"os"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
)

type MockFs struct {
	mock.Mock
}

var _ afero.Fs = &MockFs{}

// NewMockFs creates a wrapper fs to record function calls with the
// mock library, it requires a memory map fs
func NewMockChrootFs(root string) (*MockFs, *ChrootFs) {
	mockFs := &MockFs{}
	return mockFs, NewChrootFs(mockFs, root)
}

func (m *MockFs) Create(name string) (afero.File, error) {
	res := m.Called(name)
	return res.Get(0).(afero.File), res.Error(1)
}

func (m *MockFs) Mkdir(name string, perm os.FileMode) error {
	return m.Called(name, perm).Error(0)
}

func (m *MockFs) MkdirAll(name string, perm os.FileMode) error {
	return m.Called(name, perm).Error(0)
}

func (m *MockFs) Open(name string) (afero.File, error) {
	res := m.Called(name)
	log.Println(res)
	return res.Get(0).(afero.File), res.Error(1)
}

func (m *MockFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	res := m.Called(name, flag, perm)
	return res.Get(0).(afero.File), res.Error(1)
}

func (m *MockFs) Remove(name string) error {
	return m.Called(name).Error(0)
}

func (m *MockFs) RemoveAll(path string) error {
	return m.Called(path).Error(0)
}

func (m *MockFs) Rename(oldname, newname string) error {
	return m.Called(oldname, newname).Error(0)
}

func (m *MockFs) Stat(name string) (os.FileInfo, error) {
	res := m.Called(name)
	return res.Get(0).(os.FileInfo), res.Error(1)
}

func (m *MockFs) Name() string {
	return m.Called().String(0)
}

func (m *MockFs) Chmod(name string, mode os.FileMode) error {
	return m.Called(name, mode).Error(0)
}

func (m *MockFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return m.Called(name, atime, mtime).Error(0)
}
