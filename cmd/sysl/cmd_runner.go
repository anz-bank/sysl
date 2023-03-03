package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/pbutil"
	"github.com/anz-bank/sysl/pkg/sysl"

	"github.com/anz-bank/golden-retriever/pkg/gitfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

type cmdRunner struct {
	commands map[string]cmdutils.Command

	Root          string
	CloneVersion  string
	modules       []string
	parseSettings parse.Settings
}

// Run identifies the command to run, loads the Sysl modules from the input (if necessary), then
// executes the command with all of the accumulated context.
func (r *cmdRunner) Run(which string, fs afero.Fs, logger *logrus.Logger, stdin io.Reader) error {
	// splitter to parse main command from subcommand
	mainCommand := strings.Split(which, " ")[0]
	if cmd, ok := r.commands[mainCommand]; ok {
		if cmd.Name() == mainCommand {
			var modules []*sysl.Module
			var err error
			var gitRoot string

			preExecuter, ok := cmd.(cmdutils.PreExecuteCommand)
			if ok {
				err = preExecuter.PreExecute(&r.parseSettings)
				if err != nil {
					return err
				}
			}

			if r.CloneVersion != "" {
				fs, gitRoot, err = r.getClonedRepo(fs)
				if err != nil {
					return err
				}
			}

			if cmd.MaxSyslModule() > 0 {
				if len(r.modules) > 0 {
					if r.CloneVersion != "" {
						modules, err = r.loadFromClone(fs, gitRoot)
					} else {
						modules, err = r.loadFromModules(fs, logger)
					}
					// stdin may still be provided for use by commands like transform.
				} else {
					modules, err = r.loadFromStdin(stdin, fs, logger)
				}
				if err != nil {
					return err
				}
			}

			if len(modules) > cmd.MaxSyslModule() {
				return fmt.Errorf("this command can accept max %d module(s)", cmd.MaxSyslModule())
			}
			return cmd.Execute(cmdutils.ExecuteArgs{
				Command:        which,
				Modules:        modules,
				Filesystem:     fs,
				Logger:         logger,
				DefaultAppName: "",
				ModulePaths:    r.modules,
				Root:           r.Root,
				Stdin:          stdin,
			})
		}
	}
	return nil
}

func (r *cmdRunner) Configure(app *kingpin.Application) error {
	commands := []cmdutils.Command{
		&codegenCmd{},
		&databaseScriptCmd{},
		&datamodelCmd{},
		&diagramCmd{},
		&envCmd{},
		&exportCmd{},
		&importCmd{},
		&infoCmd{},
		&intsCmd{},
		&lspCmd{},
		&modDatabaseScriptCmd{},
		&protobufCmd{},
		&replCmd{},
		&sequenceDiagramCmd{},
		&templateCmd{},
		&testRigCmd{},
		&transformCmd{},
		&validateCmd{},
		&displaySummaryCmd{},
	}
	r.commands = map[string]cmdutils.Command{}

	app.Flag("root",
		"Sysl root directory for input model file. If root is not found, the module directory becomes "+
			"the root, but the module can not import with absolute paths (or imports must be relative).").StringVar(&r.Root)

	app.Flag("clone-version",
		"Before running the command it will clone the local repo into memory and checkout the specific version",
	).StringVar(&r.CloneVersion)

	app.Flag("max-import-depth",
		"Maximum depth to follow imports, including the original file (ignores any that are deeper)."+
			" 0 (default) for unlimited."+
			" eg 1 means just the original file with no imports.",
	).IntVar(&r.parseSettings.MaxImportDepth)

	app.Flag("operation-summary",
		"Currently just outputs the names of the files parsed.",
	).BoolVar(&r.parseSettings.OperationSummary)

	sort.Slice(commands, func(i, j int) bool {
		return strings.Compare(commands[i].Name(), commands[j].Name()) < 0
	})
	for _, cmd := range commands {
		r.ConfigureCmd(app, cmd)
	}

	return nil
}

