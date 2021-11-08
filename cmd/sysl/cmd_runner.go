package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/loader"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdRunner struct {
	commands map[string]cmdutils.Command

	Root    string
	modules []string
}

// Run identifies the command to run, loads the Sysl modules from the input (if necessary), then
// executes the command with all of the accumulated context.
func (r *cmdRunner) Run(which string, fs afero.Fs, logger *logrus.Logger, stdin io.Reader) error {
	// splitter to parse main command from subcommand
	mainCommand := strings.Split(which, " ")[0]
	if cmd, ok := r.commands[mainCommand]; ok {
		if cmd.Name() == mainCommand {
			var defaultAppName string
			var modules []*sysl.Module
			var err error

			if cmd.MaxSyslModule() > 0 {
				if len(r.modules) > 0 {
					modules, defaultAppName, err = r.loadFromModules(fs, logger)
					// stdin may still be provided for use by commands like transform.
				} else {
					modules, defaultAppName, err = r.loadFromStdin(stdin, fs, logger)
				}
				if err != nil {
					return err
				}
			}

			if len(modules) > cmd.MaxSyslModule() {
				return fmt.Errorf("this command can accept max %d module(s)", cmd.MaxSyslModule())
			}
			return cmd.Execute(cmdutils.ExecuteArgs{
				Command:        which,
				Modules:        modules,
				Filesystem:     fs,
				Logger:         logger,
				DefaultAppName: defaultAppName,
				ModulePaths:    r.modules,
				Root:           r.Root,
				Stdin:          stdin,
			})
		}
	}
	return nil
}

func (r *cmdRunner) Configure(app *kingpin.Application) error {
	commands := []cmdutils.Command{
		&codegenCmd{},
		&databaseScriptCmd{},
		&datamodelCmd{},
		&diagramCmd{},
		&envCmd{},
		&exportCmd{},
		&importCmd{},
		&infoCmd{},
		&intsCmd{},
		&lspCmd{},
		&modDatabaseScriptCmd{},
		&protobufCmd{},
		&replCmd{},
		&sequenceDiagramCmd{},
		&templateCmd{},
		&testRigCmd{},
		&transformCmd{},
		&validateCmd{},
	}
	r.commands = map[string]cmdutils.Command{}

	app.Flag("root",
		"sysl root directory for input model file. If root is not found, the module directory becomes "+
			"the root, but the module can not import with absolute paths (or imports must be relative).").StringVar(&r.Root)

	sort.Slice(commands, func(i, j int) bool {
		return strings.Compare(commands[i].Name(), commands[j].Name()) < 0
	})
	for _, cmd := range commands {
		r.ConfigureCmd(app, cmd)
	}

	return nil
}

func (r *cmdRunner) ConfigureCmd(app *kingpin.Application, cmd cmdutils.Command) {
	c := cmd.Configure(app)
	if cmd.MaxSyslModule() > 0 {
		c.Arg("MODULE", "input files without .sysl extension and with leading /, eg: "+
			"/project_dir/my_models combine with --root if needed").
			StringsVar(&r.modules)
	}
	r.commands[cmd.Name()] = cmd
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
		for _, elem := range c.Elements {
			if f, _ := elem.Clause.(*kingpin.FlagClause); f != nil && f.Model().Name == "help" {
				return nil // help requested, don't need to check for empty flags
			}
		}
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

// loadFromStdin attempts to load the Sysl modules for the files provided via stdin.
func (r *cmdRunner) loadFromStdin(stdin io.Reader, fs afero.Fs, logger *logrus.Logger) ([]*sysl.Module, string, error) {
	stdinFiles, err := loadStdinFiles(stdin)
	if err != nil {
		return nil, "", err
	}

	fs = afero.NewCopyOnWriteFs(fs, afero.NewMemMapFs())
	for _, f := range stdinFiles {
		r.modules = append(r.modules, f.Path)
		absPath, err := filepath.Abs(f.Path)
		if err != nil {
			return nil, "", err
		}
		err = afero.WriteFile(fs, absPath, []byte(f.Content), os.ModePerm)
		if err != nil {
			return nil, "", err
		}
	}
	return r.loadFromModules(fs, logger)
}

// loadFromModules attempts to load the Sysl modules for the files specified in r.modules.
func (r *cmdRunner) loadFromModules(fs afero.Fs, logger *logrus.Logger) ([]*sysl.Module, string, error) {
	var mods []*sysl.Module
	for _, moduleName := range r.modules {
		mod, _, err := loader.LoadSyslModule(r.Root, moduleName, fs, logger)
		if err != nil {
			return nil, "", err
		}
		mods = append(mods, mod)
	}
	return mods, "", nil
}
