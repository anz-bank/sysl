package main

import (
	"fmt"
	"regexp"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
)

type AppElement struct {
	Name     string
	Endpoint string
}

type AppDependency struct {
	Self, Target AppElement
	Statement    *sysl.Statement
}

func (dep *AppDependency) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", dep.Self.Name, dep.Self.Endpoint, dep.Target.Name, dep.Target.Endpoint)
}

func MakeAppElement(name, endpoint string) AppElement {
	return AppElement{name, endpoint}
}

func dealWithPassthrough(excludes StrSet,
	dep AppDependency, integrations *collectApplicationDependenciesState, statements []*sysl.Statement) {
	for _, stmt := range statements {
		call := stmt.GetStmt().(*sysl.Statement_Call).Call
		nextAppName := strings.Join(call.GetTarget().GetPart(), " :: ")
		nextEpName := call.GetEndpoint()
		next := MakeAppElement(nextAppName, nextEpName)
		nextDep := AppDependency{dep.Target, next, stmt}
		excludedApps := MakeStrSet(nextAppName).Intersection(excludes)
		undeterminedEps := MakeStrSet(nextEpName).Intersection(MakeStrSet(".. * <- *", "*"))
		if len(excludedApps) == 0 && len(undeterminedEps) == 0 {
			integrations.deps[nextDep.String()] = nextDep
		}
	}
}

func (ds *collectApplicationDependenciesState) FindIntegrations(apps, excludes, passthroughs StrSet,
) map[string]AppDependency {
	integrations := &collectApplicationDependenciesState{
		deps: map[string]AppDependency{},
	}
	lenPassthroughs := len(passthroughs)
	commonEndpoints := MakeStrSet(".. * <- *", "*")
	for _, dep := range ds.deps {
		appNames := dep.extractAppNames()
		endpoints := dep.extractEndpoints()
		isSubsection := appNames.IsSubset(apps)
		isSelfSubsection := MakeStrSet(dep.Self.Name).IsSubset(apps)
		isSelfPassthrough := MakeStrSet(dep.Self.Name).IsSubset(passthroughs)
		isTargetPassthrough := MakeStrSet(dep.Target.Name).IsSubset(passthroughs)
		interExcludes := appNames.Intersection(excludes)
		interEndpoints := endpoints.Intersection(commonEndpoints)
		lenInterExcludes := len(interExcludes)
		lenInterEndpoints := len(interEndpoints)
		if isSubsection && lenInterExcludes == 0 && lenInterEndpoints == 0 {
			integrations.deps[dep.String()] = dep
		}
		// deal with passthrough
		if lenPassthroughs > 0 &&
			(((isSelfSubsection && isTargetPassthrough) || (isSelfPassthrough && isTargetPassthrough)) &&
				lenInterExcludes == 0 && lenInterEndpoints == 0) {
			integrations.deps[dep.String()] = dep
			dealWithPassthrough(excludes, dep, integrations, ds.AppCalls[makeAppCallKey(dep.Target.Name, dep.Target.Endpoint)])
		}
	}

	return integrations.deps
}

func FindApps(module *sysl.Module, excludes, integrations StrSet, ds map[string]AppDependency, drawable bool) StrSet {
	output := []string{}
	appReStr := toPattern(integrations.ToSlice())
	re := regexp.MustCompile(appReStr)
	for _, dep := range ds {
		appNames := dep.extractAppNames()
		excludeApps := appNames.Intersection(excludes)
		if len(excludeApps) > 0 {
			continue
		}
		highlightApps := appNames.Intersection(integrations)
		if !drawable && len(highlightApps) == 0 {
			continue
		}
		for _, item := range appNames.ToSlice() {
			app := module.GetApps()[item]
			if drawable {
				if re.MatchString(item) &&
					!HasPattern(app.GetAttrs(), "human") {
					output = append(output, item)
				}
				continue
			}
			if !HasPattern(app.GetAttrs(), "human") {
				output = append(output, item)
			}
		}
	}

	return MakeStrSet(output...)
}

func toPattern(comp []string) string {
	return fmt.Sprintf(`^(?:%s)(?: *::|$)`, strings.Join(comp, "|"))
}

func (dep *AppDependency) extractAppNames() StrSet {
	s := StrSet{}
	s.Insert(dep.Self.Name)
	s.Insert(dep.Target.Name)
	return s
}

func (dep *AppDependency) extractEndpoints() StrSet {
	s := StrSet{}
	s.Insert(dep.Self.Endpoint)
	s.Insert(dep.Target.Endpoint)
	return s
}

func makeAppCallKey(appname, epname string) string {
	return strings.Join([]string{appname, epname}, ":")
}
