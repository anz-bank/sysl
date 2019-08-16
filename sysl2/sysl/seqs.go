package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type sequenceDiagParam struct {
	AppLabeler
	EndpointLabeler
	endpoints  []string
	title      string
	blackboxes map[string]*Upto
	appName    string
	group      string
}

func generateSequenceDiag(m *sysl.Module, p *sequenceDiagParam) (string, error) {
	w := MakeSequenceDiagramWriter(true, "skinparam maxMessageSize 250")
	v := MakeSequenceDiagramVisitor(p.AppLabeler, p.EndpointLabeler, w, m, p.appName, p.group)
	e := MakeEndpointCollectionElement(p.title, p.endpoints, p.blackboxes)

	if err := e.Accept(v); err != nil {
		return "", err
	}

	const color = "#LightBlue"
	for boxname, appset := range v.groupboxes {
		fmt.Fprintf(w, "box \"%s\" %s\n", boxname, color)
		for key := range appset {
			fmt.Fprintf(w, "\tparticipant %s\n", v.UniqueVarForAppName(key))
		}
		fmt.Fprintf(w, "end box\n")
	}

	return w.String(), nil
}

func loadApp(root string, models string) (*sysl.Module, error) {
	// Model we want to generate seqs for, the non-empty model
	mod, err := parse.Parse(models, root)
	if err != nil {
		return nil, fmt.Errorf(
			"loadApp: unable to load module:\n\troot: %s\n\tmodel: %s\nerror: %s",
			root, models, err.Error(),
		)
	}
	return mod, nil
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

func DoConstructSequenceDiagrams(cmdContextParam *CmdContextParamSeqgen) (map[string]string, error) {
	var blackboxes [][]string

	log.Debugf("root: %s\n", *cmdContextParam.root)
	log.Debugf("endpoints: %v\n", cmdContextParam.endpointsFlag)
	log.Debugf("app: %v\n", cmdContextParam.appsFlag)
	log.Debugf("endpoint_format: %s\n", *cmdContextParam.endpointFormat)
	log.Debugf("app_format: %s\n", *cmdContextParam.appFormat)
	log.Debugf("title: %s\n", *cmdContextParam.title)
	log.Debugf("modules: %s\n", *cmdContextParam.modulesFlag)
	log.Debugf("output: %s\n", *cmdContextParam.output)
	log.Debugf("loglevel: %s\n", *cmdContextParam.loglevel)

	if cmdContextParam.plantuml == nil {
		plantuml := os.Getenv("SYSL_PLANTUML")
		cmdContextParam.plantuml = &plantuml
		if *cmdContextParam.plantuml == "" {
			*cmdContextParam.plantuml = "http://localhost:8080/plantuml"
		}
	}
	log.Debugf("plantuml: %s\n", *cmdContextParam.plantuml)

	if cmdContextParam.blackboxes == nil {
		blackboxes = ParseBlackBoxesFromArgument(*cmdContextParam.blackboxesFlag)
		log.Debugf("blackbox: %s\n", *cmdContextParam.blackboxesFlag)
	} else {
		blackboxes = *cmdContextParam.blackboxes
	}

	if *cmdContextParam.isVerbose {
		*cmdContextParam.loglevel = debug
	}
	// Default info
	if level, has := logLevels[*cmdContextParam.loglevel]; has {
		log.SetLevel(level)
	}

	result := make(map[string]string)
	mod, err := loadApp(*cmdContextParam.root, *cmdContextParam.modulesFlag)
	if err != nil {
		return nil, err
	}
	if strings.Contains(*cmdContextParam.output, "%(epname)") {
		if len(blackboxes) > 0 {
			log.Warnf("Ignoring blackboxes passed from command line")
		}
		spout := MakeFormatParser(*cmdContextParam.output)
		for _, appName := range *cmdContextParam.appsFlag {
			app := mod.Apps[appName]
			bbs := TransformBlackBoxes(app.GetAttrs()["blackboxes"].GetA().GetElt())
			spseqtitle := constructFormatParser(app.GetAttrs()["seqtitle"].GetS(), *cmdContextParam.title)
			spep := constructFormatParser(app.GetAttrs()["epfmt"].GetS(), *cmdContextParam.endpointFormat)
			spapp := constructFormatParser(app.GetAttrs()["appfmt"].GetS(), *cmdContextParam.appFormat)
			keys := []string{}
			for k := range app.GetEndpoints() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			bbsAll := map[string]*Upto{}
			TransformBlackboxesToUptos(bbsAll, bbs, BBApplication)
			var sd *sequenceDiagParam
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
				if len(sdEndpoints) == 0 {
					log.Errorf("No call statements to build sequence diagram for endpoint %s", endpoint.Name)
					return result, nil
				}
				groupAttr := epAttrs["groupby"].GetS()
				if len(groupAttr) == 0 {
					groupAttr = *cmdContextParam.group
				} else if len(*cmdContextParam.group) > 0 {
					log.Warnf("Ignoring groupby passed from command line")
				}
				TransformBlackboxesToUptos(bbsAll, bbs2, BBEndpointCollection)
				sd = &sequenceDiagParam{
					endpoints:       sdEndpoints,
					AppLabeler:      spapp,
					EndpointLabeler: spep,
					title:           spseqtitle.FmtSeq(endpoint.GetName(), endpoint.GetLongName(), varrefs),
					blackboxes:      bbsAll,
					appName:         fmt.Sprintf("'%s :: %s'", appName, endpoint.GetName()),
					group:           groupAttr,
				}
				out, _ := generateSequenceDiag(mod, sd)
				for indx := range bbs2 {
					delete(bbsAll, bbs2[indx][0])
				}
				result[outputDir] = out
			}
			for bbKey, bbVal := range bbsAll {
				if bbVal.VisitCount == 0 && bbVal.ValueType == BBApplication {
					log.Warnf("blackbox '%s' not hit in app '%s'\n", bbKey, appName)
				}
			}
		}
	} else {
		if *cmdContextParam.endpointsFlag == nil {
			return result, nil
		}
		spep := constructFormatParser("", *cmdContextParam.endpointFormat)
		spapp := constructFormatParser("", *cmdContextParam.appFormat)
		bbsAll := map[string]*Upto{}
		TransformBlackboxesToUptos(bbsAll, blackboxes, BBCommandLine)
		sd := &sequenceDiagParam{
			endpoints:       *cmdContextParam.endpointsFlag,
			AppLabeler:      spapp,
			EndpointLabeler: spep,
			title:           *cmdContextParam.title,
			blackboxes:      bbsAll,
			group:           *cmdContextParam.group,
		}
		out, _ := generateSequenceDiag(mod, sd)
		for bbKey, bbVal := range bbsAll {
			if bbVal.VisitCount == 0 && bbVal.ValueType == BBCommandLine {
				log.Warnf("blackbox '%s' passed on commandline not hit\n", bbKey)
			}
		}
		result[*cmdContextParam.output] = out
	}

	return result, nil
}

