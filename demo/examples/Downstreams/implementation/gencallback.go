// implementation/gencallback.go

// gencallback.go contains all the manual config code that is used to implement the generated sysl
package implementation

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.service.anz/sysl/server-lib/common"
	"github.service.anz/sysl/server-lib/validator"
	"github.service.anz/sysl/syslbyexample/_examples/Downstreams/gen/simple"
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
	return context.WithTimeout(ctx, 1*time.Second)
}

func LoadServices(ctx context.Context) error {
	router := chi.NewRouter()

	// simpleServiceInterface is the struct which is composed of our functions we wrote in `methods.go`
	// Struct embedding is used for the Service interface (yes, not interfaces); GetStuffList and GetFoobarList are the functions we wrote
	simpleServiceInterface := simple.ServiceInterface{
		GetStuffList:  GetStuffList,
		GetFoobarList: GetFoobarList,
	}

	genCallbacks := Callback{
		UpstreamTimeout:   60 * time.Millisecond,
		DownstreamTimeout: 60 * time.Millisecond,
		RouterBasePath:    "/",
		UpstreamConfig:    nil,
	}

	// Service Handler
	// use the Example service that implements behaviour of the downstream service
	serviceHandler := simple.NewServiceHandler(genCallbacks, &simpleServiceInterface, &ExampleService{})

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
