package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
)

const ArrowColorNone = "none"
const PumlHeader = `''''''''''''''''''''''''''''''''''''''''''
''                                      ''
''  AUTOGENERATED CODE -- DO NOT EDIT!  ''
''                                      ''
''''''''''''''''''''''''''''''''''''''''''

`
const ComponentStart = `hide stereotype
scale max 16384 height
skinparam component {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
`
const StateStart = `left to right direction
scale max 16384 height
hide empty description
skinparam state {
  BackgroundColor FloralWhite
  BorderColor Black
  ArrowColor Crimson
`

type IntsParam struct {
	apps         []string
	drawableApps StrSet
	integrations []AppDependency
	app          *sysl.Application
	endpt        *sysl.Endpoint
}

type Args struct {
	title     string
	project   string
	clustered bool
	epa       bool
}

type IntsDiagramVisitor struct {
	mod           *sysl.Module
	stringBuilder *strings.Builder
	drawableApps  StrSet
	symbols       map[string]*_var
	topSymbols    map[string]*_topVar
	project       string
}

type _topVar struct {
	topLabel string
	topAlias string
}

type AppPair struct {
	Self   string
	Target string
}

type viewParams struct {
	restrictBy         string
	endptAttrs         map[string]*sysl.Attribute
	highLightColor     string
	arrowColor         string
	indirectArrowColor string
	diagramTitle       string
}

func MakeIntsDiagramVisitor(
	mod *sysl.Module, stringBuilder *strings.Builder,
	drawableApps StrSet, project string,
) *IntsDiagramVisitor {
	return &IntsDiagramVisitor{
		mod:           mod,
		stringBuilder: stringBuilder,
		drawableApps:  drawableApps,
		symbols:       map[string]*_var{},
		topSymbols:    map[string]*_topVar{},
		project:       project,
	}
}

func (v *IntsDiagramVisitor) VarManagerForComponent(appName string, nameMap map[string]string) string {
	if key, ok := nameMap[appName]; ok {
		appName = key
	}
	if s, ok := v.symbols[appName]; ok {
		return s.alias
	}

	i := len(v.symbols)
	alias := fmt.Sprintf("_%d", i)
	fp := MakeFormatParser(v.mod.Apps[v.project].GetAttrs()["appfmt"].GetS())
	attrs := getApplicationAttrs(v.mod, appName)
	controls := getSortedISOCtrlStr(attrs)
	label := fp.LabelApp(appName, controls, attrs)
	s := &_var{
		label: label,
		alias: alias,
	}
	v.symbols[appName] = s
	if _, ok := v.drawableApps[appName]; ok {
		fmt.Fprintf(v.stringBuilder, "[%s] as %s <<highlight>>\n", label, alias)
	} else {
		fmt.Fprintf(v.stringBuilder, "[%s] as %s\n", label, alias)
	}
	return s.alias
}

func (v *IntsDiagramVisitor) VarManagerForTopState(appName string) string {
	var alias, label string
	if ts, ok := v.topSymbols[appName]; ok {
		return ts.topAlias
	}
	i := len(v.topSymbols)
	alias = fmt.Sprintf("_%d", i)

	fp := MakeFormatParser(v.mod.Apps[v.project].GetAttrs()["appfmt"].GetS())
	attrs := getApplicationAttrs(v.mod, appName)
	controls := getSortedISOCtrlStr(attrs)
	label = fp.LabelApp(appName, controls, attrs)
	ts := &_topVar{
		topLabel: label,
		topAlias: alias,
	}
	v.topSymbols[appName] = ts
	if _, ok := v.drawableApps[appName]; ok {
		fmt.Fprintf(v.stringBuilder, "state \"%s\" as X%s <<highlight>> {\n", label, alias)
	} else {
		fmt.Fprintf(v.stringBuilder, "state \"%s\" as X%s {\n", label, alias)
	}

	return ts.topAlias
}

func (v *IntsDiagramVisitor) VarManagerForEPA(name string) string {
	var appName, alias, label string
	attrs := map[string]string{}

	appName = strings.Split(name, " : ")[0]
	epName := strings.Split(name, " : ")[1]

	if s, ok := v.symbols[name]; ok {
		return s.alias
	}
	i := len(v.symbols)
	alias = fmt.Sprintf("_%d", i)

	if v.mod.Apps[appName].Endpoints[epName] != nil {
		for k, v := range v.mod.Apps[appName].Endpoints[epName].Attrs {
			attrs["@"+k] = v.GetS()
		}
	}
	attrs["appname"] = epName
	fp := MakeFormatParser(v.mod.Apps[v.project].GetAttrs()["appfmt"].GetS())
	label = fp.Parse(attrs)

	s := &_var{
		label: label,
		alias: alias,
	}
	v.symbols[name] = s

	if _, ok := v.drawableApps[appName]; ok {
		fmt.Fprintf(v.stringBuilder, "  state \"%s\" as %s <<highlight>>\n", label, alias)
	} else {
		fmt.Fprintf(v.stringBuilder, "  state \"%s\" as %s\n", label, alias)
	}
	return s.alias
}

