// Package testrig provides tools for generating standalone test environment for sysl generated services
package testrig

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/ghodss/yaml"
	"github.com/spf13/afero"
)

// GenerateRig creates full set of files required to start up a test rig
// for every service, it creates main.go (optionally with DB connection) and a dockerfile to containerize it
// if the service contains DB tables, it creates sidecar postgres container and creates a schema there
// finally, it creates docker-compose.yml at the point of all (likely your project's root)
// that joins all containerized services
func GenerateRig(fs afero.Fs, templateFile string, outputDir string, modules []*sysl.Module) error {
	serviceMap, err := readUserData(fs, templateFile)
	if err != nil {
		return err
	}

	applications := make(map[string]*sysl.Application)
	for _, module := range modules {
		for appName, app := range module.GetApps() {
			applications[appName] = app
		}
	}

	composeServices := make(map[string]interface{})

	for serviceName, serviceTmpl := range serviceMap {
		if err := fs.MkdirAll(filepath.Join(outputDir, serviceName), os.ModePerm); err != nil {
			return err
		}

		if err := generateMain(fs, serviceName, serviceTmpl, outputDir, appNeedsDB(applications[serviceName])); err != nil {
			return err
		}

		if err := generateDockerfile(fs, serviceName, serviceTmpl, outputDir); err != nil {
			return err
		}

		var dependencies []string
		if appNeedsDB(applications[serviceName]) {
			dbServiceName := serviceName + "_db"
			composeServices[dbServiceName] = generateDBService(serviceName)
			dependencies = []string{dbServiceName}
		}
		composeServices[serviceName] = generateAppService(&generateAppServiceInput{
			ServiceName: serviceName,
			Vars:        serviceTmpl,
			OutputDir:   outputDir,
			DependsOn:   dependencies,
		})
	}

	return generateCompose(fs, composeServices)
}

func appNeedsDB(app *sysl.Application) bool {
	if app == nil {
		return false
	}
	patterns := syslutil.MakeStrSetFromAttr("patterns", app.Attrs)
	return patterns.Contains("DB")
}

func readUserData(fs afero.Fs, templateFileName string) (ServiceVarsMap, error) {
	templateFile, err := fs.Open(templateFileName)
	if err != nil {
		return nil, err
	}
	defer templateFile.Close()

	byteValue, err := io.ReadAll(templateFile)
	if err != nil {
		return nil, err
	}
	var vars ServiceVarsMap
	if err := json.Unmarshal(byteValue, &vars); err != nil {
		return nil, err
	}
	return vars, nil
}

type generateAppServiceInput struct {
	ServiceName string
	Vars        ServiceVars
	OutputDir   string
	DependsOn   []string
}

func generateAppService(input *generateAppServiceInput) map[string]interface{} {
	return map[string]interface{}{
		"build": map[string]interface{}{
			"context":    ".",
			"dockerfile": filepath.Join(input.OutputDir, input.ServiceName, "Dockerfile"),
		},
		"ports":      []string{fmt.Sprintf("%v:%v", input.Vars.Port, input.Vars.Port)},
		"depends_on": input.DependsOn,
	}
}

func generateDBService(serviceName string) map[string]interface{} {
	return map[string]interface{}{
		"image": "postgres:latest",
		"ports": []string{"5432:5432"},
		"volumes": []string{
			fmt.Sprintf("../gen/%v/%v.sql:/docker-entrypoint-initdb.d/%v.sql", serviceName, serviceName, serviceName),
		},
		"environment": map[string]string{
			"POSTGRES_USER":     "someuser",
			"POSTGRES_PASSWORD": "somepassword",
			"POSTGRES_DB":       "somedb",
		},
	}
}

func generateCompose(fs afero.Fs, composeServices map[string]interface{}) error {
	output, err := fs.Create("docker-compose.yml")
	if err != nil {
		return err
	}
	defer output.Close()

	data := map[string]interface{}{
		"version":  "3.3",
		"services": composeServices,
	}

	dataRaw, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := output.Write(dataRaw); err != nil {
		return err
	}
	return output.Sync()
}

func processTemplate(file afero.File, resource string, data ServiceVars, outputDir string) error {
	tmpl, err := template.New("servicesData").Parse(resource)
	if err != nil {
		return err
	}
	userData := struct {
		OutputDir string
		Service   ServiceVars
	}{
		OutputDir: outputDir,
		Service:   data,
	}
	if err := tmpl.Execute(file, userData); err != nil {
		return err
	}
	return file.Sync()
}

func generateMain(fs afero.Fs, serviceName string, serviceVars ServiceVars, outputDir string, needDB bool) error {
	output, err := fs.Create(filepath.Join(outputDir, serviceName, "main.go"))
	if err != nil {
		return err
	}
	defer output.Close()

	mainTemplate := GetMainStub()
	if needDB {
		mainTemplate = GetMainDBStub()
	}

	return processTemplate(output, mainTemplate, serviceVars, outputDir)
}

func generateDockerfile(fs afero.Fs, serviceName string, serviceVars ServiceVars, outputDir string) error {
	output, err := fs.Create(filepath.Join(outputDir, serviceName, "Dockerfile"))
	if err != nil {
		return err
	}
	defer output.Close()

	dockerTemplate := GetDockerfileStub()

	return processTemplate(output, dockerTemplate, serviceVars, outputDir)
}
