package importer

import (
	"github.com/sirupsen/logrus"
)

// MakeSQLImporter is a factory method for creating new SQL importer.
func MakeSQLImporter(logger *logrus.Logger) *ArraiImporter {
	return MakeArraiImporterImporter(ArraiImporterDir+"/sql/import_cli.arraiz", logger)
}
