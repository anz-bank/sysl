package mod

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type goModule struct {
	Path    string
	Dir     string
	Version string
}

type goModules struct{}

func (d *goModules) Get(filename, ver string, m *Modules) (mod *Module, err error) {
	if names := strings.Split(filename, "/"); len(names) > 0 {
		for i := range names[1:] {
			gogetPath := path.Join(names[:1+i]...)
			gogetPath = AppendVersion(gogetPath, ver)

			err = goGet(gogetPath)
			if err == nil {
				err = d.Load(m)
				if err != nil {
					return nil, err
				}
				mod = d.Find(filename, ver, m)
				if mod == nil {
					return nil, fmt.Errorf("error finding module of file %s", filename)
				}
				return mod, nil
			}
			logrus.Debugf("go get %s error: %s\n", gogetPath, err.Error())
		}
	}

	return nil, errors.New("no such module")
}

func (*goModules) Find(filename, ver string, m *Modules) *Module {
	for i, mod := range *m {
		if hasPathPrefix(mod.Name, filename) {
			if i == 0 && ver != "" && ver != MasterBranch {
				logrus.Warn("importing files from current folder in remote way is incorrect: use local importing instead")
			}
			if i == 0 || ver == "" || ver == MasterBranch || ver == mod.Version {
				return mod
			}
		}
	}

	return nil
}

func (*goModules) Load(m *Modules) error {
	err := goDownload()
	if err != nil {
		return err
	}

	b, err := goList()
	if err != nil {
		return err
	}

	dec := json.NewDecoder(b)
	for {
		module := &goModule{}
		if err := dec.Decode(module); err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "failed to decode modules list")
		}

		m.Add(&Module{
			Name:    module.Path,
			Dir:     module.Dir,
			Version: module.Version,
		})
	}

	return nil
}

func goGet(args ...string) error {
	if err := runGo(context.Background(), logrus.StandardLogger().Out, append([]string{"get"}, args...)...); err != nil { // nolint:lll
		return errors.Wrapf(err, "failed to get %q", args)
	}
	return nil
}

func goDownload() error {
	err := runGo(context.Background(), ioutil.Discard, "mod", "download")
	if err != nil {
		return errors.Wrap(err, "failed to download modules")
	}
	return nil
}

func goList() (io.Reader, error) {
	b := &bytes.Buffer{}
	err := runGo(context.Background(), b, "list", "-m", "-json", "all")
	if err != nil {
		return b, errors.Wrap(err, "failed to list modules")
	}
	return b, nil
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
