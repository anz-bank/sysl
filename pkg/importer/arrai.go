package importer

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ArraiImporter encapsulates glue code for calling arr.ai scripts to import specs.
type ArraiImporter struct {
	appName string
	pkg     string
	asset   []byte
	logger  *logrus.Logger
}

// MakeArraiImporterImporter returns a new ArraiImporter.
func MakeArraiImporterImporter(asset []byte, logger *logrus.Logger) *ArraiImporter {
	return &ArraiImporter{
		asset:  asset,
		logger: logger,
	}
}

// Configure allows the imported Sysl application name, package and import directories to be specified.
func (i *ArraiImporter) Configure(arg *ImporterArg) (Importer, error) {
	if arg.AppName == "" {
		return nil, errors.New("application name not provided")
	}
	i.appName = arg.AppName
	i.pkg = arg.PackageName
	return i, nil
}

// LoadFile generates a Sysl spec be invoking the arr.ai script.
func (i *ArraiImporter) LoadFile(path string) (string, error) {
	args, err := buildArraiImporterArgs(&arraiImporterArgs{
		appName:  i.appName,
		specPath: path,
		pkg:      i.pkg,
	})
	if err != nil {
		return "", err
	}

	val, err := arrai.EvaluateBundle(i.asset, args...)
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
	args, err := buildArraiImporterArgs(&arraiImporterArgs{
		appName:     i.appName,
		specContent: content,
		pkg:         i.pkg,
	})
	if err != nil {
		return "", err
	}

	val, err := arrai.EvaluateBundle(i.asset, args...)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("AppName: %s, Content: %s", i.appName, content),
			Err:      err,
			ShortMsg: err.Error(),
		}, "Executing arr.ai transform failed")
	}
	return val.String(), nil
}

type arraiImporterArgs struct {
	appName, specPath, specContent, pkg string
}

func buildArraiImporterArgs(a *arraiImporterArgs) ([]string, error) {
	args := []string{}

	if a.specContent != "" && a.specPath != "" {
		return nil, errors.New("provide only path to spec or the spec content")
	}

	if a.specContent == "" && a.specPath == "" {
		return nil, errors.New("spec not provided")
	}

	if a.appName != "" {
		args = append(args, "--app-name", a.appName)
	}

	if a.specContent != "" {
		args = append(args, "--spec", a.specContent)
	}
	if a.specPath != "" {
		args = append(args, "--input", a.specPath)
	}

	if a.pkg != "" {
		args = append(args, "--package", a.pkg)
	}

	return args, nil
}
