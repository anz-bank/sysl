package test_rig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/test_rig/template"

	"github.com/cbroglie/mustache"
	"github.com/ghodss/yaml"
)

// type TestRig struct {
// 	mainContent string
// 	dockerBlock string
// }

// type TestRigGenerator interface {
// 	generateMain(template string, vars map[string]interface{}) string
// 	generateDocker(template string, vars map[string]interface{}) string
// }

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

		block := generateAppService(serviceName, serviceTmpl)
		composeServices[serviceName] = block
		if appsWhoNeedDb[serviceName] {
			block := generateDbService(serviceName)
			composeServices[serviceName+"_db"] = block
		}
	}

	composeMap := make(map[string]interface{})
	composeMap["version"] = "3.7"
	composeMap["services"] = composeServices

	data, err := yaml.Marshal(composeMap)
	if err != nil {
		return err
	}

	composeFile, err := os.Create(filepath.Join(outputDir, "docker-compose.yml"))
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

func generateAppService(serviceName string, serviceTmpl template.Service) map[string]interface{} {
	block := make(map[string]interface{})
	block["build"] = "./" + serviceName
	block["ports"] = []string{fmt.Sprintf("%v:%v", serviceTmpl.Port, serviceTmpl.Port)}
	return block
}

func generateDbService(serviceName string) map[string]interface{} {
	block := make(map[string]interface{})
	block["image"] = "postgres:latest"
	block["ports"] = []string{"5432:5432"}
	block["volumes"] = []string{fmt.Sprintf("../gen/%v/%v.sql:/docker-entrypoint-initdb.d/%v.sql", serviceName, serviceName, serviceName)}
	environemnt := make(map[string]string)
	environemnt["POSTGRES_USER"] = "someuser"
	environemnt["POSTGRES_PASSWORD"] = "somepassword"
	environemnt["POSTGRES_DB"] = "somedb"
	block["environment"] = environemnt
	return block
}

func toGenericMap(service template.Service) (map[string]interface{}, error) {
	data, err := json.Marshal(&service)
	if err != nil {
		return nil, err
	}

	var generic map[string]interface{}
	err = json.Unmarshal(data, &generic)
	if err != nil {
		return nil, err
	}
	return generic, nil
}

func generateMain(serviceName string, service template.Service, outputDir string, needDb bool) error {
	output, err := os.Create(filepath.Join(outputDir, serviceName, "main.go"))
	defer output.Close()
	if err != nil {
		return err
	}

	genericMap, err := toGenericMap(service)
	if err != nil {
		return err
	}
	template := GetMainStub()
	if needDb {
		template = GetMainDbStub()
	}
	rendered, err := mustache.Render(template, genericMap)
	if err != nil {
		return err
	}

	n, err := output.WriteString(rendered)
	if err != nil {
		return err
	}

	log.Println("written", n, "bytes")
	output.Sync()

	return nil
}

func generateDockerfile(serviceName string, service template.Service, outputDir string) error {
	output, err := os.Create(filepath.Join(outputDir, serviceName, "Dockerfile"))
	defer output.Close()
	if err != nil {
		return err
	}

	genericMap, err := toGenericMap(service)
	if err != nil {
		return err
	}

	rendered, err := mustache.Render(GetDokerfileStub(), genericMap)
	if err != nil {
		return err
	}

	n, err := output.WriteString(rendered)
	if err != nil {
		return err
	}

	log.Println("written", n, "bytes")
	output.Sync()

	return nil
}
