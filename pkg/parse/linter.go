package parse

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	parser "github.com/anz-bank/sysl/pkg/grammar"
	"github.com/sirupsen/logrus"
)

// linterRecords records any data that is required to do linting
type linterRecords struct {
	apps  map[string]*graph
	calls *graph
}

type graph map[string]*graphData

type graphData struct {
	locations map[string]bool
	rec       *graph
}

func (s *TreeShapeListener) lint() {
	s.lintMode = true
	s.linter = newLinterRecords()
}

func newAppEndpointGraph() *graph {
	rec := new(graph)
	*rec = make(map[string]*graphData)
	return rec
}

func newLocation(location string) map[string]bool {
	loc := make(map[string]bool)
	loc[location] = true
	return loc
}

func (g *graph) recordApp(appName, location string) error {
	if app, exists := (*g)[appName]; !exists {
		(*g)[appName] = &graphData{newLocation(location), newAppEndpointGraph()}
		return nil
	} else if !app.locations[location] {
		//TODO: find better way, to handle collector
		app.locations[location] = true
		return nil
	}
	return fmt.Errorf("recordApp: app already exists: %s %s", location, appName)
}

func (g *graph) recordEndpoint(appName, endpoint, location string) error {
	if app, exists := (*g)[appName]; exists {
		if e, exists := (*app.rec)[endpoint]; !exists {
			(*app.rec)[endpoint] = &graphData{newLocation(location), nil}
			return nil
		} else if !e.locations[location] {
			e.locations[location] = true
			return nil
		}
		return fmt.Errorf("recordEndpoint: endpoint already exists: %s %s %s", location, appName, endpoint)
	}
	return fmt.Errorf("recordEndpoint: app does not exist: %s", appName)
}

func (g *graph) recordMethod(appName, endpoint, method, location string) error {
	if app, exists := (*g)[appName]; exists {
		if e, exists := (*app.rec)[endpoint]; exists {
			if e.rec == nil {
				e.rec = newAppEndpointGraph()
			}
			if _, exists := (*e.rec)[method]; !exists {
				(*e.rec)[method] = &graphData{newLocation(location), nil}
				return nil
			}
			return fmt.Errorf("recordMethod: method already exist: %s %s <- %s %s", location, appName, method, endpoint)
		}
	}
	return fmt.Errorf("recordMethod: app does not exist: %s %s", location, appName)
}

func (g *graph) recordAsCall(appName, endpoint, method, location string) error {
	g.recordApp(appName, "") //nolint:errcheck
	if method == "" {
		if err := g.recordEndpoint(appName, endpoint, location); err != nil {
			return err
		}
		return nil
	}
	g.recordEndpoint(appName, endpoint, location) //nolint:errcheck
	if err := g.recordMethod(appName, endpoint, method, location); err != nil {
		app := (*g)[appName]
		endpoints := (*app.rec)[endpoint]
		methods := (*endpoints.rec)[method]
		if methods.locations[location] {
			// this isn't possible
			return fmt.Errorf("recordAsCall: location already exists")
		}
		methods.locations[location] = true
	}
	return nil
}

func newLinterRecords() *linterRecords {
	return &linterRecords{
		make(map[string]*graph),
		newAppEndpointGraph(),
	}
}

func (s *TreeShapeListener) getApps() *graph {
	appName := s.getFullAppName()
	if apps, exists := s.linter.apps[strings.ToLower(appName)]; exists {
		return apps
	}
	rec := newAppEndpointGraph()
	s.linter.apps[strings.ToLower(appName)] = rec
	return rec
}

func (s *TreeShapeListener) createLocation(lineNum, colNum int) string {
	return fmt.Sprintf("%s:%d:%d", s.sc.filename, lineNum, colNum)
}

func (s *TreeShapeListener) recordApp(location string) {
	if !s.lintMode {
		return
	}
	appName := s.getFullAppName()
	if err := s.getApps().recordApp(appName, location); err != nil {
		logrus.Fatal(err)
	}
}

