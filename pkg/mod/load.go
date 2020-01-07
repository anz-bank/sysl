package mod

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func LoadAll() error {
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

	var modules goModules

	dec := json.NewDecoder(b)
	for {
		m := &goModule{}
		if err := dec.Decode(m); err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "failed to decode modules list")
		}

		modules = append(modules, m)
	}

	GoMods = modules
	return nil
}

func goGet(args ...string) error {
	if err := runGo(context.Background(), logrus.StandardLogger().Out, append([]string{"get"}, args...)...); err != nil {
		return errors.Wrapf(err, "failed to get %q", args)
	}
	return nil
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

func GetExternalFile(filename string) (string, error) {
	err := goGetByFilepath(filename)
	if err != nil {
		return filename, fmt.Errorf("%s not found\n", filename)
	}

	if err = LoadAll(); err != nil {
		return filename, fmt.Errorf("error loading modules: %s\n", err.Error())
	}

	mod := GoMods.GetByFilepath(filename)
	if mod == nil {
		return filename, fmt.Errorf("error finding module of file %s\n", filename)
	}

	relpath, err := filepath.Rel(mod.Path, filename)
	if err != nil {
		return filename, err
	}
	return filepath.Join(mod.Dir, relpath), nil
}
