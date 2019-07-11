package main

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type sequenceDiagParam struct {
	AppLabeler
	EndpointLabeler
	endpoints  []string
	title      string
	blackboxes map[string]Upto
	appName    string
}

func generateSequenceDiag(m *sysl.Module, p *sequenceDiagParam) (string, error) {
	w := MakeSequenceDiagramWriter(true, "skinparam maxMessageSize 250")
	v := MakeSequenceDiagramVisitor(p.AppLabeler, p.EndpointLabeler, w, m, p.appName)
	e := MakeEndpointCollectionElement(p.title, p.endpoints, p.blackboxes)
	if err := e.Accept(v); err != nil {
		return "", err
	}

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

func constructFormatParser(former, latter string) *FormatParser {
	fmtstr := former
	if former == "" {
		fmtstr = latter
	}

	return MakeFormatParser(escapeWordBoundary(fmtstr))
}

func escapeWordBoundary(src string) string {
	result, _ := json.Marshal(src)
	escapeStr := strings.Replace(string(result), `\u0008`, `\\b`, -1)
	var val string
	if err := json.Unmarshal([]byte(escapeStr), &val); err != nil {
		panic(err)
	}

	return val
}

func DoConstructSequenceDiagrams(
	stdout, stderr io.Writer, root_model, endpoint_format, app_format, title, output, modules string,
	endpoints, apps []string,
	blackboxes [][]string,
) map[string]string {
	result := make(map[string]string)
	mod := loadApp(rootModel, modules)
	if mod == nil {
		return result
	}
	if strings.Contains(output, "%(epname)") {
		spout := MakeFormatParser(output)
		for _, appName := range apps {
			app := mod.Apps[appName]
			bbs := TransformBlackBoxes(app.GetAttrs()["blackboxes"].GetA().GetElt())
			spseqtitle := constructFormatParser(app.GetAttrs()["seqtitle"].GetS(), title)
			spep := constructFormatParser(app.GetAttrs()["epfmt"].GetS(), endpointFormat)
			spapp := constructFormatParser(app.GetAttrs()["appfmt"].GetS(), appFormat)
			keys := []string{}
			for k := range app.GetEndpoints() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				endpoint := app.GetEndpoints()[k]
				epAttrs := endpoint.GetAttrs()
				outputDir := spout.FmtOutput(appName, k, endpoint.GetLongName(), epAttrs)
				bbs2 := TransformBlackBoxes(endpoint.GetAttrs()["blackboxes"].GetA().GetElt())
				varrefs := MergeAttributes(app.GetAttrs(), endpoint.GetAttrs())
				sdEndpoints := []string{}
				for _, stmt := range endpoint.GetStmt() {
					_, ok := stmt.Stmt.(*sysl.Statement_Call)
					if ok {
						parts := stmt.GetCall().GetTarget().GetPart()
						ep := stmt.GetCall().GetEndpoint()
						sdEndpoints = append(sdEndpoints, strings.Join(parts, " :: ")+" <- "+ep)
					}
				}

				bbsAll := seqs.TransformBlackboxesToUptos(bbs, seqs.BB_APPLICATION)
				bbsEndpoint := seqs.TransformBlackboxesToUptos(bbs2, seqs.BB_ENDPOINT_COLLECTION)
				for k, v := range bbsEndpoint {
					bbsAll[k] = v
				}
				sd := &sequenceDiagParam{
					endpoints:       sdEndpoints,
					AppLabeler:      spapp,
					EndpointLabeler: spep,
					title:           spseqtitle.FmtSeq(endpoint.GetName(), endpoint.GetLongName(), varrefs),
					blackboxes:      bbsAll,
					appName:         appName,
				}
				out, _ := generateSequenceDiag(stdout, stderr, mod, sd)
				result[output_dir] = out
			}
		}
	} else {
		if endpoints == nil {
			return result
		}
		spep := constructFormatParser("", endpointFormat)
		spapp := constructFormatParser("", appFormat)
		sd := &sequenceDiagParam{
			endpoints:       endpoints,
			AppLabeler:      spapp,
			EndpointLabeler: spep,
			title:           title,
			blackboxes:      seqs.TransformBlackboxesToUptos(blackboxes, seqs.BB_APPLICATION),
			appName:         "Global",
		}
		out, _ := generateSequenceDiag(stdout, stderr, mod, sd)
		result[output] = out
	}

	return result
}

// DoGenerateSequenceDiagrams generate sequence diagrams for the given model
// TODO: Remove nolint when code that uses stdout and stderr is added.
//nolint:unparam
func DoGenerateSequenceDiagrams(stdout, stderr io.Writer, args []string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	sd := kingpin.New("sd", "Generate sequence diagram")

	root := sd.Flag("root",
		"sysl root directory for input model file (default: .)",
	).Default(".").String()

	endpointFormat := sd.Flag("endpoint_format",
		"Specify the format string for sequence diagram endpoints. May include "+
			"%(epname), %(eplongname) and %(@foo) for attribute foo (default: %(epname))",
	).Default("%(epname)").String()

	appFormat := sd.Flag("app_format",
		"Specify the format string for sequence diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(appname))",
	).Default("%(appname)").String()

	title := sd.Flag("title", "diagram title").Short('t').String()

	plantuml := sd.Flag("plantuml",
		"base url of plantuml server (default: $SYSL_PLANTUML or "+
			"http://localhost:8080/plantuml see "+
			"http://plantuml.com/server.html#install for more info)",
	).Short('p').String()

	output := sd.Flag("output",
		"output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').String()

	endpointsFlag := sd.Flag("endpoint",
		"Include endpoint in sequence diagram",
	).Short('s').Strings()

	appsFlag := sd.Flag("app",
		"Include all endpoints for app in sequence diagram (currently "+
			"only works with templated --output). Use SYSL_SD_FILTERS env (a "+
			"comma-list of shell globs) to limit the diagrams generated",
	).Short('a').Strings()

	blackboxesFlag := sd.Flag("blackbox", "Apps to be treated as black boxes").Short('b').StringMap()

	loglevel := sd.Flag("log", "log level[debug,info,warn,off]").Default("warn").String()

	modulesFlag := sd.Arg("modules",
		"input files without .sysl extension and with leading /, eg: "+
			"/project_dir/my_models combine with --root if needed",
	).String()

	if _, err := sd.Parse(args[1:]); err != nil {
		return err
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
	log.Debugf("endpoints: %v\n", endpointsFlag)
	log.Debugf("app: %v\n", appsFlag)
	log.Debugf("endpoint_format: %s\n", *endpointFormat)
	log.Debugf("app_format: %s\n", *appFormat)
	log.Debugf("blackbox: %s\n", *blackboxesFlag)
	log.Debugf("title: %s\n", *title)
	log.Debugf("plantuml: %s\n", *plantuml)
	log.Debugf("modules: %s\n", *modulesFlag)
	log.Debugf("output: %s\n", *output)
	log.Debugf("loglevel: %s\n", *loglevel)

	result := DoConstructSequenceDiagrams(*root, *endpointFormat, *appFormat, *title, *output, *modulesFlag,
		*endpointsFlag, *appsFlag, ParseBlackBoxesFromArgument(*blackboxesFlag))
	for k, v := range result {
		if err := OutputPlantuml(k, *plantuml, v); err != nil {
			return err
		}
	}
	return nil
}
