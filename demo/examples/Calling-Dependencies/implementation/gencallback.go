// implementation/gencallback.go

// gencallback.go contains all the manual config code that is used to implement the generated sysl
package implementation

import (
	"context"
	"flag"
	"github.com/go-chi/chi"
	"github.service.anz/sysl/server-lib/common"
	"github.service.anz/sysl/server-lib/validator"
	"github.service.anz/sysl/syslbyexample/_examples/Calling-Dependencies/gen/mydependency"
	"github.service.anz/sysl/syslbyexample/_examples/Calling-Dependencies/gen/simple"
	"log"
	"net/http"
	"time"
)

type Callback struct {
	UpstreamTimeout   time.Duration
	DownstreamTimeout time.Duration
	RouterBasePath    string
	UpstreamConfig    validator.Validator
}

type Config struct{}

func (c Config) Validate() error {
	return nil
}

func (g Callback) AddMiddleware(ctx context.Context, r chi.Router) {
}

func (g Callback) BasePath() string {
	return "/"
}

func (g Callback) Config() validator.Validator {
	return Config{}
}

func (g Callback) HandleError(ctx context.Context, w http.ResponseWriter, kind common.Kind, message string, cause error) {
	se := common.CreateError(ctx, kind, message, cause)

	httpError := common.HandleError(ctx, se)

	httpError.WriteError(ctx, w)
}

func (g Callback) DownstreamTimeoutContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 120*time.Second)
}

func LoadServices(ctx context.Context) error {
	router := chi.NewRouter()

	// simpleServiceInterface is the struct which is composed of our functions we wrote in `methods.go`

	// Struct embedding is used for the Service interface (yes, not interfaces)
	simpleServiceInterface := simple.ServiceInterface{
		GetFoobarList: GetFoobarList,
	}

	genCallbacks := Callback{
		UpstreamTimeout:   120 * time.Second,
		DownstreamTimeout: 120 * time.Second,
		RouterBasePath:    "/",
		UpstreamConfig:    nil,
	}

	// Service Handler
	// use the Example service that implements behaviour of the downstream service
	// Here we specify that the base path (serviceURL) of the url we're hitting; https://jsonplaceholder.typicode.com
	// Sometimes this will have endpoints attached https://example.com/v2 but should never have a trailing slash
	// Note that this ServiceURL is http, as our service is served over http
	serviceHandler := simple.NewServiceHandler(genCallbacks, &simpleServiceInterface, mydependency.NewClient(http.DefaultClient, "http://jsonplaceholder.typicode.com"))

	// Service Router
	serviceRouter := simple.NewServiceRouter(genCallbacks, serviceHandler)
	serviceRouter.WireRoutes(ctx, router)

	var serverAddress string
	flag.StringVar(&serverAddress, "p", ":8080", "Specify server address")
	flag.Parse()

	log.Println("Starting Server on " + serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
	return nil
}
