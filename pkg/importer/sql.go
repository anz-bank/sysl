package importer

import (
	"github.com/anz-bank/golden-retriever/reader"
	"github.com/anz-bank/sysl/internal/bundles"
	"github.com/sirupsen/logrus"
)

// MakeSQLImporter is a factory method for creating new SQL importer.
func MakeSQLImporter(logger *logrus.Logger, reader reader.Reader) *ArraiImporter {
	return MakeArraiImporterImporter(bundles.SQLImporter, logger, reader)
}