func (v *IntsDiagramVisitor) buildClusterForEPAView(deps []AppDependency, restrictBy string) {
	clusters := map[string][]string{}
	for _, dep := range deps {
		appA := dep.Self.Name
		appB := dep.Target.Name
		epA := dep.Self.Endpoint
		epB := dep.Target.Endpoint
		_, okA := v.mod.Apps[appA].Endpoints[epA].Attrs[restrictBy]
		_, okB := v.mod.Apps[appB].Endpoints[epB].Attrs[restrictBy]
		if _, ok := v.mod.Apps[appA].Attrs[restrictBy]; !ok && restrictBy != "" {
			if _, ok := v.mod.Apps[appB].Attrs[restrictBy]; !ok {
				continue
			}
		}
		if !okA && restrictBy != "" && !okB {
			continue
		}
		clusters[appA] = append(clusters[appA], epA)
		if appA != appB && !v.mod.Apps[appA].Endpoints[epA].IsPubsub {
			clusters[appA] = append(clusters[appA], epB+" client")
		}
		clusters[appB] = append(clusters[appB], epB)
	}

	keys := []string{}
	for k := range clusters {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v.VarManagerForTopState(k)
		strSet := MakeStrSet(clusters[k]...)
		for _, m := range strSet.ToSortedSlice() {
			v.VarManagerForEPA(k + " : " + m)
		}
		v.stringBuilder.WriteString("}\n")
	}
}

func (v *IntsDiagramVisitor) buildClusterForIntsView(apps []string) map[string]string {
	nameMap := map[string]string{}
	clusters := map[string][]string{}
	for _, v := range apps {
		cluster := strings.Split(v, " :: ")
		if len(cluster) > 1 {
			clusters[cluster[0]] = append(clusters[cluster[0]], v)
		}
	}

	for k, v := range clusters {
		if len(v) <= 1 {
			delete(clusters, k)
		}
		for _, s := range v {
			nameMap[s] = strings.Split(s, " :: ")[1]
		}
	}

	for k, apps := range clusters {
		fmt.Fprintf(v.stringBuilder, "package \"%s\" {\n", k)
		for _, n := range apps {
			v.VarManagerForComponent(n, nameMap)
		}
		v.stringBuilder.WriteString("}\n")
	}

	return nameMap
}

