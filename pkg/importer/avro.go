package importer

import (
	"fmt"
	"io/ioutil"
	"strings"

	arrai2 "github.com/anz-bank/sysl/internal/arrai"
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
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return i.Load(string(bs))
}

// Load returns a Sysl spec equivalent to avroSpec.
func (i *avroImporter) Load(avroSpec string) (string, error) {
	b, err := arrai2.Asset("pkg/importer/avro/transformer_cli.arraiz")
	if err != nil {
		return "", err
	}
	val, err := arrai.EvaluateBundle(b, avroSpec, i.appName, i.pkg)
	if err != nil {
		return "", errors.Wrap(arrai.ExecutionError{
			Context:  fmt.Sprintf("import(`%s`, `%s`, `%s`)", avroSpec, i.appName, i.pkg),
			Err:      err,
			ShortMsg: "Error executing SQL dir importer",
		}, "Executing arrai transform failed")
	}
	return strings.TrimSpace(val.String()) + "\n", nil
}

func (i *avroImporter) WithAppName(appName string) Importer {
	i.appName = appName
	return i
}

// Sets the package attribute of the imported app.
func (i *avroImporter) WithPackage(pkg string) Importer {
	i.pkg = pkg
	return i
}
