package main

import (
	"errors"
	"sort"
	"strings"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

// The idea of syslRootMarker is to be similar to how .git folder works
// Might have more features using the folder but right now it's just an empty folder
const syslRootMarker = ".sysl"

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
			var appName string
			if cmd.RequireSyslModule() {
				mod, appName, err = LoadSyslModule(r.Root, r.module, fs, logger)
				if err != nil {
					return err
				}
			}

			return cmd.Execute(ExecuteArgs{Module: mod, Filesystem: fs, Logger: logger, DefaultAppName: appName})
		}
	}
	logger.Debugf("command %s was configured but not in CmdRunner's commands", which)
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
		&replCmd{},
		&modCmd{},
	}
	r.commands = map[string]Command{}

	app.Flag("root",
		"sysl root directory for input model file. If root is not found, the module directory becomes "+
			"the root, but the module can not import with absolute paths (or imports must be relative).").StringVar(&r.Root)

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
