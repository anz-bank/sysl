package main

import (
	"flag"
	"io"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/seqs"
	log "github.com/sirupsen/logrus"
)

type SimpleParser struct {
	self string
}

type sequenceDiagParam struct {
	seqs.AppLabeler
	seqs.EndpointLabeler
	endpoints  []string
	title      string
	blackboxes [][]string
}

func generateSequenceDiag(m *sysl.Module, p *sequenceDiagParam) (string, error) {
	w := seqs.MakeSequenceDiagramWriter(true, "skinparam maxMessageSize 250")
	v := seqs.MakeSequenceDiagramVisitor(p.AppLabeler, p.EndpointLabeler, w, m)
	e := seqs.MakeEndpointCollectionElement(p.title, p.endpoints, p.blackboxes)

	e.Accept(v)

	return w.String(), nil
}

type arrayFlags []string

// String implements flag.Value.
func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

// Set implements flag.Value.
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func loadApp(root string, models []string) *sysl.Module {
	// Model we want to generate seqs for, the first non-empty model
	var model string
	for _, val := range models {
		if len(val) > 0 {
			model = val
			break
		}
	}
	mod, err := Parse(model, root)
	if err == nil {
		return mod
	}
	log.Errorf("unable to load module:\n\troot: " + root + "\n\tmodel:" + model)
	return nil
}

func mergeAttributesMap(val map[string]string, attrs map[string]*sysl.Attribute) map[string]string {
	for k, v := range attrs {
		val[k] = v.GetS()
	}

	return val
}

func (sp *SimpleParser) LabelEndpoint(p *seqs.EndpointLabelerParam) string {
	initialStr := sp.self
	attrs := map[string]string{
		"epname":       p.EndpointName,
		"human":        p.Human,
		"human_sender": p.HumanSender,
		"args":         p.Args,
		"patterns":     p.Patterns,
		"controls":     p.Controls,
	}
	attrs = mergeAttributesMap(attrs, p.Attrs)

	return seqs.ParseAttributesFormat(initialStr, attrs)
}

func (sp *SimpleParser) LabelApp(appname, controls string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	valMap := map[string]string{
		"appname":  appname,
		"controls": controls,
	}
	valMap = mergeAttributesMap(valMap, attrs)

	return seqs.ParseAttributesFormat(initialStr, valMap)
}

func (sp *SimpleParser) fmtSeq(epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	valMap := map[string]string{
		"epname":     epname,
		"eplongname": eplongname,
	}
	valMap = mergeAttributesMap(valMap, attrs)

	return seqs.ParseAttributesFormat(initialStr, valMap)
}

func (sp *SimpleParser) fmtOutput(appname, epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	valMap := map[string]string{
		"appname":    appname,
		"epname":     epname,
		"eplongname": eplongname,
	}
	valMap = mergeAttributesMap(valMap, attrs)

	return seqs.ParseAttributesFormat(initialStr, valMap)
}

func constructSimpleParser(former, latter string) *SimpleParser {
	fmtstr := former
	if former == "" {
		fmtstr = latter
	}

	return &SimpleParser{self: fmtstr}
}

func DoConstructSequenceDiagrams(
	root_model, endpoint_format, app_format, title, plantuml, output string,
	endpoints, apps, modules []string,
	blackboxes [][]string,
) {
	mod := loadApp(root_model, modules)

	if strings.Contains(output, "%(epname)") {
		spout := &SimpleParser{self: output}
		for _, appName := range apps {
			app := mod.Apps[appName]
			bbs := seqs.TransformBlackBoxes(app.GetAttrs()["blackboxes"].GetA().GetElt())
			spseqtitle := constructSimpleParser(app.GetAttrs()["seqtitle"].GetS(), title)
			spep := constructSimpleParser(app.GetAttrs()["epfmt"].GetS(), endpoint_format)
			spapp := constructSimpleParser(app.GetAttrs()["appfmt"].GetS(), app_format)
			keys := []string{}
			for k := range app.GetEndpoints() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				endpoint := app.GetEndpoints()[k]
				epAttrs := endpoint.GetAttrs()
				output_dir := spout.fmtOutput(appName, k, endpoint.GetLongName(), epAttrs)
				bbs2 := seqs.TransformBlackBoxes(app.GetEndpoints()[k].GetAttrs()["blackboxes"].GetA().GetElt())
				varrefs := seqs.MergeAttributes(app.GetAttrs(), endpoint.GetAttrs())
				sdEndpoints := []string{}
				statements := endpoint.GetStmt()
				for _, stmt := range statements {
					parts := stmt.GetCall().GetTarget().GetPart()
					ep := stmt.GetCall().GetEndpoint()
					sdEndpoints = append(sdEndpoints, strings.Join(parts, " :: ")+" <- "+ep)
				}

				sd := &sequenceDiagParam{
					endpoints:       sdEndpoints,
					AppLabeler:      spapp,
					EndpointLabeler: spep,
					title:           spseqtitle.fmtSeq(app.GetEndpoints()[k].GetName(), endpoint.GetLongName(), varrefs),
					blackboxes:      append(bbs, bbs2...),
				}
				out, _ := generateSequenceDiag(mod, sd)
				seqs.OutputPlantuml(output_dir, plantuml, out)
			}
		}
	} else {
		if endpoints == nil {
			return
		}
		spep := constructSimpleParser("", endpoint_format)
		spapp := constructSimpleParser("", app_format)
		sd := &sequenceDiagParam{
			endpoints:       endpoints,
			AppLabeler:      spapp,
			EndpointLabeler: spep,
			title:           title,
			blackboxes:      blackboxes,
		}
		out, _ := generateSequenceDiag(mod, sd)
		seqs.OutputPlantuml(output, plantuml, out)
	}
}

