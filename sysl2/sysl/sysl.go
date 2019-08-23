package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/pbutil"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/anz-bank/sysl/sysl2/sysl/validate"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const debug string = "debug"

// main3 is the real main function. It takes its output streams and command-line
// arguments as parameters to support testability.
func main3(args []string) error {
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	textpbParams := configureCmdlineForPb(sysl)
	codegenParams := configureCmdlineForCodegen(sysl)
	sequenceParams := configureCmdlineForSeqgen(sysl)
	intgenParams := configureCmdlineForIntgen(sysl)
	validateParams := validate.ConfigureCmdlineForValidate(sysl)
	var selectedCommand string
	var err error
	if selectedCommand, err = sysl.Parse(args[1:]); err != nil {
		return err
	}
	switch selectedCommand {
	case "pb":
		return doGeneratePb(textpbParams)
	case "gen":
		output := GenerateCode(codegenParams)
		return outputToFiles(*codegenParams.outDir, output)
	case "sd":
		result, err := DoConstructSequenceDiagrams(sequenceParams)
		if err != nil {
			return err
		}
		for k, v := range result {
			if err := OutputPlantuml(k, *sequenceParams.plantuml, v); err != nil {
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
			err := OutputPlantuml(k, *intgenParams.plantuml, v)
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
func main2(args []string, main3 func(args []string) error) int {
	if err := main3(args); err != nil {
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
	if rc := main2(os.Args, main3); rc != 0 {
		os.Exit(rc)
	}
}

func configureCmdlineForPb(sysl *kingpin.Application) *CmdContextParamPbgen {
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

func doGeneratePb(textpbParams *CmdContextParamPbgen) error {
	logrus.Debugf("Root: %s\n", *textpbParams.root)
	logrus.Debugf("Module: %s\n", *textpbParams.modules)
	logrus.Debugf("Mode: %s\n", *textpbParams.mode)
	logrus.Debugf("Log Level: %s\n", *textpbParams.loglevel)

	format := strings.ToLower(*textpbParams.output)
	toJSON := *textpbParams.mode == "json" || *textpbParams.mode == "" && strings.HasSuffix(format, ".json")
	logrus.Debugf("%s\n", *textpbParams.modules)
	mod, err := parse.NewParser().Parse(*textpbParams.modules, *textpbParams.root)
	*textpbParams.output = strings.Trim(*textpbParams.output, " ")
	if err != nil {
		return err
	}

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
			return pbutil.JSONPB(mod, *textpbParams.output)
		}
		if *textpbParams.output == "-" {
			return pbutil.FTextPB(logrus.StandardLogger().Out, mod)
		}
		return pbutil.TextPB(mod, *textpbParams.output)
	}
	return nil
}
