package main

import (
	"fmt"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/anz-bank/sysl/sysl2/sysl/validate"
	"gopkg.in/alecthomas/kingpin.v2"
)

type codegenCmd struct {
	CmdContextParamCodegen
	outDir       string
	appName      string
	validateOnly bool
}

func (p *codegenCmd) Name() string            { return "codegen" }
func (p *codegenCmd) RequireSyslModule() bool { return true }

func (p *codegenCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate code").Alias("gen")
	cmd.Flag("root-transform", "Deprecated").Default(".").Hidden().StringVar(&p.rootTransform)
	cmd.Flag("transform", "path to transform file from the root transform directory").Required().StringVar(&p.transform)
	cmd.Flag("grammar", "path to grammar file").Required().StringVar(&p.grammar)
	cmd.Flag("app-name",
		"name of the sysl app defined in sysl model."+
			" if there are multiple apps defined in sysl model,"+
			" code will be generated only for the given app").Default("").StringVar(&p.appName)
	cmd.Flag("start", "start rule for the grammar").Default(".").StringVar(&p.start)
	cmd.Flag("outdir", "output directory").Default(".").StringVar(&p.outDir)
	cmd.Flag("validate-only", "Only Perform validation on the transform grammar").BoolVar(&p.validateOnly)
	EnsureFlagsNonEmpty(cmd, "app-name")
	return cmd
}

func (p *codegenCmd) Execute(args ExecuteArgs) error {
	if p.validateOnly {
		return validate.DoValidate(validate.Params{
			RootTransform: p.rootTransform,
			Transform:     p.transform,
			Grammar:       p.grammar,
			Start:         p.start,
			Filesystem:    args.Filesystem,
			Logger:        args.Logger,
		})
	}
	if p.appName == "" {
		if len(args.Module.Apps) > 1 {
			args.Logger.Errorf("required argument --app-name value missing")
			return fmt.Errorf("missing required argument")
		}
		p.appName = args.DefaultAppName
	}
	output, err := GenerateCode(&p.CmdContextParamCodegen, args.Module, p.appName, args.Filesystem, args.Logger)
	if err != nil {
		return err
	}
	return outputToFiles(output, syslutil.NewChrootFs(args.Filesystem, p.outDir))
}
