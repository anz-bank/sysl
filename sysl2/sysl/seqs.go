package main

import (
	"flag"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
)

type SimpleParser struct {
	self string
}

type epFmtParam struct {
	epname, human, human_sender, needs_int, args, patterns, controls string
	attrs                                                            map[string]*sysl.Attribute
}

type SfmtApp = func(appname, controls string, attrs map[string]*sysl.Attribute) string

type SfmtEP = func(p *epFmtParam) string

type sequenceDiagParam struct {
	endpoints   []string
	epfmt       SfmtEP
	appfmt      SfmtApp
	activations bool
	title       string
	blackboxes  [][]string
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

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
	logrus.Errorf("unable to load module:\n\troot: " + root + "\n\tmodel:" + model)
	return nil
}

func transformBlackBoxes(blackboxes []*sysl.Attribute) [][]string {
	bbs := [][]string{}
	logrus.Warnf("transform blackboxes: %v", blackboxes)
	for _, vals := range blackboxes {
		sub_bbs := []string{}
		for _, val := range vals.GetA().Elt {
			sub_bbs = append(sub_bbs, val.GetS())
		}
		bbs = append(bbs, sub_bbs)
	}

	return bbs
}

func parseBlackBoxesFromArgument(blackboxFlags []string) [][]string {
	bbs := [][]string{}
	for _, blackboxFlag := range blackboxFlags {
		sub_bbs := []string{}
		sub_bbs = append(sub_bbs, strings.Split(blackboxFlag, ",")...)
		bbs = append(bbs, sub_bbs)
	}

	return bbs
}

func (sp *SimpleParser) fmtEp(p *epFmtParam) string {
	initialStr := sp.self
	matchItems := findMatchItems(initialStr)
	for _, item := range matchItems {
		attr := removeWrapper(item)
		var value string
		switch attr {
		case "epname":
			value = p.epname
		case "human":
			value = p.human
		case "human_sender":
			value = p.human_sender
		case "needs_int":
			value = p.needs_int
		case "args":
			value = p.args
		case "patterns":
			value = p.patterns
		case "controls":
			value = p.controls
		default:
			value = p.attrs[attr].GetS()
		}
		initialStr = strings.Replace(initialStr, item, value, 1)
	}

	return removePercentSymbol(initialStr)
}

