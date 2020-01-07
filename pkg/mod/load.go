package mod

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func LoadAll() error {
	return nil
}

func goGet(args ...string) error {
	if err := runGo(context.Background(), logrus.StandardLogger().Out, append([]string{"get"}, args...)...); err != nil {
		return errors.Wrapf(err, "failed to get %q", args)
	}
	return nil
}

func goGetByFilepath(filename string) error {
	return nil
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
