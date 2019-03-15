package main

import (
	"flag"
	"io"
	"os"
	"sort"
	"strings"
	"sysl/sysl2/sysl/seqs"

	"github.com/anz-bank/sysl/src/proto"
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

	v.Visit(e)

	return w.String(), nil
}

type arrayFlags []string

// implement the String method of Value interface in flag.go
func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

// implement the Set method of Value interface in flag.go
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func loadApp(root string, models []string) *sysl.Module {
	// Model we want to generate seqs for
	var model string
	for _, val := range models {
		model = val
		break
	}
	mod, err := Parse(model, root)
	if err == nil {
		return mod
	}
	log.Errorf("unable to load module:\n\troot: " + root + "\n\tmodel:" + model)
	return nil
}

func (sp *SimpleParser) LabelEndpoint(p *seqs.EndpointLabelerParam) string {
	initialStr := sp.self
	matchItems := seqs.FindMatchItems(initialStr)
	for _, item := range matchItems {
		attr := seqs.RemoveWrapper(item)
		var value string
		switch attr {
		case "epname":
			value = p.EndpointName
		case "human":
			value = p.Human
		case "human_sender":
			value = p.HumanSender
		case "needs_int":
			value = p.NeedsInt
		case "args":
			value = p.Args
		case "patterns":
			value = p.Patterns
		case "controls":
			value = p.Controls
		default:
			value = p.Attrs[attr].GetS()
		}
		initialStr = strings.Replace(initialStr, item, value, 1)
	}

	return seqs.RemovePercentSymbol(initialStr)
}

func (sp *SimpleParser) LabelApp(appname, controls string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	matchItems := seqs.FindMatchItems(initialStr)
	for _, item := range matchItems {
		attr := seqs.RemoveWrapper(item)
		var value string
		switch attr {
		case "appname":
			value = appname
		case "controls":
			value = controls
		default:
			value = attrs[attr].GetS()
		}
		initialStr = strings.Replace(initialStr, item, value, 1)
	}

	return seqs.RemovePercentSymbol(initialStr)
}

func (sp *SimpleParser) fmtSeq(epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	matchItems := seqs.FindMatchItems(initialStr)
	for _, item := range matchItems {
		attr := seqs.RemoveWrapper(item)
		var value string
		switch attr {
		case "epname":
			value = epname
		case "eplongname":
			value = eplongname
		default:
			value = attrs[attr].GetS()
		}
		initialStr = strings.Replace(initialStr, item, value, 1)
	}

	return seqs.RemovePercentSymbol(initialStr)
}

func (sp *SimpleParser) fmtOutput(appname, epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	matchItems := seqs.FindMatchItems(initialStr)

	for _, item := range matchItems {
		attr := seqs.RemoveWrapper(item)
		var value string
		switch attr {
		case "appname":
			value = appname
		case "epname":
			value = epname
		case "eplongname":
			value = eplongname
		default:
			value = attrs[attr].GetS()
		}
		initialStr = strings.Replace(initialStr, item, value, 1)
	}

	return seqs.RemovePercentSymbol(initialStr)
}

func constructSimpleParser(former, latter string) *SimpleParser {
	fmtstr := former
	if former == "" {
		fmtstr = latter
	}

	return &SimpleParser{self: fmtstr}
}

func DoConstructSequenceDiagrams(root_model, endpoint_format, app_format, title, plantuml, filter, output string,
	no_activations, verbose, expire_cache, dry_run bool,
	endpoints, apps, modules []string, blackboxes [][]string) {
	mod := loadApp(root_model, modules)
	syslSdFilters, exists := os.LookupEnv("SYSL_SD_FILTERS")
	epFilters := []string{}
	if exists {
		epFilters = append(epFilters, strings.Split(syslSdFilters, ",")...)
	} else {
		epFilters = append(epFilters, "*")
	}

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
				is_continue := false
				for _, filt := range epFilters {
					log.Warn(filt)
				}

				if is_continue {
					continue
				}

				epAttrs := app.GetEndpoints()[k].GetAttrs()
				output_dir := spout.fmtOutput(appName, k, app.GetEndpoints()[k].GetLongName(), epAttrs)
				bbs2 := seqs.TransformBlackBoxes(app.GetEndpoints()[k].GetAttrs()["blackboxes"].GetA().GetElt())
				varrefs := seqs.MergeAttributes(app.GetAttrs(), app.GetEndpoints()[k].GetAttrs())
				sdEndpoints := []string{}
				statements := app.GetEndpoints()[k].GetStmt()
				for _, stmt := range statements {
					parts := stmt.GetCall().GetTarget().GetPart()
					ep := stmt.GetCall().GetEndpoint()
					sdEndpoints = append(sdEndpoints, strings.Join(parts, " :: ")+" <- "+ep)
				}

				sd := &sequenceDiagParam{
					endpoints:   sdEndpoints,
					AppLabeler:  spapp,
					EndpointLabeler: spep,
					title:       spseqtitle.fmtSeq(app.GetEndpoints()[k].GetName(), app.GetEndpoints()[k].GetLongName(), varrefs),
					blackboxes:  append(bbs, bbs2...),
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
			endpoints:   endpoints,
			AppLabeler:  spapp,
			EndpointLabeler: spep,
			title:       title,
			blackboxes:  blackboxes,
		}
		out, _ := generateSequenceDiag(mod, sd)
		seqs.OutputPlantuml(output, plantuml, out)
	}
}