func (v *IntsDiagramVisitor) generateEPAView(viewParams viewParams, params *IntsParam) string {
	v.stringBuilder.WriteString("@startuml\n")
	if viewParams.diagramTitle != "" {
		fmt.Fprintf(v.stringBuilder, "title %s\n", viewParams.diagramTitle)
	}
	v.stringBuilder.WriteString(StateStart)
	if viewParams.highLightColor != "" {
		fmt.Fprintf(v.stringBuilder, "  BackgroundColor<<highlight>> %s\n", viewParams.highLightColor)
	}
	if viewParams.arrowColor != "" {
		fmt.Fprintf(v.stringBuilder, "  ArrowColor %s\n", viewParams.arrowColor)
	}

	if viewParams.indirectArrowColor != "" && viewParams.indirectArrowColor != ArrowColorNone {
		fmt.Fprintf(v.stringBuilder, "  ArrowColor<<indirect>> %s\n", viewParams.indirectArrowColor)
	}
	v.stringBuilder.WriteString("}\n")
	v.buildClusterForEPAView(params.integrations, viewParams.restrictBy)
	var processed []string
	for _, dep := range params.integrations {
		appA := dep.Self.Name
		appB := dep.Target.Name
		epA := dep.Self.Endpoint
		epB := dep.Target.Endpoint
		_, restrictByAppA := v.mod.Apps[appA].Attrs[viewParams.restrictBy]
		_, restrictByAppB := v.mod.Apps[appB].Attrs[viewParams.restrictBy]
		_, restrictByEpA := v.mod.Apps[appA].Endpoints[epA].Attrs[viewParams.restrictBy]
		_, restrictByEpB := v.mod.Apps[appB].Endpoints[epB].Attrs[viewParams.restrictBy]
		if viewParams.restrictBy != "" && !(restrictByAppA || restrictByAppB) {
			continue
		}
		if viewParams.restrictBy != "" && !(restrictByEpA || restrictByEpB) {
			continue
		}
		matchApp := appB
		matchEp := epB
		label := ""
		needsInt := appA != matchApp

		pubSubSrcPtrns := MakeStrSetFromAttr("patterns", v.mod.Apps[appA].Endpoints[epA].Attrs)
		attrs := map[string]string{}
		for k, v := range dep.Statement.GetAttrs() {
			attrs["@"+k] = v.GetS()
		}
		var ptrns string
		var targetPatterns StrSet
		var srcPtrns StrSet
		if v.mod.Apps[matchApp].Endpoints[matchEp].Attrs["patterns"] != nil {
			targetPatterns = MakeStrSetFromAttr("patterns", v.mod.Apps[matchApp].Endpoints[matchEp].Attrs)
		} else {
			targetPatterns = MakeStrSetFromAttr("patterns", v.mod.Apps[matchApp].Attrs)
		}
		if dep.Statement.GetAttrs()["patterns"] != nil {
			srcPtrns = MakeStrSetFromAttr("patterns", dep.Statement.GetAttrs())
		} else {
			srcPtrns = pubSubSrcPtrns
		}
		if srcPtrns != nil || targetPatterns != nil {
			ptrns = strings.Join(srcPtrns.ToSlice(), ", ") + " → " + strings.Join(targetPatterns.ToSlice(), ", ")
		} else {
			ptrns = ""
		}
		attrs["patterns"] = ptrns
		if needsInt {
			attrs["needs_int"] = strconv.FormatBool(needsInt)
		}
		attrs["epname"] = params.endpt.GetName()
		attrs["eplongname"] = params.endpt.GetLongName()
		fp := MakeFormatParser(params.app.Attrs["epfmt"].GetS())
		label = fp.Parse(attrs)

		flow := strings.Join([]string{appA, epB, appB, epB}, ".")
		isPubSub := v.mod.Apps[appA].Endpoints[epA].GetIsPubsub()
		epBClient := epB + " client"

		if appA != appB {
			if label != "" {
				label = " : " + label
			}
			if isPubSub {
				fmt.Fprintf(
					v.stringBuilder,
					"%s -%s> %s%s\n", v.VarManagerForEPA(appA+" : "+epA),
					"[#blue]",
					v.VarManagerForEPA(appB+" : "+epB),
					label,
				)
			} else {
				color := ""
				if viewParams.indirectArrowColor == "" {
					color = "[#silver]-"
				} else {
					color = "[#" + viewParams.indirectArrowColor + "]-"
				}
				fmt.Fprintf(
					v.stringBuilder,
					"%s -%s> %s\n",
					v.VarManagerForEPA(appA+" : "+epA),
					color,
					v.VarManagerForEPA(appA+" : "+epBClient),
				)
				if !stringInSlice(flow, processed) {
					fmt.Fprintf(
						v.stringBuilder,
						"%s -%s> %s%s\n",
						v.VarManagerForEPA(appA+" : "+epBClient),
						"[#black]",
						v.VarManagerForEPA(appB+" : "+epB),
						label,
					)
					processed = append(processed, flow)
				}
			}
		} else {
			color := ""
			if viewParams.indirectArrowColor == "" {
				color = "[#silver]-"
			} else {
				color = "[#" + viewParams.indirectArrowColor + "]-"
			}
			fmt.Fprintf(
				v.stringBuilder,
				"%s -%s> %s%s\n",
				v.VarManagerForEPA(appA+" : "+epA),
				color,
				v.VarManagerForEPA(appB+" : "+epB),
				label,
			)
		}
	}
	v.stringBuilder.WriteString("@enduml")
	return v.stringBuilder.String()

}

func (v *IntsDiagramVisitor) drawIntsView(viewParams viewParams, params *IntsParam, nameMap map[string]string) {
	callsDrawn := map[AppPair]struct{}{}
	if viewParams.endptAttrs["view"].GetS() == "system" {
		v.drawSystemView(viewParams, params, nameMap)
	} else {
		for _, dep := range params.integrations {
			appA := dep.Self.Name
			appB := dep.Target.Name
			if appA == appB {
				continue
			}
			appPair := AppPair{
				Self:   appA,
				Target: appB,
			}
			var direct StrSet
			apps := MakeStrSet(appA, appB)
			direct = apps.Intersection(params.drawableApps)
			if _, ok := callsDrawn[appPair]; !ok {
				if len(direct) > 0 || direct != nil || viewParams.indirectArrowColor != ArrowColorNone {
					indirect := ""
					if len(direct) == 0 {
						indirect = " <<indirect>>"
					}
					fmt.Fprintf(
						v.stringBuilder,
						"%s --> %s%s\n",
						v.VarManagerForComponent(appA, nameMap),
						v.VarManagerForComponent(appB, nameMap),
						indirect,
					)
					callsDrawn[appPair] = struct{}{}
				}
			}
		}
		for _, app := range params.apps {
			for _, mixin := range v.mod.Apps[app].GetMixin2() {
				mixinName := strings.Join(mixin.Name.Part, " :: ")
				fmt.Fprintf(
					v.stringBuilder,
					"%s <|.. %s\n",
					v.VarManagerForComponent(mixinName, nameMap),
					v.VarManagerForComponent(app, nameMap),
				)
			}
		}
	}
}

