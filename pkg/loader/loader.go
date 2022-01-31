// package loader loads a source file into a sysl Module
package loader

import (
	"os"
	"path/filepath"

	"github.com/anz-bank/sysl/pkg/pbutil"

	"github.com/anz-bank/sysl/pkg/syslutil"

	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type ProjectConfiguration struct {
	Module, Root string
	RootIsFound  bool
	Fs           afero.Fs
}

func LoadSyslModule(root, filename string, fs afero.Fs, logger *logrus.Logger) (*sysl.Module, string, error) {
	logger.Debugf("Attempting to load module:%s (root:%s)", filename, root)
	projectConfig := NewProjectConfiguration()
	if err := projectConfig.ConfigureProject(root, filename, fs, logger); err != nil {
		return nil, "", err
	}

	modelParser := parse.NewParser()
	if !projectConfig.RootIsFound {
		modelParser.RestrictToLocalImport()
	}
	return parse.LoadAndGetDefaultApp(projectConfig.Module, projectConfig.Fs, modelParser)
}

// LoadSyslPb decodes a Sysl module from a protobuf message.
func LoadSyslModuleFromPb(pbPath string, fs afero.Fs) (*sysl.Module, error) {
	return pbutil.FromPB(pbPath, fs)
}

func NewProjectConfiguration() *ProjectConfiguration {
	return &ProjectConfiguration{
		Root:        "",
		Module:      "",
		RootIsFound: false,
		Fs:          nil,
	}
}

func (pc *ProjectConfiguration) ConfigureProject(root, module string, fs afero.Fs, logger *logrus.Logger) error {
	rootIsDefined := root != ""

	modulePath := module
	if rootIsDefined {
		modulePath = filepath.Join(root, module)
	}

	syslRootPath, err := FindRootFromSyslModule(modulePath, fs)
	if err != nil {
		return err
	}

	rootMarkerExists := syslRootPath != ""

	switch {
	case rootIsDefined:
		pc.RootIsFound = true
		pc.Root = root
		pc.Module = module
	case rootMarkerExists:
		pc.Root = syslRootPath

		// module has to be relative to the root
		absModulePath, err := filepath.Abs(module)
		if err != nil {
			return err
		}
		pc.Module, err = filepath.Rel(pc.Root, absModulePath)
		if err != nil {
			return err
		}
		pc.RootIsFound = true
	default:
		// uses the module directory as the root, changing the module to be relative to the root
		pc.Root = filepath.Dir(module)
		pc.Module = filepath.Base(module)
		pc.RootIsFound = false
	}

	logrus.Debugf("root is set to: %s\n", pc.Root)

	pc.Fs = syslutil.NewChrootFs(fs, pc.Root)

	return nil
}

func FindRootFromSyslModule(modulePath string, fs afero.Fs) (string, error) {
	currentPath, err := filepath.Abs(modulePath)
	if err != nil {
		return "", err
	}

	systemRoot, err := filepath.Abs(string(os.PathSeparator))
	if err != nil {
		return "", err
	}

	// Keep walking up the directories to find nearest root marker
	for {
		currentPath = filepath.Dir(currentPath)
		exists, err := afero.Exists(fs, filepath.Join(currentPath, parse.SyslRootMarker))
		reachedRoot := currentPath == systemRoot || (err != nil && os.IsPermission(err))
		switch {
		case exists:
			return currentPath, nil
		case reachedRoot:
			return "", nil
		case err != nil:
			return "", err
		}
	}
}
