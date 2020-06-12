package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"github.com/anz-bank/sysl/pkg/syslwrapper"

	"github.com/anz-bank/sysl/pkg/exporter"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/alecthomas/kingpin.v2"
)

type exportCmd struct {
	appName string
	out     string
	mode    string
	format  string
}

const (
	swaggerMode  = "swagger"
	openapi2Mode = "openapi2"
	openapi3Mode = "openapi3"
	jsonMode     = "json"
	yamlMode     = "yaml"
)

func (p *exportCmd) Name() string       { return "export" }
func (p *exportCmd) MaxSyslModule() int { return 1 }

func (p *exportCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Export sysl to external types. Supported types: Swagger,openapi2,openapi3")
	cmd.Flag("app-name", "name of the sysl App defined in the sysl model."+
		" if there are multiple Apps defined in the sysl model,"+
		" swagger will be generated only for the given app").Short('a').StringVar(&p.appName)
	cmd.Flag(
		"format",
		"format of export, supported options; (swagger | openapi2 | openapi3)",
	).Default("swagger").Short('f').StringVar(&p.mode)
	cmd.Flag("output", "output filepath.format(yaml | json) (default: %(appname).yaml)").Default(
		"%(appname).yaml").Short('o').StringVar(&p.out)
	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *exportCmd) writeSwaggerForApp(
	fs afero.Fs,
	filename string,
	syslApp *sysl.Application,
	logger *logrus.Logger,
) error {
	var output []byte
	switch p.mode {
	case openapi2Mode, swaggerMode:
		swaggerExporter := exporter.MakeSwaggerExporter(syslApp, logger)
		err := swaggerExporter.GenerateSwagger()
		if err != nil {
			logger.Warnf("Error generating Swagger/Openapi2 for the application %s", err)
			return err
		}
		output, err = swaggerExporter.SerializeOutput(p.format)
		if err != nil {
			return err
		}
	case openapi3Mode:
		mod := &sysl.Module{
			Apps: map[string]*sysl.Application{
				syslApp.Name.Part[0]: syslApp,
			},
		}
		mapper := syslwrapper.MakeAppMapper(mod)
		mapper.IndexTypes()
		simpleApps, err := mapper.Map()
		if err != nil {
			return err
		}
		openapi3Exporter := exporter.MakeOpenAPI3Exporter(simpleApps, logger)
		err = openapi3Exporter.Export()
		if err != nil {
			logger.Warnf("Error generating Openapi3 for the application %s", err)
			return err
		}
		output, err = openapi3Exporter.SerializeOutput(syslApp.Name.Part[0], p.format)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported export format")
	}

	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return err
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.Debugf("Target folder does not exist %s", err)
		if err = os.Mkdir(dir, 0755); err != nil {
			logger.Errorf("Error creating target folder; check permission")
			return err
		}
	}
	err = afero.WriteFile(fs, filename, output, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *exportCmd) Execute(args cmdutils.ExecuteArgs) error {
	err := p.determineOperationMode(p.out)
	if err != nil {
		return err
	}

	var writeCount int
	for appName, syslApp := range args.Modules[0].GetApps() {
		if appName == p.appName || p.appName == "" {
			outputFileName := cmdutils.MakeFormatParser(p.out).LabelApp(appName, "", syslApp.GetAttrs())
			if err := args.Filesystem.MkdirAll(filepath.Dir(outputFileName), os.ModePerm); err != nil {
				return err
			}
			if p.appName == "" && outputFileName == p.out {
				ext := filepath.Ext(outputFileName)
				// convert out.yaml something like out.Fooapp.yaml
				outputFileName = fmt.Sprintf("%s.%s%s", strings.TrimSuffix(outputFileName, ext), appName, ext)
			}
			args.Logger.Infof("Exporting app `%s` -> %s\n", appName, outputFileName)
			err := p.writeSwaggerForApp(args.Filesystem, outputFileName, syslApp, args.Logger)
			if err != nil {
				return err
			}
			writeCount++
		}
	}
	if writeCount == 0 {
		return fmt.Errorf("app not found in the Sysl file")
	}
	return nil
}

func (p *exportCmd) determineOperationMode(filename string) error {
	fileExtn := strings.TrimLeft(filepath.Ext(filepath.Base(filename)), ".")
	switch fileExtn {
	case jsonMode:
		p.format = jsonMode
	case yamlMode:
		p.format = yamlMode
	default:
		return fmt.Errorf("invalid output file format %s", fileExtn)
	}
	return nil
}
