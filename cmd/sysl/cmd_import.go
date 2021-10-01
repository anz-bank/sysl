package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/pkg/importer"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

type importCmd struct {
	importer.OutputData
	filename string
	format   string
	outFile  string
}

func (p *importCmd) Name() string       { return "import" }
func (p *importCmd) MaxSyslModule() int { return 0 }

func (p *importCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	optsText := buildOptionsText(importer.Formats)
	cmd := app.Command(p.Name(), "Import foreign type to Sysl. Supported types: ["+optsText+"]")
	cmd.Flag("input", "path of file to import").Short('i').Required().StringVar(&p.filename)
	cmd.Flag("app-name",
		"name of the Sysl app to define in Sysl model.").Required().Short('a').StringVar(&p.AppName)
	cmd.Flag("package",
		"name of the Sysl package to define in Sysl model.").Short('p').StringVar(&p.Package)
	cmd.Flag("output", "path of file to write the imported sysl, writes to stdout when not specified").
		Short('o').StringVar(&p.outFile)
	cmd.Flag("format", fmt.Sprintf("format of the input filename, options: [%s]. "+
		"Formats are autodetected, but this can force the use of a particular importer.", optsText)).Short('f').
		StringVar(&p.format)
	cmd.Flag("import-paths", "comma separated list of paths used to resolve imports in "+
		"the input file. Currently only used for protobuf input.").Short('I').
		StringVar(&p.ImportPaths)
	return cmd
}

func (p *importCmd) Execute(args cmdutils.ExecuteArgs) error {
	isDir, err := afero.IsDir(args.Filesystem, p.filename)
	if err != nil {
		return err
	}
	var content []byte
	if !isDir {
		content, err = afero.ReadFile(args.Filesystem, p.filename)
		if err != nil {
			return err
		}
	}
	var imp importer.Importer
	inputFilePath, err := filepath.Abs(p.filename)
	if err != nil {
		return err
	}
	imp, err = importer.Factory(inputFilePath, isDir, p.format, content, logrus.New())
	if err != nil {
		return err
	}
	imp.WithAppName(p.AppName).WithPackage(p.Package).WithImports(p.ImportPaths)

	var output string
	// TODO: Abstract this logic.
	switch imp.(type) {
	case *importer.ArraiImporter:
		output, err = imp.LoadFile(inputFilePath)
	case *importer.ProtobufImporter:
		// TODO: support `Load` for proto importer
		output, err = imp.LoadFile(p.filename)
	default:
		output, err = imp.Load(string(content))
	}
	if err != nil {
		return err
	}

	if p.outFile != "" {
		return afero.WriteFile(args.Filesystem, p.outFile, []byte(output), os.ModePerm)
	}
	_, err = fmt.Println(output)
	return err
}

func buildOptionsText(opts []importer.Format) string {
	var optionsText []string
	for _, format := range opts {
		optionsText = append(optionsText, format.Name)
	}
	sort.Strings(optionsText)
	return strings.Join(optionsText, ", ")
}
