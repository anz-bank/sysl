package main

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/codegen"

	"github.com/anz-bank/sysl/pkg/eval"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/validate"
	"gopkg.in/alecthomas/kingpin.v2"
)

type codegenCmd struct {
	cmdutils.CmdContextParamCodegen
	outDir         string
	validateOnly   bool
	enableDebugger bool
}

func (p *codegenCmd) Name() string       { return "codegen" }
func (p *codegenCmd) MaxSyslModule() int { return 1 }

func (p *codegenCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	desc := "Generate code, it has 3 required flags config, grammar and transform. If it doesn't set config, "
	desc += "grammar and transform must be set."
	cmd := app.Command(p.Name(), desc).Alias("gen")

	cmd.Flag("config",
		"config file path to set flags, it can set all runtime flags in the config file").StringVar(&p.Config)
	cmd.Flag("root-transform",
		"sysl root directory for input transform file (default: .)").
		Default(".").StringVar(&p.RootTransform)
	cmd.Flag("transform", "path to transform file from the root transform directory").StringVar(&p.Transform)
	cmd.Flag("grammar", "path to grammar file").StringVar(&p.Grammar)

	cmd.Flag("app-name",
		"name of the sysl app defined in sysl model."+
			" if there are multiple apps defined in sysl model,"+
			" code will be generated only for the given app").Default("").StringVar(&p.AppName)
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
	err := p.loadFlags()
	if err != nil {
		return err
	}
	if p.validateOnly {
		return validate.DoValidate(validate.Params{
			RootTransform: p.RootTransform,
			Transform:     p.Transform,
			Grammar:       p.Grammar,
			Start:         p.Start,
			DepPath:       p.DepPath,
			BasePath:      p.BasePath,
			Filesystem:    args.Filesystem,
			Logger:        args.Logger,
		})
	}
	if p.AppName == "" {
		if len(args.Modules[0].Apps) > 1 {
			args.Logger.Errorf("required argument --app-name value missing")
			return fmt.Errorf("missing required argument")
		}
		p.AppName = args.DefaultAppName
	}
	eval.EnableDebugger = p.enableDebugger
	output, err := GenerateCode(&p.CmdContextParamCodegen, args.Modules[0], p.AppName, args.Filesystem, args.Logger)
	if err != nil {
		return err
	}
	return outputToFiles(output, syslutil.NewChrootFs(args.Filesystem, p.outDir))
}

func (p *codegenCmd) loadFlags() error {
	err := validate.CodegenRequiredFlags(p.Config, p.Grammar, p.Transform)
	if err != nil {
		return err
	}

	if p.Config != "" {
		config, err := codegen.ReadCodeGenFlags(p.Config)
		if err != nil {
			return fmt.Errorf("failed to read config file %s", p.Config)
		}

		p.Transform = syslutil.GetNonEmpty(p.Transform, config.Transform)
		p.Grammar = syslutil.GetNonEmpty(p.Grammar, config.Grammar)
		p.DepPath = syslutil.GetNonEmpty(p.DepPath, config.DepPath)
		p.BasePath = syslutil.GetNonEmpty(p.BasePath, config.BasePath)
		p.AppName = syslutil.GetNonEmpty(p.AppName, config.AppName)
	}

	return nil
}
