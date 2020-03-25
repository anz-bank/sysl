// package loader loads a source file into a sysl Module
package loader

import (
	"os"
	"path/filepath"

	"github.com/anz-bank/sysl/pkg/mod"
	"github.com/anz-bank/sysl/pkg/parse"
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/anz-bank/sysl/pkg/syslutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const syslRootMarker = ".sysl"

type ProjectConfiguration struct {
	Module, Root string
	RootIsFound  bool
	Fs           afero.Fs
}

func LoadSyslModule(root, filename string, fs afero.Fs, logger *logrus.Logger) (*sysl.Module, string, error) {
	return LoadSyslModuleWithParserType(root, filename, fs, logger, parse.DefaultParserType)
}

func LoadSyslModuleWithParserType(root, filename string, fs afero.Fs, logger *logrus.Logger, parserType parse.ParserType) (*sysl.Module, string, error) {
	logger.Debugf("Attempting to load module:%s (root:%s)", filename, root)
	projectConfig := NewProjectConfiguration()
	if err := projectConfig.ConfigureProject(root, filename, fs, logger); err != nil {
		return nil, "", err
	}

	modelParser := parse.NewParserWithParserType(parserType)
	if !projectConfig.RootIsFound {
		modelParser.RestrictToLocalImport()
	}
	return parse.LoadAndGetDefaultApp(projectConfig.Module, projectConfig.Fs, modelParser)
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

	if rootIsDefined {
		pc.RootIsFound = true
		pc.Root = root
		pc.Module = module
		if rootMarkerExists {
			logger.Warningf("%s found in %s but will use %s instead",
				syslRootMarker, syslRootPath, pc.Root)
		} else {
			logger.Warningf("%s is not defined but root flag is defined in %s",
				syslRootMarker, pc.Root)
		}
	} else {
		if rootMarkerExists {
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
		} else {
			// uses the module directory as the root, changing the module to be relative to the root
			pc.Root = filepath.Dir(module)
			pc.Module = filepath.Base(module)
			pc.RootIsFound = false
			logger.Warningf("root and %s are undefined, %s will be used instead",
				syslRootMarker, pc.Root)
		}
	}

	pc.Fs = syslutil.NewChrootFs(fs, pc.Root)
	if mod.SyslModules {
		pc.Fs = mod.NewFs(pc.Fs)
	}

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
		exists, err := afero.Exists(fs, filepath.Join(currentPath, syslRootMarker))
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
