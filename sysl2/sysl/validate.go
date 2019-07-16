package main

import (
	"flag"
	"fmt"
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ERROR = 100
	WARN  = 200
	INFO  = 300
	UNDEF = 400
)

type ValidationMsg struct {
	Message string `json:"message"`
	MsgType int    `json:"msg_type"`
}

func logMsg(messages ...ValidationMsg) {
	for _, msg := range messages {
		formattedMsg := "[Validator]: " + msg.Message
		switch msg.MsgType {
		case ERROR:
			logrus.Error(formattedMsg)
		case WARN:
			logrus.Warn(formattedMsg)
		case INFO:
			logrus.Info(formattedMsg)
		}
	}
}

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

func validateEntryPoint(views map[string]*sysl.View, start string) []ValidationMsg {
	view, exists := views[start]

	if !exists {
		return []ValidationMsg{{
			Message: fmt.Sprintf("Entry point view: '%s' is undefined", start),
			MsgType: ERROR,
		}}
	}

	if getTypeName(view.GetRetType()) != start || isCollectionType(view.GetRetType()) {
		return []ValidationMsg{{
			Message: fmt.Sprintf("Output type of entry point view: '%s' should be '%s'", start, start),
			MsgType: ERROR,
		}}
	}

	return nil
}

func validateFileName(views map[string]*sysl.View) []ValidationMsg {
	viewName := "filename"
	view, exists := views[viewName]
	validationMsgs := make([]ValidationMsg, 0, 2)

	if !exists {
		return []ValidationMsg{{
			Message: "View 'filename' is undefined",
			MsgType: ERROR,
		}}
	}

	if getTypeName(view.GetRetType()) != "STRING" || isCollectionType(view.GetRetType()) {
		validationMsgs = append(validationMsgs, ValidationMsg{
			Message: "In view 'filename', output type should be 'string'",
			MsgType: ERROR,
		})
	}

	assignCount := 0
	for _, stmt := range view.GetExpr().GetTransform().GetStmt() {
		if stmt.GetAssign() != nil {
			if assignCount == 0 && stmt.GetAssign().GetName() != viewName {
				validationMsgs = append(validationMsgs, ValidationMsg{
					Message: "In view 'filename' : missing type: 'filename'",
					MsgType: ERROR,
				})
			} else if assignCount > 0 {
				validationMsgs = append(validationMsgs, ValidationMsg{
					Message: fmt.Sprintf("In view 'filename' : Excess assignments: '%s'", stmt.GetAssign().GetName()),
					MsgType: ERROR,
				})
			}
			assignCount++
		}
	}
	return validationMsgs
}

func getImplAttrNames(attrs map[string]*sysl.Type) map[string]struct{} {
	implAttrNames := map[string]struct{}{}

	for attrName := range attrs {
		implAttrNames[attrName] = struct{}{}
	}

	return implAttrNames
}

func compareSpecVsImpl(grammarSpec, specAttrs, implAttrs map[string]*sysl.Type,
	implAttrNames map[string]struct{}, viewName string) []ValidationMsg {
	validationMsgs := make([]ValidationMsg, 0, len(implAttrs))

	for gk, gv := range specAttrs {
		if grammarSpec[gk].GetOneOf() != nil {
			validationMsgs = append(validationMsgs, compareOneOf(grammarSpec, implAttrs, implAttrNames, viewName, gk)...)
		} else if _, exists := implAttrs[gk]; !exists {
			if !gv.GetOpt() {
				validationMsgs = append(validationMsgs, ValidationMsg{
					Message: fmt.Sprintf("In view '%s', type '%s' is missing", viewName, gk),
					MsgType: ERROR,
				})
			}
		} else {
			delete(implAttrNames, gk)
		}
	}

	for attrName := range implAttrNames {
		validationMsgs = append(validationMsgs, ValidationMsg{
			Message: fmt.Sprintf("In view '%s', excess attribute is defined: '%s'", viewName, attrName),
			MsgType: ERROR,
		})
		delete(implAttrNames, attrName)
	}

	return validationMsgs
}

