package syslutil

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

type ChrootFs struct {
	fs   afero.Fs
	root string
}

var _ afero.Fs = &ChrootFs{}

func NewChrootFs(fs afero.Fs, root string) *ChrootFs {
	return &ChrootFs{fs: fs, root: root}
}

func (fs *ChrootFs) join(name string) string {
	return path.Join(fs.root, name)
}

func (fs *ChrootFs) Create(name string) (afero.File, error) {
	return fs.fs.Create(fs.join(name))
}

func (fs *ChrootFs) Mkdir(name string, perm os.FileMode) error {
	return fs.fs.Mkdir(fs.join(name), perm)
}

func (fs *ChrootFs) MkdirAll(path string, perm os.FileMode) error {
	return fs.fs.MkdirAll(fs.join(path), perm)
}

func (fs *ChrootFs) Open(name string) (afero.File, error) {
	if err := fs.openAllowed(name); err != nil {
		return nil, err
	}
	return fs.fs.Open(fs.join(name))
}

func (fs *ChrootFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	if err := fs.openAllowed(name); err != nil {
		return nil, err
	}
	return fs.fs.OpenFile(fs.join(name), flag, perm)
}

func (fs *ChrootFs) Remove(name string) error {
	return fs.fs.Remove(fs.join(name))
}

func (fs *ChrootFs) RemoveAll(path string) error {
	return fs.fs.RemoveAll(fs.join(path))
}

func (fs *ChrootFs) Rename(oldname, newname string) error {
	return fs.fs.Rename(fs.join(oldname), fs.join(newname))
}

func (fs *ChrootFs) Stat(name string) (os.FileInfo, error) {
	return fs.fs.Stat(fs.join(name))
}

func (fs *ChrootFs) Name() string {
	return "ChrootFS"
}

func (fs *ChrootFs) Chmod(name string, mode os.FileMode) error {
	return fs.fs.Chmod(fs.join(name), mode)
}

func (fs *ChrootFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.fs.Chtimes(fs.join(name), atime, mtime)
}

func (fs *ChrootFs) openAllowed(name string) (err error) {
	modulePath := filepath.Join(fs.root, name)
	if !strings.HasPrefix(modulePath, fs.root) {
		err = errors.New("import can not go up the directory")
	}
	return
}
