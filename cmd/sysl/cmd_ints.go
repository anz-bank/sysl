package main

import (
	"regexp"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/diagrams"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const endpointWildcard = ".. * <- *"

func GenerateIntegrations(intgenParams *cmdutils.CmdContextParamIntgen,
	model *sysl.Module,
	logger *logrus.Logger) (map[string]string, error) {
	r := make(map[string]string)

	logger.Debugf("project: %v\n", intgenParams.Project)
	logger.Debugf("clustered: %t\n", intgenParams.Clustered)
	logger.Debugf("exclude: %s\n", intgenParams.Exclude)
	logger.Debugf("epa: %t\n", intgenParams.EPA)
	logger.Debugf("title: %s\n", intgenParams.Title)
	logger.Debugf("filter: %s\n", intgenParams.Filter)
	logger.Debugf("output: %s\n", intgenParams.Output)

	if len(intgenParams.Exclude) == 0 && intgenParams.Project != "" {
		intgenParams.Exclude = []string{intgenParams.Project}
	}
	excludeStrSet := syslutil.MakeStrSet(intgenParams.Exclude...)

	// The "project" app that specifies the required view of the integration
	app := model.GetApps()[intgenParams.Project]
	of := cmdutils.MakeFormatParser(intgenParams.Output)
	// Iterate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		outputDir := of.FmtOutput(intgenParams.Project, epname, endpt.GetLongName(), endpt.GetAttrs())
		if intgenParams.Filter != "" {
			re := regexp.MustCompile(intgenParams.Filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		excludes := syslutil.MakeStrSetFromAttr("exclude", endpt.GetAttrs())
		passthroughs := syslutil.MakeStrSetFromAttr("passthrough", endpt.GetAttrs())
		b := makeBuilderfromStmt(model, endpt.GetStmt(), excludeStrSet.Union(excludes), passthroughs)
		intsParam := &IntsParam{b.finalApps, b.seedAppsMap, b.depsOut, app, endpt}
		args := &Args{intgenParams.Title, intgenParams.Project, intgenParams.Clustered, intgenParams.EPA}
		r[outputDir] = GenerateView(args, intsParam, model)
	}

	return r, nil
}

type intsCmd struct {
	diagrams.Plantumlmixin
	cmdutils.CmdContextParamIntgen
}

func (p *intsCmd) Name() string       { return "integrations" }
func (p *intsCmd) MaxSyslModule() int { return 1 }

func (p *intsCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate integrations").Alias("ints")

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.Title)
	p.AddFlag(cmd)
	cmd.Flag("output",
		"output file(default: %(epname).png)").Default("%(epname).png").Short('o').StringVar(&p.Output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.Project)
	cmd.Flag("filter", "Only generate diagrams whose output paths match a pattern").StringVar(&p.Filter)
	cmd.Flag("exclude", "apps to exclude").Short('e').StringsVar(&p.Exclude)
	cmd.Flag("clustered",
		"group integration components into clusters").Short('c').Default("false").BoolVar(&p.Clustered)
	cmd.Flag("epa", "produce and EPA integration view").Default("false").BoolVar(&p.EPA)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *intsCmd) Execute(args cmdutils.ExecuteArgs) error {
	result, err := GenerateIntegrations(&p.CmdContextParamIntgen, args.Modules[0], args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(result, args.Filesystem)
}
