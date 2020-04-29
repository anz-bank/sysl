package main

import (
	"errors"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/mermaid/endpointanalysisdiagram"
	"github.com/anz-bank/sysl/pkg/mermaid/integrationdiagram"
	"github.com/anz-bank/sysl/pkg/mermaid/sequencediagram"
	"github.com/anz-bank/sysl/pkg/sysl"
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
		" diagram, specify the application name").Short('i').StringVar(&p.IntegrationValue)
	cmd.Flag("sequencediagram", "To generate the sequence diagram, specify the "+
		"application name and endpoint in this format app->endpoint").Short('s').StringVar(&p.SequenceValue)
	cmd.Flag("endpointanalysis", "To generates the endpoint analysis diagram"+
		" for the given sysl file, specify true followed by the flag").Short('e').BoolVar(&p.EndpointAnalysisValue)
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
	if p.IntegrationValue == "" {
		if p.SequenceValue != "" {
			res := strings.Split(p.SequenceValue, "->")
			return sequencediagram.GenerateSequenceDiagram(m, res[0], res[1])
		} else if p.EndpointAnalysisValue {
			return endpointanalysisdiagram.GenerateEndpointAnalysisDiagram(m)
		}
	} else {
		return integrationdiagram.GenerateIntegrationDiagram(m, p.IntegrationValue)
	}
	return "", errors.New("correct value has not been specified; please check help for more information")
}
