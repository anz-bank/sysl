package importer

import (
	"github.com/sirupsen/logrus"
)

// MakePostgresqlImporter is a factory method for creating new PostgreSQL importer.
func MakePostgresqlImporter(logger *logrus.Logger) *ArraiImporter {
	return MakeArraiImporterImporter(ArraiImporterDir+"/postgresql/import_postgresql_sql.arraiz", logger)
}
