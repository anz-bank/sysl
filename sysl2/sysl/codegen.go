package main

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/anz-bank/sysl/src/proto"
	parser "github.com/anz-bank/sysl/sysl2/naive"
	ebnfGrammar "github.com/anz-bank/sysl/sysl2/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Node can be string or node
type Node []interface{}

type codeGenOutput struct {
	filename string
	output   Node
}

func getKeyFromValueMap(v *sysl.Value, key string) *sysl.Value {
	if m := v.GetMap(); m != nil {
		return m.Items[key]
	}
	return nil
}

func processChoice(g *ebnfGrammar.Grammar, obj *sysl.Value, choice *ebnfGrammar.Choice) Node {
	var result Node

	for i, seq := range choice.Sequence {
		seqResult := Node{}
		fullScan := true
		for _, term := range seq.Term {
			switch x := term.Atom.Union.(type) {

			// String tokens dont have quantifiers
			case *ebnfGrammar.Atom_String_:
				seqResult = append(seqResult, x.String_)
			case *ebnfGrammar.Atom_Rulename:
				var ruleResult interface{}

				minc, maxc := parser.GetTermMinMaxCount(term)
				v := getKeyFromValueMap(obj, x.Rulename.Name)

				// raise error if required
				//  i.e.  no quantifier or +
				//        and missing from obj map
				if minc > 0 && v == nil {
					fullScan = false
					break
				}

				// skip if rule has
				//    quantifier == * or ?
				//    and does not exist in obj map
				if minc == 0 && v == nil {
					continue
				}

				if maxc > 1 {
					var valueList []*sysl.Value
					switch vv := v.Value.(type) {
					case *sysl.Value_List_:
						valueList = vv.List.Value
					case *sysl.Value_Set:
						valueList = vv.Set.Value
					default:
						logrus.Warnf("Expecting a collection type, got %T for rule %s", vv, x.Rulename.Name)
						fullScan = false
						break
					}
					ruleInstances := Node{}

					for _, valueItem := range valueList {
						// Drill down the rule
						node := processRule(g, valueItem, x.Rulename.Name)
						// Check post-conditions
						if len(node) == 0 {
							fullScan = false
							break
						}
						ruleInstances = append(ruleInstances, node)
					}
					ruleResult = ruleInstances
				} else { //maxc == 1
					// Drill down the rule
					if v.GetList() != nil || v.GetSet() != nil {
						logrus.Warnf("Got List or Set instead of map")
					}
					node := processRule(g, v, x.Rulename.Name)
					// Check post-conditions
					if len(node) == 0 {
						logrus.Warnf("could not process rule: ( %s )", x.Rulename.Name)
						fullScan = false
						break
					}
					if s, ok := node[0].(string); ok && len(node) == 1 {
						ruleResult = s
					} else {
						ruleResult = node
					}
				}
				seqResult = append(seqResult, ruleResult)
			case *ebnfGrammar.Atom_Choices:
				// minc, maxc := parser.GetMinMaxCount(term)
				node := processChoice(g, obj, x.Choices)
				if len(node) == 0 {
					logrus.Warnf("could not process Choice\n")
					fullScan = false
					break
				}
				seqResult = append(seqResult, node)
			default:
				logrus.Warningf("processChoice: choice %d : %T", i, x)
				panic("Unexpected atom type")
			}
			if !fullScan {
				break
			}
		}
		if fullScan {
			result = append(result, seqResult)
		}
	}
	return result
}

func processRule(g *ebnfGrammar.Grammar, obj *sysl.Value, ruleName string) Node {
	var str string
	if x := obj.GetMap(); x != nil {
		for key := range x.Items {
			str += key + ", "
		}
	}
	// logrus.Debugf("processRule: %s, obj keys (%s)", ruleName, str)
	rule := g.Rules[ruleName]
	if rule == nil {
		root := Node{}
		if IsCollectionType(obj) {
			return nil
		}
		// Should we convert int and bools to string and return?
		return append(root, obj.GetS())
	}
	root := processChoice(g, obj, rule.Choices)
	return root
}

func readGrammar(filename, grammarName, startRule string) *ebnfGrammar.Grammar {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Errorf("Unable to open grammar file: %s\nGot Error: %s", filename, err.Error())
		return nil
	}
	return parser.ParseEBNF(string(dat), grammarName, startRule)
}

