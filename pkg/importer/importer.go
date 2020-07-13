package importer

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"
)

// Importer is an interface implemented by all sysl importers
type Importer interface {
	// Load takes in a string in a format supported by an the importer
	// It returns the converted Sysl as a string
	Load(file string) (string, error)
	// WithAppName allows the exported Sysl application name to be specified
	WithAppName(appName string) Importer
	// WithPackage allows the exported Sysl package attribute to be specified
	WithPackage(packageName string) Importer
	// WithFlags allows a set of feature flags to be passed to the importer.
	WithFlags(flags Flags) Importer
}

// Flags are used to pass in feature flags into the importer.
// They are used to allow breaking changes to be introduced without impacting existing users.
// They are intended to be short-lived and are expected to be removed when
// corresponding changes in dependent repos have been made.
type Flags struct {
	MultipleErrorResponses bool // MultipleErrorResponses is used to toggle on the generation of multiple error responses
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
		return NewOpenAPIV3Importer(logger, afero.NewOsFs()), nil
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
