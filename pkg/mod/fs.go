package mod

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"sync"

	"github.com/anz-bank/pkg/mod"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

var SyslModules = os.Getenv("SYSL_MODULES") != SyslModulesOff
var ModMode = mod.ModeType(os.Getenv("SYSL_MODULES"))

const (
	SyslModulesOff = "off"
	SyslModulesOn  = "on"
)

var githubAccessToken = os.Getenv("SYSL_GITHUB_TOKEN")

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

	m, relpath, err := fs.fetchRemoteFile(name)
	if err != nil {
		return nil, err
	}

	return syslutil.NewChrootFs(afero.NewOsFs(), m.Dir).Open(relpath)
}

func (fs *Fs) OpenWithModule(name string) (afero.File, *mod.Module, error) {
	f, err := fs.Fs.Open(name)
	if err == nil {
		return f, nil, nil
	} else if !SyslModules {
		return nil, nil, fmt.Errorf("%s not found: no such file in current working directory", name)
	}

	m, relpath, err := fs.fetchRemoteFile(name)
	if err != nil {
		return nil, nil, err
	}

	f, err = syslutil.NewChrootFs(afero.NewOsFs(), m.Dir).Open(relpath)
	return f, m, err
}

func (fs *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	f, err := fs.Fs.OpenFile(name, flag, perm)
	if err == nil {
		return f, nil
	} else if !SyslModules {
		return nil, fmt.Errorf("%s not found: no such file in current working directory", name)
	}

	m, relpath, err := fs.fetchRemoteFile(name)
	if err != nil {
		return nil, err
	}

	return syslutil.NewChrootFs(afero.NewOsFs(), m.Dir).OpenFile(relpath, flag, perm)
}

func (fs *Fs) Name() string {
	return "ModSupportedFs"
}

var configModOnce sync.Once

func (fs *Fs) fetchRemoteFile(filename string) (*mod.Module, string, error) {
	name, ver, err := mergeRootAndPath(fs.root, filename)
	if err != nil {
		return nil, "", err
	}

	config(fs.root)

	m, err := mod.Retrieve(name, ver)
	if err != nil {
		switch err.(type) {
		case *mod.RateLimitError:
			return nil, "", errors.Wrap(err,
				"\033[1;36mplease set up envvar SYSL_GITHUB_TOKEN\033[0m")
		case *mod.NotFoundError:
			return nil, "", errors.Wrap(err,
				"\033[1;36mplease check whether envvar SYSL_GITHUB_TOKEN is set if it is a private repository\033[0m")
		}
		return nil, "", err
	}

	relpath, err := filepath.Rel(m.Name, name)
	if err != nil {
		return nil, "", err
	}

	return m, relpath, nil
}

func config(root string) {
	configModOnce.Do(func() {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		cacheDir := filepath.Join(usr.HomeDir, ".sysl")
		if ModMode == "" {
			ModMode = mod.GoModulesMode
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		if err := mod.Config(ModMode,
			mod.GoModulesOptions{Root: root, ModName: filepath.Base(wd)},
			mod.GitHubOptions{CacheDir: cacheDir, AccessToken: githubAccessToken}); err != nil {
			log.Fatal(err)
		}
	})
}

func mergeRootAndPath(root, name string) (string, string, error) {
	// path.Join will strip path elements of ".", so if the root is "."
	// it will still work as a go module path when prepended with "."
	root, ver := mod.ExtractVersion(root)
	var nameVer string
	name, nameVer = mod.ExtractVersion(name)

	if ver != "" && nameVer != "" && ver != nameVer {
		return "", "", fmt.Errorf("root version %s does not equal to path version %s", ver, nameVer)
	}

	if nameVer != "" {
		ver = nameVer
	}

	name = path.Join(root, name)
	return name, ver, nil
}
