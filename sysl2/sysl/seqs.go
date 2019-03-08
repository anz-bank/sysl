package main

import (
	"os"
	"regexp"
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
					sdEndpoints = append(sdEndpoints, strings.Join(parts, " :: ")+" <- "+ep)
				}

				sd := &sequenceDiagParam{
<<<<<<< HEAD
					endpoints:   sdEndpoints,
					epfmt:       SfmtEP(fmtEp),
					appfmt:      SfmtApp(fmtApp),
=======
					endpoints: sdEndpoints,
<<<<<<< HEAD
					epfmt: SfmtEP(spEp.fmtEp),
					appfmt: SfmtApp(spEp.fmtApp),
>>>>>>> Sysl generate sequence diagram
					activations: no_activations,
					title:       "",
					blackboxes:  append(bbs, bbs2...),
=======
					epfmt: SfmtEP(spep.fmtEp),
					appfmt: SfmtApp(spapp.fmtApp),
					activations: no_activations,
					title: spseqtitle.fmtSeq(app.GetEndpoints()[k].GetName(), app.GetEndpoints()[k].GetLongName(), varrefs),
					blackboxes: append(bbs, bbs2...),
>>>>>>> Sysl generate sequence diagram
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
		out, _ := generateSequenceDiag(mod, sd)

		logrus.Warnf("sd: %v", sd)
		OutputPlantuml(output, plantuml, out)
	}
}
