package endpointanalysisdiagram

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/mermaid"
	"github.com/anz-bank/sysl/pkg/sysl"
)

//externalLink keeps track of the statement-endpoint pairs we visit during execution
type externalLink struct {
	statement, endPoint string
}

//GenerateEndpointAnalysisDiagram accepts the sysl module and returns a string (and an error if any)
//The resulting string is the mermaid code for the endpoint analysis for that application and endpoint
func GenerateEndpointAnalysisDiagram(m *sysl.Module) (string, error) {
	return generateEndpointAnalysisDiagramHelper(m, &[]externalLink{}, true)
}

//Similar to the above, but accepts a slice of application names
//and returns a diagram only including applications specified
func GenerateMultipleAppEndpointAnalysisDiagram(m *sysl.Module, appNames []string) (string, error) {
	return generateMultipleAppEndpointAnalysisDiagramHelper(m, appNames, &[]externalLink{}, true)
}

//generateEndpointAnalysisDiagram is a helper which has additional arguments which need not be entered by the user
func generateEndpointAnalysisDiagramHelper(m *sysl.Module,
	externalLinks *[]externalLink, theStart bool) (string, error) {
	var result string
	if theStart {
		result = mermaid.GeneratedHeader + "graph TD\n"
	}
	count := 1
	for appName, app := range m.Apps {
		result += fmt.Sprintf(" subgraph %d[\"%s\"]\n", count, appName)
		for epName, endPoint := range app.Endpoints {
			statements := endPoint.Stmt
			result += printEndpointAnalysisStatements(m, statements, mermaid.CleanString(epName), externalLinks)
		}
		result += " end\n"
		count++
	}
	for _, eLink := range *externalLinks {
		result += fmt.Sprintf(" %s --> %s\n", eLink.statement, eLink.endPoint)
	}
	return result, nil
}

func generateMultipleAppEndpointAnalysisDiagramHelper(m *sysl.Module, appNames []string,
	externalLinks *[]externalLink, theStart bool) (string, error) {
	var result string
	if theStart {
		result = mermaid.GeneratedHeader + "graph TD\n"
	}
	count := 1
	for _, appName := range appNames {
		result += fmt.Sprintf(" subgraph %d[\"%s\"]\n", count, appName)
		endPoints := m.Apps[appName].Endpoints
		for epName, endPoint := range endPoints {
			statements := endPoint.Stmt
			result += printEndpointAnalysisStatements(m, statements, mermaid.CleanString(epName), externalLinks)
		}
		result += " end\n"
		count++
	}
	for _, eLink := range *externalLinks {
		result += fmt.Sprintf(" %s --> %s\n", eLink.statement, eLink.endPoint)
	}
	return result, nil
}

//printEndpointAnalysisStatements is used to print the mermaid code for different sysl statements
func printEndpointAnalysisStatements(m *sysl.Module, statements []*sysl.Statement,
	endPoint string, externalLinks *[]externalLink) string {
	var result string
	for _, statement := range statements {
		switch c := statement.Stmt.(type) {
		case *sysl.Statement_Call:
			appEndPoint := fmt.Sprintf("%s-%s", mermaid.CleanString(c.Call.Target.Part[0]), mermaid.CleanString(c.Call.Endpoint))
			result += fmt.Sprintf("  %s --> %s\n", endPoint, appEndPoint)
			pair := externalLink{appEndPoint, mermaid.CleanString(c.Call.Endpoint)}
			if !externalLinksContain(*externalLinks, pair) {
				*externalLinks = append(*externalLinks, pair)
			}
		case *sysl.Statement_Group:
			result += printEndpointAnalysisStatements(m, c.Group.Stmt, endPoint, externalLinks)
		case *sysl.Statement_Cond:
			result += printEndpointAnalysisStatements(m, c.Cond.Stmt, endPoint, externalLinks)
		case *sysl.Statement_Loop:
			result += printEndpointAnalysisStatements(m, c.Loop.Stmt, endPoint, externalLinks)
		case *sysl.Statement_LoopN:
			result += printEndpointAnalysisStatements(m, c.LoopN.Stmt, endPoint, externalLinks)
		case *sysl.Statement_Foreach:
			result += printEndpointAnalysisStatements(m, c.Foreach.Stmt, endPoint, externalLinks)
		case *sysl.Statement_Action:
			result += fmt.Sprintf("  %s --> %s\n", endPoint, mermaid.CleanString(c.Action.Action))
		case *sysl.Statement_Ret:
			result += fmt.Sprintf("  %s --> %s\n", endPoint, mermaid.CleanString(c.Ret.Payload))
		default:
			result += ""
		}
	}
	return result
}

//externalLinksContain checks if the statement-endpoint group have been already visited or not
func externalLinksContain(i []externalLink, ip externalLink) bool {
	for _, a := range i {
		if a == ip {
			return true
		}
	}
	return false
}
