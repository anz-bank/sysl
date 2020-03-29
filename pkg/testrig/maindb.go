package testrig

func GetMainDbStub() string {
	return `package main

import (
	"context"
	"log"
	"net/http"
	"database/sql"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	{{.Service.Name}} "{{.Service.URL}}"
	{{.Service.Impl.Name}} "{{.Service.Impl.URL}}"
)

func initDb() (*sql.DB, error) {
	psqlInfo := "host={{.Service.Name}}_db port=5432 user=someuser password=somepassword dbname=somedb sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second)

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}

func loadServices(ctx context.Context) error {
	router := chi.NewRouter()

	db, err := initDb()
	if err != nil {
		return err
	}
	defer db.Close()

	serviceInterface := {{.Service.Impl.Name}}.{{.Service.Impl.InterfaceFactory}}()

	genCallbacks := {{.Service.Impl.Name}}.{{.Service.Impl.CallbackFactory}}(db)

	serviceHandler := {{.Service.Name}}.NewServiceHandler(genCallbacks, &serviceInterface)

	serviceRouter := {{.Service.Name}}.NewServiceRouter(genCallbacks, serviceHandler)
	serviceRouter.WireRoutes(ctx, router)

	log.Println("starting {{.Service.Name}} on :{{.Service.Port}}")
	return http.ListenAndServe(":{{.Service.Port}}", router)
}

func main() {
	log.Fatal(loadServices(context.Background()))
}
`
}
