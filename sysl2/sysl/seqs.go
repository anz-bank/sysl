package main

import (
	"encoding/json"
	"flag"
	"io"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/seqs"
	log "github.com/sirupsen/logrus"
)

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

func loadApp(root string, models string) *sysl.Module {
	// Model we want to generate seqs for, the non-empty model
	mod, err := Parse(models, root)
	if err == nil {
		return mod
	}
	log.Errorf("unable to load module:\n\troot: " + root + "\n\tmodel:" + models)
	return nil
}

func constructFormatParser(former, latter string) *seqs.FormatParser {
	fmtstr := former
	if former == "" {
		fmtstr = latter
	}

	return seqs.MakeFormatParser(escapeWordBoundary(fmtstr))
}

func escapeWordBoundary(src string) string {
	result, _ := json.Marshal(src)
	escapeStr := strings.Replace(string(result), `\u0008`, `\\b`, -1)
	var val string
	json.Unmarshal([]byte(escapeStr), &val)

	return val
}

func DoConstructSequenceDiagrams(
	root_model, endpoint_format, app_format, title, output, modules string,
	endpoints, apps []string,
	blackboxes [][]string,
) map[string]string {
	result := make(map[string]string)
	mod := loadApp(root_model, modules)
	if mod == nil {
		return result
	}
	if strings.Contains(output, "%(epname)") {
		spout := seqs.MakeFormatParser(output)
		for _, appName := range apps {
			app := mod.Apps[appName]
			bbs := seqs.TransformBlackBoxes(app.GetAttrs()["blackboxes"].GetA().GetElt())
			spseqtitle := constructFormatParser(app.GetAttrs()["seqtitle"].GetS(), title)
			spep := constructFormatParser(app.GetAttrs()["epfmt"].GetS(), endpoint_format)
			spapp := constructFormatParser(app.GetAttrs()["appfmt"].GetS(), app_format)
			keys := []string{}
			for k := range app.GetEndpoints() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				endpoint := app.GetEndpoints()[k]
				epAttrs := endpoint.GetAttrs()
				output_dir := spout.FmtOutput(appName, k, endpoint.GetLongName(), epAttrs)
				bbs2 := seqs.TransformBlackBoxes(endpoint.GetAttrs()["blackboxes"].GetA().GetElt())
				varrefs := seqs.MergeAttributes(app.GetAttrs(), endpoint.GetAttrs())
				sdEndpoints := []string{}
				for _, stmt := range endpoint.GetStmt() {
					_, ok := stmt.Stmt.(*sysl.Statement_Call)
					if ok {
						parts := stmt.GetCall().GetTarget().GetPart()
						ep := stmt.GetCall().GetEndpoint()
						sdEndpoints = append(sdEndpoints, strings.Join(parts, " :: ")+" <- "+ep)
					}
				}

				sd := &sequenceDiagParam{
					endpoints:       sdEndpoints,
					AppLabeler:      spapp,
					EndpointLabeler: spep,
					title:           spseqtitle.FmtSeq(endpoint.GetName(), endpoint.GetLongName(), varrefs),
					blackboxes:      append(bbs, bbs2...),
				}
				out, _ := generateSequenceDiag(mod, sd)
				result[output_dir] = out
			}
		}
	} else {
		if endpoints == nil {
			return result
		}
		spep := constructFormatParser("", endpoint_format)
		spapp := constructFormatParser("", app_format)
		sd := &sequenceDiagParam{
			endpoints:       endpoints,
			AppLabeler:      spapp,
			EndpointLabeler: spep,
			title:           title,
			blackboxes:      blackboxes,
		}
		out, _ := generateSequenceDiag(mod, sd)
		result[output] = out
	}

	return result
}

// DoGenerateSequenceDiagrams generate sequence diagrams for the given model
func DoGenerateSequenceDiagrams(stdout, stderr io.Writer, flags *flag.FlagSet, args []string) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	var endpoints_flag, apps_flag, blackboxes_flag arrayFlags
	root := flags.String("root", ".", "sysl root directory for input model file (default: .)")
	endpoint_format := flags.String("endpoint_format", "%(epname)", "Specify the format string for sequence diagram endpoints. "+
		"May include %%(epname), %%(eplongname) and %%(@foo) for attribute foo(default: %(epname))")
	app_format := flags.String("app_format", "%(appname)", "Specify the format string for sequence diagram participants. "+
		"May include %%(appname) and %%(@foo) for attribute foo(default: %(appname))")
	title := flags.String("title", "", "diagram title")
	plantuml := flags.String("plantuml", "", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n"))
	output := flags.String("output", "%(epname).png", "output file(default: %(epname).png)")
	modules_flag := flags.String("modules", ".", strings.Join([]string{"input files without .sysl extension and with leading /",
		"eg: /project_dir/my_models",
		"combine with --root if needed"}, "\n"))
	flags.Var(&endpoints_flag, "endpoint", "Include endpoint in sequence diagram")
	flags.Var(&apps_flag, "app", "Include all endpoints for app in sequence diagram (currently "+
		"only works with templated --output). Use SYSL_SD_FILTERS env (a "+
		"comma-list of shell globs) to limit the diagrams generated")
	flags.Var(&blackboxes_flag, "blackbox", "Apps to be treated as black boxes")

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
	log.Debugf("root: %s\n", *root)
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
	log.Debugf("modules: %s\n", *modules_flag)
	log.Debugf("output: %s\n", *output)

	result := DoConstructSequenceDiagrams(*root, *endpoint_format, *app_format, *title, *output, *modules_flag,
		endpoints_flag, apps_flag, seqs.ParseBlackBoxesFromArgument(blackboxes_flag))
	for k, v := range result {
		seqs.OutputPlantuml(k, *plantuml, v)
	}
}
