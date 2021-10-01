package importer

import (
	"fmt"
	"strings"

	arrai2 "github.com/anz-bank/sysl/internal/arrai"
	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ProtobufImporter struct {
	appName     string
	pkg         string
	importPaths string
	logger      *logrus.Logger
}

// MakeProtobufImporter is a factory method for creating new Protobuf importer.
func MakeProtobufImporter(logger *logrus.Logger) *ProtobufImporter {
	return &ProtobufImporter{logger: logger}
}

func (p *ProtobufImporter) LoadFile(path string) (string, error) {
	b, err := arrai2.Asset("pkg/importer/proto/import_cli.arraiz")
	if err != nil {
		return "", err
	}
	// TODO: Make the appname optional
	val, err := arrai.EvaluateBundle(b, `--app-name`, p.appName, `--input`, path, `--import-paths`, p.importPaths)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("AppName: %s, File: %s, ImportPaths: %s", p.appName, path, p.importPaths),
			Err:      err,
			ShortMsg: err.Error(),
		}, "Executing arrai transform failed")
	}
	return strings.TrimSpace(val.String()) + "\n", nil
}

// Load returns a Sysl spec equivalent to protoSpec.
func (p *ProtobufImporter) Load(protoSpec string) (string, error) {
	b, err := arrai2.Asset("pkg/importer/proto/import_cli.arraiz")
	if err != nil {
		return "", err
	}
	// TODO: Make the appname optional
	val, err := arrai.EvaluateBundle(b, `--app-name`, p.appName, `--spec`, protoSpec, `--import-paths`, p.importPaths)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("AppName: %s, Content: %s, ImportPaths: %s", p.appName, protoSpec, p.importPaths),
			Err:      err,
			ShortMsg: err.Error(),
		}, "Executing arrai transform failed")
	}
	return strings.TrimSpace(val.String()) + "\n", nil
}

// Set the AppName of the imported app
func (p *ProtobufImporter) WithAppName(appName string) Importer {
	p.appName = appName
	return p
}

// Set the package attribute of the imported app
func (p *ProtobufImporter) WithPackage(pkg string) Importer {
	p.pkg = pkg
	return p
}

// Set the importPaths attribute of the imported app
func (p *ProtobufImporter) WithImports(importPaths string) Importer {
	p.importPaths = importPaths
	return p
}