func (sp *SimpleParser) fmtApp(appname, controls string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	matchItems := findMatchItems(initialStr)
	for _, item := range matchItems {
		attr := removeWrapper(item)
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

	return removePercentSymbol(initialStr)
}

func (sp *SimpleParser) fmtSeq(epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	matchItems := findMatchItems(initialStr)
	for _, item := range matchItems {
		attr := removeWrapper(item)
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

	return removePercentSymbol(initialStr)
}

func (sp *SimpleParser) fmtOutput(appname, epname, eplongname string, attrs map[string]*sysl.Attribute) string {
	initialStr := sp.self
	matchItems := findMatchItems(initialStr)

	for _, item := range matchItems {
		attr := removeWrapper(item)
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

	return removePercentSymbol(initialStr)
}

func findMatchItems(origin string) []string {
	re := regexp.MustCompile(`(%\(\w+\))`)
	return re.FindAllString(origin, -1)
}

func removeWrapper(origin string) string {
	replaced := strings.Replace(origin, "%(", "", 1)
	replaced = strings.Replace(replaced, ")", "", 1)
	return replaced
}

func removePercentSymbol(origin string) string {
	return strings.Replace(origin, "%", "", -1)
}

func constructSimpleParser(former, latter string) *SimpleParser {
	fmtstr := former
	if former == "" {
		fmtstr = latter
	}

	return &SimpleParser{self: fmtstr}
}

func mergeAttributes(app, edpnt map[string]*sysl.Attribute) map[string]*sysl.Attribute {
	result := make(map[string]*sysl.Attribute)
	for k, v := range app {
		result[k] = v
	}
	for k, v := range edpnt {
		result[k] = v
	}

	return result
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
	logrus.Warnf("mod: %v", mod)

	if strings.Contains(output, "%(epname)") {
		spout := &SimpleParser{self: output}
		for _, appName := range apps {
			app := mod.Apps[appName]

			bbs := transformBlackBoxes(app.GetAttrs()["blackboxes"].GetA().GetElt())
			logrus.Warnf("bbs: %v", bbs)

			spseqtitle := constructSimpleParser(app.GetAttrs()["seqtitle"].GetS(), title)

			spep := constructSimpleParser(app.GetAttrs()["epfmt"].GetS(), endpoint_format)

			spapp := constructSimpleParser(app.GetAttrs()["appfmt"].GetS(), app_format)

			keys := []string{}
			for k := range app.GetEndpoints() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys  {
				logrus.Warnf("key: %s", k)
				logrus.Warnf("val: %v", app.GetEndpoints()[k])

				is_continue := false
				for _, filt := range epFilters {
					logrus.Warn(filt)
				}

				if is_continue {
					continue
				}

				epAttrs := app.GetEndpoints()[k].GetAttrs()

				output_dir := spout.fmtOutput(appName, k, app.GetEndpoints()[k].GetLongName(), epAttrs)

				bbs2 := transformBlackBoxes(app.GetEndpoints()[k].GetAttrs()["blackboxes"].GetA().GetElt())
				logrus.Warnf("bbs2: %v", bbs2)
				logrus.Warnf("union bbs: %v", append(bbs, bbs2...))

				varrefs := mergeAttributes(app.GetAttrs(), app.GetEndpoints()[k].GetAttrs())


				sdEndpoints := []string{}
				statements := app.GetEndpoints()[k].GetStmt()
				logrus.Warnf("end points: %v", app.GetEndpoints())
				for _, stmt := range statements {
					parts := stmt.GetCall().GetTarget().GetPart()
					ep := stmt.GetCall().GetEndpoint()
					sdEndpoints = append(sdEndpoints, strings.Join(parts, " :: ") + " <- " + ep)
				}

				sd := &sequenceDiagParam{
					endpoints: sdEndpoints,
					epfmt: SfmtEP(spep.fmtEp),
					appfmt: SfmtApp(spapp.fmtApp),
					activations: no_activations,
					title: spseqtitle.fmtSeq(app.GetEndpoints()[k].GetName(), app.GetEndpoints()[k].GetLongName(), varrefs),
					blackboxes: append(bbs, bbs2...),
				}
				out, _ := generateSequenceDiag(mod, sd)
				logrus.Warnf("out: %s", out)

				logrus.Warnf("output_dir: %s", output_dir)
				OutputPlantuml(output_dir, plantuml, out)
			}
		}
	} else {
		if endpoints == nil {
			return
		}

		spEp := &SimpleParser{self: "%(appname)"}
		sd := &sequenceDiagParam{
			endpoints: endpoints,
			epfmt: SfmtEP(spEp.fmtEp),
			appfmt: SfmtApp(spEp.fmtApp),
			activations: no_activations,
			title: title,
			blackboxes: blackboxes,
		}
		out, _ := generateSequenceDiag(mod, sd)

		logrus.Warnf("sd: %v", sd)
		OutputPlantuml(output, plantuml, out)
	}
}

// DoGenerateSequenceDiagrams generate sequence diagrams for the given model
func DoGenerateSequenceDiagrams(stdout, stderr io.Writer, flags *flag.FlagSet, args []string) int {
	var endpoints_flag, apps_flag, blackboxes_flag, modules_flag arrayFlags
	root_model := flags.String("root-model", ".", "sysl root directory for input model file (default: .)")
	flags.Var(&endpoints_flag, "endpoint", "Include endpoint in sequence diagram")
	flags.Var(&apps_flag, "app", "Include all endpoints for app in sequence diagram (currently " +
		"only works with templated --output). Use SYSL_SD_FILTERS env (a " +
		"comma-list of shell globs) to limit the diagrams generated")
	no_activations := flags.Bool("no-activations", true, "Suppress sequence diagram activation bars(default: true)")
	endpoint_format := flags.String("endpoint_format", "%(epname)", "Specify the format string for sequence diagram endpoints. " +
		"May include %%(epname), %%(eplongname) and %%(@foo) for attribute foo(default: %(epname))")
	app_format := flags.String("app_format", "%(appname)", "Specify the format string for sequence diagram participants. " +
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

	flags.Parse(args[1:])
	logrus.Warnf("root_model: %s\n", *root_model)
	logrus.Warnf("endpoints: %v\n", endpoints_flag)
	logrus.Warnf("app: %v\n", apps_flag)
	logrus.Warnf("no_activations: %t\n", *no_activations)
	logrus.Warnf("endpoint_format: %s\n", *endpoint_format)
	logrus.Warnf("app_format: %s\n", *app_format)
	logrus.Warnf("blackbox: %s\n", blackboxes_flag)
	logrus.Warnf("title: %s\n", *title)
	logrus.Warnf("plantuml: %s\n", *plantuml)
	logrus.Warnf("verbose: %t\n", *verbose)
	logrus.Warnf("expire_cache: %t\n", *expire_cache)
	logrus.Warnf("dry_run: %t\n", *dry_run)
	logrus.Warnf("filter: %s\n", *filter)
	logrus.Warnf("modules: %s\n", modules_flag)
	logrus.Warnf("output: %s\n", *output)

	DoConstructSequenceDiagrams(*root_model, *endpoint_format, *app_format, *title, *plantuml, *filter, *output, *no_activations,
		*verbose, *expire_cache, *dry_run, endpoints_flag, apps_flag, modules_flag, parseBlackBoxesFromArgument(blackboxes_flag))

	return 0
}
