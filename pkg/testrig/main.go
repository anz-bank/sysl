package testrig

func GetMainStub() string {
	return `package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	{{.Service.Name}} "{{.Service.URL}}"
	{{.Service.Impl.Name}} "{{.Service.Impl.URL}}"
)

func LoadServices(ctx context.Context) error {
	router := chi.NewRouter()

	serviceInterface := {{.Service.Impl.Name}}.{{.Service.Impl.InterfaceFactory}}()

	genCallbacks := {{.Service.Impl.Name}}.{{.Service.Impl.CallbackFactory}}()

	serviceHandler := {{.Service.Name}}.NewServiceHandler(genCallbacks, &serviceInterface)

	serviceRouter := {{.Service.Name}}.NewServiceRouter(genCallbacks, serviceHandler)
	serviceRouter.WireRoutes(ctx, router)

	log.Println("starting {{.Service.Name}} on :{{.Service.Port}}")
	return http.ListenAndServe(":{{.Service.Port}}", router)
}

func main() {
	log.Fatal(LoadServices(context.Background()))
}
`
}
