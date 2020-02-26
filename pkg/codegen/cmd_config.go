package codegen

import (
	"io/ioutil"

	"github.com/anz-bank/sysl/pkg/cmdutils"
	"gopkg.in/yaml.v2"
)

func ReadCodeGenFlags(configPath string) (cmdutils.CmdContextParamCodegen, error) {
	var config cmdutils.CmdContextParamCodegen
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
