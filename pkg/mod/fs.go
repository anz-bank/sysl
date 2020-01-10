package mod

import (
	"os"
	"path/filepath"
	"time"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

type Fs struct {
	fs afero.Fs
}

func NewFs(fs afero.Fs) *Fs {
	return &Fs{fs: fs}
}

func (fs *Fs) Open(name string) (afero.File, error) {
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

func (fs *Fs) Create(name string) (afero.File, error) {
	return fs.fs.Create(name)
}

func (fs *Fs) Mkdir(name string, perm os.FileMode) error {
	return fs.fs.Mkdir(name, perm)
}

func (fs *Fs) MkdirAll(path string, perm os.FileMode) error {
	return fs.fs.MkdirAll(path, perm)
}

func (fs *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return fs.fs.OpenFile(name, flag, perm)
}

func (fs *Fs) Remove(name string) error {
	return fs.fs.Remove(name)
}

func (fs *Fs) RemoveAll(path string) error {
	return fs.fs.RemoveAll(path)
}

func (fs *Fs) Rename(oldname, newname string) error {
	return fs.fs.Rename(oldname, newname)
}

func (fs *Fs) Stat(name string) (os.FileInfo, error) {
	return fs.fs.Stat(name)
}

func (fs *Fs) Name() string {
	return "Fs"
}

func (fs *Fs) Chmod(name string, mode os.FileMode) error {
	return fs.fs.Chmod(name, mode)
}

func (fs *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.fs.Chtimes(name, atime, mtime)
}
