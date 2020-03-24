package main

import (
	"github.com/anz-bank/sysl/pkg/integrationdiagram"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"

	"gopkg.in/alecthomas/kingpin.v2"
)

type intsCmd struct {
	diagrams.Plantumlmixin
	cmdutils.CmdContextParamIntgen
}

func (p *intsCmd) Name() string       { return "integrations" }
func (p *intsCmd) MaxSyslModule() int { return 1 }

func (p *intsCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate integrations").Alias("ints")

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.Title)
	p.AddFlag(cmd)
	cmd.Flag("output",
		"output file(default: %(epname).png)").Default("%(epname).png").Short('o').StringVar(&p.Output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.Project)
	cmd.Flag("filter", "Only generate diagrams whose output paths match a pattern").StringVar(&p.Filter)
	cmd.Flag("exclude", "Apps to exclude").Short('e').StringsVar(&p.Exclude)
	cmd.Flag("clustered",
		"group integration components into clusters").Short('c').Default("false").BoolVar(&p.Clustered)
	cmd.Flag("epa", "produce and EPA integration view").Default("false").BoolVar(&p.EPA)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *intsCmd) Execute(args cmdutils.ExecuteArgs) error {
	result, err := integrationdiagram.GenerateIntegrations(&p.CmdContextParamIntgen, args.Modules[0], args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(result, args.Filesystem)
}
