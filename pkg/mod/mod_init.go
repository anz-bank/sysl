package mod

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Assumes that a go module and a sysl module is mutually exclusive.
// This function makes the assumption that the CWD is not a go module since
// we hijack this command and use go.mod and possibly go.sum to determine
// whether the current folder/project is a sysl module
func SyslModInit(modName string, logger *logrus.Logger) error {
	out, err := exec.Command("go", "mod", "init", modName).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("go mod init failed: %s", err.Error()))
	}

	logger.Info(string(out))
	return nil
}
