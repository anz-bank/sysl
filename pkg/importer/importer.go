package importer

import (
	"fmt"
	"path"

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

// Factory takes in a fileName and a file and returns an importer from the detected file type
func Factory(fileName string, file []byte, logger *logrus.Logger) (Importer, error) {
	fileType, err := GuessFileType(fileName, file, Formats)
	if err != nil {
		return nil, err
	}
	switch fileType.Name {
	case Swagger.Name:
		logger.Debugln("Detected OpenAPI2")
		return MakeOpenAPI2Importer(logger, "", path.Dir(fileName)), nil
	case OpenAPI3.Name:
		logger.Debugln("Detected OpenAPI3")
		return MakeOpenAPI3Importer(logger, "", path.Dir(fileName)), nil
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
