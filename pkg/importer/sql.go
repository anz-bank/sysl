package importer

import (
	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/sirupsen/logrus"
)

// MakeSQLImporter is a factory method for creating new SQL importer.
func MakeSQLImporter(logger *logrus.Logger) *ArraiImporter {
	return MakeArraiImporterImporter(bundles.SQLImporter, logger)
}
