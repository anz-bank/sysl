package test_rig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/test_rig/template"

	"github.com/cbroglie/mustache"
	"github.com/ghodss/yaml"
)

func appNeedsDb(app *sysl.Application) bool {
	patterns := app.GetAttrs()["patterns"]
	if patterns == nil {
		return false
	} else {
		return strings.Contains(patterns.String(), "DB")
	}
}

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
	for serviceName, _ := range serviceMap {
		appsWhoNeedDb[serviceName] = false
		if applications[serviceName] != nil {
			if appNeedsDb(applications[serviceName]) {
				appsWhoNeedDb[serviceName] = true
			}
		}
	}

	composeServices := make(map[string]interface{})

	for serviceName, serviceTmpl := range serviceMap {
		os.MkdirAll(filepath.Join(outputDir, serviceName), os.ModePerm)

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

	composeMap := make(map[string]interface{})
	composeMap["version"] = "3.3"
	composeMap["services"] = composeServices

	data, err := yaml.Marshal(composeMap)
	if err != nil {
		return err
	}

	composeFile, err := os.Create("docker-compose.yml")
	defer composeFile.Close()
	if err != nil {
		return err
	}

	_, err = composeFile.Write(data)
	if err != nil {
		return err
	}
	composeFile.Sync()

	return nil
}

func readTemplate(templateFileName string) (template.ServiceMap, error) {
	templateFile, err := os.Open(templateFileName)
	defer templateFile.Close()
	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(templateFile)
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
	environemnt := make(map[string]string)
	environemnt["POSTGRES_USER"] = "someuser"
	environemnt["POSTGRES_PASSWORD"] = "somepassword"
	environemnt["POSTGRES_DB"] = "somedb"
	block["environment"] = environemnt
	return block
}

func processTemplate(file *os.File, template string, data template.Service, outputDir string) error {
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
	file.Sync()

	return nil
}

func generateMain(serviceName string, service template.Service, outputDir string, needDb bool) error {
	output, err := os.Create(filepath.Join(outputDir, serviceName, "main.go"))
	defer output.Close()
	if err != nil {
		return err
	}

	template := GetMainStub()
	if needDb {
		template = GetMainDbStub()
	}

	return processTemplate(output, template, service, outputDir)
}

func generateDockerfile(serviceName string, service template.Service, outputDir string) error {
	output, err := os.Create(filepath.Join(outputDir, serviceName, "Dockerfile"))
	defer output.Close()
	if err != nil {
		return err
	}

	template := GetDockerfileStub()

	return processTemplate(output, template, service, outputDir)
}
