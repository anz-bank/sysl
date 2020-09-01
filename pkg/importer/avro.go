package importer

import (
	"fmt"

	"github.com/anz-bank/sysl/pkg/importer/avro"

	"github.com/arr-ai/arrai/syntax"
	"github.com/sirupsen/logrus"
)

// NewAvroImporter returns a new avroImporter
func NewAvroImporter(logger *logrus.Logger) Importer {
	return &avroImporter{
		logger: logger,
	}
}

// AvroImporter represents Avro specification importer
type avroImporter struct {
	appName string
	pkg     string
	logger  *logrus.Logger
}

// Load returns Sysl result
func (i *avroImporter) Load(avroSpec string) (string, error) {
	i.logger.Debugln("Load avro spec")

	val, err := syntax.EvaluateExpr("",
		fmt.Sprintf("%s%s", avro.AvroTransformerScript,
			fmt.Sprintf("(`%s`, `%s`, `%s`)", avroSpec, i.appName, i.pkg)))
	if err != nil {
		return "", err
	}

	return val.String(), nil
}

func (i *avroImporter) WithAppName(appName string) Importer {
	i.appName = appName
	return i
}

// Set the package attribute of the imported app
func (i *avroImporter) WithPackage(pkg string) Importer {
	i.pkg = pkg
	return i
}
