package main

import (
	"fmt"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/gohugoio/hugo/livereload"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/anz-bank/sysl-catalog/pkg/catalog"
	"github.com/anz-bank/sysl-catalog/pkg/watcher"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"gopkg.in/alecthomas/kingpin.v2"
)

type catalogCmd struct {
	port       string
	server     bool
	outputType string
	outputDir  string
}

func (p *catalogCmd) Name() string       { return "catalog" }
func (p *catalogCmd) MaxSyslModule() int { return 1 }

func (p *catalogCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate Documentation from your sysl definitions")
	cmd.Flag("port", "host and port to serve on").Default(":6900").Short('p').StringVar(&p.port)
	cmd.Flag("server", "start a server on port").Short('s').BoolVar(&p.server)
	cmd.Flag("outputType", "output type").Default("markdown").Short('t').EnumVar(&p.outputType, "html", "markdown")
	cmd.Flag("outputDir", "output directory to generate docs to").Default("/").Short('o').StringVar(&p.outputType)
	return cmd
}

func (p *catalogCmd) Execute(args cmdutils.ExecuteArgs) error {

	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		log.Fatal("Error: Set SYSL_PLANTUML env variable")
	}

	if p.server {
		handler := catalog.NewProject(args.InputFiles[0], "/"+p.outputDir, plantumlService, "html", args.Logger, args.Modules[0]).
			SetServerMode().
			EnableLiveReload()

		go watcher.WatchFile(func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Error:", r)
				}
			}()
			m, err := parse.NewParser().Parse(args.InputFiles[0], args.Filesystem)
			if err != nil {
				panic(err)
			}
			handler.Update(m)
			livereload.ForceRefresh()
			fmt.Println("Done Regenerating")
		}, path.Dir(args.InputFiles[0]))
		fmt.Println("Serving on http://localhost" + p.port)

		livereload.Initialize()
		http.HandleFunc("/livereload.js", livereload.ServeJS)
		http.HandleFunc("/livereload", livereload.Handler)
		http.Handle("/", handler)
		log.Fatal(http.ListenAndServe(p.port, nil))
		select {}
	}
	project := catalog.NewProject(args.InputFiles[0], p.outputDir, plantumlService, p.outputType, args.Logger, args.Modules[0])
	project.SetOutputFs(args.Filesystem).ExecuteTemplateAndDiagrams()
	return nil
}