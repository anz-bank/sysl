package sequencediagram

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/anz-bank/sysl/pkg/mermaid"
	"github.com/anz-bank/sysl/pkg/sysl"
)

//sequencePair keeps track of the application pairs and the associated endpoint we visit during the recursion
type sequencePair struct {
	firstApp, secondApp, endPoint string
}

const projectDir = "../../../"

var startElse = regexp.MustCompile("^else.*")
var isElse = regexp.MustCompile("^else$")
var isLoop = regexp.MustCompile("^loop.*")

//GenerateSequenceDiagram accepts an application name and an endpoint as inputs and returns a string and an error
//The resulting string is the mermaid code for the sequence diagram for that application and endpoint
func GenerateSequenceDiagram(m *sysl.Module, appname string, epname string) (string, error) {
	return generateSequenceDiagramHelper(m, appname, epname, "...", 1, &[]sequencePair{}, true)
}

//generateSequenceDiagramHelper is a helper which has additional arguments which need not be entered by the user
func generateSequenceDiagramHelper(m *sysl.Module, appName string, epName string,
	previousApp string, indent int, sequencePairs *[]sequencePair, theStart bool) (string, error) {
	var result string
	if theStart {
		result = mermaid.GeneratedHeader + "sequenceDiagram\n"
		if err := isValidAppNameAndEndpoint(m, appName, epName); err != nil {
			return "", err
		}
		result += fmt.Sprintf(" %s ->> %s: %s\n", previousApp, cleanAppName(appName), epName)
	}
	statements := m.Apps[appName].Endpoints[epName].GetStmt()
	result += printSequenceDiagramStatements(m, statements, appName, previousApp, indent, sequencePairs, theStart)
	return result, nil
}

//printSequenceDiagramStatements is where the printing takes place
//Uses a switch statement to decide what to print and what recursion needs to be done
func printSequenceDiagramStatements(m *sysl.Module, statements []*sysl.Statement, appName string,
	previousApp string, indent int, sequencePairs *[]sequencePair, theStart bool) string {
	var result string
	count := 0
	for _, statement := range statements {
		switch c := statement.Stmt.(type) {
		case *sysl.Statement_Group:
			result += fmt.Sprintf("%s%s\n", addIndent(indent), c.Group.Title)
			result += printSequenceDiagramStatements(m, c.Group.Stmt, appName, previousApp, indent+1, sequencePairs, theStart)
			isloop := isLoop.MatchString(c.Group.Title)
			iselse := isElse.MatchString(c.Group.Title)
			if isloop || iselse {
				result += fmt.Sprintf("%send\n", addIndent(indent))
			}
		case *sysl.Statement_Call:
			nextapp := c.Call.Target.Part[0]
			nextep := c.Call.Endpoint
			pair := sequencePair{appName, nextep, nextep}
			if !sequencePairsContain(*sequencePairs, pair) {
				*sequencePairs = append(*sequencePairs, pair)
				result += callStatement(appName, nextep, nextapp, indent)
				previous := appName
				out, err := generateSequenceDiagramHelper(m, nextapp, nextep, previous, indent, sequencePairs, false)
				if err != nil {
					panic("Error in generating sequence diagram; check if app names or endpoints are correct")
				}
				result += out
			}
		case *sysl.Statement_Ret:
			retEndpoint := c.Ret.Payload
			result += retStatement(appName, retEndpoint, previousApp, indent, theStart)
		case *sysl.Statement_Action:
			result += actionStatement(appName, c.Action.Action, indent)
		case *sysl.Statement_Cond:
			result += fmt.Sprintf("%salt %s\n", addIndent(indent), c.Cond.Test)
			result += printSequenceDiagramStatements(m, c.Cond.Stmt, appName, previousApp, indent+1, sequencePairs, theStart)
			if count+1 < len(statements) {
				switch temp := statements[count+1].Stmt.(type) {
				case *sysl.Statement_Group:
					if ok := startElse.MatchString(temp.Group.Title); !ok {
						result += fmt.Sprintf("%send\n", addIndent(indent))
					}
				default:
					result += fmt.Sprintf("%send\n", addIndent(indent))
				}
			} else {
				result += fmt.Sprintf("%send\n", addIndent(indent))
			}
		case *sysl.Statement_Foreach:
			result += fmt.Sprintf("%sloop %s\n", addIndent(indent), c.Foreach.Collection)
			result += printSequenceDiagramStatements(m, c.Foreach.Stmt, appName, previousApp, indent+1, sequencePairs, theStart)
			result += fmt.Sprintf("%send\n", addIndent(indent))
		case *sysl.Statement_Loop:
			result += fmt.Sprintf("%sloop %s\n", addIndent(indent), c.Loop.Criterion)
			result += printSequenceDiagramStatements(m, c.Loop.Stmt, appName, previousApp, indent+1, sequencePairs, theStart)
			result += fmt.Sprintf("%send\n", addIndent(indent))
		case *sysl.Statement_LoopN:
			result += fmt.Sprintf("%sloop %d times\n", addIndent(indent), c.LoopN.Count)
			result += printSequenceDiagramStatements(m, c.LoopN.Stmt, appName, previousApp, indent+1, sequencePairs, theStart)
			result += fmt.Sprintf("%send\n", addIndent(indent))
		default:
			result += ""
		}
		count++
	}
	return result
}

//isValidAppNameAndEndpoint checks if the entered application name and endpoint exists in the sysl module or not
func isValidAppNameAndEndpoint(m *sysl.Module, appName string, epName string) error {
	if _, ok := m.Apps[appName]; !ok {
		return errors.New("invalid app name")
	}
	if _, ok := m.Apps[appName].Endpoints[epName]; !ok {
		return errors.New("invalid endpoint")
	}
	return nil
}

//callStatement is a printer to print a call statement
func callStatement(appName string, epName string, nextApp string, indent int) string {
	return fmt.Sprintf("%s%s ->>+ %s: %s\n", addIndent(indent),
		cleanAppName(appName), cleanAppName(nextApp), epName)
}

//retStatement is a printer to print a return statement
func retStatement(appName string, epName string, previousApp string, indent int, theStart bool) string {
	if theStart {
		return fmt.Sprintf("%s%s -->> %s: %s\n", addIndent(indent),
			cleanAppName(appName), cleanAppName(previousApp), epName)
	}
	return fmt.Sprintf("%s%s -->>- %s: %s\n", addIndent(indent),
		cleanAppName(appName), cleanAppName(previousApp), epName)
}

//actionStatement is a printer to print an action statement
func actionStatement(appName string, action string, indent int) string {
	return fmt.Sprintf("%s%s ->> %s: %s\n", addIndent(indent),
		cleanAppName(appName), cleanAppName(appName), action)
}

//addIndent adds indents based on the input
func addIndent(indent int) string {
	var out string
	for i := 0; i < indent; i++ {
		out += " "
	}
	return out
}

// sequencePairsContain checks if the application-endpoint group have been already visited or not
func sequencePairsContain(s []sequencePair, sp sequencePair) bool {
	for _, a := range s {
		if a == sp {
			return true
		}
	}
	return false
}

func cleanAppName(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}
