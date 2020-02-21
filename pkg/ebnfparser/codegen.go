package ebnfparser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/hcl/strconv"

	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/eval"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

type EbnfGrammar struct {
	Rules map[string][]TermNode
	Start string
}

func ReadGrammar(fs afero.Fs, filename string, start string) (*EbnfGrammar, error) {
	grammar, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, err
	}
	g, err := ParseString(string(grammar))
	if err != nil {
		return nil, err
	}

	res := EbnfGrammar{
		Rules: map[string][]TermNode{},
		Start: start,
	}
	WalkerOps{EnterProdNode: func(node ProdNode) Stopper {
		res.Rules[node.OneIdent().String()] = node.AllTerm()
		return nil
	}}.Walk(g)

	_, has := res.Rules[start]
	if !has {
		return nil, fmt.Errorf("start rule '%s' not defined", start)
	}
	return &res, nil
}

func GenerateOutput(grammar *EbnfGrammar, value *sysl.Value, logger *logrus.Logger) (string, error) {
	return walk(grammar.Rules[grammar.Start], value, grammar.Rules, logger)
}

type walker struct {
	WalkerOps

	logger *logrus.Logger
	rules  map[string][]TermNode
	out    []string

	err   error
	value *sysl.Value
}

func walk(terms []TermNode, val *sysl.Value, rules map[string][]TermNode, logger *logrus.Logger) (string, error) {
	w := walker{
		logger: logger,
		value:  val,
		rules:  rules,
	}
	w.WalkerOps = WalkerOps{
		EnterTermNode: w.EnterTermNode,
		EnterStringNode: func(node StringNode) Stopper {
			s, err := strconv.Unquote(node.String())
			if err != nil {
				s = node.String()
			}
			s = strings.NewReplacer("\\n", "\n", "\\t", "\t").Replace(s)
			w.addToken(s)
			return nil
		},
	}

	for _, t := range terms {
		WalkTermNode(t, w.WalkerOps)
	}

	return strings.TrimRight(strings.Join(w.out, " "), " "), w.err
}

func (w walker) findRule(name string) []TermNode {
	out, has := w.rules[name]
	if has {
		return out
	}
	return nil
}

func (w *walker) addToken(s string) {
	w.out = append(w.out, s)
}

func (w *walker) EnterTermNode(node TermNode) Stopper {
	if node.OneOp() != "" {
		// handle choices
		for _, option := range node.AllTerm() {
			out, err := walk([]TermNode{option}, w.value, w.rules, w.logger)
			if err == nil {
				w.addToken(out)
				return &nodeExiter{}
			}
		}
		w.err = fmt.Errorf("none of the options were satisfied")
		return &aborter{}
	}
	min := 1
	max := 1
	switch node.OneQuant() {
	case "*":
		min = 0
		max = 0
	case "?":
		min = 0
	case "+":
		max = 0
	}
	_ = max

	if atom := node.OneAtom(); atom.Node != nil {
		if id := atom.OneRule().String(); id != "" {
			obj := getValue(id, w.value)
			if obj == nil {
				if min == 1 {
					w.err = fmt.Errorf("expected a '%s' but none found", id)
					return &aborter{}
				}
				return nil
			}
			if !eval.IsCollectionType(obj) {
				obj = eval.MakeValueList(obj)
			}
			vals := eval.GetValueSlice(obj)
			if max == 1 && len(vals) > 1 {
				w.err = fmt.Errorf("expected a single '%s' but %d found", id, len(vals))
				return &aborter{}
			} else if min == 1 && len(vals) == 0 {
				w.err = fmt.Errorf("expected a '%s' but none found", id)
				return &aborter{}
			}
			var fn func(string, *sysl.Value) error
			if next := w.findRule(id); next != nil {
				fn = w.applyRule
			} else {
				fn = func(s string, value *sysl.Value) error {
					w.addToken(getValue(s, value).GetS())
					return nil
				}
			}

			for _, v := range vals {
				err := fn(id, v)
				w.err = err
				if err != nil {
					return &aborter{}
				}
			}
		}
	}
	return nil
}

func (w *walker) applyRule(name string, value *sysl.Value) error {
	next := w.findRule(name)
	s, err := walk(next, value, w.rules, w.logger)
	if err != nil {
		return err
	}
	w.addToken(s)
	return nil
}

func getValue(name string, value *sysl.Value) *sysl.Value {
	switch v := value.Value.(type) {
	case *sysl.Value_Map_:
		if v, has := v.Map.Items[name]; has {
			return v
		}
	case *sysl.Value_S:
		return value
	default:
		// Should we convert int and bools to string and return?
		//return eval.UnaryString(value)
		return nil
	}
	return nil
}
