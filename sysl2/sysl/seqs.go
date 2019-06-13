package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/seqs"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
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
	sd := kingpin.New("sd", "Generate sequence diagram")
	root := sd.Flag("root", "sysl root directory for input model file (default: .)").Default(".").String()
	endpoint_format := sd.Flag("endpoint_format", "Specify the format string for sequence diagram endpoints. "+
		"May include %(epname), %(eplongname) and %(@foo) for attribute foo(default: %(epname))").Default("%(epname)").String()
	app_format := sd.Flag("app_format", "Specify the format string for sequence diagram participants. "+
		"May include %%(appname) and %%(@foo) for attribute foo(default: %(appname))").Default("%(appname)").String()
	title := sd.Flag("title", "diagram title").Short('t').String()
	plantuml := sd.Flag("plantuml", strings.Join([]string{"base url of plantuml server",
		"(default: $SYSL_PLANTUML or http://localhost:8080/plantuml",
		"see http://plantuml.com/server.html#install for more info)"}, "\n")).Short('p').String()
	output := sd.Flag("output", "output file(default: %(epname).png)").Default("%(epname).png").Short('o').String()
	endpoints_flag := sd.Flag("endpoint", "Include endpoint in sequence diagram").Short('s').Strings()
	apps_flag := sd.Flag("app", "Include all endpoints for app in sequence diagram (currently "+
		"only works with templated --output). Use SYSL_SD_FILTERS env (a "+
		"comma-list of shell globs) to limit the diagrams generated").Short('a').Strings()
	blackboxes_flag := sd.Flag("blackbox", "Apps to be treated as black boxes").Strings()
	loglevel := sd.Flag("log", "log level[debug,info,warn,off]").Default("warn").String()

	modules_flag := sd.Arg("modules", strings.Join([]string{"input files without .sysl extension and with leading /",
		"eg: /project_dir/my_models",
		"combine with --root if needed"}, "\n")).String()

	_, err := sd.Parse(args[1:])

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
	log.Debugf("endpoints: %v\n", endpoints_flag)
	log.Debugf("app: %v\n", apps_flag)
	log.Debugf("endpoint_format: %s\n", *endpoint_format)
	log.Debugf("app_format: %s\n", *app_format)
	log.Debugf("blackbox: %s\n", *blackboxes_flag)
	log.Debugf("title: %s\n", *title)
	log.Debugf("plantuml: %s\n", *plantuml)
	log.Debugf("modules: %s\n", *modules_flag)
	log.Debugf("output: %s\n", *output)
	log.Debugf("loglevel: %s\n", *loglevel)

	result := DoConstructSequenceDiagrams(*root, *endpoint_format, *app_format, *title, *output, *modules_flag,
		*endpoints_flag, *apps_flag, seqs.ParseBlackBoxesFromArgument(*blackboxes_flag))
	for k, v := range result {
		seqs.OutputPlantuml(k, *plantuml, v)
	}
}