func compareOneOf(grammarSpec, implAttrs map[string]*sysl.Type,
	implAttrNames map[string]struct{}, viewName, typeName string) []ValidationMsg {
	validationMsgs := make([]ValidationMsg, 0, len(implAttrs))
	matching := true

	for _, one := range grammarSpec[typeName].GetOneOf().GetType() {
		name := one.GetTypeRef().GetRef().GetPath()[0]

		if strings.Index(name, "__Choice_Combination_") == 0 {
			if len(implAttrs) == 1 {
				continue
			}

			validationMsgs = append(validationMsgs, compareSpecVsImpl(grammarSpec,
				grammarSpec[name].GetTuple().GetAttrDefs(), implAttrs, implAttrNames, viewName)...)
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
		implAttrName := ""
		for k := range implAttrs {
			implAttrName = k
		}
		validationMsgs = append(validationMsgs, ValidationMsg{
			Message: fmt.Sprintf("In view '%s', invalid choice has been made as : '%s'", viewName, implAttrName),
			MsgType: ERROR})
	}

	return validationMsgs
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

func resolveExprType(expr *sysl.Expr, viewName string) (*sysl.Type, []ValidationMsg) {
	validationMsgs := make([]ValidationMsg, 0, 1)

	switch expr.Expr.(type) {

	case *sysl.Expr_Transform_:
		if typeRef := expr.GetType().GetTypeRef(); typeRef != nil {
			return &sysl.Type{
				Type: &sysl.Type_TypeRef{
					TypeRef: &sysl.ScopedRef{
						Ref: &sysl.Scope{Path: typeRef.GetRef().GetAppname().GetPart()},
					},
				},
			}, validationMsgs
		}
		return expr.GetType(), validationMsgs
	case *sysl.Expr_Literal:
		return expr.GetType(), validationMsgs
	case *sysl.Expr_Unexpr:
		varType, messages := resolveExprType(expr.GetUnexpr().GetArg(), viewName)
		validationMsgs = append(validationMsgs, messages...)
		if !hasSameType(varType, boolType) {
			_, typeDetail := getTypeDetail(varType)
			validationMsgs = append(validationMsgs, ValidationMsg{
				Message: fmt.Sprintf("In view '%s', unary operator used with non boolean type: '%s'", viewName, typeDetail),
				MsgType: ERROR})
		}
		return boolType, validationMsgs
	}

	return noType(), validationMsgs
}

func validateTransform(specApp *sysl.Application, transform *sysl.Expr_Transform,
	viewName string, implViews map[string]*sysl.View, typeName string) []ValidationMsg {
	validationMsgs := make([]ValidationMsg, 0, 50)

	attrDefs := map[string]*sysl.Type{}

	for _, stmt := range transform.GetStmt() {
		if stmt.GetAssign() != nil {
			varName := stmt.GetAssign().GetName()

			expr := stmt.GetAssign().GetExpr()
			exprType, messages := resolveExprType(expr, viewName)
			attrDefs[varName] = exprType
			validationMsgs = append(validationMsgs, messages...)

			if innerTfm := expr.GetTransform(); innerTfm != nil {
				attrTypeName := getTypeName(exprType)
				validationMsgs = append(validationMsgs, validateTransform(specApp, innerTfm, viewName, implViews, attrTypeName)...)
			}
		}
	}

	if grammarType, exists := specApp.Types[typeName]; exists {
		validationMsgs = append(
			validationMsgs,
			compareSpecVsImpl(specApp.Types, grammarType.GetTuple().GetAttrDefs(),
				attrDefs, getImplAttrNames(attrDefs), viewName)...)
	}

	return validationMsgs
}

func validate(grammar, transform *sysl.Application, start string) []ValidationMsg {
	validationMsgs := make([]ValidationMsg, 0, 100)

	validationMsgs = append(validationMsgs, validateEntryPoint(transform.Views, start)...)
	validationMsgs = append(validationMsgs, validateFileName(transform.Views)...)

	for viewName, tfmView := range transform.Views {
		typeName := getTypeName(tfmView.GetRetType())
		validationMsgs = append(validationMsgs,
			validateTransform(grammar, tfmView.GetExpr().GetTransform(), viewName, transform.Views, typeName)...)
	}

	return validationMsgs
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

	validationMsgs := validate(grammar, transform, *start)

	logMsg(validationMsgs...)

	if len(validationMsgs) > 0 {
		return errors.New("validation failed")
	}

	logMsg(ValidationMsg{Message: "Validation success", MsgType: INFO})
	return nil
}
