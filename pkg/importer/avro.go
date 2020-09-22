package importer

import (
	"strings"

	"github.com/anz-bank/sysl/pkg/arrai"
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

// Load returns a Sysl spec equivalent to avroSpec.
func (i *avroImporter) Load(avroSpec string) (string, error) {
	i.logger.Debugln("Load avro spec")

	val, err := arrai.EvaluateScript(avroTransformerScript, avroSpec, i.appName, i.pkg)
	if err != nil {
		return "", err
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
