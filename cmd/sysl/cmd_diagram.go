package main

import (
	"errors"
	"io/ioutil"

	"github.com/anz-bank/mermaid-go/mermaid"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/mermaid/datamodeldiagram"
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
	cmd.Flag("integration",
		"Generate an integration diagram (Optional- specify --app)",
	).Default("false").Short('i').BoolVar(&p.IntegrationDiagram)
	cmd.Flag("sequence",
		"Generate a sequence diagram (Specify --app and --endpoint)",
	).Default("false").Short('s').BoolVar(&p.SequenceDiagram)
	// cmd.Flag("endpointanalysis",
	// 	"Generate an integration diagram with its endpoints (Optional- specify --app)",
	// ).Default("false").Short('p').BoolVar(&p.EndpointAnalysis)
	cmd.Flag("data", "Generate a Data model diagram").Default("false").Short('d').BoolVar(&p.DataDiagram)
	cmd.Flag("app", "Optional flag to specify specific application").Short('a').StringVar(&p.AppName)
	cmd.Flag("endpoint", "Optional flag to specify endpoint").Short('e').StringVar(&p.Endpoint)
	cmd.Flag("output",
		"Output file (Default: diagram.svg)",
	).Default("diagram.svg").Short('o').StringVar(&p.Output)
	return cmd
}

func (p *diagramCmd) Execute(args cmdutils.ExecuteArgs) error {
	out, err := callDiagramGenerator(args.Modules[0], p)
	if err != nil {
		return err
	}
	g := mermaid.Init()
	svg := g.Execute(out)
	if err := ioutil.WriteFile(p.Output, []byte(svg), 0644); err != nil {
		panic(err)
	}
	return nil
}

func callDiagramGenerator(m *sysl.Module, p *diagramCmd) (string, error) {
	switch {
	case p.IntegrationDiagram:
		if p.AppName != "" {
			return integrationdiagram.GenerateIntegrationDiagram(m, p.AppName)
		}
		return integrationdiagram.GenerateFullIntegrationDiagram(m)
	case p.SequenceDiagram:
		if p.AppName != "" && p.Endpoint != "" {
			return sequencediagram.GenerateSequenceDiagram(m, p.AppName, p.Endpoint)
		}
		return "", errors.New("empty appname/endpoint; please check help for more information")
	// case p.EndpointAnalysis:
	// 	return endpointanalysisdiagram.GenerateEndpointAnalysisDiagram(m)
	case p.DataDiagram:
		return datamodeldiagram.GenerateFullDataDiagram(m)
	default:
		return "", errors.New("correct value has not been specified; please check help for more information")
	}
}
