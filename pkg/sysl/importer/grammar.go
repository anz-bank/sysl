package importer

import (
	"bytes"

	parser "github.com/anz-bank/sysl/pkg/naive"
	ebnfGrammar "github.com/anz-bank/sysl/pkg/proto"
	"github.com/sirupsen/logrus"
)

type grammarParam struct {
	grammarInput string
	grammarName  string
}

type Grammar struct {
	grammarParam
	grammar         *ebnfGrammar.Grammar
	logger          *logrus.Logger
	types           TypeList
	incompleteTypes map[string]struct{}
}

func LoadGrammar(args OutputData, grammarInput string, logger *logrus.Logger) (out string, err error) {
	gParam := grammarParam{
		grammarInput: grammarInput,
		grammarName:  args.AppName,
	}
	return newGrammar(gParam, logger).writeSysl()
}

func newGrammar(param grammarParam, logger *logrus.Logger) *Grammar {
	return &Grammar{
		grammarParam:    param,
		logger:          logger,
		types:           TypeList{},
		incompleteTypes: map[string]struct{}{},
	}
}

func (g *Grammar) addTerm(name string, term *ebnfGrammar.Term, props *FieldList) {
	switch x := term.Atom.Union.(type) {
	case *ebnfGrammar.Atom_String_:
		// hit the leaf node
		// return
	case *ebnfGrammar.Atom_Choices:
		cname := "__Choice_" + name
		if _, found := g.types.Find(cname); !found {
			g.types.Add(g.createChoice(x.Choices, cname))
		}
		choiceField := g.quantifyUnionType(cname, cname, term.Quantifier)
		*props = append(*props, choiceField)

	case *ebnfGrammar.Atom_Rulename:
		var typeString string
		if t, found := g.types.Find(x.Rulename.GetName()); found {
			typeString = t.Name()
		} else {
			typeString = x.Rulename.Name
			g.incompleteTypes[typeString] = struct{}{}
		}
		ruleField := g.quantifyUnionType(x.Rulename.Name, typeString, term.Quantifier)
		*props = append(*props, ruleField)
	}
}

func (g *Grammar) createChoice(choice *ebnfGrammar.Choice, name string) Type {
	choiceType := Union{
		name: name,
	}
	for _, seq := range choice.Sequence {
		for _, term := range seq.Term {
			g.addTerm(name, term, &choiceType.Options)
		}
	}
	return &choiceType
}

func (g *Grammar) quantifyUnionType(fieldName, fieldType string, q *ebnfGrammar.Quantifier) Field {
	field := Field{
		Name: fieldName,
		Type: &StandardType{
			name: fieldType,
		}}

	if q != nil {
		switch q.Union.(type) {
		case *ebnfGrammar.Quantifier_Optional:
			field.Optional = true
		case *ebnfGrammar.Quantifier_ZeroPlus:
			field.Type = &Array{
				name:  fieldName,
				Items: field.Type,
			}
			field.SizeSpec = &sizeSpec{
				Min:     0,
				Max:     0,
				MaxType: OpenEnded,
			}
		case *ebnfGrammar.Quantifier_OnePlus:
			field.SizeSpec = &sizeSpec{
				Min:     1,
				Max:     0,
				MaxType: OpenEnded,
			}
		}
	}
	return field
}

func (g *Grammar) createRule(rule *ebnfGrammar.Rule, name string) Type {
	ruleType := StandardType{
		name: name,
	}
	for _, seq := range rule.Choices.Sequence {
		for _, term := range seq.Term {
			g.addTerm(name, term, &ruleType.Properties)
		}
	}

	if len(ruleType.Properties) == 0 {
		t, found := g.types.Find(name)
		if !found {
			t = &Alias{
				name:   name,
				Target: &SyslBuiltIn{name: StringTypeName},
			}
		}
		return t
	}
	return &ruleType
}

func (g *Grammar) writeSysl() (string, error) {
	g.grammar = parser.ParseEBNF(g.grammarParam.grammarInput, "", "")
	var syslBytes bytes.Buffer
	syslWriter := newWriter(&syslBytes, g.logger)
	syslWriter.DisableJSONTags = true
	info := SyslInfo{
		OutputData: OutputData{
			AppName: g.grammarName,
		},
		Title:       g.grammarName + "File",
		Description: g.grammarName + " Grammar",
	}
	for ruleName, rule := range g.grammar.GetRules() {
		if _, found := g.types.Find(ruleName); !found {
			ruleType := g.createRule(rule, ruleName)
			g.types.Add(ruleType)
		}
	}
	for name := range g.incompleteTypes {
		if _, found := g.types.Find(name); !found {
			g.types.Add(&Alias{
				name:   name,
				Target: &SyslBuiltIn{name: StringTypeName},
			})
		}
	}

	g.types.Sort()
	if err := syslWriter.Write(info, g.types, ""); err != nil {
		g.logger.Errorf("writing into buffer failed %s", err)
		return "", err
	}
	return syslBytes.String(), nil
}
