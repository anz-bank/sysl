package main

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	syslRootMarker = ".sysl"
)

type cmdRunner struct {
	commands map[string]Command

	Root   string
	module string
}

func (r *cmdRunner) Run(which string, fs afero.Fs, logger *logrus.Logger) error {
	if cmd, ok := r.commands[which]; ok {
		if cmd.Name() == which {
			var mod *sysl.Module
			var err error
			if cmd.RequireSyslModule() {
				if err = r.setProjectRoot(fs, logger); err != nil {
					return err
				}

				mod, _, err = LoadSyslModule(r.Root, r.module, fs, logger)
				if err != nil {
					return err
				}
			}

			return cmd.Execute(ExecuteArgs{mod, fs, logger})
		}
	}
	return nil
}

func (r *cmdRunner) Configure(app *kingpin.Application) error {
	commands := []Command{
		&protobuf{},
		&intsCmd{},
		&datamodelCmd{},
		&codegenCmd{},
		&sequenceDiagramCmd{},
		&importCmd{},
		&infoCmd{},
		&validateCmd{},
		&exportCmd{},
	}
	r.commands = map[string]Command{}

	app.Flag("root",
		"sysl root directory for input model file. If root is not found, the module directory becomes the root").
		Default("").StringVar(&r.Root)

	sort.Slice(commands, func(i, j int) bool {
		return strings.Compare(commands[i].Name(), commands[j].Name()) < 0
	})
	for _, cmd := range commands {
		c := cmd.Configure(app)
		if cmd.RequireSyslModule() {
			c.Arg("MODULE", "input files without .sysl extension and with leading /, eg: "+
				"/project_dir/my_models combine with --root if needed").
				Required().StringVar(&r.module)
		}
		r.commands[cmd.Name()] = cmd
	}

	return nil
}

// Helper function to validate that a set of command flags are not empty values
func EnsureFlagsNonEmpty(cmd *kingpin.CmdClause, excludes ...string) {
	inExcludes := func(s string) bool {
		for _, e := range excludes {
			if s == e {
				return true
			}
		}
		return false
	}
	fn := func(c *kingpin.ParseContext) error {
		var errorMsg strings.Builder
		for _, f := range cmd.Model().Flags {
			if inExcludes(f.Name) {
				continue
			}
			val := f.Value.String()

			if val != "" {
				val = strings.Trim(val, " ")
				if val == "" {
					errorMsg.WriteString("'" + f.Name + "'" + " value passed is empty\n")
				}
			} else if len(f.Default) > 0 {
				errorMsg.WriteString("'" + f.Name + "'" + " value passed is empty\n")
			}
		}
		if errorMsg.Len() > 0 {
			return errors.New(errorMsg.String())
		}
		return nil
	}

	cmd.PreAction(fn)
}

func (r *cmdRunner) setProjectRoot(fs afero.Fs, logger *logrus.Logger) error {
	rootIsDefined := r.Root != ""

	modulePath := r.module
	if rootIsDefined {
		modulePath = filepath.Join(r.Root, r.module)
	}

	syslRootPath, err := findRootFromSyslModule(modulePath, fs)
	if err != nil {
		return err
	}

	rootMarkerExists := syslRootPath != ""

	if rootIsDefined {
		if rootMarkerExists {
			logger.Warningf("%s found in %s but will use %s instead",
				syslRootMarker, syslRootPath, r.Root)
		} else {
			logger.Warningf("%s is not defined but root flag is defined in %s",
				syslRootMarker, r.Root)
		}
	} else {
		if rootMarkerExists {
			r.Root = syslRootPath

			// module has to be relative to the root
			absModulePath, err := filepath.Abs(r.module)
			if err != nil {
				return err
			}
			r.module, err = filepath.Rel(r.Root, absModulePath)
			if err != nil {
				return err
			}
		} else {
			// uses the module directory as the root, changing the module to be relative to the root
			r.Root = filepath.Dir(r.module)
			r.module = filepath.Base(r.module)
			logger.Warningf("root and %s are undefined, %s will be used instead",
				syslRootMarker, r.Root)
		}
	}

	return nil
}

func findRootFromSyslModule(modulePath string, fs afero.Fs) (string, error) {
	currentPath, err := filepath.Abs(modulePath)
	if err != nil {
		return "", err
	}

	systemRoot, err := filepath.Abs(string(os.PathSeparator))
	if err != nil {
		return "", err
	}

	// Keep walking up the directories to find nearest root marker
	for {
		currentPath = filepath.Dir(currentPath)
		exists, err := afero.Exists(fs, filepath.Join(currentPath, syslRootMarker))
		reachedRoot := currentPath == systemRoot || (err != nil && os.IsPermission(err))
		switch {
		case exists:
			return currentPath, nil
		case reachedRoot:
			return "", nil
		case err != nil:
			return "", err
		}
	}
}