// DoGenerateSequenceDiagrams generate sequence diagrams for the given model
func DoGenerateSequenceDiagrams(stdout, stderr io.Writer, flags *flag.FlagSet, args []string) {
	var endpoints_flag, apps_flag, blackboxes_flag, modules_flag arrayFlags
	root_model := flags.String("root-model", ".", "sysl root directory for input model file (default: .)")
	endpoint_format := flags.String("endpoint_format", "%(epname)", "Specify the format string for sequence diagram endpoints. "+
		"May include %%(epname), %%(eplongname) and %%(@foo) for attribute foo(default: %(epname))")
	app_format := flags.String("app_format", "%(appname)", "Specify the format string for sequence diagram participants. "+
		"May include %%(appname) and %%(@foo) for attribute foo(default: %(appname))")
	title := flags.String("title", "", "diagram title")
	plantuml := flags.String("plantuml", "", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n"))
	output := flags.String("output", "%(epname).png", "output file(default: %(epname).png)")
	flags.Var(&endpoints_flag, "endpoint", "Include endpoint in sequence diagram")
	flags.Var(&apps_flag, "app", "Include all endpoints for app in sequence diagram (currently "+
		"only works with templated --output). Use SYSL_SD_FILTERS env (a "+
		"comma-list of shell globs) to limit the diagrams generated")
	flags.Var(&blackboxes_flag, "blackbox", "Apps to be treated as black boxes")
	flags.Var(&modules_flag, "modules", strings.Join([]string{"input files without .sysl extension and with leading /",
		"eg: /project_dir/my_models",
		"combine with --root if needed"}, "\n"))

	// Following variables currently are unused. Keep them to align with the python version.
	filter := flags.String("filter", "", "Only generate diagrams whose output paths match a pattern")
	no_activations := flags.Bool("no-activations", true, "Suppress sequence diagram activation bars(default: true)")
	verbose := flags.Bool("verbose", false, "Report each output(default: false)")
	expire_cache := flags.Bool("expire-cache", false, "Expire cache entries to force checking against real destination(default: false)")
	dry_run := flags.Bool("dry-run", false, "Don't perform confluence uploads, but show what would have happened(default: false)")

	err := flags.Parse(args[1:])
	if err != nil {
		log.Errorf("arguments parse error: %v", err)
	}
	log.Debugf("root_model: %s\n", *root_model)
	log.Debugf("endpoints: %v\n", endpoints_flag)
	log.Debugf("app: %v\n", apps_flag)
	log.Debugf("no_activations: %t\n", *no_activations)
	log.Debugf("endpoint_format: %s\n", *endpoint_format)
	log.Debugf("app_format: %s\n", *app_format)
	log.Debugf("blackbox: %s\n", blackboxes_flag)
	log.Debugf("title: %s\n", *title)
	log.Debugf("plantuml: %s\n", *plantuml)
	log.Debugf("verbose: %t\n", *verbose)
	log.Debugf("expire_cache: %t\n", *expire_cache)
	log.Debugf("dry_run: %t\n", *dry_run)
	log.Debugf("filter: %s\n", *filter)
	log.Debugf("modules: %s\n", modules_flag)
	log.Debugf("output: %s\n", *output)

	DoConstructSequenceDiagrams(*root_model, *endpoint_format, *app_format, *title, *plantuml, *output,
		endpoints_flag, apps_flag, modules_flag, seqs.ParseBlackBoxesFromArgument(blackboxes_flag))
}
