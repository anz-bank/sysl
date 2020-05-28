package main

import (
	"errors"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/mermaid/datamodeldiagram"
	"github.com/anz-bank/sysl/pkg/mermaid/endpointanalysisdiagram"
	"github.com/anz-bank/sysl/pkg/mermaid/integrationdiagram"
	"github.com/anz-bank/sysl/pkg/mermaid/sequencediagram"
	"github.com/anz-bank/sysl/pkg/sysl"
	"gopkg.in/alecthomas/kingpin.v2"
)

type diagramCmd cmdutils.DiagramCmd

func (p *diagramCmd) Name() string { return "diagram" }

func (p *diagramCmd) MaxSyslModule() int { return 1 }

func (p *diagramCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate mermaid diagrams").Alias("md")
	cmd.Flag("integrationdiagram",
		"Generate an integration diagram (Specify the application name)",
	).Short('i').StringVar(&p.IntegrationDiagram)
	cmd.Flag("sequencediagram",
		"Generate a sequence diagram (Specify 'appname->endpoint')",
	).Short('s').StringVar(&p.SequenceDiagram)
	cmd.Flag("endpointanalysis",
		"Generate an integration diagram with its endpoints (Specify 'true')",
	).Default("false").Short('e').BoolVar(&p.EndpointAnalysis)
	cmd.Flag("datadiagram",
		"Generate a Data model diagram (Specify 'true')",
	).Default("false").Short('d').BoolVar(&p.DataDiagram)
	return cmd
}

func (p *diagramCmd) Execute(args cmdutils.ExecuteArgs) error {
	out, err := callDiagramGenerator(args.Modules[0], p)
	if err != nil {
		return err
	}
	print(out)
	return nil
}

func callDiagramGenerator(m *sysl.Module, p *diagramCmd) (string, error) {
	switch {
	case p.IntegrationDiagram != "":
		return integrationdiagram.GenerateIntegrationDiagram(m, p.IntegrationDiagram)
	case p.SequenceDiagram != "":
		res := strings.Split(p.SequenceDiagram, "->")
		return sequencediagram.GenerateSequenceDiagram(m, res[0], res[1])
	case p.EndpointAnalysis:
		return endpointanalysisdiagram.GenerateEndpointAnalysisDiagram(m)
	case p.DataDiagram:
		return datamodeldiagram.GenerateFullDiagram(m)
	default:
		return "", errors.New("correct value has not been specified; please check help for more information")
	}
}
