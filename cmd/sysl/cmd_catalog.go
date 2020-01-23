package main

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/catalog"
	"gopkg.in/alecthomas/kingpin.v2"
)

type catalogCmd struct {
	output string
	mode   string
}

func (p *catalogCmd) Name() string       { return "catalog" }
func (p *catalogCmd) MaxSyslModule() int { return 1 }

func (p *catalogCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate ui catalog of applicatons and endpoints")
	cmd.Flag("output", "output file name").Short('o').Default("-").StringVar(&p.output)
	cmd.Flag("format", "format of server, supported options; html | json").Default("html").Short('f').StringVar(&p.mode)
	return cmd
}

func (p *catalogCmd) Execute(args ExecuteArgs) error {
	args.Logger.Debugf("catalog: %+v", *p)
	p.output = strings.TrimSpace(p.output)
	p.mode = strings.TrimSpace(p.mode)
	catalogServer := catalog.CatalogServer{Port: 8080, Host: "localhost", Type: p.mode, Fs: args.Filesystem, Log: args.Logger, Module: args.Modules[0]}

	return catalogServer.Serve()
}
