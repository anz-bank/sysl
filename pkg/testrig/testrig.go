// Package testrig provides tools for generating standalone test environment for sysl generated services
package testrig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/testrig/embedded"

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
	appsWhoNeedDb := make(map[string]bool)
	for serviceName := range serviceMap {
		appsWhoNeedDb[serviceName] = false
		if applications[serviceName] != nil {
			if appNeedsDb(applications[serviceName]) {
				appsWhoNeedDb[serviceName] = true
			}
		}
	}

	composeServices := make(map[string]interface{})

	for serviceName, serviceTmpl := range serviceMap {
		if err := fs.MkdirAll(filepath.Join(outputDir, serviceName), os.ModePerm); err != nil {
			return err
		}

		if err := generateMain(fs, serviceName, serviceTmpl, outputDir, appsWhoNeedDb[serviceName]); err != nil {
			return err
		}

		if err := generateDockerfile(fs, serviceName, serviceTmpl, outputDir); err != nil {
			return err
		}

		block := generateAppService(serviceName, serviceTmpl, outputDir)
		composeServices[serviceName] = block
		if appsWhoNeedDb[serviceName] {
			block := generateDbService(serviceName)
			composeServices[serviceName+"_db"] = block
		}
	}

	return generateCompose(fs, composeServices)
}

func appNeedsDb(app *sysl.Application) bool {
	patterns := app.GetAttrs()["patterns"]
	if patterns == nil {
		return false
	}
	return strings.Contains(patterns.String(), "DB")
}

func readUserData(fs afero.Fs, templateFileName string) (embedded.ServiceMap, error) {
	templateFile, err := fs.Open(templateFileName)
	if err != nil {
		return nil, err
	}
	defer templateFile.Close()

	byteValue, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return nil, err
	}
	var vars embedded.ServiceMap
	if err := json.Unmarshal(byteValue, &vars); err != nil {
		return nil, err
	}
	return vars, nil
}

func generateAppService(serviceName string, serviceTmpl embedded.Service, outputDir string) map[string]interface{} {
	return map[string]interface{}{
		"build": map[string]interface{}{
			"context":    ".",
			"dockerfile": filepath.Join(outputDir, serviceName, "Dockerfile"),
		},
		"ports": []string{fmt.Sprintf("%v:%v", serviceTmpl.Port, serviceTmpl.Port)},
	}
}

func generateDbService(serviceName string) map[string]interface{} {
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
	if err := output.Sync(); err != nil {
		return err
	}

	return nil
}

func processTemplate(file afero.File, resource string, data embedded.Service, outputDir string) error {
	tmpl, err := template.New("servicesData").Parse(resource)
	if err != nil {
		return err
	}
	userData := struct {
		OutputDir string
		Service   embedded.Service
	}{
		OutputDir: outputDir,
		Service:   data,
	}
	if err := tmpl.Execute(file, userData); err != nil {
		return err
	}
	return file.Sync()
}

func generateMain(fs afero.Fs, serviceName string, service embedded.Service, outputDir string, needDb bool) error {
	output, err := fs.Create(filepath.Join(outputDir, serviceName, "main.go"))
	if err != nil {
		return err
	}
	defer output.Close()

	mainTemplate := embedded.GetMainStub()
	if needDb {
		mainTemplate = embedded.GetMainDbStub()
	}

	err = processTemplate(output, mainTemplate, service, outputDir)
	if err != nil {
		return err
	}

	return nil
}

func generateDockerfile(fs afero.Fs, serviceName string, service embedded.Service, outputDir string) error {
	output, err := fs.Create(filepath.Join(outputDir, serviceName, "Dockerfile"))
	if err != nil {
		return err
	}
	defer output.Close()

	dockerTemplate := embedded.GetDockerfileStub()

	err = processTemplate(output, dockerTemplate, service, outputDir)
	if err != nil {
		return err
	}

	return nil
}
