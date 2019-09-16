package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/anz-bank/sysl/sysl2/sysl/commands"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/anz-bank/sysl/sysl2/sysl/validate"
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
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")
	syslCmd.Version(Version)

	(&debugTypeData{}).add(syslCmd)

	runner := commands.Runner{}
	if err := runner.Init(syslCmd); err != nil {
		return err
	}

	flagmap := map[string][]string{}
	codegenParams := configureCmdlineForCodegen(syslCmd, flagmap)
	sequenceParams := configureCmdlineForSeqgen(syslCmd, flagmap)
	intgenParams := configureCmdlineForIntgen(syslCmd, flagmap)
	validateParams := validate.ConfigureCmdlineForValidate(syslCmd)
	datagenParams := configureCmdlineForDatagen(syslCmd)

	var selectedCommand string
	var err error
	if len(args) > 1 {
		syslCmd.Validate(generateAppargValidator(args[1], flagmap))
	}
	if selectedCommand, err = syslCmd.Parse(args[1:]); err != nil {
		return err
	}

	if runner.HasCommand(selectedCommand) {
		return runner.Run(selectedCommand, fs, logger)
	}

	switch selectedCommand {
	case "gen":
		codegenParams.rootModel = &runner.Root
		output, err := GenerateCode(codegenParams, fs, logger)
		if err != nil {
			return err
		}
		return outputToFiles(output, syslutil.NewChrootFs(fs, *codegenParams.outDir))
	case "sd":
		sequenceParams.root = &runner.Root
		result, err := DoConstructSequenceDiagrams(sequenceParams, logger)
		if err != nil {
			return err
		}
		for k, v := range result {
			if err := OutputPlantuml(k, *sequenceParams.plantuml, v, fs); err != nil {
				return err
			}
		}
		return nil
	case "ints":
		intgenParams.root = &runner.Root
		r, err := GenerateIntegrations(intgenParams)
		if err != nil {
			return err
		}
		for k, v := range r {
			err := OutputPlantuml(k, *intgenParams.plantuml, v, fs)
			if err != nil {
				return fmt.Errorf("plantuml drawing error: %v", err)
			}
		}
		return nil
	case "validate":
		return validate.DoValidate(validateParams)
	case "data":
		datagenParams.root = &runner.Root
		outmap, err := GenerateDataModels(datagenParams)
		if err != nil {
			return err
		}
		for k, v := range outmap {
			err := OutputPlantuml(k, *datagenParams.plantuml, v, fs)
			if err != nil {
				return fmt.Errorf("plantuml drawing error: %v", err)
			}
		}
		return nil
	}
	return nil
}

type debugTypeData struct {
	loglevel string
	verbose  bool
}

//nolint:unparam
func (d *debugTypeData) do(_ *kingpin.ParseContext) error {
	if d.verbose {
		d.loglevel = debug
	}
	// Default info
	if level, has := syslutil.LogLevels[d.loglevel]; has {
		logrus.SetLevel(level)
	}

	logrus.Infof("Logging: %+v", *d)
	return nil
}
func (d *debugTypeData) add(app *kingpin.Application) {

	var levels []string
	for l := range syslutil.LogLevels {
		if l != "" {
			levels = append(levels, l)
		}
	}
	app.Flag("log", fmt.Sprintf("log level: [%s]", strings.Join(levels, ","))).
		HintOptions(levels...).
		Default(levels[0]).
		StringVar(&d.loglevel)
	app.Flag("verbose", "enable verbose logging").Short('v').BoolVar(&d.verbose)
	app.PreAction(d.do)
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

func generateAppargValidator(selectedCommand string, flags map[string][]string) func(*kingpin.Application) error {
	return func(app *kingpin.Application) error {
		var errorMsg strings.Builder
		for _, longFlagName := range flags[selectedCommand] {
			if flag := app.GetCommand(selectedCommand).GetFlag(longFlagName); flag != nil {
				val := flag.Model().Value.String()
				if val != "" {
					val = strings.Trim(val, " ")
					if val == "" {
						errorMsg.WriteString("'" + longFlagName + "'" + " value passed is empty\n")
					}
				} else if len(flag.Model().Default) > 0 {
					errorMsg.WriteString("'" + longFlagName + "'" + " value passed is empty\n")
				}
			}
		}
		if errorMsg.Len() > 0 {
			return errors.New(errorMsg.String())
		}
		return nil
	}
}
