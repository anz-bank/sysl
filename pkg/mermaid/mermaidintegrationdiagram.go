package mermaid

import (
	"errors"
	"fmt"

	"github.com/anz-bank/sysl/pkg/sysl"
)

func GenerateMermaidIntegrationDiagram(m *sysl.Module, appname string) (string, error) {
	return generateMermaidIntegrationDiagramHelper(m, appname, true)
}

func generateMermaidIntegrationDiagramHelper(m *sysl.Module, appname string, thestart bool) (string, error) {
	var result string
	if thestart {
		result = "%% AUTOGENERATED CODE -- DO NOT EDIT!\n\ngraph TD\n"
		if err := validAppName(m, appname); err != nil {
			return "", err
		}
	}
	endpoints := m.Apps[appname].Endpoints
	for _, value := range endpoints {
		statements := value.Stmt
		result += printIntegrationDiagramStatements(m, statements, appname)
	}
	return result, nil
}

func printIntegrationDiagramStatements(m *sysl.Module, statements []*sysl.Statement, appname string) string {
	var result string
	for _, statement := range statements {
		switch c := statement.Stmt.(type) {
		case *sysl.Statement_Call:
			if appname != c.Call.Target.Part[0] {
				result += fmt.Sprintf(" %s --> %s\n", appname, c.Call.Target.Part[0])
				out, err := generateMermaidIntegrationDiagramHelper(m, c.Call.Target.Part[0], false)
				if err != nil {
					panic("Error in generating integration diagram; check if app name is correct")
				}
				result += out
			}
		case *sysl.Statement_Group:
			result += printIntegrationDiagramStatements(m, c.Group.Stmt, appname)
		case *sysl.Statement_Cond:
			result += printIntegrationDiagramStatements(m, c.Cond.Stmt, appname)
		case *sysl.Statement_Loop:
			result += printIntegrationDiagramStatements(m, c.Loop.Stmt, appname)
		case *sysl.Statement_LoopN:
			result += printIntegrationDiagramStatements(m, c.LoopN.Stmt, appname)
		case *sysl.Statement_Foreach:
			result += printIntegrationDiagramStatements(m, c.Foreach.Stmt, appname)
		case *sysl.Statement_Action:
			result += ""
		case *sysl.Statement_Ret:
			result += ""
		default:
			panic("Unrecognised statement type")
		}
	}
	return result
}

func validAppName(m *sysl.Module, appname string) error {
	if _, ok := m.Apps[appname]; !ok {
		return errors.New("invalid app name")
	}
	return nil
}
