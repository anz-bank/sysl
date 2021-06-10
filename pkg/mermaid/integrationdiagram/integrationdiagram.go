package integrationdiagram

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anz-bank/sysl/pkg/mermaid"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
)

// integrationPair keeps track of the application pairs we visit during the recursion
type integrationPair struct {
	firstApp, secondApp string
}

// GenerateFullIntegrationDiagram returns the full integration diagram for a sysl module
func GenerateFullIntegrationDiagram(m *sysl.Module) (string, error) {
	return generateFullIntegrationDiagramHelper(m, &[]integrationPair{})
}

// GenerateIntegrationDiagram accepts an application name as input and returns a string (and an error if any)
// The resulting string is the mermaid code for the integration diagram
func GenerateIntegrationDiagram(m *sysl.Module, appName string) (string, error) {
	return generateIntegrationDiagramHelper(m, appName, &[]integrationPair{}, true)
}

func GenerateMultipleAppIntegrationDiagram(m *sysl.Module, appNames []string) (string, error) {
	return generateMultipleAppIntegrationDiagramHelper(m, mermaid.CleanAppNames(appNames), &[]integrationPair{})
}

// generateEntireIntegrationDiagramHelper is a helper which is generates an entire integration diagram
func generateFullIntegrationDiagramHelper(m *sysl.Module,
	integrationPairs *[]integrationPair) (string, error) {
	var result string
	result = mermaid.GeneratedHeader + "graph TD\n"
	for appName, appValue := range m.Apps {
		endPoints := appValue.Endpoints
		for _, endPoint := range endPoints {
			statements := endPoint.Stmt
			result += printIntegrationDiagramStatements(m, statements, appName, integrationPairs)
		}
	}
	return result, nil
}

// generateIntegrationDiagramHelper accepts an application name and returns the respective integration diagram
func generateIntegrationDiagramHelper(m *sysl.Module, appName string,
	integrationPairs *[]integrationPair, theStart bool) (string, error) {
	var result string
	if theStart {
		result = mermaid.GeneratedHeader + "graph TD\n"
		if err := IsValidAppName(m, appName); err != nil {
			return "", err
		}
	}
	endPoints := m.Apps[appName].Endpoints
	// For every endpoint, the statements are retrieved and we pass it to the printer to print appropriate mermaid code
	for _, endPoint := range endPoints {
		statements := endPoint.Stmt
		result += printIntegrationDiagramStatements(m, statements, appName, integrationPairs)
	}
	return result, nil
}

func generateMultipleAppIntegrationDiagramHelper(m *sysl.Module, appNames []string,
	integrationPairs *[]integrationPair) (string, error) {
	var result string
	result = mermaid.GeneratedHeader + "graph TD\n"
	for _, appName := range appNames {
		if app := m.Apps[appName]; app != nil {
			endPoints := app.Endpoints
			result += printClassStatement(appName)
			for _, endPoint := range endPoints {
				statements := endPoint.Stmt
				result += printIntegrationDiagramStatements(m, statements, appName, integrationPairs)
			}
		}
	}

	// Get all applications which call an App in appNames
	for currentApp, appValue := range m.Apps {
		endPoints := appValue.Endpoints
		for _, endPoint := range endPoints {
			for _, targetApp := range appNames {
				result += printIntegrationDiagramStatementsTargetedApp(m, endPoint.Stmt, currentApp, integrationPairs, targetApp)
			}
		}
	}
	return result, nil
}

func printClassStatement(className string) string {
	return fmt.Sprintf("    %s\n", cleanAppName(className))
}

