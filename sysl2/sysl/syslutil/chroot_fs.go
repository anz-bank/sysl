package syslutil

import (
	"errors"
	"log"
	"os"
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

// NewChrootFs returns a filesystem that is rooted at root argument
func NewChrootFs(fs afero.Fs, root string) (*ChrootFs, error) {
	var err error
	if !filepath.IsAbs(root) {
		root, err = filepath.Abs(root)

		if err != nil {
			return nil, err
		}
	}
	return &ChrootFs{fs: fs, root: root}, nil
}

func (fs *ChrootFs) join(name string) string {
	volumeName := filepath.VolumeName(name)
	cleaned := filepath.ToSlash(name)
	cleaned = strings.TrimLeft(cleaned, volumeName)
	joined := filepath.Join(fs.root, cleaned)
	log.Println("Cleaned: ", cleaned, "Joined: ", joined)
	return joined
}

func (fs *ChrootFs) Create(name string) (afero.File, error) {
	data, err := fs.wrapCallWithData(name,
		func(fixedPath string) (interface{}, error) {
			return fs.fs.Create(fixedPath)
		})
	if err != nil {
		return nil, err
	}
	return data.(afero.File), err
}

func (fs *ChrootFs) Mkdir(name string, perm os.FileMode) error {
	return fs.wrapCall(name, func(fixedPath string) error {
		return fs.fs.Mkdir(fixedPath, perm)
	})
}

func (fs *ChrootFs) MkdirAll(path string, perm os.FileMode) error {
	return fs.wrapCall(path, func(fixedPath string) error {
		return fs.fs.MkdirAll(fixedPath, perm)
	})
}

func (fs *ChrootFs) Open(name string) (afero.File, error) {
	data, err := fs.wrapCallWithData(name,
		func(fixedPath string) (interface{}, error) {
			return fs.fs.Open(fixedPath)
		})
	if err != nil {
		return nil, err
	}
	return data.(afero.File), err
}

func (fs *ChrootFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	data, err := fs.wrapCallWithData(name,
		func(fixedPath string) (interface{}, error) {
			return fs.fs.OpenFile(fixedPath, flag, perm)
		})
	if err != nil {
		return nil, err
	}
	return data.(afero.File), err
}

func (fs *ChrootFs) Remove(name string) error {
	return fs.wrapCall(name, func(fixedPath string) error {
		return fs.fs.Remove(fixedPath)
	})
}

func (fs *ChrootFs) RemoveAll(path string) error {
	return fs.wrapCall(path, func(fixedPath string) error {
		return fs.fs.RemoveAll(fixedPath)
	})
}

func (fs *ChrootFs) Rename(oldname, newname string) error {
	return fs.wrapCall(oldname, func(fixedPath string) error {
		return fs.fs.Rename(fixedPath, fs.join(newname))
	})
}

func (fs *ChrootFs) Stat(name string) (os.FileInfo, error) {
	data, err := fs.wrapCallWithData(name,
		func(fixedPath string) (interface{}, error) {
			return fs.fs.Stat(fixedPath)
		})
	if err != nil {
		return nil, err
	}
	return data.(os.FileInfo), err
}

func (fs *ChrootFs) Name() string {
	return "ChrootFS"
}

func (fs *ChrootFs) Chmod(name string, mode os.FileMode) error {
	return fs.wrapCall(name, func(fixedPath string) error {
		return fs.fs.Chmod(fs.join(fixedPath), mode)
	})
}

func (fs *ChrootFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.wrapCall(name, func(fixedPath string) error {
		return fs.fs.Chtimes(fixedPath, atime, mtime)
	})
}

func (fs *ChrootFs) openAllowed(name string) error {
	absoluteRoot, err := filepath.Abs(fs.root)
	if err != nil {
		return err
	}

	relativePath, err := filepath.Rel(absoluteRoot, name)
	if err != nil {
		return err
	}

	if strings.Split(relativePath, string(os.PathSeparator))[0] == ".." {
		return errors.New("permission denied, file outside root")
	}

	return err
}

func (fs *ChrootFs) wrapCall(path string, fn func(string) error) error {
	filename := fs.join(path)
	if err := fs.openAllowed(filename); err != nil {
		return err
	}
	return fn(filename)
}

func (fs *ChrootFs) wrapCallWithData(
	path string,
	fn func(string) (interface{}, error)) (interface{}, error) {
	filename := fs.join(path)
	if err := fs.openAllowed(filename); err != nil {
		return nil, err
	}
	return fn(filename)
}
