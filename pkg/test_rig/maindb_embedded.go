package test_rig

import (
	_ "github.com/lib/pq"
)

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

	{{name}} "{{url}}"
	{{impl.name}} "{{impl.url}}"
)

func initDb() (*sql.DB, error) {
	psqlInfo := "host={{name}}_db port=5432 user=someuser password=somepassword dbname=somedb sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

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

	serviceInterface := {{impl.name}}.{{impl.interface_factory}}()

	genCallbacks := {{impl.name}}.{{impl.callback_factory}}(db)

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
	loadServices(context.Background())
}
`
}
