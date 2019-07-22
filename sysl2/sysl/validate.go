package main

import (
	"flag"
	"fmt"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type TypeData struct {
	refType *sysl.Type
	tuple   *sysl.Type
}

type Validator struct {
	grammar     *sysl.Application
	transform   *sysl.Application
	assignTypes map[string]TypeData
	letTypes    map[string]TypeData
	messages    map[string][]Msg
}

func DoValidate(flags *flag.FlagSet, args []string) error {
	rootTransform := flags.String("root-transform", ".", "sysl root directory for input transform file (default: .)")
	transformFile := flags.String("transform", ".", "transform.sysl")
	grammarFile := flags.String("grammar", "", "grammar.g")
	start := flags.String("start", "", "start rule for the grammar")

	if err := flags.Parse(args[2:]); err != nil {
		return err
	}

	logrus.Infof("root-transform: %s\n", *rootTransform)
	logrus.Infof("transform: %s\n", *transformFile)
	logrus.Infof("grammar: %s\n", *grammarFile)
	logrus.Infof("start: %s\n", *start)

	grammar, err := loadGrammar(*grammarFile)
	if err != nil {
		return err
	}

	p := NewParser()
	transform, err := loadTransform(*rootTransform, *transformFile, p)
	if err != nil {
		return err
	}

	validator := NewValidator(grammar, transform, p)
	validator.Validate(*start)

	for viewName, messages := range validator.GetMessages() {
		NewMsg(TitleViewName, []string{viewName}).logMsg()
		for _, message := range messages {
			message.logMsg()
		}
	}

	if len(validator.GetMessages()) > 0 {
		NewMsg(ErrValidationFailed, nil).logMsg()
		return errors.New("validation failed")
	}

	NewMsg(InfoValidatedSuccessfully, nil).logMsg()

	return nil
}

func (v *Validator) Validate(start string) {
	v.validateEntryPoint(start)
	v.validateFileName()
	v.validateViews()
}

func (v *Validator) validateEntryPoint(start string) {
	view, exists := v.transform.Views[start]

	if !exists {
		v.messages["EntryPoint"] = append(v.messages["EntryPoint"], *NewMsg(ErrEntryPointUndefined, []string{start}))
		return
	}

	if getTypeName(view.GetRetType()) != start || isCollectionType(view.GetRetType()) {
		v.messages["EntryPoint"] = append(v.messages["EntryPoint"],
			*NewMsg(ErrInvalidEntryPointReturn, []string{start, start}))
	}
}

func (v *Validator) validateFileName() {
	viewName := "filename"
	view, exists := v.transform.Views[viewName]

	if !exists {
		v.messages[viewName] = append(v.messages[viewName], *NewMsg(ErrUndefinedView, []string{viewName}))
		return
	}

	if getTypeName(view.GetRetType()) != "STRING" || isCollectionType(view.GetRetType()) {
		v.messages[viewName] = append(v.messages[viewName], *NewMsg(ErrInvalidReturn, []string{viewName, "string"}))
	}

	assignCount := 0
	for _, stmt := range view.GetExpr().GetTransform().GetStmt() {
		if stmt.GetAssign() != nil {
			if assignCount == 0 && stmt.GetAssign().GetName() != viewName {
				v.messages[viewName] = append(v.messages[viewName],
					*NewMsg(ErrMissingReqField, []string{viewName, viewName, "string"}))
			} else if assignCount > 0 {
				v.messages[viewName] = append(v.messages[viewName],
					*NewMsg(ErrExcessAttr, []string{stmt.GetAssign().GetName(), viewName, "string"}))
			}
		}
		assignCount++
	}
}

func (v *Validator) validateViews() {
	for viewName, resolvedTypes := range v.assignTypes {
		typeName := getTypeName(resolvedTypes.refType)
		resolvedType := resolvedTypes.tuple
		if grammarType, exists := v.grammar.Types[typeName]; exists {
			switch t := grammarType.Type.(type) {
			case *sysl.Type_Tuple_:
				v.compareTuple(t.Tuple, resolvedType.GetTuple(),
					getAttrNames(resolvedType.GetTuple().GetAttrDefs()), viewName, typeName)
			default:
				fmt.Println("[validate.validateViews] Unhandled grammar type")
			}

		}
	}

	for viewName, resolvedTypes := range v.letTypes {
		typeName := getTypeName(resolvedTypes.refType)
		resolvedType := resolvedTypes.tuple
		if grammarType, exists := v.grammar.Types[typeName]; exists {
			switch t := grammarType.Type.(type) {
			case *sysl.Type_Tuple_:
				v.compareTuple(t.Tuple, resolvedType.GetTuple(),
					getAttrNames(resolvedType.GetTuple().GetAttrDefs()), viewName, typeName)
			default:
				fmt.Println("[validate.validateViews] Unhandled grammar type")
			}

		}
	}
}

func (v *Validator) compareTuple(
	specTuple, implTuple *sysl.Type_Tuple,
	implAttrNames map[string]struct{},
	viewName, specTupleName string) {
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
					*NewMsg(ErrMissingReqField, []string{gk, viewName, specTupleName}))
			}
		} else {
			delete(implAttrNames, gk)
		}
	}

	for attrName := range implAttrNames {
		v.messages[viewName] = append(v.messages[viewName],
			*NewMsg(ErrExcessAttr, []string{attrName, viewName, specTupleName}))
		delete(implAttrNames, attrName)
	}
}