func (s *TreeShapeListener) getFullAppName() string {
	// handle subpackage
	return strings.Join(s.currentApp().Name.Part, "::")
}

func (s *TreeShapeListener) recordEndpoint(endpoint, location string) {
	if !s.lintMode {
		return
	}
	appName := s.getFullAppName()
	if err := s.getApps().recordEndpoint(appName, endpoint, location); err != nil {
		logrus.Fatal(err)
	}
}

func (s *TreeShapeListener) recordMethod(endpoint, method, location string) {
	if !s.lintMode {
		return
	}

	appName := s.getFullAppName()
	s.getApps().recordEndpoint(appName, endpoint, location) //nolint:errcheck
	if err := s.getApps().recordMethod(appName, endpoint, method, location); err != nil {
		logrus.Warn(err)
	}
}

func (s *TreeShapeListener) recordCall(appName, endpoint, method, location string) {
	if !s.lintMode {
		return
	}
	if err := s.linter.calls.recordAsCall(appName, endpoint, method, location); err != nil {
		logrus.Warn(err)
	}
}

func (s *TreeShapeListener) lintEndpoint() {
	appNotExistLog := func(location, appName, call string) {
		logrus.Warnf("lint %s: Application '%s' does not exist for call '%s'", location, appName, call)
	}
	lint := func(appName, endpoint, method, call string, locations map[string]bool) {
		for location := range locations {
			apps, exists := s.linter.apps[strings.ToLower(appName)]
			if !exists {
				appNotExistLog(location, appName, call)
				continue
			}
			if app, exists := (*apps)[appName]; exists {
				if endpoints, exists := (*app.rec)[endpoint]; exists {
					// if method is empty string, it is linting simple endpoint, not REST endpoint
					if method != "" {
						if _, exists = (*endpoints.rec)[method]; exists {
							continue
						}
						logrus.Warnf("lint %s: Method '%s' does not exist for call '%s'", location, method, call)
						continue
					}
				}
				logrus.Warnf("lint %s: Endpoint '%s' does not exist for call '%s'", location, endpoint, call)
				continue
			}
			appNotExistLog(location, appName, call)
		}
	}

	for appName, appData := range *s.linter.calls {
		for endpoint, endpointData := range *appData.rec {
			// if method maps == nil, it is not REST endpoint
			if endpointData.rec == nil {
				lint(
					appName,
					endpoint,
					"",
					fmt.Sprintf("%s <- %s", appName, endpoint),
					endpointData.locations,
				)
				continue
			}

			for method, methodData := range *endpointData.rec {
				lint(
					appName,
					endpoint,
					method,
					fmt.Sprintf("%s <- %s %s", appName, method, endpoint),
					methodData.locations,
				)
			}
		}
	}
}

func (s *TreeShapeListener) lintAppDefs() {
	for _, apps := range s.linter.apps {
		if len(*apps) > 1 {
			appDefs := make([]string, 0, len(*apps))
			for a, data := range *apps {
				locations := make([]string, 0, len(data.locations))
				for location := range data.locations {
					locations = append(locations, location)
				}
				appDefs = append(appDefs, fmt.Sprintf("%s:%s", a, strings.Join(locations, ", ")))
			}
			sort.Strings(appDefs)
			logrus.Warnf(
				"lint: case-sensitive redefinitions detected:\n%s",
				strings.Join(appDefs, "\n"),
			)
		}
	}
}

func (s *TreeShapeListener) lintRetStmt(payload string, ctx *parser.Ret_stmtContext) {
	if !regexp.MustCompile(`<:`).MatchString(payload) {
		if regexp.MustCompile(`(ok|error|[1-5][0-9][0-9])`).MatchString(payload) {
			return
		}
		logrus.Warnf(
			"lint %s: 'return %s' not supported, use 'return ok <: %[2]s' instead",
			s.createLocation(ctx.GetStart().GetLine(), ctx.GetStart().GetColumn()),
			payload)
	}
}
