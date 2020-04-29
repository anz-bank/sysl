package main

import (
	"github.com/anz-bank/sysl/pkg/cmdutils"
	//"github.com/anz-bank/sysl/pkg/mermaid/sequencediagram"
	"github.com/anz-bank/sysl/pkg/mermaid/integrationdiagram"
	//"github.com/anz-bank/sysl/pkg/mermaid/endpointanalysisdiagram"
	"gopkg.in/alecthomas/kingpin.v2"
)

type diagramCmd struct {
	cmdutils.DiagramCmd
}

func (p *diagramCmd) Name() string { return "diagram" }

func (p *diagramCmd) MaxSyslModule() int { return 1 }

func (p *diagramCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate mermaid diagrams").Alias("md")
	cmd.Flag("integrationdiagram", "To generate the integration"+
		" diagram, specify the application name").Short('i').StringVar(&p.Value)
	cmd.Flag("endpointanalysisdiagram", "Generates the endpoint"+
		" analysis diagram for the given sysl file").Short('e').StringVar(&p.Value)
	cmd.Flag("sequencediagram", "To generate the sequence diagram, specify the "+
		"application name and endpoint in this format app->endpoint").Short('s').StringVar(&p.Value)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *diagramCmd) Execute(args cmdutils.ExecuteArgs) error {
	outString, err := integrationdiagram.GenerateIntegrationDiagram(args.Modules[0], p.Value)
	if err != nil {
		return err
	}
	print(outString)
	return nil
}
