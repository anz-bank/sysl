package importer

import (
	"github.com/sirupsen/logrus"
)

// MakeMySQLDirImporter is a factory method for creating new PostgreSQL directory importer.
func MakeMySQLDirImporter(logger *logrus.Logger) *ArraiImporter {
	return MakeArraiImporterImporter(ArraiImporterDir+"/mysql/import_migrations.arraiz", logger)
}
