package main

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/ui"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var uiFields = `
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

type uiCmd struct {
	host   string
	fields string
	grpcui bool
}

func (p *uiCmd) Name() string       { return "ui" }
func (p *uiCmd) MaxSyslModule() int { return 1 }

func (p *uiCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Starts the Sysl UI which displays a visual view of Apps defined in Sysl.")
	cmd.Flag("host", "host and port to serve on").Default(":8080").Short('h').StringVar(&p.host)
	cmd.Flag("fields", "fields to display on the UI, separated by comma").Default(uiFields).Short('f').StringVar(&p.fields) //nolint:lll
	cmd.Flag("grpcui", "enables the grpcUI handlers").BoolVar(&p.grpcui)
	return cmd
}

func (p *uiCmd) Execute(args cmdutils.ExecuteArgs) error {
	args.Logger.Debugf("ui: %+v", *p)
	syslUI := ui.SyslUI{
		Host:    p.host,
		Fields:  strings.Split(p.fields, ","),
		Fs:      args.Filesystem,
		Log:     args.Logger,
		Modules: args.Modules,
		GRPCUI:  p.grpcui,
	}
	args.Logger.SetLevel(logrus.InfoLevel)
	server, err := syslUI.GenerateServer()
	if err != nil {
		return err
	}
	err = server.Setup()
	if err != nil {
		return err
	}
	return server.Serve()
}
