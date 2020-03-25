package main

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/anz-bank/sysl/pkg/eval"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/validate"
	"gopkg.in/alecthomas/kingpin.v2"
)

type codegenCmd struct {
	cmdutils.CmdContextParamCodegen
	outDir         string
	appName        string
	validateOnly   bool
	enableDebugger bool
}

func (p *codegenCmd) Name() string       { return "codegen" }
func (p *codegenCmd) MaxSyslModule() int { return 1 }

func (p *codegenCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate code").Alias("gen")
	cmd.Flag("root-transform",
		"sysl root directory for input transform file (default: .)").
		Default(".").StringVar(&p.RootTransform)
	cmd.Flag("transform", "path to transform file from the root transform directory").Required().StringVar(&p.Transform)
	cmd.Flag("grammar", "path to grammar file").Required().StringVar(&p.Grammar)
	cmd.Flag("app-name",
		"name of the sysl App defined in the sysl model."+
			" if there are multiple Apps defined in the sysl model,"+
			" code will be generated only for the given app").Default("").StringVar(&p.appName)
	cmd.Flag("start", "start rule for the grammar").Default(".").StringVar(&p.Start)
	cmd.Flag("outdir", "output directory").Default(".").StringVar(&p.outDir)
	cmd.Flag("dep-path", "path passed to sysl transform").Default("").StringVar(&p.DepPath)
	cmd.Flag("basepath", "base path for ReST output").Default("").StringVar(&p.BasePath)
	cmd.Flag("validate-only", "Only Perform validation on the transform grammar").BoolVar(&p.validateOnly)
	cmd.Flag("disable-validator", "Disable validation on the transform grammar").
		Default("false").BoolVar(&p.DisableValidator)
	cmd.Flag("debugger", "Enable the evaluation debugger on error").Default("false").BoolVar(&p.enableDebugger)
	EnsureFlagsNonEmpty(cmd, "app-name", "basepath", "dep-path")
	return cmd
}

func (p *codegenCmd) Execute(args cmdutils.ExecuteArgs) error {
	if p.validateOnly {
		return validate.DoValidate(validate.Params{
			RootTransform: p.RootTransform,
			Transform:     p.Transform,
			Grammar:       p.Grammar,
			Start:         p.Start,
			DepPath:       p.DepPath,
			BasePath:      p.BasePath,
			ParserType:    args.ParserType,
			Filesystem:    args.Filesystem,
			Logger:        args.Logger,
		})
	}
	if p.appName == "" {
		if len(args.Modules[0].Apps) > 1 {
			args.Logger.Errorf("required argument --app-name value missing")
			return fmt.Errorf("missing required argument")
		}
		p.appName = args.DefaultAppName
	}
	eval.EnableDebugger = p.enableDebugger
	output, err := GenerateCode(&p.CmdContextParamCodegen, args.Modules[0], p.appName, args.Filesystem, args.Logger, args.ParserType)
	if err != nil {
		return err
	}
	return outputToFiles(output, syslutil.NewChrootFs(args.Filesystem, p.outDir))
}
