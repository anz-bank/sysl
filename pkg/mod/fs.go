package mod

import (
	"os"
	"path/filepath"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/spf13/afero"
)

type Fs struct {
	afero.Fs
	root string
}

func NewFs(fs afero.Fs, root string) *Fs {
	return &Fs{Fs: syslutil.NewChrootFs(fs, root), root: root}
}

func (fs *Fs) Open(name string) (afero.File, error) {
	f, err := fs.Fs.Open(name)
	if err == nil {
		return f, nil
	} else if !SyslModules {
		return nil, err
	}

	// filepath.Join will strip path elements of ".", so if the root is "."
	// it will still work as a go module path when prepended with "."
	name = filepath.Join(fs.root, name)
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

func (fs *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.Fs.OpenFile(name, flag, perm)
	if err == nil {
		return f, nil
	} else if !SyslModules {
		return nil, err
	}

	// filepath.Join will strip path elements of ".", so if the root is "."
	// it will still work as a go module path when prepended with "."
	name = filepath.Join(fs.root, name)
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

func (fs *Fs) Name() string {
	return "ModSupportedFs"
}
