package mod

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type goModule struct {
	Path string
	Dir  string
}

func goGetByFilepath(filename string) error {
	dir := filepath.Dir(filename)
	re := regexp.MustCompile(`^\.|\.\.|/$`)

	for !re.MatchString(dir) {
		logrus.Debug(dir)
		err := goGet(dir)
		if err == nil {
			return nil
		}
		dir = filepath.Dir(dir)
	}

	return errors.New("No such module")
}

func goGet(args ...string) error {
	if err := runGo(context.Background(), logrus.StandardLogger().Out, append([]string{"get"}, args...)...); err != nil {
		return errors.Wrapf(err, "failed to get %q", args)
	}
	return nil
}

func (m *Modules) Load() error {
	out := ioutil.Discard

	err := runGo(context.Background(), out, "mod", "download")
	if err != nil {
		return errors.Wrap(err, "failed to download modules")
	}

	b := &bytes.Buffer{}
	err = runGo(context.Background(), b, "list", "-m", "-json", "all")
	if err != nil {
		return errors.Wrap(err, "failed to list modules")
	}

	dec := json.NewDecoder(b)
	for {
		goMod := &goModule{}
		if err := dec.Decode(goMod); err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "failed to decode modules list")
		}

		m.Add(&Module{
			Name: goMod.Path,
			Dir:  goMod.Dir,
		})
	}

	return nil
}

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