func (v *IntsDiagramVisitor) drawSystemView(viewParams viewParams, params *IntsParam, nameMap map[string]string) {
	callsDrawn := map[AppPair]struct{}{}
	for _, dep := range params.integrations {
		appA := dep.Self.Name
		appB := dep.Target.Name
		if appA == appB {
			continue
		}
		appPair := AppPair{
			Self:   appA,
			Target: appB,
		}
		var direct []string
		if _, ok := params.drawableApps[appA]; ok {
			direct = append(direct, appA)
		}
		if _, ok := params.drawableApps[appB]; ok {
			direct = append(direct, appB)
		}
		appA = strings.Split(appA, " :: ")[0]
		appB = strings.Split(appB, " :: ")[0]
		if _, ok := callsDrawn[appPair]; !ok {
			if len(direct) > 0 || direct != nil || viewParams.indirectArrowColor != ArrowColorNone {
				indirect := ""
				if len(direct) == 0 {
					indirect = " <<indirect>>"
				}
				fmt.Fprintf(
					v.stringBuilder,
					"%s --> %s%s\n",
					v.VarManagerForComponent(appA, nameMap),
					v.VarManagerForComponent(appB, nameMap),
					indirect,
				)
				callsDrawn[appPair] = struct{}{}
			}
		}
	}
}

func (v *IntsDiagramVisitor) generateIntsView(args *Args, viewParams viewParams, params *IntsParam) string {

	v.stringBuilder.WriteString("@startuml\n")
	if viewParams.diagramTitle != "" {
		fmt.Fprintf(v.stringBuilder, "title %s\n", viewParams.diagramTitle)
	}
	v.stringBuilder.WriteString(ComponentStart)
	if viewParams.highLightColor != "" {
		fmt.Fprintf(v.stringBuilder, "  BackgroundColor<<highlight>> %s\n", viewParams.highLightColor)
	}
	if viewParams.arrowColor != "" {
		fmt.Fprintf(v.stringBuilder, "  ArrowColor %s\n", viewParams.arrowColor)
	}

	if viewParams.indirectArrowColor != "" && viewParams.indirectArrowColor != ArrowColorNone {
		fmt.Fprintf(v.stringBuilder, "  ArrowColor<<indirect>> %s\n", viewParams.indirectArrowColor)
	}
	v.stringBuilder.WriteString("}\n")
	nameMap := map[string]string{}
	if args.clustered || viewParams.endptAttrs["view"].GetS() == "clustered" {
		nameMap = v.buildClusterForIntsView(params.apps)
	}
	v.drawIntsView(viewParams, params, nameMap)
	v.stringBuilder.WriteString("@enduml")
	return v.stringBuilder.String()
}

func GenerateView(args *Args, params *IntsParam, mod *sysl.Module) string {
	var stringBuilder strings.Builder
	var titleParser *FormatParser
	v := MakeIntsDiagramVisitor(mod, &stringBuilder, params.drawableApps, args.project)
	restrictBy := ""
	if params.endpt.Attrs["restrict_by"] != nil {
		restrictBy = params.endpt.Attrs["restrict_by"].GetS()
	}

	appAttrs := params.app.Attrs
	endptAttrs := params.endpt.Attrs
	highLightColor := appAttrs["highlight_color"].GetS()
	arrowColor := appAttrs["arrow_color"].GetS()
	indirectArrowColor := appAttrs["indirect_arrow_color"].GetS()

	attrs := map[string]string{
		"epname":     params.endpt.Name,
		"eplongname": params.endpt.LongName,
	}
	title := args.title
	if appAttrs["title"].GetS() != "" {
		title = appAttrs["title"].GetS()
	}
	titleParser = MakeFormatParser(title)
	diagramTitle := titleParser.Parse(attrs)

	viewParams := &viewParams{
		restrictBy:         restrictBy,
		endptAttrs:         endptAttrs,
		highLightColor:     highLightColor,
		arrowColor:         arrowColor,
		indirectArrowColor: indirectArrowColor,
		diagramTitle:       diagramTitle,
	}
	v.stringBuilder.WriteString(PumlHeader)

	if args.epa || endptAttrs["view"].GetS() == "epa" {
		return v.generateEPAView(*viewParams, params)
	}
	return v.generateIntsView(args, *viewParams, params)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
