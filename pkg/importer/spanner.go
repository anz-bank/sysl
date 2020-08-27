package importer

import (
	"context"
	"fmt"

	"github.com/arr-ai/arrai/syntax"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ArraiTransformError encapsulates detailed error msgs from arrai runtime.
type ArraiTransformError struct {
	Context  string
	Err      error
	ShortMsg string
}

func (e ArraiTransformError) Error() string { return e.Context + ": " + e.Err.Error() }

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

// Load takes in a string in a format supported by an the importer.
// It returns the converted Sysl as a string.
func (s *Spanner) Load(filePath string) (string, error) {
	importParams := fmt.Sprintf("import(`%s`, `%s`, `%s`)", filePath, s.appName, s.pkg)
	syslFile, err := syntax.EvaluateExpr(context.Background(), "", fmt.Sprintf("%s.%s", importSpannerScript, importParams))
	if err != nil {
		return "", errors.Wrap(ArraiTransformError{
			Context:  importParams,
			Err:      err,
			ShortMsg: "Error executing sql importer",
		}, "Executing arrai transform failed")
	}
	return syslFile.String(), nil
}
