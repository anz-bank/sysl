package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gohugoio/hugo/livereload"

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
	cmd.Flag("outputDir", "output directory to generate docs to").Default("/").Short('o').StringVar(&p.outputDir)
	return cmd
}

func (p *catalogCmd) Execute(args cmdutils.ExecuteArgs) error {
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		log.Fatal("Error: Set SYSL_PLANTUML env variable")
	}

	if p.server {
		handler := catalog.NewProject(args.DefaultAppName,
			p.outputDir,
			plantumlService,
			"html",
			args.Logger, args.Modules[0]).
			SetServerMode().EnableLiveReload()

		go watcher.WatchFile(func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Error:", r)
				}
			}()
			handler.Update(args.Modules[0])
			livereload.ForceRefresh()
			fmt.Println("Done Regenerating")
		}, folderFromPath(args.ModulePaths)...)
		fmt.Println("Serving on http://localhost" + p.port)

		livereload.Initialize()
		http.HandleFunc("/livereload.js", livereload.ServeJS)
		http.HandleFunc("/livereload", livereload.Handler)
		http.Handle("/", handler)
		log.Fatal(http.ListenAndServe(p.port, nil))
		select {}
	}
	project := catalog.NewProject(
		args.DefaultAppName,
		p.outputDir,
		plantumlService,
		p.outputType,
		args.Logger,
		args.Modules[0])
	project.SetOutputFs(args.Filesystem).ExecuteTemplateAndDiagrams()
	return nil
}

func folderFromPath(files []string) []string {
	var folders []string
	for _, v := range files {
		folders = append(folders, path.Dir(v))
	}
	return folders
}
