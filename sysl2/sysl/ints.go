package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type collectApplicationDependenciesState struct {
	m            *sysl.Module
	excludes     StrSet
	passthroughs StrSet
	drawableApps StrSet
	apps         StrSet
	deps         map[string]AppDependency
}

func collectApplicationDependencies(
	mod *sysl.Module,
	excludes, passthroughs, highlights, integrations StrSet,
) (drawableApps, apps StrSet, dependencies map[string]AppDependency) {
	v := &collectApplicationDependenciesState{
		m:            mod,
		excludes:     excludes,
		passthroughs: passthroughs,
		drawableApps: highlights,
		apps:         highlights.Clone(),
		deps:         map[string]AppDependency{},
	}
	applications := v.m.GetApps()
	for app := range integrations {
		for epname, endpt := range applications[app].GetEndpoints() {
			v.handleStatement(app, epname, endpt.GetStmt())
		}
	}
	return v.drawableApps, v.apps, v.deps
}

func (v *collectApplicationDependenciesState) handleStatement(appname, epname string, stmts []*sysl.Statement) {
	for _, stmt := range stmts {
		switch t := stmt.GetStmt().(type) {
		case *sysl.Statement_Call:
			targetName := getAppName(t.Call.GetTarget())
			if added := v.addAppDependency(appname, epname, targetName, t.Call.Endpoint, stmt); added {
				if v.passthroughs.Contains(targetName) && t.Call.Endpoint != ".. * <- *" {
					v.handleStatement(targetName, t.Call.Endpoint, v.m.GetApps()[targetName].GetEndpoints()[t.Call.Endpoint].Stmt)
				}
			}
		case *sysl.Statement_Action, *sysl.Statement_Ret:
			continue
		case *sysl.Statement_Cond:
			v.handleStatement(appname, epname, t.Cond.GetStmt())
		case *sysl.Statement_Loop:
			v.handleStatement(appname, epname, t.Loop.GetStmt())
		case *sysl.Statement_LoopN:
			v.handleStatement(appname, epname, t.LoopN.GetStmt())
		case *sysl.Statement_Foreach:
			v.handleStatement(appname, epname, t.Foreach.GetStmt())
		case *sysl.Statement_Group:
			v.handleStatement(appname, epname, t.Group.GetStmt())
		case *sysl.Statement_Alt:
			for _, choice := range t.Alt.GetChoice() {
				v.handleStatement(appname, epname, choice.GetStmt())
			}
		default:
			panic("No statement!")
		}
	}
}

func (v *collectApplicationDependenciesState) addAppDependency(
	sourceApp, sourceEndpt, targetApp, targetEndpt string,
	stmt *sysl.Statement,
) bool {
	if sourceEndpt == ".. * <- *" || sourceEndpt == "*" {
		return false
	}
	k := fmt.Sprintf("%s:%s:%s:%s", sourceApp, sourceEndpt, targetApp, targetEndpt)
	if _, has := v.deps[k]; has {
		return false
	}
	if v.excludes.Contains(sourceApp) || v.excludes.Contains(targetApp) {
		return false
	}
	if v.drawableApps.Contains(sourceApp) || v.drawableApps.Contains(targetApp) {
		if !HasPattern(v.m.GetApps()[sourceApp].GetAttrs(), "human") {
			v.apps.Insert(sourceApp)
		}
		if !HasPattern(v.m.GetApps()[targetApp].GetAttrs(), "human") {
			v.apps.Insert(targetApp)
		}
	}
	if !((v.apps.Contains(sourceApp) && v.apps.Contains(targetApp)) ||
		(v.apps.Contains(sourceApp) && v.passthroughs.Contains(targetApp)) ||
		v.passthroughs.Contains(sourceApp)) {
		return false
	}
	v.deps[k] = AppDependency{
		Self:      AppElement{Name: sourceApp, Endpoint: sourceEndpt},
		Target:    AppElement{Name: targetApp, Endpoint: targetEndpt},
		Statement: stmt,
	}
	return true
}

func GenerateIntegrations(
	rootModel, title, output, project, filter, modules string,
	exclude []string,
	clustered, epa bool,
) map[string]string {
	r := make(map[string]string)

	mod := loadApp(rootModel, modules)

	if len(exclude) == 0 && project != "" {
		exclude = []string{project}
	}
	excludeStrSet := MakeStrSet(exclude...)

	// The "project" app that specifies the required view of the integration
	app := mod.GetApps()[project]
	of := MakeFormatParser(output)
	// Interate over each endpoint within the selected project
	for epname, endpt := range app.GetEndpoints() {
		excludes := MakeStrSetFromAttr("exclude", endpt.GetAttrs())
		passthroughs := MakeStrSetFromAttr("passthrough", endpt.GetAttrs())
		integrations := MakeStrSetFromActionStatement(endpt.GetStmt())
		outputDir := of.FmtOutput(project, epname, endpt.GetLongName(), endpt.GetAttrs())
		if filter != "" {
			re := regexp.MustCompile(filter)
			if !re.MatchString(outputDir) {
				continue
			}
		}
		drawableApps := findDrawableApps(mod, integrations)
		drawableApps, apps, deps := collectApplicationDependencies(
			mod, excludeStrSet.Union(excludes), passthroughs, drawableApps, integrations)
		intsParam := &IntsParam{apps.ToSlice(), drawableApps, deps, app, endpt}
		args := &Args{title, project, clustered, epa}
		r[outputDir] = GenerateView(args, intsParam, mod)
	}

	return r
}

func findDrawableApps(mod *sysl.Module, apps StrSet) StrSet {
	drawable := StrSet{}
	for appName := range apps {
		app := mod.GetApps()[appName]
		if !HasPattern(app.GetAttrs(), "human") {
			drawable.Insert(appName)
		}
	}

	return drawable
}

func DoGenerateIntegrations(args []string) {
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
	log.Debugf("root: %s\n", *root)
	log.Debugf("project: %v\n", project)
	log.Debugf("clustered: %t\n", *clustered)
	log.Debugf("exclude: %s\n", *exclude)
	log.Debugf("epa: %t\n", *epa)
	log.Debugf("title: %s\n", *title)
	log.Debugf("plantuml: %s\n", *plantuml)
	log.Debugf("filter: %s\n", *filter)
	log.Debugf("modules: %s\n", *modules)
	log.Debugf("output: %s\n", *output)
	log.Debugf("loglevel: %s\n", *loglevel)

	r := GenerateIntegrations(*root, *title, *output, *project, *filter, *modules, *exclude, *clustered, *epa)
	for k, v := range r {
		err := OutputPlantuml(k, *plantuml, v)
		if err != nil {
			log.Errorf("plantuml drawing error: %v", err)
		}
	}
}
