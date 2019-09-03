package validate

import (
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/anz-bank/sysl/sysl2/sysl/msg"
	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Validator struct {
	grammar     *sysl.Application
	transform   *sysl.Application
	assignTypes map[string]parse.TypeData
	letTypes    map[string]parse.TypeData
	messages    map[string][]msg.Msg
}

type CmdContextParamValidate struct {
	rootTransform *string
	transform     *string
	grammar       *string
	start         *string
	isVerbose     *bool
	loglevel      *string
}

func DoValidate(validateParams *CmdContextParamValidate) error {
	logrus.Debugf("root-transform: %s\n", *validateParams.rootTransform)
	logrus.Debugf("transform: %s\n", *validateParams.transform)
	logrus.Debugf("grammar: %s\n", *validateParams.grammar)
	logrus.Debugf("start: %s\n", *validateParams.start)
	logrus.Debugf("loglevel: %s\n", *validateParams.loglevel)
	logrus.Debugf("isVerbose: %v\n", *validateParams.isVerbose)

	// Default info
	if level, has := syslutil.LogLevels[*validateParams.loglevel]; has {
		logrus.SetLevel(level)
	}
	if *validateParams.isVerbose {
		*validateParams.loglevel = "debug"
	}

	fs := afero.NewOsFs()

	grammar, err := LoadGrammar(*validateParams.grammar, fs)
	if err != nil {
		return err
	}

	parser := parse.NewParser()
	transform, err := loadTransform(
		*validateParams.transform,
		syslutil.NewChrootFs(fs, *validateParams.rootTransform),
		parser,
	)
	if err != nil {
		return err
	}

	validator := NewValidator(grammar, transform, parser)
	validator.Validate(*validateParams.start)
	validator.LogMessages()

	if len(validator.GetMessages()) > 0 {
		msg.NewMsg(msg.ErrValidationFailed, nil).LogMsg()
		return errors.New("validation failed")
	}

	msg.NewMsg(msg.InfoValidatedSuccessfully, nil).LogMsg()

	return nil
}

func (v *Validator) Validate(start string) {
	v.validateEntryPoint(start)
	v.validateFileName()
	v.validateViews()
}

func (v *Validator) LogMessages() {
	for viewName, messages := range v.GetMessages() {
		msg.NewMsg(msg.TitleViewName, []string{viewName}).LogMsg()
		for _, message := range messages {
			message.LogMsg()
		}
	}
}

func (v *Validator) validateEntryPoint(start string) {
	view, exists := v.transform.Views[start]

	if !exists {
		v.messages["EntryPoint"] = append(v.messages["EntryPoint"], *msg.NewMsg(msg.ErrEntryPointUndefined, []string{start}))
		return
	}

	if getTypeName(view.GetRetType()) != start || isCollectionType(view.GetRetType()) {
		v.messages["EntryPoint"] = append(v.messages["EntryPoint"],
			*msg.NewMsg(msg.ErrInvalidEntryPointReturn, []string{start, start}))
	}
}

func (v *Validator) validateFileName() {
	viewName := "filename"
	view, exists := v.transform.Views[viewName]

	if !exists {
		v.messages[viewName] = append(v.messages[viewName], *msg.NewMsg(msg.ErrUndefinedView, []string{viewName}))
		return
	}

	if getTypeName(view.GetRetType()) != "STRING" || isCollectionType(view.GetRetType()) {
		v.messages[viewName] = append(v.messages[viewName], *msg.NewMsg(msg.ErrInvalidReturn, []string{viewName, "string"}))
	}

	hasFileNameAssign := false

	for _, stmt := range view.GetExpr().GetTransform().GetStmt() {
		if stmt.GetAssign() != nil {
			hasFileNameAssign = hasFileNameAssign || stmt.GetAssign().GetName() == viewName
			if stmt.GetAssign().GetName() != viewName {
				v.messages[viewName] = append(v.messages[viewName],
					*msg.NewMsg(msg.ErrExcessAttr, []string{stmt.GetAssign().GetName(), viewName, "string"}))
			}
		}
	}

	if !hasFileNameAssign {
		v.messages[viewName] = append(v.messages[viewName],
			*msg.NewMsg(msg.ErrMissingReqField, []string{viewName, viewName, "string"}))
	}
}

func (v *Validator) validateViews() {
	for viewName, resolvedTypes := range v.assignTypes {
		typeName := getTypeName(resolvedTypes.RefType)
		resolvedType := resolvedTypes.Tuple
		if grammarType, exists := v.grammar.Types[typeName]; exists {
			switch t := grammarType.Type.(type) {
			case *sysl.Type_Tuple_:
				v.compareTuple(t.Tuple, resolvedType.GetTuple(),
					getAttrNames(resolvedType.GetTuple().GetAttrDefs()), viewName, typeName)
			default:
				logrus.Warnf("[validate.validateViews] Unhandled grammar type: %T", t)
			}

		}
	}

	for viewName, resolvedTypes := range v.letTypes {
		typeName := getTypeName(resolvedTypes.RefType)
		resolvedType := resolvedTypes.Tuple
		if grammarType, exists := v.grammar.Types[typeName]; exists {
			switch t := grammarType.Type.(type) {
			case *sysl.Type_Tuple_:
				v.compareTuple(t.Tuple, resolvedType.GetTuple(),
					getAttrNames(resolvedType.GetTuple().GetAttrDefs()), viewName, typeName)
			default:
				logrus.Warnf("[validate.validateViews] Unhandled grammar type: %T", t)
			}

		}
	}
}

func (v *Validator) compareTuple(
	specTuple, implTuple *sysl.Type_Tuple,
	implAttrNames map[string]struct{},
	viewName, specTupleName string,
) {
	grammarSpec := v.grammar.Types

	specAttrs := specTuple.GetAttrDefs()
	implAttrs := implTuple.GetAttrDefs()

	for ikey, ival := range implTuple.GetAttrDefs() {
		if ival.GetTuple() == nil {
			continue
		}

		if grammarType, exists := grammarSpec[ikey]; exists {
			v.compareTuple(
				grammarType.GetTuple(), ival.GetTuple(), getAttrNames(ival.GetTuple().GetAttrDefs()), viewName, ikey)
		}
	}

	for gk, gv := range specAttrs {
		if specOneOf := grammarSpec[gk].GetOneOf(); specOneOf != nil {
			v.compareOneOf(specOneOf, implTuple, implAttrNames, viewName, specTupleName)
		} else if _, exists := implAttrs[gk]; !exists {
			if !gv.GetOpt() {
				v.messages[viewName] = append(v.messages[viewName],
					*msg.NewMsg(msg.ErrMissingReqField, []string{gk, viewName, specTupleName}))
			}
		} else {
			delete(implAttrNames, gk)
		}
	}

	for attrName := range implAttrNames {
		v.messages[viewName] = append(v.messages[viewName],
			*msg.NewMsg(msg.ErrExcessAttr, []string{attrName, viewName, specTupleName}))
		delete(implAttrNames, attrName)
	}
}

func (v *Validator) compareOneOf(
	specOneOf *sysl.Type_OneOf,
	implTuple *sysl.Type_Tuple,
	implAttrNames map[string]struct{},
	viewName, specTupleName string,
) {
	implAttrs := implTuple.GetAttrDefs()
	matching := true
	grammarSpec := v.grammar.Types

	for _, one := range specOneOf.GetType() {
		name := one.GetTypeRef().GetRef().GetPath()[0]

		if strings.Index(name, "__Choice_Combination_") == 0 {
			if len(implAttrs) == 1 {
				continue
			}
			v.compareTuple(grammarSpec[name].GetTuple(), implTuple, implAttrNames, viewName, specTupleName)
			break
		} else if _, matching = implAttrs[name]; matching {
			delete(implAttrNames, name)
			break
		}
	}

	if !matching {
		var implAttrNames []string
		for k := range implAttrs {
			implAttrNames = append(implAttrNames, k)
		}
		v.messages[viewName] = append(v.messages[viewName],
			*msg.NewMsg(msg.ErrInvalidOption, []string{viewName, strings.Join(implAttrNames, ","), specTupleName}))
	}
}

func (v *Validator) GetMessages() map[string][]msg.Msg {
	return v.messages
}

func NewValidator(grammar, transform *sysl.Application, parser *parse.Parser) *Validator {
	return &Validator{
		grammar:     grammar,
		transform:   transform,
		assignTypes: parser.GetAssigns(),
		letTypes:    parser.GetLets(),
		messages:    parser.GetMessages()}
}

func getTypeName(syslType *sysl.Type) string {
	if syslType == nil {
		return "Unknown"
	}

	switch t := syslType.Type.(type) {
	case *sysl.Type_Primitive_:
		return t.Primitive.String()
	case *sysl.Type_Sequence:
		if typeName := t.Sequence.GetPrimitive().String(); typeName != "NO_Primitive" {
			return typeName
		}
		return t.Sequence.GetTypeRef().GetRef().GetPath()[0]
	case *sysl.Type_TypeRef:
		if t.TypeRef.GetRef().GetAppname() != nil {
			return t.TypeRef.GetRef().GetAppname().GetPart()[0]
		}
		return t.TypeRef.GetRef().GetPath()[0]
	default:
		return "Unknown"
	}
}

func isCollectionType(syslType *sysl.Type) bool {
	switch syslType.Type.(type) {
	case *sysl.Type_Set, *sysl.Type_Sequence, *sysl.Type_List_, *sysl.Type_Map_:
		return true
	default:
		return false
	}
}

func getAttrNames(attrs map[string]*sysl.Type) map[string]struct{} {
	implAttrNames := map[string]struct{}{}

	for attrName := range attrs {
		implAttrNames[attrName] = struct{}{}
	}

	return implAttrNames
}

func loadTransform(transformFile string, fs afero.Fs, p *parse.Parser) (*sysl.Application, error) {
	transform, name, err := parse.LoadAndGetDefaultApp(transformFile, fs, p)
	if err != nil {
		return nil, err
	}
	return transform.GetApps()[name], nil
}

// LoadGrammar loads sysl conversion of given grammar.
// eg: if grammarFile is ./foo/bar.g, this tries to load ./foo/bar.sysl
func LoadGrammar(grammarFile string, fs afero.Fs) (*sysl.Application, error) {
	tokens := strings.Split(grammarFile, "/")
	rootGrammar := strings.Join(tokens[:len(tokens)-1], "/")
	grammarFileName := tokens[len(tokens)-1]

	tokens = strings.Split(grammarFileName, ".")
	tokens[len(tokens)-1] = "sysl"
	grammarSyslFile := strings.Join(tokens, ".")
	p := parse.NewParser()

	grammar, name, err := parse.LoadAndGetDefaultApp(grammarSyslFile, syslutil.NewChrootFs(fs, rootGrammar), p)
	if err != nil {
		return nil, err
	}
	return grammar.GetApps()[name], nil
}

// ConfigureCmdlineForValidate configures commandline params related to validate
func ConfigureCmdlineForValidate(sysl *kingpin.Application) *CmdContextParamValidate {
	validate := sysl.Command("validate", "Validate transform")
	return &CmdContextParamValidate{
		rootTransform: validate.Flag("root-transform",
			"sysl root directory for input transform file (default: .)").Default(".").String(),
		transform: validate.Flag("transform", "grammar.g").Default(".").String(),
		grammar:   validate.Flag("grammar", "grammar.sysl").Default(".").String(),
		start:     validate.Flag("start", "start rule for the grammar").Default(".").String(),
		loglevel:  validate.Flag("log", "log level[debug,info,warn,off]").Default("info").String(),
		isVerbose: validate.Flag("verbose", "show output").Short('v').Default("false").Bool(),
	}
}
