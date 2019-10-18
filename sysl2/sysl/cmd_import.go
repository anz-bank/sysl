package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/anz-bank/sysl/sysl2/sysl/importer"
	"gopkg.in/alecthomas/kingpin.v2"
)

type importCmd struct {
	importer.OutputData
	filename string
	outfile  string
	mode     string
}

func (p *importCmd) Name() string            { return "import" }
func (p *importCmd) RequireSyslModule() bool { return false }

func (p *importCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Import foreign type to sysl. Supported types: Swagger and xsd")
	cmd.Flag("input", "input filename").Short('i').Required().StringVar(&p.filename)
	cmd.Flag("app-name",
		"name of the sysl app to define in sysl model.").Short('a').Required().StringVar(&p.AppName)
	cmd.Flag("package",
		"name of the sysl package to define in sysl model.").Short('p').Required().StringVar(&p.Package)
	cmd.Flag("output", "output filename").Default("output.sysl").Short('o').StringVar(&p.outfile)

	opts := []string{"auto", "swagger", "xsd"}
	cmd.Flag("format", fmt.Sprintf("format of the input filename, options: [%s]", strings.Join(opts, ", "))).
		Short('f').
		Default("auto").
		HintOptions(opts...).
		EnumVar(&p.mode, opts...)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

const (
	modeSwagger = "swagger"
	modeXSD     = "xsd"
)

func (p *importCmd) Execute(args ExecuteArgs) error {
	data, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return err
	}

	if p.mode == "auto" {
		p.mode = guessFileType(p.filename)
	}
	var imp importer.Func
	switch p.mode {
	case modeSwagger:
		args.Logger.Infof("Using swagger importer\n")
		imp = importer.LoadSwaggerText
	case modeXSD:
		args.Logger.Infof("Using xsd importer\n")
		imp = importer.LoadXSDText
	default:
		args.Logger.Fatalf("Unsupported input format: %s\n", p.mode)
	}
	output, err := imp(p.OutputData, string(data), args.Logger)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p.outfile, []byte(output), 0644)
}

func guessFileType(filename string) string {
	parts := strings.Split(filename, ".")
	switch ext := parts[len(parts)-1]; ext {
	case "xml", "xsd":
		return modeXSD
	case "yaml", "yml", "json":
		return modeSwagger
	default:
		return ext
	}
}
