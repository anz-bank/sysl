package importer

import (
	"github.com/sirupsen/logrus"
)

func NewAvroImporter(logger *logrus.Logger) avroimporter {
	return avroimporter{
		logger: logger,
	}
}

// avroimporter represents Avro specification importer
type avroimporter struct {
	appName string
	pkg     string
	logger  *logrus.Logger
}

// Load returns
func (i avroimporter) Load(avroSpec string) (string, error) {
	i.logger.Debugln("Load avro spec")
	return "Sysl hello", nil
}

func (i avroimporter) WithAppName(appName string) Importer {
	i.appName = appName
	return i
}

// Set the package attribute of the imported app
func (i avroimporter) WithPackage(pkg string) Importer {
	i.pkg = pkg
	return i
}
