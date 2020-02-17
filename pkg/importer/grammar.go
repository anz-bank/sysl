package importer

/* FIXME
type grammarParam struct {
	grammarInput string
	grammarName  string
}

type Grammar struct {
	grammarParam
	grammar *parser.Grammar
	logger  *logrus.Logger
	types   TypeList
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
	}
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
	allTypes := syslutil.MakeStrSet()
	for _, rule := range g.grammar.GetRules() {
		rp := makeRule(rule)
		allTypes = allTypes.Union(rp.deps)
		for _, t := range rp.types.Items() {
			g.types.Add(t)
			allTypes.Remove(t.Name())
		}
	}

	for k := range allTypes {
		if _, found := g.types.Find(k); !found {
			g.types.Add(&Alias{name: k, Target: StringAlias})
		}
	}

	g.linkTypes()
	sort.SliceStable(g.types.Items(), func(i, j int) bool {
		// We want to sort so __prefix types end up under the rule they are linked with
		a := g.types.Items()[i].Name()
		b := g.types.Items()[j].Name()
		switch {
		case strings.HasPrefix(a, "__") && strings.HasPrefix(b, "__"):
			return strings.Compare(a, b) < 0
		case strings.HasPrefix(a, "__"):
			return strings.Compare(a[2:], b) < 0
		case strings.HasPrefix(b, "__"):
			return strings.Compare(a, b[2:]) < 0
		}
		return strings.Compare(a, b) < 0
	})

	if err := syslWriter.Write(info, g.types); err != nil {
		g.logger.Errorf("writing into buffer failed %s", err)
		return "", err
	}
	return syslBytes.String(), nil
}

func makeRule(rule *parser.Rule) ruleProcessor {
	rp := ruleProcessor{
		rule: rule,
		deps: syslutil.MakeStrSet(),
	}

	root := rp.createChoice(rule.Choices)
	if val, ok := rp.types.Find(root); ok {
		switch x := val.(type) {
		case *Union:
			x.name = rule.Name.Name
			rp.deps.Remove(root)
			return rp
		case *StandardType:
			x.name = rule.Name.Name
			rp.deps.Remove(root)
			return rp
		}
	}
	rp.deps.Insert(root)

	t := &Alias{name: rule.Name.Name, Target: &SyslBuiltIn{root}}
	rp.types.Add(t)

	return rp
}

type ruleProcessor struct {
	rule    *parser.Rule
	deps    syslutil.StrSet // Names of external rules which this rule depends on
	types   TypeList
	counter int
}

func (rp *ruleProcessor) newName() string {
	out := fmt.Sprintf("__%s_%02d", rp.rule.Name.Name, rp.counter)
	rp.counter++
	return out
}

func (rp *ruleProcessor) createChoice(c *parser.Choice) string {
	if len(c.Sequence) == 1 {
		return rp.buildSequence(c.Sequence[0])
	}
	t := &Union{
		name:    rp.newName(),
		Options: FieldList{},
	}
	for _, f := range c.Sequence {
		t.Options = append(t.Options, Field{Name: rp.buildSequence(f)})
	}
	rp.deps.Insert(t.name)
	rp.types.Add(t)
	return t.name
}

func makeSizedField(name string, q *parser.Quantifier) *Field {
	field := Field{
		Name: name,
	}
	nested := &Array{
		name: name,
	}
	if q == nil {
		return nil
	}
	switch q.Union.(type) {
	case *parser.Quantifier_Optional:
		field.Optional = true
	case *parser.Quantifier_ZeroPlus:
		field.Type = nested
		field.SizeSpec = &sizeSpec{
			Min:     0,
			Max:     0,
			MaxType: OpenEnded,
		}
	case *parser.Quantifier_OnePlus:
		field.Type = nested
		field.SizeSpec = &sizeSpec{
			Min:     1,
			Max:     0,
			MaxType: OpenEnded,
		}
	}
	return &field
}

func (rp *ruleProcessor) buildSequence(s *parser.Sequence) string {
	var fields []term
	for _, t := range s.Term {
		fields = append(fields, rp.buildTerm(t))
	}
	if len(fields) == 1 && fields[0].field == nil {
		rp.deps.Insert(fields[0].name)
		return fields[0].name
	}
	t := &StandardType{
		name:       rp.newName(),
		Properties: FieldList{},
	}
	for _, f := range fields {
		if len(fields) > 1 && f.name == StringTypeName {
			continue
		}
		if f.field == nil {
			f.field = &Field{Name: f.name}
		}
		t.Properties = append(t.Properties, *f.field)
		rp.deps.Insert(f.name)
	}
	rp.types.Add(t)
	rp.deps.Insert(t.name)
	return t.name
}

type term struct {
	name  string // name of the type being pointed to
	field *Field
}

func (rp *ruleProcessor) buildTerm(t *parser.Term) term {
	res := term{}
	switch a := t.Atom.Union.(type) {
	case *parser.Atom_String_:
		res.name = StringTypeName
	case *parser.Atom_Rulename:
		res.name = a.Rulename.Name
	case *parser.Atom_Choices:
		res.name = rp.createChoice(a.Choices)
	default:
		logrus.Fatal("Unhandled Term type")
	}
	if t.Quantifier != nil {
		res.field = makeSizedField(res.name, t.Quantifier)
	}
	return res
}

func linkType(t *Type, types TypeList, cache *map[string]Type) {
	find := func(name string) Type {
		if val, ok := (*cache)[name]; ok {
			return val
		} else if val, ok := types.Find(name); ok {
			(*cache)[name] = val
			return val
		}
		return StringAlias
	}
	switch x := (*t).(type) {
	case *Union:
		for i := range x.Options {
			if x.Options[i].Type == nil {
				x.Options[i].Type = find(x.Options[i].Name)
			} else {
				linkType(&x.Options[i].Type, types, cache)
			}
		}
	case *Alias:
		if x.Target == nil {
			x.Target = find(x.name)
		}
	case *StandardType:
		for i := range x.Properties {
			if x.Properties[i].Type == nil {
				x.Properties[i].Type = find(x.Properties[i].Name)
			} else {
				linkType(&x.Properties[i].Type, types, cache)
			}
		}
	case *Array:
		x.Items = find(x.name)
	}
}

func (g *Grammar) linkTypes() {
	cache := map[string]Type{}
	for i := range g.types.Items() {
		t := g.types.types[i]
		name := g.types.types[i].Name()
		cache[name] = t
		linkType(&t, g.types, &cache)
	}
}
*/
