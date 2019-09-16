package commands

import (
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type Runner struct {
	commands []Command

	Root   string
	module string
}

func (r *Runner) Run(which string, fs afero.Fs, logger *logrus.Logger) error {
	for _, cmd := range r.commands {
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

func (r *Runner) Init(app *kingpin.Application) error {
	r.commands = []Command{
		&protobuf{},
	}

	app.Flag("root",
		"sysl root directory for input model file (default: .)").
		Default(".").StringVar(&r.Root)

	for _, cmd := range r.commands {
		c := cmd.Init(app)
		if cmd.RequireSyslModule() {
			c.Arg("MODULE", "input files without .sysl extension and with leading /, eg: "+
				"/project_dir/my_models combine with --root if needed").
				Required().StringVar(&r.module)
		}
	}

	return nil
}

// Temp until all the commands move in here
func (r *Runner) HasCommand(which string) bool {
	for _, cmd := range r.commands {
		if strings.EqualFold(cmd.Name(), which) {
			return true
		}
	}
	return false
}
