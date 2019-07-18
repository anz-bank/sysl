package main

import (
	"flag"
	"fmt"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func noType() *sysl.Type {
	return &sysl.Type{
		Type: &sysl.Type_NoType_{
			NoType: &sysl.Type_NoType{},
		},
	}
}

func getTypeName(syslType *sysl.Type) string {
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

func validateEntryPoint(views map[string]*sysl.View, start string) []Msg {
	view, exists := views[start]

	if !exists {
		return []Msg{*NewMsg(ErrEntryPointUndefined, []string{start})}
	}

	if getTypeName(view.GetRetType()) != start || isCollectionType(view.GetRetType()) {
		return []Msg{*NewMsg(ErrInvalidEntryPointReturn, []string{start, start})}
	}

	return nil
}

func validateFileName(views map[string]*sysl.View) []Msg {
	viewName := "filename"
	view, exists := views[viewName]
	messages := make([]Msg, 0, 2)

	if !exists {
		return []Msg{*NewMsg(ErrUndefinedView, []string{viewName})}
	}

	if getTypeName(view.GetRetType()) != "STRING" || isCollectionType(view.GetRetType()) {
		messages = append(messages, *NewMsg(ErrInvalidReturn, []string{viewName, "string"}))
	}

	assignCount := 0
	for _, stmt := range view.GetExpr().GetTransform().GetStmt() {
		if stmt.GetAssign() != nil {
			if assignCount == 0 && stmt.GetAssign().GetName() != viewName {
				messages = append(messages, *NewMsg(ErrMissingReqField, []string{viewName, viewName, "string"}))
			} else if assignCount > 0 {
				messages = append(messages, *NewMsg(ErrExcessAttr, []string{stmt.GetAssign().GetName(), viewName, "string"}))
			}
		}
		assignCount++
	}

	return messages
}

func getImplAttrNames(attrs map[string]*sysl.Type) map[string]struct{} {
	implAttrNames := map[string]struct{}{}

	for attrName := range attrs {
		implAttrNames[attrName] = struct{}{}
	}

	return implAttrNames
}

func compareTuple(grammarSpec map[string]*sysl.Type,
	specTuple, implTuple *sysl.Type_Tuple,
	implAttrNames map[string]struct{},
	viewName, specTupleName string) []Msg {

	specAttrs := specTuple.GetAttrDefs()
	implAttrs := implTuple.GetAttrDefs()
	messages := make([]Msg, 0, len(implAttrs))

	for gk, gv := range specAttrs {
		if specOneOf := grammarSpec[gk].GetOneOf(); specOneOf != nil {
			messages = append(messages,
				compareOneOf(grammarSpec, specOneOf, implTuple, implAttrNames, viewName, specTupleName)...)
		} else if _, exists := implAttrs[gk]; !exists {
			if !gv.GetOpt() {
				messages = append(messages, *NewMsg(ErrMissingReqField, []string{gk, viewName, specTupleName}))
			}
		} else {
			delete(implAttrNames, gk)
		}
	}

	for attrName := range implAttrNames {
		messages = append(messages, *NewMsg(ErrExcessAttr, []string{attrName, viewName, specTupleName}))
		delete(implAttrNames, attrName)
	}

	return messages
}

func compareOneOf(grammarSpec map[string]*sysl.Type,
	specOneOf *sysl.Type_OneOf,
	implTuple *sysl.Type_Tuple,
	implAttrNames map[string]struct{},
	viewName, specTupleName string) []Msg {

	implAttrs := implTuple.GetAttrDefs()
	messages := make([]Msg, 0, len(implAttrs))
	matching := true

	for _, one := range specOneOf.GetType() {
		name := one.GetTypeRef().GetRef().GetPath()[0]

		if strings.Index(name, "__Choice_Combination_") == 0 {
			if len(implAttrs) == 1 {
				continue
			}

			messages = append(messages, compareTuple(grammarSpec,
				grammarSpec[name].GetTuple(), implTuple, implAttrNames, viewName, specTupleName)...)
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
		messages = append(messages,
			*NewMsg(ErrInvalidOption, []string{viewName, strings.Join(implAttrNames, ","), specTupleName}))
	}

	return messages
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
	}

	return false
}

func resolveExprType(expr *sysl.Expr, viewName string) (*sysl.Type, []Msg) {
	messages := make([]Msg, 0, 1)

	switch e := expr.Expr.(type) {

	case *sysl.Expr_Transform_:
		tfmType := expr.GetType()
		if typeRef := tfmType.GetTypeRef(); typeRef != nil {
			if tfmType.GetTypeRef().GetRef().GetPath() == nil && len(tfmType.GetTypeRef().GetRef().GetAppname().GetPart()) == 1 {
				tfmType.GetTypeRef().GetRef().Path = tfmType.GetTypeRef().GetRef().GetAppname().GetPart()
			}

			return tfmType, messages
		}
		return expr.GetType(), messages
	case *sysl.Expr_Literal:
		return expr.GetType(), messages
	case *sysl.Expr_Unexpr:
		varType, messages := resolveExprType(expr.GetUnexpr().GetArg(), viewName)
		messages = append(messages, messages...)
		switch e.Unexpr.GetOp() {
		case sysl.Expr_UnExpr_NOT, sysl.Expr_UnExpr_INV:
			if !hasSameType(varType, boolType) {
				_, typeDetail := getTypeDetail(varType)
				messages = append(messages, *NewMsg(ErrInvalidUnary, []string{viewName, typeDetail}))
			}
			return boolType, messages
		case sysl.Expr_UnExpr_NEG, sysl.Expr_UnExpr_POS:
			if !hasSameType(varType, intType) {
				_, typeDetail := getTypeDetail(varType)
				messages = append(messages, *NewMsg(ErrInvalidUnary, []string{viewName, typeDetail}))
			}
			return intType, messages
		}
	}

	return noType(), messages
}

func validateTransform(specApp *sysl.Application,
	transform *sysl.Expr_Transform,
	viewName string,
	implViews map[string]*sysl.View,
	typeName string) []Msg {

	messages := make([]Msg, 0, len(transform.GetStmt()))
	newTuple := &sysl.Type_Tuple{
		AttrDefs: map[string]*sysl.Type{},
	}
	attrDefs := newTuple.AttrDefs

	for _, stmt := range transform.GetStmt() {
		if stmt.GetAssign() != nil {
			varName := stmt.GetAssign().GetName()

			expr := stmt.GetAssign().GetExpr()
			exprType, messages1 := resolveExprType(expr, viewName)
			attrDefs[varName] = exprType
			messages = append(messages, messages1...)

			if innerTfm := expr.GetTransform(); innerTfm != nil {
				attrTypeName := getTypeName(exprType)
				messages = append(messages, validateTransform(specApp, innerTfm, viewName, implViews, attrTypeName)...)
			}
		}
	}

	if grammarType, exists := specApp.Types[typeName]; exists {
		switch t := grammarType.Type.(type) {
		case *sysl.Type_Tuple_:
			messages = append(
				messages,
				compareTuple(specApp.Types, t.Tuple, newTuple, getImplAttrNames(attrDefs), viewName, typeName)...)
		default:
			fmt.Println("[validate.validateTransform] Unhandled grammar type")
		}

	}

	return messages
}

func validate(grammar, transform *sysl.Application, start string) []Msg {
	messages := make([]Msg, 0, 10)

	messages = append(messages, validateEntryPoint(transform.Views, start)...)
	messages = append(messages, validateFileName(transform.Views)...)

	for viewName, tfmView := range transform.Views {
		typeName := getTypeName(tfmView.GetRetType())
		messages = append(messages,
			validateTransform(grammar, tfmView.GetExpr().GetTransform(), viewName, transform.Views, typeName)...)
	}

	return messages
}

func loadTransform(rootTransform, transformFile string) (*sysl.Application, error) {
	transform, name := loadAndGetDefaultApp(rootTransform, transformFile)

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

	grammar, name := loadAndGetDefaultApp("", grammarSyslFile)
	if grammar == nil {
		err := errors.New("Unable to load grammar-sysl")
		return nil, err
	}
	return grammar.GetApps()[name], nil
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

	transform, err := loadTransform(*rootTransform, *transformFile)
	if err != nil {
		return err
	}

	messages := validate(grammar, transform, *start)

	for _, message := range messages {
		message.logMsg()
	}

	if len(messages) > 0 {
		NewMsg(ErrValidationFailed, nil).logMsg()
		return errors.New("validation failed")
	}

	NewMsg(InfoValidatedSuccessfully, nil).logMsg()

	return nil
}
