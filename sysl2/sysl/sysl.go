package main

import (
	"fmt"
	sysl "github.com/anz-bank/sysl/src/proto"
	"os"
	"strings"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Version contains the sysl binary version
//nolint:gochecknoglobals
var Version = "unspecified"

const debug string = "debug"

// main3 is the real main function. It takes its output streams and command-line
// arguments as parameters to support testability.
func main3(args []string, fs afero.Fs, logger *logrus.Logger) error {
	syslCmd, cmdctx := initCommandLineArgs(logger)

	var selectedCommand string
	var err error

	if selectedCommand, err = syslCmd.Parse(args[1:]); err != nil {
		return err
	}

	for _, cmd := range cmdctx.commands {
		if cmd.Name() == selectedCommand {
			var mod *sysl.Module
			var modelAppName string
			if cmd.RequireSyslModule() {
				logger.Infof("Attempting to load module:%s (root:%s)", cmdctx.module, cmdctx.root)
				modelParser := parse.NewParser()
				mod, modelAppName, err = parse.LoadAndGetDefaultApp(cmdctx.module, syslutil.NewChrootFs(fs, cmdctx.root), modelParser)
				if err != nil {
					return err
				}
			}

			return cmd.Execute(ExecuteArgs{mod, modelAppName, fs, logger})
		}
	}


	return nil
}

type debugTypeData struct {
	loglevel string
	verbose  bool
	logger *logrus.Logger
}

func (d *debugTypeData) do(_ *kingpin.ParseContext) error {
	if d.verbose {
		d.loglevel = debug
	}
	// Default info
	if level, has := syslutil.LogLevels[d.loglevel]; has {
		d.logger.SetLevel(level)
	}

	d.logger.Infof("Logging: %+v", *d)
	return nil
}
func (d *debugTypeData) add(app *kingpin.Application, logger *logrus.Logger) {

	d.logger = logger

	var levels []string
	for l := range syslutil.LogLevels {
		if l != "" {
			levels = append(levels, l)
		}
	}
	app.Flag("log", fmt.Sprintf("log level: [%s]", strings.Join(levels, ","))).
		Default(levels[0]).
		EnumVar(&d.loglevel, levels...)
	app.Flag("verbose", "enable verbose logging").Short('v').BoolVar(&d.verbose)
	app.PreAction(d.do)
}

type cmdContext struct {
	root string
	module string
	commands []Command
}

func initCommandLineArgs(logger *logrus.Logger) (*kingpin.Application, *cmdContext) {
	app := kingpin.New("sysl", "System Modelling Language Toolkit")
	app.Version(Version)


	var ctx cmdContext

	(&debugTypeData{}).add(app, logger)

	ctx.commands = []Command{
		&pbCommand{},
		&genCmd{},
		&sdCmd{},
		&intsCmd{},
		&dmCmd{},
		&validateCmd{},
	}


	app.Flag("root",
		"sysl root directory for input model file (default: .)").
		Default(".").StringVar(&ctx.root)

	for _, cmd := range ctx.commands {
		c := cmd.Init(app)
		if cmd.RequireSyslModule() {
			c.Arg("MODULE", "input files without .sysl extension and with leading /, eg: "+
				"/project_dir/my_models combine with --root if needed").
				Required().StringVar(&ctx.module)
		}
	}

	return app, &ctx
}

// main2 calls main3 and handles any errors it returns. It takes its output
// streams and command-line arguments and even main3 as parameters to support
// testability.
func main2(
	args []string,
	fs afero.Fs,
	logger *logrus.Logger,
	main3 func(args []string, fs afero.Fs, logger *logrus.Logger) error,
) int {
	if err := main3(args, fs, logger); err != nil {
		logger.Errorln(err.Error())
		if err, ok := err.(parse.Exit); ok {
			return err.Code
		}
		return 1
	}
	return 0
}

// main is as small as possible to minimise its no-coverage footprint.
func main() {
	if rc := main2(os.Args, afero.NewOsFs(), logrus.StandardLogger(), main3); rc != 0 {
		os.Exit(rc)
	}
}
