package validate

import (
	"os"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/pkg/msg"
	"github.com/anz-bank/sysl/pkg/parse"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var depPathRegEx = regexp.MustCompile("^[[:alnum:]]+([.][[:alnum:]]+)*(/[[:alnum:]]+)*[/]?$")
var basePathRegEx = regexp.MustCompile("^/[[:alnum:]]+(/[[:alnum:]]+)*$")

type Validator struct {
	grammar     *sysl.Application
	transform   *sysl.Application
	assignTypes map[string]parse.TypeData
	letTypes    map[string]parse.TypeData
	messages    map[string][]msg.Msg
}

type Params struct {
	RootTransform string
	Transform     string
	Grammar       string
	DepPath       string
	Start         string
	BasePath      string
	ParserType    parse.ParserType
	Filesystem    afero.Fs
	Logger        *logrus.Logger
}

func DoValidate(validateParams Params) error {
	logrus.Debugf("root-transform: %s\n", validateParams.RootTransform)
	logrus.Debugf("transform: %s\n", validateParams.Transform)
	logrus.Debugf("dep-path: %s\n", validateParams.DepPath)
	logrus.Debugf("grammar: %s\n", validateParams.Grammar)
	logrus.Debugf("start: %s\n", validateParams.Start)
	logrus.Debugf("basepath: %s\n", validateParams.BasePath)
	logrus.Debugf("parser: %s\n", validateParams.ParserType)

	grammar, err := LoadGrammarWithParserType(validateParams.Grammar, validateParams.Filesystem,
		validateParams.ParserType)
	if err != nil {
		return err
	}

	parser := parse.NewParserWithParserType(validateParams.ParserType)
	transform, err := loadTransform(
		validateParams.Transform,
		syslutil.NewChrootFs(validateParams.Filesystem, validateParams.RootTransform),
		parser,
	)
	if err != nil {
		return err
	}

	validator := NewValidator(grammar, transform, parser)
	validator.Validate(validateParams.Start, validateParams.DepPath, validateParams.BasePath)
	validator.LogMessages()

	if len(validator.GetMessages()) > 0 {
		msg.NewMsg(msg.ErrValidationFailed, nil).LogMsg()
		return errors.New("validation failed")
	}

	msg.NewMsg(msg.InfoValidatedSuccessfully, nil).LogMsg()

	return nil
}

func (v *Validator) Validate(start, depPath, basepath string) {
	v.validateEntryPoint(start)
	v.validateDependencyPath(depPath)
	v.validateBasePath(basepath)
	v.validateFileName()
	v.validateViews()
	v.validateReturn()
}

func (v *Validator) LogMessages() {
	for viewName, messages := range v.GetMessages() {
		msg.NewMsg(msg.TitleViewName, []string{viewName}).LogMsg()
		for _, message := range messages {
			message.LogMsg()
		}
	}
}

func (v *Validator) validateDependencyPath(depPath string) {
	if depPath == "" {
		return
	}

	if match := depPathRegEx.MatchString(depPath); !match {
		v.messages["DepPath"] = append(v.messages["DepPath"],
			*msg.NewMsg(msg.ErrDepPathInvalid, []string{depPath, depPathRegEx.String()}))
	}
}

func (v *Validator) validateBasePath(basepath string) {
	if basepath == "" || basepath == "/" {
		return
	}

	if match := basePathRegEx.MatchString(basepath); !match {
		v.messages["BasePath"] = append(v.messages["BasePath"],
			*msg.NewMsg(msg.ErrBasePathInvalid, []string{basepath, basePathRegEx.String()}))
	}
}

func (v *Validator) validateEntryPoint(start string) {
	view, exists := v.transform.Views[start]

	if !exists {
		v.messages["EntryPoint"] = append(v.messages["EntryPoint"], *msg.NewMsg(msg.ErrEntryPointUndefined, []string{start}))
		return
	}

	typeName := getTypeName(view.GetRetType())
	if typeName != start || isCollectionType(view.GetRetType()) {
		v.messages["EntryPoint"] = append(v.messages["EntryPoint"],
			*msg.NewMsg(msg.ErrInvalidEntryPointReturn, []string{typeName, start}))
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
		v.messages[viewName] = append(v.messages[viewName],
			*msg.NewMsg(msg.ErrInvalidReturn, []string{viewName, "Expected type is string"}))
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

func (v *Validator) validateReturn() {
	for viewName, view := range v.transform.GetViews() {
		v.validateTfmReturn(viewName, view.GetExpr(), view.GetRetType())
	}
}

func (v *Validator) validateTfmReturn(viewName string, expr *sysl.Expr, retType *sysl.Type) {
	if tfm, ok := expr.Expr.(*sysl.Expr_Transform_); ok {
		for _, stmt := range tfm.Transform.GetStmt() {
			switch s := stmt.Stmt.(type) {
			case *sysl.Expr_Transform_Stmt_Assign_:
				v.validateTfmReturn(viewName, s.Assign.GetExpr(), s.Assign.GetExpr().Type)
			case *sysl.Expr_Transform_Stmt_Let:
				v.validateTfmReturn(viewName, s.Let.GetExpr(), s.Let.GetExpr().Type)
			}
		}

		switch retType.Type.(type) {
		case *sysl.Type_Sequence, *sysl.Type_Set:
			if tfm.Transform.Scopevar == "." {
				typeName, typeDetail := syslutil.GetTypeDetail(retType)
				v.messages[viewName] = append(v.messages[viewName],
					*msg.NewMsg(msg.ErrInvalidReturn, []string{viewName, "Expected a " + typeName + " of " + typeDetail}))
			}
		default:
			if tfm.Transform.Scopevar != "." {
				_, typeDetail := syslutil.GetTypeDetail(retType)
				v.messages[viewName] = append(v.messages[viewName],
					*msg.NewMsg(msg.ErrInvalidReturn, []string{viewName, "Expected a single " + typeDetail}))
			}
		}
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
	return LoadGrammarWithParserType(grammarFile, fs, parse.DefaultParserType)
}

func LoadGrammarWithParserType(grammarFile string, fs afero.Fs, parserType parse.ParserType) (*sysl.Application, error) {
	tokens := strings.Split(grammarFile, string(os.PathSeparator))
	rootGrammar := strings.Join(tokens[:len(tokens)-1], string(os.PathSeparator))
	grammarFileName := tokens[len(tokens)-1]

	tokens = strings.Split(grammarFileName, ".")
	tokens[len(tokens)-1] = "sysl"
	grammarSyslFile := strings.Join(tokens, ".")
	p := parse.NewParserWithParserType(parserType)

	grammar, name, err := parse.LoadAndGetDefaultApp(grammarSyslFile,
		syslutil.NewChrootFs(fs, rootGrammar), p)
	if err != nil {
		return nil, err
	}
	return grammar.GetApps()[name], nil
}
