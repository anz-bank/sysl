package mod

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

const syslModuleName = "syslmodules"

func SyslModInit(logger *logrus.Logger) error {
	// makes the assumption that the CWD is not a go module since we hijack this command
	out, err := exec.Command("go", "mod", "init", syslModuleName).CombinedOutput()
	if err != nil {
		return err
	}

	logger.Info(string(out))
	return nil
}
