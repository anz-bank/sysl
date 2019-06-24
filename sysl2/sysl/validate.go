package main

import (
	"fmt"
	"github.com/anz-bank/sysl/src/proto"
	"github.com/sirupsen/logrus"
	"strings"
)

type content struct {
	Name     string
	TypeRef  string
	Children map[string]content
}

type View struct {
	Name       string
	Body       map[string]content
	OutputType string
	RawData    *sysl.View
}

const (
	ERROR = 100
	WARN  = 200
	INFO  = 300
	UNDEF = 400
)

type ValidationMsg struct {
	message string
	msgType int
}

type Views map[string]View

func typeRef(expr *sysl.Expr) string {
	if expr.GetLiteral() != nil {
		return "string"
	}
	if expr.GetList() != nil {
		return "list"
	}
	if call := expr.GetCall(); call != nil {
		if _, exists := GoFuncMap[call.GetFunc()]; exists {
			return "string"
		}
		return "View"
	}
	if varName := expr.GetName(); varName != "" {
		return varName
	}
	if expr.GetBinexpr() != nil {
		return "BinExpr"
	}
	if expr.GetIfelse() != nil {
		return "Conditional"
	}
	if expr.GetTransform() != nil {
		return "Transform"
	}
	return "Unknown"
}

func readView(stmt *sysl.Expr_Transform_Stmt) content {
	if tfm := stmt.GetAssign().GetExpr().GetTransform(); tfm != nil {
		children := map[string]content{}

		for _, tfmStmt := range tfm.GetStmt() {
			if tfmStmt.GetAssign() != nil {
				child := readView(tfmStmt)
				children[child.Name] = child
			}
		}

		return content{
			Name:     stmt.GetAssign().GetName(),
			TypeRef:  typeRef(stmt.GetAssign().GetExpr()),
			Children: children,
		}
	} else {
		return content{
			Name:    stmt.GetAssign().GetName(),
			TypeRef: typeRef(stmt.GetAssign().GetExpr()),
		}
	}
}

func outputType(syslType *sysl.Type) string {
	if typeName := syslType.GetPrimitive().String(); typeName != "NO_Primitive" {
		return typeName
	} else if typeSeq := syslType.GetSequence(); typeSeq != nil {
		if typeName := typeSeq.GetPrimitive().String(); typeName != "NO_Primitive" {
			return "sequence of " + typeName
		}
		return "sequence of " + typeSeq.GetTypeRef().GetRef().GetPath()[0]
	} else if typeRef := syslType.GetTypeRef(); typeRef != nil {
		return typeRef.GetRef().GetAppname().GetPart()[0]
	} else {
		return "Unknown"
	}
}

func validate(start string, transform *sysl.Module, grm *sysl.Module) {
	views := Views{}

	for _, tfm := range transform.GetApps() {
		for viewName, tfmView := range tfm.Views {
			view := View{
				Name:       viewName,
				Body:       map[string]content{},
				OutputType: outputType(tfmView.GetRetType()),
				RawData:    tfmView,
			}

			for _, stmt := range tfmView.GetExpr().GetTransform().GetStmt() {
				if stmt.GetAssign() != nil {
					view.Body[stmt.GetAssign().GetName()] = readView(stmt)
				}
			}
			views[viewName] = view
		}
	}

	logMsg(validateEntryPoint(views, start))
	logMsg(validateFileName(views))
	logMsg(validateViews(views, grm)...)
}

func logMsg(messages ...ValidationMsg) {
	for _, msg := range messages {
		if msg.msgType == ERROR {
			logrus.Error(msg.message)
		} else if msg.msgType == WARN {
			logrus.Warn(msg.message)
		} else if msg.msgType == INFO {
			logrus.Info(msg.message)
		}
	}
}

func validateEntryPoint(views Views, start string) ValidationMsg {
	view, exists := views[start]

	if !exists {
		return ValidationMsg{
			message: fmt.Sprintf("[Validator]: Entry point view: %s is undefined", start),
			msgType: ERROR,
		}
	}
	if view.OutputType != start {
		return ValidationMsg{
			message: fmt.Sprintf("[Validator]: Output type of entry point view: %s should be %s", start, start),
			msgType: ERROR,
		}
	}
	return ValidationMsg{
		message: "",
		msgType: UNDEF,
	}
}

func validateFileName(views Views) ValidationMsg {
	viewName := "filename"
	view, exists := views[viewName]

	if !exists {
		return ValidationMsg{
			message: "[Validator]: view: filename is undefined",
			msgType: ERROR,
		}
	}
	if view.OutputType != "STRING" {
		return ValidationMsg{
			message: "[Validator]: Output type of view: filename should be string",
			msgType: ERROR,
		}
	}
	if _, exists := view.Body[viewName]; !exists {
		return ValidationMsg{
			message: "[Validator]: In view filename : Missing type: filename",
			msgType: ERROR,
		}
	}
	return ValidationMsg{
		message: "",
		msgType: UNDEF,
	}
}

func validateViews(views Views, grm *sysl.Module) []ValidationMsg {
	var errMsgs []ValidationMsg
	for _, grammarSysl := range grm.GetApps() {

		for _, view := range views {
			outputType := view.OutputType
			typeRef := strings.Replace(outputType, "sequence of", "",-1)
			typeRef = strings.TrimSpace(typeRef)

			if grammarType, exists := grammarSysl.Types[typeRef]; exists {
				errMsgs = append(errMsgs, validateRetType(view, grammarType, outputType != typeRef)...)
			}
		}
	}

	return errMsgs
}

func validateRetType(view View, grammarType *sysl.Type, isSequence bool) []ValidationMsg {
	var errMsgs []ValidationMsg

	for k := range grammarType.GetTuple().GetAttrDefs() {
		if _, exists := view.Body[k]; !exists && !strings.Contains(k, "Choice") {
			errMsgs = append(errMsgs, ValidationMsg{
				message: fmt.Sprintf("[Validator]: In view %s : Missing type: %v", view.Name, k),
				msgType: ERROR,
			})
		}
	}

	if isSequence && view.RawData.GetExpr().GetTransform().GetScopevar() == "." {
		errMsgs = append(errMsgs, ValidationMsg{
			message: fmt.Sprintf("[Validator]: In view %s : Expects sequence type as return value", view.Name),
			msgType: ERROR,
		})
	}

	return errMsgs
}
