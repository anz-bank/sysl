package importer

import (
	"fmt"
	"os"
	"strings"

	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/anz-bank/sysl/pkg/arrai"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// NewAvroImporter returns a new avroImporter.
func NewAvroImporter(logger *logrus.Logger) Importer {
	return &avroImporter{logger: logger}
}

// avroImporter represents Avro specification importer.
type avroImporter struct {
	appName string
	pkg     string
	logger  *logrus.Logger
}

func (i *avroImporter) LoadFile(path string) (string, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return i.Load(string(bs))
}

// Load returns a Sysl spec equivalent to avroSpec.
func (i *avroImporter) Load(avroSpec string) (string, error) {
	if i.appName == "" {
		return "", errors.New("application name not provided")
	}
	b := bundles.Transformer.Bytes()
	val, err := arrai.EvaluateBundle(b, avroSpec, i.appName, i.pkg)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("import(`%s`, `%s`, `%s`)", avroSpec, i.appName, i.pkg),
			Err:      err,
			ShortMsg: err.Error(),
		}, "Executing arrai transform failed")
	}
	return strings.TrimSpace(val.String()) + "\n", nil
}

// Configure allows the imported Sysl application name, package and import directories to be specified.
func (i *avroImporter) Configure(arg *ImporterArg) (Importer, error) {
	if arg.AppName == "" {
		return nil, errors.New("application name not provided")
	}
	i.appName = arg.AppName
	i.pkg = arg.PackageName
	return i, nil
}
