package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const endpointWildcard = ".. * <- *"

func GenerateIntegrations(
	rootModel, title, output, project, filter, modules string,
	exclude []string,
	clustered, epa bool,
) (map[string]string, error) {
	r := make(map[string]string)

	mod, err := loadApp(rootModel, modules)
	if err != nil {
		return nil, err
	}

	if len(exclude) == 0 && project != "" {
		exclude = []string{project}
	}
	excludeStrSet := MakeStrSet(exclude...)

	// The "project" app that specifies the required view of the integration
	app := mod.GetApps()[project]
	of := MakeFormatParser(output)
	// Iterate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		outputDir := of.FmtOutput(project, epname, endpt.GetLongName(), endpt.GetAttrs())
		if filter != "" {
			re := regexp.MustCompile(filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		excludes := MakeStrSetFromAttr("exclude", endpt.GetAttrs())
		passthroughs := MakeStrSetFromAttr("passthrough", endpt.GetAttrs())
		b := makeBuilderfromStmt(mod, endpt.GetStmt(), excludeStrSet.Union(excludes), passthroughs)
		intsParam := &IntsParam{b.finalApps, b.seedAppsMap, b.depsOut, app, endpt}
		args := &Args{title, project, clustered, epa}
		r[outputDir] = GenerateView(args, intsParam, mod)
	}

	return r, nil
}

func DoGenerateIntegrations(args []string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	ints := kingpin.New("ints", "Generate integrations")
	root := ints.Flag("root", "sysl root directory for input model file (default: .)").Default(".").String()
	title := ints.Flag("title", "diagram title").Short('t').String()
	plantuml := ints.Flag("plantuml", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n")).Short('p').String()
	output := ints.Flag("output", "output file(default: %(epname).png)").Default("%(epname).png").Short('o').String()
	project := ints.Flag("project", "project pseudo-app to render").Short('j').String()
	filter := ints.Flag("filter", "Only generate diagrams whose output paths match a pattern").String()
	modules := ints.Arg("modules", strings.Join([]string{"input files without .sysl extension and with leading /",
		"eg: /project_dir/my_models",
		"combine with --root if needed"}, "\n")).String()
	exclude := ints.Flag("exclude", "apps to exclude").Short('e').Strings()
	clustered := ints.Flag("clustered", "group integration components into clusters").Short('c').Default("false").Bool()
	epa := ints.Flag("epa", "produce and EPA integration view").Default("false").Bool()
	loglevel := ints.Flag("log", "log level[debug,info,warn,off]").Default("warn").String()
	verbose := ints.Flag("verbose", "show output").Short('v').Default("false").Bool()

	_, err := ints.Parse(args[1:])
	if err != nil {
		log.Errorf("arguments parse error: %v", err)
	}
	if level, has := defaultLevel[*loglevel]; has {
		log.SetLevel(level)
	}
	if *plantuml == "" {
		*plantuml = os.Getenv("SYSL_PLANTUML")
		if *plantuml == "" {
			*plantuml = "http://localhost:8080/plantuml"
		}
	}
	log.Infof("root: %s\n", *root)
	log.Infof("project: %v\n", *project)
	log.Infof("clustered: %t\n", *clustered)
	log.Infof("exclude: %s\n", *exclude)
	log.Infof("epa: %t\n", *epa)
	log.Infof("title: %s\n", *title)
	log.Infof("plantuml: %s\n", *plantuml)
	log.Infof("filter: %s\n", *filter)
	log.Infof("modules: %s\n", *modules)
	log.Infof("output: %s\n", *output)
	log.Infof("loglevel: %s\n", *loglevel)

	r, err := GenerateIntegrations(*root, *title, *output, *project, *filter, *modules, *exclude, *clustered, *epa)
	if err != nil {
		return err
	}
	for k, v := range r {
		if *verbose {
			fmt.Println(k)
		}
		err := OutputPlantuml(k, *plantuml, v)
		if err != nil {
			return fmt.Errorf("plantuml drawing error for %#v: %v", k, err)
		}
	}
	return nil
}
