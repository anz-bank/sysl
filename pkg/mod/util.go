package mod

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func runGo(ctx context.Context, out io.Writer, args ...string) error {
	cmd := exec.CommandContext(ctx, "go", args...)

	wd, err := os.Getwd()
	if err != nil {
		return errors.Errorf("get current working directory error: %s\n", err.Error())
	}
	cmd.Dir = wd

	errbuf := new(bytes.Buffer)
	cmd.Stderr = errbuf
	cmd.Stdout = out

	logrus.Debugf("running command `go %v`\n", strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee.Err == exec.ErrNotFound {
			return nil
		}

		_, ok := err.(*exec.ExitError)
		if !ok {
			return errors.Errorf("failed to execute 'go %v': %s %T", args, err, err)
		}

		// Too old Go version
		if strings.Contains(errbuf.String(), "flag provided but not defined") {
			return nil
		}
		return errors.Errorf("go command failed: %s", errbuf)
	}

	return nil
}
