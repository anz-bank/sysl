package main

import (
	"os"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const endpointWildcard = ".. * <- *"

func GenerateIntegrations(intgenParams *CmdContextParamIntgen) (map[string]string, error) {
	r := make(map[string]string)

	log.Debugf("root: %s\n", *intgenParams.root)
	log.Debugf("project: %v\n", *intgenParams.project)
	log.Debugf("clustered: %t\n", *intgenParams.clustered)
	log.Debugf("exclude: %s\n", *intgenParams.exclude)
	log.Debugf("epa: %t\n", *intgenParams.epa)
	log.Debugf("title: %s\n", *intgenParams.title)
	log.Debugf("filter: %s\n", *intgenParams.filter)
	log.Debugf("modules: %s\n", *intgenParams.modules)
	log.Debugf("output: %s\n", *intgenParams.output)

	if *intgenParams.plantuml == "" {
		plantuml := os.Getenv("SYSL_PLANTUML")
		intgenParams.plantuml = &plantuml
		if *intgenParams.plantuml == "" {
			*intgenParams.plantuml = localPlantuml
		}
	}
	log.Debugf("plantuml: %s\n", *intgenParams.plantuml)

	mod, err := loadApp(*intgenParams.modules, syslutil.NewChrootFs(afero.NewOsFs(), *intgenParams.root))
	if err != nil {
		return nil, err
	}

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

func configureCmdlineForIntgen(sysl *kingpin.Application, flagmap map[string][]string) *CmdContextParamIntgen {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	flagmap["ints"] = []string{"root", "title", "plantuml", "output", "project", "filter", "exclude", "epa"}
	ints := sysl.Command("ints", "Generate integrations")
	returnValues := &CmdContextParamIntgen{}

	returnValues.root = ints.Flag("root", "sysl root directory for input model file (default: .)").Default(".").String()
	returnValues.title = ints.Flag("title", "diagram title").Short('t').String()
	returnValues.plantuml = ints.Flag("plantuml", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n")).Short('p').String()
	returnValues.output = ints.Flag("output",
		"output file(default: %(epname).png)").Default("%(epname).png").Short('o').String()
	returnValues.project = ints.Flag("project", "project pseudo-app to render").Short('j').String()
	returnValues.filter = ints.Flag("filter", "Only generate diagrams whose output paths match a pattern").String()
	returnValues.modules = ints.Arg("modules",
		strings.Join([]string{"input files without .sysl extension and with leading /",
			"eg: /project_dir/my_models",
			"combine with --root if needed"}, "\n")).String()
	returnValues.exclude = ints.Flag("exclude", "apps to exclude").Short('e').Strings()
	returnValues.clustered = ints.Flag("clustered",
		"group integration components into clusters").Short('c').Default("false").Bool()
	returnValues.epa = ints.Flag("epa", "produce and EPA integration view").Default("false").Bool()

	return returnValues
}
