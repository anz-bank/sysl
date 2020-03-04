package template

func GetMainStub() string {
	return `package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	{{name}} "{{url}}"
	{{impl.name}} "{{impl.url}}"
)

func LoadServices(ctx context.Context) error {
	router := chi.NewRouter()

	serviceInterface := {{impl.name}}.{{impl.interface_factory}}()

	genCallbacks := {{impl.name}}.{{impl.callback_factory}}()

	// serviceHandler := simple.NewServiceHandler(genCallbacks, &serviceInterface, mydependency.NewClient(http.DefaultClient, "http://jsonplaceholder.typicode.com"))
	serviceHandler := {{name}}.NewServiceHandler(genCallbacks, &serviceInterface)

	// Service Router
	serviceRouter := {{name}}.NewServiceRouter(genCallbacks, serviceHandler)
	serviceRouter.WireRoutes(ctx, router)

	log.Println("starting {{name}} on :{{port}}")
	log.Fatal(http.ListenAndServe(":{{port}}", router))
	return nil
}

func main() {
	LoadServices(context.Background())
}
`
}
