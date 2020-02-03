package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

// The sysl command consults environment variables for configuration.
// If an environment variable is unset, the sysl command uses a sensible
// default setting.
// 	SYSL_PLANTUML
// 		URL of PlantUML server. Sysl depends upon
// 		[PlantUML](http://plantuml.com/) for diagram generation.
// 	SYSL_MODULES
// 		Whether the sysl modules is enabled.
// 		Enable by default, set to "off" to disable sysl modules.
const KnownEnv = `
	SYSL_MODULES
	SYSL_PLANTUML
`

type envCmd struct{}

func (c *envCmd) Name() string       { return "env" }
func (c *envCmd) MaxSyslModule() int { return 0 }

func (c *envCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	return app.Command(c.Name(), "Print sysl environment information.")
}

func (c *envCmd) Execute(args ExecuteArgs) error {
	for _, e := range strings.Fields(KnownEnv) {
		v := os.Getenv(e)
		fmt.Printf("%s=\"%s\"\n", e, v)
	}

	return nil
}
