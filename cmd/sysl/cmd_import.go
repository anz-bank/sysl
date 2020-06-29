package main

import (
	"fmt"
	"io/ioutil"
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
	outfile  string
}

func (p *importCmd) Name() string       { return "import" }
func (p *importCmd) MaxSyslModule() int { return 0 }

func (p *importCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	optsText := buildOptionsText(importer.Formats)
	cmd := app.Command(p.Name(), "Import foreign type to sysl. Supported types: ["+optsText+"]")
	cmd.Flag("input", "input filename").Short('i').Required().StringVar(&p.filename)
	cmd.Flag("app-name",
		"name of the sysl app to define in sysl model.").Required().Short('a').StringVar(&p.AppName)
	cmd.Flag("package",
		"name of the sysl package to define in sysl model.").Short('p').StringVar(&p.Package)
	cmd.Flag("output", "output filename").Default("output.sysl").Short('o').StringVar(&p.outfile)
	cmd.Flag("format", fmt.Sprintf("format of the input filename, options: [%s]."+
		"NOTE: This flag is deprecated. Please remove this from your scripts, as formats are now autodetected",
		optsText)).String()

	return cmd
}

func (p *importCmd) Execute(args cmdutils.ExecuteArgs) error {
	data, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return err
	}

	var imp importer.Importer
	inputFilePath, err := filepath.Abs(p.filename)
	if err != nil {
		return err
	}
	imp, err = importer.Factory(inputFilePath, data, logrus.New())
	if err != nil {
		return err
	}
	imp.WithAppName(p.AppName).WithPackage(p.Package)
	output, err := imp.Load(string(data))
	if err != nil {
		return err
	}

	return afero.WriteFile(args.Filesystem, p.outfile, []byte(output), os.ModePerm)
}

func buildOptionsText(opts []importer.Format) string {
	var optionsText []string
	for _, format := range opts {
		optionsText = append(optionsText, format.Name)
	}
	sort.Strings(optionsText)
	return strings.Join(optionsText, ", ")
}