// DoGenerateSequenceDiagrams generate sequence diagrams for the given model
func DoGenerateSequenceDiagrams(stdout, stderr io.Writer, flags *flag.FlagSet, args []string) int {
	var endpoints_flag, apps_flag, blackboxes_flag, modules_flag arrayFlags
	root_model := flags.String("root-model", ".", "sysl root directory for input model file (default: .)")
	flags.Var(&endpoints_flag, "endpoint", "Include endpoint in sequence diagram")
	flags.Var(&apps_flag, "app", "Include all endpoints for app in sequence diagram (currently "+
		"only works with templated --output). Use SYSL_SD_FILTERS env (a "+
		"comma-list of shell globs) to limit the diagrams generated")
	no_activations := flags.Bool("no-activations", true, "Suppress sequence diagram activation bars(default: true)")
	endpoint_format := flags.String("endpoint_format", "%(epname)", "Specify the format string for sequence diagram endpoints. "+
		"May include %%(epname), %%(eplongname) and %%(@foo) for attribute foo(default: %(epname))")
	app_format := flags.String("app_format", "%(appname)", "Specify the format string for sequence diagram participants. "+
		"May include %%(appname) and %%(@foo) for attribute foo(default: %(appname))")
	flags.Var(&blackboxes_flag, "blackbox", "Apps to be treated as black boxes")
	title := flags.String("title", "", "diagram title")
	plantuml := flags.String("plantuml", "", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n"))
	verbose := flags.Bool("verbose", false, "Report each output(default: false)")
	expire_cache := flags.Bool("expire-cache", false, "Expire cache entries to force checking against real destination(default: false)")
	dry_run := flags.Bool("dry-run", false, "Don't perform confluence uploads, but show what would have happened(default: false)")
	filter := flags.String("filter", "", "Only generate diagrams whose output paths match a pattern")
	flags.Var(&modules_flag, "modules", strings.Join([]string{"input files without .sysl extension and with leading /",
		"eg: /project_dir/my_models",
		"combine with --root if needed"}, "\n"))
	output := flags.String("output", "%(epname).png", "output file(default: %(epname).png)")

	err := flags.Parse(args[1:])
	if err != nil {
		log.Errorf("arguments parse error: %v", err)
	}
	log.Warnf("root_model: %s\n", *root_model)
	log.Warnf("endpoints: %v\n", endpoints_flag)
	log.Warnf("app: %v\n", apps_flag)
	log.Warnf("no_activations: %t\n", *no_activations)
	log.Warnf("endpoint_format: %s\n", *endpoint_format)
	log.Warnf("app_format: %s\n", *app_format)
	log.Warnf("blackbox: %s\n", blackboxes_flag)
	log.Warnf("title: %s\n", *title)
	log.Warnf("plantuml: %s\n", *plantuml)
	log.Warnf("verbose: %t\n", *verbose)
	log.Warnf("expire_cache: %t\n", *expire_cache)
	log.Warnf("dry_run: %t\n", *dry_run)
	log.Warnf("filter: %s\n", *filter)
	log.Warnf("modules: %s\n", modules_flag)
	log.Warnf("output: %s\n", *output)

	DoConstructSequenceDiagrams(*root_model, *endpoint_format, *app_format, *title, *plantuml, *filter, *output, *no_activations,
		*verbose, *expire_cache, *dry_run, endpoints_flag, apps_flag, modules_flag, seqs.ParseBlackBoxesFromArgument(blackboxes_flag))

	return 0
}
