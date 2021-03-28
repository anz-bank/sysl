package importer

import (
	"github.com/sirupsen/logrus"
)

// MakePostgresqlDirImporter is a factory method for creating new PostgreSQL directory importer.
func MakePostgresqlDirImporter(logger *logrus.Logger) *ArraiImporter {
	return MakeArraiImporterImporter(ArraiImporterDir+"/postgresql/import_migrations.arraiz", logger)
}
