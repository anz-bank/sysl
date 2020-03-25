package main

import (
	"os"
	"path/filepath"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/anz-bank/sysl/pkg/transforms"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

type templateCmd struct {
	rootTemplate string
	template     string
	appName      []string
	start        string
	outDir       string
}

func (p *templateCmd) Name() string       { return "template" }
func (p *templateCmd) MaxSyslModule() int { return 1 }

func (p *templateCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Apply a model to a template for custom text output").Alias("tmpl")

	cmd.Flag("root-template",
		"sysl root directory for input template file (default: .)").
		Default(".").StringVar(&p.rootTemplate)
	cmd.Flag("template", "path to template file from the root transform directory").Required().StringVar(&p.template)
	cmd.Flag("app-name",
		"name of the sysl app defined in the sysl model."+
			" if there are multiple Apps defined in the sysl model,"+
			" code will be generated only for the given app").
		Short('a').Default("").StringsVar(&p.appName)
	cmd.Flag("start", "start rule for the template").Required().StringVar(&p.start)
	cmd.Flag("outdir", "output directory").Short('o').Default(".").StringVar(&p.outDir)
	EnsureFlagsNonEmpty(cmd, "app-name")
	return cmd
}

func (p *templateCmd) Execute(args cmdutils.ExecuteArgs) error {
	tmplFs := syslutil.NewChrootFs(args.Filesystem, p.rootTemplate)
	tfmParser := parse.NewParserWithParserType(args.ParserType)
	tx, transformAppName, err := parse.LoadAndGetDefaultApp(p.template, tmplFs, tfmParser)
	if err != nil {
		return err
	}

	t, err := transforms.NewWorker(tx, transformAppName, p.start)
	if err != nil {
		return err
	}

	output := t.Apply(args.Modules[0], p.appName...)

	for filename, data := range output {
		if _, err := os.Stat(p.outDir); os.IsNotExist(err) {
			if err = os.Mkdir(p.outDir, 0755); err != nil {
				return errors.Wrap(err, "Error creating output folder; check permission")
			}
		}
		if err := afero.WriteFile(tmplFs, filepath.Join(p.outDir, filename), []byte(data.GetS()), 0644); err != nil {
			return err
		}
	}
	return nil
}
