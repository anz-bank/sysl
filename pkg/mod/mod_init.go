package mod

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Assumes that a go module and a sysl module is mutually exclusive.
// This function makes the assumption that the CWD is not a go module since
// we hijack this command and use go.mod and possibly go.sum to determine
// whether the current folder/project is a sysl module
func SyslModInit(modName string) error {
	err := runGo(context.Background(), ioutil.Discard, "mod", "init", modName)
	if err != nil {
		return errors.New(fmt.Sprintf("go mod init failed: %s", err.Error()))
	}

	return nil
}
