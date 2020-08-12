package mod

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

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
		return nil, fmt.Errorf("%s not found: no such file in current working directory", name)
	}

	mod, relpath, err := fs.fetchRemoteFile(name)
	if err != nil {
		return nil, err
	}

	return syslutil.NewChrootFs(afero.NewOsFs(), mod.Dir).Open(relpath)
}

func (fs *Fs) OpenWithModule(name string) (afero.File, *Module, error) {
	f, err := fs.Fs.Open(name)
	if err == nil {
		return f, nil
	} else if !SyslModules {
		return nil, fmt.Errorf("%s not found: no such file in current working directory", name)
	}

	mod, relpath, err := fs.fetchRemoteFile(name)
	if err != nil {
		return nil, err
	}

	f, err := syslutil.NewChrootFs(afero.NewOsFs(), mod.Dir).Open(relpath)
	return f, mod, err
}

func (fs *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.Fs.OpenFile(name, flag, perm)
	if err == nil {
		return f, nil
	} else if !SyslModules {
		return nil, fmt.Errorf("%s not found: no such file in current working directory", name)
	}

	mod, relpath, err := fs.fetchRemoteFile(name)
	if err != nil {
		return nil, err
	}

	return syslutil.NewChrootFs(afero.NewOsFs(), mod.Dir).OpenFile(relpath, flag, perm)
}

func (fs *Fs) Name() string {
	return "ModSupportedFs"
}

func (fs *Fs) fetchRemoteFile(filename string) (*Module, string, error) {
	name, ver, err := mergeRootAndPath(fs.root, filename)
	if err != nil {
		return nil, "", err
	}

	mod, err := Retrieve(name, ver)
	if err != nil {
		return nil, "", err
	}

	relpath, err := filepath.Rel(mod.Name, name)
	if err != nil {
		return nil, "", err
	}

	return mod, relpath, nil
}

func mergeRootAndPath(root, name string) (string, string, error) {
	// path.Join will strip path elements of ".", so if the root is "."
	// it will still work as a go module path when prepended with "."
	root, ver := ExtractVersion(root)
	var nameVer string
	name, nameVer = ExtractVersion(name)

	if ver != "" && nameVer != "" && ver != nameVer {
		return "", "", fmt.Errorf("root version %s does not equal to path version %s", ver, nameVer)
	}

	if nameVer != "" {
		ver = nameVer
	}

	name = path.Join(root, name)
	return name, ver, nil
}

func ExtractVersion(path string) (newpath, ver string) {
	newpath = path
	s := strings.Split(path, "@")
	if len(s) > 1 {
		ver = s[len(s)-1]
		newpath = path[:len(path)-len(ver)-1]
	}
	return
}

func AppendVersion(path, ver string) string {
	if ver == "" {
		return path
	}
	return path + "@" + ver
}
