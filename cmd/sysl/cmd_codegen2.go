package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/anz-bank/sysl/pkg/sysl"

	"github.com/anz-bank/sysl/pkg/cmdutils"

	"github.com/anz-bank/sysl/pkg/syslutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

type codegen2Cmd struct {
	cmdutils.CmdContextParamCodegen
	outDir         string
	appName        string
	validateOnly   bool
	enableDebugger bool

	transforms map[string]transformData
}

type transformData struct {
	grammar      string
	grammarStart string
	transforms   []string
}

func (t transformData) exec(codegenParams cmdutils.CmdContextParamCodegen,
	model *sysl.Module, modelAppName string,
	fs afero.Fs, outdir string, logger *logrus.Logger) error {
	outfs := getOutFs(fs, outdir)
	for _, xform := range t.transforms {
		codegenParams.Transform = xform
		output, err := GenerateCode(&codegenParams, model, modelAppName, fs, logger)
		if err != nil {
			return err
		}
		if outputToFiles(output, outfs) != nil {
			return err
		}
	}
	return nil
}

func getOutFs(fs afero.Fs, outdir string) afero.Fs {
	_, err := fs.Stat(outdir)
	if os.IsNotExist(err) {
		fs.MkdirAll(outdir, os.ModePerm)
	}
	return syslutil.NewChrootFs(fs, outdir)
}

func (p *codegen2Cmd) Name() string       { return "syslgen" }
func (p *codegen2Cmd) MaxSyslModule() int { return 1 }

func (p *codegen2Cmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate code via a sysl control file").Alias("gen2")
	cmd.Flag("root-transform",
		"sysl root directory for input transform file (default: .)").
		Default(".").StringVar(&p.RootTransform)
	cmd.Flag("outdir", "output directory").Default(".").StringVar(&p.outDir)
	cmd.Flag("debugger", "Enable the evaluation debugger on error").Default("false").BoolVar(&p.enableDebugger)
	return cmd
}

func (p *codegen2Cmd) Execute(args cmdutils.ExecuteArgs) error {
	p.transforms = map[string]transformData{}

	// init transform data
	for _, app := range args.Modules[0].GetApps() {
		if syslutil.HasPattern(app.GetAttrs(), "codegen") {
			for epName, ep := range app.GetEndpoints() {
				xform := transformData{}
				if t, ok := ep.GetAttrs()["targets"]; ok && t != nil {
					for _, val := range t.GetA().Elt {
						xform.transforms = append(xform.transforms, val.GetS())
					}
				}
				xform.grammar = app.GetAttrs()["grammar"].GetS()
				xform.grammarStart = app.GetAttrs()["grammar_start"].GetS()
				p.transforms[epName] = xform
			}
		}
	}

	// now find the apps to execute
	for appName, app := range args.Modules[0].GetApps() {
		if cfg, err := getCodegenConfig(app.GetAttrs()); err == nil {
			if x, ok := p.transforms[cfg.transform]; ok {
				params := cmdutils.CmdContextParamCodegen{
					RootTransform:    p.RootTransform,
					Grammar:          x.grammar,
					Start:            x.grammarStart,
					DepPath:          "",
					BasePath:         cfg.basepath,
					DisableValidator: false,
				}
				x.exec(params, args.Modules[0], appName, args.Filesystem, filepath.Join(p.outDir, cfg.outdir), args.Logger)
			}
		}
	}
	return nil
}

type codegencfg struct {
	transform, basepath, outdir string
}

func getCodegenConfig(attrs map[string]*sysl.Attribute) (codegencfg, error) {
	cfg := codegencfg{}
	if a, ok := attrs["codegen_transform"]; ok && a != nil {
		cfg.transform = a.GetS()
	} else {
		return cfg, fmt.Errorf("missing @codegen_transform")
	}
	if a, ok := attrs["codegen_dirname"]; ok && a != nil {
		cfg.outdir = a.GetS()
	}

	if a, ok := attrs["codegen_basepath"]; ok && a != nil {
		cfg.basepath = a.GetS()
	}

	return cfg, nil
}
