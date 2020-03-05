package main

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/catalog"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var catalogFields = `
team,
team.slack,
owner.name,
owner.email,
file.version,
release.version,
description,
deploy.env1.url,
deploy.sit1.url,
deploy.sit2.url,
deploy.qa.url,
deploy.prod.url,
repo.url,
docs.url,
type`

type catalogCmd struct {
	host   string
	fields string
}

func (p *catalogCmd) Name() string       { return "catalog" }
func (p *catalogCmd) MaxSyslModule() int { return 1 }

func (p *catalogCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Starts the Sysl UI which visually presents your applications.")
	cmd.Flag("host", "host and port to serve on").Default(":8080").Short('h').StringVar(&p.host)
	cmd.Flag("fields", "fields to display on the UI, separated by comma").Default(catalogFields).Short('f').StringVar(&p.fields) //nolint:lll
	return cmd
}

func (p *catalogCmd) Execute(args cmdutils.ExecuteArgs) error {
	args.Logger.Debugf("catalog: %+v", *p)
	syslCatalog := catalog.SyslUI{
		Host:    p.host,
		Fields:  strings.Split(p.fields, ","),
		Fs:      args.Filesystem,
		Log:     args.Logger,
		Modules: args.Modules,
	}
	args.Logger.SetLevel(logrus.InfoLevel)
	server, err := syslCatalog.GenerateServer()
	if err != nil {
		return err
	}
	err = server.Setup()
	if err != nil {
		return err
	}
	return server.Serve()
}
