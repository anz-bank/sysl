package main

import (
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

type genCmd struct {
	rootTransform string
	transform     string
	grammar       string
	start         string
	outDir        string
}

func (p *genCmd) Name() string { return "gen" }
func (p *genCmd) RequireSyslModule() bool { return true }

func (p *genCmd) Init(app *kingpin.Application) *kingpin.CmdClause {

	cmd := app.Command(p.Name(), "Generate code").Alias("codegen")
	cmd.Flag("root-transform",
		"sysl root directory for input transform file (default: .)").
		Default(".").StringVar(&p.rootTransform)
	cmd.Flag("transform", "grammar.g").Required().StringVar(&p.transform)
	cmd.Flag("grammar", "grammar.g").Required().StringVar(&p.grammar)
	cmd.Flag("start", "start rule for the grammar").Default(".").StringVar(&p.start)
	cmd.Flag("outdir", "output directory").Default(".").StringVar(&p.outDir)
	return cmd
}

func (p *genCmd) Execute(args ExecuteArgs) error {

	codegenParams := &CmdContextParamCodegen{
		model:         args.module,
		modelAppName:  args.modAppName,
		rootTransform: &p.rootTransform,
		transform:     &p.transform,
		grammar:       &p.grammar,
		start:         &p.start,
		outDir:        &p.outDir,
	}
	output, err := GenerateCode(codegenParams, args.fs, args.logger)
	if err != nil {
		return err
	}
	return outputToFiles(output, syslutil.NewChrootFs(args.fs, *codegenParams.outDir))
}
