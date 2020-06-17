package importer

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Importer interface {
	Load(file string) (string, error)
	WithAppName(appName string) Importer
	WithPackage(packageName string) Importer
}

var Formats = []Format{
	Grammar,
	OpenAPI3,
	Swagger,
	XSD,
}

// Factory takes in an absolute filePath and a file and returns an importer from the detected file type
func Factory(filePath string, file []byte, logger *logrus.Logger) (Importer, error) {
	fileType, err := GuessFileType(filePath, file, Formats)
	if err != nil {
		return nil, err
	}
	switch fileType.Name {
	case Swagger.Name:
		logger.Debugln("Detected OpenAPI2")
		return MakeOpenAPI2Importer(logger, "", filePath), nil
	case OpenAPI3.Name:
		logger.Debugln("Detected OpenAPI3")
		return MakeOpenAPI3Importer(logger, "", filePath), nil
	case XSD.Name:
		logger.Debugln("Detected XSD")
		return MakeXSDImporter(logger), nil
	case Grammar.Name:
		logger.Debugln("Detected Grammar file")
		return nil, fmt.Errorf("importer disabled for: %s", fileType.Name)
	default:
		return nil, fmt.Errorf("an importer does not exist for: %s", fileType.Name)
	}
}
