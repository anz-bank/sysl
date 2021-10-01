package importer

import (
	"fmt"

	arrai2 "github.com/anz-bank/sysl/internal/arrai"
	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const ArraiImporterDir = "pkg/importer"

// ArraiImporter encapsulates glue code for calling arr.ai scripts to import specs.
type ArraiImporter struct {
	appName     string
	pkg         string
	asset       string
	importPaths string
	logger      *logrus.Logger
}

// MakeArraiImporterImporter returns a new ArraiImporter.
func MakeArraiImporterImporter(asset string, logger *logrus.Logger) *ArraiImporter {
	return &ArraiImporter{
		asset:  asset,
		logger: logger,
	}
}

// WithAppName allows the exported Sysl application name to be specified.
func (i *ArraiImporter) WithAppName(appName string) Importer {
	i.appName = appName
	return i
}

// WithPackage allows the exported Sysl package attribute to be specified.
func (i *ArraiImporter) WithPackage(packageName string) Importer {
	i.pkg = packageName
	return i
}

// Set the importPaths attribute of the imported app
func (i *ArraiImporter) WithImports(importPaths string) Importer {
	i.importPaths = importPaths
	return i
}

// LoadFile generates a Sysl spec be invoking the arr.ai script.
func (i *ArraiImporter) LoadFile(path string) (string, error) {
	b, err := arrai2.Asset(i.asset)
	if err != nil {
		return "", err
	}
	// TODO: Make the appname optional
	val, err := arrai.EvaluateBundle(b, `--app-name`, i.appName, `--input`, path)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("import(`%s`, `%s`)", i.appName, path),
			Err:      err,
			ShortMsg: err.Error(),
		}, "Executing arr.ai transform failed")
	}
	return val.String(), nil
}

// Load generates a Sysl spec given the content of an input file.
func (i *ArraiImporter) Load(content string) (string, error) {
	b, err := arrai2.Asset(i.asset)
	if err != nil {
		return "", err
	}
	// TODO: Make the appname optional
	val, err := arrai.EvaluateBundle(b, `--app-name`, i.appName, `--spec`, content)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("AppName: %s, Content: %s", i.appName, content),
			Err:      err,
			ShortMsg: err.Error(),
		}, "Executing arr.ai transform failed")
	}
	return val.String(), nil
}
