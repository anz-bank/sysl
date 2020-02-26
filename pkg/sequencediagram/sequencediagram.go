package sequencediagram

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
)

type SequenceDiagParam struct {
	cmdutils.AppLabeler
	cmdutils.EndpointLabeler
	Endpoints  []string
	Title      string
	Blackboxes map[string]*cmdutils.Upto
	AppName    string
	Group      string
}

func GenerateSequenceDiag(m *sysl.Module, p *SequenceDiagParam, logger *logrus.Logger) (string, error) {
	w := cmdutils.MakeSequenceDiagramWriter(true, "skinparam maxMessageSize 250")
	v := cmdutils.MakeSequenceDiagramVisitor(p.AppLabeler, p.EndpointLabeler, w, m, p.AppName, p.Group, logger)
	e := cmdutils.MakeEndpointCollectionElement(p.Title, p.Endpoints, p.Blackboxes)

	if err := e.Accept(v); err != nil {
		return "", err
	}

	const color = "#LightBlue"
	for boxname, appset := range v.Groupboxes {
		fmt.Fprintf(w, "box \"%s\" %s\n", boxname, color)
		for key := range appset {
			fmt.Fprintf(w, "\tparticipant %s\n", v.UniqueVarForAppName(key))
		}
		fmt.Fprintf(w, "end box\n")
	}

	return w.String(), nil
}

func ConstructFormatParser(former, latter string) *cmdutils.FormatParser {
	fmtstr := former
	if former == "" {
		fmtstr = latter
	}

	return cmdutils.MakeFormatParser(EscapeWordBoundary(fmtstr))
}

func EscapeWordBoundary(src string) string {
	result, err := json.Marshal(src)
	syslutil.PanicOnError(err)
	escapeStr := strings.Replace(string(result), `\u0008`, `\\b`, -1)
	var val string
	err = json.Unmarshal([]byte(escapeStr), &val)
	syslutil.PanicOnError(err)

	return val
}

func DoConstructSequenceDiagrams(
	cmdContextParam *cmdutils.CmdContextParamSeqgen,
	model *sysl.Module,
	logger *logrus.Logger,
) (map[string]string, error) {
	var blackboxes [][]string

	logger.Debugf("endpoints: %v\n", cmdContextParam.EndpointsFlag)
	logger.Debugf("app: %v\n", cmdContextParam.AppsFlag)
	logger.Debugf("endpoint_format: %s\n", cmdContextParam.EndpointFormat)
	logger.Debugf("app_format: %s\n", cmdContextParam.AppFormat)
	logger.Debugf("title: %s\n", cmdContextParam.Title)
	logger.Debugf("output: %s\n", cmdContextParam.Output)

	if len(cmdContextParam.Blackboxes) == 0 {
		blackboxes = cmdutils.ParseBlackBoxesFromArgument(cmdContextParam.BlackboxesFlag)
		logger.Debugf("blackbox: %s\n", cmdContextParam.BlackboxesFlag)
	} else {
		blackboxes = cmdContextParam.Blackboxes
	}

	result := make(map[string]string)

	if strings.Contains(cmdContextParam.Output, "%(epname)") {
		if len(blackboxes) > 0 {
			logger.Warnf("Ignoring blackboxes passed from command line")
		}
		spout := cmdutils.MakeFormatParser(cmdContextParam.Output)
		for _, appName := range cmdContextParam.AppsFlag {
			app := model.Apps[appName]
			bbs := cmdutils.TransformBlackBoxes(app.GetAttrs()["blackboxes"].GetA().GetElt())
			spseqtitle := ConstructFormatParser(app.GetAttrs()["seqtitle"].GetS(), cmdContextParam.Title)
			spep := ConstructFormatParser(app.GetAttrs()["epfmt"].GetS(), cmdContextParam.EndpointFormat)
			spapp := ConstructFormatParser(app.GetAttrs()["appfmt"].GetS(), cmdContextParam.AppFormat)
			keys := []string{}
			for k := range app.GetEndpoints() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			bbsAll := map[string]*cmdutils.Upto{}
			cmdutils.TransformBlackboxesToUptos(bbsAll, bbs, cmdutils.BBApplication)
			var sd *SequenceDiagParam
			for _, k := range keys {
				endpoint := app.GetEndpoints()[k]
				epAttrs := endpoint.GetAttrs()
				outputDir := spout.FmtOutput(appName, k, endpoint.GetLongName(), epAttrs)
				bbs2 := cmdutils.TransformBlackBoxes(endpoint.GetAttrs()["blackboxes"].GetA().GetElt())
				varrefs := cmdutils.MergeAttributes(app.GetAttrs(), endpoint.GetAttrs())
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
					return nil, fmt.Errorf("no call statements to build sequence diagram for endpoint %s",
						endpoint.Name)
				}
				groupAttr := epAttrs["groupby"].GetS()
				if len(groupAttr) == 0 {
					groupAttr = cmdContextParam.Group
				} else if len(cmdContextParam.Group) > 0 {
					logger.Warnf("Ignoring groupby passed from command line")
				}
				cmdutils.TransformBlackboxesToUptos(bbsAll, bbs2, cmdutils.BBEndpointCollection)
				sd = &SequenceDiagParam{
					Endpoints:       sdEndpoints,
					AppLabeler:      spapp,
					EndpointLabeler: spep,
					Title:           spseqtitle.FmtSeq(endpoint.GetName(), endpoint.GetLongName(), varrefs),
					Blackboxes:      bbsAll,
					AppName:         fmt.Sprintf("'%s :: %s'", appName, endpoint.GetName()),
					Group:           groupAttr,
				}
				out, err := GenerateSequenceDiag(model, sd, logger)
				if err != nil {
					return nil, err
				}
				for indx := range bbs2 {
					delete(bbsAll, bbs2[indx][0])
				}
				result[outputDir] = out
			}
			for bbKey, bbVal := range bbsAll {
				if bbVal.VisitCount == 0 && bbVal.ValueType == cmdutils.BBApplication {
					logger.Warnf("blackbox '%s' not hit in app '%s'\n", bbKey, appName)
				}
			}
		}
	} else {
		if len(cmdContextParam.EndpointsFlag) == 0 {
			return result, nil
		}
		spep := ConstructFormatParser("", cmdContextParam.EndpointFormat)
		spapp := ConstructFormatParser("", cmdContextParam.AppFormat)
		bbsAll := map[string]*cmdutils.Upto{}
		cmdutils.TransformBlackboxesToUptos(bbsAll, blackboxes, cmdutils.BBCommandLine)
		sd := &SequenceDiagParam{
			Endpoints:       cmdContextParam.EndpointsFlag,
			AppLabeler:      spapp,
			EndpointLabeler: spep,
			Title:           cmdContextParam.Title,
			Blackboxes:      bbsAll,
			Group:           cmdContextParam.Group,
		}
		out, err := GenerateSequenceDiag(model, sd, logger)
		if err != nil {
			return nil, err
		}
		for bbKey, bbVal := range bbsAll {
			if bbVal.VisitCount == 0 && bbVal.ValueType == cmdutils.BBCommandLine {
				logger.Warnf("blackbox '%s' passed on commandline not hit\n", bbKey)
			}
		}
		result[cmdContextParam.Output] = out
	}

	return result, nil
}
