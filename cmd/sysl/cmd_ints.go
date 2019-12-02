package main

import (
	"regexp"

	sysl "github.com/anz-bank/sysl/pkg/proto_old"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const endpointWildcard = ".. * <- *"

func GenerateIntegrations(intgenParams *CmdContextParamIntgen,
	model *sysl.Module,
	logger *logrus.Logger) (map[string]string, error) {
	r := make(map[string]string)

	logger.Debugf("project: %v\n", intgenParams.project)
	logger.Debugf("clustered: %t\n", intgenParams.clustered)
	logger.Debugf("exclude: %s\n", intgenParams.exclude)
	logger.Debugf("epa: %t\n", intgenParams.epa)
	logger.Debugf("title: %s\n", intgenParams.title)
	logger.Debugf("filter: %s\n", intgenParams.filter)
	logger.Debugf("output: %s\n", intgenParams.output)

	if len(intgenParams.exclude) == 0 && intgenParams.project != "" {
		intgenParams.exclude = []string{intgenParams.project}
	}
	excludeStrSet := syslutil.MakeStrSet(intgenParams.exclude...)

	// The "project" app that specifies the required view of the integration
	app := model.GetApps()[intgenParams.project]
	of := MakeFormatParser(intgenParams.output)
	// Iterate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		outputDir := of.FmtOutput(intgenParams.project, epname, endpt.GetLongName(), endpt.GetAttrs())
		if intgenParams.filter != "" {
			re := regexp.MustCompile(intgenParams.filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		excludes := syslutil.MakeStrSetFromAttr("exclude", endpt.GetAttrs())
		passthroughs := syslutil.MakeStrSetFromAttr("passthrough", endpt.GetAttrs())
		b := makeBuilderfromStmt(model, endpt.GetStmt(), excludeStrSet.Union(excludes), passthroughs)
		intsParam := &IntsParam{b.finalApps, b.seedAppsMap, b.depsOut, app, endpt}
		args := &Args{intgenParams.title, intgenParams.project, intgenParams.clustered, intgenParams.epa}
		r[outputDir] = GenerateView(args, intsParam, model)
	}

	return r, nil
}

type intsCmd struct {
	plantumlmixin
	CmdContextParamIntgen
}

func (p *intsCmd) Name() string            { return "integrations" }
func (p *intsCmd) RequireSyslModule() bool { return true }

func (p *intsCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate integrations").Alias("ints")

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)
	p.AddFlag(cmd)
	cmd.Flag("output",
		"output file(default: %(epname).png)").Default("%(epname).png").Short('o').StringVar(&p.output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.project)
	cmd.Flag("filter", "Only generate diagrams whose output paths match a pattern").StringVar(&p.filter)
	cmd.Flag("exclude", "apps to exclude").Short('e').StringsVar(&p.exclude)
	cmd.Flag("clustered",
		"group integration components into clusters").Short('c').Default("false").BoolVar(&p.clustered)
	cmd.Flag("epa", "produce and EPA integration view").Default("false").BoolVar(&p.epa)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *intsCmd) Execute(args ExecuteArgs) error {
	result, err := GenerateIntegrations(&p.CmdContextParamIntgen, args.Module, args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(result, args.Filesystem)
}