func (v *Validator) compareOneOf(
	specOneOf *sysl.Type_OneOf,
	implTuple *sysl.Type_Tuple,
	implAttrNames map[string]struct{},
	viewName, specTupleName string) {

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
		} else {
			if _, exists := implAttrs[name]; !exists {
				matching = false
			} else {
				matching = true
				delete(implAttrNames, name)
				break
			}
		}
	}

	if !matching {
		var implAttrNames []string
		for k := range implAttrs {
			implAttrNames = append(implAttrNames, k)
		}
		v.messages[viewName] = append(v.messages[viewName],
			*NewMsg(ErrInvalidOption, []string{viewName, strings.Join(implAttrNames, ","), specTupleName}))
	}
}

func (v *Validator) GetMessages() map[string][]Msg {
	return v.messages
}

func NewValidator(grammar, transform *sysl.Application, parser *Parser) *Validator {
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

func hasSameType(type1 *sysl.Type, type2 *sysl.Type) bool {
	if type1 == nil || type2 == nil {
		return false
	}

	switch type1.GetType().(type) {
	case *sysl.Type_Primitive_:
		return type1.GetPrimitive() == type2.GetPrimitive()
	case *sysl.Type_TypeRef:
		if type2.GetTypeRef() != nil {
			ref1 := type1.GetTypeRef().GetRef()
			ref2 := type2.GetTypeRef().GetRef()

			if ref1.GetAppname() != nil && ref2.GetAppname() != nil {
				return ref1.GetAppname().GetPart()[0] == ref2.GetAppname().GetPart()[0]
			} else if ref1.GetPath() != nil && ref2.GetPath() != nil {
				return ref1.GetPath()[0] == ref2.GetPath()[0]
			}
		}
	case *sysl.Type_Tuple_:
		return type2.GetTuple() != nil
	}

	return false
}

func loadTransform(rootTransform, transformFile string, p *Parser) (*sysl.Application, error) {
	transform, name := loadAndGetDefaultApp(rootTransform, transformFile, p)

	if transform == nil {
		err := errors.New("Unable to load transform")
		return nil, err
	}

	return transform.GetApps()[name], nil
}

func loadGrammar(grammarFile string) (*sysl.Application, error) {
	tokens := strings.Split(grammarFile, ".")
	tokens[len(tokens)-1] = "sysl"
	grammarSyslFile := strings.Join(tokens, ".")
	p := NewParser()

	grammar, name := loadAndGetDefaultApp("", grammarSyslFile, p)
	if grammar == nil {
		err := errors.New("Unable to load grammar-sysl")
		return nil, err
	}
	return grammar.GetApps()[name], nil
}
