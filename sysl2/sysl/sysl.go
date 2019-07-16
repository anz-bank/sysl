package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//nolint:gochecknoglobals
var defaultLevel = map[string]logrus.Level{
	"":      logrus.ErrorLevel,
	"off":   logrus.ErrorLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
}

func (e exit) Error() string {
	return e.message
}

// main3 is the real main function. It takes its output streams and command-line
// arguments as parameters to support testability.
func main3(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.PanicOnError)

	if len(args) > 1 && args[1] == "validate" {
		return DoValidate(flags, args)
	}

	switch filepath.Base(args[0]) {
	case "syslgen":
		return DoGenerateCode(flags, args)
	case "sd":
		return DoGenerateSequenceDiagrams(args)
	case "ints":
		DoGenerateIntegrations(args)
		return nil
	}
	root := flags.String("root", ".", "sysl root directory for input files (default: .)")
	output := flags.String("o", "", "output file name")
	mode := flags.String("mode", "textpb", "output mode")
	loglevel := flags.String("log", "warn", "log level[debug,info,warn,off]")

	//nolint:errcheck
	flags.Parse(args[1:])
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	switch *mode {
	case "", "textpb", "json":
	default:
		return fmt.Errorf("invalid -mode %#v", *mode)
	}

	if level, has := defaultLevel[*loglevel]; has {
		logrus.SetLevel(level)
	} else {
		return fmt.Errorf("invalid -log %#v", *loglevel)
	}

	filename := flags.Arg(0)

	log.Infof("Args: %v\n", flags.Args())
	log.Infof("Root: %s\n", *root)
	log.Infof("Module: %s\n", filename)
	log.Infof("Mode: %s\n", *mode)
	log.Infof("Log Level: %s\n", *loglevel)
	format := strings.ToLower(*output)
	toJSON := *mode == "json" || *mode == "" && strings.HasSuffix(format, ".json")
	log.Infof("%s\n", filename)
	mod, err := Parse(filename, *root)
	if err != nil {
		return err
	}
	if mod != nil {
		if toJSON {
			if *output == "-" {
				return FJSONPB(log.StandardLogger().Out, mod)
			}
			return JSONPB(mod, *output)
		}
		if *output == "-" {
			return FTextPB(log.StandardLogger().Out, mod)
		}
		return TextPB(mod, *output)
	}
	return nil
}

// main2 calls main3 and handles any errors it returns. It takes its output
// streams and command-line arguments and even main3 as parameters to support
// testability.
func main2(args []string, main3 func(args []string) error,
) int {
	if err := main3(args); err != nil {
		log.Errorln(err.Error())
		if err, ok := err.(exit); ok {
			return err.code
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
