package setting

import (
	"encoding/json"
	"log"
	"os"
)

const SYSL_PARAMS = "SYSL_PARAMS"

func readDefaultParamsSetting() map[string]string {
	// get environment variable - SYSL_SETTING
	paramsFile := os.Getenv("/Users/ezhang/Projects/ANZx/params.json")
	file, err := os.Open(paramsFile)

	if err != nil {
		// read setting and populate map
		outmap := make(map[string]string)
		if err := json.NewDecoder(file).Decode(&outmap); err != nil {
			log.Fatal(err)
		}

		return outmap
	}

	return nil
}

func ReadCodeGenDefaultParamsSetting() map[string]string {
	return readDefaultParamsSetting()
}
