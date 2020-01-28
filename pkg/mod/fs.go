package mod

import (
	"os"
	"path/filepath"
	"time"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

type Fs struct {
	source afero.Fs
}

func NewFs(fs afero.Fs) *Fs {
	return &Fs{source: fs}
}

func (fs *Fs) Open(name string) (afero.File, error) {
	f, err := fs.source.Open(name)
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

	return afero.NewOsFs().Open(filepath.Join(mod.Dir, relpath))
}

func (fs *Fs) Create(name string) (afero.File, error) {
	return fs.source.Create(name)
}

func (fs *Fs) Mkdir(name string, perm os.FileMode) error {
	return fs.source.Mkdir(name, perm)
}

func (fs *Fs) MkdirAll(path string, perm os.FileMode) error {
	return fs.source.MkdirAll(path, perm)
}

func (fs *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.source.OpenFile(name, flag, perm)
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

	return syslutil.NewChrootFs(afero.NewOsFs(), mod.Dir).OpenFile(relpath, flag, perm)
}

func (fs *Fs) Remove(name string) error {
	return fs.source.Remove(name)
}

func (fs *Fs) RemoveAll(path string) error {
	return fs.source.RemoveAll(path)
}

func (fs *Fs) Rename(oldname, newname string) error {
	return fs.source.Rename(oldname, newname)
}

func (fs *Fs) Stat(name string) (os.FileInfo, error) {
	return fs.source.Stat(name)
}

func (fs *Fs) Name() string {
	return "ModSupportedFs"
}

func (fs *Fs) Chmod(name string, mode os.FileMode) error {
	return fs.source.Chmod(name, mode)
}

func (fs *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.source.Chtimes(name, atime, mtime)
}
