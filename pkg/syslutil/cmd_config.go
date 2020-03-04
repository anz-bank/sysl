package syslutil

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func ReadCMDFlags(configPath string) ([]string, error) {
	var flags []string
	data, ferr := ioutil.ReadFile(configPath)
	if ferr != nil {
		return flags, ferr
	}
	str := string(data)
	flags = strings.Fields(strings.TrimSpace(str))

	return flags, nil
}

func PopulateCMDFlagsFromFile(cmdArgs []string) ([]string, error) {
	if len(cmdArgs) < 3 {
		return nil, fmt.Errorf("command arguments are not enough")
	}

	var fileFlag string
	for _, flag := range cmdArgs[2:] {
		if strings.HasPrefix(flag, "@") {
			fileFlag = flag
			break
		}
	}

	flags, err := ReadCMDFlags(strings.Replace(fileFlag, "@", "", 1))
	if err != nil {
		return nil, err
	}

	return append(cmdArgs[0:2], flags...), err
}
