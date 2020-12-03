package importer

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"
)

// Importer is an interface implemented by all sysl importers
type Importer interface {
	// Load reads in a file from path and returns the generated Sysl.
	LoadFile(path string) (string, error)
	// Load takes in a string in a format supported by an the importer
	// It returns the converted Sysl as a string.
	Load(file string) (string, error)
	// WithAppName allows the exported Sysl application name to be specified.
	WithAppName(appName string) Importer
	// WithPackage allows the exported Sysl package attribute to be specified.
	WithPackage(packageName string) Importer
}

var Formats = []Format{
	Grammar,
	OpenAPI3,
	Swagger,
	XSD,
	Avro,
	SpannerSQL,
	SpannerSQLDir,
	Postgres,
	PostgresDir,
}

// Factory takes in an absolute path and its contents (if path is a file) and returns an importer
// for the detected file type.
func Factory(path string, isDir bool, format string, content []byte, logger *logrus.Logger) (Importer, error) {
	var fileType Format
	if format != "" {
		for _, f := range Formats {
			if format == f.Name {
				fileType = f
				break
			}
		}
		if fileType.Name == "" {
			return nil, fmt.Errorf("an importer does not exist for %s", format)
		}
	} else {
		ft, err := GuessFileType(path, isDir, content, Formats)
		if err != nil {
			return nil, err
		}
		fileType = ft
	}

	switch fileType.Name {
	case Swagger.Name:
		logger.Debugln("Detected OpenAPI2")
		return MakeOpenAPI2Importer(logger, "", path), nil
	case OpenAPI3.Name:
		logger.Debugln("Detected OpenAPI3")
		return NewOpenAPIV3Importer(logger, afero.NewOsFs()), nil
	case XSD.Name:
		logger.Debugln("Detected XSD")
		return MakeXSDImporter(logger), nil
	case Grammar.Name:
		logger.Debugln("Detected grammar file")
		return nil, fmt.Errorf("importer disabled for: %s", fileType.Name)
	case Avro.Name:
		logger.Debugln("Detected Avro")
		return NewAvroImporter(logger), nil
	case SpannerSQL.Name:
		logger.Debugln("Detected Spanner SQL file")
		return MakeSpannerImporter(logger), nil
	case SpannerSQLDir.Name:
		logger.Debugln("Detected Spanner SQL directory")
		return MakeSpannerDirImporter(logger), nil
	case Postgres.Name:
		logger.Debugln("Detected PostgreSQL file")
		return MakePostgresqlImporter(logger), nil
	case PostgresDir.Name:
		logger.Debugln("Detected PostgreSQL directory")
		return MakePostgresqlDirImporter(logger), nil
	default:
		return nil, fmt.Errorf("an importer does not exist for %s", fileType.Name)
	}
}
