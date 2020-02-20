package config

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

const CODEGEN_CONFIG_PATH = "/Users/ezhang/Projects/ANZx/codegen_config.yml"

type CodeGenConfig struct {
	Grammar   string
	Transform string
}

func ReadCodeGenDefaultParamsConfig() CodeGenConfig {
	data, _ := ioutil.ReadFile(CODEGEN_CONFIG_PATH)
	str := string(data)
	fmt.Println(str)
	var config CodeGenConfig
	yaml.Unmarshal(data, &config)
	return config
}
