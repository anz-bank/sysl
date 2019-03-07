package main

import (
	"os"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/seqs"
	"github.com/sirupsen/logrus"
)

<<<<<<< HEAD
type sequenceDiagParam struct {
	seqs.Labeler
	endpoints  []string
	title      string
	blackboxes [][]string
=======
type SimpleParser struct {
	self string
}

type epFmtParam struct {
	epname, human, human_sender, needs_int, args, patterns, controls string
	attrs                                                            map[string]*sysl.Attribute
>>>>>>> Sysl generate sequence diagram
}

func generateSequenceDiag(m *sysl.Module, p *sequenceDiagParam) (string, error) {
	w := seqs.MakeSequenceDiagramWriter(true, "skinparam maxMessageSize 250")

	v := seqs.MakeSequenceDiagramVisitor(p.Labeler, w, m)

	e := seqs.MakeEndpointCollectionElement(p.title, p.endpoints, p.blackboxes)

	v.Visit(e)

	return w.String(), nil
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
	if blackboxes != nil {
		for _, vals := range blackboxes {
			sub_bbs := []string{}
			for _, val := range vals.GetA().Elt {
				sub_bbs = append(sub_bbs, val.GetS())
			}
			bbs = append(bbs, sub_bbs)
		}
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

	return ""
}

func (sp *SimpleParser) fmtApp(appname, controls string, attrs map[string]*sysl.Attribute) string {

	return ""
}

func (sp *SimpleParser) fmtSeq() string {

	return ""
}

func (sp *SimpleParser) fmtOutput() string {

	return ""
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
		//var out_fmt interface{}
		for _, appName := range apps {
			app := mod.Apps[appName]

			bbs := transformBlackBoxes(app.Attrs["blackboxes"].GetA().GetElt())
			logrus.Warnf("bbs: %v", bbs)

			//var seqtitle interface{}

			//var epfmt interface{}

			//var appfmt interface{}

			keys := []string{}
			for k := range app.GetEndpoints() {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				logrus.Warnf("key: %s", k)
				logrus.Warnf("val: %v", app.GetEndpoints()[k])

				is_continue := false
				for _, filt := range epFilters {
					logrus.Warn(filt)
				}

				if is_continue {
					continue
				}

				//attrs :=
				bbs2 := transformBlackBoxes(app.GetEndpoints()[k].GetAttrs()["blackboxes"].GetA().GetElt())
				logrus.Warnf("bbs2: %v", bbs2)
				logrus.Warnf("union bbs: %v", append(bbs, bbs2...))

				//varrefs := map[string]string{}

				sdEndpoints := []string{}
				statements := app.GetEndpoints()[k].GetStmt()
				logrus.Warnf("end points: %v", app.GetEndpoints())
				for _, stmt := range statements {
					parts := stmt.GetCall().GetTarget().GetPart()
					ep := stmt.GetCall().GetEndpoint()
					sdEndpoints = append(sdEndpoints, strings.Join(parts, " :: ")+" <- "+ep)
				}

				spEp := &SimpleParser{self: "%(appname)"}
				sd := &sequenceDiagParam{
<<<<<<< HEAD
					endpoints:   sdEndpoints,
					epfmt:       SfmtEP(fmtEp),
					appfmt:      SfmtApp(fmtApp),
=======
					endpoints: sdEndpoints,
					epfmt: SfmtEP(spEp.fmtEp),
					appfmt: SfmtApp(spEp.fmtApp),
>>>>>>> Sysl generate sequence diagram
					activations: no_activations,
					title:       "",
					blackboxes:  append(bbs, bbs2...),
				}
				logrus.Warnf("sd: %v", sd)
			}
		}
	} else {
		if endpoints == nil {
			return
		}

		spEp := &SimpleParser{self: "%(appname)"}
		sd := &sequenceDiagParam{
<<<<<<< HEAD
			endpoints:   endpoints,
			epfmt:       SfmtEP(fmtEp),
			appfmt:      SfmtApp(fmtApp),
=======
			endpoints: endpoints,
			epfmt: SfmtEP(spEp.fmtEp),
			appfmt: SfmtApp(spEp.fmtApp),
>>>>>>> Sysl generate sequence diagram
			activations: no_activations,
			title:       title,
			blackboxes:  blackboxes,
		}

		logrus.Warnf("sd: %v", sd)
	}
}
