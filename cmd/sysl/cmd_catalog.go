package main

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/catalog"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var catalogFields = `team,
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
type,
confluence.url`

type catalogCmd struct {
	host   string
	fields string
}

func (p *catalogCmd) Name() string       { return "catalog" }
func (p *catalogCmd) MaxSyslModule() int { return 1 }

func (p *catalogCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate ui catalog of applications and endpoints")
	cmd.Flag("host", "host to serve on").Default(":8080").Short('h').StringVar(&p.host)
	cmd.Flag("fields", "fields to display on the UI, seperated by comma").Default(catalogFields).Short('f').StringVar(&p.fields)
	return cmd
}

func (p *catalogCmd) Execute(args ExecuteArgs) error {
	args.Logger.Debugf("catalog: %+v", *p)
	catalogServer := catalog.Server{
		Host:    p.host,
		Fields:  strings.Split(p.fields, ","),
		Fs:      args.Filesystem,
		Log:     args.Logger,
		Modules: args.Modules,
		Port:    ":8080",
	}
	args.Logger.SetLevel(logrus.InfoLevel)
	err := catalogServer.Serve()
	args.Logger.Info(err)
	return err
}
