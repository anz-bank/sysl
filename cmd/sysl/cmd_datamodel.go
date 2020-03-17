package main

import (
	"github.com/anz-bank/sysl/pkg/datamodeldiagram"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"

	"gopkg.in/alecthomas/kingpin.v2"
)

type datamodelCmd struct {
	diagrams.Plantumlmixin
	cmdutils.CmdContextParamDatagen
}

func (p *datamodelCmd) Name() string       { return "datamodel" }
func (p *datamodelCmd) MaxSyslModule() int { return 1 }

func (p *datamodelCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models").Alias("data")
	cmd.Flag("class_format",
		"Specify the format string for data diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))",
	).Default("%(classname)").StringVar(&p.ClassFormat)

	cmd.Flag("title", "Diagram title").Short('t').StringVar(&p.Title)

	p.AddFlag(cmd)

	cmd.Flag("output",
		"Output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').StringVar(&p.Output)
	cmd.Flag("project", "Project pseudo-app to render").Short('j').StringVar(&p.Project)
	cmd.Flag("direct", "Process data model directly without project manner").Short('d').BoolVar(&p.Direct)
	cmd.Flag("filter", "Only generate diagrams whose names match a pattern").Short('f').StringVar(&p.Filter)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *datamodelCmd) Execute(args cmdutils.ExecuteArgs) error {
	outmap, err := datamodeldiagram.GenerateDataModels(&p.CmdContextParamDatagen, args.Modules[0], args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(outmap, args.Filesystem)
}
