package mod

import (
	"path/filepath"
	"strings"
	"time"
)

type goModule struct {
	Path     string     // module path
	Version  string     // module version
	Versions []string   // available module versions (with -versions)
	Replace  *goModule  // replaced by this module
	Time     *time.Time // time version was created
	Update   *goModule  // available update, if any (with -u)
	Main     bool       // is this the main module?
	Indirect bool       // is this module only an indirect dependency of main module?
	Dir      string     // directory holding files for this module, if any
	GoMod    string     // path to go.mod file for this module, if any
}

type goModules []*goModule

var GoMods goModules

func (modules goModules) GetByFilepath(p string) *goModule {
	if modules == nil {
		return nil
	}

	p = filepath.Clean(p)

	for _, m := range modules {
		if strings.HasPrefix(p, filepath.Clean(m.Path)) {
			return m
		}
	}

	return nil
}
