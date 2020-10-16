package importer

import (
	"fmt"

	arrai2 "github.com/anz-bank/sysl/internal/arrai"
	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Spanner encapsulates glue code for calling arrai scripts which in turn ingests spanner sql.
type Spanner struct {
	appName string
	pkg     string
	logger  *logrus.Logger
}

// MakeSpannerImporter is a factory method for creating new spanner sql importer.
func MakeSpannerImporter(logger *logrus.Logger) *Spanner {
	return &Spanner{
		logger: logger,
	}
}

// WithAppName allows the exported Sysl application name to be specified.
func (s *Spanner) WithAppName(appName string) Importer {
	s.appName = appName
	return s
}

// WithPackage allows the exported Sysl package attribute to be specified.
func (s *Spanner) WithPackage(packageName string) Importer {
	s.pkg = packageName
	return s
}

// LoadFile generates a Sysl spec equivalent to the sum of the statements in the specified SQL file.
func (s *Spanner) LoadFile(path string) (string, error) {
	b, err := arrai2.Asset("pkg/importer/spanner/import_spanner_sql.arraiz")
	if err != nil {
		return "", err
	}
	// TODO: Make the appname optional
	val, err := arrai.EvaluateBundle(b, path, s.appName, s.pkg)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("import(`%s`, `%s`, `%s`)", path, s.appName, s.pkg),
			Err:      err,
			ShortMsg: "Error executing SQL importer",
		}, "Executing arrai transform failed")
	}
	return val.String(), nil
}

func (s *Spanner) Load(content string) (string, error) {
	panic("not implemented")
}
