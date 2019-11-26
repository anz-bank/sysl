package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto_old"
	"github.com/anz-bank/sysl/src/sysl/parse"
	"github.com/anz-bank/sysl/src/sysl/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Version   - Binary version
// GitCommit - Commit SHA of the source code
// BuildDate - Binary build date
// BuildOS   - Operating System used to build binary
//nolint:gochecknoglobals
var (
	Version   = "unspecified"
	GitCommit = "unspecified"
	BuildDate = "unspecified"
	BuildOS   = "unspecified"
)

const debug string = "debug"

type projectConfiguration struct {
	module, root string
	rootIsFound  bool
	fs           afero.Fs
}

func LoadSyslModule(root, filename string, fs afero.Fs, logger *logrus.Logger) (*sysl.Module, string, error) {
	logger.Debugf("Attempting to load module:%s (root:%s)", filename, root)
	projectConfig := newProjectConfiguration()
	if err := projectConfig.configureProject(root, filename, fs, logger); err != nil {
		return nil, "", err
	}

	modelParser := parse.NewParser()
	if !projectConfig.rootIsFound {
		modelParser.RestrictToLocalImport()
	}
	return parse.LoadAndGetDefaultApp(projectConfig.module, projectConfig.fs, modelParser)
}

func newProjectConfiguration() *projectConfiguration {
	return &projectConfiguration{
		root:        "",
		module:      "",
		rootIsFound: false,
		fs:          nil,
	}
}

func (pc *projectConfiguration) configureProject(root, module string, fs afero.Fs, logger *logrus.Logger) error {
	rootIsDefined := root != ""

	modulePath := module
	if rootIsDefined {
		modulePath = filepath.Join(root, module)
	}

	syslRootPath, err := findRootFromSyslModule(modulePath, fs)
	if err != nil {
		return err
	}

	rootMarkerExists := syslRootPath != ""

	if rootIsDefined {
		pc.rootIsFound = true
		pc.root = root
		pc.module = module
		if rootMarkerExists {
			logger.Warningf("%s found in %s but will use %s instead",
				syslRootMarker, syslRootPath, pc.root)
		} else {
			logger.Warningf("%s is not defined but root flag is defined in %s",
				syslRootMarker, pc.root)
		}
	} else {
		if rootMarkerExists {
			pc.root = syslRootPath

			// module has to be relative to the root
			absModulePath, err := filepath.Abs(module)
			if err != nil {
				return err
			}
			pc.module, err = filepath.Rel(pc.root, absModulePath)
			if err != nil {
				return err
			}
			pc.rootIsFound = true
		} else {
			// uses the module directory as the root, changing the module to be relative to the root
			pc.root = filepath.Dir(module)
			pc.module = filepath.Base(module)
			pc.rootIsFound = false
			logger.Warningf("root and %s are undefined, %s will be used instead",
				syslRootMarker, pc.root)
		}
	}

	pc.fs = syslutil.NewChrootFs(fs, pc.root)
	return nil
}

func findRootFromSyslModule(modulePath string, fs afero.Fs) (string, error) {
	currentPath, err := filepath.Abs(modulePath)
	if err != nil {
		return "", err
	}

	systemRoot, err := filepath.Abs(string(os.PathSeparator))
	if err != nil {
		return "", err
	}

	// Keep walking up the directories to find nearest root marker
	for {
		currentPath = filepath.Dir(currentPath)
		exists, err := afero.Exists(fs, filepath.Join(currentPath, syslRootMarker))
		reachedRoot := currentPath == systemRoot || (err != nil && os.IsPermission(err))
		switch {
		case exists:
			return currentPath, nil
		case reachedRoot:
			return "", nil
		case err != nil:
			return "", err
		}
	}
}

// main3 is the real main function. It takes its output streams and command-line
// arguments as parameters to support testability.
func main3(args []string, fs afero.Fs, logger *logrus.Logger) error {
	syslCmd := kingpin.New("sysl", "System Modelling Language Toolkit")
	syslCmd.Version(Version)
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

	return runner.Run(selectedCommand, fs, logger)
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
