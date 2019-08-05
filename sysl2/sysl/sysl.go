package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/pbutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

//nolint:gochecknoglobals
var defaultLevel = map[string]logrus.Level{
	"":      logrus.ErrorLevel,
	"off":   logrus.ErrorLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
}

const debug string = "debug"

// main3 is the real main function. It takes its output streams and command-line
// arguments as parameters to support testability.
func main3(args []string) error {
	sysl := kingpin.New("sysl", "System Modelling Language Toolkit")
	textpbParams, err1 := configureCmdlineForPb(sysl)
	if err1 != nil {
		return err1
	}
	codegenParams := configureCmdlineForCodegen(sysl)
	sequenceParams := configureCmdlineForSeqdiaggen(sysl)
	intgenParams := configureCmdlineForIntgen(sysl)
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
			if *intgenParams.verbose {
				fmt.Println(k)
			}
			err := OutputPlantuml(k, *intgenParams.plantuml, v)
			if err != nil {
				return fmt.Errorf("plantuml drawing error: %v", err)
			}
		}
		return nil
	}
	return nil
}

// main2 calls main3 and handles any errors it returns. It takes its output
// streams and command-line arguments and even main3 as parameters to support
// testability.
func main2(args []string, main3 func(args []string) error,
) int {
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

func configureCmdlineForPb(sysl *kingpin.Application) (*CmdContextParamPbgen, error) {
	textpb := sysl.Command("pb", "Generate textpb/json")
	returnValues := &CmdContextParamPbgen{}

	returnValues.root = textpb.Flag("root",
		"sysl root directory for input model file (default: .)",
	).Default(".").String()

	returnValues.output = textpb.Flag("output", "output file name").Short('o').String()
	returnValues.mode = textpb.Flag("mode", "output mode").Default("textpb").String()

	returnValues.loglevel = textpb.Flag("log",
		"log level[debug,info,warn,off]").Default("warn").String()

	returnValues.verbose = textpb.Flag("verbose", "show output").Short('v').Default("false").Bool()

	returnValues.modules = textpb.Arg("modules", "input files without .sysl extension and with leading /, eg: "+
		"/project_dir/my_models combine with --root if needed",
	).String()

	switch *returnValues.mode {
	case "", "textpb", "json":
	default:
		return nil, fmt.Errorf("invalid -mode %#v", *returnValues.mode)
	}

	return returnValues, nil
}

func doGeneratePb(textpbParams *CmdContextParamPbgen) error {
	logrus.Infof("Root: %s\n", *textpbParams.root)
	logrus.Infof("Module: %s\n", *textpbParams.modules)
	logrus.Infof("Mode: %s\n", *textpbParams.mode)
	logrus.Infof("Log Level: %s\n", *textpbParams.loglevel)

	format := strings.ToLower(*textpbParams.output)
	toJSON := *textpbParams.mode == "json" || *textpbParams.mode == "" && strings.HasSuffix(format, ".json")
	logrus.Infof("%s\n", *textpbParams.modules)
	mod, err := parse.Parse(*textpbParams.modules, *textpbParams.root)

	if *textpbParams.verbose {
		*textpbParams.loglevel = debug
	}
	if level, has := defaultLevel[*textpbParams.loglevel]; has {
		logrus.SetLevel(level)
	} else {
		return fmt.Errorf("invalid -log %#v", *textpbParams.loglevel)
	}

	if err != nil {
		return err
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
