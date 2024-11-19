package importer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/anz-bank/sysl/pkg/arrai"
)

type ProtobufImporter struct {
	pkg         string
	importPaths string
	shallow     bool
	logger      *logrus.Logger
}

// MakeProtobufImporter is a factory method for creating new Protobuf importer.
func MakeProtobufImporter(logger *logrus.Logger) *ProtobufImporter {
	return &ProtobufImporter{logger: logger}
}

func (p *ProtobufImporter) LoadFile(path string) (string, error) {
	b := bundles.ProtoImporter

	args, err := buildImporterArgs(&protobufImporterArgs{
		specPath:    path,
		importPaths: p.importPaths,
		shallow:     p.shallow,
	})

	if err != nil {
		return "", err
	}

	val, err := arrai.EvaluateBundle(b, args...)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("File: %s, ImportPaths: %s", path, p.importPaths),
			Err:      err,
			ShortMsg: err.Error(),
		}, "Executing arrai transform failed")
	}
	return strings.TrimSpace(val.String()) + "\n", nil
}

// Load returns a Sysl spec equivalent to protoSpec.
func (p *ProtobufImporter) Load(protoSpec string) (string, error) {
	b := bundles.ProtoImporter

	args, err := buildImporterArgs(&protobufImporterArgs{
		specContent: protoSpec,
		importPaths: p.importPaths,
		shallow:     p.shallow,
	})

	if err != nil {
		return "", err
	}

	val, err := arrai.EvaluateBundle(b, args...)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("Content: %s, ImportPaths: %s", protoSpec, p.importPaths),
			Err:      err,
			ShortMsg: err.Error(),
		}, "Executing arrai transform failed")
	}
	return strings.TrimSpace(val.String()) + "\n", nil
}

// Configure allows the imported Sysl application name, package and import directories to be specified.
func (p *ProtobufImporter) Configure(arg *ImporterArg) (Importer, error) {
	p.pkg = arg.PackageName
	p.importPaths = arg.Imports
	p.shallow = arg.Shallow
	return p, nil
}

type protobufImporterArgs struct {
	importPaths, specPath, specContent string
	shallow                            bool
}

func buildImporterArgs(a *protobufImporterArgs) ([]string, error) {
	args := []string{}

	if a.specContent != "" && a.specPath != "" {
		return nil, errors.New("provide only path to spec or the spec content")
	}

	if a.specContent == "" && a.specPath == "" {
		return nil, errors.New("spec not provided")
	}

	if a.specContent != "" {
		args = append(args, "--spec", a.specContent)
	}
	if a.specPath != "" {
		args = append(args, "--input", a.specPath)
	}

	if a.importPaths != "" {
		args = append(args, "--import-paths", a.importPaths)
	}

	if a.shallow {
		args = append(args, "--shallow")
	}

	return args, nil
}
