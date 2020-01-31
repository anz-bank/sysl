package cfg

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Version   - Binary version
// GitCommit - Commit SHA of the source code
// BuildDate - Binary build date
// BuildOS   - Operating System used to build binary
//nolint:gochecknoglobals
var (
	Version   = "unspecified"
	GitCommit = "unspecified"
	BuildDate = "unspecified"
	BuildOS   = runtime.GOOS
)

//nolint:gochecknoinits
func init() {
	info, err := latestGitReleaseInfo()
	if err != nil {
		logrus.Errorf("Get latest git release info error: %s\n", err.Error())
		return
	}

	Version = info[0]
	GitCommit = info[1]
	BuildDate = info[2]
}

func latestGitReleaseInfo() ([]string, error) {
	c1 := exec.Command("git", "describe", "--match=v[0-9]*", "--tags", "--abbrev=0")
	c2 := exec.Command("xargs", "-I{}", "sh", "-c", "echo $1 $(git log -1 --pretty='%h %cd' --date=short {})", "sh", "{}")

	var err error

	c2.Stdin, err = c1.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	c2.Stdout = &b

	err = c2.Start()
	if err != nil {
		return nil, err
	}

	err = c1.Run()
	if err != nil {
		return nil, err
	}

	err = c2.Wait()
	if err != nil {
		return nil, err
	}

	return strings.Fields(b.String()), nil
}
