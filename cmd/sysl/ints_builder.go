package main

import (
	"sort"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

type intsBuilder struct {
	m            *sysl.Module
	seedApps     []string
	seedAppsMap  syslutil.StrSet
	passthroughs syslutil.StrSet
	excludes     syslutil.StrSet
	finalApps    []string
	finalAppsMap syslutil.StrSet
	deps         syslutil.StrSet
	depsOut      []AppDependency
}

func sortedSlice(endpts map[string]*sysl.Endpoint) []string {
	s := []string{}
	for epname := range endpts {
		s = append(s, epname)
	}
	sort.Strings(s)
	return s
}

func makeBuilderfromStmt(m *sysl.Module, stmts []*sysl.Statement, excludes, passthroughs syslutil.StrSet) *intsBuilder {
	b := &intsBuilder{
		m:            m,
		seedApps:     []string{},
		excludes:     excludes,
		passthroughs: passthroughs,
		finalApps:    []string{},
		finalAppsMap: syslutil.StrSet{},
		deps:         syslutil.StrSet{},
		depsOut:      []AppDependency{},
	}

	collector := endpointWildcard
	apps := b.m.GetApps()

	// build list to start from
	// aka seedApps
	for _, stmt := range stmts {
		if a, ok := stmt.Stmt.(*sysl.Statement_Action); ok {
			app := apps[a.Action.Action]
			if app != nil && !syslutil.HasPattern(app.GetAttrs(), "human") {
				b.seedApps = append(b.seedApps, a.Action.Action)
				b.finalApps = append(b.finalApps, a.Action.Action)
			}
		}
	}

	b.seedAppsMap = syslutil.MakeStrSet(b.seedApps...)
	// add direct dependencies of seedApps
	// dont include targetApp if its in exclude
	// if targetApp in passthrough, only add apps that are called by this endpoint
	// and do this recursively till the app is in passthrough
	for _, appname := range b.seedApps {
		app := apps[appname]
		for _, epname := range sortedSlice(app.GetEndpoints()) {
			if epname != collector {
				endpt := app.GetEndpoints()[epname]
				processCalls(appname, epname, endpt.GetStmt(), b.processExcludeAndPassthrough)
			}
		}
	}

	// look for apps that are dependent on seedApps
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
				processCalls(appname, epname, endpt.GetStmt(), b.myCallers)
			}
		}
	}

	b.finalAppsMap = syslutil.MakeStrSet(b.finalApps...)
	// connect all the apps that we have so far
	for _, appname := range b.finalApps {
		app := apps[appname]
		for _, epname := range sortedSlice(app.GetEndpoints()) {
			if epname != collector {
				endpt := app.GetEndpoints()[epname]
				processCalls(appname, epname, endpt.GetStmt(), b.indirectCalls)
			}
		}
	}

	return b
}

func (b *intsBuilder) indirectCalls(sourceApp, epname string, t *sysl.Statement) {
	call := t.GetCall()
	targetApp := syslutil.GetAppName(call.GetTarget())
	if !b.finalAppsMap.Contains(targetApp) {
		return
	}
	if syslutil.HasPattern(b.m.GetApps()[targetApp].GetAttrs(), "human") {
		return
	}
	b.addCall(sourceApp, epname, t)
}

func (b *intsBuilder) myCallers(sourceApp, epname string, t *sysl.Statement) {
	if b.excludes.Contains(sourceApp) {
		return
	}
	call := t.GetCall()
	targetApp := syslutil.GetAppName(call.GetTarget())
	if !b.seedAppsMap.Contains(targetApp) {
		return
	}
	if syslutil.HasPattern(b.m.GetApps()[targetApp].GetAttrs(), "human") {
		return
	}
	b.addCall(sourceApp, epname, t)
	b.finalApps = append(b.finalApps, sourceApp)
}

func (b *intsBuilder) walkPassthrough(appname, epname string) {
	if b.passthroughs.Contains(appname) {
		endpt := b.m.GetApps()[appname].GetEndpoints()[epname]
		processCalls(appname, epname, endpt.GetStmt(), b.processExcludeAndPassthrough)
	}
}

// meant to be called only for initial seed apps
// i.e. expect all calls to be added, except if the targetApp is in excluded list
func (b *intsBuilder) processExcludeAndPassthrough(sourceApp, epname string, t *sysl.Statement) {
	call := t.GetCall()
	targetApp := syslutil.GetAppName(call.GetTarget())
	if /* b.excludes.Contains(sourceApp) || */ b.excludes.Contains(targetApp) {
		return
	}
	if syslutil.HasPattern(b.m.GetApps()[targetApp].GetAttrs(), "human") {
		return
	}

	b.addCall(sourceApp, epname, t)
	b.finalApps = append(b.finalApps, targetApp)
	b.walkPassthrough(targetApp, call.Endpoint)
}

func (b *intsBuilder) addCall(appname, epname string, t *sysl.Statement) {
	call := t.GetCall()
	targetApp := syslutil.GetAppName(call.GetTarget())
	dep := AppDependency{
		Self:      AppElement{Name: appname, Endpoint: epname},
		Target:    AppElement{Name: targetApp, Endpoint: call.Endpoint},
		Statement: t,
	}

	k := dep.String()
	if _, has := b.deps[k]; has {
		return
	}
	b.deps.Insert(k)
	b.depsOut = append(b.depsOut, dep)
}

type callHandler func(appname, epname string, stmt *sysl.Statement)

func processCalls(appname, epname string, stmts []*sysl.Statement, fn callHandler) {
	for _, stmt := range stmts {
		switch t := stmt.GetStmt().(type) {
		case *sysl.Statement_Call:
			fn(appname, epname, stmt)
		case *sysl.Statement_Action, *sysl.Statement_Ret:
			continue
		case *sysl.Statement_Cond:
			processCalls(appname, epname, t.Cond.GetStmt(), fn)
		case *sysl.Statement_Loop:
			processCalls(appname, epname, t.Loop.GetStmt(), fn)
		case *sysl.Statement_LoopN:
			processCalls(appname, epname, t.LoopN.GetStmt(), fn)
		case *sysl.Statement_Foreach:
			processCalls(appname, epname, t.Foreach.GetStmt(), fn)
		case *sysl.Statement_Group:
			processCalls(appname, epname, t.Group.GetStmt(), fn)
		case *sysl.Statement_Alt:
			for _, choice := range t.Alt.GetChoice() {
				processCalls(appname, epname, choice.GetStmt(), fn)
			}
		default:
			panic("No statement!")
		}
	}
}
