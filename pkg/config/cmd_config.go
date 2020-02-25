package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type CmdContextParamCodegen struct {
	Transform string
	Grammar   string
	DepPath   string
	BasePath  string
	AppName   string
}

func ReadCodeGenFlags(configPath string) (CmdContextParamCodegen, error) {
	var config CmdContextParamCodegen
	data, ferr := ioutil.ReadFile(configPath)
	if ferr != nil {
		return config, ferr
	}
	yerr := yaml.Unmarshal(data, &config)
	if yerr != nil {
		return config, yerr
	}

	return config, nil
}
