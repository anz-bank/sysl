// Package testrig provides tools for generating standalone test environment for sysl generated services
package testrig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/testrig/template"

	"github.com/cbroglie/mustache"
	"github.com/ghodss/yaml"
)

// GenerateRig creates full set of files required to start up a test rig
// for every service, it creates main.go (optionally with DB connection) and a dockerfile to containerize it
// if the service contains DB tables, it creates sidecar postgres container and creates a schema there
// finally, it creates docker-compose.yml at the point of all (likely your project's root) that joins all containerized services
func GenerateRig(templateFile string, outputDir string, modules []*sysl.Module) error {
	serviceMap, err := readTemplate(templateFile)
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
		err = os.MkdirAll(filepath.Join(outputDir, serviceName), os.ModePerm)
		if err != nil {
			return err
		}

		err = generateMain(serviceName, serviceTmpl, outputDir, appsWhoNeedDb[serviceName])
		if err != nil {
			return err
		}

		err = generateDockerfile(serviceName, serviceTmpl, outputDir)
		if err != nil {
			return err
		}

		block := generateAppService(serviceName, serviceTmpl, outputDir)
		composeServices[serviceName] = block
		if appsWhoNeedDb[serviceName] {
			block := generateDbService(serviceName)
			composeServices[serviceName+"_db"] = block
		}
	}

	composeFile, err := os.Create("docker-compose.yml")
	if err != nil {
		return err
	}
	defer composeFile.Close()
	return generateCompose(composeFile, composeServices)
}

func appNeedsDb(app *sysl.Application) bool {
	patterns := app.GetAttrs()["patterns"]
	if patterns == nil {
		return false
	}
	return strings.Contains(patterns.String(), "DB")
}

func readTemplate(templateFileName string) (template.ServiceMap, error) {
	templateFile, err := os.Open(templateFileName)
	if err != nil {
		return nil, err
	}
	defer templateFile.Close()

	byteValue, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return nil, err
	}
	var vars template.ServiceMap
	err = json.Unmarshal([]byte(byteValue), &vars)
	if err != nil {
		return nil, err
	}

	return vars, nil
}

func generateAppService(serviceName string, serviceTmpl template.Service, outputDir string) map[string]interface{} {
	build := make(map[string]interface{})
	build["context"] = "."
	build["dockerfile"] = filepath.Join(outputDir, serviceName, "Dockerfile")
	block := make(map[string]interface{})
	block["build"] = build
	block["ports"] = []string{fmt.Sprintf("%v:%v", serviceTmpl.Port, serviceTmpl.Port)}
	return block
}

func generateDbService(serviceName string) map[string]interface{} {
	block := make(map[string]interface{})
	block["image"] = "postgres:latest"
	block["ports"] = []string{"5432:5432"}
	// piece of magic to make Postgres execute our script on startup
	block["volumes"] = []string{fmt.Sprintf("../gen/%v/%v.sql:/docker-entrypoint-initdb.d/%v.sql", serviceName, serviceName, serviceName)}
	environment := make(map[string]string)
	environment["POSTGRES_USER"] = "someuser"
	environment["POSTGRES_PASSWORD"] = "somepassword"
	environment["POSTGRES_DB"] = "somedb"
	block["environment"] = environment
	return block
}

func generateCompose(file *os.File, composeServices map[string]interface{}) error {
	data := make(map[string]interface{})
	data["version"] = "3.3"
	data["services"] = composeServices

	dataRaw, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	_, err = file.Write(dataRaw)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func processMustacheTemplate(file *os.File, template string, data template.Service, outputDir string) error {
	// mustache needs generic map
	rawData, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	var genericData map[string]interface{}
	err = json.Unmarshal(rawData, &genericData)
	if err != nil {
		return err
	}
	genericData["outputDir"] = outputDir

	renderedTemplate, err := mustache.Render(template, genericData)
	if err != nil {
		return err
	}
	_, err = file.WriteString(renderedTemplate)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func generateMain(serviceName string, service template.Service, outputDir string, needDb bool) error {
	output, err := os.Create(filepath.Join(outputDir, serviceName, "main.go"))
	if err != nil {
		return err
	}
	defer output.Close()

	mainTemplate := template.GetMainStub()
	if needDb {
		mainTemplate = template.GetMainDbStub()
	}

	err = processMustacheTemplate(output, mainTemplate, service, outputDir)
	if err != nil {
		return err
	}

	return nil
}

func generateDockerfile(serviceName string, service template.Service, outputDir string) error {
	output, err := os.Create(filepath.Join(outputDir, serviceName, "Dockerfile"))
	if err != nil {
		return err
	}
	defer output.Close()

	dockerTemplate := template.GetDockerfileStub()

	err = processMustacheTemplate(output, dockerTemplate, service, outputDir)
	if err != nil {
		return err
	}

	return nil
}
