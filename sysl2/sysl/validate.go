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

func viewOutput(syslType *sysl.Type) (string, bool) {
	if typeName := syslType.GetPrimitive().String(); typeName != "NO_Primitive" {
		return typeName, false
	} else if typeSeq := syslType.GetSequence(); typeSeq != nil {
		if typeName := typeSeq.GetPrimitive().String(); typeName != "NO_Primitive" {
			return typeName, true
		}
		return typeSeq.GetTypeRef().GetRef().GetPath()[0], true
	} else if typeRef := syslType.GetTypeRef(); typeRef != nil {
		if typeRef.GetRef().GetAppname() != nil {
			return typeRef.GetRef().GetAppname().GetPart()[0], false
		}
		return typeRef.GetRef().GetPath()[0], false
	}
	return "Unknown", false
}

func validateEntryPoint(views map[string]*sysl.View, start string) []ValidationMsg {
	view, exists := views[start]

	if !exists {
		return []ValidationMsg{{
			Message: fmt.Sprintf("Entry point view: '%s' is undefined", start),
			MsgType: ERROR,
		}}
	}
	if t, seq := viewOutput(view.GetRetType()); t != start || seq {
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
	if t, seq := viewOutput(view.GetRetType()); t != "STRING" || seq {
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
	validationMsgs := make([]ValidationMsg, 0, 20)

	for gk, gv := range specAttrs {
		isOneOf := strings.Index(gk, "__Choice_") == 0 && grammarSpec[gk].GetOneOf() != nil

		if isOneOf {
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
	validationMsgs := make([]ValidationMsg, 0, 10)
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

func hasSameType(type1 *sysl.Type, type2 *sysl.Type) (bool, string) {
	if type1 == nil || type2 == nil {
		return false, "nil"
	}
	expectedType := "Unresolved"

	switch type1.GetType().(type) {
	case *sysl.Type_Primitive_:
		expectedType = type1.GetPrimitive().String()
		if _, ok := type2.GetType().(*sysl.Type_Primitive_); ok {
			return type1.GetPrimitive().String() == type2.GetPrimitive().String(), expectedType
		}

	case *sysl.Type_TypeRef:
		switch {
		case type2.GetTypeRef() != nil:
			ref1 := type1.GetTypeRef().GetRef()
			ref2 := type2.GetTypeRef().GetRef()

			if ref1.GetAppname() != nil && ref2.GetAppname() != nil {
				expectedType = ref1.GetAppname().GetPart()[0]
				return ref1.GetAppname().GetPart()[0] == ref2.GetAppname().GetPart()[0], expectedType
			} else if ref1.GetPath() != nil && ref2.GetPath() != nil {
				expectedType = ref1.GetPath()[0]
				return ref1.GetPath()[0] == ref2.GetPath()[0], expectedType
			}
		case type1.GetTypeRef().GetRef().GetAppname() != nil:
			expectedType = type1.GetTypeRef().GetRef().GetAppname().GetPart()[0]
		case type1.GetTypeRef().GetRef().GetPath() != nil:
			expectedType = type1.GetTypeRef().GetRef().GetPath()[0]
		}
	}

	return false, expectedType
}

func resolveVariableType(expr *sysl.Expr, viewName string) (*sysl.Type, []ValidationMsg) {
	validationMsgs := make([]ValidationMsg, 0, 10)

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
		varType, messages := resolveVariableType(expr.GetUnexpr().GetArg(), viewName)
		validationMsgs = append(validationMsgs, messages...)
		if isSame, expected := hasSameType(varType, boolType); !isSame {
			validationMsgs = append(validationMsgs, ValidationMsg{
				Message: fmt.Sprintf("In view '%s', unary operator used with non boolean type: '%s'", viewName, expected),
				MsgType: ERROR})
		}
		return boolType, validationMsgs
	}

	return noType(), validationMsgs
}

func validateTypes(specModule *sysl.Module, expression *sysl.Expr,
	viewName string, implViews map[string]*sysl.View, out string) []ValidationMsg {
	validationMsgs := make([]ValidationMsg, 0, 50)

	attrDefs := map[string]*sysl.Type{}

	for _, stmt := range expression.GetTransform().GetStmt() {
		if stmt.GetAssign() != nil {
			typeName := stmt.GetAssign().GetName()

			expr := stmt.GetAssign().GetExpr()
			attrType, messages := resolveVariableType(expr, viewName)
			attrDefs[typeName] = attrType
			validationMsgs = append(validationMsgs, messages...)

			if expr.GetTransform() != nil {
				attrTypeName, _ := viewOutput(attrType)
				validationMsgs = append(validationMsgs, validateTypes(specModule, expr, viewName, implViews, attrTypeName)...)
			}
		}
	}

	for _, grammarSysl := range specModule.GetApps() {
		if grammarType, exists := grammarSysl.Types[out]; exists {
			validationMsgs = append(
				validationMsgs,
				compareSpecVsImpl(grammarSysl.Types, grammarType.GetTuple().GetAttrDefs(),
					attrDefs, getImplAttrNames(attrDefs), viewName)...)
		}
	}

	return validationMsgs
}

func validateTransform(rootTransform, transform, grammar, start string) []ValidationMsg {
	validationMsgs := make([]ValidationMsg, 0, 100)

	tx, _ := loadAndGetDefaultApp(rootTransform, transform)

	grammarSyslFile := grammar[:len(grammar)-1] + "sysl"
	grammarSysl, _ := loadAndGetDefaultApp("", grammarSyslFile)

	if grammarSysl == nil {
		panic(errors.Errorf("Unable to load grammar-sysl"))
	}

	for _, tfm := range tx.GetApps() {

		validationMsgs = append(validationMsgs, validateEntryPoint(tfm.Views, start)...)
		validationMsgs = append(validationMsgs, validateFileName(tfm.Views)...)

		for viewName, tfmView := range tfm.Views {
			out, _ := viewOutput(tfmView.GetRetType())
			validationMsgs = append(validationMsgs, validateTypes(grammarSysl, tfmView.GetExpr(), viewName, tfm.Views, out)...)
		}
	}

	return validationMsgs
}

func DoValidateTransform(flags *flag.FlagSet, args []string) error {
	rootTransform := flags.String("root-transform", ".", "sysl root directory for input transform file (default: .)")
	transform := flags.String("transform", ".", "transform.sysl")
	grammar := flags.String("grammar", ".", "grammar.g")
	start := flags.String("start", ".", "start rule for the grammar")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	logrus.Infof("root-transform: %s\n", *rootTransform)
	logrus.Infof("transform: %s\n", *transform)
	logrus.Infof("grammar: %s\n", *grammar)
	logrus.Infof("start: %s\n", *start)

	validationMsgs := validateTransform(*rootTransform, *transform, *grammar, *start)

	logMsg(validationMsgs...)

	if len(validationMsgs) > 0 {
		return errors.New("validation failed")
	}

	logMsg(ValidationMsg{Message: "Validation success", MsgType: INFO})
	return nil
}
