package mod

import (
	"os"
	"path/filepath"
	"time"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

type ModSupportedFs struct {
	fs afero.Fs
}

func NewModSupportedFs(fs afero.Fs) *ModSupportedFs {
	return &ModSupportedFs{fs: fs}
}

func (fs *ModSupportedFs) Open(name string) (afero.File, error) {
	f, err := fs.fs.Open(name)
	if err == nil {
		return f, nil
	}

	mod, err := Find(name)
	if err != nil {
		return nil, err
	}
	relpath, err := filepath.Rel(mod.Name, name)
	if err != nil {
		return nil, err
	}

	return syslutil.NewChrootFs(afero.NewOsFs(), mod.Dir).Open(relpath)
}

func (fs *ModSupportedFs) Create(name string) (afero.File, error) {
	return fs.fs.Create(name)
}

func (fs *ModSupportedFs) Mkdir(name string, perm os.FileMode) error {
	return fs.fs.Mkdir(name, perm)
}

func (fs *ModSupportedFs) MkdirAll(path string, perm os.FileMode) error {
	return fs.fs.MkdirAll(path, perm)
}

func (fs *ModSupportedFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return fs.fs.OpenFile(name, flag, perm)
}

func (fs *ModSupportedFs) Remove(name string) error {
	return fs.fs.Remove(name)
}

func (fs *ModSupportedFs) RemoveAll(path string) error {
	return fs.fs.RemoveAll(path)
}

func (fs *ModSupportedFs) Rename(oldname, newname string) error {
	return fs.fs.Rename(oldname, newname)
}

func (fs *ModSupportedFs) Stat(name string) (os.FileInfo, error) {
	return fs.fs.Stat(name)
}

func (fs *ModSupportedFs) Name() string {
	return "ModSupportedFs"
}

func (fs *ModSupportedFs) Chmod(name string, mode os.FileMode) error {
	return fs.fs.Chmod(name, mode)
}

func (fs *ModSupportedFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.fs.Chtimes(name, atime, mtime)
}
