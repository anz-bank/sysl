package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/buger/goterm"
	"github.com/gohugoio/hugo/livereload"

	"github.com/anz-bank/sysl-catalog/pkg/catalog"
	"github.com/anz-bank/sysl-catalog/pkg/watcher"
	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/loader"
	"github.com/anz-bank/sysl/pkg/sysl"
	"gopkg.in/alecthomas/kingpin.v2"
)

type catalogCmd struct {
	port              string
	server            bool
	outputType        string
	outputDir         string
	templates         string
	noCSS             bool
	disableLiveReload bool
	noImages          bool
	enableMermaid     bool
}

func (p *catalogCmd) Name() string       { return "catalog" }
func (p *catalogCmd) MaxSyslModule() int { return 1 }

func (p *catalogCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate Documentation from your sysl definitions")
	cmd.Flag("port", "host and port to serve on").Default(":6900").Short('p').StringVar(&p.port)
	cmd.Flag("server", "start a server on port").Short('s').BoolVar(&p.server)
	cmd.Flag("outputType", "output type").Default("markdown").Short('t').EnumVar(&p.outputType, "html", "markdown")
	cmd.Flag("outputDir", "output directory to generate docs to").Default("/").Short('o').StringVar(&p.outputDir)
	cmd.Flag("templates", "templates, separated by commas ").StringVar(&p.templates)
	cmd.Flag("noCSS", "disable adding css to served html").BoolVar(&p.noCSS)
	cmd.Flag("disableLiveReload", "diable live reload").Default("false").BoolVar(&p.disableLiveReload)
	cmd.Flag("noImages", "don't create images").Default("false").BoolVar(&p.noImages)
	cmd.Flag("mermaid", "use mermaid diagrams where possible").Default("false").BoolVar(&p.enableMermaid)
	return cmd
}

func (p *catalogCmd) Execute(args cmdutils.ExecuteArgs) error {
	plantumlService := os.Getenv("SYSL_PLANTUML")
	if plantumlService == "" {
		log.Fatal("Error: Set SYSL_PLANTUML env variable")
	}

	if !p.server {
		catalog.NewProject(
			args.ModulePaths[0],
			plantumlService,
			p.outputType,
			args.Logger,
			args.Modules[0],
			args.Filesystem,
			p.outputDir,
			false).
			WithTemplateFs(args.Filesystem, strings.ReplaceAll(p.templates, ",", "")).
			SetOptions(p.noCSS, p.noImages, "").
			Run()
		return nil
	}

	handler := catalog.NewProject(
		args.ModulePaths[0],
		plantumlService,
		"html",
		args.Logger,
		nil,
		nil,
		p.outputDir,
		false).
		WithTemplateFs(args.Filesystem, strings.ReplaceAll(p.templates, ",", "")).
		ServerSettings(p.noCSS, !p.disableLiveReload, true)
	goterm.Clear()
	PrintToPosition(1, "Serving on http://localhost"+p.port)
	go watcher.WatchFile(func(i interface{}) {
		PrintToPosition(3, "Regenerating")
		m, err := func() (m *sysl.Module, err error) {
			defer func() {
				if r := recover(); r != nil {
					m = nil
					fmt.Println("Error:", r)
				}
			}()
			m, _, err = loader.LoadSyslModule("", args.ModulePaths[0], args.Filesystem, args.Logger)
			if err != nil {
				return nil, err
			}

			return
		}()
		if err != nil {
			PrintToPosition(4, err)
		}
		handler.Update(m)
		livereload.ForceRefresh()
		PrintToPosition(2, i)
		PrintToPosition(4, goterm.RESET_LINE)
		PrintToPosition(3, goterm.RESET_LINE)
		PrintToPosition(3, "Done Regenerating")
	}, folderFromPath(args.ModulePaths)...)

	livereload.Initialize()
	http.HandleFunc("/livereload.js", livereload.ServeJS)
	http.HandleFunc("/livereload", livereload.Handler)
	http.Handle("/", handler)
	log.Fatal(http.ListenAndServe(p.port, nil))
	select {}
}

func PrintToPosition(y int, i interface{}) {
	goterm.MoveCursor(1, y)
	goterm.Print(i)
	goterm.Flush()
}

func folderFromPath(files []string) []string {
	var folders []string
	for _, v := range files {
		folders = append(folders, path.Dir(v))
	}
	return folders
}
