package main

import (
	"errors"
	"path/filepath"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	syslRootMarker = ".SYSL_ROOT"
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
				err = r.rootHandler(fs, logger)
				if err != nil {
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

func (r *cmdRunner) rootHandler(fs afero.Fs, logger *logrus.Logger) error {
	moduleAbsolutePath, err := filepath.Abs(r.module)
	if err != nil {
		return err
	}

	syslRootPath, err := r.findRootFromASyslModule(moduleAbsolutePath, fs)
	if err != nil {
		return err
	}

	rootIsDefined := r.Root != ""
	syslRootExists := syslRootPath != ""

	if rootIsDefined && syslRootExists {
		logger.WithFields(logrus.Fields{
			"root":      r.Root,
			"SYSL_ROOT": syslRootPath,
		}).Warningf("root is defined even though %s exists\n", syslRootMarker)
	} else if rootIsDefined && !syslRootExists {
		logger.Warningf("%s is not defined but root flag is defined in %s and will be used", syslRootMarker, r.Root)
	} else if !rootIsDefined && !syslRootExists {
		logger.Errorf("root and %s are undefined", syslRootMarker)
		return errors.New("root is not defined")
	}

	// r.Root must be relative to the current directory
	absoluteRoot, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	if !rootIsDefined && syslRootExists {
		relativeSyslRootPath, err := filepath.Rel(absoluteRoot, syslRootPath)
		if err != nil {
			return err
		}
		r.Root = relativeSyslRootPath
	}

	// the absolute path of the currently used root directory
	realAbsolutePath, err := filepath.Abs(r.Root)
	if err != nil {
		return err
	}

	// module path is relative to the currently used root path
	relativePathModule, err := filepath.Rel(realAbsolutePath, moduleAbsolutePath)
	if err != nil {
		return err
	}
	r.module = relativePathModule

	return nil
}

func (r *cmdRunner) findRootFromASyslModule(absolutePath string, fs afero.Fs) (string, error) {
	// Takes the closest root marker
	currentPath := absolutePath

	for {
		// Keep walking up the directories
		currentPath = filepath.Dir(currentPath)

		rootPath := filepath.Join(currentPath, syslRootMarker)

		if exists, err := afero.Exists(fs, rootPath); err != nil {
			return "", err
		} else if exists {
			break
		}
		if currentPath == "." || currentPath == "/" {
			return "", nil
		}
	}

	// returned path is always an absolute path
	return currentPath, nil
}

func (r *cmdRunner) Configure(app *kingpin.Application) error {
	app.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)

	commands := []Command{
		&protobuf{},
		&intsCmd{},
		&datamodelCmd{},
		&codegenCmd{},
		&sequenceDiagramCmd{},
		&validateCmd{},
		&importSwaggerCmd{},
		&infoCmd{},
	}
	r.commands = map[string]Command{}

	app.Flag("root",
		"sysl root directory for input model file (default: .)").
		Default("").StringVar(&r.Root)

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
