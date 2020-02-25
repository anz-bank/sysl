package main

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/config"
	"github.com/anz-bank/sysl/pkg/eval"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/validate"
	"gopkg.in/alecthomas/kingpin.v2"
)

type codegenCmd struct {
	CmdContextParamCodegen
	outDir         string
	appName        string
	validateOnly   bool
	enableDebugger bool
	config         string
}

func (p *codegenCmd) Name() string       { return "codegen" }
func (p *codegenCmd) MaxSyslModule() int { return 1 }

func (p *codegenCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	desc := "Generate code, it has 3 required flags config, grammar and transform. If it doesn't set config, "
	desc += "grammar and transform must be set."
	cmd := app.Command(p.Name(), desc).Alias("gen")
	//
	cmd.Flag("config",
		"config file path to set flags, it can set all runtime flags in the config file").StringVar(&p.config)
	cmd.Flag("grammar", "path to grammar file").StringVar(&p.grammar)
	cmd.Flag("transform", "path to transform file from the root transform directory").StringVar(&p.transform)
	//

	cmd.Flag("root-transform",
		"sysl root directory for input transform file (default: .)").
		Default(".").StringVar(&p.rootTransform)
	cmd.Flag("app-name",
		"name of the sysl app defined in sysl model."+
			" if there are multiple apps defined in sysl model,"+
			" code will be generated only for the given app").Default("").StringVar(&p.appName)
	cmd.Flag("start", "start rule for the grammar").Default(".").StringVar(&p.start)
	cmd.Flag("outdir", "output directory").Default(".").StringVar(&p.outDir)
	cmd.Flag("dep-path", "path passed to sysl transform").Default("").StringVar(&p.depPath)
	cmd.Flag("basepath", "base path for ReST output").Default("").StringVar(&p.basePath)
	cmd.Flag("validate-only", "Only Perform validation on the transform grammar").BoolVar(&p.validateOnly)
	cmd.Flag("disable-validator", "Disable validation on the transform grammar").
		Default("false").BoolVar(&p.disableValidator)
	cmd.Flag("debugger", "Enable the evaluation debugger on error").Default("false").BoolVar(&p.enableDebugger)
	EnsureFlagsNonEmpty(cmd, "app-name", "basepath", "dep-path")
	return cmd
}

func (p *codegenCmd) Execute(args ExecuteArgs) error {
	err := p.loadFlags()
	if err != nil {
		return err
	}

	if p.validateOnly {
		return validate.DoValidate(validate.Params{
			RootTransform: p.rootTransform,
			Transform:     p.transform,
			Grammar:       p.grammar,
			Start:         p.start,
			DepPath:       p.depPath,
			BasePath:      p.basePath,
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
	output, err := GenerateCode(&p.CmdContextParamCodegen, args.Modules[0], p.appName, args.Filesystem, args.Logger)
	if err != nil {
		return err
	}
	return outputToFiles(output, syslutil.NewChrootFs(args.Filesystem, p.outDir))
}

func (p *codegenCmd) loadFlags() error {
	err := validate.CodeggenRequiredFlags(p.config, p.grammar, p.transform)
	if err != nil {
		return err
	}

	if p.config != "" {
		config, err := config.ReadCodeGenFlags(p.config)
		if err != nil {
			return fmt.Errorf("failed to read config file %s", p.config)
		}

		p.transform = syslutil.ResetVal(p.transform, config.Transform)
		p.grammar = syslutil.ResetVal(p.grammar, config.Grammar)
		p.depPath = syslutil.ResetVal(p.depPath, config.DepPath)
		p.basePath = syslutil.ResetVal(p.basePath, config.BasePath)
		p.appName = syslutil.ResetVal(p.appName, config.AppName)
	}

	return nil
}
