// server.go contains all the manual config code that is used to implement the generated sysl
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anz-bank/sysl/demo/examples/Code-Generation/gen/jsonplaceholder"
	"github.com/anz-bank/sysl/demo/examples/Code-Generation/gen/simple"
	"github.com/anz-bank/sysl/demo/examples/Code-Generation/pkg/defaultcallback"
	"github.com/go-chi/chi"
)

func LoadServices(ctx context.Context) error {
	router := chi.NewRouter()

	// simpleServiceInterface is the struct which is composed of our functions we wrote in `methods.go`
	// Struct embedding is used for the Service interface (yes, not interfaces)
	simpleServiceInterface := simple.ServiceInterface{
		GetFoobarList: GetFoobarList,
		Get:           GetHandler,
	}

	// Default callback behaviour
	genCallbacks := defaultcallback.DefaultCallback()

	serviceHandler := simple.NewServiceHandler(
		genCallbacks,
		&simpleServiceInterface,
		jsonplaceholder.NewClient(http.DefaultClient, "http://jsonplaceholder.typicode.com"))

	// Service Router
	serviceRouter := simple.NewServiceRouter(genCallbacks, serviceHandler)
	serviceRouter.WireRoutes(ctx, router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddress := ":" + port

	log.Println("Starting Server on " + serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, router))
	return nil
}

// GetHandler refers to the endpoint in our sysl file
func GetHandler(ctx context.Context, req *simple.GetRequest, client simple.GetClient) (*simple.Welcome, error) {
	welcome := simple.Welcome{
		Content: "Hello World!",
	}
	return &welcome, nil
}

// GetFoobarList refers to the endpoint in our sysl file
func GetFoobarList(ctx context.Context, req *simple.GetFoobarListRequest, client simple.GetFoobarListClient) (*jsonplaceholder.TodosResponse, error) {
	// Here we can make a request on the client object which was generated from the call to "myDownstream" in the sysl model
	// We will get the id equal to one, which was generated from out {id} from /todos/{id<:int}
	ans, err := client.GetTodos(ctx, &jsonplaceholder.GetTodosRequest{ID: 1})
	fmt.Println(ans)
	return ans, err
}
