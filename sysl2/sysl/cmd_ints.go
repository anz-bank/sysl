package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const endpointWildcard = ".. * <- *"

func GenerateIntegrations(intgenParams *CmdContextParamIntgen) (map[string]string, error) {
	r := make(map[string]string)

	log.Debugf("project: %v\n", *intgenParams.project)
	log.Debugf("clustered: %t\n", *intgenParams.clustered)
	log.Debugf("exclude: %s\n", *intgenParams.exclude)
	log.Debugf("epa: %t\n", *intgenParams.epa)
	log.Debugf("title: %s\n", *intgenParams.title)
	log.Debugf("filter: %s\n", *intgenParams.filter)
	log.Debugf("output: %s\n", *intgenParams.output)

	if *intgenParams.plantuml == "" {
		plantuml := os.Getenv("SYSL_PLANTUML")
		intgenParams.plantuml = &plantuml
		if *intgenParams.plantuml == "" {
			*intgenParams.plantuml = localPlantuml
		}
	}
	log.Debugf("plantuml: %s\n", *intgenParams.plantuml)


	mod := intgenParams.model

	if len(*intgenParams.exclude) == 0 && *intgenParams.project != "" {
		*intgenParams.exclude = []string{*intgenParams.project}
	}
	excludeStrSet := syslutil.MakeStrSet(*intgenParams.exclude...)

	// The "project" app that specifies the required view of the integration
	app := mod.GetApps()[*intgenParams.project]
	of := MakeFormatParser(*intgenParams.output)
	// Iterate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		outputDir := of.FmtOutput(*intgenParams.project, epname, endpt.GetLongName(), endpt.GetAttrs())
		if *intgenParams.filter != "" {
			re := regexp.MustCompile(*intgenParams.filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		excludes := syslutil.MakeStrSetFromAttr("exclude", endpt.GetAttrs())
		passthroughs := syslutil.MakeStrSetFromAttr("passthrough", endpt.GetAttrs())
		b := makeBuilderfromStmt(mod, endpt.GetStmt(), excludeStrSet.Union(excludes), passthroughs)
		intsParam := &IntsParam{b.finalApps, b.seedAppsMap, b.depsOut, app, endpt}
		args := &Args{*intgenParams.title, *intgenParams.project, *intgenParams.clustered, *intgenParams.epa}
		r[outputDir] = GenerateView(args, intsParam, mod)
	}

	return r, nil
}

type intsCmd struct {
	title     string
	output    string
	project   string
	filter    string
	exclude   []string
	clustered bool
	epa       bool
	plantuml  string
}

func (p *intsCmd) Name() string { return "ints" }
func (p *intsCmd) RequireSyslModule() bool { return true }

func (p *intsCmd) Init(app *kingpin.Application) *kingpin.CmdClause {

	cmd := app.Command(p.Name(), "Generate integrations")

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)
	cmd.Flag("plantuml", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n")).Short('p').StringVar(&p.plantuml)
	cmd.Flag("output",
		"output file(default: %(epname).png)").Default("%(epname).png").Short('o').StringVar(&p.output)
	cmd.Flag("project", "project pseudo-app to render").Short('j').StringVar(&p.project)
	cmd.Flag("filter", "Only generate diagrams whose output paths match a pattern").StringVar(&p.filter)
	cmd.Flag("exclude", "apps to exclude").Short('e').StringsVar(&p.exclude)
	cmd.Flag("clustered",
		"group integration components into clusters").Short('c').Default("false").BoolVar(&p.clustered)
	cmd.Flag("epa", "produce and EPA integration view").Default("false").BoolVar(&p.epa)

	return cmd
}

func (p *intsCmd) Execute(args ExecuteArgs) error {

	intgenParams := &CmdContextParamIntgen{
		model:          args.module,
		modelAppName:   args.modAppName,
		title:     &p.title,
		output:    &p.output,
		project:   &p.project,
		filter:    &p.filter,
		exclude:   &p.exclude,
		clustered: &p.clustered,
		epa:       &p.epa,
		plantuml:  &p.plantuml,
	}
	r, err := GenerateIntegrations(intgenParams)
	if err != nil {
		return err
	}
	for k, v := range r {
		err := OutputPlantuml(k, *intgenParams.plantuml, v, args.fs)
		if err != nil {
			return fmt.Errorf("plantuml drawing error: %v", err)
		}
	}
	return nil
}
