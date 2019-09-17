package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	PlantUMLEnvVar  = "SYSL_PLANTUML"
	PlantUMLDefault = "http://localhost:8080/plantuml"
)

type plantumlmixin struct {
	value string
}

func (p *plantumlmixin) AddFlag(cmd *kingpin.CmdClause) {
	cmd.Flag("plantuml",
		"base url of plantuml server (default: "+PlantUMLEnvVar+" or "+
			PlantUMLDefault+" see "+
			"http://plantuml.com/server.html#install for more info)",
	).Short('p').StringVar(&p.value)
}

func (p *plantumlmixin) Value() string {
	if p.value == "" {
		p.value = os.Getenv(PlantUMLEnvVar)
		if p.value == "" {
			p.value = PlantUMLDefault
		}
	}
	return p.value
}
