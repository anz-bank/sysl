package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/anz-bank/sysl/pkg/cfg"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

// main3 is the real main function. It takes its output streams and command-line
// arguments as parameters to support testability.
func main3(args []string, fs afero.Fs, logger *logrus.Logger, stdin io.Reader, stdout io.Writer) error {
	flags, err := syslutil.PopulateCMDFlagsFromFile(args)
	if err == nil && len(flags) > 0 {
		// apply flags in file
		args = flags
	}

	syslCmd := kingpin.New("sysl", "System Specification Language Toolkit")
	syslCmd.Version(fmt.Sprintf("sysl %s %s", Version, BuildOS))
	syslCmd.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)

	(&debugTypeData{}).add(syslCmd)

	runner := cmdRunner{}
	if err := runner.Configure(syslCmd); err != nil {
		return err
	}

	selectedCommand, err := syslCmd.Parse(args[1:])
	if err != nil {
		return err
	}

	// clean paths for multiplatform compatibility
	for i, val := range runner.modules {
		runner.modules[i] = path.Clean(val)
	}

	return runner.Run(selectedCommand, fs, logger, stdin, stdout)
}

type debugTypeData struct {
	loglevel string
	verbose  bool
}

//nolint:unparam
func (d *debugTypeData) do(_ *kingpin.ParseContext) error {
	if d.verbose {
		d.loglevel = cfg.LogLevel
	}
	// Default info
	syslutil.SetLogLevel(d.loglevel)

	logrus.Debugf("Logging: %+v", *d)
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
		Default("warn").
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
	stdin io.Reader,
	stdout io.Writer,
	main3 func(args []string, fs afero.Fs, logger *logrus.Logger, stdin io.Reader, stdout io.Writer) error,
) int {
	if err := main3(args, fs, logger, stdin, stdout); err != nil {
		arraiErr, ok := errors.Cause(err).(arrai.ExecutionError)
		var exitCode = 1
		if ok {
			logger.Debugln(err)
			logger.Errorln(arraiErr.ShortMsg)
		} else {
			if err, ok := err.(syslutil.Exit); ok {
				exitCode = err.Code
			}

			switch exitCode {
			case 0:
				logger.Infoln(err.Error())
			case 2:
				logger.Warnln(err.Error())
			default:
				logger.Errorln(err.Error())
			}
		}
		return exitCode
	}
	return 0
}

// main is as small as possible to minimise its no-coverage footprint.
func main() {
	if rc := main2(os.Args, afero.NewOsFs(), logrus.StandardLogger(), os.Stdin, os.Stdout, main3); rc != 0 {
		os.Exit(rc)
	}
}
