package importer

import (
	"github.com/sirupsen/logrus"
)

// MakeMySQLImporter is a factory method for creating new MySQL SQL importer.
func MakeMySQLImporter(logger *logrus.Logger) *ArraiImporter {
	return MakeArraiImporterImporter("pkg/importer/mysql/import_mysql_sql.arraiz", logger)
}