// applyTranformToModel loads applies the transform to input model
func applyTranformToModel(modelName, transformAppName, viewName string, model, transform *sysl.Module) *sysl.Value {
	modelApp := model.Apps[modelName]
	view := transform.Apps[transformAppName].Views[viewName]
	if view == nil {
		panic(errors.Errorf("Cannot execute missing view: %s, in app %s", viewName, transformAppName))
	}
	s := Scope{}
	s.AddApp("app", modelApp)
	var result *sysl.Value
	// assume args are
	//  app <: sysl.App and
	//  type <: sysl.Type
	if len(view.Param) >= 2 {
		result = MakeValueList()
		var tNames []string
		for tName := range modelApp.Types {
			tNames = append(tNames, tName)
		}
		sort.Strings(tNames)
		for _, tName := range tNames {
			t := modelApp.Types[tName]
			s["typeName"] = MakeValueString(tName)
			s["type"] = typeToValue(t)
			appendItemToValueList(result.GetList(), EvalView(transform, transformAppName, viewName, s))
		}
	} else {
		result = EvalView(transform, transformAppName, viewName, s)
	}

	return result
}

// Serialize serializes node to string
func Serialize(w io.Writer, delim string, node Node) {
	for _, n := range node {
		switch x := n.(type) {
		case string:
			io.WriteString(w, x+delim)
		case Node:
			Serialize(w, delim, x)
		}
	}
}

// return the one and only app defined in the module
func getDefaultAppName(mod *sysl.Module) string {
	for app := range mod.Apps {
		return app
	}
	return ""
}

func loadAndGetDefaultApp(root, model string) (*sysl.Module, string) {
	// Model we want to generate code for
	mod, err := Parse(model, root)
	if err == nil {
		modelAppName := getDefaultAppName(mod)
		return mod, modelAppName
	}
	logrus.Errorf("unable to load module:\n\troot: " + root + "\n\tmodel:" + model)
	return nil, ""
}

// GenerateCode transform input sysl model to code in the target language described by
// grammar and a sysl transform
func GenerateCode(root_model, model, root_transform, transform, grammar, start string) []*codeGenOutput {
	var codeOutput []*codeGenOutput

	mod, modelAppName := loadAndGetDefaultApp(root_model, model)
	tx, transformAppName := loadAndGetDefaultApp(root_transform, transform)
	g := readGrammar(grammar, "gen", start)
	if g == nil {
		panic(errors.Errorf("Unable to load grammar"))
	}
	fileNames := applyTranformToModel(modelAppName, transformAppName, "filename", mod, tx)
	if fileNames == nil {
		return nil
	}
	result := applyTranformToModel(modelAppName, transformAppName, g.Start, mod, tx)
	if result.GetMap() != nil {
		filename := fileNames.GetMap().Items["filename"].GetS()
		logrus.Println(filename)
		r := processRule(g, result, g.Start)
		codeOutput = append(codeOutput, &codeGenOutput{filename, r})
	} else {
		fileValues := fileNames.GetList().Value
		for i, v := range result.GetList().Value {
			filename := fileValues[i].GetMap().Items["filename"].GetS()
			r := processRule(g, v, g.Start)
			codeOutput = append(codeOutput, &codeGenOutput{filename, r})
		}
	}
	return codeOutput
}

func outputToFiles(outDir string, output []*codeGenOutput) {
	for _, o := range output {
		f, err := os.Create(outDir + "/" + o.filename)
		if f == nil {
			logrus.Errorf("Unable to open file: %s\nGot error:\n%s", f.Name(), err.Error())
			continue
		}
		logrus.Warningln("Writing file: " + f.Name())
		Serialize(f, " ", o.output)
		f.Close()
	}
}

// DoGenerateCode generate code for the given model, using transform
// and the grammar of the target language
func DoGenerateCode(stdout, stderr io.Writer, flags *flag.FlagSet, args []string) int {
	root_model := flags.String("root-model", ".", "sysl root directory for input model file (default: .)")
	root_transform := flags.String("root-transform", ".", "sysl root directory for input transform file (default: .)")
	model := flags.String("model", ".", "model.sysl")
	transform := flags.String("transform", ".", "transform.sysl")
	grammar := flags.String("grammar", ".", "grammar.g")
	start := flags.String("start", ".", "start rule for the grammar")
	outDir := flags.String("outdir", ".", "output directory")

	flags.Parse(args[1:])
	logrus.Warnf("root_model: %s\n", *root_model)
	logrus.Warnf("root_transform: %s\n", *root_transform)
	logrus.Warnf("model: %s\n", *model)
	logrus.Warnf("transform: %s\n", *transform)
	logrus.Warnf("grammar: %s\n", *grammar)
	logrus.Warnf("start: %s\n", *start)
	output := GenerateCode(*root_model, *model, *root_transform, *transform, *grammar, *start)
	outputToFiles(*outDir, output)
	return 0
}
