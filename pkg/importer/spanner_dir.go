package importer

import (
	"github.com/sirupsen/logrus"
)

// MakeSpannerDirImporter is a factory method for creating new Spanner SQL directory importer.
func MakeSpannerDirImporter(logger *logrus.Logger) *ArraiImporter {
	return MakeArraiImporterImporter("pkg/importer/spanner/import_migrations.arraiz", logger)
}
