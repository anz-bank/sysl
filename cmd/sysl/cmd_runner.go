package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
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

func (r *cmdRunner) Run(which string, fs afero.Fs, logger *logrus.Logger) error {
	// splitter to parse main command from subcommand
	mainCommand := strings.Split(which, " ")[0]
	if cmd, ok := r.commands[mainCommand]; ok {
		if cmd.Name() == mainCommand {
			var module *sysl.Module
			var err error
			var appName string
			var mods []*sysl.Module

			if cmd.MaxSyslModule() > 0 {
				for _, moduleName := range r.modules {
					module, appName, err = loader.LoadSyslModule(r.Root, moduleName, fs, logger)
					if err != nil {
						return err
					}
					mods = append(mods, module)
				}
			}

			if len(mods) > cmd.MaxSyslModule() {
				logger.Error("this command can accept max " + strconv.Itoa(cmd.MaxSyslModule()) + " module(s).")
				return fmt.Errorf("this command can accept max " + strconv.Itoa(cmd.MaxSyslModule()) + " module(s).")
			}
			return cmd.Execute(cmdutils.ExecuteArgs{Command: which, Modules: mods, Filesystem: fs,
				Logger: logger, DefaultAppName: appName, ModulePaths: r.modules})
		}
	}
	return nil
}

func (r *cmdRunner) Configure(app *kingpin.Application) error {
	commands := []cmdutils.Command{
		&protobuf{},
		&intsCmd{},
		&datamodelCmd{},
		&databaseScriptCmd{},
		&modDatabaseScriptCmd{},
		&codegenCmd{},
		&sequenceDiagramCmd{},
		&diagramCmd{},
		&importCmd{},
		&infoCmd{},
		&validateCmd{},
		&exportCmd{},
		&replCmd{},
		&envCmd{},
		&templateCmd{},
		&modCmd{},
		&testRigCmd{},
	}
	r.commands = map[string]cmdutils.Command{}

	app.Flag("root",
		"sysl root directory for input model file. If root is not found, the module directory becomes "+
			"the root, but the module can not import with absolute paths (or imports must be relative).").StringVar(&r.Root)

	sort.Slice(commands, func(i, j int) bool {
		return strings.Compare(commands[i].Name(), commands[j].Name()) < 0
	})
	for _, cmd := range commands {
		c := cmd.Configure(app)
		if cmd.MaxSyslModule() > 0 {
			c.Arg("MODULE", "input files without .sysl extension and with leading /, eg: "+
				"/project_dir/my_models combine with --root if needed").
				Required().StringsVar(&r.modules)
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
