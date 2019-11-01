package main

import (
	"fmt"

	sysl "github.com/anz-bank/sysl/src/proto"

	"github.com/anz-bank/sysl/sysl2/sysl/parse"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

type templateCmd struct {
	rootTemplate string
	template     string
	appName      string
	start        string
	outDir       string
}

func (p *templateCmd) Name() string            { return "template" }
func (p *templateCmd) RequireSyslModule() bool { return true }

func (p *templateCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Apply a model to a template for custom text output").Alias("tmpl")

	cmd.Flag("root-template",
		"sysl root directory for input template file (default: .)").
		Default(".").StringVar(&p.rootTemplate)
	cmd.Flag("template", "path to template file from the root transform directory").Required().StringVar(&p.template)
	cmd.Flag("app-name",
		"name of the sysl app defined in sysl model."+
			" if there are multiple apps defined in sysl model,"+
			" code will be generated only for the given app").
		Short('a').Default("").StringVar(&p.appName)
	cmd.Flag("start", "start rule for the template").Default(".").StringVar(&p.start)
	cmd.Flag("outdir", "output directory").Default(".").StringVar(&p.outDir)
	EnsureFlagsNonEmpty(cmd, "app-name")
	return cmd
}

type modData struct {
	appName string
	mod     *sysl.Module
}

func (p *templateCmd) Execute(args ExecuteArgs) error {

	tmplFs := syslutil.NewChrootFs(args.Filesystem, p.rootTemplate)
	tfmParser := parse.NewParser()
	tx, transformAppName, err := parse.LoadAndGetDefaultApp(p.template, tmplFs, tfmParser)
	if err != nil {
		return err
	}

	tmpl := modData{
		appName: transformAppName,
		mod:     tx,
	}

	if p.appName == "" {
		for k, _ := range args.Module.GetApps() {
			p.run(tmpl, modData{appName: k, mod: args.Module})
		}
	} else {
		p.run(tmpl, modData{appName: p.appName, mod: args.Module})
	}

	return nil
}

func (p *templateCmd) run(template, model modData) error {

	res, err := applyTransformToModel(model.appName, template.appName, p.start, model.mod, template.mod)
	if err != nil {
		return err
	}

	fmt.Printf(res.GetS())

	return err
}
