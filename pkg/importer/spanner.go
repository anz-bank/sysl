package importer

import (
	"github.com/sirupsen/logrus"
)

// MakeSpannerImporter is a factory method for creating new Spanner SQL importer.
func MakeSpannerImporter(logger *logrus.Logger) *ArraiImporter {
	return MakeArraiImporterImporter("pkg/importer/spanner/import_spanner_sql.arraiz", logger)
}
