package roothandler

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path/filepath"
)

const (
	syslRootMarker = ".SYSL_ROOT"
)

type RootFinder interface {
	Root() (string, error)
}

type RootStatus struct {
	root string
	rootIsDefined, rootMarkerExists bool // makes decision on filesystem to use
}

func NewRootStatus(root string) *RootStatus{
	return &RootStatus{root: root}
}

func (r *RootStatus) RootHandler(modulePath string, fs afero.Fs, logger *logrus.Logger) error {

	moduleAbsolutePath, err := filepath.Abs(modulePath)
	if err != nil {
		return err
	}

	syslRootPath, err := r.findRootFromASyslModule(moduleAbsolutePath, fs)
	if err != nil {
		return err
	}

	r.rootIsDefined = r.root != ""
	r.rootMarkerExists = syslRootPath != ""

	switch {
	case r.rootIsDefined && r.rootMarkerExists:
		logger.WithFields(logrus.Fields{
			"root":      r.root,
			"SYSL_ROOT": syslRootPath,
		}).Warningf("root is defined even though %s exists\n", syslRootMarker)
	case r.rootIsDefined && !r.rootMarkerExists:
		logger.Warningf("%s is not defined but root flag is defined in %s and will be used", syslRootMarker, r.root)
	case !r.rootIsDefined && !r.rootMarkerExists:
		logger.Errorf("root and %s are undefined", syslRootMarker)
		r.root = ""
		return nil
	}

	// r.Root must be relative to the current directory
	currentUserDirectory, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	if !r.rootIsDefined && r.rootMarkerExists {
		relativeSyslRootPath, err := filepath.Rel(currentUserDirectory, syslRootPath)
		if err != nil {
			return err
		}
		r.root = relativeSyslRootPath
	}

	return nil
}

func (r *RootStatus) findRootFromASyslModule(absolutePath string, fs afero.Fs) (string, error) {
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
func (r *RootStatus) Root() (string, error) {
	if r.root == "" {
		return "", errors.New("root is not defined")
	}
	return r.root, nil
}
