package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	rootUndefined  = "\000"
	syslRootMarker = ".SYSL_ROOT"
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
			var err error
			if cmd.RequireSyslModule() {
				mod, _, err = LoadSyslModule(r.Root, r.module, fs, logger)
				if err != nil {
					return err
				}
			}

			return cmd.Execute(ExecuteArgs{mod, fs, logger})
		}
	}
	return nil
}

func (r *cmdRunner) rootHandler(fs afero.Fs, logger *logrus.Logger) (err error) {

	if ok, err := afero.Exists(fs, r.module); !ok {
		return errors.New("Sysl module not found")
	} else if err != nil {
		return err
	}

	syslRootPath, syslRootExists, err := r.findRootFromASyslModule(fs, logger)
	if err != nil {
		return err
	}

	if syslRootExists && r.Root == syslRootPath {
		return
	}

	if r.rootIsDefined() && syslRootExists {
		logger.WithFields(logrus.Fields{
			"root":      r.Root,
			"SYSL_ROOT": syslRootPath,
		}).Warningln(fmt.Sprint("root is defined even though .SYSL_ROOT exists"))
	} else if !r.rootIsDefined() && !syslRootExists {
		logger.Errorln("root is not defined and .SYSL_ROOT can not be found")
	} else if !r.rootIsDefined() && syslRootExists {
		// TODO: log this?
		r.Root = syslRootPath
	}

	return
}

func (r *cmdRunner) findRootFromASyslModule(fs afero.Fs, logger *logrus.Logger) (syslRootPath string, syslRootExists bool, err error) {
	absolutePath, err := filepath.Abs(r.module)
	if err != nil {
		return
	}
	// paths := make(chan string)

	// go func(){
	// 	path,
	// }
	syslRootPath, err = r.walkUpParentDirectory(fs, absolutePath)
	if (err != nil){
		return
	}
	if (syslRootPath == ""){
		return 
	}

	return
}

func (r *cmdRunner) walkUpParentDirectory(fs afero.Fs, syslModuleAbsolutePath string) (syslRootPath string, err error) {

	for syslModuleAbsolutePath != "."{
		if ok, err := afero.DirExists(fs, filepath.Join(filepath.Dir(syslModuleAbsolutePath), syslRootMarker)); ok{
			return syslModuleAbsolutePath, nil
		} else if err != nil {
			return "", err
		}

		syslModuleAbsolutePath = filepath.Dir(syslModuleAbsolutePath)
	}

	return "", nil
}

func (r *cmdRunner) walkDownDirectory(syslModuleAbsolutePath string) (syslRootPath string, err error){

	return "", nil
}

func (r *cmdRunner) rootIsDefined() bool { return r.Root != rootUndefined }

func (r *cmdRunner) Configure(app *kingpin.Application) error {
	// app.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)

	commands := []Command{
		&protobuf{},
		&intsCmd{},
		&datamodelCmd{},
		&codegenCmd{},
		&sequenceDiagramCmd{},
		&importSwaggerCmd{},
		&infoCmd{},
		&validateCmd{},
	}
	r.commands = map[string]Command{}

	app.Flag("root",
		"sysl root directory for input model file (default: .)").
		Default(rootUndefined).StringVar(&r.Root)
	// TODO: Show warning based on whether root is defined or not
	for _, cmd := range commands {
		c := cmd.Configure(app)
		// TODO: find root based on the sysl module
		if cmd.RequireSyslModule() {
			c.Arg("MODULE", "input files without .sysl extension and with leading /, eg: "+
				"/project_dir/my_models combine with --root if needed").
				Required().StringVar(&r.module)
		}
		r.commands[cmd.Name()] = cmd
	}

	return nil
}

// Helper function to validate that a set of command flags are not empty values
func EnsureFlagsNonEmpty(cmd *kingpin.CmdClause, excludes ...string) {
	inExcludes := func(s string) bool {
		for _, e := range excludes {
			if s == e {
				return true
			}
		}
		return false
	}
	fn := func(c *kingpin.ParseContext) error {
		var errorMsg strings.Builder
		for _, f := range cmd.Model().Flags {
			if inExcludes(f.Name) {
				continue
			}
			val := f.Value.String()

			if val != "" {
				val = strings.Trim(val, " ")
				if val == "" {
					errorMsg.WriteString("'" + f.Name + "'" + " value passed is empty\n")
				}
			} else if len(f.Default) > 0 {
				errorMsg.WriteString("'" + f.Name + "'" + " value passed is empty\n")
			}
		}
		if errorMsg.Len() > 0 {
			return errors.New(errorMsg.String())
		}
		return nil
	}

	cmd.PreAction(fn)
}
