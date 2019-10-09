package main

import (
	"errors"
	"sort"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/roothandler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	syslRootMarker = ".SYSL_ROOT"
)

type cmdRunner struct {
	commands map[string]Command

	RootStatus *roothandler.RootStatus
	root       string
	module     string
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

				mod, _, err = LoadSyslModule(r, r.module, fs, logger)
				if err != nil {
					return err
				}
			}

			// TODO: what filesystem to pass for execution?
			return cmd.Execute(ExecuteArgs{mod, fs, logger})
		}
	}
	return nil
}

func (r *cmdRunner) rootHandler(fs afero.Fs, logger *logrus.Logger) error {
	r.RootStatus = roothandler.NewRootStatus(r.root)
	return r.RootStatus.RootHandler(r.module, fs, logger)
}

func (r *cmdRunner) Configure(app *kingpin.Application) error {
	app.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)

	commands := []Command{
		&protobuf{},
		&intsCmd{},
		&datamodelCmd{},
		&codegenCmd{},
		&sequenceDiagramCmd{},
		&importCmd{},
		&infoCmd{},
		&validateCmd{},
		&swaggerExportCmd{},
	}
	r.commands = map[string]Command{}

	app.Flag("root",
		"sysl root directory for input model file (default: .)").
		Default("").StringVar(&r.root)

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