// printIntegrationDiagramStatements is where the printing takes place
// Uses a switch statement to decide what to print and what recursion needs to be done
func printIntegrationDiagramStatements(m *sysl.Module, statements []*sysl.Statement,
	appName string, integrationPairs *[]integrationPair) string {
	var result string
	for _, statement := range statements {
		switch c := statement.Stmt.(type) {
		case *sysl.Statement_Call:
			nextApp := syslutil.GetAppName(c.Call.Target)
			pair := integrationPair{appName, nextApp}
			if !integrationPairsContain(*integrationPairs, pair) {
				*integrationPairs = append(*integrationPairs, pair)
				result += fmt.Sprintf(" %s[\"%s\"] --> %s[\"%s\"]\n",
					cleanAppName(appName), appName, cleanAppName(nextApp), nextApp)
				out, err := generateIntegrationDiagramHelper(m, nextApp, integrationPairs, false)
				if err != nil {
					panic("Error in generating integration diagram; check if app name is correct")
				}
				result += out
			}
		case *sysl.Statement_Group:
			result += printIntegrationDiagramStatements(m, c.Group.Stmt, appName, integrationPairs)
		case *sysl.Statement_Cond:
			result += printIntegrationDiagramStatements(m, c.Cond.Stmt, appName, integrationPairs)
		case *sysl.Statement_Loop:
			result += printIntegrationDiagramStatements(m, c.Loop.Stmt, appName, integrationPairs)
		case *sysl.Statement_LoopN:
			result += printIntegrationDiagramStatements(m, c.LoopN.Stmt, appName, integrationPairs)
		case *sysl.Statement_Foreach:
			result += printIntegrationDiagramStatements(m, c.Foreach.Stmt, appName, integrationPairs)
		case *sysl.Statement_Action:
			result += ""
		case *sysl.Statement_Ret:
			result += ""
		default:
			result += ""
		}
	}
	return result
}

// printIntegrationDiagramStatements is where the printing takes place
// Uses a switch statement to decide what to print and what recursion needs to be done
func printIntegrationDiagramStatementsTargetedApp(m *sysl.Module, statements []*sysl.Statement,
	appName string, integrationPairs *[]integrationPair, targetAppName string) string {
	var result string
	for _, statement := range statements {
		switch c := statement.Stmt.(type) {
		case *sysl.Statement_Call:
			nextApp := syslutil.GetAppName(c.Call.Target)
			pair := integrationPair{appName, nextApp}
			if !integrationPairsContain(*integrationPairs, pair) && nextApp == targetAppName {
				*integrationPairs = append(*integrationPairs, pair)
				result += fmt.Sprintf(" %s[\"%s\"] --> %s[\"%s\"]\n",
					cleanAppName(appName), appName, cleanAppName(nextApp), nextApp)
				out, err := generateIntegrationDiagramHelper(m, nextApp, integrationPairs, false)
				if err != nil {
					panic("Error in generating integration diagram; check if app name is correct")
				}
				result += out
			}
		case *sysl.Statement_Group:
			result += printIntegrationDiagramStatementsTargetedApp(m, c.Group.Stmt, appName, integrationPairs, targetAppName)
		case *sysl.Statement_Cond:
			result += printIntegrationDiagramStatementsTargetedApp(m, c.Cond.Stmt, appName, integrationPairs, targetAppName)
		case *sysl.Statement_Loop:
			result += printIntegrationDiagramStatementsTargetedApp(m, c.Loop.Stmt, appName, integrationPairs, targetAppName)
		case *sysl.Statement_LoopN:
			result += printIntegrationDiagramStatementsTargetedApp(m, c.LoopN.Stmt, appName, integrationPairs, targetAppName)
		case *sysl.Statement_Foreach:
			result += printIntegrationDiagramStatementsTargetedApp(m, c.Foreach.Stmt, appName, integrationPairs, targetAppName)
		case *sysl.Statement_Action:
			result += ""
		case *sysl.Statement_Ret:
			result += ""
		default:
			result += ""
		}
	}
	return result
}

// IsValidAppName checks if the entered application name exists in the sysl module or not
func IsValidAppName(m *sysl.Module, appName string) error {
	if _, ok := m.Apps[appName]; !ok {
		return errors.New("invalid app name")
	}
	return nil
}

// integrationPairsContain checks if the application couple have been already visited or not
func integrationPairsContain(i []integrationPair, ip integrationPair) bool {
	for _, a := range i {
		if a == ip {
			return true
		}
	}
	return false
}

func cleanAppName(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}
