package main

import (
	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
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
			var modelAppName string
			var err error
			if cmd.RequireSyslModule() {
				logger.Infof("Attempting to load module:%s (root:%s)", r.module, r.Root)
				modelParser := parse.NewParser()
				mod, modelAppName, err = parse.LoadAndGetDefaultApp(r.module, syslutil.NewChrootFs(fs, r.Root), modelParser)
				if err != nil {
					return err
				}
			}

			return cmd.Execute(ExecuteArgs{mod, modelAppName, fs, logger})
		}
	}
	return nil
}

func (r *cmdRunner) Configure(app *kingpin.Application) error {
	commands := []Command{
		&protobuf{},

		&validateCmd{},
	}
	r.commands = map[string]Command{}

	app.Flag("root",
		"sysl root directory for input model file (default: .)").
		Default(".").StringVar(&r.Root)

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

// Temp until all the commands move in here
func (r *cmdRunner) HasCommand(which string) bool {

	_, ok := r.commands[which]
	return ok
}
