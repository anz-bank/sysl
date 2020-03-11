// Package catalog takes a sysl module with attributes defined (catalogFields)
// and serves a webserver listing the applications and endpoints
// It also uses GRPCUI and Redoc in order to generate an interactive page to interact with all the endpoints
// GRPC currently uses server reflection TODO: Support gpcui directly from swagger files
package catalog

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"net/http"

	"github.com/anz-bank/sysl/pkg/sysl"
)

// Server to set context of catalog
// Todo: Simplify this

type SyslUI struct {
	Fs       afero.Fs       // Required
	Log      *logrus.Logger // Required
	Modules  []*sysl.Module // Required
	Fields   []string       // Required
	Host     string         // Required
	GRPCUI   bool
	BasePath string
}

func (s *SyslUI) GenerateServer() (*Server, error) {
	pkger.Include("/ui/build")

	if s.BasePath == "" {
		s.BasePath = "/"
	}
	syslApps := s.getAllApps()

	webServices := s.buildWebServices(syslApps)
	if len(webServices) == 0 {
		return nil, fmt.Errorf(`no services registered: Make sure ~grpc or ~rest are in the service definitions`) //nolint:lll
	}

	jsonBytes, err := json.Marshal(webServices)
	if err != nil {
		s.Log.Error("Error marshalling services", err)
		return nil, err
	}

	apiDocs := s.buildAPIDocs(syslApps)
	docHandlers := make(map[string]*http.Handler, 10)
	for _, doc := range apiDocs {
		docHandlers[doc.name] = &doc.handler
	}

	return MakeServer(docHandlers, apiDocs, jsonBytes, s.Log, s.Host), nil
}

func (s *SyslUI) buildAPIDocs(apps []*sysl.Application) []*APIDoc {
	var apiDocs []*APIDoc
	for _, a := range apps {
		if serviceType(a) == "" {
			continue
		}

		b := MakeAPIDocBuilder(a, s.Log, s.GRPCUI)
		newDoc, err := b.BuildAPIDoc()
		if err != nil {
			s.Log.Errorf("Error importing %s", a.GetName().GetPart()[0])
			continue
		}
		apiDocs = append(apiDocs, newDoc)
	}
	return apiDocs
}

func (s *SyslUI) buildWebServices(apps []*sysl.Application) []*WebService {
	var webServices []*WebService
	for _, a := range apps {
		// Remove this for webServices??
		if serviceType(a) == "" {
			s.Log.Debugf("Skipping %s due to undefined type", a.GetName().GetPart()[0])
			continue
		}
		newService, err := BuildWebService(a)
		if err != nil {
			s.Log.Errorf("Error importing %s", a.GetName().GetPart()[0])
			continue
		}
		webServices = append(webServices, newService)
	}
	return webServices
}

func (s *SyslUI) getAllApps() []*sysl.Application {
	var syslApps []*sysl.Application
	for _, m := range s.Modules {
		for _, a := range m.GetApps() {
			syslApps = append(syslApps, a)
		}
	}
	return syslApps
}

type Server struct {
	docHandlers map[string]*http.Handler
	spaHandler  spaHandler // Handler for the react app
	serviceJSON []byte     // JSON representation of all the webservices
	apiDocs     []*APIDoc
	router      *mux.Router
	log         *logrus.Logger
	host        string
}

// MakeServer returns a WebServer
//nolint:lll
func MakeServer(docHandlers map[string]*http.Handler, apiDocs []*APIDoc, serviceJSON []byte, log *logrus.Logger, host string) *Server {
	return &Server{
		docHandlers: docHandlers,
		spaHandler:  spaHandler{staticPath: "/ui/build", indexPath: "/ui/build/index.html"},
		serviceJSON: serviceJSON,
		apiDocs:     apiDocs,
		router:      mux.NewRouter(),
		log:         log,
		host:        host,
	}
}

// Serve Runs the command and runs a webserver on catalogURL of a list of endpoints in the sysl file
func (s *Server) Serve() error {
	s.log.Infof("Serving Sysl UI at : %s", s.host)
	log.Fatal(http.ListenAndServe(s.host, s))
	return nil
}

func (s *Server) Setup() error {
	return s.routes()
}

// ServeHTTP calls the http handler from the router
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.log.Debugln(r.URL)
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() error {
	s.router.HandleFunc("/rest/spec/{service}", s.handleRestSpec)
	s.router.PathPrefix("/grpc/{service}").HandlerFunc(s.handleAPIDoc)
	s.router.HandleFunc("/rest/{service}", s.handleAPIDoc)
	s.router.HandleFunc("/data/services.json", s.handleJSONServices)
	s.router.PathPrefix("/").Handler(s.spaHandler)

	return nil
}

func (s *Server) handleRestSpec(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["service"]
	for _, apidoc := range s.apiDocs {
		if apidoc.name == serviceName {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, err := w.Write(apidoc.spec)
			if err != nil {
				panic(err)
			}
			return
		}
	}
}

func (s *Server) handleAPIDoc(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["service"]
	serviceHandler, ok := s.docHandlers[serviceName]
	if !ok {
		s.log.Error("Handler not found")
		return
	}
	(*serviceHandler).ServeHTTP(w, r)
}

// handleJSONServices registers a handler for a JSON list of services
func (s *Server) handleJSONServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write(s.serviceJSON)
	if err != nil {
		panic(err)
	}
}
