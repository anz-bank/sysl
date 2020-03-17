package integrationdiagram

import (
	"sort"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

const endpointWildcard = ".. * <- *"

type IntsBuilder struct {
	M            *sysl.Module
	SeedApps     []string
	SeedAppsMap  syslutil.StrSet
	Passthroughs syslutil.StrSet
	Excludes     syslutil.StrSet
	FinalApps    []string
	FinalAppsMap syslutil.StrSet
	Deps         syslutil.StrSet
	DepsOut      []AppDependency
}

func sortedSlice(endpts map[string]*sysl.Endpoint) []string {
	s := []string{}
	for epname := range endpts {
		s = append(s, epname)
	}
	sort.Strings(s)
	return s
}

func MakeBuilderfromStmt(m *sysl.Module, stmts []*sysl.Statement, excludes, passthroughs syslutil.StrSet) *IntsBuilder {
	b := &IntsBuilder{
		M:            m,
		SeedApps:     []string{},
		Excludes:     excludes,
		Passthroughs: passthroughs,
		FinalApps:    []string{},
		FinalAppsMap: syslutil.StrSet{},
		Deps:         syslutil.StrSet{},
		DepsOut:      []AppDependency{},
	}

	collector := endpointWildcard
	apps := b.M.GetApps()

	// build list to start from
	// aka SeedApps
	for _, stmt := range stmts {
		if a, ok := stmt.Stmt.(*sysl.Statement_Action); ok {
			app := apps[a.Action.Action]
			if app != nil && !syslutil.HasPattern(app.GetAttrs(), "human") {
				b.SeedApps = append(b.SeedApps, a.Action.Action)
				b.FinalApps = append(b.FinalApps, a.Action.Action)
			}
		}
	}

	b.SeedAppsMap = syslutil.MakeStrSet(b.SeedApps...)
	// add direct dependencies of SeedApps
	// dont include targetApp if its in exclude
	// if targetApp in passthrough, only add apps that are called by this endpoint
	// and do this recursively till the app is in passthrough
	for _, appname := range b.SeedApps {
		app := apps[appname]
		for _, epname := range sortedSlice(app.GetEndpoints()) {
			if epname != collector {
				endpt := app.GetEndpoints()[epname]
				ProcessCalls(appname, epname, endpt.GetStmt(), b.ProcessExcludeAndPassthrough)
			}
		}
	}

	// look for apps that are dependent on SeedApps
	appsNames := make([]string, 0, len(apps))

	for appname := range apps {
		appsNames = append(appsNames, appname)
	}
	sort.Strings(appsNames)
	for _, appname := range appsNames {
		app := apps[appname]
		for _, epname := range sortedSlice(app.GetEndpoints()) {
			if epname != collector {
				endpt := app.GetEndpoints()[epname]
				ProcessCalls(appname, epname, endpt.GetStmt(), b.MyCallers)
			}
		}
	}

	b.FinalAppsMap = syslutil.MakeStrSet(b.FinalApps...)
	// connect all the apps that we have so far
	for _, appname := range b.FinalApps {
		app := apps[appname]
		for _, epname := range sortedSlice(app.GetEndpoints()) {
			if epname != collector {
				endpt := app.GetEndpoints()[epname]
				ProcessCalls(appname, epname, endpt.GetStmt(), b.IndirectCalls)
			}
		}
	}

	return b
}

func (b *IntsBuilder) IndirectCalls(sourceApp, epname string, t *sysl.Statement) {
	call := t.GetCall()
	targetApp := syslutil.GetAppName(call.GetTarget())
	if !b.FinalAppsMap.Contains(targetApp) {
		return
	}
	if syslutil.HasPattern(b.M.GetApps()[targetApp].GetAttrs(), "human") {
		return
	}
	b.AddCall(sourceApp, epname, t)
}

func (b *IntsBuilder) MyCallers(sourceApp, epname string, t *sysl.Statement) {
	if b.Excludes.Contains(sourceApp) {
		return
	}
	call := t.GetCall()
	targetApp := syslutil.GetAppName(call.GetTarget())
	if !b.SeedAppsMap.Contains(targetApp) {
		return
	}
	if syslutil.HasPattern(b.M.GetApps()[targetApp].GetAttrs(), "human") {
		return
	}
	b.AddCall(sourceApp, epname, t)
	b.FinalApps = append(b.FinalApps, sourceApp)
}

func (b *IntsBuilder) WalkPassthrough(appname, epname string) {
	if b.Passthroughs.Contains(appname) {
		endpt := b.M.GetApps()[appname].GetEndpoints()[epname]
		ProcessCalls(appname, epname, endpt.GetStmt(), b.ProcessExcludeAndPassthrough)
	}
}

// meant to be called only for initial seed apps
// i.e. expect all calls to be added, except if the targetApp is in excluded list
func (b *IntsBuilder) ProcessExcludeAndPassthrough(sourceApp, epname string, t *sysl.Statement) {
	call := t.GetCall()
	targetApp := syslutil.GetAppName(call.GetTarget())
	if /* b.Excludes.Contains(sourceApp) || */ b.Excludes.Contains(targetApp) {
		return
	}
	if syslutil.HasPattern(b.M.GetApps()[targetApp].GetAttrs(), "human") {
		return
	}

	b.AddCall(sourceApp, epname, t)
	b.FinalApps = append(b.FinalApps, targetApp)
	b.WalkPassthrough(targetApp, call.Endpoint)
}

func (b *IntsBuilder) AddCall(appname, epname string, t *sysl.Statement) {
	call := t.GetCall()
	targetApp := syslutil.GetAppName(call.GetTarget())
	dep := AppDependency{
		Self:      AppElement{Name: appname, Endpoint: epname},
		Target:    AppElement{Name: targetApp, Endpoint: call.Endpoint},
		Statement: t,
	}

	k := dep.String()
	if _, has := b.Deps[k]; has {
		return
	}
	b.Deps.Insert(k)
	b.DepsOut = append(b.DepsOut, dep)
}

type callHandler func(appname, epname string, stmt *sysl.Statement)

func ProcessCalls(appname, epname string, stmts []*sysl.Statement, fn callHandler) {
	for _, stmt := range stmts {
		switch t := stmt.GetStmt().(type) {
		case *sysl.Statement_Call:
			fn(appname, epname, stmt)
		case *sysl.Statement_Action, *sysl.Statement_Ret:
			continue
		case *sysl.Statement_Cond:
			ProcessCalls(appname, epname, t.Cond.GetStmt(), fn)
		case *sysl.Statement_Loop:
			ProcessCalls(appname, epname, t.Loop.GetStmt(), fn)
		case *sysl.Statement_LoopN:
			ProcessCalls(appname, epname, t.LoopN.GetStmt(), fn)
		case *sysl.Statement_Foreach:
			ProcessCalls(appname, epname, t.Foreach.GetStmt(), fn)
		case *sysl.Statement_Group:
			ProcessCalls(appname, epname, t.Group.GetStmt(), fn)
		case *sysl.Statement_Alt:
			for _, choice := range t.Alt.GetChoice() {
				ProcessCalls(appname, epname, choice.GetStmt(), fn)
			}
		default:
			panic("No statement!")
		}
	}
}
