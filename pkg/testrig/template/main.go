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

	serviceHandler := {{name}}.NewServiceHandler(genCallbacks, &serviceInterface)

	serviceRouter := {{name}}.NewServiceRouter(genCallbacks, serviceHandler)
	serviceRouter.WireRoutes(ctx, router)

	log.Println("starting {{name}} on :{{port}}")
	return http.ListenAndServe(":{{port}}", router)
}

func main() {
	log.Fatal(LoadServices(context.Background()))
}
`
}
