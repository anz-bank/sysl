package roothandler

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	syslRootMarker = ".sysl"
)

// RootHandler handles all operations about root and argument module
type RootHandler interface {
	Root() string
	Module() string
	RootIsFound() bool
	ImportAllowed(modulePath string) (bool, error)
}

type rootStatus struct {
	root          string
	module        string
	rootIsDefined bool
}

// NewRootHandler returns a root handler
func NewRootHandler(root, module string, fs afero.Fs, logger *logrus.Logger) (RootHandler, error) {
	rh := &rootStatus{root: root, module: module, rootIsDefined: true}
	err := rh.handleRoot(fs, logger)
	if err != nil {
		return nil, err
	}
	return rh, nil
}

func (r *rootStatus) handleRoot(fs afero.Fs, logger *logrus.Logger) error {

	moduleAbsolutePath, err := filepath.Abs(r.module)
	if err != nil {
		return err
	}
	syslRootPath, err := findRootFromASyslModule(moduleAbsolutePath, fs)
	if err != nil {
		return err
	}

	rootIsDefined := r.root != ""
	rootMarkerExists := syslRootPath != ""

	// r.Root must be relative to the current directory
	currentUserDirectory, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	switch {
	case rootIsDefined:
		absoluteRootFlag, err := filepath.Abs(r.root)
		if err != nil {
			return err
		}
		relativeRootPath, err := filepath.Rel(currentUserDirectory, absoluteRootFlag)
		if err != nil {
			return err
		}
		r.root = relativeRootPath

		if rootMarkerExists {
			logger.WithFields(logrus.Fields{
				"root":      absoluteRootFlag,
				"SYSL_ROOT": syslRootPath,
			}).Warningf("root is defined even though %s exists\n", syslRootMarker)
		} else {
			warningFormat := "%s is not defined but root flag is defined in %s and will be used"
			logger.Warningf(warningFormat, syslRootMarker, absoluteRootFlag)
		}

	case !rootIsDefined && !rootMarkerExists:
		relativeModulePath, err := filepath.Rel(currentUserDirectory, moduleAbsolutePath)
		if err != nil {
			return err
		}
		relativeModulePath = filepath.Dir(relativeModulePath)
		logger.Warningf("root and %s are undefined, %s will be used instead", syslRootMarker, relativeModulePath)
		// uses the module directory as the root
		r.root = relativeModulePath
		r.rootIsDefined = false

	case !rootIsDefined && rootMarkerExists:
		relativeSyslRootPath, err := filepath.Rel(currentUserDirectory, syslRootPath)
		if err != nil {
			return err
		}
		r.root = relativeSyslRootPath
	}

	absoluteRootPath, err := filepath.Abs(r.root)
	if err != nil {
		return err
	}

	relativeModulePath, err := filepath.Rel(absoluteRootPath, moduleAbsolutePath)
	if err != nil {
		return err
	}

	r.module = relativeModulePath
	return nil
}

func findRootFromASyslModule(absolutePath string, fs afero.Fs) (string, error) {
	// Takes the closest root marker
	currentPath := absolutePath

	for {
		// Keep walking up the directories
		currentPath = filepath.Dir(currentPath)

		rootPath := filepath.Join(currentPath, syslRootMarker)

		if exists, err := afero.Exists(fs, rootPath); err != nil {
			return "", err
		} else if exists {
			break
		}
		if currentPath == "." || currentPath == "/" {
			return "", nil
		}
	}

	// returned path is always an absolute path
	return currentPath, nil
}

// Root lazily loads the root and returns an error if root is not defined
func (r *rootStatus) Root() string {
	return r.root
}

// RootIsFound shows whether a root marker or a root flag is defined
func (r *rootStatus) RootIsFound() bool {
	return r.rootIsDefined
}

func (r *rootStatus) Module() string {
	return r.module
}

func (r *rootStatus) ImportAllowed(modulePath string) (bool, error) {
	modulePath = filepath.Clean(modulePath)

	if strings.Contains(modulePath, "..") {
		return false, errors.New("import does not allow \"..\"")
	}

	return true, nil
}
