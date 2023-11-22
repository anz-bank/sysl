package main

import (
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/parse"
	"gopkg.in/alecthomas/kingpin.v2"
)

type displaySummaryCmd struct{}

func (p *displaySummaryCmd) Name() string       { return "display-summary" }
func (p *displaySummaryCmd) MaxSyslModule() int { return 1 }

func (p *displaySummaryCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Display summary without parsing (currently just includes)")

	return cmd
}

func (p *displaySummaryCmd) PreExecute(settings *parse.Settings) error {
	settings.OperationSummary = true
	settings.NoParsing = true

	return nil
}

func (p *displaySummaryCmd) Execute(_ cmdutils.ExecuteArgs) error {
	// Nothing to do here, the runner loads the sysl file automatically. If we got here the file was successfully loaded
	return nil
}
