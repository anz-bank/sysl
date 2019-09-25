package main

import (
	"github.com/anz-bank/sysl/sysl2/sysl/validate"
	"gopkg.in/alecthomas/kingpin.v2"
)

type validateCmd struct {
	rootTransform string
	transform     string
	grammar       string
	start         string
}

func (p *validateCmd) Name() string            { return "validate" }
func (p *validateCmd) RequireSyslModule() bool { return false }

func (p *validateCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {

	cmd := app.Command(p.Name(), "Validate transform")
	cmd.Flag("root-transform", "sysl root directory for input transform file (default: .)").
		Default(".").StringVar(&p.rootTransform)
	cmd.Flag("transform", "grammar.g").Required().StringVar(&p.transform)
	cmd.Flag("grammar", "grammar.sysl").Required().StringVar(&p.grammar)
	cmd.Flag("start", "start rule for the grammar").Default(".").StringVar(&p.start)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *validateCmd) Execute(args ExecuteArgs) error {

	return validate.DoValidate(validate.Params{
		RootTransform: p.rootTransform,
		Transform:     p.transform,
		Grammar:       p.grammar,
		Start:         p.start,
		Filesystem:    args.Filesystem,
		Logger:        args.Logger,
	})
}
