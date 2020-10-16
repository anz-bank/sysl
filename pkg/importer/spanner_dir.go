package importer

import (
	"fmt"

	arrai2 "github.com/anz-bank/sysl/internal/arrai"
	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// SpannerDir encapsulates glue code for calling arrai scripts which in turn ingests spanner sql.
type SpannerDir struct {
	appName string
	pkg     string
	logger  *logrus.Logger
}

// MakeSpannerImporter is a factory method for creating new spanner sql importer.
func MakeSpannerDirImporter(logger *logrus.Logger) *SpannerDir {
	return &SpannerDir{
		logger: logger,
	}
}

// WithAppName allows the exported Sysl application name to be specified.
func (s *SpannerDir) WithAppName(appName string) Importer {
	s.appName = appName
	return s
}

// WithPackage allows the exported Sysl package attribute to be specified.
func (s *SpannerDir) WithPackage(packageName string) Importer {
	s.pkg = packageName
	return s
}

// LoadFile takes a path to a directory of migration SQL scripts (.up.sql) and generates a Sysl spec
// representing the sum of statements in those scripts.
func (s *SpannerDir) LoadFile(path string) (string, error) {
	b, err := arrai2.Asset("pkg/importer/spanner/import_migrations.arraiz")
	if err != nil {
		return "", err
	}
	// TODO: Make the appname optional
	val, err := arrai.EvaluateBundle(b, path, s.appName, s.pkg)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("import(`%s`, `%s`, `%s`)", path, s.appName, s.pkg),
			Err:      err,
			ShortMsg: "Error executing SQL dir importer",
		}, "Executing arrai transform failed")
	}
	return val.String(), nil
}

// Load is not implemented, since SpannerDir importing requires a directory.
func (s *SpannerDir) Load(string) (string, error) {
	panic("not implemented")
}