func (r *cmdRunner) ConfigureCmd(app *kingpin.Application, cmd cmdutils.Command) {
	c := cmd.Configure(app)
	if cmd.MaxSyslModule() > 0 {
		c.Arg("MODULE", "input files without .sysl extension and with leading /, eg: "+
			"/project_dir/my_models combine with --root if needed").
			StringsVar(&r.modules)
	}
	r.commands[cmd.Name()] = cmd
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
		for _, elem := range c.Elements {
			if f, _ := elem.Clause.(*kingpin.FlagClause); f != nil && f.Model().Name == "help" {
				return nil // help requested, don't need to check for empty flags
			}
		}
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

// loadFromStdin attempts to load the Sysl modules for the files provided via stdin.
func (r *cmdRunner) loadFromStdin(stdin io.Reader, fs afero.Fs, logger *logrus.Logger) ([]*sysl.Module, error) {
	src, err := io.ReadAll(stdin)
	if err != nil {
		return nil, err
	}

	stdinFiles, filesErr := loadStdinFiles(src)
	if filesErr != nil {
		mod, pbErr := pbutil.FromPBByteContents("stdin.pb", src)
		if pbErr != nil {
			return nil, fmt.Errorf(
				"failed to parse module from stdin: content is not valid file list JSON or compiled pb\n├─ JSON: %v\n└─ pb: %v",
				filesErr, pbErr)
		}
		return []*sysl.Module{mod}, nil
	}

	fs = afero.NewCopyOnWriteFs(fs, afero.NewMemMapFs())
	for _, f := range stdinFiles {
		r.modules = append(r.modules, f.Path)
		absPath, err := filepath.Abs(f.Path)
		if err != nil {
			return nil, err
		}
		err = afero.WriteFile(fs, absPath, []byte(f.Content), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return r.loadFromModules(fs, logger)
}

// loadFromModules attempts to load the Sysl modules for the files specified in r.modules.
func (r *cmdRunner) loadFromModules(fs afero.Fs, logger *logrus.Logger) ([]*sysl.Module, error) {
	var mods []*sysl.Module
	for _, moduleName := range r.modules {
		mod, _, err := loader.LoadSyslModuleWithSettings(r.Root, moduleName, fs, logger, r.parseSettings)
		if err != nil {
			return nil, err
		}
		mods = append(mods, mod)
	}
	return mods, nil
}

// loadFromClone attempts to load the Sysl modules for the files specified in r.modules and assumes that fs is a
// cloned version of the local repo.
func (r *cmdRunner) loadFromClone(fs afero.Fs, gitRoot string) ([]*sysl.Module, error) {
	var mods []*sysl.Module
	for _, moduleName := range r.modules {
		var err error

		// make moduleName relative to the repo
		moduleName, err = filepath.Abs(moduleName)
		if err != nil {
			return nil, err
		}
		moduleName, err = filepath.Rel(gitRoot, moduleName)
		if err != nil {
			return nil, err
		}

		modelParser := parse.NewParser()
		modelParser.Set(r.parseSettings)
		mod, err := modelParser.ParseFromFs(moduleName, fs)
		if err != nil {
			return nil, err
		}

		mods = append(mods, mod)
	}

	return mods, nil
}

func (r *cmdRunner) getClonedRepo(fs afero.Fs) (afero.Fs, string, error) {
	gitRoot := r.Root
	if gitRoot == "" {
		var err error
		gitRoot, err = loader.FindRootFromSyslModule("model.sysl", fs, parse.GitRootMarker)
		if err != nil || gitRoot == "" {
			return nil, "", fmt.Errorf("couldn't find local repo to clone: %w", err)
		}
	}

	// clone to in-memory storage
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:          gitRoot,
		SingleBranch: false,
	})
	if err != nil {
		return nil, "", fmt.Errorf("clone failed: %w", err)
	}

	// get long git ref
	ref, err := repo.ResolveRevision(plumbing.Revision(r.CloneVersion))
	if err != nil {
		return nil, "", fmt.Errorf("repo.ResolveRevision failed: %w", err)
	}

	// get the commit
	commit, err := repo.CommitObject(*ref)
	if err != nil {
		return nil, "", fmt.Errorf("repo.CommitObject failed: %w", err)
	}

	return gitfs.NewGitMemFs(commit), gitRoot, nil
}
