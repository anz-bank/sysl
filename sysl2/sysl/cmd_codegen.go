package main

import (
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

type codegenCmd struct {
	CmdContextParamCodegen
	outDir string
}

func (p *codegenCmd) Name() string            { return "codegen" }
func (p *codegenCmd) RequireSyslModule() bool { return true }

func (p *codegenCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {

	cmd := app.Command(p.Name(), "Generate code").Alias("gen")
	cmd.Flag("root-transform",
		"sysl root directory for input transform file (default: .)").
		Default(".").StringVar(&p.rootTransform)
	cmd.Flag("transform", "grammar.g").Required().StringVar(&p.transform)
	cmd.Flag("grammar", "grammar.g").Required().StringVar(&p.grammar)
	cmd.Flag("start", "start rule for the grammar").Default(".").StringVar(&p.start)
	cmd.Flag("outdir", "output directory").Default(".").StringVar(&p.outDir)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *codegenCmd) Execute(args ExecuteArgs) error {

	output, err := GenerateCode(&p.CmdContextParamCodegen, args.Module, args.ModuleAppName, args.Filesystem, args.Logger)
	if err != nil {
		return err
	}
	return outputToFiles(output, syslutil.NewChrootFs(args.Filesystem, p.outDir))
}
