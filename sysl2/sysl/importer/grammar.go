package importer

import (
	"bytes"
	"fmt"

	parser "github.com/anz-bank/sysl/sysl2/naive"
	ebnfGrammar "github.com/anz-bank/sysl/sysl2/proto"
	"github.com/sirupsen/logrus"
)

type grammarParam struct {
	grammarInput string
	grammarName  string
}

type Grammar struct {
	grammarParam
	grammar    *ebnfGrammar.Grammar
	logger     *logrus.Logger
	types      TypeList
	visitedMap map[string]bool
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
		grammarParam: param,
		logger:       logger,
		types:        TypeList{},
		visitedMap:   map[string]bool{},
	}
}

func (g *Grammar) createChoice(choice *ebnfGrammar.Choice, name string) Type {
	choiceType := Union{
		name: name,
	}
	// mark the node as visited if not visited before
	if !g.visitedMap[name] {
		g.visitedMap[name] = true
	}
	for _, seq := range choice.Sequence {
		for _, term := range seq.Term {
			switch x := term.Atom.Union.(type) {
			case *ebnfGrammar.Atom_String_:
				// hit the leaf node
				// return
			case *ebnfGrammar.Atom_Choices:
				choiceType.Attributes = append(choiceType.Attributes, "__Choice_"+name)

				// visit the unvisited node
				if !g.visitedMap["__Choice_"+name] {
					childChoiceType := g.createChoice(x.Choices, "__Choice_"+name)
					g.types.types = append([]Type{childChoiceType}, g.types.types...)
				}
			case *ebnfGrammar.Atom_Rulename:
				var childRuleType Type
				choiceType.Attributes = append(choiceType.Attributes, x.Rulename.GetName())

				// visit the unvisited node
				if !g.visitedMap[x.Rulename.GetName()] {
					if rule, exists := g.grammar.Rules[x.Rulename.GetName()]; exists {
						childRuleType = g.createRule(g.grammar.Rules[rule.Name.Name], rule.Name.Name)
					} else {
						childRuleType = g.createEnumType(x.Rulename.GetName())
					}
					g.types.types = append([]Type{childRuleType}, g.types.types...)
				}
			default:
				g.logger.Warnf("unexpected term type")
			}
		}
	}
	var baseType Type = &choiceType
	return baseType
}

func (g *Grammar) quantifyUnionType(fieldName, fieldType string, quant *ebnfGrammar.Quantifier) Field {
	if quant == nil {
		return Field{
			Name: fieldName,
			Type: &StandardType{
				name: fieldType,
			}}
	}
	switch quant.Union.(type) {
	case *ebnfGrammar.Quantifier_Optional:
		fieldType = fmt.Sprintf("%s?", fieldType)
	case *ebnfGrammar.Quantifier_ZeroPlus:
		fieldType = fmt.Sprintf("sequence of %s?", fieldType)
	case *ebnfGrammar.Quantifier_OnePlus:
		fieldName = fmt.Sprintf("%s(1..)", fieldName)
	default:
		g.logger.Warnf("unexpected quantifier type")
	}
	return Field{
		Name: fieldName,
		Type: &StandardType{
			name: fieldType,
		}}
}

func (g *Grammar) createRule(rule *ebnfGrammar.Rule, name string) Type {
	ruleType := StandardType{
		name: name,
	}
	if !g.visitedMap[name] {
		g.visitedMap[name] = true
	}
	for _, seq := range rule.Choices.Sequence {
		for _, term := range seq.Term {
			switch x := term.Atom.Union.(type) {
			case *ebnfGrammar.Atom_String_:
				// hit the leaf node
				// return
			case *ebnfGrammar.Atom_Choices:
				choiceField := g.quantifyUnionType("__Choice_"+name, "__Choice_"+name, term.Quantifier)
				ruleType.Properties = append(ruleType.Properties, choiceField)

				// visit the unvisited node
				if !g.visitedMap["__Choice_"+name] {
					childChoiceType := g.createChoice(x.Choices, "__Choice_"+name)
					g.types.types = append([]Type{childChoiceType}, g.types.types...)
				}
			case *ebnfGrammar.Atom_Rulename:
				var childRuleType Type
				typeString := "string"
				if _, exists := g.grammar.Rules[x.Rulename.GetName()]; exists {
					typeString = x.Rulename.Name
				}
				ruleField := g.quantifyUnionType(x.Rulename.Name, typeString, term.Quantifier)
				ruleType.Properties = append(ruleType.Properties, ruleField)

				// visit the unvisited node
				if !g.visitedMap[x.Rulename.GetName()] {
					if rule, exists := g.grammar.Rules[x.Rulename.GetName()]; exists {
						childRuleType = g.createRule(g.grammar.Rules[rule.Name.Name], rule.Name.Name)
						g.types.types = append([]Type{childRuleType}, g.types.types...)
					}
				}
			default:
				g.logger.Warnf("unexpected term type")
			}
		}
	}
	var baseType Type = &ruleType
	return baseType
}

func (g *Grammar) writeSysl() (string, error) {
	g.grammar = parser.ParseEBNF(g.grammarParam.grammarInput, "", "")
	var syslBytes bytes.Buffer
	syslWriter := newWriter(&syslBytes, g.logger)
	info := SyslInfo{
		OutputData: OutputData{
			AppName: g.grammarName,
		},
		Title:       g.grammarName + "File",
		Description: g.grammarName + " Grammar",
	}
	for keyRuleName, keyRule := range g.grammar.GetRules() {
		if _, visited := g.visitedMap[keyRuleName]; !visited {
			ruleType := g.createRule(keyRule, keyRuleName)
			g.types.types = append([]Type{ruleType}, g.types.types...)
			g.visitedMap[keyRuleName] = true
		}
	}
	g.types.Sort()
	if err := syslWriter.Write(info, g.types, ""); err != nil {
		g.logger.Errorf("writing into buffer failed %s", err)
		return "", err
	}
	return syslBytes.String(), nil
}

func (g *Grammar) createEnumType(typeName string) Type {
	ruleEnum := Enum{
		name: typeName,
	}
	var baseType Type = &ruleEnum
	return baseType
}
