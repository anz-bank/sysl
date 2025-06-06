package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/kingpin/v2"
	"github.com/anz-bank/sysl/pkg/arrai/transform"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/arr-ai/arrai/pkg/shell"
	"github.com/arr-ai/arrai/pkg/test"
	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/arrai/syntax"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type transformCmd struct {
	transformFile string
	outFile       string
	testFile      string
	debug         bool
}

func (p *transformCmd) Name() string       { return "transform" }
func (p *transformCmd) MaxSyslModule() int { return 99 }

func (p *transformCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Transform one or more Sysl models using an arr.ai script")
	cmd.Flag("script",
		"path of arr.ai script file with a transform function (accepts and returns specific structure, see docs)").
		Short('s').Required().StringVar(&p.transformFile)
	cmd.Flag("output",
		"path of file to write the formatted transform output, writes to stdout when not specified").
		Short('o').StringVar(&p.outFile)
	cmd.Flag("tests",
		"path of arr.ai script file with test function that accepts the transform output").
		Short('t').StringVar(&p.testFile)
	cmd.Flag("debug",
		"enable arr.ai shell debug mode when transform script fails").
		Short('d').BoolVar(&p.debug)
	return cmd
}

func (p *transformCmd) Execute(args cmdutils.ExecuteArgs) error {
	var err error
	var result rel.Value

	input, err := transform.BuildTransformInput(args.Modules, args.ModulePaths)
	if err != nil {
		return err
	}

	// Expands '~' since it's not automatically expanded by the shell like '$ENV' notations are.
	if strings.HasPrefix(p.transformFile, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		p.transformFile = homeDir + p.transformFile[1:]
	}

	var scriptBytes []byte

	exists, err := afero.Exists(args.Filesystem, p.transformFile)
	scriptPath := "temp.arrai"
	switch {
	case exists:
		scriptBytes, err = afero.ReadFile(args.Filesystem, p.transformFile)
		scriptPath = p.transformFile
	case p.transformFile == "-":
		scriptBytes, err = io.ReadAll(args.Stdin)
	case p.transformFile[0] == '\\':
		scriptBytes = []byte(p.transformFile)
	case regexp.MustCompile(`^[\w-]+\.[\w-]+/[\w-]+/[\w-]+/`).MatchString(p.transformFile):
		scriptBytes = []byte(fmt.Sprintf("//{%s}", p.transformFile))
	default:
		err = fmt.Errorf("the specified --script is neither '-' (for stdin), a local file, a remote file " +
			"(in the form of 'github.com/org/repo/path/to/file') or an inline script function")
	}

	if err == nil {
		result, err = transform.EvalWithParam(scriptBytes, scriptPath, input)
	}

	if err != nil {
		if p.debug {
			// Prints the error first and create the debugging session.
			// This returns nil so that it does not print the error messages twice
			logrus.Error(err)
			shell.CreateDebugSession(err)
			return nil
		}
		return err
	}

	if p.testFile == "" {
		return outputTransformResult(args.Filesystem, result, p.outFile)
	}
	return runTransformTests(args.Filesystem, result, p.testFile)
}

// outputTransformResult outputs the supplied transform result into the specified outfile or to stdout if none is
// specified, and sets the process exit code, and displayed a message.
func outputTransformResult(fs afero.Fs, transformResult rel.Value, outfile string) error {
	resultTuple, ok := transformResult.(rel.Tuple)
	if !ok {
		return fmt.Errorf("result of transform must be a tuple")
	}

	var outcome = map[string]interface{}{}
	if outcomeVal, found := resultTuple.Get("outcome"); found {
		outcome = outcomeVal.Export(context.Background()).(map[string]interface{})
	}

	var exitCode int
	var message string

	switch name := outcome["name"]; name {
	case "success", nil:
		exitCode = 0
		message = ""
	case "failure":
		exitCode = 1
		message = "transformation failed"
	case "warning":
		exitCode = 2
		message = "transformation completed with warning"
	default:
		return fmt.Errorf(
			"unrecognized outcome name, not one of the valid values ('success', 'failure' or 'warning'): %s", name)
	}

	if val, found := outcome["exitCode"]; found {
		exitCode = int(val.(float64))
	}
	if val, found := outcome["message"]; found {
		message = val.(string)
	}

	if output, found := resultTuple.Get("output"); found {
		var prettyResult string
		if _, ok := output.(rel.String); ok {
			prettyResult = output.String()
		} else {
			var err error
			prettyResult, err = syntax.PrettifyString(output, 0)
			if err != nil {
				return err
			}
		}

		if !strings.HasSuffix(prettyResult, "\n") {
			prettyResult += "\n"
		}

		if outfile == "" || outfile == "-" {
			fmt.Printf("%s", prettyResult)
		} else {
			return afero.WriteFile(fs, outfile, []byte(prettyResult), os.ModePerm)
		}
	}

	return syslutil.Exitf(exitCode, message)
}

// runTransformTests executes a transform test file against the results of a supplied transform result, and prints the
// test report. It also sets the exit code in the same way that the 'arrai test' command would.
func runTransformTests(fs afero.Fs, transformResult rel.Value, testFilePath string) error {
	scriptBytes, err := afero.ReadFile(fs, testFilePath)
	if err != nil {
		return err
	}

	testFile, err := transform.RunTests(scriptBytes, testFilePath, transformResult)
	if err != nil {
		return err
	}

	return test.Report(os.Stdout, []test.File{testFile})
}
