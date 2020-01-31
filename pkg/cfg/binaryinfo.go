package cfg

import (
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

//nolint:gosec
func latestGitReleaseInfo() ([]string, error) {
	tag, err := exec.Command("git", "describe", "--tags", "--abbrev=0", "--match=v[0-9]*").Output()
	if err != nil {
		return nil, err
	}
	info := make([]string, 0)
	info = append(info, string(tag[:len(tag)-1]))

	out, err := exec.Command("git", "log", "-1", "--pretty=%h %cd", "--date=short", info[0]).Output()
	if err != nil {
		return nil, err
	}
	s := strings.Split(string(out[:len(out)-1]), " ")

	return append(info, s...), nil
}
