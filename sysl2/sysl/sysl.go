package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/pbutil"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/anz-bank/sysl/sysl2/sysl/validate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

const debug string = "debug"

// main3 is the real main function. It takes its output streams and command-line
// arguments as parameters to support testability.
func main3(args []string, fs afero.Fs) error {
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	flagmap := map[string][]string{}
	textpbParams := configureCmdlineForPb(sysl, flagmap)
	codegenParams := configureCmdlineForCodegen(sysl, flagmap)
	sequenceParams := configureCmdlineForSeqgen(sysl, flagmap)
	intgenParams := configureCmdlineForIntgen(sysl, flagmap)
	validateParams := validate.ConfigureCmdlineForValidate(sysl)
	var selectedCommand string
	var err error
	sysl.Validate(generateAppargValidator(args[1], flagmap))
	if selectedCommand, err = sysl.Parse(args[1:]); err != nil {
		return err
	}
	switch selectedCommand {
	case "pb":
		return doGeneratePb(textpbParams, fs)
	case "gen":
		output, err := GenerateCode(codegenParams, fs)
		if err != nil {
			return err
		}
		return outputToFiles(output, syslutil.NewChrootFs(fs, *codegenParams.outDir))
	case "sd":
		result, err := DoConstructSequenceDiagrams(sequenceParams)
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
		r, err := GenerateIntegrations(intgenParams)
		if err != nil {
			return err
		}
		for k, v := range r {
			if *intgenParams.isVerbose {
				logrus.Debugf(k)
			}
			err := OutputPlantuml(k, *intgenParams.plantuml, v, fs)
			if err != nil {
				return fmt.Errorf("plantuml drawing error: %v", err)
			}
		}
		return nil
	case "validate":
		return validate.DoValidate(validateParams)
	}
	return nil
}

// main2 calls main3 and handles any errors it returns. It takes its output
// streams and command-line arguments and even main3 as parameters to support
// testability.
func main2(args []string, fs afero.Fs, main3 func(args []string, fs afero.Fs) error) int {
	if err := main3(args, fs); err != nil {
		logrus.Errorln(err.Error())
		if err, ok := err.(parse.Exit); ok {
			return err.Code
		}
		return 1
	}
	return 0
}

// main is as small as possible to minimise its no-coverage footprint.
func main() {
	if rc := main2(os.Args, afero.NewOsFs(), main3); rc != 0 {
		os.Exit(rc)
	}
}

func configureCmdlineForPb(sysl *kingpin.Application, flagmap map[string][]string) *CmdContextParamPbgen {
	flagmap["pb"] = []string{"root", "output", "mode"}
	textpb := sysl.Command("pb", "Generate textpb/json")
	returnValues := &CmdContextParamPbgen{}

	returnValues.root = textpb.Flag("root",
		"sysl root directory for input model file (default: .)",
	).Default(".").String()

	returnValues.output = textpb.Flag("output", "output file name").Short('o').String()
	returnValues.mode = textpb.Flag("mode", "output mode").Default("textpb").String()

	returnValues.loglevel = textpb.Flag("log",
		"log level[debug,info,warn,off]").Default("info").String()

	returnValues.isVerbose = textpb.Flag("verbose", "show output").Short('v').Default("false").Bool()

	returnValues.modules = textpb.Arg("modules", "input files without .sysl extension and with leading /, eg: "+
		"/project_dir/my_models combine with --root if needed",
	).String()

	return returnValues
}

func doGeneratePb(textpbParams *CmdContextParamPbgen, fs afero.Fs) error {
	logrus.Debugf("Root: %s\n", *textpbParams.root)
	logrus.Debugf("Module: %s\n", *textpbParams.modules)
	logrus.Debugf("Mode: %s\n", *textpbParams.mode)
	logrus.Debugf("Log Level: %s\n", *textpbParams.loglevel)

	format := strings.ToLower(*textpbParams.output)
	toJSON := *textpbParams.mode == "json" || *textpbParams.mode == "" && strings.HasSuffix(format, ".json")
	logrus.Debugf("%s\n", *textpbParams.modules)
	mod, err := parse.NewParser().Parse(*textpbParams.modules, syslutil.NewChrootFs(fs, *textpbParams.root))
	if err != nil {
		return err
	}
	*textpbParams.output = strings.Trim(*textpbParams.output, " ")

	if *textpbParams.isVerbose {
		*textpbParams.loglevel = debug
	}
	// Default info
	if level, has := syslutil.LogLevels[*textpbParams.loglevel]; has {
		logrus.SetLevel(level)
	}

	switch *textpbParams.mode {
	case "", "textpb", "json":
	default:
		return fmt.Errorf("invalid -mode %#v", *textpbParams.mode)
	}

	if mod != nil {
		if toJSON {
			if *textpbParams.output == "-" {
				return pbutil.FJSONPB(logrus.StandardLogger().Out, mod)
			}
			return pbutil.JSONPB(mod, *textpbParams.output, fs)
		}
		if *textpbParams.output == "-" {
			return pbutil.FTextPB(logrus.StandardLogger().Out, mod)
		}
		return pbutil.TextPB(mod, *textpbParams.output, fs)
	}
	return nil
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