func configureCmdlineForSeqgen(sysl *kingpin.Application, flagmap map[string][]string) *CmdContextParamSeqgen {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	flagmap["sd"] = []string{"root", "endpoint_format", "app_format", "title", "plantuml", "output",
		"groupby", "endpoint", "app"}
	sd := sysl.Command("sd", "Generate sequence diagram")
	returnValues := &CmdContextParamSeqgen{}

	returnValues.root = sd.Flag("root",
		"sysl root directory for input model file (default: .)",
	).Default(".").String()

	returnValues.endpointFormat = sd.Flag("endpoint_format",
		"Specify the format string for sequence diagram endpoints. May include "+
			"%(epname), %(eplongname) and %(@foo) for attribute foo (default: %(epname))",
	).Default("%(epname)").String()

	returnValues.appFormat = sd.Flag("app_format",
		"Specify the format string for sequence diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(appname))",
	).Default("%(appname)").String()

	returnValues.title = sd.Flag("title", "diagram title").Short('t').String()

	returnValues.plantuml = sd.Flag("plantuml",
		"base url of plantuml server (default: $SYSL_PLANTUML or "+
			"http://localhost:8080/plantuml see "+
			"http://plantuml.com/server.html#install for more info)",
	).Short('p').String()

	returnValues.output = sd.Flag("output",
		"output file (default: %(epname).png)",
	).Default("%(epname).png").Short('o').String()

	returnValues.endpointsFlag = sd.Flag("endpoint",
		"Include endpoint in sequence diagram",
	).Short('s').Strings()

	returnValues.appsFlag = sd.Flag("app",
		"Include all endpoints for app in sequence diagram (currently "+
			"only works with templated --output). Use SYSL_SD_FILTERS env (a "+
			"comma-list of shell globs) to limit the diagrams generated",
	).Short('a').Strings()

	returnValues.blackboxesFlag = sd.Flag("blackbox",
		"Input blackboxes in the format App <- Endpoint=Some description, "+
			"repeat '-b App <- Endpoint=Some description' to set multiple blackboxes",
	).Short('b').StringMap()

	returnValues.group = sd.Flag("groupby", "Enter the groupby attribute (apps having "+
		"the same attribute value are grouped together in one box").Short('g').String()

	returnValues.loglevel = sd.Flag("log", "log level[debug,info,warn,off]").Default("info").String()
	returnValues.isVerbose = sd.Flag("verbose", "show output").Short('v').Default("false").Bool()

	returnValues.modulesFlag = sd.Arg("modules",
		"input files without .sysl extension and with leading /, eg: "+
			"/project_dir/my_models combine with --root if needed",
	).String()

	return returnValues
}
